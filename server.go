package hg

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"codeberg.org/reiver/go-field"
)

// ListenAndServe listens on the TCP network address `addr` and then spawns a call to the ServeMercury method on the `handler` to serve each incoming connection.
//
// For a very simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-hg"
//	)
//	
//	func main() {
//		
//		//@TODO: In your code, you would probably want to use a different handler.
//		var handler hg.Handler = hg.DebugHandler
//		
//		err := hg.ListenAndServe(":1961", handler)
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}

// Serve accepts an incoming Mercury Protocol client connection on the net.Listener `listener`.
//
// For a very simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-hg"
//
//		"net"
//	)
//	
//	func main() {
//		
//		listener, err := net.Listen("tcp", ":1961")
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//		
//		//@TODO: In your code, you would probably want to use a different handler.
//		var handler hg.Handler = hg.DebugHandler
//		
//		err = hg.Serve(listener, handler)
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func Serve(listener net.Listener, handler Handler) error {

	server := &Server{Handler: handler}
	return server.Serve(listener)
}

// Server is a Mercury Protocol server.
//
// For a simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-hg"
//	)
//	
//	func main() {
//	
//		var handler hg.Handler = hg.DebugHandler
//	
//		server := &hg.Server{
//			Addr:":1961",
//			Handler:handler,
//		}
//	
//		err := server.ListenAndServe()
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
type Server struct {
	Addr string // TCP address to listen on; if empty defaults to ":1961"
	Handler Handler // handler to invoke; if nil defaults to hg.DebugServer
	Logger Logger

	// ReadTimeout is the maximum duration for reading the request line.
	// If zero, no timeout is applied and a slow client can hold a
	// connection open indefinitely (slowloris DoS vector).
	ReadTimeout time.Duration

	// HandlerTimeout is the default maximum duration for a handler to complete.
	// If a handler implements TimeoutHandler, its Timeout value takes precedence.
	// If zero, no timeout is applied and handlers may run indefinitely.
	HandlerTimeout time.Duration

	// MaxConnections limits the number of concurrent connections the server will handle.
	// If zero or negative, no limit is applied.
	MaxConnections int

	shutdownChOnce sync.Once
	shutdownOnce   sync.Once
	shutdownCh     chan struct{}
	activeConns    sync.WaitGroup

	connsMu sync.Mutex
	conns   map[net.Conn]struct{}
}

// ListenAndServe listens on the TCP network address 'server.Addr' and then spawns a call to the ServeMercury
// method on the 'server.Handler' to serve each incoming connection.
//
// For a simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-hg"
//	)
//	
//	func main() {
//	
//		var handler hg.Handler = hg.DebugHandler
//	
//		server := &hg.Server{
//			Addr:":1961",
//			Handler:handler,
//		}
//	
//		err := server.ListenAndServe()
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func (receiver *Server) ListenAndServe() error {
	if nil == receiver {
		return ErrNilReceiver
	}

	addr := receiver.Addr
	if "" == addr {
		addr = fmt.Sprintf(":%d", DefaultTCPPort)
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	return receiver.Serve(listener)
}

// Serve accepts an incoming Mercury Protocol client connection on the net.Listener ‘listener’.
//
// For a simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-hg"
//	)
//	
//	func main() {
//	
//		listener, err := net.Listen("tcp", ":1961")
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	
//		var handler hg.Handler = hg.DebugHandler
//	
//		server := &hg.Server{
//			Handler:handler,
//		}
//	
//		err = server.Serve(listener)
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func (receiver *Server) Serve(listener net.Listener) error {
	if nil == receiver {
		return ErrNilReceiver
	}

	receiver.initShutdownChannel()

	// A Server cannot be reused after Shutdown has been called.
	select {
	case <-receiver.shutdownCh:
		return ErrServerShutdown
	default:
	}

	defer listener.Close() // Safety net for non-shutdown exits. During shutdown, listener is also closed by the shutdown goroutine to unblock Accept(); the double-close is harmless.

	log := receiver.logger().Begin()
	defer log.End()

	handler := receiver.Handler
	if nil == handler {
		handler = DebugHandler
		log.Debug(field.S("defaulted handler to DebugHandler."))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up a semaphore to limit concurrent connections, if configured.
	var sem chan struct{}
	if 0 < receiver.MaxConnections {
		sem = make(chan struct{}, receiver.MaxConnections)
	}

	// When Shutdown is called, cancel the server context and close the listener to unblock Accept.
	go func() {
		select {
		case <-receiver.shutdownCh:
			cancel()
			listener.Close()
		case <-ctx.Done():
		}
	}()

	for {
		// Wait for a new client connection.
		log.Debug(
			field.S("listening"),
			field.Stringer("addr", listener.Addr()),
		)

		// If MaxConnections is set, wait for a slot before accepting.
		if nil != sem {
			select {
			case sem <- struct{}{}:
				// acquired a slot
			case <-receiver.shutdownCh:
				log.Debug(field.S("shutting down"))
				return nil
			}
		}

		conn, err := listener.Accept()
		if nil != err {
			// Release the semaphore slot since we didn't spawn a handler.
			if nil != sem {
				<-sem
			}

			// If shutdown was requested, this is a clean exit.
			select {
			case <-receiver.shutdownCh:
				log.Debug(field.S("shutting down"))
				return nil
			default:
			}

//@TODO: Could try to recover from certain kinds of errors. Maybe waiting a while before trying again.
			log.Error(
				field.S("error while listening"),
				field.Stringer("addr", listener.Addr()),
				field.E(err),
			)
			return err
		}
		log.Debug(
			field.S("received new connection"),
			field.Stringer("remote-addr", conn.RemoteAddr()),
		)

		// Handle the new client connection by spawning
		// a new goroutine.
		receiver.trackConn(conn)
		receiver.activeConns.Add(1)
		go func() {
			defer receiver.activeConns.Done()
			defer receiver.untrackConn(conn)
			if nil != sem {
				defer func(){ <-sem }()
			}
			handle(ctx, log, conn, handler, receiver.ReadTimeout, receiver.HandlerTimeout)
		}()
		log.Debug(
			field.S("spawned handler to handle connection"),
			field.Stringer("remote-addr", conn.RemoteAddr()),
		)
	}
}

func (receiver *Server) initShutdownChannel() {
	receiver.shutdownChOnce.Do(func() {
		receiver.shutdownCh = make(chan struct{})
	})
}

// Shutdown gracefully shuts down the server.
//
// It stops accepting new connections, then waits for active connections to finish.
// The provided context controls how long Shutdown is willing to wait — if the context
// expires before all connections are done, Shutdown returns the context's error.
//
// Shutdown is safe to call multiple times — only the first call triggers the shutdown.
func (receiver *Server) Shutdown(ctx context.Context) error {
	if nil == receiver {
		return ErrNilReceiver
	}

	receiver.initShutdownChannel()

	receiver.shutdownOnce.Do(func() {
		close(receiver.shutdownCh)
	})

	done := make(chan struct{})
	go func() {
		receiver.activeConns.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		receiver.closeAllConns()
		return ctx.Err()
	}
}

// Wait blocks until all active connection handlers have finished.
//
// When Shutdown's context expires, Shutdown force-closes all connections and returns
// immediately — but the handler goroutines may still be running their cleanup (defers,
// logging, untracking). Wait blocks until those goroutines have fully exited.
//
// Typical usage:
//
//	err := server.Shutdown(ctx)
//	if err != nil {
//		// Shutdown timed out, but goroutines are still unwinding.
//		server.Wait()
//	}
//
// If Shutdown completed without error (all connections drained before the context
// expired), Wait returns immediately.
func (receiver *Server) Wait() {
	if nil == receiver {
		return
	}

	receiver.activeConns.Wait()
}

func (receiver *Server) trackConn(conn net.Conn) {
	receiver.connsMu.Lock()
	defer receiver.connsMu.Unlock()

	if nil == receiver.conns {
		receiver.conns = make(map[net.Conn]struct{})
	}
	receiver.conns[conn] = struct{}{}
}

func (receiver *Server) untrackConn(conn net.Conn) {
	receiver.connsMu.Lock()
	defer receiver.connsMu.Unlock()

	delete(receiver.conns, conn)
}

func (receiver *Server) closeAllConns() {
	receiver.connsMu.Lock()
	defer receiver.connsMu.Unlock()

	for conn := range receiver.conns {
		conn.Close()
	}
}

func (receiver *Server) logger() Logger {

	var lgr Logger
	func(){
		if nil == receiver {
			return
		}

		lgr = receiver.Logger
	}()

	return mustlogger(lgr)
}

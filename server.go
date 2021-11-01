package hg

import (
	"net"
)

// ListenAndServe listens on the TCP network address `addr` and then spawns a call to the ServeMercuary method on the `handler` to serve each incoming connection.
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

// Serve accepts an incoming Mercuary Protocol client connection on the net.Listener `listener`.
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
//		err := hg.Serve(listener, handler)
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
//		var handler hg.Handler = hg.EchoHandler
//	
//		server := &telnet.Server{
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
func (server *Server) ListenAndServe() error {

	addr := server.Addr
	if "" == addr {
		addr = ":1961"
	}

	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	return server.Serve(listener)
}

// Serve accepts an incoming Mercury Protocol client connection on the net.Listener ‘listener’.
func (server *Server) Serve(listener net.Listener) error {

	defer listener.Close()

	logger := server.logger()

	handler := server.Handler
	if nil == handler {
		handler = DebugHandler
		logger.Log("defaulted handler to DebugHandler.")
	}

	for {
		// Wait for a new client connection.
		logger.Logf("listening at %q.", listener.Addr())
		conn, err := listener.Accept()
		if err != nil {
//@TODO: Could try to recover from certain kinds of errors. Maybe waiting a while before trying again.
			logger.Errorf("error while listing at %q: %s", listener.Addr(), err)
			return err
		}
		logger.Logf("received new connection from %q.", conn.RemoteAddr())

		// Handle the new client connection by spawning
		// a new goroutine.
		go handle(logger, conn, handler)
		logger.Logf("spawned handler to handle connection from %q.", conn.RemoteAddr())
	}
}

func (server *Server) logger() Logger {

	var lgr Logger
	func(){
		if nil == server {
			return
		}

		lgr = server.Logger
	}()

	return mustlogger(lgr)
}

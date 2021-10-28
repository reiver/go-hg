package hg

import (
	"net"
)

// Server is a Mercury Protocol server.
type Server struct {
	Addr string // TCP address to listen on; if empty defaults to ":1961"
	Handler Handler // handler to invoke; if nil defaults to hg.DebugServer
	Logger Logger
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

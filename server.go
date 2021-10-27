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

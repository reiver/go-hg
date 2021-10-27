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

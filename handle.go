package hg

import (
	"net"
)

// The handle() function handles an incoming Mercury request using the handler passed to it.
func handle(logger Logger, conn net.Conn, handler Handler) {


	defer func(){
		if r := recover(); nil != r {
			if nil != logger {
				logger.Errorf("Recovered from panic(): (%T) %v", r, r)
			}
		}
	}()

	var r Request
	{
		err := r.Parse(conn)
		if nil != err {
			logger.Errorf("problem parsing request: %s", err)
			return
		}
	}
	logger.Logf("request = %q", r)

	var rw internalResponseWriter
	{
		rw.writer = conn
	}

	var w ResponseWriter = &rw

	handler.ServeMercury(w, r)
}

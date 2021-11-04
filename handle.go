package hg

import (
	"net"
)

// The handle() function handles an incoming Mercury request using the handler passed to it.
func handle(logger Logger, conn net.Conn, handler Handler) {

	logger.Trace("hg.handle: BEGIN")
	defer logger.Trace("hg.handle: END")

	defer func(logger Logger) {
		logger.Tracef("hg.handle: will close conn %T %#v", conn, conn)
		err := conn.Close()
		if nil != err {
			logger.Errorf("hg.handle: problem with closing network connection: %s", err)
		}
		logger.Tracef("hg.handle: closed connection %T %#v", conn, conn)
	}(logger)

	var request Request // This is set later; but need it here for the panic()-recover().

	defer func(logger Logger){
		if r := recover(); nil != r {
			if nil != logger {
				logger.Errorf("hg.handle: recovered from panic() from request %q: (%T) %v", request, r, r)
			}
		}
	}(logger)

	{
		err := request.Parse(conn)
		if nil != err {
			logger.Errorf("hg.handle: problem parsing request: %s", err)
			return
		}
	}
	logger.Logf("hg.handle: request = %q", request)

	var rw internalResponseWriter
	{
		rw.Writer = conn
		rw.Logger = logger
	}

	var w ResponseWriter = &rw


	if nil == handler {
		ServeTemporaryFailure(w, request)
		return
	}
	handler.ServeMercury(w, request)
}

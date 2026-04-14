package hg

import (
	"net"

	"codeberg.org/reiver/go-field"
)

// The handle() function handles an incoming Mercury request using the handler passed to it.
func handle(logger Logger, conn net.Conn, handler Handler) {

	log := logger.Begin()
	defer log.End()

	defer func(log Logger) {
		log.Trace(
			field.S("will close conn"),
			field.FormattedString("connection-type", "%T", conn),
			field.FormattedString("connection", "%#v", conn),
		)
		err := conn.Close()
		if nil != err {
			log.Error(
				field.S("problem with closing network connection"),
				field.E(err),
			)
		}
		log.Trace(
			field.S("closed connection"),
			field.FormattedString("connection-type", "%T", conn),
			field.FormattedString("connection", "%#v", conn),
		)
	}(log)

	var request Request // This is set later; but need it here for the panic()-recover().

	defer func(log Logger){
		if r := recover(); nil != r {
			if nil != log {
				log.Error(
					field.S("recovered from panic() from request"),
					field.Stringer("request", request),
					field.FormattedString("recovered-type", "%T", r),
					field.FormattedString("recovered", "%v", r),
				)
			}
		}
	}(log)

	{
		err := request.Parse(conn)
		if nil != err {
			log.Error(
				field.S("problem parsing request"),
				field.E(err),
			)
			return
		}
	}
	log.Debug(field.Stringer("request", request))

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

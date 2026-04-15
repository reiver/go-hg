package hg

import (
	"context"
	"net"
	"time"

	"codeberg.org/reiver/go-field"

	"github.com/reiver/go-hg/internal/io2"
)

// The handle() function handles an incoming Mercury request using the handler passed to it.
func handle(ctx context.Context, logger Logger, conn net.Conn, handler Handler) {

	log := logger.Begin()
	defer log.End()

	if nil == conn {
		log.Error(
			field.S("nil net.Conn"),
		)
		return
	}

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
		writer := io2.CreateWriter(conn)
		if nil == writer {
			log.Error(
				field.S("nil writer"),
			)
			return
		}

		rw.writer = writer
		rw.logger = logger
	}

	var w ResponseWriter = &rw

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if nil == handler {
		var timeout time.Duration = 1*time.Minute

		ctxWithTimeout, _ := context.WithTimeout(ctx, timeout)

		ServeTemporaryFailure(ctxWithTimeout, w, request)
		return
	}

	handler.ServeMercury(ctx, w, request)
}

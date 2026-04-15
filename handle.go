package hg

import (
	"context"
	"net"
	"time"

	"codeberg.org/reiver/go-field"

	"github.com/reiver/go-hg/internal/io2"
)

// The handle() function handles an incoming Mercury request using the handler passed to it.
func handle(ctx context.Context, logger Logger, conn net.Conn, handler Handler, readTimeout time.Duration, handlerTimeout time.Duration) {
	if nil == ctx {
		ctx = context.Background()
	}
	if nil == logger {
		logger = internalDiscardLogger{}
	}

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
	var rw internalResponseWriter // This is set later; but need it here for the panic()-recover().
	var w ResponseWriter // This is set later; but need it here for the panic()-recover().

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

			// Only send an error response if the header hasn't been written yet.
			// If it has, the response is already partially sent — closing the
			// connection is the least-bad option.
			if nil != w && !rw.headerwritten {
				ServeTemporaryFailure(ctx, w, request)
			}
		}
	}(log)

	{
		readCtx := ctx
		if readTimeout > 0 {
			var readCancel context.CancelFunc
			readCtx, readCancel = context.WithTimeout(ctx, readTimeout)
			defer readCancel()
		}

		reader := io2.ClassicReader(readCtx, io2.CreateReader(conn))

		err := request.Parse(reader)
		if nil != err {
			log.Error(
				field.S("problem parsing request"),
				field.E(err),
			)
			return
		}
	}
	log.Debug(field.Stringer("request", request))

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

	w = &rw

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if nil == handler {
		var timeout time.Duration = 1*time.Minute

		ctxWithTimeout, cancelTimeout := context.WithTimeout(ctx, timeout)
		defer cancelTimeout()

		ServeTemporaryFailure(ctxWithTimeout, w, request)
		return
	}

	// Determine the handler timeout.
	// If the handler implements TimeoutHandler, its Timeout takes precedence.
	// Otherwise, fall back to the server-wide handlerTimeout.
	{
		timeout := handlerTimeout
		if th, casted := handler.(TimeoutHandler); casted {
			timeout = th.Timeout()
		}

		if 0 < timeout {
			var handlerCancel context.CancelFunc
			ctx, handlerCancel = context.WithTimeout(ctx, timeout)
			defer handlerCancel()
		}
	}

	handler.ServeMercury(ctx, w, request)
}

package hg

import (
	"bufio"
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
			//
			// Use a fresh context because the handler's context may already be
			// cancelled (e.g., handler timeout expired). The short timeout is
			// a safety net so we don't block on a stuck write.
			if nil != w && !rw.headerwritten {
				panicCtx, panicCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer panicCancel()
				if err := ServeTemporaryFailure(panicCtx, w); nil != err {
					log.Error(
						field.S("problem sending temporary-failure response after panic"),
						field.Stringer("request", request),
						field.E(err),
					)
				}
			}
		}
	}(log)

	{
		// Set a single, absolute read deadline on the connection for the
		// entire request-line read. This protects against slowloris attacks
		// where a client sends data very slowly to hold the connection open.
		if readTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(readTimeout))
		}

		// Buffer the request line from the network in one shot.
		// The buffer is sized to maxrequest so ReadSlice
		// returns bufio.ErrBufferFull if the line exceeds the limit.
		br := bufio.NewReaderSize(conn, maxrequest)
		line, err := br.ReadSlice('\n')

		// Clear the read deadline now that we have the request bytes.
		if readTimeout > 0 {
			conn.SetReadDeadline(time.Time{})
		}

		if nil != err {
			log.Error(
				field.S("problem reading request"),
				field.E(err),
			)
			return
		}

		// Parse the in-memory bytes (UTF-8 validation, \r\n check, etc.).
		err = request.Parse(line)
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

		if err := ServeTemporaryFailure(ctxWithTimeout, w, request.RequestValue()); nil != err {
			log.Error(
				field.S("problem sending temporary-failure response for nil handler"),
				field.Stringer("request", request),
				field.E(err),
			)
		}
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

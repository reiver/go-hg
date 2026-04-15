package hg

import (
	"context"
	"io"
	"net"
	"time"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// DialAndCall makes a TCP connection to the TCP address given by 'addr',
// and (speaking the Mercury Protocol) sends the request given by 'request'.
//
// The context 'ctx' controls the lifetime of the dial and the request write.
// If 'ctx' is nil, context.Background() is used.
// To apply a timeout, use context.WithTimeout or context.WithDeadline.
//
// What is given by 'addr' might be something like: "11.22.33.44:1961", or "example.com:1961"
//
// What is given by 'request' might be a Request containing something like: "mercury://example.com/path/to/file.txt\r\n"
//
// A example of using this might be:
//
//	var uri string = "mercury://example.com/once/twice/thrice/fource.gmni"
//
//	var request hg.Request
//	err := request.Parse(uri)
//	if nil != err {
//		return err
//	}
//
//	var addr string = "example.com:1961"
//
//	ctx := context.Background()
//
//	rr, err := hg.DialAndCall(ctx, addr, request)
//
// See also:
//
//	• [Call]
func DialAndCall(ctx context.Context, addr string, request Request) (ResponseReader, error) {

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return nil, erorr.Wrap(errs, "could not dial and call for mercury protocol",
			field.String("addr", addr),
			field.Stringer("request", request),
		)
	}

	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if nil != err {
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return nil, erorr.Wrap(errs, "could not dial for mercury protocol call",
				field.String("addr", addr),
				field.Stringer("request", request),
			)
		}
		var errs error = erorr.Errors{ErrDialError, err}
		return nil, erorr.Wrap(errs, "could not dial for mercury protocol call",
			field.String("addr", addr),
			field.Stringer("request", request),
		)
	}

	rr, err := Call(ctx, conn, request)
	if nil != err {
		// The connection is not owned by a ResponseReader on the error path,
		// so nothing else will close it. Intentionally discarding the error
		// from Close().
		conn.Close()
		return nil, erorr.Wrap(err, "could not dial and call for mercury protocol",
			field.String("addr", addr),
			field.Stringer("request", request),
		)
	}
	return rr, nil
}

// Call uses the TCP connection provided by 'conn' and (speaking the Mercury Protocol) sends the request given by 'request'.
//
// The context 'ctx' controls the lifetime of the request write.
// If 'ctx' is nil, context.Background() is used.
// If 'ctx' has a deadline, it is applied as a write deadline on the connection.
//
// What is given by 'request' might be a Request containing something like: "mercury://example.com/path/to/file.txt"
//
// A example of using this might be:
//
//	conn, err := net.Dial("tcp", addr)
//	if nil != err {
//		return err
//	}
//
//	ctx := context.Background()
//
//	rr, err := hg.Call(ctx, conn, request)
//
// See also:
//
//	• [DialAndCall]
func Call(ctx context.Context, conn net.Conn, request Request) (ResponseReader, error) {
	if nil == conn {
		return nil, ErrNilNetworkConnection
	}
	if request.IsNothing() {
		return nil, ErrRequestIsNothing
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return nil, erorr.Wrap(errs, "could not make mercury protocol call",
			field.Stringer("request", request),
			field.Stringer("conn-remote-addr", conn.RemoteAddr()),
		)
	}

	if deadline, ok := ctx.Deadline(); ok {
		// Intentionally ignoring the error from SetWriteDeadline —
		// not all net.Conn implementations support deadlines, and
		// the write itself will surface any real failures.
		conn.SetWriteDeadline(deadline)
		defer conn.SetWriteDeadline(time.Time{})
	}

	_, err := io.WriteString(conn, request.String())
	if nil != err {
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return nil, erorr.Wrap(errs, "could not write mercury protocol request",
				field.Stringer("request", request),
				field.Stringer("conn-remote-addr", conn.RemoteAddr()),
			)
		}
		var errs error = erorr.Errors{ErrWriteError, err}
		return nil, erorr.Wrap(errs, "could not write mercury protocol request",
			field.Stringer("request", request),
			field.Stringer("conn-remote-addr", conn.RemoteAddr()),
		)
	}

	var rr internalResponseReader
	{
		rr.conn = conn
	}

	return &rr, nil
}

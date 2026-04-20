package hg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"time"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// DialAndCallURL makes a TCP connection to the TCP address implied by the URL.
// and (speaking the Mercury Protocol or the Gemini Protocol) sends the request implied by the URL.
//
// The context 'ctx' controls the lifetime of the dial and the request write.
// If 'ctx' is nil, context.Background() is used.
// To apply a timeout, use [context.WithTimeout] or [context.WithDeadline].
//
// The 'tlsHandler' controls server certificate verification and client certificate selection for gemini:// (TLS) URLs.
// For mercury:// URLs it is ignored.
// Passing nil means: accept any server certificate, no client cert.
//
// An example of using this for a Gemini URL might be:
//
//	var url string = "gemini://example.com/once/twice/thrice/fource.gmni"
//	
//	ctx := context.Background()
//	
//	rr, err := hg.DialAndCallURL(ctx, url, nil)
//
// Using it with a Mercury URL looks the same:
//
//	var url string = "mercury://example.com/once/twice/thrice/fource.gmni"
//	
//	ctx := context.Background()
//	
//	rr, err := hg.DialAndCallURL(ctx, url, nil)
//
// See also:
//
//	• [Call]
//	• [DialAndCall]
//	• [DialAndCallTLS]
//	• [TLSConfig]
//	• [TLSHandler]
func DialAndCallURL(ctx context.Context, url string, tlsHandler TLSHandler) (ResponseReader, error) {
	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return nil, erorr.Wrap(errs, "was told not to dial-and-call URL",
			field.String("url", url),
		)
	}

	var request Request
	{
		err := request.Parse(url)
		if nil != err {
			return nil, erorr.Wrap(err, "failed to parse URL, that would have been used to create request, that would have been used to dial-and-call URL",
				field.String("url", url),
			)
		}
	}

	var addr string
	{
		var found bool
		addr, found = request.TCPAddr()
		if !found {
			var err error = ErrBadTCPAddr

			return nil, erorr.Wrap(err, "failed to infer TCP address request created from URL, that would have been used to dial-and-call URL",
				field.String("url", url),
			)
		}
	}

	{
		var scheme string = request.Scheme()

		switch scheme {
		case Scheme:
			return DialAndCall(ctx, addr, request)
		case SchemeTLS:
			return DialAndCallTLS(ctx, addr, request, tlsHandler)
		default:
			var err error = ErrSchemeUnsupported

			return nil, erorr.Wrap(err, "failed to dial-and-call URL because of unsupported URL scheme",
				field.String("url", url),
				field.String("scheme", scheme),
			)
		}
	}
}

// DialAndCall makes a TCP connection to the TCP address given by 'addr',
// and (speaking the Mercury Protocol) sends the request given by 'request'.
//
// The context 'ctx' controls the lifetime of the dial and the request write.
// If 'ctx' is nil, context.Background() is used.
// To apply a timeout, use [context.WithTimeout] or [context.WithDeadline].
//
// What is given by 'addr' might be something like: "11.22.33.44:1961", or "example.com:1961"
//
// What is given by 'request' might be a Request containing something like: "mercury://example.com/path/to/file.txt\r\n"
//
// An example of using this might be:
//
//	var uri string = "mercury://example.com/once/twice/thrice/fource.gmni"
//	
//	var request hg.Request
//	err := request.Parse(uri)
//	if nil != err {
//		return err
//	}
//	
//	addr, found := request.TCPAddr()
//	if !found {
//		return errBadRequest
//	}
//	
//	ctx := context.Background()
//	
//	rr, err := hg.DialAndCall(ctx, addr, request)
//
// See also:
//
//	• [Call]
//	• [DialAndCallURL]
//	• [DialAndCallTLS]
func DialAndCall(ctx context.Context, addr string, request Request) (ResponseReader, error) {
	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return nil, erorr.Wrap(errs, "was told not dial-and-call for Mercury Protocol",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if nil != err {
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return nil, erorr.Wrap(errs, "was told not dial for Mercury Protocol call",
				field.String("tcp-address", addr),
				field.Stringer("request", request),
			)
		}
		var errs error = erorr.Errors{ErrDialError, err}
		return nil, erorr.Wrap(errs, "could not dial for Mercury Protocol call",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	rr, err := Call(ctx, conn, request)
	if nil != err {
		// If we got there then the call will never be able to close the net.Conn, so we have to close it ourselves.
		// Note that we are intentionally discarding the error from the net.Conn.Close() method.
		conn.Close()
		return nil, erorr.Wrap(err, "could not dial-and-call for Mercury Protocol",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}
	return rr, nil
}

// DialAndCallTLS makes a TLS connection to the TCP address given by 'addr',
// and sends the request given by 'request'.
//
// This can be used to speak protocols that layer on TLS (such as the Gemini Protocol).
//
// The context 'ctx' controls the lifetime of the dial and the request write.
// If 'ctx' is nil, context.Background() is used.
// To apply a timeout, use [context.WithTimeout] or [context.WithDeadline].
//
// The 'tlsHandler' controls server certificate verification and client certificate
// selection. Passing nil means: accept any server certificate, no client cert.
//
// What is given by 'addr' might be something like: "11.22.33.44:1965", or "example.com:1965"
//
// What is given by 'request' might be a Request containing something like: "gemini://example.com/path/to/file.txt\r\n"
//
// An example of using this might be:
//
//	var uri string = "gemini://example.com/once/twice/thrice/fource.gmni"
//	
//	var request hg.Request
//	err := request.Parse(uri)
//	if nil != err {
//		return err
//	}
//	
//	addr, found := request.TCPAddr()
//	if !found {
//		return errBadRequest
//	}
//	
//	ctx := context.Background()
//	
//	rr, err := hg.DialAndCallTLS(ctx, addr, request, nil)
//
// See also:
//
//	• [Call]
//	• [DialAndCall]
//	• [DialAndCallURL]
//	• [TLSConfig]
//	• [TLSHandler]
func DialAndCallTLS(ctx context.Context, addr string, request Request, tlsHandler TLSHandler) (ResponseReader, error) {

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return nil, erorr.Wrap(errs, "was told not dial-and-call over tls",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	hostname, _, err := net.SplitHostPort(addr)
	if nil != err {
		return nil, erorr.Wrap(err, "could not parse host-name from TCP-address",
			field.String("tcp-address", addr),
		)
	}

	var tlsConf tls.Config
	var clientCert *tls.Certificate

	tlsConf.InsecureSkipVerify = true
	tlsConf.ServerName = hostname

	if nil != tlsHandler {
		tlsConf.VerifyPeerCertificate = func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			return tlsHandler.VerifyServerCertificate(hostname, rawCerts)
		}

		clientCert = tlsHandler.ClientCertificate(hostname)
		if nil != clientCert {
			tlsConf.Certificates = []tls.Certificate{*clientCert}
		}
	}

	var dialer net.Dialer

	if deadline, ok := ctx.Deadline(); ok {
		dialer.Deadline = deadline
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if nil != err {
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return nil, erorr.Wrap(errs, "was told not dial for Gemini Protocol call",
				field.String("tcp-address", addr),
				field.Stringer("request", request),
			)
		}
		var errs error = erorr.Errors{ErrDialError, err}
		return nil, erorr.Wrap(errs, "could not dial for Gemini Protocol call",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	tlsConn := tls.Client(conn, &tlsConf)
	err = tlsConn.HandshakeContext(ctx)
	if nil != err {
		conn.Close()
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return nil, erorr.Wrap(errs, "was told not to complete TLS handshake for Gemini Protocol",
				field.String("tcp-address", addr),
				field.Stringer("request", request),
			)
		}
		var errs error = erorr.Errors{ErrDialError, err}
		return nil, erorr.Wrap(errs, "could not complete TLS handshake for Gemini Protocol",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	rr, err := call(ctx, tlsConn, request)
	if nil != err {
		// The connection is not owned by a ResponseReader on the error path,
		// so nothing else will close it. Intentionally discarding the error
		// from Close().
		tlsConn.Close()
		return nil, erorr.Wrap(err, "could not dial-and-call over TLS for Gemini Protocol",
			field.String("tcp-address", addr),
			field.Stringer("request", request),
		)
	}

	rr.clientCert = clientCert
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
// An example of using this might be:
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
//	• [DialAndCallURL]
//	• [DialAndCallTLS]
func Call(ctx context.Context, conn net.Conn, request Request) (ResponseReader, error) {
	rr, err := call(ctx, conn, request)
	if nil != err {
		return nil, err
	}
	return rr, nil
}

// call is the internal implementation of Call that returns *internalResponseReader
// directly, allowing callers (such as DialAndCallTLS) to set fields on the struct
// before returning it as a ResponseReader interface.
func call(ctx context.Context, conn net.Conn, request Request) (*internalResponseReader, error) {
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
		return nil, erorr.Wrap(errs, "was told not to make Mercury Protocol or Gemini Protocol call",
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
			return nil, erorr.Wrap(errs, "was told not to write Mercury Protocol or Gemini Protocol request",
				field.Stringer("request", request),
				field.Stringer("conn-remote-addr", conn.RemoteAddr()),
			)
		}
		var errs error = erorr.Errors{ErrWriteError, err}
		return nil, erorr.Wrap(errs, "could not write Mercury Protocol or Gemini Protocol request",
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

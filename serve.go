package hg

import (
	"context"
	"fmt"
)

// 10 INPUT
//
// This function sends a “10 INPUT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var prompt string = "Pick a number between 1 and 10"
//		
//		hg.ServeInput(ctx, w, prompt)
//		
//		// ...
//		
//	}
func ServeInput(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusInput

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaInput
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 11 SENSITIVE INPUT
//
// This function sends a “11 SENSITIVE INPUT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var prompt string = "Pick enter your password"
//		
//		hg.ServeSensitiveInput(ctx, w, prompt)
//		
//		// ...
//		
//	}
func ServeSensitiveInput(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusSensitiveInput

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaSensitiveInput
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 30 REDIRECT - TEMPORARY
//
// This function sends a “30 REDIRECT - TEMPORARY” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		//var url string = "/apple/banana/cherry.txt"
//		//var url string = "documents/info.txt"
//		var url string = "mercury://example.com/once/twice/thrice/fource.txt"
//		
//		hg.ServeRedirectTemporary(ctx, w, url)
//		
//		// ...
//		
//	}
func ServeRedirectTemporary(ctx context.Context, w ResponseWriter, target string) error {
	const statuscode = StatusRedirectTemporary

	if nil == w {
		return ErrNilResponseWriter
	}

	_, err := w.WriteHeader(ctx, statuscode, target)
	return err
}

// 31 REDIRECT - PERMANENT
//
// This function sends a “31 REDIRECT - PERMANENT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		//var url string = "/apple/banana/cherry.txt"
//		//var url string = "documents/info.txt"
//		var url string = "mercury://example.com/once/twice/thrice/fource.txt"
//		
//		hg.ServeRedirectPermanent(ctx, w, url)
//		
//		// ...
//		
//	}
func ServeRedirectPermanent(ctx context.Context, w ResponseWriter, target string) error {
	const statuscode = StatusRedirectPermanent

	if nil == w {
		return ErrNilResponseWriter
	}

	_, err := w.WriteHeader(ctx, statuscode, target)
	return err
}

// 40 TEMPORARY FAILURE
//
// This function sends a “40 TEMPORARY FAILURE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "try again later"
//		
//		hg.ServeTemporaryFailure(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeTemporaryFailure(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusTemporaryFailure

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaTemporaryFailure
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 41 SERVER UNAVAILABLE
//
// This function sends a “41 SERVER UNAVAILABLE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "we are upgrading the server"
//		
//		hg.ServeServerUnavailable(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeServerUnavailable(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusServerUnavailable

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaServerUnavailable
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 42 CGI ERROR
//
// This function sends a “42 CGI ERROR” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "the program being run just had an unexpected fatal error"
//		
//		hg.ServeCGIError(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeCGIError(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusCGIError

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaCGIError
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 43 PROXY ERROR
//
// This function sends a “43 PROXY ERROR” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "the proxy server providing TLS encryption errored out"
//		
//		hg.ServeProxyError(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeProxyError(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusProxyError

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaProxyError
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 44 SLOW DOWN
//
// This function sends a “44 SLOW DOWN” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var numberOfSecondsToWait uint = 8
//		
//		hg.ServeSlowDown(ctx, w, numberOfSecondsToWait)
//		
//		// ...
//		
//	}
func ServeSlowDown(ctx context.Context, w ResponseWriter, numberOfSecondsToWait uint) error {
	const statuscode = StatusSlowDown

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string = fmt.Sprintf("%d", numberOfSecondsToWait)

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 50 PERMANENT FAILURE
//
// This function sends a “50 PERMANENT FAILURE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "someone deleted the database"
//		
//		hg.ServePermanentFailure(ctx, w, message)
//		
//		// ...
//		
//	}
func ServePermanentFailure(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusPermanentFailure

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaPermanentFailure
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 51 NOT FOUND
//
// This function sends a “51 NOT FOUND” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "this is not the gem-page you are looking for"
//		
//		hg.ServeNotFound(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeNotFound(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusNotFound

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaNotFound
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 52 GONE
//
// This function sends a “52 GONE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "he's dead jim"
//		
//		hg.ServeGone(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeGone(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusGone

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaGone
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 53 PROXY REQUEST REFUSED
//
// This function sends a “53 PROXY REQUEST REFUSED” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "you did not enter a number"
//		
//		hg.ServeProxyRequestRefused(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeProxyRequestRefused(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusProxyRequestRefused

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaProxyRequestRefused
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

// 59 BAD REQUEST
//
// This function sends a “59 BAD REQUEST” Mercury Protocol response.
//
// Example Usage
//
// This is how one might use this helper-function:
//
//	func ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "you did not enter a number"
//		
//		hg.ServeBadRequest(ctx, w, message)
//		
//		// ...
//		
//	}
func ServeBadRequest(ctx context.Context, w ResponseWriter, a ...any) error {
	const statuscode = StatusBadRequest

	if nil == w {
		return ErrNilResponseWriter
	}

	var meta string
	{
		switch {
		case 0 < len(a):
			meta = fmt.Sprint(a...)
		default:
			meta = DefaultMetaBadRequest
		}
	}

	_, err := w.WriteHeader(ctx, statuscode, meta)
	return err
}

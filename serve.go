package hg

import (
	"fmt"
)

// 10 INPUT
//
// This function sends a “10 INPUT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var prompt string = "Pick a number between 1 and 10"
//		
//		hg.ServeInput(w, prompt)
//		
//		// ...
//		
//	}
func ServeInput(w ResponseWriter, meta string) {
	const statuscode = StatusInput

	if nil == w {
		return
	}

	w.WriteHeader(statuscode, meta)
}

// 11 SENSITIVE INPUT
//
// This function sends a “11 SENSITIVE INPUT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var prompt string = "Pick enter your password"
//		
//		hg.ServeSensitiveInput(w, prompt)
//		
//		// ...
//		
//	}
func ServeSensitiveInput(w ResponseWriter, meta string) {
	const statuscode = StatusSensitiveInput

	if nil == w {
		return
	}

	w.WriteHeader(statuscode, meta)
}

// 30 REDIRECT - TEMPORARY
//
// This function sends a “30 REDIRECT - TEMPORARY” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		//var url string = "/apple/banana/cherry.txt"
//		//var url string = "documents/info.txt"
//		var url string = "mercury://example.com/once/twice/thrice/fource.txt"
//		
//		hg.ServeRedirectTemporary(w, url)
//		
//		// ...
//		
//	}
func ServeRedirectTemporary(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectTemporary

	if nil == w {
		return
	}

	w.WriteHeader(statuscode, meta)
}

// 31 REDIRECT - PERMANENT
//
// This function sends a “31 REDIRECT - PERMANENT” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		//var url string = "/apple/banana/cherry.txt"
//		//var url string = "documents/info.txt"
//		var url string = "mercury://example.com/once/twice/thrice/fource.txt"
//		
//		hg.ServeRedirectPermanent(w, url)
//		
//		// ...
//		
//	}
func ServeRedirectPermanent(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectPermanent

	if nil == w {
		return
	}

	w.WriteHeader(statuscode, meta)
}

// 40 TEMPORARY FAILURE
//
// This function sends a “40 TEMPORARY FAILURE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "try again later"
//		
//		hg.ServeTemporaryFailure(w, message)
//		
//		// ...
//		
//	}
func ServeTemporaryFailure(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusTemporaryFailure

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 41 SERVER UNAVAILABLE
//
// This function sends a “41 SERVER UNAVAILABLE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "we are upgrading the server"
//		
//		hg.ServeServerUnavailable(w, message)
//		
//		// ...
//		
//	}
func ServeServerUnavailable(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusServerUnavailable

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 42 CGI ERROR
//
// This function sends a “42 CGI ERROR” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "the program being run just had an unexpected fatal error"
//		
//		hg.ServeCGIError(w, message)
//		
//		// ...
//		
//	}
func ServeCGIError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusCGIError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 43 PROXY ERROR
//
// This function sends a “43 PROXY ERROR” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "the proxy server providing TLS encryption errored out"
//		
//		hg.ServeProxyError(w, message)
//		
//		// ...
//		
//	}
func ServeProxyError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 44 SLOW DOWN
//
// This function sends a “44 SLOW DOWN” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var numberOfSecondsToWait uint = 8
//		
//		hg.ServeSlowDown(w, numberOfSecondsToWait)
//		
//		// ...
//		
//	}
func ServeSlowDown(w ResponseWriter, numberOfSecondsToWait uint) {
	const statuscode = StatusSlowDown

	if nil == w {
		return
	}

	var meta string = fmt.Sprintf("%d", numberOfSecondsToWait)

	w.WriteHeader(statuscode, meta)
}

// 50 PERMANENT FAILURE
//
// This function sends a “50 PERMANENT FAILURE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "someone deleted the database"
//		
//		hg.ServePermanentFailure(w, message)
//		
//		// ...
//		
//	}
func ServePermanentFailure(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusPermanentFailure

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 51 NOT FOUND
//
// This function sends a “51 NOT FOUND” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "this is not the gem-page you are looking for"
//		
//		hg.ServeNotFound(w, message)
//		
//		// ...
//		
//	}
func ServeNotFound(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusNotFound

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 52 GONE
//
// This function sends a “52 GONE” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "he's dead jim"
//		
//		hg.ServeGone(w, message)
//		
//		// ...
//		
//	}
func ServeGone(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusGone

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 53 PROXY REQUEST REFUSED
//
// This function sends a “53 PROXY REQUEST REFUSED” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "you did not enter a number"
//		
//		hg.ServeProxyRequestRefused(w, message)
//		
//		// ...
//		
//	}
func ServeProxyRequestRefused(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyRequestRefused

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 59 BAD REQUEST
//
// This function sends a “59 BAD REQUEST” Mercury Protocol response.
//
// Example Usage
//
// This is how one might is this helper-function:
//
//	func ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		
//		// ...
//		
//		var message string = "you did not enter a number"
//		
//		hg.ServeBadRequest(w, message)
//		
//		// ...
//		
//	}
func ServeBadRequest(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusBadRequest

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

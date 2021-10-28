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

// 40
func ServeTemporaryFailure(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusTemporaryFailure

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 41
func ServeServerUnavailable(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusServerUnavailable

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 42
func ServeCGIError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusCGIError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 43
func ServeProxyError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 44
func ServeSlowDown(w ResponseWriter, meta string) {
	const statuscode = StatusSlowDown

	if nil == w {
		return
	}

	w.WriteHeader(statuscode, meta)
}

// 50
func ServePermanentFailure(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusPermanentFailure

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 51
func ServeNotFound(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusNotFound

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 52
func ServeGone(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusGone

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 53
func ServeProxyRequestRefused(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyRequestRefused

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

// 59
func ServeBadRequest(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusBadRequest

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	w.WriteHeader(statuscode, meta)
}

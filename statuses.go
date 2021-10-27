package hg

import (
	"fmt"
	"io"
)

const (
	StatusInput              = 10
	StatusSensitiveInput     = 11

	StatusSuccess            = 20

	StatusRedirectTemporary  = 30
	StatusRedirectPermanent  = 31

	StatusTemporaryFailure   = 40
	StatusServerUnavailable  = 41
	StatusCGIError           = 42
	StatusProxyError         = 43
	StatusSlowDown           = 44

	StatusPermanentFailture   = 50
	StatusNotFound            = 51
	StatusGone                = 52
	StatusProxyRequestRefused = 53
	StatusBadRequest          = 59
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
//		hg.Input(w, prompt)
//		
//		// ...
//		
//	}
func Input(w ResponseWriter, meta string) {
	const statuscode = StatusInput

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
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
//		hg.SensitiveInput(w, prompt)
//		
//		// ...
//		
//	}
func SensitiveInput(w ResponseWriter, meta string) {
	const statuscode = StatusSensitiveInput

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
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
//		hg.RedirectTemporary(w, url)
//		
//		// ...
//		
//	}
func RedirectTemporary(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectTemporary

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
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
//		hg.RedirectPermanent(w, url)
//		
//		// ...
//		
//	}
func RedirectPermanent(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectPermanent

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 40
func TemporaryFailure(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusTemporaryFailure

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 41
func ServerUnavailable(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusServerUnavailable

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 42
func CGIError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusCGIError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 43
func ProxyError(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyError

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 44
func SlowDown(w ResponseWriter, meta string) {
	const statuscode = StatusSlowDown

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 50
func PermanentFailture(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusPermanentFailture

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 51
func NotFound(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusNotFound

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 52
func Gone(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusGone

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 53
func ProxyRequestRefused(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusProxyRequestRefused

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 59
func BadRequest(w ResponseWriter, a ...interface{}) {
	const statuscode = StatusBadRequest

	if nil == w {
		return
	}

	var meta string = fmt.Sprint(a...)

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

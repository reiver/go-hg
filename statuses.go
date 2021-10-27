package hg

import (
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

	StatusPermanentFailture  = 50
	StatusNotFound           = 51
	StatusGone               = 52
	StatusProxyRequestFailed = 53
	StatusBadRequest         = 59
)

// 10
func Input(w ResponseWriter, meta string) {
	const statuscode = StatusInput

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 11
func SensitiveInput(w ResponseWriter, meta string) {
	const statuscode = StatusSensitiveInput

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 30
func RedirectTemporary(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectTemporary

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 31
func RedirectPermanent(w ResponseWriter, meta string) {
	const statuscode = StatusRedirectPermanent

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 40
func TemporaryFailure(w ResponseWriter, meta string) {
	const statuscode = StatusTemporaryFailure

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 41
func ServerUnavailable(w ResponseWriter, meta string) {
	const statuscode = StatusServerUnavailable

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 42
func CGIError(w ResponseWriter, meta string) {
	const statuscode = StatusCGIError

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 43
func ProxyError(w ResponseWriter, meta string) {
	const statuscode = StatusProxyError

	if nil == w {
		return
	}

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
func PermanentFailture(w ResponseWriter, meta string) {
	const statuscode = StatusPermanentFailture

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 51
func NotFound(w ResponseWriter, meta string) {
	const statuscode = StatusNotFound

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 52
func Gone(w ResponseWriter, meta string) {
	const statuscode = StatusGone

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 53
func ProxyRequestFailed(w ResponseWriter, meta string) {
	const statuscode = StatusProxyRequestFailed

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

// 59
func BadRequest(w ResponseWriter, meta string) {
	const statuscode = StatusBadRequest

	if nil == w {
		return
	}

	io.WriteString(w.Meta(), meta)
	w.WriteHeader(statuscode)
}

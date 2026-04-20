package hg

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrBadStatusCode             = erorr.Error("bad status code")
	ErrBadResponseHeaderMeta     = erorr.Error("bad response header meta")
	ErrBadTCPAddr                = erorr.Error("bad TCP address")
	ErrCannotParse               = erorr.Error("cannot parse")
	ErrContextDone               = erorr.Error("context done")
	ErrDialError                 = erorr.Error("dial error")
	ErrNilNetworkConnection      = erorr.Error("nil network connection")
	ErrNilReceiver               = erorr.Error("nil receiver")
	ErrNilResponseReader         = erorr.Error("nil response reader")
	ErrNilResponseWriter         = erorr.Error("nil response writer")
	ErrRequestIsNothing          = erorr.Error("request is nothing")
	ErrResponseHeaderMetaTooBig  = erorr.Error("response header meta too big")
	ErrSchemeUnsupported         = erorr.Error("scheme unsupported")
	ErrServerCertificateNotFound = erorr.Error("server did not present any certificate(s)")
	ErrServerShutdown            = erorr.Error("server shutdown")
	ErrTargetTypeUnsupported     = erorr.Error("target type unsupported")
	ErrWriteError                = erorr.Error("write error")
)

const (
	errBadResponse          = erorr.Error("hg: bad response")
	errExpectedLineFeed     = erorr.Error("hg: expected line-feed character")
	errHeaderAlreadyRead    = erorr.Error("hg: header already read")
	errHeaderAlreadyWritten = erorr.Error("hg: header already written")
	errNilReader            = erorr.Error("hg: nil reader")
	errNilWriter            = erorr.Error("hg: nil io.Writer")
	errNilSource            = erorr.Error("hg: nil source")
	errNothing              = erorr.Error("hg: nothing")
	errRequestTooLong       = erorr.Error("hg: request too long")
	errRuneError            = erorr.Error("hg: rune error")
)

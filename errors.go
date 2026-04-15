package hg

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrCannotParse              = erorr.Error("cannot parse")
	ErrContextDone              = erorr.Error("context done")
	ErrDialError                = erorr.Error("dial error")
	ErrNilNetworkConnection     = erorr.Error("nil network connection")
	ErrRequestIsNothing         = erorr.Error("request is nothing")
	ErrResponseHeaderMetaTooBig = erorr.Error("response header meta too big")
	ErrTargetTypeUnsupported    = erorr.Error("target type unsupported")
	ErrWriteError               = erorr.Error("write error")
)

const (
	errBadResponse             = erorr.Error("hg: bad response")
	errBadStatusCode           = erorr.Error("hg: bad status code")
	errExpectedLineFeed        = erorr.Error("hg: expected line-feed character")
	errHeaderAlreadyRead       = erorr.Error("hg: header already read")
	errHeaderAlreadyWritten    = erorr.Error("hg: header already written")
	errNilReader               = erorr.Error("hg: nil reader")
	errNilReceiver             = erorr.Error("hg: nil receiver")
	errNilWriter               = erorr.Error("hg: nil io.Writer")
	errNilSource               = erorr.Error("hg: nil source")
	errNothing                 = erorr.Error("hg: nothing")
	errRuneError               = erorr.Error("hg: rune error")
)

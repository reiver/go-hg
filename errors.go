package hg

import (
	"errors"
)

var (
	errBadResponse             = errors.New("hg: bad response")
	errBadStatusCode           = errors.New("hg: bad status code")
	errExpectedLineFeed        = errors.New("hg: expected line-feed character")
	errHeaderAlreadyRead       = errors.New("hg: header already read")
	errHeaderAlreadyWritten    = errors.New("hg: header already written")
	errNilNetworkConnection    = errors.New("hg: nil network connection")
	errNilReader               = errors.New("hg: nil reader")
	errNilReceiver             = errors.New("hg: nil receiver")
	errNilWriter               = errors.New("hg: nil io.Writer")
	errNilSource               = errors.New("hg: nil source")
	errNothing                 = errors.New("hg: nothing")
	errRuneError               = errors.New("hg: rune error")
)

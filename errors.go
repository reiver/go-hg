package hg

import (
	"errors"
)

var (
	errExpectedLineFeed = errors.New("hg: expected line-feed character")
	errNilReceiver = errors.New("hg: nil receiver")
	errNilWriter   = errors.New("hg: nil io.Writer")
	errNilSource   = errors.New("hg: nil source")
	errNothing     = errors.New("hg: nothing")
	errRuneError   = errors.New("hg: rune error")
)

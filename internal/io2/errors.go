package io2

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrNilDeadLinedReader = erorr.Error("nil dead-lined reader")
	ErrNilDeadLinedWriter = erorr.Error("nil dead-lined writer")
	ErrNilReader          = erorr.Error("nil reader")
	ErrNilReceiver        = erorr.Error("nil receiver")
	ErrNilWriter          = erorr.Error("nil writer")
)

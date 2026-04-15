package io2

import (
	"codeberg.org/reiver/go-erorr"
)

const (
	ErrNilDeadLinedWriter = erorr.Error("nil dead-lined writer")
	ErrNilReceiver        = erorr.Error("nil receiver")
	ErrNilWriter          = erorr.Error("nil writer")
)

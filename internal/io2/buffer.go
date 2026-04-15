package io2

import (
	"context"
	"strings"

	"codeberg.org/reiver/go-erorr"
)

type Buffer struct {
	builder strings.Builder
}

var _ Writer = &Buffer{}

func (receiver *Buffer) String() string {
	if nil == receiver {
		return ""
	}

	return receiver.builder.String()
}

func (receiver *Buffer) Write(ctx context.Context, p []byte) (n int, err error) {
	if nil == receiver {
		var nada int
		return nada, ErrNilReceiver
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		var nada int
		return nada, erorr.Wrap(err, "failed to write due to invalid context")
	}

	return receiver.builder.Write(p)
}

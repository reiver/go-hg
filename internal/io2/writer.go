package io2

import (
	"context"
	"io"
	"time"

	"codeberg.org/reiver/go-erorr"
)

// DeadLinedWriter represents an io.Writer with a (particular type of) SetWriteDeadline method.
//
// One example of this is net.Conn.
type DeadLinedWriter interface {
	io.Writer
	SetWriteDeadline(t time.Time) error
}

// Writer represents a modernized version of [io.Writer], where the Write method has a [context.Context] as its first parameter.
type Writer interface {
	Write(ctx context.Context, p []byte) (n int, err error)
}

type internalWriter struct {
	deadLinedWriter DeadLinedWriter
}

var _ Writer = internalWriter{}

// CreateWriter returns a [Writer] based on a [DeadLinedWriter] (such as [net.Conn]).
func CreateWriter(dw DeadLinedWriter) Writer {
	if nil == dw {
		return nil
	}

	return internalWriter{
		deadLinedWriter: dw,
	}
}

func (receiver internalWriter) Write(ctx context.Context, p []byte) (n int, err error) {
	if nil == receiver.deadLinedWriter {
		var nada int
		return nada, ErrNilDeadLinedWriter
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		var nada int
		return nada, erorr.Wrap(err, "failed to write due to invalid context")
	}

	if deadline, ok := ctx.Deadline(); ok {
		// Intentionally mostly ignoring the error from SetWriteDeadline —
		// not all net.Conn implementations support deadlines, and
		// the write itself will surface any real failures.
		if nil == receiver.deadLinedWriter.SetWriteDeadline(deadline) {
			defer receiver.deadLinedWriter.SetWriteDeadline(time.Time{})
		}
	}

	if err := ctx.Err(); err != nil {
		var nada int
		return nada, erorr.Wrap(err, "failed to write due to invalid context")
	}
	return receiver.deadLinedWriter.Write(p)
}

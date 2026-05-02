package io2

import (
	"context"
	"io"
	"time"

	"codeberg.org/reiver/go-erorr"
)

// DeadLinedReader represents an io.Reader with a (particular type of) SetReadDeadline method.
//
// One example of this is net.Conn.
type DeadLinedReader interface {
	io.Reader
	SetReadDeadline(t time.Time) error
}

// Reader represents a modernized version of [io.Reader], where the Read method has a [context.Context] as its first parameter.
type Reader interface {
	Read(ctx context.Context, p []byte) (n int, err error)
}

type internalReader struct {
	deadLinedReader DeadLinedReader
}

var _ Reader = internalReader{}

// CreateReader returns a [Reader] based on a [DeadLinedReader] (such as [net.Conn]).
func CreateReader(dr DeadLinedReader) Reader {
	if nil == dr {
		return nil
	}

	return internalReader{
		deadLinedReader: dr,
	}
}

func (receiver internalReader) Read(ctx context.Context, p []byte) (n int, err error) {
	if nil == receiver.deadLinedReader {
		var nada int
		return nada, ErrNilDeadLinedReader
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		var nada int
		return nada, erorr.Wrap(err, "failed to read due to invalid context")
	}

	if deadline, ok := ctx.Deadline(); ok {
		// Intentionally mostly ignoring the error from SetReadDeadline —
		// not all net.Conn implementations support deadlines, and
		// the read itself will surface any real failures.
		//
		// Note: the defer clears the deadline to zero after the read completes.
		// This means any deadline previously set on the underlying conn by an
		// outer layer will be erased. This is acceptable because io2.Reader
		// owns the deadline for the duration of each Read call, and net.Conn
		// provides no way to read back the current deadline to restore it.
		if nil == receiver.deadLinedReader.SetReadDeadline(deadline) {
			defer receiver.deadLinedReader.SetReadDeadline(time.Time{})
		}
	}

	if err := ctx.Err(); err != nil {
		var nada int
		return nada, erorr.Wrap(err, "failed to read due to invalid context")
	}
	return receiver.deadLinedReader.Read(p)
}

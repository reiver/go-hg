package hg

import (
	"context"
	"encoding"
	"io"
	"net"
	"strings"
	"time"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-utf8s"
)

// ResponseReader is used by a Handler to read a Mercury Protocol response.
type ResponseReader interface {
	io.Closer
	Read(ctx context.Context, data []byte) (int, error)
	Reader(ctx context.Context) io.Reader
	ReadHeader(ctx context.Context, statusCode *int, meta any) (int, error)
}

var _ ResponseReader = &internalResponseReader{}

// internalResponseReader is used to create a ResponseReader around a net.Conn.
type internalResponseReader struct {
	conn net.Conn
	headerread bool
}

func (receiver *internalResponseReader) Close() error {
	if nil == receiver {
		return errNilReceiver
	}

	var conn net.Conn = receiver.conn
	{
		if nil == conn {
			return nil
		}
	}

	return conn.Close()
}

func (receiver *internalResponseReader) Read(ctx context.Context, data []byte) (n int, err error) {
	if nil == receiver {
		return 0, errNilReceiver
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return 0, erorr.Wrap(errs, "could not read mercury protocol response")
	}

	var conn net.Conn = receiver.conn
	if nil == conn {
		return 0, errNilReader
	}

	// Set up deadline and cancellation goroutine.
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetReadDeadline(deadline)
	}
	if nil != ctx.Done() {
		done := make(chan struct{})
		exited := make(chan struct{})
		defer func() {
			close(done)
			<-exited
			conn.SetReadDeadline(time.Time{})
		}()
		go func() {
			defer close(exited)
			select {
			case <-ctx.Done():
				conn.SetReadDeadline(time.Now())
			case <-done:
			}
		}()
	} else if _, ok := ctx.Deadline(); ok {
		// Deadline but no cancel — still need to clear after.
		defer conn.SetReadDeadline(time.Time{})
	}

	// Auto-read header if not yet read.
	// Uses internal readHeader (no context management) to stay
	// within this method's deadline/goroutine scope.
	if !receiver.headerread {
		var statusCode int
		var meta string

		_, err := receiver.readHeader(&statusCode, &meta)
		if nil != err {
			if ctxErr := ctx.Err(); nil != ctxErr {
				var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
				return 0, erorr.Wrap(errs, "could not read mercury protocol response header")
			}
			return 0, err
		}

		if StatusSuccess != statusCode {
			return 0, ErrorResponse(statusCode, meta)
		}
	}

	// Read body bytes.
	{
		m, err := conn.Read(data)
		n += m
		if nil != err {
			if ctxErr := ctx.Err(); nil != ctxErr {
				var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
				return n, erorr.Wrap(errs, "could not read mercury protocol response body")
			}
			return n, err
		}
	}

	return n, nil
}

// readHeader does pure header parsing with no context/deadline management.
// Both the public ReadHeader and the auto-read path in Read call this method.
func (receiver *internalResponseReader) readHeader(statusCode *int, meta any) (n int, err error) {

	// The Gemini Protocol spec (which the Mercury Protocol is based on) says the meta SHOULD be a maximum of 1024 bytes.
	// We are allowing more than that.
	const maxmeta = 1024 * 2

	if nil == receiver {
		return 0, errNilReceiver
	}

	var reader io.Reader = receiver.conn
	if nil == reader {
		return 0, errNilReader
	}

	if receiver.headerread {
		return 0, errHeaderAlreadyRead
	}

	var mostSignificant rune
	{
		r, size, err := utf8s.ReadRune(reader)
		n += size
		if io.EOF == err {
			return n, errBadResponse
		}
		if nil != err {
			return n, err
		}
		if utf8s.RuneError == r {
			return n, errRuneError
		}

		mostSignificant = r
	}

	var leastSignificant rune
	{
		r, size, err := utf8s.ReadRune(reader)
		n += size
		if io.EOF == err {
			return n, errBadResponse
		}
		if nil != err {
			return n, err
		}
		if utf8s.RuneError == r {
			return n, errRuneError
		}

		leastSignificant = r
	}

	{
		result, valid := statuscodeFromRunes(mostSignificant, leastSignificant)
		if !valid {
			return n, errBadResponse
		}

		if nil != statusCode {
			*statusCode = result
		}
	}

	{
		r, size, err := utf8s.ReadRune(reader)
		n += size
		if io.EOF == err {
			return n, errBadResponse
		}
		if nil != err {
			return n, err
		}
		if utf8s.RuneError == r {
			return n, errRuneError
		}

		if ' ' != r {
			return n, errBadResponse
		}
	}

	{
		var storage strings.Builder

		for {
			r, size, err := utf8s.ReadRune(reader)
			n += size
			if io.EOF == err {
				return n, errBadResponse
			}
			if nil != err {
				return n, err
			}
			if utf8s.RuneError == r {
				return n, errRuneError
			}

			if '\r' == r {
		/////////////// BREAK
				break
			}
			if maxmeta < storage.Len() {
				var err error = ErrResponseHeaderMetaTooBig
				return n, erorr.Wrap(err, "response header meta too big",
					field.Uint64("max", maxmeta),
				)
			}

			storage.WriteRune(r)
		}

		if nil != meta {

			var value string = storage.String()

			switch casted := meta.(type) {
			case encoding.TextUnmarshaler:
				var p []byte = []byte(value)

				err := casted.UnmarshalText(p)
				if nil != err {
					return n, err
				}
			case *string:
				*casted = value
			case *[]byte:
				var p []byte = []byte(value)

				*casted = p
			case *[]rune:
				var p []rune = []rune(value)

				*casted = p
			default:
				var err error = ErrTargetTypeUnsupported
				return n, erorr.Wrap(err, "cannot load value for meta into variable of provided type",
					field.FormattedString("meta-type", "%T", meta),
					field.Any("meta", meta),
				)
			}
		}
	}

	{
		r, size, err := utf8s.ReadRune(reader)
		n += size
		if io.EOF == err {
			return n, errBadResponse
		}
		if nil != err {
			return n, err
		}
		if utf8s.RuneError == r {
			return n, errRuneError
		}

		if '\n' != r {
			return n, errBadResponse
		}
	}

	receiver.headerread = true
	return n, nil
}

// ReadHeader is the public context-aware wrapper.
func (receiver *internalResponseReader) ReadHeader(ctx context.Context, statusCode *int, meta any) (n int, err error) {
	if nil == receiver {
		return 0, errNilReceiver
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var errs error = erorr.Errors{ErrContextDone, ctxErr}
		return 0, erorr.Wrap(errs, "could not read mercury protocol response header")
	}

	var conn net.Conn = receiver.conn
	if nil == conn {
		return 0, errNilReader
	}

	// Set up deadline and cancellation goroutine (same pattern as Read).
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetReadDeadline(deadline)
	}
	if nil != ctx.Done() {
		done := make(chan struct{})
		exited := make(chan struct{})
		defer func() {
			close(done)
			<-exited
			conn.SetReadDeadline(time.Time{})
		}()
		go func() {
			defer close(exited)
			select {
			case <-ctx.Done():
				conn.SetReadDeadline(time.Now())
			case <-done:
			}
		}()
	} else if _, ok := ctx.Deadline(); ok {
		defer conn.SetReadDeadline(time.Time{})
	}

	n, err = receiver.readHeader(statusCode, meta)
	if nil != err {
		if ctxErr := ctx.Err(); nil != ctxErr {
			var errs error = erorr.Errors{ErrContextDone, ctxErr, err}
			return n, erorr.Wrap(errs, "could not read mercury protocol response header")
		}
		return n, err
	}

	return n, nil
}

func (receiver *internalResponseReader) Reader(ctx context.Context) io.Reader {
	if nil == receiver {
		return &internalReaderAdapter{}
	}

	if nil == ctx {
		ctx = context.Background()
	}

	return &internalReaderAdapter{
		rr:  receiver,
		ctx: ctx,
	}
}

type internalReaderAdapter struct {
	rr  *internalResponseReader
	ctx context.Context
}

func (receiver *internalReaderAdapter) Read(data []byte) (int, error) {
	return receiver.rr.Read(receiver.ctx, data)
}

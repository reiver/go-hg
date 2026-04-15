package hg

import (
	"context"
	"fmt"
	"io"
	"strings"

	"codeberg.org/reiver/go-field"

	"github.com/reiver/go-hg/internal/io2"
)

// ResponseWriter is used by a Handler to construct a Mercury Protocol response.
//
// For example:
//
//	func serveMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//
//		// ...
//
//	}
//
// Notice that the first parameter is a ResponseWriter.
type ResponseWriter interface {
	Write(ctx context.Context, p []byte) (n int, err error)
	Writer(ctx context.Context) io.Writer
	WriteHeader(ctx context.Context, statusCode int, meta any) (int, error)
}

var _ ResponseWriter = &internalResponseWriter{}

// internalResponseWriter is used to create a ResponseWriter around a io.Writer (such as a net.Conn).
type internalResponseWriter struct {
	writer io2.Writer
	logger Logger
	headerwritten bool
}

func (receiver *internalResponseWriter) Writer(ctx context.Context) io.Writer {
	return io2.ClassicWriter(ctx, receiver.writer)
}

func (receiver *internalResponseWriter) Write(ctx context.Context, data []byte) (n int, err error) {

	if nil == receiver {
		return 0, ErrNilReceiver
	}

	log := mustlogger(receiver.logger).Begin()
	defer log.End()

	if !receiver.headerwritten {
		_, err = receiver.WriteHeader(ctx, StatusSuccess, "application/octet-stream")
		if nil != err {
			return 0, err
		}
	}
	if len(data) <= 0 {
		return 0, nil
	}

	var writer io2.Writer = receiver.writer
	if nil == writer {
		var err error = errNilWriter

		log.Error(
			field.S("failed to write Mercury Protocol response body"),
			field.E(err),
		)

		return 0, err
	}

	{
		m, err := writer.Write(ctx, data)
		n += m
		if nil != err {
			log.Error(
				field.S("failed to write Mercury Protocol response body"),
				field.E(err),
			)
			return n, err
		}
	}

	return n, nil
}

func (receiver *internalResponseWriter) WriteHeader(ctx context.Context, statusCode int, meta any) (int, error) {
	if nil == receiver {
		return 0, ErrNilReceiver
	}

	log := mustlogger(receiver.logger).Begin()
	defer log.End()

	if statusCode < 0 || 100 <= statusCode {
		var err error = ErrBadStatusCode

		log.Error(
			field.S("failed to write Mercury Protocol response header"),
			field.E(err),
		)

		return 0, err
	}

	if receiver.headerwritten {
		var err error = errHeaderAlreadyWritten

		log.Error(
			field.S("failed to write Mercury Protocol response header (again)"),
			field.E(err),
		)

		return 0, err
	}

	var writer io2.Writer = receiver.writer
	if nil == writer {
		var err error = errNilWriter

		log.Error(
			field.S("failed to write Mercury Protocol response header"),
			field.E(err),
		)

		return 0, err
	}

	var header strings.Builder
	{
		fmt.Fprintf(&header, "%02d %s\r\n", statusCode, meta)
		log.Trace(field.S("wrote header"))
	}

	var n int
	{
		var err error

		n, err = io.WriteString(receiver.Writer(ctx), header.String())
		if nil != err {
			log.Error(
				field.S("failed to write Mercury Protocol response header"),
				field.E(err),
			)
			return n, err
		}
		receiver.headerwritten = true
	}

	return n, nil
}

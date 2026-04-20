package hg

import (
	"context"
	"fmt"
	"io"
	"strings"

	"codeberg.org/reiver/go-erorr"
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
	WriteHeader(ctx context.Context, statusCode int, meta string) (int, error)
}

var _ ResponseWriter = &internalResponseWriter{}

// internalResponseWriter is used to create a ResponseWriter around a io.Writer (such as a net.Conn).
type internalResponseWriter struct {
	writer        io2.Writer
	logger        Logger
	headerwritten bool
}

func (receiver *internalResponseWriter) Writer(ctx context.Context) io.Writer {
	if nil == receiver {
		return nil
	}

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

func (receiver *internalResponseWriter) WriteHeader(ctx context.Context, statusCode int, meta string) (int, error) {
	if nil == receiver {
		return 0, ErrNilReceiver
	}

	log := mustlogger(receiver.logger).Begin()
	defer log.End()

	if maxmeta < len(meta) {
		var err error = ErrResponseHeaderMetaTooBig

		const msg string = "response header meta too big"

		log.Error(
			field.S(msg),
			field.String("meta-preview", metaPreview(meta)),
			field.Int("meta-len", len(meta)),
			field.Uint("max", maxmeta),
			field.E(err),
		)

		return 0, erorr.Wrap(err, msg,
			field.String("meta-preview", metaPreview(meta)),
			field.Int("meta-len", len(meta)),
			field.Uint("max", maxmeta),
		)
	}

	// Deal with the case where there is a "\r\n" (or "\r" or "\n") in the `meta` string.
	{
		if 0 <= strings.IndexAny(meta, "\r\n") {
			var err error = ErrBadResponseHeaderMeta

			const msg string = "response header meta contains carriage-return or line-feed"

			log.Error(
				field.S(msg),
				field.String("meta-preview", metaPreview(meta)),
				field.Int("meta-len", len(meta)),
				field.E(err),
			)

			return 0, erorr.Wrap(err, msg,
				field.String("meta-preview", metaPreview(meta)),
				field.Int("meta-len", len(meta)),
			)
		}
	}

	if nil == ctx {
		ctx = context.Background()
	}

	if ctxErr := ctx.Err(); nil != ctxErr {
		var err error = erorr.Errors{ErrContextDone, ctxErr}

		const msg string = "failed to write Mercury Protocol response header"

		log.Error(
			field.S(msg),
			field.E(err),
		)

		return 0, erorr.Wrap(err, msg)
	}

	if statusCode < 0 || 100 <= statusCode {
		var err error = ErrBadStatusCode

		const msg string = "failed to write Mercury Protocol response header"

		log.Error(
			field.S(msg),
			field.E(err),
		)

		return 0, erorr.Wrap(err, msg)
	}

	if receiver.headerwritten {
		var err error = errHeaderAlreadyWritten

		const msg string = "failed to write Mercury Protocol response header (again)"

		log.Error(
			field.S(msg),
			field.E(err),
		)

		return 0, erorr.Wrap(err, msg)
	}

	var writer io2.Writer = receiver.writer
	if nil == writer {
		var err error = errNilWriter

		const msg string = "failed to write Mercury Protocol response header"

		log.Error(
			field.S(msg),
			field.E(err),
		)

		return 0, erorr.Wrap(err, msg)
	}

	var headerBuffer [maxrequest]byte
	var header []byte = headerBuffer[0:0]
	{
		header = appendHeader(header, statusCode, meta)
		log.Trace(field.S("wrote header to buffer"))
	}

	var n int
	{
		var err error

		n, err = receiver.writer.Write(ctx, header)
		if nil != err {
			const msg string = "failed to write Mercury Protocol response header"

			log.Error(
				field.S(msg),
				field.E(err),
			)
			return n, erorr.Wrap(err, msg)
		}
		receiver.headerwritten = true
	}

	return n, nil
}

func appendHeader(p []byte, statusCode int, meta string) []byte {
	p = append(p, fmt.Sprintf("%02d ", statusCode)...)
	p = append(p, meta...)
	p = append(p, "\r\n"...)
	return p
}

// metaPreview returns a log-safe preview of the meta string, capped at 64 bytes
// and stripped of any invalid UTF-8 (which would otherwise happen if the byte
// boundary fell mid-rune).
func metaPreview(meta string) string {
	return strings.ToValidUTF8(meta[:min(64, len(meta))], "")
}

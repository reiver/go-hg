package hg

import (
	"fmt"
	"io"
	"strings"

	"codeberg.org/reiver/go-field"
)

// ResponseWriter is used by a Handler to construct a Mercury Protocol response.
//
// For example:
//
//	func serveMercury(w hg.ResponseWriter, r hg.Request) {
//
//		// ...
//
//	}
//
// Notice that the first parameter is a ResponseWriter.
type ResponseWriter interface {
	io.Writer
	WriteHeader(statusCode int, meta any) (int, error)
}

var _ ResponseWriter = &internalResponseWriter{}

// internalResponseWriter is used to create a ResponseWriter around a io.Writer (such as a net.Conn).
type internalResponseWriter struct {
	Writer io.Writer
	Logger Logger
	headerwritten bool
}

func (receiver *internalResponseWriter) Write(data []byte) (n int, err error) {

	if nil == receiver {
		return 0, errNilReceiver
	}

	log := mustlogger(receiver.Logger).Begin()
	defer log.End()

	if !receiver.headerwritten {
		var m int

		m, err = receiver.WriteHeader(StatusSuccess, "application/octet-stream")
		n += m
		if nil != err {
			return n, err
		}
	}
	if len(data) <= 0 {
		return 0, nil
	}

	var writer io.Writer = receiver.Writer
	if nil == writer {
		return 0, errNilWriter
	}

	{
		m, err := writer.Write(data)
		n += m
		if nil != err {
			return n, err
		}
	}

	return n, nil
}

func (receiver *internalResponseWriter) WriteHeader(statusCode int, meta any) (int, error) {

	if nil == receiver {
		return 0, errNilReceiver
	}

	log := mustlogger(receiver.Logger).Begin()
	defer log.End()

	if statusCode < 0 || 100 <= statusCode {
		return 0, errBadStatusCode
	}

	if receiver.headerwritten {
		log.Error(field.S("header already written"))
		return 0, errHeaderAlreadyWritten
	}

	var writer io.Writer = receiver.Writer
	if nil == writer {
		log.Error(field.S("nil writer"))
		return 0, errNilWriter
	}

	var header strings.Builder
	{
		fmt.Fprintf(&header, "%02d %s\r\n", statusCode, meta)
		log.Trace(field.S("wrote header"))
	}

	var n int
	{
		var err error

		n, err = io.WriteString(writer, header.String())
		if nil != err {
			log.Error(
				field.S("error writing string"),
				field.E(err),
			)
			return n, err
		}
		receiver.headerwritten = true
	}

	return n, nil
}

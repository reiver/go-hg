package hg

import (
	"fmt"
	"io"
	"strings"
)

// ResponseWriter is used by a Handler to construct a Mercury response.
type ResponseWriter interface {
	io.Writer
	WriteHeader(statusCode int, meta interface{}) (int, error)
}

var _ ResponseWriter = &internalResponseWriter{}

// internalResponseWriter is used to create a ResponseWriter around a io.Writer (such as a net.Conn).
type internalResponseWriter struct {
	writer io.Writer
	headerwritten bool
}

func (receiver *internalResponseWriter) Write(data []byte) (n int, err error) {

	if nil == receiver {
		return 0, errNilReceiver
	}
	if !receiver.headerwritten {
		return 0, errHeaderNotAlreadyWritten
	}
	if len(data) <= 0 {
		return 0, nil
	}

	var writer io.Writer = receiver.writer
	if nil == writer {
		return 0, errNilWriter
	}

	return writer.Write(data)
}

func (receiver *internalResponseWriter) WriteHeader(statusCode int, meta interface{}) (int, error) {

	if nil == receiver {
		return 0, errNilReceiver
	}
	if receiver.headerwritten {
		return 0, errHeaderAlreadyWritten
	}

	var writer io.Writer = receiver.writer
	if nil == writer {
		return 0, errNilWriter
	}

	var header strings.Builder
	{
		fmt.Fprintf(&header, "%d %s\r\n", statusCode, meta)
	}

	var n int
	{
		var err error

		n, err = io.WriteString(writer, header.String())
		if nil != err {
			return n, err
		}
		receiver.headerwritten = true
	}

	return n, nil
}

package hg

import (
	"fmt"
	"io"
	"strings"
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
		var m int

		m, err = receiver.WriteHeader(StatusSuccess, "text/gemini")
		n += m
		if nil != err {
			return n, err
		}
	}
	if len(data) <= 0 {
		return 0, nil
	}

	var writer io.Writer = receiver.writer
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

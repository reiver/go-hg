package hg

import (
	"io"
)

// ResponseWriter is used by a Handler to construct a Mercury response.
type ResponseWriter interface {
	io.Writer
	Meta() io.Writer
	WriteHeader(statusCode int)
}

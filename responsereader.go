package hg

import (
	"github.com/reiver/go-utf8s"

	"encoding"
	"fmt"
	"io"
	"strings"
)

// ResponseReader is used by a Handler to read a Mercury Protocol response.
type ResponseReader interface {
	io.Reader
	ReadHeader(statusCode *int, meta interface{}) (int, error)
}

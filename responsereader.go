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

var _ ResponseReader = &internalResponseReader{}

// internalResponseReader is used to create a ResponseReader around a io.Reader (such as a net.Conn).
type internalResponseReader struct {
	reader io.Reader
	headerread bool
}

func (receiver *internalResponseReader) Read(data []byte) (n int, err error) {
	if nil == receiver {
		return 0, errNilReceiver
	}

	var reader io.Reader = receiver.reader
	if nil == reader {
		return 0, errNilReader
	}

	if !receiver.headerread {
		var statusCode int
		var meta string

		m, err := receiver.ReadHeader(&statusCode, &meta)
		n += m
		if nil != err {
			return n, err
		}

		if StatusSuccess != statusCode {
			return n, ErrorResponse(statusCode, meta)
		}
	}

	{
		m, err := reader.Read(data)
		n += m
		if nil != err {
			return n, err
		}
	}

	return n, nil
}

func (receiver *internalResponseReader) ReadHeader(statusCode *int, meta interface{}) (n int, err error) {

	// The Gemini Protocol spec (which the Mercury Protocol is based on) says the meta SHOULD be a maximum of 1024 bytes.
	// We are allowing more than that.
	const maxmeta = 1024 * 2

	if nil == receiver {
		return 0, errNilReceiver
	}

	var reader io.Reader = receiver.reader
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
		if nil == io.EOF {
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
		if nil == io.EOF {
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
		if nil == io.EOF {
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
			if nil == io.EOF {
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
				return n, errorf("response header meta too big — max=%d", maxmeta)
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
				return n, fmt.Errorf("hg: cannot load value for meta into type %T", meta)
			}
		}
	}

	{
		r, size, err := utf8s.ReadRune(reader)
		n += size
		if nil == io.EOF {
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

	return n, nil
}

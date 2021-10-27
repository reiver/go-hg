package hg

import (
	"github.com/reiver/go-utf8s"

	"bytes"
	"fmt"
	"io"
	"strings"
)

// Request represents a Mercury Protocol request — either received by a server, or being sent by a client.
//
// A client may create a Mercury Protocol request with code similar to the following:
//
//	var request hg.Request
//	
//	err := request.Parse("mercury://example.com/apple/banana/cherry.txt")
//
// A server would receive the request as a parameter to the ServeMercury method:
//
//	func (receiver Type) ServeMercury(w hg.ResponseWriter, r hg.Request) {
//		// ...
//	}
type Request struct {
	loaded bool
	value string
}

// nothing returns the ‘nothing’ value for hg.Request; hg.Request is an option-type.
func nothing() Request {
	return Request{}
}

// something returns a ‘something’ value for hg.Request; hg.Request is an option-type.
func something(value string) Request {
	return Request{
		loaded:true,
		value:value + "\r\n",
	}
}

// Parse parses the input ‘value’ and if valid sets the value of the request.
//
// Note that ‘value’ should NOT include the trailing carriage-return and line-feed.
//
// Example Usage:
//
//	var request hg.Request
//	
//	err := request.Parse("mercury://example.com/apple/banana/cherry.txt")
func (receiver *Request) Parse(src interface{}) error {
	if nil == receiver {
		return errNilReceiver
	}
	if nil == src {
		return errNilSource
	}

	var reader io.Reader
	switch casted := src.(type) {
	case io.Reader:
		reader = casted
	case string:
		reader = strings.NewReader(casted)
	case []byte:
		reader = bytes.NewReader(casted)
	case []rune:
		reader = strings.NewReader(string(casted))
	default:
		return fmt.Errorf("hg: cannot parse from type %T", src)
	}

	return receiver.parse(reader)
}


func (receiver *Request) parse(reader io.Reader) error {
	if nil == receiver {
		return errNilReceiver
	}
	if nil == reader {
		return errNilSource
	}

	var storage strings.Builder
	{
		for {
			r, _, err := utf8s.ReadRune(reader)
			if io.EOF == err {
		/////////////// BREAK
				break
			}
			if nil != err {
				return err
			}
			if utf8s.RuneError == r {
				return errRuneError
			}

			if '\r' == r {
				r, _, err := utf8s.ReadRune(reader)
				if io.EOF == err {
					return errExpectedLineFeed
				}
				if nil != err {
					return err
				}
				if utf8s.RuneError == r {
					return errRuneError
				}

				if '\n' == r {
		/////////////////////// BREAK
					break
				}

				return fmt.Errorf("hg: expected a line-feed character but instead got %q", string(r))
			}

			storage.WriteRune(r)
		}
	}

	*receiver = something(storage.String())

	return nil
}

// String returns the full value of the Mercury request.
// Note that this included the trailing carriage-return and line-feed.
//
// String makes Request fit the fmt.Stringer interface.
func (receiver Request) String() string {
	if nothing() == receiver {
		return "⧼nothing⧽"
	}

	return receiver.value
}

// MarshalText makes Request fit the encoding.TextMarshaler interface.
func (receiver Request) MarshalText() ([]byte, error) {
	if nothing() == receiver {
		return nil, errNothing
	}

	return []byte(receiver.value), nil
}

// UnmarshalText  makes Request fit the encoding.TextUnmarshaler interface.
func (receiver *Request) UnmarshalText(text []byte) error {
	if nil == receiver {
		return errNilReceiver
	}

	var value string = string(text)

	return receiver.Parse(value)
}

// WriteTo writers the value of the Mercury request (including the trailing carriage-return and line-feed) to ‘w’ until there's no more to write or when an error occurs.
// The return value ‘n’ is the number of bytes written.
// Any error encountered during the write is also returned.
func (receiver Request) WriteTo(w io.Writer) (int64, error) {
	if nothing() == receiver {
		return 0, errNothing
	}
	if nil == w {
		return 0, errNilWriter
	}

	n, err := io.WriteString(w, receiver.value)

	var n64 int64 = int64(n)

	return n64, err
}

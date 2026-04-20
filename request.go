package hg

import (
	"bytes"
	"io"
	"net"
	gourl "net/url"
	"strings"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
	"github.com/reiver/go-opt"
	"github.com/reiver/go-utf8s"
	"golang.org/x/net/idna"
)

const (
	requestEOL = "\r\n"
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
//	func (receiver Type) ServeMercury(ctx context.Context, w hg.ResponseWriter, r hg.Request) {
//		// ...
//	}
//
// NOTE that the .Parse() methods will accept non-mercury URIs. For example:
//
//	err := request.Parse("gemini://example.com/apple/banana/cherry.txt")
type Request struct {
	optional opt.Optional[string]
}

// someURL returns a ‘something’ value for hg.Request; hg.Request is an option-type.
func someURL(value string) Request {
	return Request{
		optional: opt.Something(value + requestEOL),
	}
}

// RequestValue returns the value of the request without the trailing "\r\n"
//
// For example, if the full value of the request was:
//
//	"mercury://example.com/path/to/file.txt\r\n"
//
// Then RequestValue would return:
//
//	"mercury://example.com/path/to/file.txt"
func (receiver Request) RequestValue() string {
	value, found := receiver.optional.Get()
	if !found {
		return ""
	}
	return value[:len(value)-len(requestEOL)]
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
func (receiver *Request) Parse(src any) error {
	if nil == receiver {
		return ErrNilReceiver
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
		var err error = ErrCannotParse
		return erorr.Wrap(err, "mercury request value cannot be parsed",
			field.FormattedString("type", "%T", src),
		)
	}

	return receiver.parse(reader)
}

func (receiver *Request) parse(reader io.Reader) error {
	if nil == receiver {
		return ErrNilReceiver
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

				return erorr.Stamp("expected a line-feed character but instead got something else",
					field.String("character", string(r)),
				)
			}

			storage.WriteRune(r)

			if (maxrequest-len(requestEOL)) < storage.Len() {
				return errRequestTooLong
			}
		}
	}

	*receiver = someURL(storage.String())

	return nil
}

func (receiver Request) IsNothing() bool {
	return receiver.optional.IsNothing()
}

// Scheme returns the Scheme of the URL/URI/IRI in the request.
//
// For example, if the request is:
//
//	"mercury://example.com/once/twice/thrice/fource.gmni\r\n"
//
// Then scheme would return:
//
//	"mercury"
//
// And, for example, if the request is:
//
//	"gemini://example.com/once/twice/thrice/fource.gmni\r\n"
//
// Then scheme would return:
//
//	"gemini"
func (receiver Request) Scheme() string {
	value, found := receiver.optional.Get()
	if !found {
		return ""
	}

	var length int
	{
		length = len(Scheme)+len(":")
		if newLength := len(SchemeTLS)+len(":"); length < newLength {
			length = newLength
		}
	}
	if length < 0 {
		return ""
	}

	if len(value) < length {
		return ""
	}
	var str string = strings.ToLower(value[:length])

	{
		const needle string = ":"

		var index int = strings.Index(str, needle)
		if index < len(needle) {
			return ""
		}

		str = str[:index]
	}

	return str
}

// String returns the full value of the Mercury request.
// Note that this included the trailing carriage-return and line-feed.
//
// String makes Request fit the fmt.Stringer interface.
func (receiver Request) String() string {
	value, found := receiver.optional.Get()
	if !found {
		return "⧼nothing⧽"
	}

	return value
}

// TCPAddr returns the TCP-address that is embedded in the request.
//
// Example usage for the Mercury Protocol:
//
//	var req hg.Request
//	
//	// ...
//	
//	addr, found := req.TCPAddr()
//	if !found {
//		return errBadRequest
//	}
//	
//	rr, err := hg.DialAndCall(ctx, addr, req)
//
//
// Example usage for the Gemini Protocol:
//
//	var req hg.Request
//	
//	// ...
//	
//	addr, found := req.TCPAddr()
//	if !found {
//		return errBadRequest
//	}
//	
//	rr, err := hg.DialAndCallTLS(ctx, addr, req)
//
// See also:
//
//	• [DialAndCall]
//	• [DialAndCallTLS]
func (receiver Request) TCPAddr() (string, bool) {
	value, found := receiver.optional.Get()
	if !found {
		return "", false
	}

	var str string = value

	// Remove the "\r\n" at the end.
	// We assume it is there without verifying.
	if len(str) < len(requestEOL) {
		return "", false
	}
	str = str[:len(str)-len(requestEOL)]

	url, err := gourl.Parse(str)
	if nil != err {
		return "", false
	}

	var hostName string = url.Hostname()
	if "" == hostName {
		return "", false
	}
	{
		var err error
		hostName, err = idna.ToASCII(hostName)
		if nil != err {
			return "", false
		}
	}
	if "" == hostName {
		return "", false
	}
	hostName = strings.ToLower(hostName)

	var tcpPort string = url.Port()
	if "" == tcpPort {
		switch url.Scheme {
		case Scheme:
			tcpPort = DefaultTCPPortString
		case SchemeTLS:
			tcpPort = DefaultTCPPortTLSString
		default:
			return "", false
		}
	}

	return net.JoinHostPort(hostName, tcpPort), true
}

// MarshalText makes Request fit the encoding.TextMarshaler interface.
func (receiver Request) MarshalText() ([]byte, error) {
	value, found := receiver.optional.Get()
	if !found {
		return nil, errNothing
	}

	return []byte(value), nil
}

// UnmarshalText  makes Request fit the encoding.TextUnmarshaler interface.
func (receiver *Request) UnmarshalText(text []byte) error {
	if nil == receiver {
		return ErrNilReceiver
	}

	var value string = string(text)

	return receiver.Parse(value)
}

// WriteTo writers the value of the Mercury request (including the trailing carriage-return and line-feed) to ‘w’ until there's no more to write or when an error occurs.
// The return value ‘n’ is the number of bytes written.
// Any error encountered during the write is also returned.
func (receiver Request) WriteTo(w io.Writer) (int64, error) {
	value, found := receiver.optional.Get()
	if !found {
		return 0, errNothing
	}
	if nil == w {
		return 0, errNilWriter
	}

	n, err := io.WriteString(w, value)

	var n64 int64 = int64(n)

	return n64, err
}

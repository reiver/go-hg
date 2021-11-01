package hg

import (
	"io"
	"net"
)

// DialAndCall makes a TCP connection to the TCP address given given by ‘addr’,
// and (speaking the Mercury Protocol) sends the request given by ‘request’.
//
// What is given by ‘addr’ might be something like: "11.22.33.44:1961", or "example.com:1961"
//
// What is given by ‘request’ might be a Request containing something like: "mercury://example.com/path/to/file.txt"
func DialAndCall(addr string, request Request) (ResponseReader, error) {

	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}

	return Call(conn, request)
}

// Call uses the TCP connection provided by ‘conn’ and (speaking the Mercury Protocol) sends the request given by ‘request’.
//
// What is given by ‘request’ might be a Request containing something like: "mercury://example.com/path/to/file.txt"
//
// Note that the net.Conn hat provides the TCP connection can be created with code similar to:
//
//	conn, err := net.Dial("tcp", addr)
//
// Where what is given by ‘addr’ might be something like: "11.22.33.44:1961", or "example.com:1961"
func Call(conn net.Conn, request Request) (ResponseReader, error) {

	if nil == conn {
		return nil, errNilNetworkConnection
	}

	_, err := io.WriteString(conn, request.String())
	if nil != err {
		return nil, err
	}

	var rr internalResponseReader
	{
		rr.rc = conn
	}

	return &rr, nil
}

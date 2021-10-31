package hg

import (
	"io"
	"net"
)

func DialAndCall(addr string, request Request) (ResponseReader, error) {

	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}

	return Call(conn, request)
}

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

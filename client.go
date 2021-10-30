package hg

import (
	"io"
	"net"
)

type Client struct {
	Caller Caller
}

func (client *Client) Call(conn net.Conn, request Request) error {
	if nil == client {
		return errNilReceiver
	}

	var caller Caller = client.Caller
	if nil == caller {
		return errNilCaller
	}

	_, err := io.WriteString(conn, request.String())
	if nil != err {
		return err
	}

	var rr internalResponseReader
	{
		rr.reader = conn
	}

	caller.CallMercury(&rr, request)

	return nil
}

func (client *Client) DialAndCall(addr string, request Request) error {
	if nil == client {
		return errNilReceiver
	}

	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return err
	}
	defer conn.Close()

	return client.Call(conn, request)
}

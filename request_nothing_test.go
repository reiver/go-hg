package hg_test

import (
	"github.com/reiver/go-hg"

	"context"
	"errors"
	"net"
	"testing"
)

// TestCall_nothingRequest confirms that calling Call with an uninitialized
// (nothing) Request returns ErrRequestIsNothing without writing anything.
func TestCall_nothingRequest(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request hg.Request // uninitialized — the "nothing" value

	_, err := hg.Call(context.Background(), client, request)
	if nil == err {
		t.Fatal("expected error with nothing request, got nil")
	}

	if !errors.Is(err, hg.ErrRequestIsNothing) {
		t.Fatalf("expected ErrRequestIsNothing, got: %v", err)
	}
}

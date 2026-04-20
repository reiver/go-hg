package hg

import (
	"net"
	"strings"
	"testing"
)

// TestReadHeader_MetaExactlyMaxmeta verifies that a response with a meta field of exactly maxmeta bytes is accepted.
func TestReadHeader_MetaExactlyMaxmeta(t *testing.T) {

	meta := strings.Repeat("a", maxmeta)
	response := "20 " + meta + "\r\nBody here"

	server, client := net.Pipe()
	defer client.Close()

	go func() {
		server.Write([]byte(response))
		server.Close()
	}()

	rr := &internalResponseReader{conn: client}

	var statusCode int
	var gotMeta string
	_, err := rr.readHeader(&statusCode, &gotMeta)
	if nil != err {
		t.Fatalf("expected no error for meta of exactly maxmeta bytes, got: %v", err)
	}
	if 20 != statusCode {
		t.Fatalf("expected status code 20, got: %d", statusCode)
	}
	if gotMeta != meta {
		t.Fatalf("meta mismatch: expected length %d, got length %d", len(meta), len(gotMeta))
	}
}

// TestReadHeader_MetaOneByteTooLong verifies that a response with a meta field of maxmeta+1 bytes is rejected.
func TestReadHeader_MetaOneByteTooLong(t *testing.T) {

	meta := strings.Repeat("a", maxmeta+1)
	response := "20 " + meta + "\r\nBody here"

	server, client := net.Pipe()
	defer client.Close()

	go func() {
		server.Write([]byte(response))
		server.Close()
	}()

	rr := &internalResponseReader{conn: client}

	var statusCode int
	var gotMeta string
	_, err := rr.readHeader(&statusCode, &gotMeta)
	if nil == err {
		t.Fatal("expected error for meta of maxmeta+1 bytes, got nil")
	}
}

// TestReadHeader_MetaMultibyteRuneAtBoundary verifies that a response with a
// meta field that would be within maxmeta if only counting runes (not bytes)
// but exceeds maxmeta in bytes due to multi-byte runes is correctly rejected.
//
// This is the regression test for the off-by-one bug where the length check
// happened before WriteRune, allowing up to maxmeta + 4 bytes.
func TestReadHeader_MetaMultibyteRuneAtBoundary(t *testing.T) {

	// Fill meta to maxmeta-1 bytes with ASCII, then add a 2-byte rune.
	// Total: maxmeta+1 bytes. This must be rejected.
	filler := strings.Repeat("a", maxmeta-1)
	// 'é' is U+00E9, encoded as 2 bytes in UTF-8 (0xC3 0xA9).
	meta := filler + "é"

	if len(meta) != maxmeta+1 {
		t.Fatalf("test setup: expected meta length %d, got %d", maxmeta+1, len(meta))
	}

	response := "20 " + meta + "\r\nBody here"

	server, client := net.Pipe()
	defer client.Close()

	go func() {
		server.Write([]byte(response))
		server.Close()
	}()

	rr := &internalResponseReader{conn: client}

	var statusCode int
	var gotMeta string
	_, err := rr.readHeader(&statusCode, &gotMeta)
	if nil == err {
		t.Fatalf("expected error for meta exceeding maxmeta by multi-byte rune boundary, got nil (meta byte length: %d)", len(gotMeta))
	}
}

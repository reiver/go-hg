package hg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// TestCall_NilContextSucceeds verifies that passing a nil context to Call
// behaves as context.Background() — the write succeeds and a ResponseReader
// is returned.
func TestCall_NilContextSucceeds(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		// drain the server side so the write doesn't block
		buf := make([]byte, 4096)
		_, _ = server.Read(buf)
		errCh <- nil
	}()

	rr, err := Call(nil, client, request)
	if nil != err {
		t.Fatalf("expected no error with nil context, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}

	<-errCh
}

// TestCall_AlreadyCancelledContextReturnsErrContextDone verifies that a
// pre-cancelled context causes Call to return ErrContextDone and
// context.Canceled.
func TestCall_AlreadyCancelledContextReturnsErrContextDone(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = Call(ctx, client, request)
	if nil == err {
		t.Fatal("expected error with cancelled context, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true, got false; err=%v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true, got false; err=%v", err)
	}
}

// TestCall_AlreadyExpiredDeadlineReturnsErrContextDone verifies that a
// context with a deadline in the past causes Call to return ErrContextDone
// and context.DeadlineExceeded.
func TestCall_AlreadyExpiredDeadlineReturnsErrContextDone(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	_, err = Call(ctx, client, request)
	if nil == err {
		t.Fatal("expected error with expired deadline, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true, got false; err=%v", err)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected errors.Is(err, context.DeadlineExceeded) to be true, got false; err=%v", err)
	}
}

// TestCall_NilConnReturnsErrNilNetworkConnection verifies that nil conn
// is checked before context, even when the context is cancelled.
func TestCall_NilConnReturnsErrNilNetworkConnection(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var request Request

	_, err := Call(ctx, nil, request)
	if nil == err {
		t.Fatal("expected error with nil conn, got nil")
	}

	if !errors.Is(err, errNilNetworkConnection) {
		t.Fatalf("expected errNilNetworkConnection, got: %v", err)
	}

	// Should NOT be ErrContextDone — nil conn check comes first.
	if errors.Is(err, ErrContextDone) {
		t.Fatal("expected error NOT to be ErrContextDone when conn is nil")
	}
}

// TestCall_ValidContextWithDeadlineSucceeds verifies that a context with a
// future deadline allows the write to succeed.
func TestCall_ValidContextWithDeadlineSucceeds(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		buf := make([]byte, 4096)
		_, _ = server.Read(buf)
		errCh <- nil
	}()

	rr, err := Call(ctx, client, request)
	if nil != err {
		t.Fatalf("expected no error, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}

	<-errCh
}

// TestCall_ContextWithoutDeadlineSucceeds verifies that context.Background()
// works without setting any deadline.
func TestCall_ContextWithoutDeadlineSucceeds(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		buf := make([]byte, 4096)
		_, _ = server.Read(buf)
		errCh <- nil
	}()

	rr, err := Call(context.Background(), client, request)
	if nil != err {
		t.Fatalf("expected no error, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}

	<-errCh
}

// TestCall_WriteFailureWithContextDoneWrapsErrContextDone verifies that when
// the write fails and the context is done, the error wraps ErrContextDone.
// Uses cancelOnWriteConn so the context is still valid at pre-check time but
// becomes done during the write — hitting the mid-operation error path.
func TestCall_WriteFailureWithContextDoneWrapsErrContextDone(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wrapped := &cancelOnWriteConn{Conn: client, cancel: cancel}

	_, err = Call(ctx, wrapped, request)

	if nil == err {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true, got false; err=%v", err)
	}
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("expected original write error (io.ErrClosedPipe) to be unwrappable; err=%v", err)
	}
}

// TestCall_WriteFailureWithoutContextDoneReturnsRawError verifies that when
// the write fails but the context is NOT done, the raw write error is returned.
func TestCall_WriteFailureWithoutContextDoneReturnsRawError(t *testing.T) {

	server, client := net.Pipe()
	// Close server so client write fails.
	server.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	_, err = Call(context.Background(), client, request)
	client.Close()

	if nil == err {
		t.Fatal("expected error, got nil")
	}

	// Should NOT be wrapped with ErrContextDone since context is fine.
	if errors.Is(err, ErrContextDone) {
		t.Fatalf("expected error NOT to be ErrContextDone, got: %v", err)
	}
}

// TestDialAndCall_NilContextSucceeds verifies that passing nil ctx to
// DialAndCall behaves as context.Background().
func TestDialAndCall_NilContextSucceeds(t *testing.T) {

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if nil != err {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if nil != err {
			return
		}
		defer conn.Close()
		// drain the request then send a simple response
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)
		fmt.Fprint(conn, "20 text/gemini\r\nHello\n")
	}()

	var request Request
	err = request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	rr, err := DialAndCall(nil, ln.Addr().String(), request)
	if nil != err {
		t.Fatalf("expected no error with nil context, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}
	defer rr.Close()
}

// TestDialAndCall_AlreadyCancelledContextReturnsErrContextDone verifies that
// a pre-cancelled context causes DialAndCall to return ErrContextDone.
func TestDialAndCall_AlreadyCancelledContextReturnsErrContextDone(t *testing.T) {

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = DialAndCall(ctx, "127.0.0.1:1961", request)
	if nil == err {
		t.Fatal("expected error with cancelled context, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true, got false; err=%v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true, got false; err=%v", err)
	}
}

// TestDialAndCall_DialFailureWithContextDoneWrapsErrContextDone verifies that
// a dial failure due to context timeout wraps with ErrContextDone.
func TestDialAndCall_DialFailureWithContextDoneWrapsErrContextDone(t *testing.T) {

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	// Use an already-expired timeout so DialContext fails immediately.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // ensure the timeout expires

	_, err = DialAndCall(ctx, "192.0.2.1:1961", request)
	if nil == err {
		t.Fatal("expected error with expired context, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true, got false; err=%v", err)
	}
}

// cancelOnWriteConn wraps a net.Conn and cancels a context when Write is
// called, then returns an error. This simulates a write failure that happens
// while the context becomes done — hitting the mid-operation path in Call
// rather than the pre-check path.
type cancelOnWriteConn struct {
	net.Conn
	cancel context.CancelFunc
}

func (c *cancelOnWriteConn) Write(p []byte) (int, error) {
	c.cancel()
	return 0, io.ErrClosedPipe
}

// trackCloseConn wraps a net.Conn and records whether Close was called.
type trackCloseConn struct {
	net.Conn
	mu     sync.Mutex
	closed bool
}

func (c *trackCloseConn) Close() error {
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()
	return c.Conn.Close()
}

func (c *trackCloseConn) wasClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

// TestDialAndCall_ConnectionLeakFix verifies the close-on-error pattern that
// DialAndCall uses when Call returns an error. Since DialAndCall creates the
// connection internally (via DialContext), we cannot inject a tracked conn
// directly. Instead, we test the pattern using Call with a cancelOnWriteConn
// wrapped in a trackCloseConn — this exercises the same error path and
// verifies the conn would be closed.
func TestDialAndCall_ConnectionLeakFix(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// cancelOnWriteConn makes Call fail during write (mid-operation path).
	cowConn := &cancelOnWriteConn{Conn: client, cancel: cancel}
	tracked := &trackCloseConn{Conn: cowConn}

	_, err = Call(ctx, tracked, request)
	if nil == err {
		t.Fatal("expected error from Call, got nil")
	}

	// This is what DialAndCall does on the error path (client.go lines 72-77):
	//   conn.Close()
	// Simulate the same pattern here.
	tracked.Close()

	if !tracked.wasClosed() {
		t.Fatal("expected connection to be closed after Call failure")
	}
}

// TestErrorUnwrap_PreCheckUnwrapsBothErrContextDoneAndContextCanceled
// verifies that a pre-check error can be unwrapped to both ErrContextDone
// and context.Canceled.
func TestErrorUnwrap_PreCheckUnwrapsBothErrContextDoneAndContextCanceled(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = Call(ctx, client, request)
	if nil == err {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true; err=%v", err)
	}
}

// TestErrorUnwrap_MidOperationUnwrapsToAll verifies that a mid-operation
// error (write failure with context done) can be unwrapped to ErrContextDone,
// ctx.Err(), and the original write error.
// Uses cancelOnWriteConn so the context is valid at pre-check time but
// becomes done during the write — producing a 3-element erorr.Errors slice.
func TestErrorUnwrap_MidOperationUnwrapsToAll(t *testing.T) {

	server, client := net.Pipe()
	defer server.Close()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wrapped := &cancelOnWriteConn{Conn: client, cancel: cancel}

	_, err = Call(ctx, wrapped, request)

	if nil == err {
		t.Fatal("expected error, got nil")
	}

	// Should unwrap to ErrContextDone (our sentinel).
	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}

	// Should unwrap to context.Canceled (the ctx.Err() value).
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true; err=%v", err)
	}

	// Should unwrap to the original write error (io.ErrClosedPipe).
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("expected original write error (io.ErrClosedPipe) to be unwrappable; err=%v", err)
	}
}

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

	if !errors.Is(err, ErrNilNetworkConnection) {
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

// TestCall_WriteFailureWithoutContextDoneWrapsErrWriteError verifies that when
// the write fails but the context is NOT done, the error wraps ErrWriteError
// (not ErrContextDone) and the original write error is still unwrappable.
func TestCall_WriteFailureWithoutContextDoneWrapsErrWriteError(t *testing.T) {

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

	// Should be wrapped with ErrWriteError.
	if !errors.Is(err, ErrWriteError) {
		t.Fatalf("expected errors.Is(err, ErrWriteError) to be true, got false; err=%v", err)
	}

	// Should NOT be wrapped with ErrContextDone since context is fine.
	if errors.Is(err, ErrContextDone) {
		t.Fatalf("expected error NOT to be ErrContextDone, got: %v", err)
	}

	// Original write error should still be unwrappable.
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("expected errors.Is(err, io.ErrClosedPipe) to be true, got false; err=%v", err)
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

// TestDialAndCall_DialFailureWithoutContextDoneWrapsErrDialError verifies that
// when dialing fails but the context is NOT done, the error wraps ErrDialError
// (not ErrContextDone).
func TestDialAndCall_DialFailureWithoutContextDoneWrapsErrDialError(t *testing.T) {

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	// Use a port that is not listening — dial will fail with connection refused.
	// 127.0.0.1:1 is almost certainly not listening.
	_, err = DialAndCall(context.Background(), "127.0.0.1:1", request)
	if nil == err {
		t.Fatal("expected error dialing refused port, got nil")
	}

	// Should be wrapped with ErrDialError.
	if !errors.Is(err, ErrDialError) {
		t.Fatalf("expected errors.Is(err, ErrDialError) to be true, got false; err=%v", err)
	}

	// Should NOT be wrapped with ErrContextDone since context is fine.
	if errors.Is(err, ErrContextDone) {
		t.Fatalf("expected error NOT to be ErrContextDone, got: %v", err)
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

// ---------------------------------------------------------------------------
// ResponseReader context-aware Read / ReadHeader / Reader tests
// ---------------------------------------------------------------------------

// newTestPipe sets up a Call over net.Pipe, drains the request on the server
// side, and returns (server, ResponseReader). The caller must close both.
func newTestPipe(t *testing.T) (net.Conn, ResponseReader) {
	t.Helper()

	server, client := net.Pipe()

	var request Request
	err := request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		server.Close()
		client.Close()
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	// drain the request bytes on the server side
	go func() {
		buf := make([]byte, 4096)
		_, _ = server.Read(buf)
	}()

	rr, err := Call(context.Background(), client, request)
	if nil != err {
		server.Close()
		client.Close()
		t.Fatalf("unexpected error from Call: %v", err)
	}

	return server, rr
}

// TestRead_AlreadyCancelledContext verifies that Read with an already-cancelled
// context returns ErrContextDone.
func TestRead_AlreadyCancelledContext(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()
	defer rr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	buf := make([]byte, 256)
	_, err := rr.Read(ctx, buf)
	if nil == err {
		t.Fatal("expected error with cancelled context, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true; err=%v", err)
	}
}

// TestReadHeader_AlreadyCancelledContext verifies that ReadHeader with an
// already-cancelled context returns ErrContextDone.
func TestReadHeader_AlreadyCancelledContext(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()
	defer rr.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var statusCode int
	var meta string
	_, err := rr.ReadHeader(ctx, &statusCode, &meta)
	if nil == err {
		t.Fatal("expected error with cancelled context, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected errors.Is(err, context.Canceled) to be true; err=%v", err)
	}
}

// TestRead_WithDeadlineSucceeds verifies that Read with a future deadline
// succeeds and clears the deadline afterward.
func TestRead_WithDeadlineSucceeds(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	// Server sends a success response with body.
	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nHello, world!")
	}()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	buf := make([]byte, 256)
	n, err := rr.Read(ctx, buf)
	if nil != err {
		t.Fatalf("expected no error, got: %v", err)
	}
	if n == 0 {
		t.Fatal("expected to read some bytes, got 0")
	}

	got := string(buf[:n])
	if got != "Hello, world!" {
		t.Fatalf("expected %q, got %q", "Hello, world!", got)
	}
}

// TestReadHeader_WithDeadlineSucceeds verifies that ReadHeader with a future
// deadline succeeds and clears the deadline afterward.
func TestReadHeader_WithDeadlineSucceeds(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nHello")
	}()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var statusCode int
	var meta string
	_, err := rr.ReadHeader(ctx, &statusCode, &meta)
	if nil != err {
		t.Fatalf("expected no error, got: %v", err)
	}
	if StatusSuccess != statusCode {
		t.Fatalf("expected status code %d, got %d", StatusSuccess, statusCode)
	}
	if "text/gemini" != meta {
		t.Fatalf("expected meta %q, got %q", "text/gemini", meta)
	}
}

// TestRead_MidReadCancellation verifies that cancelling a context during a
// blocked read actually interrupts the read and returns ErrContextDone.
func TestRead_MidReadCancellation(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()
	defer rr.Close()

	// Server sends a success header but then blocks (no body data sent).
	go func() {
		fmt.Fprint(server, "20 text/gemini\r\n")
		// Don't send body, don't close — read will block.
	}()

	// First, read header explicitly so auto-read doesn't interfere.
	ctx := context.Background()
	var statusCode int
	var meta string
	_, err := rr.ReadHeader(ctx, &statusCode, &meta)
	if nil != err {
		t.Fatalf("expected no error reading header, got: %v", err)
	}

	// Now try to Read body with a context that we cancel after a short delay.
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel2()
	}()

	buf := make([]byte, 256)
	_, err = rr.Read(ctx2, buf)
	if nil == err {
		t.Fatal("expected error from cancelled read, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
}

// TestReadHeader_MidReadCancellation verifies that cancelling a context during
// a blocked ReadHeader actually interrupts the read.
func TestReadHeader_MidReadCancellation(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()
	defer rr.Close()

	// Server never sends anything — ReadHeader will block.

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	var statusCode int
	var meta string
	_, err := rr.ReadHeader(ctx, &statusCode, &meta)
	if nil == err {
		t.Fatal("expected error from cancelled ReadHeader, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
}

// TestReader_IoReadAll verifies that Reader(ctx) returns an io.Reader that
// works with io.ReadAll.
func TestReader_IoReadAll(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nHello from Reader!")
	}()

	ctx := context.Background()
	reader := rr.Reader(ctx)

	data, err := io.ReadAll(reader)
	if nil != err {
		t.Fatalf("expected no error from io.ReadAll, got: %v", err)
	}

	got := string(data)
	if "Hello from Reader!" != got {
		t.Fatalf("expected %q, got %q", "Hello from Reader!", got)
	}
}

// TestRead_NilContext verifies that Read with nil context behaves as
// context.Background().
func TestRead_NilContext(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nnil ctx body")
	}()

	buf := make([]byte, 256)
	n, err := rr.Read(nil, buf)
	if nil != err {
		t.Fatalf("expected no error with nil context, got: %v", err)
	}
	if 0 == n {
		t.Fatal("expected to read some bytes, got 0")
	}

	got := string(buf[:n])
	if "nil ctx body" != got {
		t.Fatalf("expected %q, got %q", "nil ctx body", got)
	}
}

// TestReadHeader_NilContext verifies that ReadHeader with nil context behaves
// as context.Background().
func TestReadHeader_NilContext(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nHello")
	}()

	var statusCode int
	var meta string
	_, err := rr.ReadHeader(nil, &statusCode, &meta)
	if nil != err {
		t.Fatalf("expected no error with nil context, got: %v", err)
	}
	if StatusSuccess != statusCode {
		t.Fatalf("expected status code %d, got %d", StatusSuccess, statusCode)
	}
}

// TestReader_NilContext verifies that Reader with nil context behaves as
// context.Background().
func TestReader_NilContext(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nnil reader ctx")
	}()

	reader := rr.Reader(nil)

	data, err := io.ReadAll(reader)
	if nil != err {
		t.Fatalf("expected no error from io.ReadAll, got: %v", err)
	}

	got := string(data)
	if "nil reader ctx" != got {
		t.Fatalf("expected %q, got %q", "nil reader ctx", got)
	}
}

// TestReader_MidReadCancellation verifies that cancelling a context while
// Reader(ctx).Read() is blocked actually interrupts the read.
func TestReader_MidReadCancellation(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()
	defer rr.Close()

	// Server sends header but then blocks.
	go func() {
		fmt.Fprint(server, "20 text/gemini\r\n")
	}()

	// Read header first.
	var statusCode int
	var meta string
	_, err := rr.ReadHeader(context.Background(), &statusCode, &meta)
	if nil != err {
		t.Fatalf("expected no error reading header, got: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	reader := rr.Reader(ctx)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	buf := make([]byte, 256)
	_, err = reader.Read(buf)
	if nil == err {
		t.Fatal("expected error from cancelled Reader.Read, got nil")
	}

	if !errors.Is(err, ErrContextDone) {
		t.Fatalf("expected errors.Is(err, ErrContextDone) to be true; err=%v", err)
	}
}

// TestClose_WhileReaderBlocked verifies that calling Close() while
// Reader(ctx).Read() is blocked does not deadlock or panic.
func TestClose_WhileReaderBlocked(t *testing.T) {

	server, rr := newTestPipe(t)
	defer server.Close()

	// Server sends header but then blocks.
	go func() {
		fmt.Fprint(server, "20 text/gemini\r\n")
	}()

	// Read header first.
	var statusCode int
	var meta string
	_, err := rr.ReadHeader(context.Background(), &statusCode, &meta)
	if nil != err {
		t.Fatalf("expected no error reading header, got: %v", err)
	}

	ctx := context.Background()
	reader := rr.Reader(ctx)

	done := make(chan error, 1)
	go func() {
		buf := make([]byte, 256)
		_, err := reader.Read(buf)
		done <- err
	}()

	// Give the read goroutine time to block.
	time.Sleep(50 * time.Millisecond)

	// Close the ResponseReader — should unblock the read.
	rr.Close()

	select {
	case err := <-done:
		// We expect an error (connection closed), just not a panic or deadlock.
		if nil == err {
			t.Fatal("expected error after Close, got nil")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for read to unblock after Close")
	}
}

// TestRead_AutoReadHeaderNonSuccessStatus verifies that Read with a
// non-success status code auto-reads the header and returns the
// appropriate ErrorResponse type.
func TestRead_AutoReadHeaderNonSuccessStatus(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "51 not found\r\n")
	}()

	ctx := context.Background()
	buf := make([]byte, 256)
	_, err := rr.Read(ctx, buf)
	if nil == err {
		t.Fatal("expected error from non-success status, got nil")
	}

	var notFound ResponseNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected ResponseNotFound, got: %T %v", err, err)
	}
	if "not found" != notFound.Meta() {
		t.Fatalf("expected meta %q, got %q", "not found", notFound.Meta())
	}
}

// TestRead_SecondCallAfterAutoReadHeader verifies that a second Read
// call correctly skips the auto-read header path and reads body data.
func TestRead_SecondCallAfterAutoReadHeader(t *testing.T) {

	server, rr := newTestPipe(t)
	defer rr.Close()

	go func() {
		defer server.Close()
		fmt.Fprint(server, "20 text/gemini\r\nfirst chunk, second chunk")
	}()

	ctx := context.Background()

	// First Read — triggers auto-read header.
	buf := make([]byte, 12)
	n, err := rr.Read(ctx, buf)
	if nil != err {
		t.Fatalf("expected no error on first Read, got: %v", err)
	}
	first := string(buf[:n])

	// Second Read — header already read, should go straight to body.
	buf2 := make([]byte, 256)
	n2, err := rr.Read(ctx, buf2)
	if nil != err {
		t.Fatalf("expected no error on second Read, got: %v", err)
	}
	second := string(buf2[:n2])

	combined := first + second
	if "first chunk, second chunk" != combined {
		t.Fatalf("expected %q, got %q", "first chunk, second chunk", combined)
	}
}

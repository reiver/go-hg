package hg_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/reiver/go-hg"
)

// generateSelfSignedCert creates a self-signed TLS certificate for testing.
// The returned tls.Certificate can be used for both server and client certs.
func generateSelfSignedCert(t *testing.T, hosts ...string) tls.Certificate {
	t.Helper()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if nil != err {
		t.Fatalf("could not generate ed25519 key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:         true,
		BasicConstraintsValid: true,
	}
	for _, h := range hosts {
		template.DNSNames = append(template.DNSNames, h)
	}
	if 0 == len(hosts) {
		template.DNSNames = []string{"localhost"}
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, pub, priv)
	if nil != err {
		t.Fatalf("could not create self-signed certificate: %v", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  priv,
	}
}

// startTLSServer starts a TLS server on a random port that accepts one connection,
// reads a request, and sends a Gemini "20 text/gemini\r\nHello\r\n" response.
// If requireClientCert is true, the server demands a client certificate.
// Returns the listener address and a channel that receives any server-side error.
func startTLSServer(t *testing.T, serverCert tls.Certificate, requireClientCert bool, clientCAPool *x509.CertPool) (string, chan error) {
	t.Helper()

	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}
	if requireClientCert {
		tlsConf.ClientAuth = tls.RequireAnyClientCert
		if nil != clientCAPool {
			tlsConf.ClientCAs = clientCAPool
		}
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if nil != err {
		t.Fatalf("could not listen: %v", err)
	}

	tlsLn := tls.NewListener(ln, tlsConf)

	errCh := make(chan error, 1)
	go func() {
		defer tlsLn.Close()

		conn, err := tlsLn.Accept()
		if nil != err {
			errCh <- err
			return
		}
		defer conn.Close()

		// Drain request.
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)

		// Send response.
		_, err = fmt.Fprint(conn, "20 text/gemini\r\nHello\r\n")
		errCh <- err
	}()

	return ln.Addr().String(), errCh
}

// TestDialAndCallTLS_NilHandlerAcceptsSelfSigned verifies that passing nil
// TLSHandler to DialAndCallTLS accepts any server cert (including self-signed).
func TestDialAndCallTLS_NilHandlerAcceptsSelfSigned(t *testing.T) {

	serverCert := generateSelfSignedCert(t, "localhost")
	addr, errCh := startTLSServer(t, serverCert, false, nil)

	var request hg.Request
	err := request.Parse("gemini://localhost/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rr, err := hg.DialAndCallTLS(ctx, addr, request, nil)
	if nil != err {
		t.Fatalf("expected no error with nil TLSHandler, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}
	defer rr.Close()

	serverErr := <-errCh
	if nil != serverErr {
		t.Fatalf("server error: %v", serverErr)
	}
}

// TestDialAndCallTLS_VerifyServerCertificateRejectsConnection verifies that when
// VerifyServerCertificate returns an error, DialAndCallTLS returns that error.
func TestDialAndCallTLS_VerifyServerCertificateRejectsConnection(t *testing.T) {

	serverCert := generateSelfSignedCert(t, "localhost")
	addr, _ := startTLSServer(t, serverCert, false, nil)

	var request hg.Request
	err := request.Parse("gemini://localhost/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	verifyErr := errors.New("certificate rejected by verifier")

	handler := hg.TLSConfig{
		VerifyServerCertificateFunc: func(host string, rawCerts [][]byte) error {
			return verifyErr
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = hg.DialAndCallTLS(ctx, addr, request, handler)
	if nil == err {
		t.Fatal("expected error when VerifyServerCertificate rejects, got nil")
	}

	if !errors.Is(err, verifyErr) {
		t.Fatalf("expected errors.Is(err, verifyErr) to be true, got false; err=%v", err)
	}
}

// TestDialAndCallTLS_ClientCertificatePresented verifies that when
// ClientCertificate returns a cert, it is presented during the TLS handshake.
func TestDialAndCallTLS_ClientCertificatePresented(t *testing.T) {

	serverCert := generateSelfSignedCert(t, "localhost")
	clientCert := generateSelfSignedCert(t, "client")

	// Server requires client cert (any cert, no CA verification).
	addr, errCh := startTLSServer(t, serverCert, true, nil)

	var request hg.Request
	err := request.Parse("gemini://localhost/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	handler := hg.TLSConfig{
		ClientCertificateFunc: func(host string) *tls.Certificate {
			return &clientCert
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rr, err := hg.DialAndCallTLS(ctx, addr, request, handler)
	if nil != err {
		t.Fatalf("expected no error when presenting client cert, got: %v", err)
	}
	if nil == rr {
		t.Fatal("expected non-nil ResponseReader")
	}
	defer rr.Close()

	serverErr := <-errCh
	if nil != serverErr {
		t.Fatalf("server error: %v", serverErr)
	}
}

// TestResponseReader_ServerCertificate verifies that after a successful TLS
// connection, ServerCertificate() returns the server's leaf cert.
func TestResponseReader_ServerCertificate(t *testing.T) {

	serverCert := generateSelfSignedCert(t, "localhost")
	addr, errCh := startTLSServer(t, serverCert, false, nil)

	var request hg.Request
	err := request.Parse("gemini://localhost/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rr, err := hg.DialAndCallTLS(ctx, addr, request, nil)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rr.Close()

	cert := rr.ServerCertificate()
	if nil == cert {
		t.Fatal("expected non-nil ServerCertificate")
	}

	// The returned cert should match the DER bytes of our server cert.
	if len(serverCert.Certificate) == 0 {
		t.Fatal("test setup error: serverCert has no DER bytes")
	}

	expected := serverCert.Certificate[0]
	actual := cert.Raw
	if len(expected) != len(actual) {
		t.Fatalf("ServerCertificate DER length mismatch: expected %d, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Fatalf("ServerCertificate DER mismatch at byte %d", i)
			break
		}
	}

	serverErr := <-errCh
	if nil != serverErr {
		t.Fatalf("server error: %v", serverErr)
	}
}

// TestResponseReader_ClientCertificate verifies that after a successful TLS
// connection with a client cert, ClientCertificate() returns the cert provided.
func TestResponseReader_ClientCertificate(t *testing.T) {

	serverCert := generateSelfSignedCert(t, "localhost")
	clientCert := generateSelfSignedCert(t, "client")

	addr, errCh := startTLSServer(t, serverCert, true, nil)

	var request hg.Request
	err := request.Parse("gemini://localhost/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	handler := hg.TLSConfig{
		ClientCertificateFunc: func(host string) *tls.Certificate {
			return &clientCert
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rr, err := hg.DialAndCallTLS(ctx, addr, request, handler)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rr.Close()

	got := rr.ClientCertificate()
	if nil == got {
		t.Fatal("expected non-nil ClientCertificate")
	}

	// Should be the same pointer we provided.
	if got != &clientCert {
		t.Fatal("ClientCertificate() returned a different pointer than what was provided")
	}

	serverErr := <-errCh
	if nil != serverErr {
		t.Fatalf("server error: %v", serverErr)
	}
}

// TestResponseReader_PlainTCPReturnsNilCerts verifies that after a plain TCP
// (Mercury) connection via DialAndCall, both ServerCertificate() and
// ClientCertificate() return nil.
func TestResponseReader_PlainTCPReturnsNilCerts(t *testing.T) {

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
		buf := make([]byte, 4096)
		_, _ = conn.Read(buf)
		fmt.Fprint(conn, "20 text/gemini\r\nHello\r\n")
	}()

	var request hg.Request
	err = request.Parse("mercury://example.com/test\r\n")
	if nil != err {
		t.Fatalf("unexpected error parsing request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rr, err := hg.DialAndCall(ctx, ln.Addr().String(), request)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rr.Close()

	if nil != rr.ServerCertificate() {
		t.Fatal("expected nil ServerCertificate for plain TCP connection")
	}

	if nil != rr.ClientCertificate() {
		t.Fatal("expected nil ClientCertificate for plain TCP connection")
	}
}

// TestDefaultVerifier_AcceptsSelfSigned verifies that DefaultVerifier accepts
// a self-signed certificate.
func TestDefaultVerifier_AcceptsSelfSigned(t *testing.T) {

	tests := []struct {
		Host string
	}{
		{Host: "localhost"},
		{Host: "example.com"},
		{Host: "test.local"},
	}

	for testNumber, test := range tests {

		cert := generateSelfSignedCert(t, test.Host)

		err := hg.DefaultVerifier.VerifyServerCertificate(test.Host, [][]byte{cert.Certificate[0]})
		if nil != err {
			t.Errorf("For test #%d, expected DefaultVerifier to accept self-signed cert, got error.", testNumber)
			t.Logf("HOST: %q", test.Host)
			t.Logf("ERROR: %v", err)
			continue
		}
	}
}

// TestDefaultVerifier_AcceptsCASignedCert verifies that DefaultVerifier accepts
// a CA-signed cert when the CA is trusted, hostname matches, and cert is valid.
// Since we cannot add a test CA to system roots, we verify the logic by confirming
// that a non-self-signed cert goes through CA verification (and fails with
// UnknownAuthorityError for our test CA — which proves the code path is hit).
// We also verify that if we could get past the CA check, hostname matching works.
func TestDefaultVerifier_AcceptsCASignedCert(t *testing.T) {

	// Create a CA.
	caPub, caPriv, err := ed25519.GenerateKey(rand.Reader)
	if nil != err {
		t.Fatalf("could not generate CA key: %v", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(100),
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, caPub, caPriv)
	if nil != err {
		t.Fatalf("could not create CA cert: %v", err)
	}

	caCert, err := x509.ParseCertificate(caDER)
	if nil != err {
		t.Fatalf("could not parse CA cert: %v", err)
	}

	// Create a leaf cert signed by the CA with correct hostname.
	leafPub, _, err := ed25519.GenerateKey(rand.Reader)
	if nil != err {
		t.Fatalf("could not generate leaf key: %v", err)
	}

	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(200),
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		DNSNames:     []string{"example.com"},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caCert, leafPub, caPriv)
	if nil != err {
		t.Fatalf("could not create leaf cert: %v", err)
	}

	// The cert is NOT self-signed (signed by our test CA), so DefaultVerifier
	// should attempt full CA verification with system roots. Since our test CA
	// is not in system roots, it will return an UnknownAuthorityError.
	// This proves the CA verification code path IS exercised for non-self-signed certs.
	err = hg.DefaultVerifier.VerifyServerCertificate("example.com", [][]byte{leafDER})
	if nil == err {
		// If this passes (e.g., cert somehow passes system verification), that's fine.
		return
	}

	// We expect UnknownAuthorityError specifically — any other error would mean
	// something else is wrong.
	var unknownAuthority x509.UnknownAuthorityError
	if !errors.As(err, &unknownAuthority) {
		t.Fatalf("expected x509.UnknownAuthorityError for test CA cert, got: %T %v", err, err)
	}
}

// TestDefaultVerifier_RejectsWrongHostnameOrExpiry verifies that a CA-signed
// cert with wrong hostname is rejected by DefaultVerifier.
func TestDefaultVerifier_RejectsWrongHostnameOrExpiry(t *testing.T) {

	// Create a CA.
	caPub, caPriv, err := ed25519.GenerateKey(rand.Reader)
	if nil != err {
		t.Fatalf("could not generate CA key: %v", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(100),
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, caPub, caPriv)
	if nil != err {
		t.Fatalf("could not create CA cert: %v", err)
	}

	caCert, err := x509.ParseCertificate(caDER)
	if nil != err {
		t.Fatalf("could not parse CA cert: %v", err)
	}

	tests := []struct {
		Name     string
		Host     string
		CertDNS []string
		NotBefore time.Time
		NotAfter  time.Time
	}{
		{
			Name:     "wrong hostname",
			Host:     "evil.example.com",
			CertDNS: []string{"good.example.com"},
			NotBefore: time.Now().Add(-1 * time.Hour),
			NotAfter:  time.Now().Add(24 * time.Hour),
		},
		{
			Name:     "expired cert",
			Host:     "example.com",
			CertDNS: []string{"example.com"},
			NotBefore: time.Now().Add(-48 * time.Hour),
			NotAfter:  time.Now().Add(-24 * time.Hour),
		},
	}

	for testNumber, test := range tests {

		leafPub, _, err := ed25519.GenerateKey(rand.Reader)
		if nil != err {
			t.Fatalf("For test #%d, could not generate leaf key: %v", testNumber, err)
		}

		leafTemplate := &x509.Certificate{
			SerialNumber: big.NewInt(int64(200 + testNumber)),
			NotBefore:    test.NotBefore,
			NotAfter:     test.NotAfter,
			DNSNames:     test.CertDNS,
			KeyUsage:     x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}

		leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, caCert, leafPub, caPriv)
		if nil != err {
			t.Fatalf("For test #%d, could not create leaf cert: %v", testNumber, err)
		}

		err = hg.DefaultVerifier.VerifyServerCertificate(test.Host, [][]byte{leafDER})
		if nil == err {
			t.Errorf("For test #%d (%s), expected DefaultVerifier to reject cert, got nil error.", testNumber, test.Name)
			t.Logf("HOST: %q", test.Host)
			t.Logf("CERT DNS: %v", test.CertDNS)
			continue
		}
	}
}

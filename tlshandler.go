package hg

import (
	"crypto/tls"
	"crypto/x509"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// TLSHandler handles TLS-layer decisions during connection setup.
//
// Used only for Gemini Protocol.
// Ignored for the Mercury Protocol.
//
// See also:
//
//	• [DialAndCallTLS]
//	• [DialAndCallURL]
type TLSHandler interface {

	// VerifyServerCertificate is called after the TLS handshake with the server's certificate chain (DER-encoded, leaf first).
	// rawCerts may be empty if the server sends no certificates.
	// Return nil to accept the connection, or an error to reject it.
	VerifyServerCertificate(hostname string, rawCerts [][]byte) error

	// ClientCertificate is called when a client certificate is needed.
	// 'hostname' is the host being connected to without TCP-port.
	// Return nil to proceed without presenting one.
	ClientCertificate(hostname string) *tls.Certificate
}

// TLSConfig is a [TLSHandler] implementation built from function fields.
// Nil fields are treated as no-ops (accept any server, no client cert).
type TLSConfig struct {
	VerifyServerCertificateFunc func(hostname string, rawCerts [][]byte) error
	ClientCertificateFunc       func(hostname string) *tls.Certificate
}

func (receiver TLSConfig) VerifyServerCertificate(hostname string, rawCerts [][]byte) error {
	if nil == receiver.VerifyServerCertificateFunc {
		return nil
	}
	return receiver.VerifyServerCertificateFunc(hostname, rawCerts)
}

func (receiver TLSConfig) ClientCertificate(hostname string) *tls.Certificate {
	if nil == receiver.ClientCertificateFunc {
		return nil
	}
	return receiver.ClientCertificateFunc(hostname)
}

// defaultVerifier is the internal type for DefaultVerifier.
type defaultVerifier struct{}

// DefaultVerifier is a TLSHandler that applies CA verification only when the server's certificate chains to a system-trusted CA.
// Self-signed certs are accepted without CA verification (for TOFU workflows).
//
// Behavior:
//
//	• Self-signed cert (CheckSignatureFrom succeeds) -> accepted
//	• CA-signed cert, valid -> accepted
//	• CA-signed cert, expired or wrong hostname -> rejected
//	• No certificates presented -> rejected
var DefaultVerifier TLSHandler = defaultVerifier{}

func (defaultVerifier) VerifyServerCertificate(hostname string, rawCerts [][]byte) error {
	if len(rawCerts) <= 0 {
		var err error = ErrServerCertificateNotFound

		return erorr.Wrap(err, "failed to verify server certificates",
			field.String("hostname", hostname),
		)
	}

	cert, err := x509.ParseCertificate(rawCerts[0])
	if nil != err {
		return err
	}

	// Check if the certificate is self-signed by verifying it signed itself.
	// A self-signed cert has Issuer == Subject AND validates its own signature.
	if nil == cert.CheckSignatureFrom(cert) {
		// If we got here, then this server certificate is self-signed.
		// We accept self-signed server certificates unconditionally.
		//
		// This is NOT TOFU (Trust On First Use), but an alternative to it.
		//
		// Implementations that want TOFU should use (or implement) a different [TLSHandler].
		return nil
	}

	// Not self-signed — full CA verification with system roots.
	_, err = cert.Verify(x509.VerifyOptions{
		DNSName: hostname,
	})
	return err
}

func (defaultVerifier) ClientCertificate(hostname string) *tls.Certificate {
	return nil
}

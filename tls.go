package hg

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"time"

	"codeberg.org/reiver/go-erorr"
)

// GenerateClientTLSConfig generates a self-signed TLS client certificate and returns
// a *tls.Config containing it.
//
// This is useful when speaking protocols that use TLS client certificates for identity
// (such as the Gemini Protocol). In the Gemini Protocol, a server may respond with a
// "60 CLIENT CERTIFICATE REQUIRED" status code, indicating the client must present a
// certificate. Gemini clients routinely generate throwaway or per-site self-signed
// certificates on the fly when this happens.
//
// Without this helper, generating a self-signed client certificate requires importing
// crypto/x509, crypto/ed25519, math/big, time, etc. and writing ~30-40 lines of
// boilerplate. GenerateClientTLSConfig handles all of that.
//
// Example usage:
//
//	rr, err := hg.DialAndCallTLS(ctx, addr, request, nil)
//	// ... get back ResponseCertificateRequired ...
//
//	tlsConf, err := hg.GenerateClientTLSConfig()
//	rr, err = hg.DialAndCallTLS(ctx, addr, request, tlsConf)
//
// See also:
//
//	• [DialAndCallTLS]
//	• [ResponseCertificateRequired]
func GenerateClientTLSConfig() (*tls.Config, error) {

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if nil != err {
		return nil, erorr.Wrap(err, "could not generate ed25519 key pair for tls client certificate")
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
	if nil != err {
		return nil, erorr.Wrap(err, "could not create self-signed tls client certificate")
	}

	cert := tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  privateKey,
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

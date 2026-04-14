package hg

import (
	"fmt"
)

// Constants for the Mercury Protocol status codes.
//
// Can use, for example, with ResponseWriter's WriteHeader method.
//
// For example:
//
//	func ServerMercury(w hg.ResponseWriter, r hg.Request) {
//		w.WriteHeader(hg.StatusNotFound, "uh oh!")
//	}
const (
	StatusInput               = 10
	StatusSensitiveInput      = 11

	StatusSuccess             = 20

	StatusRedirectTemporary   = 30
	StatusRedirectPermanent   = 31

	StatusTemporaryFailure    = 40
	StatusServerUnavailable   = 41
	StatusCGIError            = 42
	StatusProxyError          = 43
	StatusSlowDown            = 44

	StatusPermanentFailure    = 50
	StatusNotFound            = 51
	StatusGone                = 52
	StatusProxyRequestRefused = 53
	StatusBadRequest          = 59

	// Not exported since, unlike Gemini Protocol, Mercury Protocol doesn't support certficates.
	statusCertificateRequired      = 60
	statusCertificateNotAuthorized = 61
	statusCertificateNotValid      = 62
)

const (
	StatusTextInput               = "input"                 // 10
	StatusTextSensitiveInput      = "sensitive‐input"       // 11

	StatusTextSuccess             = "success"               // 20

	StatusTextRedirectTemporary   = "temporary-redirection" // 30
	StatusTextRedirectPermanent   = "permanent-redirection" // 31

	StatusTextTemporaryFailure    = "temporary‐failure"     // 40
	StatusTextServerUnavailable   = "server‐unavailable"    // 41
	StatusTextCGIError            = "cgi‐error"             // 42
	StatusTextProxyError          = "proxy‐error"           // 43
	StatusTextSlowDown            = "slow-down"             // 44

	StatusTextPermanentFailure    = "permanent‐failure"     // 50
	StatusTextNotFound            = "not‐found"             // 51
	StatusTextGone                = "gone"                  // 52
	StatusTextProxyRequestRefused = "proxy‐request‐refused" // 53
	StatusTextBadRequest          = "bad‐request"           // 59

	// Not exported since, unlike Gemini Protocol, Mercury Protocol doesn't support certficates.
	statusTextCertificateRequired      = "certificate-required" // 60
	statusTextCertificateNotAuthorized = "certificate-not-authorized" // 61
	statusTextCertificateNotValid      = "certificate-not-valid" // 62
)

func StatusText(code int) string {
	switch code {
	case StatusInput:              // 10
		return StatusTextInput
	case StatusSensitiveInput:     // 11
		return StatusTextSensitiveInput

	case StatusRedirectTemporary:   // 30
		return StatusTextRedirectTemporary
	case StatusRedirectPermanent:   // 31
		return StatusTextRedirectPermanent

	case StatusSuccess:            // 20
		return StatusTextSuccess

	case StatusTemporaryFailure:    // 40
		return StatusTextTemporaryFailure
	case StatusServerUnavailable:   // 41
		return StatusTextServerUnavailable
	case StatusCGIError:            // 42
		return StatusTextCGIError
	case StatusProxyError:          // 43
		return StatusTextProxyError
	case StatusSlowDown:            // 44
		return StatusTextSlowDown

	case StatusPermanentFailure:    // 50
		return StatusTextPermanentFailure
	case StatusNotFound:            // 51
		return StatusTextNotFound
	case StatusGone:                // 52
		return StatusTextGone
	case StatusProxyRequestRefused: // 53
		return StatusTextProxyRequestRefused
	case StatusBadRequest:          // 59
		return StatusTextBadRequest

	case statusCertificateRequired:      // 60
		return statusTextCertificateRequired
	case statusCertificateNotAuthorized: // 61
		return statusTextCertificateNotAuthorized
	case statusCertificateNotValid:      // 62
		return statusTextCertificateNotValid

	default:
		return fmt.Sprintf("meta-%d", code)

	}
}

package hg

import (
	"fmt"
)

// ResponseInput represents a Mercury Protocol “10 INPUT” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseInput               struct {meta string} // 10
// ResponseSensitiveInput represents a Mercury Protocol “11 SENSITIVE INPUT” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseSensitiveInput      struct {meta string} // 11

// ResponseRedirectTemporary represents a Mercury Protocol “30 REDIRECT ‐ TEMPORARY” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseRedirectTemporary   struct {meta string} // 30
// ResponseRedirectPermanent represents a Mercury Protocol “31 REDIRECT ‐ PERMANENT” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseRedirectPermanent   struct {meta string} // 31

// ResponseTemporaryFailure represents a Mercury Protocol “40 TEMPORARY FAILURE” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseTemporaryFailure    struct {meta string} // 40
// ResponseServerUnavailable represents a Mercury Protocol “41 SERVER UNAVAILABLE” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseServerUnavailable   struct {meta string} // 41
// ResponseCGIError represents a Mercury Protocol “42 CGI ERROR” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseCGIError            struct {meta string} // 42
// ResponseProxyError represents a Mercury Protocol “43 PROXY ERROR” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseProxyError          struct {meta string} // 43
// ResponseSlowDown represents a Mercury Protocol “44 SLOW DOWN” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseSlowDown            struct {meta string} // 44

// ResponsePermanentFailure represents a Mercury Protocol “50 PERMANENT FAILURE” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponsePermanentFailure    struct {meta string} // 50
// ResponseNotFound represents a Mercury Protocol “51 NOT FOUND” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseNotFound            struct {meta string} // 51
// ResponseGone represents a Mercury Protocol “52 GONE” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseGone                struct {meta string} // 52
// ResponseProxyRequestRefused represents a Mercury Protocol “53 PROXY REQUEST REFUSED” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseProxyRequestRefused struct {meta string} // 53
// ResponseBadRequest represents a Mercury Protocol “59 BAD REQUEST” response. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type ResponseBadRequest          struct {meta string} // 59

// UnknownResponse represents a Mercury Protocol unknown response that this package doesn't have a type for. You might get this from hg.ErrorResponse() or called the .Read() method on a hg.ResponseReader
type UnknownResponse             struct {meta string ; statusCode int}



func (receiver ResponseInput)               Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 10
func (receiver ResponseSensitiveInput)      Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 11

func (receiver ResponseRedirectTemporary)   Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 30
func (receiver ResponseRedirectPermanent)   Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 31

func (receiver ResponseTemporaryFailure)    Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 40
func (receiver ResponseServerUnavailable)   Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 41
func (receiver ResponseCGIError)            Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 42
func (receiver ResponseProxyError)          Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 43
func (receiver ResponseSlowDown)            Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 44

func (receiver ResponsePermanentFailure)    Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 50
func (receiver ResponseNotFound)            Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 51
func (receiver ResponseGone)                Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 52
func (receiver ResponseProxyRequestRefused) Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 53
func (receiver ResponseBadRequest)          Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 59

func (receiver UnknownResponse)             Error() string {return fmt.Sprintf("hg: response error — status‐code=%d meta=%q", receiver.statusCode,   receiver.meta)}



func (receiver ResponseInput)               StatusCode() int {return 10}
func (receiver ResponseSensitiveInput)      StatusCode() int {return 11}

func (receiver ResponseRedirectTemporary)   StatusCode() int {return 30}
func (receiver ResponseRedirectPermanent)   StatusCode() int {return 31}

func (receiver ResponseTemporaryFailure)    StatusCode() int {return 40}
func (receiver ResponseServerUnavailable)   StatusCode() int {return 41}
func (receiver ResponseCGIError)            StatusCode() int {return 42}
func (receiver ResponseProxyError)          StatusCode() int {return 43}
func (receiver ResponseSlowDown)            StatusCode() int {return 44}

func (receiver ResponsePermanentFailure)    StatusCode() int {return 50}
func (receiver ResponseNotFound)            StatusCode() int {return 51}
func (receiver ResponseGone)                StatusCode() int {return 52}
func (receiver ResponseProxyRequestRefused) StatusCode() int {return 53}
func (receiver ResponseBadRequest)          StatusCode() int {return 59}

func (receiver UnknownResponse)             StatusCode() int {return receiver.statusCode}



func (receiver ResponseInput)               Meta() string {return receiver.meta}
func (receiver ResponseSensitiveInput)      Meta() string {return receiver.meta}

func (receiver ResponseRedirectTemporary)   Meta() string {return receiver.meta}
func (receiver ResponseRedirectPermanent)   Meta() string {return receiver.meta}

func (receiver ResponseTemporaryFailure)    Meta() string {return receiver.meta}
func (receiver ResponseServerUnavailable)   Meta() string {return receiver.meta}
func (receiver ResponseCGIError)            Meta() string {return receiver.meta}
func (receiver ResponseProxyError)          Meta() string {return receiver.meta}
func (receiver ResponseSlowDown)            Meta() string {return receiver.meta}

func (receiver ResponsePermanentFailure)    Meta() string {return receiver.meta}
func (receiver ResponseNotFound)            Meta() string {return receiver.meta}
func (receiver ResponseGone)                Meta() string {return receiver.meta}
func (receiver ResponseProxyRequestRefused) Meta() string {return receiver.meta}
func (receiver ResponseBadRequest)          Meta() string {return receiver.meta}

func (receiver UnknownResponse)             Meta() string {return receiver.meta}



func (receiver ResponseInput)               StatusText() string {return "input"}
func (receiver ResponseSensitiveInput)      StatusText() string {return "sensitive input"}

func (receiver ResponseRedirectTemporary)   StatusText() string {return "redirect ‐ temporary"}
func (receiver ResponseRedirectPermanent)   StatusText() string {return "redirect ‐ permanent"}

func (receiver ResponseTemporaryFailure)    StatusText() string {return "temporary failure"}
func (receiver ResponseServerUnavailable)   StatusText() string {return "server unavailable"}
func (receiver ResponseCGIError)            StatusText() string {return "cgi error"}
func (receiver ResponseProxyError)          StatusText() string {return "proxy error"}
func (receiver ResponseSlowDown)            StatusText() string {return "slow down"}
func (receiver ResponseNotFound)            StatusText() string {return "not found"}
func (receiver ResponseGone)                StatusText() string {return "gone"}
func (receiver ResponseProxyRequestRefused) StatusText() string {return "proxy request refused"}
func (receiver ResponseBadRequest)          StatusText() string {return "bad request"}

func (receiver UnknownResponse)             StatusText() string {return ""}

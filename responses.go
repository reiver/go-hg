package hg

import (
	"fmt"
)

type ResponseInput               struct {meta string} // 10
type ResponseSensitiveInput      struct {meta string} // 11

type ResponseRedirectTemporary   struct {meta string} // 30
type ResponseRedirectPermanent   struct {meta string} // 31

type ResponseTemporaryFailure    struct {meta string} // 40
type ResponseServerUnavailable   struct {meta string} // 41
type ResponseCGIError            struct {meta string} // 42
type ResponseProxyError          struct {meta string} // 43
type ResponseSlowDown            struct {meta string} // 44

type ResponsePermanentFailure    struct {meta string} // 50
type ResponseNotFound            struct {meta string} // 51
type ResponseGone                struct {meta string} // 52
type ResponseProxyRequestRefused struct {meta string} // 53
type ResponseBadRequest          struct {meta string} // 59



func (receiver ResponseInput)               Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 10
func (receiver ResponseSensitiveInput)      Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 11

func (receiver ResponseRedirectTemporary)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 30
func (receiver ResponseRedirectPermanent)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 31

func (receiver ResponseTemporaryFailure)    Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 40
func (receiver ResponseServerUnavailable)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 41
func (receiver ResponseCGIError)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 42
func (receiver ResponseProxyError)          Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 43
func (receiver ResponseSlowDown)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 44

func (receiver ResponsePermanentFailure)    Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 50
func (receiver ResponseNotFound)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 51
func (receiver ResponseGone)                Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 52
func (receiver ResponseProxyRequestRefused) Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 53
func (receiver ResponseBadRequest)          Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.meta)} // 59



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

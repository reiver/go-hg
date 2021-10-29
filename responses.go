package hg

import (
	"fmt"
)

type ResponseInput               struct {Meta string} // 10
type ResponseSensitiveInput      struct {Meta string} // 11

type ResponseRedirectTemporary   struct {Meta string} // 30
type ResponseRedirectPermanent   struct {Meta string} // 31

type ResponseTemporaryFailure    struct {Meta string} // 40
type ResponseServerUnavailable   struct {Meta string} // 41
type ResponseCGIError            struct {Meta string} // 42
type ResponseProxyError          struct {Meta string} // 43
type ResponseSlowDown            struct {Meta string} // 44

type ResponsePermanentFailure    struct {Meta string} // 50
type ResponseNotFound            struct {Meta string} // 51
type ResponseGone                struct {Meta string} // 52
type ResponseProxyRequestRefused struct {Meta string} // 53
type ResponseBadRequest          struct {Meta string} // 59



func (receiver ResponseInput)               Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 10
func (receiver ResponseSensitiveInput)      Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 11

func (receiver ResponseRedirectTemporary)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 30
func (receiver ResponseRedirectPermanent)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 31

func (receiver ResponseTemporaryFailure)    Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 40
func (receiver ResponseServerUnavailable)   Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 41
func (receiver ResponseCGIError)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 42
func (receiver ResponseProxyError)          Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 43
func (receiver ResponseSlowDown)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 44

func (receiver ResponsePermanentFailure)    Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 50
func (receiver ResponseNotFound)            Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 51
func (receiver ResponseGone)                Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 52
func (receiver ResponseProxyRequestRefused) Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 53
func (receiver ResponseBadRequest)          Error() string {return fmt.Sprintf("hg: response error — status-code=%d meta=%q", receiver.StatusCode(), receiver.Meta)} // 59



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

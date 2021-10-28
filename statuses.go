package hg

const (
	StatusInput              = 10
	StatusSensitiveInput     = 11

	StatusSuccess            = 20

	StatusRedirectTemporary  = 30
	StatusRedirectPermanent  = 31

	StatusTemporaryFailure   = 40
	StatusServerUnavailable  = 41
	StatusCGIError           = 42
	StatusProxyError         = 43
	StatusSlowDown           = 44

	StatusPermanentFailure    = 50
	StatusNotFound            = 51
	StatusGone                = 52
	StatusProxyRequestRefused = 53
	StatusBadRequest          = 59
)

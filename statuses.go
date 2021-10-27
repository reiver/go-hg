package hg

const (
	StatusInput              = 10
	StatusSensitiveInput     = 11

	StatusSuccess            = 20

	StatusRedirectTemporary  = 30
	StatusRedirecPermanent   = 31

	StatusTemporaryFailure   = 40
	StatusServerUnavailable  = 41
	StatusCGIError           = 42
	StatusProxyError         = 43
	StatusSlowDown           = 44

	StatusPermanentFailture  = 50
	StatusNotFound           = 51
	StatusGone               = 52
	StatusProxyRequestFailed = 53
	StatusBadRequest         = 59
)

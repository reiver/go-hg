package hg

// Constants that can be used as default values for the Mercury Protocol's meta.
//
// Can use, for example, with ResponseWriter's WriteHeader method.
//
// For example:
//
//	hg.ServeNotFound(w, hg.DefaultMetaNotFound)
//
// Also for example:
//
//	hg.ServeTemporaryFailure(w, hg.DefaultMetaTemporaryFailure)
const (
	DefaultMetaInput               = "input"                 // 10
	DefaultMetaSensitiveInput      = "sensitive‐input"       // 11

	DefaultMetaSuccess             = "success"               // 20

	DefaultMetaTemporaryFailure    = "temporary‐failure"     // 40
	DefaultMetaServerUnavailable   = "server‐unavailable"    // 41
	DefaultMetaCGIError            = "cgi‐error"             // 42
	DefaultMetaProxyError          = "proxy‐error"           // 43
	DefaultMetaSlowDown            = "3"                     // 44

	DefaultMetaPermanentFailure    = "permanent‐failure"     // 50
	DefaultMetaNotFound            = "not‐found"             // 51
	DefaultMetaGone                = "gone"                  // 52
	DefaultMetaProxyRequestRefused = "proxy‐request‐refused" // 53
	DefaultMetaBadRequest          = "bad‐request"           // 59
)

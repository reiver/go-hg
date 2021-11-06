package hg

// These are constants that can be used as default values for a Mercury Protocol's response‐header's meta.
//
// For example usage:
//
//	hg.ServeNotFound(w, hg.DefaultMetaNotFound)
//
// Also for another example usage:
//
//	hg.ServeTemporaryFailure(w, hg.DefaultMetaTemporaryFailure)
//
// To understand these —
//
// The Mercury Protocol is based on Gemini Protocol.
// And therefore a Mercury Protocol response‐header's structure is defined in the Gemini Protocol's specification.
// In the Gemini Protocol specification, the response‐header is described as follows:
//
//	<STATUS><SPACE><META><CR><LF>
//
// In Go code, this (the Mercury Protocol's response‐header) is equivalent to:
//
//	twoDigitStatusNumericalCode + " " + meta + "\r\n"
//
// For most Mercury Protocol response types, the value of the response‐header's ‘meta’ is likely cosmetic.
// And possibly, no human will ever see them (depending on whether the client software presents them to the user or not).
// The following constants provide useful default values for these cosmetic meta value's, that can make a programmer's life easier when developing a Mercury Protcol client or server:
//
// • DefaultMetaSuccess             = "success"               // 20
//
// • DefaultMetaTemporaryFailure    = "temporary‐failure"     // 40
//
// • DefaultMetaServerUnavailable   = "server‐unavailable"    // 41
//
// • DefaultMetaCGIError            = "cgi‐error"             // 42
//
// • DefaultMetaProxyError          = "proxy‐error"           // 43
//
// • DefaultMetaPermanentFailure    = "permanent‐failure"     // 50
//
// • DefaultMetaNotFound            = "not‐found"             // 51
//
// • DefaultMetaGone                = "gone"                  // 52
//
// • DefaultMetaProxyRequestRefused = "proxy‐request‐refused" // 53
//
// • DefaultMetaBadRequest          = "bad‐request"           // 59
//
// Two of the of these default response‐header's ‘meta’ are (not cosmetic but are) shown to the user.
// The programmer SHOULD create their own message; but just in case the don't, these default values exist:
//
// • DefaultMetaInput               = "input"                 // 10
//
// • DefaultMetaSensitiveInput      = "sensitive‐input"       // 11
//
// In addition to these, one of these default response‐header is functional.
//
// DefaultMetaSlowDown            = "3"                     // 44
//
// This also SHOULD be chosen by the programmer; but again just in case the don't, a default values exists.
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

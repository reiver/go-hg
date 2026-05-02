package hg

const (
	defaultfilename = "default.gmni"
)

const (
	defaultmediatype = "application/octet-stream"
)

// maxrequest is the maximum number of bytes this package is allowing for a Mercury Protocol request — request-uri + "\r\n".
// The Gemini Protocol spec (which the Mercury Protocol is based on) says the request URI SHOULD be a maximum of 1024 bytes.
// We are allowing twice that.
const maxrequest = 1024 * 2

const maxmeta = maxrequest - 2
//                           2 == len("\r\n")

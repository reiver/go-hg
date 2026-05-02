package hg

import (
	paths "path"
)

// reqpath2fspath takes the path from the URL in a Mercury Protocol request,
// and (does some calculations and) returns the corresponding filesystem path.
//
// So, for example, if the Mercury Protocol request was:
//
//	"mercury://example.com/apple/banana/cherry\r\n"
//
// Then the URL from this request would be:
//
//	"mercury://example.com/apple/banana/cherry"
//
// And the request-path from this URL would be:
//
//	"/apple/banana/cherry"
//
// And then the filesystem-path would be:
//
//	"apple/banana/cherry"
func reqpath2fspath(reqpath string) (string, bool) {

	var p string = reqpath
	{
		p = paths.Clean(p)

		if !paths.IsAbs(p) {
			return "", false
		}

		if "/" == p {
			p = paths.Join(p, defaultfilename)
		}
	}

	var fspath string
	{
		// Strip the leading '/'
		fspath = p[1:]
	}

	return fspath, true
}

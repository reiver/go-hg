package hg

import (
	"path/filepath"
	"unicode/utf8"
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

	var path string = reqpath
	{
		path = filepath.FromSlash(path)

		path = filepath.Clean(path)

		if ! filepath.IsAbs(path) {
			return "", false
		}

		if string(filepath.Separator) == path {
			path = filepath.Join(path, defaultfilename)
		}
	}

	var fspath string
	{
		var size int = utf8.RuneLen(filepath.Separator)
		if size <= 0 {
			return "", false
		}

		fspath = path[size:]
	}

	return fspath, true
}

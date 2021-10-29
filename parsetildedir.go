package hg

import (
	paths "path"
	"strings"
)

// parsetildedir extracts the username and sub-path from a path of the from:
//
//	/~joeglow
//	/~joeglow/
//	/~joeglow/apple
//	/~joeglow/apple/
//	/~joeglow/apple/banana
//	/~joeglow/apple/banana/
//	/~joeglow/apple/banana/cherry.txt
//	/~charles
//	/~charles/once/twice/thrice/fouce
func parsetildedir(path string) (username string, subpath string, valid bool) {

	var str string = paths.Clean(path)

	{
		const prefix = "/~"
		if !strings.HasPrefix(str, prefix) {
			return "", "", false
		}

		str = str[len(prefix):]
	}

	{
		index := strings.Index(str, "/")

		switch {
		case index < 0:
			username = str
			subpath  = "/"
		default:
			username = str[:index]
			subpath = str[index:]
		}
	}

	return username, subpath, true
}

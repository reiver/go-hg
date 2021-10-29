package hg

import (
	"net/url"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"unicode/utf8"
)

// Mercury based tilde (~) capsule sites.
//
// Makes things like this:
//
//	mercury://example.com/~username/
//
// Get mapped to:
//
//	/home/username/mercury_public/default.gmni
//
// And makes things like this:
//
//	mercury://example.com/~username/once/twice/thrice/fource.txt
//
// Get mapped to:
//
//	/home/username/mercury_public/once/twice/thrice/fource.txt
const UserDirHandler internalUserDirHandler = internalUserDirHandler(0)

type internalUserDirHandler int

var _ Handler = internalUserDirHandler(0)

func (internalUserDirHandler) ServeMercury(w ResponseWriter, r Request) {

	requestValue := r.RequestValue()

	var uri *url.URL
	{
		var err error

		uri, err = url.Parse(requestValue)
		if nil != err {
			ServeBadRequest(w)
			return
		}

		if nil == uri {
			ServeTemporaryFailure(w)
			return
		}
	}

	var (
		requestpath string
		explicitDir  bool
	)
	{
		{
			r, _ := utf8.DecodeLastRuneInString(uri.Path)
			if utf8.RuneError == r {
				ServeBadRequest(w)
				return
			}

			if '/' == r {
				explicitDir = true
			}
		}

		requestpath = path.Clean(uri.Path)

		if "" == requestpath {
			ServeTemporaryFailure(w)
			return
		}
	}

	var (
		username string
		subpath  string
	)
	{
		var valid bool

		username, subpath, valid = parsetildedir(requestpath)
		if !valid {
			ServeNotFound(w)
			return
		}

		if "" == username {
			ServeTemporaryFailure(w)
			return
		}
		if "" == subpath {
			ServeTemporaryFailure(w)
			return
		}
	}

	var homedir string
	{
		u, err := user.Lookup(username)
		if nil != err {
			switch err.(type) {
			case user.UnknownUserError:
				ServeNotFound(w)
				return
			default:
				ServeTemporaryFailure(w)
				return
			}
			return
		}

		homedir = u.HomeDir

		if "" == homedir {
			ServeTemporaryFailure(w)
			return
		}
	}

	var targetpath string
	{
		const publicDir = "mercury_public"

		targetpath = path.Join(homedir, publicDir, subpath)

		targetpath = path.Clean(targetpath)

		if "" == targetpath {
			ServeTemporaryFailure(w)
			return
		}
	}

	{
		fi, err := os.Stat(targetpath)
		if nil != err {
			ServeNotFound(w)
			return
		}

		mode := fi.Mode()

		switch {
		case mode.IsDir():
			if !explicitDir {
				uri.Path += "/"
				ServeRedirectTemporary(w, uri.String())
				return
			}

			targetpath = path.Join(targetpath, "default.gmni")

		case mode.IsRegular():
			// Nothing here.
		default:
			ServeNotFound(w)
			return
		}

		var file *os.File
		{
			var err error

			file, err = os.Open(targetpath)
			if nil != err {
				ServeTemporaryFailure(w)
				return
			}
			defer func() {
				err := file.Close()
				if nil != err {
					
				}
			}()
		}

		var mediatype string
		{
			var magic [512]byte

			_, err := file.Read(magic[:])
			if nil != err {
				ServeTemporaryFailure(w)
				return
			}
			{
				_, err := file.Seek(0,0)
				if nil != err {
					ServeTemporaryFailure(w)
					return
				}
			}

			mediatype = http.DetectContentType(magic[:])

			switch mediatype {
			case "application/octet-stream":

				extension := path.Ext(targetpath)

				switch extension {
				case ".gmi", ".gmni":
					mediatype = "text/gemini"
				}
			}
		}

		w.WriteHeader(StatusSuccess, mediatype)
		io.Copy(w, file)
		return
	}
}

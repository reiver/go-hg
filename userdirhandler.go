package hg

import (
	"context"
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

func (internalUserDirHandler) ServeMercury(ctx context.Context, w ResponseWriter, r Request) {

	requestValue := r.RequestValue()

	var uri *url.URL
	{
		var err error

		uri, err = url.Parse(requestValue)
		if nil != err {
			ServeBadRequest(ctx, w)
			return
		}

		if nil == uri {
			ServeTemporaryFailure(ctx, w)
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
				ServeBadRequest(ctx, w)
				return
			}

			if '/' == r {
				explicitDir = true
			}
		}

		requestpath = path.Clean(uri.Path)

		if "" == requestpath {
			ServeTemporaryFailure(ctx, w)
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
			ServeNotFound(ctx, w)
			return
		}

		if "" == username {
			ServeTemporaryFailure(ctx, w)
			return
		}
		if "" == subpath {
			ServeTemporaryFailure(ctx, w)
			return
		}
	}

	var homedir string
	{
		u, err := user.Lookup(username)
		if nil != err {
			switch err.(type) {
			case user.UnknownUserError:
				ServeNotFound(ctx, w)
				return
			default:
				ServeTemporaryFailure(ctx, w)
				return
			}
		}

		homedir = u.HomeDir

		if "" == homedir {
			ServeTemporaryFailure(ctx, w)
			return
		}
	}

	var targetpath string
	{
		const publicDir = "mercury_public"

		targetpath = path.Join(homedir, publicDir, subpath)

		targetpath = path.Clean(targetpath)

		if "" == targetpath {
			ServeTemporaryFailure(ctx, w)
			return
		}
	}

	{
		fi, err := os.Stat(targetpath)
		if nil != err {
			ServeNotFound(ctx, w)
			return
		}

		mode := fi.Mode()

		switch {
		case mode.IsDir():
			if !explicitDir {
				uri.Path += "/"
				ServeRedirectTemporary(ctx, w, uri.String())
				return
			}

			targetpath = path.Join(targetpath, defaultfilename)

		case mode.IsRegular():
			// Nothing here.
		default:
			ServeNotFound(ctx, w)
			return
		}

		var file *os.File
		{
			var err error

			file, err = os.Open(targetpath)
			if nil != err {
				ServeTemporaryFailure(ctx, w)
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
				ServeTemporaryFailure(ctx, w)
				return
			}
			{
				_, err := file.Seek(0,0)
				if nil != err {
					ServeTemporaryFailure(ctx, w)
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

		w.WriteHeader(ctx, StatusSuccess, mediatype)
		io.Copy(w.Writer(ctx), file)
		return
	}
}

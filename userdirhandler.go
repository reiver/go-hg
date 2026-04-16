package hg

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"codeberg.org/reiver/go-field"
)

// UserDirHandler is a Mercury Protocol handler for tilde (~) capsule sites.
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
type UserDirHandler struct {
	// AllowSymLinks controls whether symbolic links inside a user's mercury_public
	// directory are followed. If false (the default), requests that resolve to a
	// symlink or pass through a symlink are rejected with a not-found response.
	AllowSymLinks bool

	Logger Logger
}

var _ Handler = &UserDirHandler{}

func (receiver *UserDirHandler) ServeMercury(ctx context.Context, w ResponseWriter, r Request) {
	if nil == receiver {
		return
	}

	log := mustlogger(receiver.Logger).Begin()
	defer log.End()

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

	const publicDir = "mercury_public"

	var targetpath string
	{
		targetpath = filepath.Join(homedir, publicDir, subpath)

		targetpath = filepath.Clean(targetpath)

		if "" == targetpath {
			ServeTemporaryFailure(ctx, w)
			return
		}
	}

	// If symlinks are not allowed, resolve the real path and verify it is still
	// under the user's mercury_public directory.
	var allowedPrefix string
	if nil == receiver || !receiver.AllowSymLinks {
		allowedPrefix = filepath.Join(homedir, publicDir) + string(filepath.Separator)
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

			targetpath = filepath.Join(targetpath, defaultfilename)

		case mode.IsRegular():
			// Nothing here.
		default:
			ServeNotFound(ctx, w)
			return
		}

		if "" != allowedPrefix {
			resolved, err := filepath.EvalSymlinks(targetpath)
			if nil != err {
				ServeNotFound(ctx, w)
				return
			}
			if !strings.HasPrefix(resolved, allowedPrefix) {
				ServeNotFound(ctx, w)
				return
			}
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
					log.Error(
						field.S("could not close file"),
						field.String("path", targetpath),
						field.Stringer("request", r),
						field.E(err),
					)
				}
			}()
		}

		var mediatype string
		{
			var magic [512]byte

			n, err := file.Read(magic[:])
			if nil != err && err != io.EOF {
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

			mediatype = http.DetectContentType(magic[:n])

			switch mediatype {
			case "application/octet-stream":

				extension := filepath.Ext(targetpath)

				switch extension {
				case ".gmi", ".gmni", ".gemini":
					mediatype = "text/gemini"
				}
			}
		}

		if _, headerErr := w.WriteHeader(ctx, StatusSuccess, mediatype); nil != headerErr {
			log.Error(
				field.S("problem writing Mercury Protocol header"),
				field.Stringer("request", r),
				field.E(headerErr),
			)
			return
		}

		if _, copyErr := io.Copy(w.Writer(ctx), file); nil != copyErr {
			log.Error(
				field.S("problem writing Mercury Protocol body by copying file inferred from request"),
				field.String("path", targetpath),
				field.Stringer("request", r),
				field.E(copyErr),
			)
		}
		return
	}
}

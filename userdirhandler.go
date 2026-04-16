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
			if serveErr := ServeBadRequest(ctx, w); nil != serveErr {
				log.Error(
					field.S("problem sending bad-request response"),
					field.Stringer("request", r),
					field.E(serveErr),
				)
			}
			return
		}

		if nil == uri {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
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
				if err := ServeBadRequest(ctx, w); nil != err {
					log.Error(
						field.S("problem sending bad-request response"),
						field.String("request-value", requestValue),
						field.E(err),
					)
				}
				return
			}

			if '/' == r {
				explicitDir = true
			}
		}

		requestpath = path.Clean(uri.Path)

		if "" == requestpath {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
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
			if err := ServeNotFound(ctx, w); nil != err {
				log.Error(
					field.S("problem sending not-found response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		if "" == username {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}
		if "" == subpath {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}
	}

	var homedir string
	{
		u, err := user.Lookup(username)
		if nil != err {
			switch err.(type) {
			case user.UnknownUserError:
				if serveErr := ServeNotFound(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending not-found response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
				return
			default:
				if serveErr := ServeTemporaryFailure(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending temporary-failure response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
				return
			}
		}

		homedir = u.HomeDir

		if "" == homedir {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}
	}

	const publicDir = "mercury_public"

	var targetpath string
	{
		targetpath = filepath.Join(homedir, publicDir, subpath)

		targetpath = filepath.Clean(targetpath)

		if "" == targetpath {
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
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
			if serveErr := ServeNotFound(ctx, w); nil != serveErr {
				log.Error(
					field.S("problem sending not-found response"),
					field.Stringer("request", r),
					field.E(serveErr),
				)
			}
			return
		}

		mode := fi.Mode()

		switch {
		case mode.IsDir():
			if !explicitDir {
				uri.Path += "/"
				if err := ServeRedirectTemporary(ctx, w, uri.String()); nil != err {
					log.Error(
						field.S("problem sending redirect-temporary response"),
						field.Stringer("request", r),
						field.E(err),
					)
				}
				return
			}

			targetpath = filepath.Join(targetpath, defaultfilename)

		case mode.IsRegular():
			// Nothing here.
		default:
			if err := ServeNotFound(ctx, w); nil != err {
				log.Error(
					field.S("problem sending not-found response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		// Symlink check: resolve the real path and verify it stays within the
		// user's mercury_public directory.
		//
		// Note: there is an inherent TOCTOU (time-of-check-time-of-use) race here —
		// the filesystem could change between EvalSymlinks and the os.Open below.
		// The correct fix would be to open the file first, then verify the resolved
		// path of the already-opened file descriptor (e.g., via fstat + readlink on
		// /proc/self/fd). However, Go's standard library does not provide a way to
		// resolve symlinks from an open *os.File — filepath.EvalSymlinks operates on
		// paths, not file descriptors. The race window is small and exploitation
		// requires write access to the user's mercury_public directory, which already
		// grants the ability to serve arbitrary content.
		if "" != allowedPrefix {
			resolved, err := filepath.EvalSymlinks(targetpath)
			if nil != err {
				if serveErr := ServeNotFound(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending not-found response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
				return
			}
			if !strings.HasPrefix(resolved, allowedPrefix) {
				if serveErr := ServeNotFound(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending not-found response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
				return
			}
		}

		var file *os.File
		{
			var err error

			file, err = os.Open(targetpath)
			if nil != err {
				if serveErr := ServeTemporaryFailure(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending temporary-failure response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
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
				if serveErr := ServeTemporaryFailure(ctx, w); nil != serveErr {
					log.Error(
						field.S("problem sending temporary-failure response"),
						field.Stringer("request", r),
						field.E(serveErr),
					)
				}
				return
			}
			{
				_, err := file.Seek(0,0)
				if nil != err {
					if serveErr := ServeTemporaryFailure(ctx, w); nil != serveErr {
						log.Error(
							field.S("problem sending temporary-failure response"),
							field.Stringer("request", r),
							field.E(serveErr),
						)
					}
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

package hg

import (
	"context"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"

	"codeberg.org/reiver/go-field"
)

// FileSystemHandler is used to create a Mercury Protocol server that serves files from a fs.FS file system.
//
// For example usage:
//
//	var fshandler hg.FileSystemHandler
//	fshandler.Root = os.DirFS("/path/to/mercury/root")
//	
//	var handler hg.Handler = &fshandler
//	
//	err := hg.ListenAndServe(":1961", handler)
type FileSystemHandler struct {
	Root fs.FS
	Logger Logger
}

func (receiver FileSystemHandler) ServeMercury(ctx context.Context, w ResponseWriter, r Request) {

	log := mustlogger(receiver.Logger).Begin()
	defer log.End()

	var root fs.FS
	{
		root = receiver.Root

		if nil == root {
			log.Error(field.S("nil root fs.FS (file system)"))
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		log.Trace(
			field.S("have root fs.FS (file system)"),
			field.FormattedString("file-system-type", "%T", root),
		)
	}

	var requestValue string
	{
		requestValue = r.RequestValue()

		log.Trace(
			field.String("request-value", requestValue),
		)
	}

	var uri *url.URL
	{
		var err error

		uri, err = url.Parse(requestValue)
		if nil != err {
			log.Error(
				field.S("could not parse request-value as URL"),
				field.String("request-value", requestValue),
				field.E(err),
			)
			if err := ServeBadRequest(ctx, w, "could not parse URL"); nil != err {
				log.Error(
					field.S("problem sending bad-request response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		log.Trace(
			field.S("URL parsed"),
			field.String("request-value", requestValue),
		)
	}

	{
		actualScheme := uri.Scheme

		log.Trace(
			field.String("expected-URL-scheme", Scheme),
			field.String("expected-URL-scheme-TLS", SchemeTLS),
			field.String("actual-URL-scheme", actualScheme),
			field.Stringer("uri", uri),
		)

		if Scheme != actualScheme && SchemeTLS != actualScheme {
			log.Error(
				field.S("the actual scheme in the URL from the request is not what was expected"),
				field.String("expected-URL-scheme", Scheme),
				field.String("expected-URL-scheme-TLS", SchemeTLS),
				field.String("actual-scheme", actualScheme),
				field.Stringer("uri", uri),
			)
			if err := ServeBadRequest(ctx, w, "unsupported URL scheme"); nil != err {
				log.Error(
					field.S("problem sending bad-request response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		log.Trace(
			field.S("URL scheme accepted"),
			field.String("scheme", actualScheme),
			field.Stringer("uri", uri),
		)
	}

	var reqpath string
	{
		reqpath = uri.Path

		if "" == reqpath {
			log.Error(
				field.S("the path from the request is empty"),
				field.String("path", reqpath),
				field.Stringer("request", r),
			)
			if err := ServeBadRequest(ctx, w, "empty request path"); nil != err {
				log.Error(
					field.S("problem sending bad-request response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}

		log.Trace(
			field.String("request-path", reqpath),
		)
	}

	var fspath string
	var file fs.File
	{
		var valid bool

		fspath, valid = reqpath2fspath(reqpath)
		if !valid {
			log.Error(
				field.S("request-path is invalid."),
				field.String("request-path", reqpath),
			)
			if err := ServeBadRequest(ctx, w, "invalid URL path"); nil != err {
				log.Error(
					field.S("problem sending bad-request response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}
		log.Trace(
			field.String("filesystem-path", fspath),
		)

		{
			var err error
			file, err = root.Open(fspath)
			if nil != err {
				var notfound bool = os.IsNotExist(err)
				log.Error(
					field.S("could not open file"),
					field.FormattedString("error-type", "%T", err),
					field.E(err),
				)
				log.Error(
					field.Bool("does-not-exists", notfound),
				)

				switch {
				case notfound:
					if err := ServeNotFound(ctx, w); nil != err {
						log.Error(
							field.S("problem sending not-found response"),
							field.Stringer("request", r),
							field.E(err),
						)
					}
					return
				default:
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
			if nil == file || fs.File(nil) == file {
				log.Error(
					field.S("file is nil"),
					field.FormattedString("file", "%#v", file),
				)
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
		log.Trace(
			field.S("file opened"),
			field.String("path", fspath),
		)
		defer func(f fs.File, path string) {
			err := f.Close()
			if nil != err {
				log.Error(
					field.S("could not close path from request"),
					field.String("path", path),
					field.Stringer("request", r),
					field.E(err),
				)
			}
		}(file, fspath)

		fileinfo, err := file.Stat()
		if nil != err {
			log.Error(
				field.S("could not get the file-status of the filesystem-path from the request:"),
				field.String("filesystem-path", fspath),
				field.Stringer("request", r),
				field.E(err),
			)
			if err := ServeTemporaryFailure(ctx, w); nil != err {
				log.Error(
					field.S("problem sending temporary-failure response"),
					field.Stringer("request", r),
					field.E(err),
				)
			}
			return
		}
		log.Trace(field.FormattedString("fileinfo", "%#v", fileinfo))

		if fileinfo.IsDir() {
			log.Trace(
				field.S("filesystem-path is a directory"),
				field.String("filesystem-path", fspath),
			)

			defaultpath := path.Join(fspath, defaultfilename)

			defaultfile, err := root.Open(defaultpath)
			if nil != err {
				var notfound bool = os.IsNotExist(err)
				log.Error(
					field.S("could not open defaulted file"),
					field.String("default-path", defaultpath),
					field.FormattedString("error-type", "%T", err),
					field.E(err),
				)
				log.Error(
					field.Bool("does-not-exists", notfound),
				)

				switch {
				case notfound:
					if err := ServeNotFound(ctx, w); nil != err {
						log.Error(
							field.S("problem sending not-found response"),
							field.Stringer("request", r),
							field.E(err),
						)
					}
					return
				default:
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
			defer func(f fs.File, path string) {
				err := f.Close()
				if nil != err {
					log.Error(
						field.S("could not close default file from request"),
						field.String("default-path", path),
						field.Stringer("request", r),
						field.E(err),
					)
				}
			}(defaultfile, defaultpath)

			file = defaultfile
			fspath = defaultpath
		} else {
			log.Trace(
				field.S("filesystem-path s NOT a directory"),
				field.String("filesystem-path", fspath),
			)
		}
	}

	var mediatype string
	{
		mediatype = infermediatype(fspath)

		log.Trace(field.String("media-type", mediatype))
	}

	{
		_, err := w.WriteHeader(ctx, StatusSuccess, mediatype)
		if nil != err {
			log.Error(
				field.S("problem writing Mercury Protocol header"),
				field.E(err),
			)
			return
		}
		log.Trace(field.S("Mercury Protocol header written"))
	}

	{
		log.Trace(
			field.S("BEGIN copying"),
			field.String("filesystem-path", fspath),
		)
		log.Trace(
			field.FormattedString("response-writer-type", "%T", w),
			field.FormattedString("response-writer", "%#v", w),
			field.FormattedString("file-type", "%T", file),
			field.FormattedString("file", "%#v", file),
		)

		_, err := io.Copy(w.Writer(ctx), file)
		if nil != err{
			log.Error(
				field.S("could not copy file inferred from request"),
				field.Stringer("request", r),
				field.E(err),
			)
			return
		}
		log.Trace(
			field.S("END copying"),
			field.String("filesystem-path", fspath),
		)

		log.Debug(field.F("%q %q", r, fspath))
	}
}

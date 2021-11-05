package hg

import (
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
)

type FileSystemHandler struct {
	Root fs.FS
	Logger Logger
}

func (receiver FileSystemHandler) ServeMercury(w ResponseWriter, r Request) {

	var logger Logger = mustlogger(receiver.Logger)

	logger.Trace("hg.FileSystemHandler.ServeMercury: BEGIN")
	defer logger.Trace("hg.FileSystemHandler.ServeMercury: END")

	var root fs.FS
	{
		root = receiver.Root

		if nil == root {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: nil root fs.FS (file system)")
			ServeTemporaryFailure(w)
			return
		}

		logger.Tracef("hg.FileSystemHandler.ServeMercury: have root fs.FS (file system) of type %T", root)
	}

	var requestValue string
	{
		requestValue = r.RequestValue()

		logger.Tracef("hg.FileSystemHandler.ServeMercury: request-value=%q", requestValue)
	}

	var uri *url.URL
	{
		var err error

		uri, err = url.Parse(requestValue)
		if nil != err {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: could not parse request-value %q as URL: %s", requestValue, err)
			ServeBadRequest(w, "could not parse URL")
			return
		}

		logger.Trace("hg.FileSystemHandler.ServeMercury: URL parsed")
	}

	{
		const expectedScheme = "mercury"
		actualScheme := uri.Scheme

		logger.Tracef("hg.FileSystemHandler.ServeMercury: expected-URL-scheme=%q actual-URL-scheme=%q", expectedScheme, actualScheme)

		if expectedScheme != actualScheme  {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: the actual scheme in the URL from the request (%q) is not what was expected — expected=%q actual=%q", r, expectedScheme, actualScheme)
			ServeBadRequest(w, "unsupported URL scheme")
			return
		}

		logger.Trace("hg.FileSystemHandler.ServeMercury: URL scheme accepted")
	}

	var reqpath string
	{
		reqpath = uri.Path

		if "" == reqpath {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: the path (%q) from the request (%q) is empty", reqpath, r)
			ServeBadRequest(w, "empty request path")
			return
		}

		logger.Tracef("hg.FileSystemHandler.ServeMercury: request-path=%q", reqpath)
	}

	var fspath string
	var file fs.File
	{
		var valid bool

		fspath, valid = reqpath2fspath(reqpath)
		if !valid {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: request-path (%q) is invalid.", reqpath)
			ServeBadRequest(w, "invalid URL path")
			return
		}
		logger.Tracef("hg.FileSystemHandler.ServeMercury: filesystem-path=%q", fspath)

		{
			var err error
			file, err = root.Open(fspath)
			if nil != err {
				var notfound bool = os.IsNotExist(err)
				logger.Errorf("hg.FileSystemHandler.ServeMercury: could not open file: (%T) %s", err, err)
				logger.Errorf("hg.FileSystemHandler.ServeMercury: does-not-exists=%t: ", notfound)

				switch {
				case notfound:
					ServeNotFound(w)
					return
				default:
					ServeTemporaryFailure(w)
					return
				}
				return
			}
			if nil == file || fs.File(nil) == file {
				logger.Errorf("hg.FileSystemHandler.ServeMercury: file is nil: %#v", file)
				ServeTemporaryFailure(w)
				return
			}
		}
		logger.Tracef("hg.FileSystemHandler.ServeMercury: file %q opened", fspath)
		defer func(f fs.File, path string) {
			err := f.Close()
			if nil != err {
				logger.Errorf("hg.FileSystemHandler.ServeMercury: could not close path (%q) from request (%q): %s", path, r, err)
			}
		}(file, fspath)

		fileinfo, err := file.Stat()
		if nil != err {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: could not get the file-status of the filesystem-path (%q) from the request (%q): %s", fspath, r, err)
			ServeTemporaryFailure(w)
			return
		}
		logger.Tracef("hg.FileSystemHandler.ServeMercury: fileinfo=%v", fileinfo)

		if fileinfo.IsDir() {
			logger.Tracef("hg.FileSystemHandler.ServeMercury: filesystem-path %q is a directory", fspath)

			defaultpath := filepath.Join(fspath, defaultfilename)

			defaultfile, err := root.Open(defaultpath)
			if nil != err {
				var notfound bool = os.IsNotExist(err)
				logger.Errorf("hg.FileSystemHandler.ServeMercury: could not open defaulted file: (%T) %s", err, err)
				logger.Errorf("hg.FileSystemHandler.ServeMercury: does-not-exists=%t: ", notfound)

				switch {
				case notfound:
					ServeNotFound(w)
					return
				default:
					ServeTemporaryFailure(w)
					return
				}
				return
			}
			defer func() {
				err := file.Close()
				if nil != err {
					logger.Errorf("hg.FileSystemHandler.ServeMercury: could not close default file (%q) from request (%q): %s", defaultpath, r, err)
				}
			}()

			file = defaultfile
			fspath = defaultpath
		} else {
			logger.Tracef("hg.FileSystemHandler.ServeMercury: filesystem-path %q is NOT a directory", fspath)
		}
	}

	var mediatype string
	{
		mediatype = infermediatype(fspath)

		logger.Tracef("hg.FileSystemHandler.ServeMercury: media-type = %q", mediatype)
	}

	{
		_, err := w.WriteHeader(StatusSuccess, mediatype)
		if nil != err {
			logger.Errorf("hg.FileSystemHandler.ServeMercury: problem writing Mercury Protocol header: %s", err)
			// intentially not returning here.
		}
		logger.Trace("hg.FileSystemHandler.ServeMercury: Mercury Protocol header written")
	}

	{
		logger.Tracef("hg.FileSystemHandler.ServeMercury: filesystem-path %q BEGIN copying", fspath)
		logger.Tracef("hg.FileSystemHandler.ServeMercury: w = (%T) %#v", w, w)
		logger.Tracef("hg.FileSystemHandler.ServeMercury: file = (%T) %#v", file, file)
		_, err := io.Copy(w, file)
		if nil != err{
			logger.Errorf("hg.FileSystemHandler.ServeMercury: could not copy file inferred from request (%q): %s", r,  err)
			return
		}
		logger.Tracef("hg.FileSystemHandler.ServeMercury: filesystem-path %q END copying", fspath)

		logger.Log("%q %q",r, fspath)
	}
}

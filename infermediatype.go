package hg

import (
	"mime"
	"path/filepath"
)

func init() {
	mime.AddExtensionType(".gmi", "text/gemini")
	mime.AddExtensionType(".gmni", "text/gemini")
	mime.AddExtensionType(".gemini", "text/gemini")

	mime.AddExtensionType(".text", "text/plain")
}

// infermediatype infers the media-type from the filename or filesystem-path.
//
// If it cannot infer it, then it returns an empty string ("").
func infermediatype(s string) string {

	var fileextension string
	{
		fileextension  = filepath.Ext(s)
	}

	{
		if "" == fileextension {
			return defaultmediatype
		}
	}

	var mediatype string
	{
		mediatype = mime.TypeByExtension(fileextension)
	}

	if "" == mediatype {
		return defaultmediatype
	}

	return mediatype
}

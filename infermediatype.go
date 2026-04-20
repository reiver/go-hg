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
func infermediatype(fileName string) string {

	var fileExtension string
	{
		fileExtension  = filepath.Ext(fileName)
	}

	{
		if "" == fileExtension {
			return defaultmediatype
		}
	}

	var mediatype string
	{
		mediatype = mime.TypeByExtension(fileExtension)
	}

	//@TODO: should we had a http.DetectContentType() check if the fileExtension check fails?

	if "" == mediatype {
		return defaultmediatype
	}

	return mediatype
}

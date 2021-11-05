package hg

import (
	"testing"
)

func TestInfermediatype(t *testing.T) {

	tests := []struct{
		Path string
		Expected string
	}{
		{
			"one.txt",
			"text/plain; charset=utf-8",
		},
		{
			"one/two.txt",
			"text/plain; charset=utf-8",
		},
		{
			"one/two/three.txt",
			"text/plain; charset=utf-8",
		},



		{
			"one.text",
			"text/plain; charset=utf-8",
		},
		{
			"one/two.text",
			"text/plain; charset=utf-8",
		},
		{
			"one/two/three.text",
			"text/plain; charset=utf-8",
		},



		{
			"one.htm",
			"text/html; charset=utf-8",
		},
		{
			"one/two.htm",
			"text/html; charset=utf-8",
		},
		{
			"one/two/three.htm",
			"text/html; charset=utf-8",
		},



		{
			"one.html",
			"text/html; charset=utf-8",
		},
		{
			"one/two.html",
			"text/html; charset=utf-8",
		},
		{
			"one/two/three.html",
			"text/html; charset=utf-8",
		},



		{
			"one.gmi",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two.gmi",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two/three.gmi",
			"text/gemini; charset=utf-8",
		},



		{
			"one.gmni",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two.gmni",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two/three.gmni",
			"text/gemini; charset=utf-8",
		},



		{
			"one.gemini",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two.gemini",
			"text/gemini; charset=utf-8",
		},
		{
			"one/two/three.gemini",
			"text/gemini; charset=utf-8",
		},



		{
			"logo.png",
			"image/png",
		},
		{
			"to/logo.png",
			"image/png",
		},
		{
			"path/to/logo.png",
			"image/png",
		},
		{
			"the/path/to/logo.png",
			"image/png",
		},



		{
			"photo.jpeg",
			"image/jpeg",
		},
		{
			"images/photo.jpeg",
			"image/jpeg",
		},
		{
			"profile/images/photo.jpeg",
			"image/jpeg",
		},



		{
			"profile/images/photo.qwertyasdfgh",
			"application/octet-stream",
		},
		{
			"README",
			"application/octet-stream",
		},
	}

	for testNumber, test := range tests {

		actual := infermediatype(test.Path)

		if expected := test.Expected; expected != actual {
			t.Errorf("For test #%d, the actual inferred media-type is not what was expeceted.", testNumber)
			t.Logf("PATH: %q", test.Path)
			t.Logf("EXPECTED MEDIA-TYPE: %q", expected)
			t.Logf("ACTUAL   MEDIA-TYPE: %q", actual)
			continue
		}
	}
}

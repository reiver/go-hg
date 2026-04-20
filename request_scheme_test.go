package hg_test

import (
	"testing"

	"github.com/reiver/go-hg"
)

func TestRequest_Scheme(t *testing.T) {
	tests := []struct{
		Link     string
		Expected string
	}{
		{},



		{
			Link:     "mercury://example.com",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com",
			Expected: "mercury",
		},



		{
			Link:     "mercury://example.com/",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com/",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com/",
			Expected: "mercury",
		},



		{
			Link:     "mercury://example.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},



		{
			Link:     "mercury://example.com:10101",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com:10101",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com:10101",
			Expected: "mercury",
		},



		{
			Link:     "mercury://example.com:10101/",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com:10101/",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com:10101/",
			Expected: "mercury",
		},



		{
			Link:     "mercury://example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},
		{
			Link:     "Mercury://Example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},
		{
			Link:     "MERCURY://EXAMPLE.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "mercury",
		},



		{
			Link:     "mercury://😈.example",
			Expected: "mercury",
		},



		{
			Link:     "gemini://example.com",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com",
			Expected: "gemini",
		},



		{
			Link:     "gemini://example.com/",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com/",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com/",
			Expected: "gemini",
		},



		{
			Link:     "gemini://example.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},



		{
			Link:     "gemini://example.com:10101",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com:10101",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com:10101",
			Expected: "gemini",
		},



		{
			Link:     "gemini://example.com:10101/",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com:10101/",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com:10101/",
			Expected: "gemini",
		},



		{
			Link:     "gemini://example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},
		{
			Link:     "Gemini://Example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},
		{
			Link:     "GEMINI://EXAMPLE.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			Expected: "gemini",
		},



		{
			Link:     "gemini://😈.example",
			Expected: "gemini",
		},



		{
			Link:     "http://host.example",
			Expected: "http",
		},
		{
			Link:     "http://host.example/",
			Expected: "http",
		},
		{
			Link:     "http://host.example/once/twice/thrice/fource.html",
			Expected: "http",
		},

		{
			Link:     "Http://host.example",
			Expected: "http",
		},
		{
			Link:     "Http://host.example/",
			Expected: "http",
		},
		{
			Link:     "Http://host.example/once/twice/thrice/fource.html",
			Expected: "http",
		},

		{
			Link:     "HTTP://host.example",
			Expected: "http",
		},
		{
			Link:     "HTTP://host.example/",
			Expected: "http",
		},
		{
			Link:     "HTTP://host.example/once/twice/thrice/fource.html",
			Expected: "http",
		},
	}

	for testNumber, test := range tests {

		var request hg.Request

		err := request.Parse(test.Link)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
			t.Logf("ERROR: %s", err)
			t.Logf("LINK: %q", test.Link)
			continue
		}

		actual := request.Scheme()

		expected := test.Expected

		if expected != actual {
			t.Errorf("For test #%d, the actual 'scheme' value is not what was expected.", testNumber)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			t.Logf("LINK: %q", test.Link)
			continue
		}
	}
}

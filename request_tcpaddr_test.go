package hg_test

import (
	"testing"

	"github.com/reiver/go-hg"
)

func TestRequest_TCPAddr(t *testing.T) {
	tests := []struct{
		Link            string
		ExpectedTCPAddr string
		ExpectedFound   bool
	}{
		{},



		{
			Link:  "mercury://example.com",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://example.com/",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com/",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com/",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://example.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1961",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://example.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://example.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "Mercury://Example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:  "MERCURY://EXAMPLE.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:  "mercury://😈.example",
			ExpectedTCPAddr: "xn--m28h.example:1961",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com/",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com/",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com/",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:1965",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com:10101",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com:10101/",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link  : "Gemini://Example.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},
		{
			Link:   "GEMINI://EXAMPLE.com:10101/apple/BANANA/Cherry/dAtE.gmni",
			ExpectedTCPAddr: "example.com:10101",
			ExpectedFound: true,
		},



		{
			Link:   "gemini://😈.example",
			ExpectedTCPAddr: "xn--m28h.example:1965",
			ExpectedFound: true,
		},



		{
			Link:   "http://host.example",
			ExpectedTCPAddr: "",
			ExpectedFound: false,
		},
		{
			Link:   "http://host.example/",
			ExpectedTCPAddr: "",
			ExpectedFound: false,
		},
		{
			Link:   "http://host.example/once/twice/thrice/fource.html",
			ExpectedTCPAddr: "",
			ExpectedFound: false,
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

		actualTCPAddr, actualFound := request.TCPAddr()

		{
			actual := actualFound
			expected := test.ExpectedFound

			if expected != actual {
				t.Errorf("For test #%d, the actual 'found' value is not what was expected.", testNumber)
				t.Logf("EXPECTED: %t", expected)
				t.Logf("ACTUAL:   %t", actual)
				t.Logf("LINK: %q", test.Link)
				continue
			}
		}

		{
			actual := actualTCPAddr
			expected := test.ExpectedTCPAddr

			if expected != actual {
				t.Errorf("For test #%d, the actual 'TCP-address' value is not what was expected.", testNumber)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("LINK: %q", test.Link)
				continue
			}
		}
	}
}

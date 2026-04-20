package hg_test

import (
	"github.com/reiver/go-hg"

	"strings"

	"testing"
)

func TestRequest_Parse_success(t *testing.T) {

	tests := []struct{
		Src interface{}
		Expected string
	}{
		{
			Src:      "",
			Expected: "\r\n",
		},









		{
			Src:      "hello-world",
			Expected: "hello-world\r\n",
		},
		{
			Src:      "hello-world\r\n",
			Expected: "hello-world\r\n",
		},
		{
			Src:      "hello-world\r\nsecond-line",
			Expected: "hello-world\r\n",
		},
		{
			Src:      "hello-world\r\nsecond-line\r\n",
			Expected: "hello-world\r\n",
		},



		{
			Src: strings.NewReader("hello-world"),
			Expected:              "hello-world\r\n",
		},
		{
			Src: strings.NewReader("hello-world\r\n"),
			Expected:              "hello-world\r\n",
		},
		{
			Src: strings.NewReader("hello-world\r\nsecond-line"),
			Expected:              "hello-world\r\n",
		},
		{
			Src: strings.NewReader("hello-world\r\nsecond-line\r\n"),
			Expected:              "hello-world\r\n",
		},



		{
			Src: []byte("hello-world"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []byte("hello-world\r\n"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []byte("hello-world\r\nsecond-line"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []byte("hello-world\r\nsecond-line\r\n"),
			Expected:   "hello-world\r\n",
		},



		{
			Src: []rune("hello-world"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []rune("hello-world\r\n"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []rune("hello-world\r\nsecond-line"),
			Expected:   "hello-world\r\n",
		},
		{
			Src: []rune("hello-world\r\nsecond-line\r\n"),
			Expected:   "hello-world\r\n",
		},









		{
			Src:      "mercury://something/apple/banana/cherry.txt?query",
			Expected: "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src:      "mercury://something/apple/banana/cherry.txt?query\r\n",
			Expected: "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src:      "mercury://something/apple/banana/cherry.txt?query\r\nsecond-line",
			Expected: "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src:      "mercury://something/apple/banana/cherry.txt?query\r\nsecond-line\r\n",
			Expected: "mercury://something/apple/banana/cherry.txt?query\r\n",
		},



		{
			Src: strings.NewReader("mercury://something/apple/banana/cherry.txt?query"),
			Expected:              "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: strings.NewReader("mercury://something/apple/banana/cherry.txt?query\r\n"),
			Expected:              "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: strings.NewReader("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line"),
			Expected:              "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: strings.NewReader("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line\r\n"),
			Expected:              "mercury://something/apple/banana/cherry.txt?query\r\n",
		},



		{
			Src: []byte("mercury://something/apple/banana/cherry.txt?query"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []byte("mercury://something/apple/banana/cherry.txt?query\r\n"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []byte("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []byte("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line\r\n"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},



		{
			Src: []rune("mercury://something/apple/banana/cherry.txt?query"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []rune("mercury://something/apple/banana/cherry.txt?query\r\n"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []rune("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},
		{
			Src: []rune("mercury://something/apple/banana/cherry.txt?query\r\nsecond-line\r\n"),
			Expected:   "mercury://something/apple/banana/cherry.txt?query\r\n",
		},









		{
			Src:      "۰۱۲۳۴۵۶۷۸۹",
			Expected: "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src:      "۰۱۲۳۴۵۶۷۸۹\r\n",
			Expected: "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src:      "۰۱۲۳۴۵۶۷۸۹\r\nsecond-line",
			Expected: "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src:      "۰۱۲۳۴۵۶۷۸۹\r\nsecond-line\r\n",
			Expected: "۰۱۲۳۴۵۶۷۸۹\r\n",
		},



		{
			Src: strings.NewReader("۰۱۲۳۴۵۶۷۸۹"),
			Expected:              "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: strings.NewReader("۰۱۲۳۴۵۶۷۸۹\r\n"),
			Expected:              "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: strings.NewReader("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line"),
			Expected:              "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: strings.NewReader("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line\r\n"),
			Expected:              "۰۱۲۳۴۵۶۷۸۹\r\n",
		},



		{
			Src: []byte("۰۱۲۳۴۵۶۷۸۹"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []byte("۰۱۲۳۴۵۶۷۸۹\r\n"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []byte("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []byte("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line\r\n"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},



		{
			Src: []rune("۰۱۲۳۴۵۶۷۸۹"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []rune("۰۱۲۳۴۵۶۷۸۹\r\n"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []rune("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
		{
			Src: []rune("۰۱۲۳۴۵۶۷۸۹\r\nsecond-line\r\n"),
			Expected:   "۰۱۲۳۴۵۶۷۸۹\r\n",
		},
	}

	for testNumber, test := range tests {

		var request hg.Request

		err := request.Parse(test.Src)
		if nil != err {
			t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
			t.Logf("ERROR: (%T) %s", err, err)
			t.Logf("SRC:      %#v", test.Src)
			t.Logf("EXPECTED: %q", test.Expected)
			continue
		}

		if expected, actual := test.Expected, request.String(); expected != actual {
			t.Errorf("For test #%d, the actual value was not what was expected.", testNumber)
			t.Logf("SRC:      %q", test.Src)
			t.Logf("EXPECTED: %q", expected)
			t.Logf("ACTUAL:   %q", actual)
			continue
		}
	}
}

func TestRequest_Parse_failure(t *testing.T) {

	tests := []struct{
		Src interface{}
	}{
		{
			Src: "\r",
		},
		{
			Src: "hello-world\r",
		},
		{
			Src: "mercury://something/apple/banana/cherry.txt?query\r",
		},
		{
			Src: "۰۱۲۳۴۵۶۷۸۹\r",
		},



		{
			Src: strings.NewReader("\r"),
		},
		{
			Src: strings.NewReader("hello-world\r"),
		},
		{
			Src: strings.NewReader("mercury://something/apple/banana/cherry.txt?query\r"),
		},
		{
			Src: strings.NewReader("۰۱۲۳۴۵۶۷۸۹\r"),
		},



		{
			Src: []byte("\r"),
		},
		{
			Src: []byte("hello-world\r"),
		},
		{
			Src: []byte("mercury://something/apple/banana/cherry.txt?query\r"),
		},
		{
			Src: []byte("۰۱۲۳۴۵۶۷۸۹\r"),
		},


		{
			Src: []rune("\r"),
		},
		{
			Src: []rune("hello-world\r"),
		},
		{
			Src: []rune("mercury://something/apple/banana/cherry.txt?query\r"),
		},
		{
			Src: []rune("۰۱۲۳۴۵۶۷۸۹\r"),
		},








		{
			Src: 12.34,
		},
		{
			Src: true,
		},
		{
			Src: struct{
				X int64
				Y int64
			}{
				X:0,
				Y:1,
			},
		},
	}

	for testNumber, test := range tests {

		var request hg.Request

		err := request.Parse(test.Src)
		if nil == err {
			t.Errorf("For test #%d, expected an error but did not actually get one.", testNumber)
			t.Logf("SRC: %#v", test.Src)
			continue
		}
	}
}

func TestRequest_Parse_failOffByOne(t *testing.T) {

	var uri string
	{
		// "mercury://x.example/" is 20 bytes.
		// "\r\n" is 2 bytes long
		// Fill the rest with 'a' to reach a request of exactly 2049 bytes total (and a URL of exactly 2047 bytes in total).
		const prefix = "mercury://x.example/"
		const targetLen = 2049

		uri = prefix + strings.Repeat("a", targetLen-len(prefix)-len("\r\n"))

		{
			actual := len(uri)
			expected := targetLen-len("\r\n")

			if expected != actual {
				t.Fatalf("TEST SETUP: the actual URI length is not what was expected")
				t.Logf("EXPECTED-LENGTH: %d", expected)
				t.Logf("ACTUAL-LENGTH:   %d", actual)
				return
			}
		}
	}

	var request hg.Request
	err := request.Parse(uri)

	// This SHOULD fail because the URI exceeds the 2048-byte limit.
	// But if there is an off-by-one error, it will succeed.
	if nil == err {
		t.Errorf("Expected an error but didn't actually get one")
		t.Logf("RequestValue length: %d", len(request.RequestValue()))
		return
	}
}

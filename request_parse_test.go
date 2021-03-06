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
			Src:      "????????????????????",
			Expected: "????????????????????\r\n",
		},
		{
			Src:      "????????????????????\r\n",
			Expected: "????????????????????\r\n",
		},
		{
			Src:      "????????????????????\r\nsecond-line",
			Expected: "????????????????????\r\n",
		},
		{
			Src:      "????????????????????\r\nsecond-line\r\n",
			Expected: "????????????????????\r\n",
		},



		{
			Src: strings.NewReader("????????????????????"),
			Expected:              "????????????????????\r\n",
		},
		{
			Src: strings.NewReader("????????????????????\r\n"),
			Expected:              "????????????????????\r\n",
		},
		{
			Src: strings.NewReader("????????????????????\r\nsecond-line"),
			Expected:              "????????????????????\r\n",
		},
		{
			Src: strings.NewReader("????????????????????\r\nsecond-line\r\n"),
			Expected:              "????????????????????\r\n",
		},



		{
			Src: []byte("????????????????????"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []byte("????????????????????\r\n"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []byte("????????????????????\r\nsecond-line"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []byte("????????????????????\r\nsecond-line\r\n"),
			Expected:   "????????????????????\r\n",
		},



		{
			Src: []rune("????????????????????"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []rune("????????????????????\r\n"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []rune("????????????????????\r\nsecond-line"),
			Expected:   "????????????????????\r\n",
		},
		{
			Src: []rune("????????????????????\r\nsecond-line\r\n"),
			Expected:   "????????????????????\r\n",
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
			Src: "????????????????????\r",
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
			Src: strings.NewReader("????????????????????\r"),
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
			Src: []byte("????????????????????\r"),
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
			Src: []rune("????????????????????\r"),
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

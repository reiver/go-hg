package hg

import (
	"fmt"
	"strings"

	"testing"
)

func TestResponseWriter_WriteHeader_success(t *testing.T) {

	tests := []struct{
		StatusCode int
		Meta interface{}
		Expected string
	}{
		{
			StatusCode: 10,
			Meta:        "",
			Expected: "10 \r\n",
		},
		{
			StatusCode: 11,
			Meta:        "",
			Expected: "11 \r\n",
		},
		{
			StatusCode: 20,
			Meta:        "",
			Expected: "20 \r\n",
		},
		{
			StatusCode: 30,
			Meta:        "",
			Expected: "30 \r\n",
		},
		{
			StatusCode: 31,
			Meta:        "",
			Expected: "31 \r\n",
		},
		{
			StatusCode: 40,
			Meta:        "",
			Expected: "40 \r\n",
		},
		{
			StatusCode: 41,
			Meta:        "",
			Expected: "41 \r\n",
		},
		{
			StatusCode: 42,
			Meta:        "",
			Expected: "42 \r\n",
		},
		{
			StatusCode: 43,
			Meta:        "",
			Expected: "43 \r\n",
		},
		{
			StatusCode: 44,
			Meta:        "",
			Expected: "44 \r\n",
		},
		{
			StatusCode: 50,
			Meta:        "",
			Expected: "50 \r\n",
		},
		{
			StatusCode: 51,
			Meta:        "",
			Expected: "51 \r\n",
		},
		{
			StatusCode: 52,
			Meta:        "",
			Expected: "52 \r\n",
		},
		{
			StatusCode: 53,
			Meta:        "",
			Expected: "53 \r\n",
		},
		{
			StatusCode: 59,
			Meta:        "",
			Expected: "59 \r\n",
		},
		{
			StatusCode: 60,
			Meta:        "",
			Expected: "60 \r\n",
		},
		{
			StatusCode: 61,
			Meta:        "",
			Expected: "61 \r\n",
		},
		{
			StatusCode: 62,
			Meta:        "",
			Expected: "62 \r\n",
		},



		{
			StatusCode: 10,
			Meta:        "ONCE twice THRICE fource",
			Expected: "10 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 11,
			Meta:        "ONCE twice THRICE fource",
			Expected: "11 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 20,
			Meta:        "ONCE twice THRICE fource",
			Expected: "20 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 30,
			Meta:        "ONCE twice THRICE fource",
			Expected: "30 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 31,
			Meta:        "ONCE twice THRICE fource",
			Expected: "31 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 40,
			Meta:        "ONCE twice THRICE fource",
			Expected: "40 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 41,
			Meta:        "ONCE twice THRICE fource",
			Expected: "41 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 42,
			Meta:        "ONCE twice THRICE fource",
			Expected: "42 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 43,
			Meta:        "ONCE twice THRICE fource",
			Expected: "43 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 44,
			Meta:        "ONCE twice THRICE fource",
			Expected: "44 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 50,
			Meta:        "ONCE twice THRICE fource",
			Expected: "50 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 51,
			Meta:        "ONCE twice THRICE fource",
			Expected: "51 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 52,
			Meta:        "ONCE twice THRICE fource",
			Expected: "52 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 53,
			Meta:        "ONCE twice THRICE fource",
			Expected: "53 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 59,
			Meta:        "ONCE twice THRICE fource",
			Expected: "59 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 60,
			Meta:        "ONCE twice THRICE fource",
			Expected: "60 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 61,
			Meta:        "ONCE twice THRICE fource",
			Expected: "61 ONCE twice THRICE fource\r\n",
		},
		{
			StatusCode: 62,
			Meta:        "ONCE twice THRICE fource",
			Expected: "62 ONCE twice THRICE fource\r\n",
		},
	}

	for testNumber, test := range tests {

		var storage strings.Builder

		var rw internalResponseWriter
		{
			rw.Writer = &storage
		}

		var responsewriter ResponseWriter = &rw

		var n int
		{
			var statusCode int = test.StatusCode
			var meta string = fmt.Sprint(test.Meta)

			var err error

			n, err = responsewriter.WriteHeader(statusCode, meta)

			if nil != err {
				t.Errorf("For test #%d, did not expect an error but actually got one.", testNumber)
				t.Logf("STATUS-CODE:       %d", test.StatusCode)
				t.Logf("META:                %q", test.Meta)
				t.Logf("EXPECTED WRITTEN: %q", test.Expected)
				t.Logf("ACTUAL   WRITTEN: %q", storage.String())
				t.Logf("EXPECTED N: %d", len(test.Expected))
				t.Logf("ACTUAL   N: %d", n)
				t.Logf("ERROR: (%T) %s", err, err)
				continue
			}
		}

		if expected, actual := len(test.Expected), n; expected != actual {
			t.Errorf("For test #%d, the actual number of bytes written is not what was expected.", testNumber)
			t.Logf("STATUS-CODE:       %d", test.StatusCode)
			t.Logf("META:                %q", test.Meta)
			t.Logf("EXPECTED WRITTEN: %q", test.Expected)
			t.Logf("ACTUAL   WRITTEN: %q", storage.String())
			t.Logf("EXPECTED N: %d", expected)
			t.Logf("ACTUAL   N: %d", actual)
			continue
		}

		if expected, actual := test.Expected, storage.String(); expected != actual {
			t.Errorf("For test #%d, the value of what was actually written is not what was expected.", testNumber)
			t.Logf("STATUS-CODE:       %d", test.StatusCode)
			t.Logf("META:                %q", test.Meta)
			t.Logf("EXPECTED WRITTEN: %q", test.Expected)
			t.Logf("ACTUAL   WRITTEN: %q", storage.String())
			t.Logf("EXPECTED N: %d", len(test.Expected))
			t.Logf("ACTUAL   N: %d", n)
			continue
		}
	}
}

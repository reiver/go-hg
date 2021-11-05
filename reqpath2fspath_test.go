package hg

import (
	"testing"
)

func TestReqpath2fspath_success(t *testing.T) {

	tests := []struct{
		RequestPath string
		Expected string
	}{
		{
			RequestPath: "/",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/default.gmni",
			Expected:     "default.gmni",
		},



		{
			RequestPath: "/apple",
			Expected:     "apple",
		},
		{
			RequestPath: "/apple/banana",
			Expected:     "apple/banana",
		},
		{
			RequestPath: "/apple/banana/cherry",
			Expected:     "apple/banana/cherry",
		},



		{
			RequestPath: "/..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../../../",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../../../..",
			Expected:     "default.gmni",
		},
		{
			RequestPath: "/../../../../../../../../../../",
			Expected:     "default.gmni",
		},
	}

	for testNumber, test := range tests {

		actual, valid := reqpath2fspath(test.RequestPath)
		if !valid {
			t.Errorf("For test #%d, expected result to be valid, but actually wasn't.", testNumber)
			t.Logf("REQUEST-PATH: %q", test.RequestPath)
			t.Logf("EXPECETED:     %q", test.Expected)
			t.Logf("ACTUAL:        %q", actual)
			continue
		}

		if expected := test.Expected; expected != actual {
			t.Errorf("For test #%d, the actual filesystem-path wasn't what was expected.", testNumber)
			t.Logf("REQUEST-PATH: %q", test.RequestPath)
			t.Logf("EXPECETED:     %q", expected)
			t.Logf("ACTUAL:        %q", actual)
			continue
		}
	}
}

func TestReqpath2fspath_failure(t *testing.T) {

	tests := []struct{
		RequestPath string
	}{
		{
			RequestPath: "apple",
		},
		{
			RequestPath: "apple/banana",
		},
		{
			RequestPath: "apple/banana/cherry",
		},



		{
			RequestPath: "apple/default.gmni",
		},
		{
			RequestPath: "apple/banana/default.gmni",
		},
		{
			RequestPath: "apple/banana/cherry/default.gmni",
		},



		{
			RequestPath: "something.gmni",
		},
		{
			RequestPath: ".gmni",
		},



		{
			RequestPath: "..",
		},
		{
			RequestPath: "../",
		},
		{
			RequestPath: "../..",
		},
		{
			RequestPath: "../../",
		},
		{
			RequestPath: "../../..",
		},
		{
			RequestPath: "../../../",
		},
		{
			RequestPath: "../../../..",
		},
		{
			RequestPath: "../../../../",
		},
		{
			RequestPath: "../../../../..",
		},
		{
			RequestPath: "../../../../../",
		},
		{
			RequestPath: "../../../../../..",
		},
		{
			RequestPath: "../../../../../../",
		},
		{
			RequestPath: "../../../../../../..",
		},
		{
			RequestPath: "../../../../../../../",
		},
		{
			RequestPath: "../../../../../../../..",
		},
		{
			RequestPath: "../../../../../../../../",
		},
		{
			RequestPath: "../../../../../../../../..",
		},
		{
			RequestPath: "../../../../../../../../../",
		},
		{
			RequestPath: "../../../../../../../../../..",
		},
		{
			RequestPath: "../../../../../../../../../../",
		},
	}

	for testNumber, test := range tests {

		actual, valid := reqpath2fspath(test.RequestPath)
		if valid {
			t.Errorf("For test #%d, expected result to be invalid, but actually wasn't.", testNumber)
			t.Logf("REQUEST-PATH: %q", test.RequestPath)
			t.Logf("ACTUAL:        %q", actual)
			continue
		}
	}
}

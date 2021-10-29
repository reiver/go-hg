package hg

import (
	"testing"
)

func TestParsetildedir_success(t *testing.T) {

	tests := []struct{
		Path string
		ExpectedUsername string
		ExpectedSubPath  string
	}{
		{
			Path:           "/~joeblow",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},
		{
			Path:          "//~joeblow",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},



		{
			Path:           "/~joeblow/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},
		{
			Path:          "//~joeblow/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},
		{
			Path:           "/~joeblow//",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},
		{
			Path:          "//~joeblow//",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/",
		},



		{
			Path:           "/~joeblow/apple",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple",
		},
		{
			Path:          "//~joeblow/apple",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple",
		},
		{
			Path:           "/~joeblow//apple",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},
		{
			Path:          "//~joeblow//apple",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},



		{
			Path:           "/~joeblow/apple/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple",
		},
		{
			Path:          "//~joeblow/apple/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple",
		},
		{
			Path:           "/~joeblow//apple/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},
		{
			Path:           "/~joeblow/apple//",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple",
		},
		{
			Path:          "//~joeblow//apple/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},
		{
			Path:           "/~joeblow//apple//",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},
		{
			Path:          "//~joeblow//apple//",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:          "/apple",
		},



		{
			Path:           "/~joeblow/apple/banana",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple/banana",
		},



		{
			Path:           "/~joeblow/apple/banana/",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple/banana",
		},



		{
			Path:           "/~joeblow/apple/banana/cherry",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple/banana/cherry",
		},



		{
			Path:           "/~joeblow/apple/banana/cherry.txt",
			ExpectedUsername: "joeblow",
			ExpectedSubPath:         "/apple/banana/cherry.txt",
		},
	}

	for testNumber, test := range tests {

		actualUsername, actualSubPath, successful := parsetildedir(test.Path)
		if !successful {
			t.Errorf("For test #%d, the parsing was no successful", testNumber)
			t.Logf("PATH: %q", test.Path)
			t.Logf("SUCCESSFUL: %t", successful)
			t.Logf("EXPECTED USERNAME: %q", test.ExpectedUsername)
			t.Logf("ACTUAL   USERNAME: %q", actualUsername)
			t.Logf("EXPECTED SUB-PATH: %q", test.ExpectedSubPath)
			t.Logf("ACTUAL   SUB-PATH: %q", actualSubPath)
			continue
		}

		if expected, actual := test.ExpectedUsername, actualUsername; expected != actual  {
			t.Errorf("For test #%d, the actual usernam is not what was expected.", testNumber)
			t.Logf("PATH: %q", test.Path)
			t.Logf("SUCCESSFUL: %t", successful)
			t.Logf("EXPECTED USERNAME: %q", test.ExpectedUsername)
			t.Logf("ACTUAL   USERNAME: %q", actualUsername)
			t.Logf("EXPECTED SUB-PATH: %q", test.ExpectedSubPath)
			t.Logf("ACTUAL   SUB-PATH: %q", actualSubPath)
			continue
		}

		if expected, actual := test.ExpectedSubPath, actualSubPath; expected != actual  {
			t.Errorf("For test #%d, the actual usernam is not what was expected.", testNumber)
			t.Logf("PATH: %q", test.Path)
			t.Logf("SUCCESSFUL: %t", successful)
			t.Logf("EXPECTED USERNAME: %q", test.ExpectedUsername)
			t.Logf("ACTUAL   USERNAME: %q", actualUsername)
			t.Logf("EXPECTED SUB-PATH: %q", test.ExpectedSubPath)
			t.Logf("ACTUAL   SUB-PATH: %q", actualSubPath)
			continue
		}
	}
}

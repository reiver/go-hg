package hg

import (
	"testing"
)

func TestStatuscodeFromRunes_success(t *testing.T) {

	tests := []struct{
		MostSignificant  rune
		LeastSignificant rune
		Expected int
	}{
		{
			MostSignificant: '0',
			LeastSignificant: '0',
			Expected:          0,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '1',
			Expected:          1,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '2',
			Expected:          2,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '3',
			Expected:          3,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '4',
			Expected:          4,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '5',
			Expected:          5,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '6',
			Expected:          6,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '7',
			Expected:          7,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '8',
			Expected:          8,
		},
		{
			MostSignificant: '0',
			LeastSignificant: '9',
			Expected:          9,
		},









		{
			MostSignificant: '1',
			LeastSignificant: '0',
			Expected:         10,
		},
		{
			MostSignificant: '1',
			LeastSignificant: '1',
			Expected:         11,
		},



		{
			MostSignificant: '2',
			LeastSignificant: '0',
			Expected:         20,
		},



		{
			MostSignificant: '3',
			LeastSignificant: '0',
			Expected:         30,
		},
		{
			MostSignificant: '3',
			LeastSignificant: '1',
			Expected:         31,
		},



		{
			MostSignificant: '4',
			LeastSignificant: '0',
			Expected:         40,
		},
		{
			MostSignificant: '4',
			LeastSignificant: '1',
			Expected:         41,
		},
		{
			MostSignificant: '4',
			LeastSignificant: '2',
			Expected:         42,
		},
		{
			MostSignificant: '4',
			LeastSignificant: '3',
			Expected:         43,
		},
		{
			MostSignificant: '4',
			LeastSignificant: '4',
			Expected:         44,
		},



		{
			MostSignificant: '5',
			LeastSignificant: '0',
			Expected:         50,
		},
		{
			MostSignificant: '5',
			LeastSignificant: '1',
			Expected:         51,
		},
		{
			MostSignificant: '5',
			LeastSignificant: '2',
			Expected:         52,
		},
		{
			MostSignificant: '5',
			LeastSignificant: '3',
			Expected:         53,
		},

		{
			MostSignificant: '5',
			LeastSignificant: '9',
			Expected:         59,
		},



		{
			MostSignificant: '6',
			LeastSignificant: '0',
			Expected:         60,
		},
		{
			MostSignificant: '6',
			LeastSignificant: '1',
			Expected:         61,
		},
		{
			MostSignificant: '6',
			LeastSignificant: '2',
			Expected:         62,
		},









		{
			MostSignificant: '8',
			LeastSignificant: '9',
			Expected:         89,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '0',
			Expected:         90,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '1',
			Expected:         91,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '2',
			Expected:         92,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '3',
			Expected:         93,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '4',
			Expected:         94,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '5',
			Expected:         95,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '6',
			Expected:         96,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '7',
			Expected:         97,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '8',
			Expected:         98,
		},
		{
			MostSignificant: '9',
			LeastSignificant: '9',
			Expected:         99,
		},
	}

	for testNumber, test := range tests {

		actual, valid := statuscodeFromRunes(test.MostSignificant, test.LeastSignificant)
		if !valid {
			t.Errorf("For test #%d, did not expect result to invalid", testNumber)
			t.Logf(" MOST-SIGNIFICANT: %q", string(test.MostSignificant))
			t.Logf("LEAST-SIGNIFICANT: %q", string(test.LeastSignificant))
			t.Logf("EXPECTED: %d", test.Expected)
			continue
		}

		if expected := test.Expected; expected != actual {
			t.Errorf("For test #%d, the actual status-code was not what was expected.", testNumber)
			t.Logf(" MOST-SIGNIFICANT: %q", string(test.MostSignificant))
			t.Logf("LEAST-SIGNIFICANT: %q", string(test.LeastSignificant))
			t.Logf("EXPECTED: %d", expected)
			t.Logf("ACTUAL:   %d", actual)
			continue
		}
	}
}

func TestStatuscodeFromRunes_failure(t *testing.T) {

	tests := []struct{
		MostSignificant  rune
		LeastSignificant rune
	}{
		{
			MostSignificant: 0,
			LeastSignificant: 0,
		},



		{
			MostSignificant: 0,
			LeastSignificant: 0,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 1,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 2,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 3,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 4,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 5,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 6,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 7,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 8,
		},
		{
			MostSignificant: 0,
			LeastSignificant: 9,
		},



		{
			MostSignificant: 1,
			LeastSignificant: 0,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 1,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 2,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 3,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 4,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 5,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 6,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 7,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 8,
		},
		{
			MostSignificant: 1,
			LeastSignificant: 9,
		},



		{
			MostSignificant: 2,
			LeastSignificant: 0,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 1,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 2,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 3,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 4,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 5,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 6,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 7,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 8,
		},
		{
			MostSignificant: 2,
			LeastSignificant: 9,
		},



		{
			MostSignificant: 3,
			LeastSignificant: 0,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 1,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 2,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 3,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 4,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 5,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 6,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 7,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 8,
		},
		{
			MostSignificant: 3,
			LeastSignificant: 9,
		},



		{
			MostSignificant: 9,
			LeastSignificant: 0,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 1,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 2,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 3,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 4,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 5,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 6,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 7,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 8,
		},
		{
			MostSignificant: 9,
			LeastSignificant: 9,
		},



		{
			MostSignificant: '/',
			LeastSignificant: '0',
		},
		{
			MostSignificant: '0',
			LeastSignificant: '/',
		},



		{
			MostSignificant: ':',
			LeastSignificant: '9',
		},
		{
			MostSignificant: '9',
			LeastSignificant: ':',
		},



		{
			MostSignificant: 'a',
			LeastSignificant: 'b',
		},



		{
			MostSignificant: 0x7F,
			LeastSignificant: 0x7F,
		},
	}

	for testNumber, test := range tests {

		actual, valid := statuscodeFromRunes(test.MostSignificant, test.LeastSignificant)
		if valid {
			t.Errorf("For test #%d, did not expect result to valid", testNumber)
			t.Logf(" MOST-SIGNIFICANT: %q", string(test.MostSignificant))
			t.Logf("LEAST-SIGNIFICANT: %q", string(test.LeastSignificant))
			t.Logf("ACTUAL: %d", actual)
			continue
		}
	}
}

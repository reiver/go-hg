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
v

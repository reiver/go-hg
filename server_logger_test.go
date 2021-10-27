package hg

import (
	"testing"
)

func TestServer_logger(t *testing.T) {

	tests := []struct{
		Server Server
		Expected Logger
	}{
		{
			Server: Server{},
			Expected: internalDiscardLogger{},
		},



		{
			Server: Server{Logger:nil},
			Expected: internalDiscardLogger{},
		},
		{
			Server: Server{Logger:Logger(nil)},
			Expected: internalDiscardLogger{},
		},



		{
			Server: Server{Logger:testLogger{}},
			Expected: testLogger{},
		},
		{
			Server: Server{Logger:new(testLogger)},
			Expected: &testLogger{},
		},
	}

	for testNumber, test := range tests  {

		var logger Logger = test.Server.logger()

		if expected, actual := test.Expected, logger; expected != actual {
			t.Errorf("For test #%d, the actual logger is not what was expected.", testNumber)
			t.Logf("SERVER: %#v", test.Server)
			t.Logf("EXPECTED: (%T) %#v", expected, expected)
			t.Logf("ACTUAL:   (%T) %#v", actual, actual)
			continue
		}

	}
}

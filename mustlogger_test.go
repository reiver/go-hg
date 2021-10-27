package hg

import (
	"testing"
)

func TestMustlogger_nil(t *testing.T) {

	tests := []struct{
		Logger Logger
	}{
		{
			Logger: nil,
		},
		{
			Logger: Logger(nil),
		},
	}

	for testNumber, test := range tests {

		var logger Logger = mustlogger(test.Logger)

		if nil == logger || Logger(nil) == logger {
			t.Errorf("For test #%d, expected a valid logger but did not actually get one.", testNumber)
			t.Logf("LOGGER: (%T) %$v", logger, logger)
			continue
		}

		switch logger.(type) {
		case internalDiscardLogger:
			// Nothing here.
		case *internalDiscardLogger:
			// Nothing here.
		default:
			t.Errorf("For test #%d, the actual type of the logger was not what was expected.", testNumber)
			t.Log("EXPECTED: hg.internalDiscardLogger")
			t.Logf("ACTUAL:   %T", logger)
			continue
		}
	}
}

func TestMustlogger_valid(t *testing.T) {

	tests := []struct{
		Logger Logger
	}{
		{
			Logger: testLogger{},
		},
		{
			Logger: new(testLogger),
		},
	}

	for testNumber, test := range tests {

		var logger Logger = mustlogger(test.Logger)

		if nil == logger ||  Logger(nil) == logger {
			t.Errorf("For test #%d, expected a valid logger but did not actually get one.", testNumber)
			t.Logf("LOGGER: (%T) %$v", logger, logger)
			continue
		}

	}
}

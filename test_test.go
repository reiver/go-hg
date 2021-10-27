package hg

type testLogger struct {}
func (testLogger) Error(...interface{}) {}
func (testLogger) Errorf(string, ...interface{}) {}
func (testLogger) Log(...interface{}) {}
func (testLogger) Logf(string, ...interface{}) {}

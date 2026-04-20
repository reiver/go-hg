package hg

type testLogger struct {}
func (testLogger) Begin(fields ...Field) Logger {return testLogger{}}
func (testLogger) End(fields ...Field) {}
func (testLogger) Error(fields ...Field) {}
func (testLogger) Debug(fields ...Field) {}
func (testLogger) Trace(fields ...Field) {}

var _ Logger = testLogger{}

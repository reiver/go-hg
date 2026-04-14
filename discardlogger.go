package hg

type internalDiscardLogger struct{}

var _ Logger = internalDiscardLogger{}

func (internalDiscardLogger) Begin(fields ...Field) Logger {return internalDiscardLogger{}}

func (internalDiscardLogger) Debug(fields ...Field) {}

func (internalDiscardLogger) End(fields ...Field) {}

func (internalDiscardLogger) Error(fields ...Field) {}

func (internalDiscardLogger) Trace(fields ...Field) {}

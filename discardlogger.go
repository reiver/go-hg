package hg

type internalDiscardLogger struct{}

func (internalDiscardLogger) Log(...interface{}) {}
func (internalDiscardLogger) Logf(string, ...interface{}) {}

func (internalDiscardLogger) Error(...interface{}) {}
func (internalDiscardLogger) Errorf(string, ...interface{}) {}

func (internalDiscardLogger) Trace(...interface{}) {}
func (internalDiscardLogger) Tracef(string, ...interface{}) {}

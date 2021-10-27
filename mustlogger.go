package hg

func mustlogger(logger Logger) Logger {
	if nil == logger {
		logger = internalDiscardLogger{}
	}

	return logger
}

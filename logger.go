package hg

type Logger interface {
	Error(...interface{})
	Errorf(string, ...interface{})

	Log(...interface{})
	Logf(string, ...interface{})
}

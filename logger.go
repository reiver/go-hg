package hg

import (
	"codeberg.org/reiver/go-field"
)

type Field = field.StringlyField

type Logger interface {
	Begin(fields ...Field) Logger
	End(fields ...Field)
	Error(fields ...Field)
	Debug(fields ...Field)
	Trace(fields ...Field)
}

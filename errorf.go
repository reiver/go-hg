package hg

import (
	"fmt"
)

// errorf is similar to fmt.Errorf(), except it prefixes the error with the name of this package.
func errorf(format string, a ...interface{}) error {
	const prefix = "hg: "

	format = prefix + format

	return fmt.Errorf(format, a...)
}

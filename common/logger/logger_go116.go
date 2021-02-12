// +build go1.16

package logger

import (
	"io"
)

func newLogger(w io.Writer) *logger {
	if w == nil {
		w = io.Discard
	}
	l := &logger{
		w: w,
	}
	return l
}

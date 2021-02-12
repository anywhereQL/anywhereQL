// +build !go1.16

package logger

import (
	"io"
	"io/ioutil"
)

func newLogger(w io.Writer) *logger {
	if w == nil {
		w = ioutil.Discard
	}
	l := &logger{
		w: w,
	}
	return l
}

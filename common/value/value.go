package value

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/anywhereQL/anywhereQL/common/helper"
)

type Type int

const (
	UNKNOWN Type = iota
	INTEGER
)

func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "Unknown value"
	case INTEGER:
		return "Ingeer value"
	default:
		return "Error Unknwo value type"
	}
}

type Value struct {
	Type Type
	Int  int64
}

func Convert(s string) (Value, error) {
	r := []rune(s)
	if !helper.IsDigit(r[0]) {
		return Value{Type: UNKNOWN}, errors.New(fmt.Sprintf("Cannot Convert %s", s))
	}
	for _, ch := range r {
		if !helper.IsDigit(ch) {
			return Value{Type: UNKNOWN}, errors.New(fmt.Sprintf("Cannot Convert %s", s))
		}
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Value{Type: UNKNOWN}, err
	}
	return Value{Type: INTEGER, Int: v}, nil
}

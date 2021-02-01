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
	FLOAT
)

func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "Unknown value"
	case INTEGER:
		return "Ingeer value"
	case FLOAT:
		return "Floating point value"
	default:
		return "Error Unknwo value type"
	}
}

type Value struct {
	Type   Type
	Int    int64
	Float  float64
	PartI  int64
	PartF  int64
	FDigit int
}

func Convert(s string) (Value, error) {
	r := []rune(s)
	isFloating := false
	if r[0] == '.' {
		isFloating = true
	}
	for _, ch := range r {
		if !helper.IsDigit(ch) && (ch == '.' && isFloating == false) {
			isFloating = true
		} else if !helper.IsDigit(ch) && ch != '.' {
			return Value{Type: UNKNOWN}, errors.New(fmt.Sprintf("Cannot Convert %s", s))
		}
	}
	if !isFloating {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return Value{Type: UNKNOWN}, err
		}
		return Value{Type: INTEGER, Int: v}, nil
	} else {
		vi := []rune{}
		vf := []rune{}
		isFloatingPart := false
		for _, ch := range r {
			if ch == '.' {
				isFloatingPart = true
				continue
			}
			if !isFloatingPart {
				vi = append(vi, ch)
			} else {
				vf = append(vf, ch)
			}
		}
		fDigit := len(vf)
		partI, _ := strconv.ParseInt(string(vi), 10, 64)
		partF, _ := strconv.ParseInt(string(vf), 10, 64)

		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return Value{Type: UNKNOWN}, err
		}
		return Value{Type: FLOAT, Float: v, PartF: partF, PartI: partI, FDigit: fDigit}, nil
	}
}

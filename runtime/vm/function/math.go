package function

import (
	"fmt"
	"math"

	"github.com/anywhereQL/anywhereQL/common/value"
)

func Abs(args []value.Value) (value.Value, error) {
	r := value.Value{}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg too long")
	}
	v := args[0]
	switch v.Type {
	case value.INTEGER:
		r.Type = value.INTEGER
		if v.Int < 0 {
			r.Int = -1 * v.Int
		} else {
			r.Int = v.Int
		}
	case value.FLOAT:
		r.Type = value.FLOAT
		r.Float = math.Abs(v.Float)
	default:
		return r, fmt.Errorf("Args is unknown type")
	}

	return r, nil
}

func Sign(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.INTEGER,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg too long")
	}
	v := args[0]
	switch v.Type {
	case value.INTEGER:
		if v.Int < 0 {
			r.Int = -1
		} else if v.Int > 0 {
			r.Int = 1
		} else {
			r.Int = 0
		}
	case value.FLOAT:
		if v.Float < 0 {
			r.Int = -1
		} else if v.Float > 0 {
			r.Int = 1
		} else {
			r.Int = 0
		}
	default:
		return r, fmt.Errorf("Args is unknown type")
	}

	return r, nil
}

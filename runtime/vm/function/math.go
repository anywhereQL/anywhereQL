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

func Floor(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.INTEGER,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg too long")
	}

	v := args[0]
	switch v.Type {
	case value.INTEGER:
		r.Int = v.Int
	case value.FLOAT:
		r.Int = int64(math.Floor(v.Float))
	}
	return r, nil
}

func Ceil(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.INTEGER,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg too long")
	}

	v := args[0]
	switch v.Type {
	case value.INTEGER:
		r.Int = v.Int
	case value.FLOAT:
		r.Int = int64(math.Ceil(v.Float))
	}
	return r, nil
}

func Round(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	d := 0
	var v value.Value
	if len(args) == 1 {
		v = args[0]
	} else if len(args) == 2 {
		if args[0].Type != value.INTEGER {
			return r, fmt.Errorf("Arg Type is ignore")
		}
		d = int(args[0].Int)
		v = args[1]
	} else {
		return r, fmt.Errorf("Arg length is mismatch")
	}

	var op1 float64
	switch v.Type {
	case value.INTEGER:
		op1 = float64(v.Int)
	case value.FLOAT:
		op1 = v.Float
	}
	t := op1 * math.Pow(10, float64(d))
	rnd := math.Round(t)
	rnd = rnd / math.Pow(10, float64(d))

	r.Float = rnd
	return r, nil
}

func Trunc(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	d := 0
	var v value.Value
	if len(args) == 1 {
		v = args[0]
	} else if len(args) == 2 {
		if args[0].Type != value.INTEGER {
			return r, fmt.Errorf("Arg Type is ignore")
		}
		d = int(args[0].Int)
		v = args[1]
	} else {
		return r, fmt.Errorf("Arg length is mismatch")
	}

	var op1 float64
	switch v.Type {
	case value.INTEGER:
		op1 = float64(v.Int)
	case value.FLOAT:
		op1 = v.Float
	}
	t := op1 * math.Pow(10, float64(d))
	trn := math.Trunc(t)
	trn = trn / math.Pow(10, float64(d))

	r.Float = trn
	return r, nil
}

func Pow(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 2 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64
	var op2 float64

	if args[1].Type == value.INTEGER {
		op1 = float64(args[1].Int)
	} else if args[1].Type == value.FLOAT {
		op1 = args[1].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	if args[0].Type == value.INTEGER {
		op2 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op2 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Pow(op1, op2)
	return r, nil
}

func Sqrt(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Sqrt(op1)
	return r, nil
}

func Exp(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Exp(op1)
	return r, nil
}

func Ln(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Log(op1)
	return r, nil
}

func Log10(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Log10(op1)
	return r, nil
}

func Sin(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Sin(op1)
	return r, nil
}

func Cos(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Cos(op1)
	return r, nil
}

func Cosh(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Cosh(op1)
	return r, nil
}

func Tan(args []value.Value) (value.Value, error) {
	r := value.Value{
		Type: value.FLOAT,
	}

	if len(args) != 1 {
		return r, fmt.Errorf("Arg length mismatch")
	}
	var op1 float64

	if args[0].Type == value.INTEGER {
		op1 = float64(args[0].Int)
	} else if args[0].Type == value.FLOAT {
		op1 = args[0].Float
	} else {
		return r, fmt.Errorf("Arg type unknown")
	}

	r.Float = math.Tan(op1)
	return r, nil
}

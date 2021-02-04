package function

import (
	"strings"

	"github.com/anywhereQL/anywhereQL/common/value"
)

type CallFunction func([]value.Value) (value.Value, error)

var funcs = map[string]CallFunction{
	// Math funcs
	"abs":      Abs,
	"sign":     Sign,
	"sgn":      Sign,
	"floor":    Floor,
	"ceil":     Ceil,
	"ceiling":  Ceil,
	"round":    Round,
	"trunc":    Trunc,
	"truncate": Trunc,
	"pow":      Pow,
	"power":    Pow,
	"sqrt":     Sqrt,
	"exp":      Exp,
	"ln":       Ln,
	"log10":    Log10,
	"sin":      Sin,
	"cos":      Cos,
	"cosh":     Cosh,
	"tan":      Tan,
	"acos":     Acos,
	"asin":     Asin,
	"atan":     Atan,
	"atan2":    Atan2,
	"cot":      Cot,
	"degrees":  Degrees,
	"radians":  Radians,
	"pi":       Pi,
	"rand":     Rand,
	"greatest": Greatest,
	"least":    Least,
}

func LookupFunction(name string) CallFunction {
	if f, exists := funcs[strings.ToLower(name)]; exists {
		return f
	}
	return nil
}

package function

import (
	"strings"

	"github.com/anywhereQL/anywhereQL/common/value"
)

type CallFunction func([]value.Value) (value.Value, error)

var funcs = map[string]CallFunction{
	"abs": Abs,
	"sign": Sign,

}

func LookupFunction(name string) CallFunction {
	if f, exists := funcs[strings.ToLower(name)]; exists {
		return f
	}
	return nil
}

package function

import (
	"strings"

	"github.com/anywhereQL/anywhereQL/common/value"
)

type callFunction func([]value.Value) (value.Value, error)

var funcs = map[string]callFunction{}

func LookupFunction(name string) callFunction {
	if f, exists := funcs[strings.ToLower(name)]; exists {
		return f
	}
	return nil
}

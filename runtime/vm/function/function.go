package function

import "github.com/anywhereQL/anywhereQL/common/result"

type callFunction func([]interface{}) result.Value

var funcs = map[string]callFunction{}

func LookupFunction(name string) callFunction {
	if f, exists := funcs[name]; exists {
		return f
	}
	return nil
}

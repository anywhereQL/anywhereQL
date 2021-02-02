package vm

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/value"
)

type stack []value.Value

func (s *stack) size() int {
	return len(*s)
}

func (s *stack) empty() bool {
	return s.size() == 0
}

func newStack() *stack {
	s := new(stack)
	return s
}

func (s *stack) push(v value.Value) {
	*s = append(*s, v)
}

func (s *stack) pop() (value.Value, error) {
	if s.empty() {
		return value.Value{}, fmt.Errorf("stack underflow error")
	}

	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, nil
}

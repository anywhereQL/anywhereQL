package vm

import "fmt"

type stack []VMValue

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

func (s *stack) push(v VMValue) {
	*s = append(*s, v)
}

func (s *stack) pop() (VMValue, error) {
	if s.empty() {
		return VMValue{}, fmt.Errorf("stack underflow error")
	}

	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v, nil
}

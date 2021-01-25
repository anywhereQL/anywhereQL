package vm

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/result"
	"github.com/anywhereQL/anywhereQL/runtime/vm/function"
)

type OpeType int

const (
	_ OpeType = iota
	PUSH
	POP
	ADD
	SUB
	MUL
	DIV
	MOD
	STORE
	CALL
)

func (o OpeType) String() string {
	switch o {
	case PUSH:
		return "PUSH"
	case POP:
		return "POP"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MUL:
		return "MUL"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case CALL:
		return "CALL"
	case STORE:
		return "STORE"
	default:
		return "Unknwon Operation"
	}
}

type ValueType int

const (
	_ ValueType = iota
	NA
	Integer
	String
)

func (v ValueType) String() string {
	switch v {
	case NA:
		return "N/A"
	case Integer:
		return "Integer"
	case String:
		return "String"
	default:
		return "Unknown"
	}
}

type VMValue struct {
	Type     ValueType
	Integral int64
	String   string
}

type VMCode struct {
	Operator OpeType
	Operand1 VMValue
	Operand2 VMValue
}

func (c VMCode) String() string {
	s := ""
	s = fmt.Sprintf("%s", c.Operator)

	if c.Operand1.Type != NA {
		switch c.Operand1.Type {
		case Integer:
			s = fmt.Sprintf("%s %d", s, c.Operand1.Integral)
		case String:
			s = fmt.Sprintf("%s %s", s, c.Operand1.String)
		}
	}

	if c.Operand2.Type != NA {
		switch c.Operand2.Type {
		case Integer:
			s = fmt.Sprintf("%s %d", s, c.Operand2.Integral)
		case String:
			s = fmt.Sprintf("%s %s", s, c.Operand2.String)
		}
	}

	return s
}

func Run(codes []VMCode) ([]result.Value, error) {
	s := newStack()
	cols := []result.Value{}

	for _, code := range codes {
		switch code.Operator {
		case PUSH:
			s.push(code.Operand1)
		case ADD:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			v := VMValue{
				Type:     Integer,
				Integral: ope1.Integral + ope2.Integral,
			}
			s.push(v)
		case SUB:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			v := VMValue{
				Type:     Integer,
				Integral: ope1.Integral - ope2.Integral,
			}
			s.push(v)
		case MUL:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			v := VMValue{
				Type:     Integer,
				Integral: ope1.Integral * ope2.Integral,
			}
			s.push(v)
		case DIV:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			if ope2.Integral == 0 {
				return []result.Value{}, fmt.Errorf("Div by 0")
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			v := VMValue{
				Type:     Integer,
				Integral: ope1.Integral / ope2.Integral,
			}
			s.push(v)
		case MOD:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			if ope2.Integral == 0 {
				return []result.Value{}, fmt.Errorf("Div by 0")
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			v := VMValue{
				Type:     Integer,
				Integral: ope1.Integral % ope2.Integral,
			}
			s.push(v)
		case CALL:
			args := []interface{}{}

			argsN, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			for i := 0; int64(i) < argsN.Integral; i++ {
				v, err := s.pop()
				if err != nil {
					return []result.Value{}, err
				}
				switch v.Type {
				case Integer:
					args = append(args, v.Integral)
				case String:
					args = append(args, v.String)
				}
			}

			call := function.LookupFunction(code.Operand1.String)
			if call == nil {
				return []result.Value{}, fmt.Errorf("Function(%s) is not implement", code.Operand1.String)
			}
			r := call(args)
			var vr VMValue
			switch r.Type {
			case result.Integral:
				vr = VMValue{
					Type:     Integer,
					Integral: r.Integral,
				}
			}
			s.push(vr)
		case STORE:
			v, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			if v.Type == Integer {
				cols = append(cols, result.Value{Type: result.Integral, Integral: v.Integral})
			}
		}
	}
	return cols, nil
}

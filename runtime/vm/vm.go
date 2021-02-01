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
	Float
	Decimal
	String
)

func (v ValueType) String() string {
	switch v {
	case NA:
		return "N/A"
	case Integer:
		return "Integer"
	case Float:
		return "Float"
	case Decimal:
		return "Decimal"
	case String:
		return "String"
	default:
		return "Unknown"
	}
}

type VMValue struct {
	Type     ValueType
	Integral int64
	Float    float64
	String   string

	PartI  int64
	PartF  int64
	FDigit int
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
		case Float:
			s = fmt.Sprintf("%s %f", s, c.Operand1.Float)
		case String:
			s = fmt.Sprintf("%s %s", s, c.Operand1.String)
		}
	}

	if c.Operand2.Type != NA {
		switch c.Operand2.Type {
		case Integer:
			s = fmt.Sprintf("%s %d", s, c.Operand2.Integral)
		case Float:
			s = fmt.Sprintf("%s %f", s, c.Operand2.Float)
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
			if ope1.Type == Integer && ope2.Type == Integer {
				v := VMValue{
					Type:     Integer,
					Integral: ope1.Integral + ope2.Integral,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float + ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Integer {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float + float64(ope2.Integral),
				}
				s.push(v)
			} else if ope1.Type == Integer && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: float64(ope1.Integral) + ope2.Float,
				}
				s.push(v)
			} else {
				return []result.Value{}, fmt.Errorf("Unknown Operation: %s + %s", ope1.Type, ope2.Type)
			}

		case SUB:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			if ope1.Type == Integer && ope2.Type == Integer {
				v := VMValue{
					Type:     Integer,
					Integral: ope1.Integral - ope2.Integral,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float - ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Integer {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float - float64(ope2.Integral),
				}
				s.push(v)
			} else if ope1.Type == Integer && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: float64(ope1.Integral) - ope2.Float,
				}
				s.push(v)
			} else {
				return []result.Value{}, fmt.Errorf("Unknown Operation: %s - %s", ope1.Type, ope2.Type)
			}

		case MUL:
			ope2, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			ope1, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			if ope1.Type == Integer && ope2.Type == Integer {
				v := VMValue{
					Type:     Integer,
					Integral: ope1.Integral * ope2.Integral,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float * ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Integer {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float * float64(ope2.Integral),
				}
				s.push(v)
			} else if ope1.Type == Integer && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: float64(ope1.Integral) * ope2.Float,
				}
				s.push(v)
			} else {
				return []result.Value{}, fmt.Errorf("Unknown Operation: %s * %s", ope1.Type, ope2.Type)
			}

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
			if ope1.Type == Integer && ope2.Type == Integer {
				v := VMValue{
					Type:     Integer,
					Integral: ope1.Integral / ope2.Integral,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float / ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == Float && ope2.Type == Integer {
				v := VMValue{
					Type:  Float,
					Float: ope1.Float / float64(ope2.Integral),
				}
				s.push(v)
			} else if ope1.Type == Integer && ope2.Type == Float {
				v := VMValue{
					Type:  Float,
					Float: float64(ope1.Integral) / ope2.Float,
				}
				s.push(v)
			} else {
				return []result.Value{}, fmt.Errorf("Unknown Operation: %s / %s", ope1.Type, ope2.Type)
			}

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
			if ope1.Type == Integer && ope2.Type == Integer {
				v := VMValue{
					Type:     Integer,
					Integral: ope1.Integral % ope2.Integral,
				}
				s.push(v)
			} else {
				return []result.Value{}, fmt.Errorf("Unknown Operation: %s %% %s", ope1.Type, ope2.Type)
			}
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
				case Float:
					args = append(args, v.Float)
				case String:
					args = append(args, v.String)
				default:
					return []result.Value{}, fmt.Errorf("Unknwon Argument Type: %s", v.Type)
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
			case result.Float:
				vr = VMValue{
					Type:   Float,
					Float:  r.Float,
					PartI:  r.PartI,
					PartF:  r.PartF,
					FDigit: r.FDigit,
				}
			}
			s.push(vr)
		case STORE:
			v, err := s.pop()
			if err != nil {
				return []result.Value{}, err
			}
			switch v.Type {
			case Integer:
				cols = append(cols, result.Value{Type: result.Integral, Integral: v.Integral})
			case Float:
				cols = append(cols, result.Value{Type: result.Float, Float: v.Float})
			case Decimal:
				cols = append(cols, result.Value{Type: result.Decimal, PartI: v.PartI, PartF: v.PartF, FDigit: v.FDigit})
			}
		}
	}
	return cols, nil
}

package vm

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/storage/virtual"
	"github.com/anywhereQL/anywhereQL/runtime/vm/function"
)

type ExprOpeType int

const (
	_ ExprOpeType = iota
	NA
	PUSH
	POP
	ADD
	SUB
	MUL
	DIV
	MOD
	STORE
	CALL
	PICK
)

func (o ExprOpeType) String() string {
	switch o {
	case NA:
		return "Not impliement"
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
	case PICK:
		return "PICK"
	default:
		return "Unknwon Operation"
	}
}

type ExprVMCode struct {
	Operator ExprOpeType
	Operand1 value.Value
	Operand2 value.Value
}

func (c ExprVMCode) String() string {
	s := ""
	s = fmt.Sprintf("%s", c.Operator)

	if c.Operand1.Type != value.NA {
		switch c.Operand1.Type {
		case value.INTEGER:
			s = fmt.Sprintf("%s %d", s, c.Operand1.Int)
		case value.FLOAT:
			s = fmt.Sprintf("%s %f", s, c.Operand1.Float)
		case value.STRING:
			s = fmt.Sprintf("%s %s", s, c.Operand1.String)
		case value.COLUMN:
			s = fmt.Sprintf("%s %s.%s", s, c.Operand1.Column.TableID, c.Operand1.Column.Column)
		case value.BOOL:
			if c.Operand1.Bool.True {
				s = fmt.Sprintf("%s TRUE", s)
			} else {
				s = fmt.Sprintf("%s FALSE", s)
			}
		}
	}

	if c.Operand2.Type != value.NA {
		switch c.Operand2.Type {
		case value.INTEGER:
			s = fmt.Sprintf("%s %d", s, c.Operand2.Int)
		case value.FLOAT:
			s = fmt.Sprintf("%s %f", s, c.Operand2.Float)
		case value.STRING:
			s = fmt.Sprintf("%s %s", s, c.Operand2.String)
		case value.COLUMN:
			s = fmt.Sprintf("%s %s.%s", s, c.Operand1.Column.TableID, c.Operand1.Column.Column)
		case value.BOOL:
			if c.Operand2.Bool.True {
				s = fmt.Sprintf("%s TRUE", s)
			} else {
				s = fmt.Sprintf("%s FALSE", s)
			}
		}
	}

	return s
}

func ExprRun(codes []ExprVMCode, line int) (value.Value, error) {
	s := newStack()
	col := value.Value{}

	for _, code := range codes {
		switch code.Operator {
		case PUSH:
			s.push(code.Operand1)
		case ADD:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int + ope2.Int,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float + ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float + float64(ope2.Int),
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: float64(ope1.Int) + ope2.Float,
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s + %s", ope1.Type, ope2.Type)
			}

		case SUB:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int - ope2.Int,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float - ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float - float64(ope2.Int),
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: float64(ope1.Int) - ope2.Float,
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s - %s", ope1.Type, ope2.Type)
			}

		case MUL:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int * ope2.Int,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float * ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float * float64(ope2.Int),
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: float64(ope1.Int) * ope2.Float,
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s * %s", ope1.Type, ope2.Type)
			}

		case DIV:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope2.Int == 0 {
					return col, fmt.Errorf("Div by 0")
				}
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int / ope2.Int,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope2.Float == 0 {
					return col, fmt.Errorf("Div by 0")
				}
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float / ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope2.Int == 0 {
					return col, fmt.Errorf("Div by 0")
				}
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float / float64(ope2.Int),
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if ope2.Float == 0 {
					return col, fmt.Errorf("Div by 0")
				}
				v := value.Value{
					Type:  value.FLOAT,
					Float: float64(ope1.Int) / ope2.Float,
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s / %s", ope1.Type, ope2.Type)
			}

		case MOD:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope2.Int == 0 {
				return col, fmt.Errorf("Div by 0")
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int % ope2.Int,
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s %% %s", ope1.Type, ope2.Type)
			}
		case CALL:
			args := []value.Value{}

			argsN, err := s.pop()
			if err != nil {
				return col, err
			}
			for i := 0; int64(i) < argsN.Int; i++ {
				v, err := s.pop()
				if err != nil {
					return col, err
				}
				args = append(args, v)
			}

			call := function.LookupFunction(code.Operand1.String)
			if call == nil {
				return col, fmt.Errorf("Function(%s) is not implement", code.Operand1.String)
			}
			r, err := call(args)
			if err != nil {
				return col, err
			}
			s.push(r)
		case STORE:
			v, err := s.pop()
			if err != nil {
				return col, err
			}
			col = v
		case PICK:
			eng := virtual.VirtualStorage
			r, err := eng.GetValue(code.Operand1.Column.TableID, code.Operand1.Column.Column, line)
			if err != nil {
				return col, err
			}
			s.push(r)
		case NA:
			panic("")
		}
	}
	return col, nil
}

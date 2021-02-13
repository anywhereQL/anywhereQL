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
	EQ
	NEQ
	LT
	LTE
	GT
	GTE
	AND
	OR
	NOT
	CAST
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
	case EQ:
		return "EQUAL"
	case NEQ:
		return "NOT EQUAL"
	case LT:
		return "LESS THAN"
	case LTE:
		return "LESS THAN EQUAL"
	case GT:
		return "GREATER THAN"
	case GTE:
		return "GREATER THAN EQUAL"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case CAST:
		return "CAST"
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
		case value.NULL:
			s = fmt.Sprintf("%s NULL", s)
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
		case value.NULL:
			s = fmt.Sprintf("%s NULL", s)
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
			if ope2.Int == 0 && ope2.Float == 0.0 {
				return col, fmt.Errorf("Div by 0")
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				v := value.Value{
					Type: value.INTEGER,
					Int:  ope1.Int / ope2.Int,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float / ope2.Float,
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				v := value.Value{
					Type:  value.FLOAT,
					Float: ope1.Float / float64(ope2.Int),
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
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

		case EQ:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int == ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float == ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float == float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) == ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.STRING && ope2.Type == value.STRING {
				if ope1.String == ope2.String {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.BOOL && ope2.Type == value.BOOL {
				if (ope1.Bool.True == true && ope2.Bool.True == true) || (ope1.Bool.False == true && ope2.Bool.False) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope2.Type == value.NULL {
				if ope1.Type == value.NULL {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s = %s", ope1.Type, ope2.Type)
			}

		case NEQ:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int != ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float != ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float != float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) != ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.STRING && ope2.Type == value.STRING {
				if ope1.String != ope2.String {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope2.Type == value.NULL {
				if ope1.Type == value.NULL {
					v.Bool.False = true
				} else {
					v.Bool.True = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s <> %s", ope1.Type, ope2.Type)
			}

		case LT:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int < ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float < ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float < float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) < ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s < %s", ope1.Type, ope2.Type)
			}

		case LTE:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int <= ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float <= ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float <= float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) <= ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s <= %s", ope1.Type, ope2.Type)
			}

		case GT:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int > ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float > ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float > float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) > ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s > %s", ope1.Type, ope2.Type)
			}

		case GTE:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.INTEGER && ope2.Type == value.INTEGER {
				if ope1.Int >= ope2.Int {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.FLOAT {
				if ope1.Float >= ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.FLOAT && ope2.Type == value.INTEGER {
				if ope1.Float >= float64(ope2.Int) {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else if ope1.Type == value.INTEGER && ope2.Type == value.FLOAT {
				if float64(ope1.Int) >= ope2.Float {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s >= %s", ope1.Type, ope2.Type)
			}

		case AND:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.BOOL && ope2.Type == value.BOOL {
				if ope1.Bool.True == true && ope2.Bool.True == true {
					v.Bool.True = true
				} else {
					v.Bool.False = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s AND %s", ope1.Type, ope2.Type)
			}

		case OR:
			ope2, err := s.pop()
			if err != nil {
				return col, err
			}
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.BOOL && ope2.Type == value.BOOL {
				if ope1.Bool.False == true && ope2.Bool.False == true {
					v.Bool.False = true
				} else {
					v.Bool.True = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: %s OR %s", ope1.Type, ope2.Type)
			}

		case NOT:
			ope1, err := s.pop()
			if err != nil {
				return col, err
			}
			v := value.Value{
				Type: value.BOOL,
				Bool: value.Bool{},
			}
			if ope1.Type == value.BOOL {
				if ope1.Bool.True == true {
					v.Bool.False = true
				}
				if ope1.Bool.False == true {
					v.Bool.True = true
				}
				s.push(v)
			} else {
				return col, fmt.Errorf("Unknown Operation: NOT %s", ope1.Type)
			}

		case CAST:
			target, err := s.pop()
			if err != nil {
				return col, err
			}

			switch code.Operand1.Type {
			case value.INTEGER:
				if target.Type == value.FLOAT {
					target.Type = value.INTEGER
					target.Int = int64(target.Float)
				} else if target.Type == value.STRING {
					val, err := value.Convert(target.String)
					if err != nil {
						return col, err
					}
					if val.Type == value.INTEGER {
						target.Type = value.INTEGER
						target.Int = val.Int
					} else if val.Type == value.FLOAT {
						target.Type = value.INTEGER
						target.Int = int64(val.Float)
					} else {
						return col, fmt.Errorf("Cannot cast")
					}
				}
			case value.FLOAT:
				if target.Type == value.INTEGER {
					target.Type = value.FLOAT
					target.Float = float64(target.Int)
				} else if target.Type == value.STRING {
					val, err := value.Convert(target.String)
					if err != nil {
						return col, err
					}
					if val.Type == value.INTEGER {
						target.Type = value.FLOAT
						target.Float = float64(val.Int)
					} else if val.Type == value.FLOAT {
						target.Type = value.FLOAT
						target.Float = val.Float
					} else {
						return col, fmt.Errorf("Cannot cast")
					}
				}
			case value.STRING:
				if target.Type == value.FLOAT {
					target.Type = value.STRING
					target.String = fmt.Sprintf("%f", target.Float)
				} else if target.Type == value.INTEGER {
					target.Type = value.STRING
					target.String = fmt.Sprintf("%d", target.Int)
				}
			}
			s.push(target)

		case NA:
			panic("")
		}
	}
	return col, nil
}

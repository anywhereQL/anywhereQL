package planner

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/vm"
	"github.com/google/uuid"
)

func Translate(expr *ast.Expression) []vm.ExprVMCode {
	codes := translateExpr(expr)
	codes = append(codes, vm.ExprVMCode{Operator: vm.STORE})
	return codes
}

func translateExpr(expr *ast.Expression) []vm.ExprVMCode {
	codes := []vm.ExprVMCode{}
	v := value.Value{}
	if expr.Literal != nil {
		if expr.Literal.Numeric != nil {
			switch expr.Literal.Numeric.Type {
			case ast.N_INT:
				v.Type = value.INTEGER
				v.Int = expr.Literal.Numeric.Integral
			case ast.N_FLOAT:
				v.Type = value.FLOAT
				v.Float = expr.Literal.Numeric.Float
				v.PartF = expr.Literal.Numeric.PartF
				v.PartI = expr.Literal.Numeric.PartI
				v.FDigit = expr.Literal.Numeric.FDigit
			}
		} else if expr.Literal.String != nil {
			v.Type = value.STRING
			v.String = expr.Literal.String.Value
		} else if expr.Literal.Bool != nil {
			v.Type = value.BOOL
			if expr.Literal.Bool.True {
				v.Bool.True = true
			} else if expr.Literal.Bool.False {
				v.Bool.False = true
			}
		} else if expr.Literal.NULL == true {
			v.Type = value.NULL
		}

		c := vm.ExprVMCode{
			Operator: vm.PUSH,
			Operand:  v,
		}
		codes = append(codes, c)
	} else if expr.BinaryOperation != nil {
		cl := translateExpr(expr.BinaryOperation.Left)
		codes = append(codes, cl...)
		cr := translateExpr(expr.BinaryOperation.Right)
		codes = append(codes, cr...)

		var c vm.ExprVMCode
		switch expr.BinaryOperation.Operator {
		case ast.B_PLUS:
			c = vm.ExprVMCode{
				Operator: vm.ADD,
			}
		case ast.B_MINUS:
			c = vm.ExprVMCode{
				Operator: vm.SUB,
			}
		case ast.B_ASTERISK:
			c = vm.ExprVMCode{
				Operator: vm.MUL,
			}
		case ast.B_SOLIDAS:
			c = vm.ExprVMCode{
				Operator: vm.DIV,
			}
		case ast.B_PERCENT:
			c = vm.ExprVMCode{
				Operator: vm.MOD,
			}
		case ast.B_EQUAL:
			c = vm.ExprVMCode{
				Operator: vm.EQ,
			}
		case ast.B_NOT_EQUAL:
			c = vm.ExprVMCode{
				Operator: vm.NEQ,
			}
		case ast.B_GREATER_THAN:
			c = vm.ExprVMCode{
				Operator: vm.GT,
			}
		case ast.B_GREATER_THAN_EQUAL:
			c = vm.ExprVMCode{
				Operator: vm.GTE,
			}
		case ast.B_LESS_THAN:
			c = vm.ExprVMCode{
				Operator: vm.LT,
			}
		case ast.B_LESS_THAN_EQUAL:
			c = vm.ExprVMCode{
				Operator: vm.LTE,
			}
		case ast.B_AND:
			c = vm.ExprVMCode{
				Operator: vm.AND,
			}
		case ast.B_OR:
			c = vm.ExprVMCode{
				Operator: vm.OR,
			}

		default:
			return codes
		}
		codes = append(codes, c)
	} else if expr.UnaryOperation != nil {
		c := translateExpr(expr.UnaryOperation.Expr)
		codes = append(codes, c...)
		switch expr.UnaryOperation.Operator {
		case ast.U_MINUS:
			codes = append(codes, vm.ExprVMCode{Operator: vm.PUSH, Operand: value.Value{Type: value.INTEGER, Int: -1}})
			codes = append(codes, vm.ExprVMCode{Operator: vm.MUL})
		case ast.U_NOT:
			codes = append(codes, vm.ExprVMCode{Operator: vm.NOT})
		}
	} else if expr.FunctionCall != nil {
		for _, arg := range expr.FunctionCall.Args {
			c := translateExpr(&arg)
			codes = append(codes, c...)
		}
		codes = append(codes, vm.ExprVMCode{Operator: vm.PUSH, Operand: value.Value{Type: value.INTEGER, Int: int64(len(expr.FunctionCall.Args))}})
		codes = append(codes, vm.ExprVMCode{Operator: vm.CALL, Operand: value.Value{Type: value.STRING, String: expr.FunctionCall.Name}})
	} else if expr.Column != nil {
		v := value.Value{
			Type: value.COLUMN,
			Column: value.Column{
				Column:  expr.Column.Column,
				TableID: expr.Column.Table.ID,
			},
		}
		codes = append(codes, vm.ExprVMCode{Operator: vm.PICK, Operand: v})
	} else if expr.Cast != nil {
		v := value.Value{}
		c := translateExpr(expr.Cast.Expr)
		codes = append(codes, c...)
		switch expr.Cast.Type {
		case ast.T_INT:
			v.Type = value.INTEGER
		case ast.T_FLOAT:
			v.Type = value.FLOAT
		case ast.T_STRING:
			v.Type = value.STRING
		}
		codes = append(codes, vm.ExprVMCode{Operator: vm.CAST, Operand: v})
	} else if expr.Case != nil {
		endMark := uuid.New()
		if expr.Case.Value != nil {
			for _, ca := range expr.Case.CaseValues {
				exp := &ast.Expression{
					BinaryOperation: &ast.BinaryOpe{
						Left:     expr.Case.Value,
						Right:    ca.Condition,
						Operator: ast.B_EQUAL,
					},
				}
				c := translateExpr(exp)
				codes = append(codes, c...)
				tOpe := translateExpr(ca.Result)
				tOpe = append(tOpe, vm.ExprVMCode{Operator: vm.JMPL, Operand: value.Value{
					Type:   value.STRING,
					String: endMark.String(),
				}})

				falseOpe := vm.ExprVMCode{
					Operator: vm.JMPNC,
					Operand: value.Value{
						Type: value.INTEGER,
						Int:  int64(len(tOpe)),
					},
				}

				codes = append(codes, falseOpe)
				codes = append(codes, tOpe...)
			}
		} else {
			for _, ca := range expr.Case.CaseValues {
				c := translateExpr(ca.Condition)
				r := translateExpr(ca.Result)
				r = append(r, vm.ExprVMCode{Operator: vm.JMPL, Operand: value.Value{
					Type:   value.STRING,
					String: endMark.String(),
				}})

				codes = append(codes, c...)
				codes = append(codes, vm.ExprVMCode{Operator: vm.JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(r))}})
				codes = append(codes, r...)
			}
		}

		if expr.Case.ElseValue != nil {
			c := translateExpr(expr.Case.ElseValue)
			codes = append(codes, c...)
		} else {
			codes = append(codes, vm.ExprVMCode{Operator: vm.PUSH, Operand: value.Value{
				Type: value.NULL,
			}})
		}
		codes = append(codes, vm.ExprVMCode{Operator: vm.LABEL, Operand: value.Value{
			Type:   value.STRING,
			String: endMark.String(),
		}})
	} else if expr.Between != nil {
		b := translateExpr(expr.Between.Src)
		b = append(b, translateExpr(expr.Between.Begin)...)
		b = append(b, vm.ExprVMCode{Operator: vm.GTE})

		e := translateExpr(expr.Between.Src)
		e = append(e, translateExpr(expr.Between.End)...)
		e = append(e, vm.ExprVMCode{Operator: vm.LTE})

		t := []vm.ExprVMCode{}
		t = append(t, vm.ExprVMCode{Operator: vm.PUSH, Operand: value.Value{Type: value.BOOL, Bool: value.Bool{True: true}}})
		t = append(t, vm.ExprVMCode{Operator: vm.JMP, Operand: value.Value{Type: value.INTEGER, Int: 1}})

		f := []vm.ExprVMCode{}
		f = append(f, vm.ExprVMCode{Operator: vm.PUSH, Operand: value.Value{Type: value.BOOL, Bool: value.Bool{False: true}}})

		codes = append(codes, b...)
		codes = append(codes, vm.ExprVMCode{Operator: vm.JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(e) + len(t) + 1)}})
		codes = append(codes, e...)
		codes = append(codes, vm.ExprVMCode{Operator: vm.JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(t))}})
		codes = append(codes, t...)
		codes = append(codes, f...)
		if expr.Between.Not {
			codes = append(codes, vm.ExprVMCode{Operator: vm.SWAP})
		}
	}
	return codes
}

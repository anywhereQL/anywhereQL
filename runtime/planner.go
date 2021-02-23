package runtime

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/google/uuid"
)

func (r *Runtime) Translate(expr *ast.Expression) []ExprVMCode {
	codes := r.translateExpr(expr)
	codes = append(codes, ExprVMCode{Operator: STORE})
	return codes
}

func (r *Runtime) translateExpr(expr *ast.Expression) []ExprVMCode {
	codes := []ExprVMCode{}
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

		c := ExprVMCode{
			Operator: PUSH,
			Operand:  v,
		}
		codes = append(codes, c)
	} else if expr.BinaryOperation != nil {
		cl := r.translateExpr(expr.BinaryOperation.Left)
		codes = append(codes, cl...)
		cr := r.translateExpr(expr.BinaryOperation.Right)
		codes = append(codes, cr...)

		var c ExprVMCode
		switch expr.BinaryOperation.Operator {
		case ast.B_PLUS:
			c = ExprVMCode{
				Operator: ADD,
			}
		case ast.B_MINUS:
			c = ExprVMCode{
				Operator: SUB,
			}
		case ast.B_ASTERISK:
			c = ExprVMCode{
				Operator: MUL,
			}
		case ast.B_SOLIDAS:
			c = ExprVMCode{
				Operator: DIV,
			}
		case ast.B_PERCENT:
			c = ExprVMCode{
				Operator: MOD,
			}
		case ast.B_EQUAL:
			c = ExprVMCode{
				Operator: EQ,
			}
		case ast.B_NOT_EQUAL:
			c = ExprVMCode{
				Operator: NEQ,
			}
		case ast.B_GREATER_THAN:
			c = ExprVMCode{
				Operator: GT,
			}
		case ast.B_GREATER_THAN_EQUAL:
			c = ExprVMCode{
				Operator: GTE,
			}
		case ast.B_LESS_THAN:
			c = ExprVMCode{
				Operator: LT,
			}
		case ast.B_LESS_THAN_EQUAL:
			c = ExprVMCode{
				Operator: LTE,
			}
		case ast.B_AND:
			c = ExprVMCode{
				Operator: AND,
			}
		case ast.B_OR:
			c = ExprVMCode{
				Operator: OR,
			}

		default:
			return codes
		}
		codes = append(codes, c)
	} else if expr.UnaryOperation != nil {
		c := r.translateExpr(expr.UnaryOperation.Expr)
		codes = append(codes, c...)
		switch expr.UnaryOperation.Operator {
		case ast.U_MINUS:
			codes = append(codes, ExprVMCode{Operator: PUSH, Operand: value.Value{Type: value.INTEGER, Int: -1}})
			codes = append(codes, ExprVMCode{Operator: MUL})
		case ast.U_NOT:
			codes = append(codes, ExprVMCode{Operator: NOT})
		}
	} else if expr.FunctionCall != nil {
		for _, arg := range expr.FunctionCall.Args {
			c := r.translateExpr(&arg)
			codes = append(codes, c...)
		}
		codes = append(codes, ExprVMCode{Operator: PUSH, Operand: value.Value{Type: value.INTEGER, Int: int64(len(expr.FunctionCall.Args))}})
		codes = append(codes, ExprVMCode{Operator: CALL, Operand: value.Value{Type: value.STRING, String: expr.FunctionCall.Name}})
	} else if expr.Column != nil {
		v := value.Value{
			Type: value.COLUMN,
			Column: value.Column{
				Schema: expr.Column.Table.Schema,
				DB:     expr.Column.Table.DB,
				Table:  expr.Column.Table.Table,
				Column: expr.Column.Column,
			},
		}
		codes = append(codes, ExprVMCode{Operator: PICK, Operand: v})
	} else if expr.Cast != nil {
		v := value.Value{}
		c := r.translateExpr(expr.Cast.Expr)
		codes = append(codes, c...)
		switch expr.Cast.Type {
		case ast.T_INT:
			v.Type = value.INTEGER
		case ast.T_FLOAT:
			v.Type = value.FLOAT
		case ast.T_STRING:
			v.Type = value.STRING
		}
		codes = append(codes, ExprVMCode{Operator: CAST, Operand: v})
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
				c := r.translateExpr(exp)
				codes = append(codes, c...)
				tOpe := r.translateExpr(ca.Result)
				tOpe = append(tOpe, ExprVMCode{Operator: JMPL, Operand: value.Value{
					Type:   value.STRING,
					String: endMark.String(),
				}})

				falseOpe := ExprVMCode{
					Operator: JMPNC,
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
				c := r.translateExpr(ca.Condition)
				r := r.translateExpr(ca.Result)
				r = append(r, ExprVMCode{Operator: JMPL, Operand: value.Value{
					Type:   value.STRING,
					String: endMark.String(),
				}})

				codes = append(codes, c...)
				codes = append(codes, ExprVMCode{Operator: JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(r))}})
				codes = append(codes, r...)
			}
		}

		if expr.Case.ElseValue != nil {
			c := r.translateExpr(expr.Case.ElseValue)
			codes = append(codes, c...)
		} else {
			codes = append(codes, ExprVMCode{Operator: PUSH, Operand: value.Value{
				Type: value.NULL,
			}})
		}
		codes = append(codes, ExprVMCode{Operator: LABEL, Operand: value.Value{
			Type:   value.STRING,
			String: endMark.String(),
		}})
	} else if expr.Between != nil {
		b := r.translateExpr(expr.Between.Src)
		b = append(b, r.translateExpr(expr.Between.Begin)...)
		b = append(b, ExprVMCode{Operator: GTE})

		e := r.translateExpr(expr.Between.Src)
		e = append(e, r.translateExpr(expr.Between.End)...)
		e = append(e, ExprVMCode{Operator: LTE})

		t := []ExprVMCode{}
		t = append(t, ExprVMCode{Operator: PUSH, Operand: value.Value{Type: value.BOOL, Bool: value.Bool{True: true}}})
		t = append(t, ExprVMCode{Operator: JMP, Operand: value.Value{Type: value.INTEGER, Int: 1}})

		f := []ExprVMCode{}
		f = append(f, ExprVMCode{Operator: PUSH, Operand: value.Value{Type: value.BOOL, Bool: value.Bool{False: true}}})

		codes = append(codes, b...)
		codes = append(codes, ExprVMCode{Operator: JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(e) + len(t) + 1)}})
		codes = append(codes, e...)
		codes = append(codes, ExprVMCode{Operator: JMPNC, Operand: value.Value{Type: value.INTEGER, Int: int64(len(t))}})
		codes = append(codes, t...)
		codes = append(codes, f...)
		if expr.Between.Not {
			codes = append(codes, ExprVMCode{Operator: SWAP})
		}
	}
	return codes
}

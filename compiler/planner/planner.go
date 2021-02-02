package planner

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/vm"
)

func Translate(a *ast.AST) []vm.VMCode {
	codes := []vm.VMCode{}

	for _, sql := range a.SQL {
		for _, col := range sql.SELECTStatement.SELECT.SelectColumns {
			c := translateSelectColumn(col)
			codes = append(codes, c...)

			s := vm.VMCode{
				Operator: vm.STORE,
				Operand1: value.Value{
					Type: value.NA,
				},
			}
			codes = append(codes, s)
		}
	}
	return codes
}

func translateSelectColumn(c ast.SelectColumn) []vm.VMCode {
	codes := translateExpression(c.Expression)
	return codes
}

func translateExpression(expr *ast.Expression) []vm.VMCode {
	codes := []vm.VMCode{}
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
			c := vm.VMCode{
				Operator: vm.PUSH,
				Operand1: v,
			}
			codes = append(codes, c)
		}
		return codes
	} else if expr.BinaryOperation != nil {
		cl := translateExpression(expr.BinaryOperation.Left)
		codes = append(codes, cl...)
		cr := translateExpression(expr.BinaryOperation.Right)
		codes = append(codes, cr...)

		var c vm.VMCode
		switch expr.BinaryOperation.Operator {
		case ast.B_PLUS:
			c = vm.VMCode{
				Operator: vm.ADD,
			}
		case ast.B_MINUS:
			c = vm.VMCode{
				Operator: vm.SUB,
			}
		case ast.B_ASTERISK:
			c = vm.VMCode{
				Operator: vm.MUL,
			}
		case ast.B_SOLIDAS:
			c = vm.VMCode{
				Operator: vm.DIV,
			}
		case ast.B_PERCENT:
			c = vm.VMCode{
				Operator: vm.MOD,
			}
		default:
			return codes
		}
		codes = append(codes, c)
		return codes
	} else if expr.UnaryOperation != nil {
		c := translateExpression(expr.UnaryOperation.Expr)
		codes = append(codes, c...)
		switch expr.UnaryOperation.Operator {
		case ast.U_MINUS:
			codes = append(codes, vm.VMCode{Operator: vm.PUSH, Operand1: value.Value{Type: value.INTEGER, Int: -1}})
			codes = append(codes, vm.VMCode{Operator: vm.MUL})
		}
		return codes
	} else if expr.FunctionCall != nil {
		for _, arg := range expr.FunctionCall.Args {
			c := translateExpression(&arg)
			codes = append(codes, c...)
		}
		codes = append(codes, vm.VMCode{Operator: vm.PUSH, Operand1: value.Value{Type: value.INTEGER, Int: int64(len(expr.FunctionCall.Args))}})
		codes = append(codes, vm.VMCode{Operator: vm.CALL, Operand1: value.Value{Type: value.STRING, String: expr.FunctionCall.Name}})
		return codes
	}
	return codes
}

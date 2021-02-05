package debug

import (
	"fmt"
	"io"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
	"github.com/anywhereQL/anywhereQL/runtime/vm"
)

func PrintToken(out io.Writer, tokens token.Tokens) {
	for n, token := range tokens {
		fmt.Fprintf(out, "[%d] Type: %s Literal: %s\n", n, token.Type, token.Literal)
	}
}

func PrintVC(out io.Writer, vc []vm.VMCode) {
	for n, c := range vc {
		fmt.Fprintf(out, "[%d] %s\n", n, c)
	}
}

func PrintAST(out io.Writer, ast *ast.AST) {
	for _, s := range ast.SQL {
		if s.SELECTStatement != nil {
			printSELECT(out, s.SELECTStatement)
		}
	}
}

func printSELECT(out io.Writer, ss *ast.SELECTStatement) {
	if ss.SELECT != nil {
		printSELECTColumns(out, ss.SELECT)
	}
	if ss.FROM != nil {
	}
}

func printSELECTColumns(out io.Writer, cols *ast.SELECTClause) {
	for _, col := range cols.SelectColumns {
		printColumn(out, col)
	}
}

func printColumn(out io.Writer, col ast.SelectColumn) {
	if col.Expression != nil {
		printExpression(out, "", col.Expression)
	}
}

func printExpression(out io.Writer, sep string, expr *ast.Expression) {
	if expr.UnaryOperation != nil {
		fmt.Fprintf(out, "%sUnary Operation (%s)\n", sep, expr.UnaryOperation.Operator)
		printExpression(out, sep+" ", expr.UnaryOperation.Expr)
	}
	if expr.BinaryOperation != nil {
		fmt.Fprintf(out, "%sBinary Operation (%s)\n", sep, expr.BinaryOperation.Operator)
		fmt.Fprintf(out, "%sLeft:\n", sep)
		printExpression(out, sep+" ", expr.BinaryOperation.Left)
		fmt.Fprintf(out, "%sRight:\n", sep)
		printExpression(out, sep+" ", expr.BinaryOperation.Right)
	}
	if expr.Literal != nil {
		if expr.Literal.Numeric != nil {
			switch expr.Literal.Numeric.Type {
			case ast.N_INT:
				fmt.Fprintf(out, "%s%d\n", sep, expr.Literal.Numeric.Integral)
			case ast.N_FLOAT:
				fmt.Fprintf(out, "%s%f\n", sep, expr.Literal.Numeric.Float)
			}
		} else if expr.Literal.String != nil {
			fmt.Fprintf(out, "%s%s", sep, expr.Literal.String.Value)
		}
	}
	if expr.Column != nil {
		if expr.Column.Schema != "" {
			fmt.Fprintf(out, "%s:Schema: %s", sep, expr.Column.Schema)
		}
		if expr.Column.DB != "" {
			fmt.Fprintf(out, "%s:DB: %s", sep, expr.Column.DB)
		}
		if expr.Column.Table != "" {
			fmt.Fprintf(out, "%s:Table: %s", sep, expr.Column.Table)
		}
		if expr.Column.Column != "" {
			fmt.Fprintf(out, "%s:Column: %s", sep, expr.Column.Column)
		}
	}
}

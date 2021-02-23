package debug

import (
	"fmt"
	"io"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
)

func PrintToken(out io.Writer, tokens token.Tokens) {
	for n, token := range tokens {
		fmt.Fprintf(out, "[%d] Type: %s Literal: %s\n", n, token.Type, token.Literal)
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
		printFROMClause(out, ss.FROM)
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

func printFROMClause(out io.Writer, from *ast.FROMClause) {
	printTable(out, from.Table)
	for _, tbl := range from.Joined {
		printTable(out, tbl.Table)
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
		} else if expr.Literal.NULL == true {
			fmt.Fprintf(out, "%sNULL", sep)
		} else if expr.Literal.Bool != nil {
			if expr.Literal.Bool.True {
				fmt.Fprintf(out, "%sTRUE", sep)
			}
			if expr.Literal.Bool.False {
				fmt.Fprintf(out, "%sFALSE", sep)
			}
		}
	}
	if expr.Column != nil {
		fmt.Fprintf(out, "%sColumn (schema: %s, database: %s, table: %s, column: %s)\n", sep, expr.Column.Table.Schema, expr.Column.Table.DB, expr.Column.Table.Table, expr.Column.Column)
	}
}

func printTable(out io.Writer, tbl *ast.Table) {
	fmt.Fprintf(out, "Table (schema:%s, database: %s, table: %s)\n", tbl.Schema, tbl.DB, tbl.Table)
}

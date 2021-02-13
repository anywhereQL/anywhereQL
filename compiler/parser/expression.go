package parser

import (
	"fmt"
	"strings"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
	"github.com/anywhereQL/anywhereQL/common/value"
)

func (p *parser) parseExpression(precedence int) (*ast.Expression, error) {
	expr := &ast.Expression{}
	unary, exists := p.unaryParseFunc[p.currentToken.Type]
	if !exists {
		return expr, fmt.Errorf("Unknown Unary Operator: %s", p.currentToken.Literal)
	}

	left, err := unary()
	if err != nil {
		return expr, err
	}

	for !(p.getNextToken().Type == token.EOS || p.getNextToken().Type == token.K_FROM) && precedence < p.getNextTokenPrecedence() {
		binary, exists := p.binaryParseFunc[p.getNextToken().Type]
		if !exists {
			return left, nil
		}
		p.readToken()
		left, err = binary(left)
		if err != nil {
			return left, err
		}
	}
	return left, nil
}

func (p *parser) parseNumber() (*ast.Expression, error) {
	expr := &ast.Expression{
		Literal: &ast.Literal{
			Numeric: &ast.Numeric{
				Integral: p.currentToken.Value.Int,
				Float:    p.currentToken.Value.Float,
				FDigit:   p.currentToken.Value.FDigit,
				PartF:    p.currentToken.Value.PartF,
				PartI:    p.currentToken.Value.PartI,
			},
		},
	}
	switch p.currentToken.Value.Type {
	case value.INTEGER:
		expr.Literal.Numeric.Type = ast.N_INT
	case value.FLOAT:
		expr.Literal.Numeric.Type = ast.N_FLOAT
	default:
		return expr, fmt.Errorf("Unknwon Value Type: %s", p.currentToken.Type)
	}
	return expr, nil
}

func (p *parser) parseIdent() (*ast.Expression, error) {
	if p.getNextToken().Type == token.S_LPAREN {
		return p.parseFunctionCallExpr()
	}
	return p.parseColumnExpr()
}

func (p *parser) parseBinaryExpr(left *ast.Expression) (*ast.Expression, error) {
	expr := &ast.Expression{
		BinaryOperation: &ast.BinaryOpe{
			Left: left,
		},
	}

	switch p.currentToken.Type {
	case token.S_PLUS:
		expr.BinaryOperation.Operator = ast.B_PLUS
	case token.S_MINUS:
		expr.BinaryOperation.Operator = ast.B_MINUS
	case token.S_ASTERISK:
		expr.BinaryOperation.Operator = ast.B_ASTERISK
	case token.S_SOLIDAS:
		expr.BinaryOperation.Operator = ast.B_SOLIDAS
	case token.S_PERCENT:
		expr.BinaryOperation.Operator = ast.B_PERCENT
	case token.S_EQUAL:
		expr.BinaryOperation.Operator = ast.B_EQUAL
	case token.S_NOT_EQUAL:
		expr.BinaryOperation.Operator = ast.B_NOT_EQUAL
	case token.S_GREATER_THAN:
		expr.BinaryOperation.Operator = ast.B_GREATER_THAN
	case token.S_GREATER_THAN_EQUAL:
		expr.BinaryOperation.Operator = ast.B_GREATER_THAN_EQUAL
	case token.S_LESS_THAN:
		expr.BinaryOperation.Operator = ast.B_LESS_THAN
	case token.S_LESS_THAN_EQUAL:
		expr.BinaryOperation.Operator = ast.B_LESS_THAN_EQUAL
	case token.K_AND:
		expr.BinaryOperation.Operator = ast.B_AND
	case token.K_OR:
		expr.BinaryOperation.Operator = ast.B_OR

	default:
		return expr, fmt.Errorf("Unknown Binary Operator: %s", p.currentToken.Literal)
	}
	precedences := p.getCurrentTokenPrecedence()

	p.readToken()

	ex, err := p.parseExpression(precedences)
	if err != nil {
		return expr, err
	}
	expr.BinaryOperation.Right = ex
	return expr, nil
}

func (p *parser) parseGroupedExpr() (*ast.Expression, error) {
	expr := &ast.Expression{}
	p.readToken()

	ex, err := p.parseExpression(LOWEST)
	if err != nil {
		return expr, err
	}
	if p.getNextToken().Type != token.S_RPAREN {
		return expr, fmt.Errorf("Unknwon Grouped Expression")
	}
	p.readToken()
	return ex, nil
}

func (p *parser) parsePrefixExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		UnaryOperation: &ast.UnaryOpe{},
	}
	switch p.currentToken.Type {
	case token.S_PLUS:
		expr.UnaryOperation.Operator = ast.U_PLUS
	case token.S_MINUS:
		expr.UnaryOperation.Operator = ast.U_MINUS
	case token.K_NOT:
		expr.UnaryOperation.Operator = ast.U_NOT
	default:
		return expr, fmt.Errorf("Unknwon Prefix Operator: %s", p.currentToken.Literal)
	}

	pre := p.getCurrentTokenPrecedence()
	p.readToken()
	ex, err := p.parseExpression(pre)
	if err != nil {
		return expr, err
	}

	expr.UnaryOperation.Expr = ex

	return expr, nil
}

func (p *parser) parseFunctionCallExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		FunctionCall: &ast.FunctionCall{},
	}

	expr.FunctionCall.Name = strings.ToUpper(p.currentToken.Literal)
	p.readToken()

	for {
		p.readToken()
		if p.currentToken.Type == token.S_RPAREN {
			break
		}
		ex, err := p.parseExpression(LOWEST)
		if err != nil {
			return expr, err
		}
		expr.FunctionCall.Args = append(expr.FunctionCall.Args, *ex)

		p.readToken()
		if p.currentToken.Type == token.S_RPAREN {
			break
		}
		if p.currentToken.Type != token.S_COMMA {
			return expr, fmt.Errorf("Unknwon Token must be COMMA or RPAREN, but %s", p.currentToken.Literal)
		}
	}
	return expr, nil
}

func (p *parser) parseString() (*ast.Expression, error) {
	expr := &ast.Expression{
		Literal: &ast.Literal{
			String: &ast.String{
				Value: p.currentToken.Value.String,
			},
		},
	}
	return expr, nil
}

func (p *parser) parseColumnExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		Column: &ast.Column{},
	}
	literal := []string{}
	for {
		if p.currentToken.Type == token.IDENT {
			literal = append(literal, p.currentToken.Literal)
		}
		if p.getNextToken().Type != token.S_PERIOD {
			break
		}
		p.readToken()
	}
	if len(literal) >= 5 {
		return expr, fmt.Errorf("Unknown Column Expression")
	} else if len(literal) == 4 {
		expr.Column.Table.Schema = literal[0]
		expr.Column.Table.DB = literal[1]
		expr.Column.Table.Table = literal[2]
		expr.Column.Column = literal[3]
	} else if len(literal) == 3 {
		expr.Column.Table.DB = literal[0]
		expr.Column.Table.Table = literal[1]
		expr.Column.Column = literal[2]
	} else if len(literal) == 2 {
		expr.Column.Table.Table = literal[0]
		expr.Column.Column = literal[1]
	} else {
		expr.Column.Column = literal[0]
	}

	return expr, nil
}

func (p *parser) parseBoolExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		Literal: &ast.Literal{
			Bool: &ast.Bool{},
		},
	}
	if p.currentToken.Type == token.K_TRUE {
		expr.Literal.Bool.True = true
	} else if p.currentToken.Type == token.K_FALSE {
		expr.Literal.Bool.False = true
	}
	return expr, nil
}

func (p *parser) parseIsExpr(left *ast.Expression) (*ast.Expression, error) {
	expr := &ast.Expression{
		BinaryOperation: &ast.BinaryOpe{
			Left: left,
			Right: &ast.Expression{
				Literal: &ast.Literal{
					NULL: true,
				},
			},
		},
	}

	if p.currentToken.Type == token.K_IS {
		if p.getNextToken().Type == token.K_NULL {
			p.readToken()
			expr.BinaryOperation.Operator = ast.B_EQUAL
		} else {
			if p.getNextToken().Type == token.K_NOT {
				p.readToken()
				if p.getNextToken().Type != token.K_NULL {
					return expr, fmt.Errorf("Unknown IS Expr")
				}
				p.readToken()
				expr.BinaryOperation.Operator = ast.B_NOT_EQUAL
			}
		}
	} else if p.currentToken.Type == token.K_ISNULL {
		expr.BinaryOperation.Operator = ast.B_EQUAL
	} else if p.currentToken.Type == token.K_NOTNULL {
		expr.BinaryOperation.Operator = ast.B_NOT_EQUAL
	}
	return expr, nil
}

func (p *parser) parseNullExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		Literal: &ast.Literal{
			NULL: true,
		},
	}
	return expr, nil
}

func (p *parser) parseCastExpr() (*ast.Expression, error) {
	expr := &ast.Expression{
		Cast: &ast.Cast{},
	}
	p.readToken()
	if p.currentToken.Type != token.S_LPAREN {
		return nil, fmt.Errorf("Unknown Cast format")
	}
	p.readToken()
	ex, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	expr.Cast.Expr = ex
	p.readToken()
	if p.currentToken.Type != token.K_AS {
		return nil, fmt.Errorf("Unknown Cast format")
	}
	p.readToken()

	switch p.currentToken.Type {
	case token.K_INT, token.K_INTEGER:
		expr.Cast.Type = ast.T_INT
	case token.K_FLOAT, token.K_DOUBLE:
		expr.Cast.Type = ast.T_FLOAT
	case token.K_STRING:
		expr.Cast.Type = ast.T_STRING
	default:
		return nil, fmt.Errorf("Unknown cast type")
	}
	p.readToken()
	if p.currentToken.Type != token.S_RPAREN {
		return nil, fmt.Errorf("Unknown Cast format")
	}

	return expr, nil
}

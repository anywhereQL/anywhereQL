package parser

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
)

type (
	unaryOpeFunction  func() (*ast.Expression, error)
	binaryOpeFunction func(*ast.Expression) (*ast.Expression, error)
)

const (
	_ int = iota
	LOWEST
	AND_OR
	NOT
	COMPARE
	SUM
	PRODUCT
	HIGHEST
)

var precedences = map[token.Type]int{
	token.S_PLUS:               SUM,
	token.S_MINUS:              SUM,
	token.S_ASTERISK:           PRODUCT,
	token.S_SOLIDAS:            PRODUCT,
	token.S_PERCENT:            PRODUCT,
	token.S_EQUAL:              COMPARE,
	token.S_NOT_EQUAL:          COMPARE,
	token.S_GREATER_THAN:       COMPARE,
	token.S_GREATER_THAN_EQUAL: COMPARE,
	token.S_LESS_THAN:          COMPARE,
	token.S_LESS_THAN_EQUAL:    COMPARE,
	token.K_NOT:                NOT,
	token.K_AND:                AND_OR,
	token.K_OR:                 AND_OR,
	token.K_IS:                 COMPARE,
	token.K_ISNULL:             COMPARE,
	token.K_NOTNULL:            COMPARE,
	token.K_BETWEEN:            COMPARE,
	token.K_IN:                 COMPARE,
}

type parser struct {
	tokens       token.Tokens
	currentToken token.Token
	pos          int

	unaryParseFunc  map[token.Type]unaryOpeFunction
	binaryParseFunc map[token.Type]binaryOpeFunction
}

func Parse(tokens token.Tokens) (*ast.AST, error) {
	a := &ast.AST{}
	p := new(tokens)
	sql, err := p.parse()
	if err != nil {
		return a, err
	}
	a.SQL = sql
	return a, nil
}

func new(tokens token.Tokens) *parser {
	p := &parser{
		tokens: tokens,
	}
	p.readToken()

	p.unaryParseFunc = make(map[token.Type]unaryOpeFunction)
	p.binaryParseFunc = make(map[token.Type]binaryOpeFunction)

	p.unaryParseFunc[token.NUMBER] = p.parseNumber
	p.unaryParseFunc[token.S_LPAREN] = p.parseGroupedExpr
	p.unaryParseFunc[token.S_PLUS] = p.parsePrefixExpr
	p.unaryParseFunc[token.S_MINUS] = p.parsePrefixExpr
	p.unaryParseFunc[token.IDENT] = p.parseIdent
	p.unaryParseFunc[token.STRING] = p.parseString
	p.unaryParseFunc[token.K_TRUE] = p.parseBoolExpr
	p.unaryParseFunc[token.K_FALSE] = p.parseBoolExpr
	p.unaryParseFunc[token.K_NOT] = p.parsePrefixExpr
	p.unaryParseFunc[token.K_NULL] = p.parseNullExpr
	p.unaryParseFunc[token.K_CAST] = p.parseCastExpr
	p.unaryParseFunc[token.K_CASE] = p.parseCaseExpr

	p.binaryParseFunc[token.S_PLUS] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_MINUS] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_ASTERISK] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_SOLIDAS] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_PERCENT] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_EQUAL] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_NOT_EQUAL] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_GREATER_THAN] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_GREATER_THAN_EQUAL] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_LESS_THAN] = p.parseBinaryExpr
	p.binaryParseFunc[token.S_LESS_THAN_EQUAL] = p.parseBinaryExpr
	p.binaryParseFunc[token.K_AND] = p.parseBinaryExpr
	p.binaryParseFunc[token.K_OR] = p.parseBinaryExpr
	p.binaryParseFunc[token.K_IS] = p.parseIsExpr
	p.binaryParseFunc[token.K_ISNULL] = p.parseIsExpr
	p.binaryParseFunc[token.K_NOTNULL] = p.parseIsExpr
	p.binaryParseFunc[token.K_NOT] = p.parseNotExpr
	p.binaryParseFunc[token.K_BETWEEN] = p.parseBetweenExpr
	p.binaryParseFunc[token.K_IN] = p.parseInExpr

	return p
}

func (p *parser) readToken() {
	if p.pos >= len(p.tokens) {
		p.currentToken = token.Token{
			Type: token.EOS,
		}
		return
	}
	p.currentToken = p.tokens[p.pos]
	p.pos++
	return
}

func (p *parser) getNextToken() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{
			Type: token.EOS,
		}
	}
	return p.tokens[p.pos]
}

func (p *parser) getCurrentTokenPrecedence() int {
	if p.pos > len(p.tokens) {
		return LOWEST
	}
	if p, exists := precedences[p.currentToken.Type]; exists {
		return p
	}
	return LOWEST
}

func (p *parser) getNextTokenPrecedence() int {
	if p.pos+1 > len(p.tokens) {
		return LOWEST
	}
	if p, exists := precedences[p.getNextToken().Type]; exists {
		return p
	}
	return LOWEST
}

func (p *parser) parse() (*ast.SQL, error) {
	var sql *ast.SQL
	for {
		if p.currentToken.Type == token.K_SELECT {
			ss, err := p.parseSELECTStatement()
			if err != nil {
				return nil, err
			}
			sql = &ast.SQL{
				SELECTStatement: ss,
			}
		} else {
			return nil, fmt.Errorf("Unexpected Token %s", p.currentToken.Literal)
		}
		if p.currentToken.Type == token.S_SEMICOLON {
			p.readToken()
			break
		}
		if p.currentToken.Type == token.EOS {
			break
		}
	}
	return sql, nil
}

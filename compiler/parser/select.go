package parser

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
)

func (p *parser) parseSELECTStatement() (*ast.SELECTStatement, error) {
	statement := &ast.SELECTStatement{}

	if p.currentToken.Type == token.K_SELECT {
		selectClause, err := p.parseSELECTClause()
		if err != nil {
			return statement, err
		}
		statement.SELECT = selectClause
	} else {
		return statement, fmt.Errorf("SELECT missing")
	}
	return statement, nil
}

func (p *parser) parseSELECTClause() (*ast.SELECTClause, error) {
	clause := &ast.SELECTClause{}
	p.readToken()

	cols, err := p.parseSelectColumns()
	if err != nil {
		return clause, err
	}
	clause.SelectColumns = cols
	return clause, nil
}

func (p *parser) parseSelectColumns() ([]ast.SelectColumn, error) {
	cols := []ast.SelectColumn{}
	loop := true
	for {
		switch p.currentToken.Type {
		case token.EOS, token.S_SEMICOLON:
			loop = false
			break
		case token.S_COMMA:
			p.readToken()
		default:
			expr, err := p.parseExpression(LOWEST)
			if err != nil {
				return cols, err
			}
			cols = append(cols, ast.SelectColumn{Expression: expr})
			p.readToken()
		}
		if !loop {
			break
		}
	}
	return cols, nil
}
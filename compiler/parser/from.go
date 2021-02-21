package parser

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
)

func (p *parser) parseFROMClause() (*ast.FROMClause, error) {
	clause := &ast.FROMClause{
		Joined: []ast.JoinedTable{},
	}
	p.readToken()

	tbl, err := p.parseTable()
	if err != nil {
		return clause, err
	}
	clause.Table = tbl

	for {
		jt := ast.JoinedTable{}
		if !(p.currentToken.Type == token.K_INNER || p.currentToken.Type == token.S_COMMA || p.currentToken.Type == token.K_LEFT || p.currentToken.Type == token.K_RIGHT || p.currentToken.Type == token.K_FULL) {
			break
		}
		if p.currentToken.Type == token.S_COMMA {
			jt.Type = ast.CROSS
		} else {
			if p.currentToken.Type == token.K_INNER {
				jt.Type = ast.INNER
				p.readToken()
			} else if p.currentToken.Type == token.K_FULL {
				jt.Type = ast.FULL
				if p.getNextToken().Type == token.K_OUTER {
					p.readToken()
				}
				p.readToken()
			} else if p.currentToken.Type == token.K_LEFT {
				jt.Type = ast.LEFT
				if p.getNextToken().Type == token.K_OUTER {
					p.readToken()
				}
				p.readToken()
			} else if p.currentToken.Type == token.K_RIGHT {
				jt.Type = ast.RIGHT
				if p.getNextToken().Type == token.K_OUTER {
					p.readToken()
				}
				p.readToken()
			}

			if p.currentToken.Type != token.K_JOIN {
				return clause, fmt.Errorf("Unknown JOIN format")
			}
			p.readToken()
		}

		tbl, err := p.parseTable()
		if err != nil {
			return clause, err
		}
		jt.Table = tbl

		if jt.Type != ast.CROSS {
			if p.currentToken.Type != token.K_ON {
				return clause, fmt.Errorf("Unknown JOIN format")
			}
			p.readToken()
			expr, err := p.parseExpression(LOWEST)
			if err != nil {
				return clause, err
			}
			jt.Condition = expr
		}
		clause.Joined = append(clause.Joined, jt)
		p.readToken()
	}

	return clause, nil
}

func (p *parser) parseTable() (*ast.Table, error) {
	tbl := &ast.Table{}
	literal := []string{}

	for {
		if p.currentToken.Type == token.IDENT {
			literal = append(literal, p.currentToken.Literal)
		} else {
			break
		}
		p.readToken()
		if p.currentToken.Type != token.S_PERIOD {
			break
		}
		p.readToken()
	}
	if len(literal) >= 4 {
		return tbl, fmt.Errorf("Unknown Table Expression")
	} else if len(literal) == 3 {
		tbl.Schema = literal[0]
		tbl.DB = literal[1]
		tbl.Table = literal[2]
	} else if len(literal) == 2 {
		tbl.DB = literal[0]
		tbl.Table = literal[1]
	} else {
		tbl.Table = literal[0]
	}

	if p.currentToken.Type == token.K_AS {
		p.readToken()
	}
	if p.currentToken.Type == token.IDENT {
		tbl.Alias = p.currentToken.Literal
		p.readToken()
	}

	return tbl, nil
}

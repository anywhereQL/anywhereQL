package parser

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/token"
)

func (p *parser) parseFROMClause() (*ast.FROMClause, error) {
	clause := &ast.FROMClause{}
	p.readToken()

	tbl, err := p.parseTable()
	if err != nil {
		return clause, err
	}
	clause.Table = tbl
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

	return tbl, nil
}

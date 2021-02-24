package runtime

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/compiler/lexer"
	"github.com/anywhereQL/anywhereQL/compiler/parser"
	"github.com/anywhereQL/anywhereQL/runtime/storage/virtual"
)

type ColumnInfo struct {
	Name     string
	RefCount int
}

type Runtime struct {
	columns  map[ast.Table]map[string]ColumnInfo
	rColumns map[string][]ast.Table
	alias    map[string]ast.Table
	tableID  map[ast.Table]string
}

func New() *Runtime {
	return &Runtime{
		columns:  make(map[ast.Table]map[string]ColumnInfo),
		rColumns: map[string][]ast.Table{},
		tableID:  map[ast.Table]string{},
		alias:    map[string]ast.Table{},
	}
}

func (r *Runtime) Start(sql string) ([][]value.Value, error) {
	tokens := lexer.Lex(sql)

	a, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return r.runSQL(a.SQL)
}

func (r *Runtime) runSQL(s *ast.SQL) ([][]value.Value, error) {
	if err := r.analyzeSQL(s); err != nil {
		return nil, err
	}
	r.readTables()
	var tID string
	var err error
	if s.SELECTStatement.FROM != nil {
		tID, err = r.joinTables(s.SELECTStatement.FROM)
		if err != nil {
			return nil, err
		}
	}
	if s.SELECTStatement.WHERE != nil {
		if tID != "" {
			tID, err = r.filter(tID, s.SELECTStatement.WHERE)
			if err != nil {
				return nil, err
			}
		}
	}
	return r.runSelectClause(s.SELECTStatement.SELECT, tID)
}

func (r *Runtime) runSelectClause(s *ast.SELECTClause, tID string) ([][]value.Value, error) {
	ret := [][]value.Value{}
	eng := virtual.VirtualStorage
	tl := eng.GetLine(tID)
	if tl == 0 {
		rl := []value.Value{}
		for _, col := range s.SelectColumns {
			vc := r.Translate(col.Expression)
			v, err := r.ExprRun(vc, nil)
			if err != nil {
				return ret, err
			}
			rl = append(rl, v)
		}
		ret = append(ret, rl)
	} else {
		for l := 0; l < tl; l++ {
			lineValue, err := eng.GetLineValue(tID, l)
			if err != nil {
				return ret, err
			}
			rl := []value.Value{}
			for _, col := range s.SelectColumns {
				vc := r.Translate(col.Expression)
				v, err := r.ExprRun(vc, lineValue)
				if err != nil {
					return ret, err
				}
				rl = append(rl, v)
			}
			ret = append(ret, rl)
		}
	}
	return ret, nil
}

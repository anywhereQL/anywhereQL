package virtual

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
)

type Engine struct {
	table map[string][]map[ast.Column]value.Value
}

var VirtualStorage = &Engine{
	table: make(map[string][]map[ast.Column]value.Value),
}

func (e *Engine) WriteTable(table string, values []map[ast.Column]value.Value) error {
	e.table[table] = values
	return nil
}

func (e *Engine) GetLine(tableID string) int {
	return len(e.table[tableID])
}

func (e *Engine) GetValue(tableID string, column ast.Column, line int) (value.Value, error) {
	val := e.table[tableID][line][column]

	return val, nil
}

func (e *Engine) GetLineValue(tableID string, line int) (map[ast.Column]value.Value, error) {
	vals := e.table[tableID][line]
	return vals, nil
}

func (e *Engine) GetAllColumnInfo(tableID string) ([]ast.Column, error) {
	ret := []ast.Column{}
	for key := range e.table[tableID][0] {
		ret = append(ret, key)
	}
	return ret, nil
}

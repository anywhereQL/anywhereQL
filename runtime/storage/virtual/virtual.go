package virtual

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
)

type Engine struct {
	table          map[string][]map[string]value.Value
	tableList      map[string]map[string]map[string]string
	tableAliasList map[string]string
}

var VirtualStorage = &Engine{
	table: make(map[string][]map[string]value.Value),
}

func (e *Engine) WriteTable(table string, values []map[string]value.Value) error {
	tbl := storage.Table{
		Name: table,
		Path: "",
	}
	cols := []storage.Column{}

	for k, v := range values[0] {
		cols = append(cols, storage.Column{
			Name: k,
			Type: v.Type,
		})
	}
	tbl.Columns = cols

	e.table[table] = values
	return nil
}

func (e *Engine) GetLine(tableID string) int {
	return len(e.table[tableID])
}

func (e *Engine) GetValue(tableID, column string, line int) (value.Value, error) {
	val := e.table[tableID][line][column]

	return val, nil
}
func (e *Engine) SetTableList(schema, db, table, alias, vTable string) error {

	if _, exists := e.tableList[schema]; !exists {
		e.tableList[schema] = make(map[string]map[string]string)
	}
	if _, exists := e.tableList[schema][db]; !exists {
		e.tableList[schema][db] = make(map[string]string)
	}
	if v, exists := e.tableList[schema][db][table]; exists {
		if v != vTable {
			return fmt.Errorf("Already exists")
		}
		return nil
	}
	e.tableList[schema][db][table] = vTable
	return nil
}

func (e *Engine) GetTable(schema, db, table string) (string, bool) {
	if v, exists := e.tableList[schema][db][table]; exists {
		return v, true
	}
	return "", false
}

func (e *Engine) SetTableAlias(alias, vTable string) error {
	if _, exists := e.tableAliasList[alias]; exists {
		return fmt.Errorf("Already exists")
	}
	e.tableAliasList[alias] = vTable
	return nil
}

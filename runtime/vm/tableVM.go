package vm

import (
	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
	"github.com/anywhereQL/anywhereQL/runtime/storage/virtual"
)

func TableRun(tbl *ast.FROMClause, rTable map[string]ast.Table) (string, error) {
	eng := virtual.VirtualStorage
	sInst := storage.GetInstance()

	lastTable := ""
	for tID, tbl := range rTable {
		srcEng := sInst.GetEngine(tbl.Schema)
		values, _ := (*srcEng).GetTableValues(tbl.DB, tbl.Table)
		eng.WriteTable(tID, values)
		lastTable = tID
	}

	return lastTable, nil
}

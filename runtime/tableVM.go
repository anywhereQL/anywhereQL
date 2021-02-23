package runtime

import (
	"github.com/google/uuid"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
	"github.com/anywhereQL/anywhereQL/runtime/storage/virtual"
)

func (r *Runtime) readTables() {
	eng := virtual.VirtualStorage

	for tbl, cols := range r.columns {
		tblID := uuid.New()
		tableValues, err := r.readTable(tbl)
		if err != nil {
			return
		}
		vals := []map[ast.Column]value.Value{}
		for _, v := range tableValues {
			line := make(map[ast.Column]value.Value)
			for k, col := range cols {
				if col.RefCount == 0 {
					continue
				}
				kc := ast.Column{
					Column: k,
					Table:  tbl,
				}
				line[kc] = v[k]
			}
			vals = append(vals, line)
		}
		eng.WriteTable(tblID.String(), vals)
		r.tableID[tbl] = tblID.String()
	}
}

func (r *Runtime) readTable(tbl ast.Table) ([]map[string]value.Value, error) {
	sInst := storage.GetInstance()
	srcEng := sInst.GetEngine(tbl.Schema)
	return (*srcEng).GetTableValues(tbl.DB, tbl.Table)
}

func (r *Runtime) joinTables(fc *ast.FROMClause) (string, error) {
	leftTableID := r.tableID[ast.Table{Schema: fc.Table.Schema, DB: fc.Table.DB, Table: fc.Table.Table}]
	for _, jt := range fc.Joined {
		rightTableID := r.tableID[ast.Table{Schema: jt.Table.Schema, DB: jt.Table.DB, Table: jt.Table.Table}]
		lID, err := r.join(leftTableID, rightTableID, jt.Type, jt.Condition)
		if err != nil {
			return "", err
		}
		leftTableID = lID
	}
	return leftTableID, nil
}

func (r *Runtime) join(left, right string, tp ast.JoinType, cond *ast.Expression) (string, error) {
	eng := virtual.VirtualStorage
	ll := eng.GetLine(left)
	rl := eng.GetLine(right)
	tblValues := []map[ast.Column]value.Value{}
	if tp == ast.INNER {
		for ln := 0; ln < ll; ln++ {
			for rn := 0; rn < rl; rn++ {
				lv, _ := eng.GetLineValue(left, ln)
				rv, _ := eng.GetLineValue(right, rn)

				vals := make(map[ast.Column]value.Value)
				for k, v := range lv {
					vals[k] = v
				}
				for k, v := range rv {
					vals[k] = v
				}
				vc := r.Translate(cond)
				ret, err := r.ExprRun(vc, vals)
				if err != nil {
					return "", err
				}
				if ret.Type != value.BOOL {
					continue
				}
				if ret.Bool.True == true {
					tblValues = append(tblValues, vals)
				}
			}
		}
	} else if tp == ast.LEFT {
		for ln := 0; ln < ll; ln++ {
			isFound := false
			lv, _ := eng.GetLineValue(left, ln)
			for rn := 0; rn < rl; rn++ {
				rv, _ := eng.GetLineValue(right, rn)

				vals := make(map[ast.Column]value.Value)
				for k, v := range lv {
					vals[k] = v
				}
				for k, v := range rv {
					vals[k] = v
				}
				vc := r.Translate(cond)
				ret, err := r.ExprRun(vc, vals)
				if err != nil {
					return "", err
				}
				if ret.Type != value.BOOL {
					continue
				}
				if ret.Bool.True == true {
					tblValues = append(tblValues, vals)
					isFound = true
				}
			}
			if !isFound {
				line := make(map[ast.Column]value.Value)
				for k, v := range lv {
					line[k] = v
				}
				if rl != 0 {
					rv, _ := eng.GetLineValue(right, 0)
					for k := range rv {
						line[k] = value.Value{
							Type: value.NULL,
						}
					}
				}
				tblValues = append(tblValues, line)
			}
		}
	} else if tp == ast.RIGHT {
		return r.join(right, left, ast.LEFT, cond)
	}
	tID := uuid.New().String()
	eng.WriteTable(tID, tblValues)
	return tID, nil
}

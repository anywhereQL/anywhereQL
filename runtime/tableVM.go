package runtime

import (
	"fmt"
	"reflect"

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
	} else if tp == ast.FULL {
		le, err := r.join(left, right, ast.LEFT, cond)
		if err != nil {
			return "", err
		}
		ri, err := r.join(left, right, ast.RIGHT, cond)
		if err != nil {
			return "", err
		}
		return r.union(le, ri)
	}
	tID := uuid.New().String()
	eng.WriteTable(tID, tblValues)
	return tID, nil
}

func (r *Runtime) union(tbl1, tbl2 string) (string, error) {
	fmt.Printf("UINON %s %s\n", tbl1, tbl2)
	eng := virtual.VirtualStorage
	l1 := eng.GetLine(tbl1)
	l2 := eng.GetLine(tbl2)
	tblValues := []map[ast.Column]value.Value{}

	for ln1 := 0; ln1 < l1; ln1++ {
		v1, _ := eng.GetLineValue(tbl1, ln1)
		tblValues = append(tblValues, v1)
	}
	for ln2 := 0; ln2 < l2; ln2++ {
		v2, _ := eng.GetLineValue(tbl2, ln2)
		isSame := false
		for ln1 := 0; ln1 < l1; ln1++ {
			v1, _ := eng.GetLineValue(tbl1, ln1)
			if len(v1) != len(v2) {
				return "", fmt.Errorf("Column length mismatch")
			}
			isColumnSame := true
			for k, v := range v1 {
				if !reflect.DeepEqual(v, v2[k]) {
					isColumnSame = false
					break
				}
			}
			if isColumnSame == true {
				isSame = true
			}
		}
		if isSame == false {
			tblValues = append(tblValues, v2)
		}
	}

	tID := uuid.New().String()
	eng.WriteTable(tID, tblValues)
	return tID, nil
}

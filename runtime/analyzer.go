package runtime

import (
	"fmt"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/config"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
)

func (r *Runtime) analyzeSQL(s *ast.SQL) error {
	if s.SELECTStatement.FROM != nil {
		if err := r.analyzeTable(s.SELECTStatement.FROM); err != nil {
			return err
		}
	}
	for _, col := range s.SELECTStatement.SELECT.SelectColumns {
		r.analyzeExpr(col.Expression)
	}
	if s.SELECTStatement.WHERE != nil {
		r.analyzeExpr(s.SELECTStatement.WHERE)
	}
	return nil
}

func (r *Runtime) analyzeTable(f *ast.FROMClause) error {
	if f.Table == nil {
		return nil
	}
	r.getTableInfo(f.Table)

	for _, jt := range f.Joined {
		if _, exists := r.columns[*jt.Table]; !exists {
			r.getTableInfo(jt.Table)
		}
		if jt.Type != ast.CROSS {
			if err := r.analyzeExpr(jt.Condition); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runtime) getTableInfo(tbl *ast.Table) {
	if tbl.Schema == "" {
		tbl.Schema = config.DBConfig.DefaultSchema
	}
	if tbl.DB == "" {
		tbl.DB = config.DBConfig.DefaultDB
	}

	t := ast.Table{
		Table:  tbl.Table,
		DB:     tbl.DB,
		Schema: tbl.Schema,
	}
	sInst := storage.GetInstance()
	srcEng := sInst.GetEngine(tbl.Schema)
	cMap := make(map[string]ColumnInfo)
	for _, k := range (*srcEng).GetColumns(tbl.DB, tbl.Table) {
		cMap[k] = ColumnInfo{Name: k, RefCount: 0}
		r.rColumns[k] = append(r.rColumns[k], t)
	}
	r.columns[t] = cMap
	if tbl.Alias != "" {
		r.alias[tbl.Alias] = t
	}
	return
}

func (r *Runtime) analyzeExpr(expr *ast.Expression) error {
	if expr.Column != nil {
		cIDs := r.rColumns[expr.Column.Column]
		if expr.Column.Table.Schema == "" && expr.Column.Table.DB == "" && expr.Column.Table.Table == "" {
			if len(cIDs) != 1 {
				return fmt.Errorf("Ambious Column %s", expr.Column.Column)
			}
			expr.Column.Table = cIDs[0]
			ci := r.columns[expr.Column.Table][expr.Column.Column]
			ci.RefCount += 1
			r.columns[expr.Column.Table][expr.Column.Column] = ci
			return nil
		} else if expr.Column.Table.Schema == "" && expr.Column.Table.DB == "" {
			if t, exists := r.alias[expr.Column.Table.Table]; exists {
				expr.Column.Table = t
				ci := r.columns[expr.Column.Table][expr.Column.Column]
				ci.RefCount += 1
				r.columns[expr.Column.Table][expr.Column.Column] = ci
				return nil
			}
			st := []ast.Table{}
			for _, tbl := range cIDs {
				if tbl.Table == expr.Column.Table.Table {
					st = append(st, tbl)
				}
			}
			if len(st) != 1 {
				return fmt.Errorf("Ambious Column %s", expr.Column.Column)
			}
			expr.Column.Table = st[0]
			ci := r.columns[expr.Column.Table][expr.Column.Column]
			ci.RefCount += 1
			r.columns[expr.Column.Table][expr.Column.Column] = ci
			return nil
		} else if expr.Column.Table.Schema == "" {
			st := []ast.Table{}
			for _, tbl := range cIDs {
				if tbl.Table == expr.Column.Table.Table && tbl.DB == expr.Column.Table.DB {
					st = append(st, tbl)
				}
			}
			if len(st) != 1 {
				return fmt.Errorf("Ambious Column %s", expr.Column.Column)
			}
			expr.Column.Table = st[0]
			ci := r.columns[expr.Column.Table][expr.Column.Column]
			ci.RefCount += 1
			r.columns[expr.Column.Table][expr.Column.Column] = ci
			return nil
		} else {
			st := []ast.Table{}
			for _, tbl := range cIDs {
				if tbl.Table == expr.Column.Table.Table && tbl.DB == expr.Column.Table.DB && tbl.Schema == expr.Column.Table.Schema {
					st = append(st, tbl)
				}
			}
			if len(st) != 1 {
				return fmt.Errorf("Ambious Column %s", expr.Column.Column)
			}
			expr.Column.Table = st[0]
			ci := r.columns[expr.Column.Table][expr.Column.Column]
			ci.RefCount += 1
			r.columns[expr.Column.Table][expr.Column.Column] = ci
			return nil
		}
	} else if expr.BinaryOperation != nil {
		if err := r.analyzeExpr(expr.BinaryOperation.Left); err != nil {
			return err
		}
		if err := r.analyzeExpr(expr.BinaryOperation.Right); err != nil {
			return err
		}
	} else if expr.FunctionCall != nil {
		for _, col := range expr.FunctionCall.Args {
			if err := r.analyzeExpr(&col); err != nil {
				return err
			}
		}
	} else if expr.UnaryOperation != nil {
		if err := r.analyzeExpr(expr.UnaryOperation.Expr); err != nil {
			return err
		}
	} else if expr.Cast != nil {
		if err := r.analyzeExpr(expr.Cast.Expr); err != nil {
			return err
		}
	} else if expr.Case != nil {
		if expr.Case.Value != nil {
			if err := r.analyzeExpr(expr.Case.Value); err != nil {
				return err
			}
		}
		for _, ca := range expr.Case.CaseValues {
			if err := r.analyzeExpr(ca.Condition); err != nil {
				return err
			}
			if err := r.analyzeExpr(ca.Result); err != nil {
				return err
			}
		}
		if expr.Case.ElseValue != nil {
			if err := r.analyzeExpr(expr.Case.ElseValue); err != nil {
				return err
			}
		}
	} else if expr.Between != nil {
		if err := r.analyzeExpr(expr.Between.Src); err != nil {
			return err
		}
		if err := r.analyzeExpr(expr.Between.Begin); err != nil {
			return err
		}
		if err := r.analyzeExpr(expr.Between.End); err != nil {
			return err
		}
	}
	return nil
}

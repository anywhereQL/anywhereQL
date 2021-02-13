package runtime

import (
	"os"

	"github.com/google/uuid"

	"github.com/anywhereQL/anywhereQL/common/ast"
	"github.com/anywhereQL/anywhereQL/common/config"
	"github.com/anywhereQL/anywhereQL/common/debug"
	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/compiler/lexer"
	"github.com/anywhereQL/anywhereQL/compiler/parser"
	"github.com/anywhereQL/anywhereQL/compiler/planner"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
	"github.com/anywhereQL/anywhereQL/runtime/storage/virtual"
	"github.com/anywhereQL/anywhereQL/runtime/vm"
)

type Runtime struct {
	tables    map[ast.Table]string
	columns   map[string][]string
	revTables map[string]ast.Table
}

func New() *Runtime {
	return &Runtime{
		tables:    make(map[ast.Table]string),
		columns:   make(map[string][]string),
		revTables: make(map[string]ast.Table),
	}
}

func (r *Runtime) Start(sql string) ([][]value.Value, error) {
	ret := [][]value.Value{}
	tokens := lexer.Lex(sql)

	a, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	if err := r.pickupTable(a); err != nil {
		return [][]value.Value{}, err
	}
	if err := r.rewriteColumn(a); err != nil {
		return [][]value.Value{}, err
	}

	debug.PrintAST(os.Stdout, a)

	for _, s := range a.SQL {
		lastID, err := vm.TableRun(s.SELECTStatement.FROM, r.revTables)
		if err != nil {
			return nil, err
		}

		ln := 1
		if lastID != "" {
			vEng := virtual.VirtualStorage
			ln = vEng.GetLine(lastID)
		}

		for l := 0; l < ln; l++ {
			lineVal := []value.Value{}
			for _, col := range s.SELECTStatement.SELECT.SelectColumns {
				vc := planner.Translate(col.Expression)
				debug.PrintExprVC(os.Stdout, vc)
				rs, err := vm.ExprRun(vc, l)
				if err != nil {
					return nil, err
				}
				lineVal = append(lineVal, rs)
			}
			ret = append(ret, lineVal)
		}
	}
	return ret, nil
}

func (r *Runtime) pickupTable(a *ast.AST) error {
	tables := make(map[ast.Table]string)
	columns := make(map[string][]string)
	revTables := make(map[string]ast.Table)

	sInst := storage.GetInstance()

	for _, sql := range a.SQL {
		if sql.SELECTStatement.FROM == nil {
			continue
		}
		tbl := sql.SELECTStatement.FROM.Table
		t := ast.Table{}
		t.Schema = tbl.Schema
		if tbl.Schema == "" {
			t.Schema = config.DBConfig.DefaultSchema
		}
		t.DB = tbl.DB
		if tbl.DB == "" {
			t.DB = config.DBConfig.DefaultDB
		}
		t.Table = tbl.Table

		if _, exists := tables[t]; !exists {
			uuidv4 := uuid.New()
			tables[t] = uuidv4.String()
			eng := sInst.GetEngine(t.Schema)
			cols := (*eng).GetColumns(t.DB, t.Table)
			for _, c := range cols {
				columns[c] = append(columns[c], uuidv4.String())
			}
			revTables[uuidv4.String()] = t
		}
	}
	r.tables = tables
	r.columns = columns
	r.revTables = revTables
	return nil
}

func (r *Runtime) rewriteColumn(a *ast.AST) error {
	for _, sql := range a.SQL {
		for _, col := range sql.SELECTStatement.SELECT.SelectColumns {
			if err := r.scanExpr(col.Expression); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Runtime) scanExpr(expr *ast.Expression) error {
	if expr.Column != nil {
		if err := r.scanColumn(expr.Column); err != nil {
			return err
		}
	} else if expr.BinaryOperation != nil {
		if err := r.scanExpr(expr.BinaryOperation.Left); err != nil {
			return err
		}
		if err := r.scanExpr(expr.BinaryOperation.Right); err != nil {
			return err
		}
	} else if expr.FunctionCall != nil {
		for _, arg := range expr.FunctionCall.Args {
			if err := r.scanExpr(&arg); err != nil {
				return err
			}
		}
	} else if expr.UnaryOperation != nil {
		if err := r.scanExpr(expr.UnaryOperation.Expr); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) scanColumn(col *ast.Column) error {
	tIDs := r.columns[col.Column]
	if len(tIDs) == 1 {
		col.Table.ID = tIDs[0]
	} else {
		sch := col.Table.Schema
		db := col.Table.DB
		tbl := col.Table.Table
		if sch == "" {
			sch = config.DBConfig.DefaultSchema
		}
		if db == "" {
			db = config.DBConfig.DefaultDB
		}

		tID := r.tables[ast.Table{Table: tbl, Schema: sch, DB: db}]
		col.Table.ID = tID
	}
	return nil
}

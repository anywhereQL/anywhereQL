package storage

import (
	"github.com/anywhereQL/anywhereQL/common/value"
)

type Engine interface {
	GetEngineName() string
	GetDatabase() map[string]DB
	GetTableValues(db, table string) ([]map[string]value.Value, error)
	GetValue(db, table, column string, line int) (value.Value, error)
	GetColumns(db, table string) []string
	WriteTable(db, table string, values []map[string]value.Value) error
}

type DB struct {
	Name   string
	Path   string
	Tables map[string]Table
}

type Table struct {
	Name    string
	Path    string
	Columns []Column
}

type Column struct {
	Name string
	Type value.Type
}

type Storage struct {
	Engine map[string]*Engine
}

var instance = &Storage{
	Engine: make(map[string]*Engine),
}

func GetInstance() *Storage {
	return instance
}

func (s *Storage) Add(schema string, engine Engine) {
	s.Engine[schema] = &engine
}

func (s *Storage) GetEngine(schema string) *Engine {
	return s.Engine[schema]
}

type DBInfo struct {
	Schema string
	DB     string
}

func (s *Storage) GetDBInfoFromDB(db string) []DBInfo {
	ret := []DBInfo{}
	for schema, eng := range s.Engine {
		for d := range (*eng).GetDatabase() {
			if d == db {
				ret = append(ret, DBInfo{Schema: schema, DB: d})
			}
		}
	}
	return ret
}

type TableInfo struct {
	Schema string
	DB     string
	Table  string
}

func (s *Storage) GetDBInfoFromTable(table string) []TableInfo {
	ret := []TableInfo{}
	for schema, eng := range s.Engine {
		for db, di := range (*eng).GetDatabase() {
			for tbl := range di.Tables {
				if tbl == table {
					ret = append(ret, TableInfo{Schema: schema, DB: db, Table: tbl})
				}
			}
		}
	}
	return ret
}

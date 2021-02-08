package aq

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/anywhereQL/anywhereQL/common/value"
	"github.com/anywhereQL/anywhereQL/runtime/storage"
)

type Engine struct {
	DB map[string]storage.DB
}

func (e Engine) GetEngineName() string {
	return "AQDB Engine"
}

func (e Engine) GetDatabase() map[string]storage.DB {
	return e.DB
}

type dbMetaInfo struct {
	Version float32     `json:"version"`
	Name    string      `json:"name"`
	Tables  []tableInfo `json:"tables"`
}

type tableInfo struct {
	Path string `json:"path"`
}

type tableMetaInfo struct {
	Version float32      `json:"version"`
	Name    string       `json:"name"`
	Path    string       `json:"file"`
	Columns []columnInfo `json:"columns"`
}

type columnInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Definition string `json:"definition"`
}

func Start(path ...string) (*Engine, error) {
	e := &Engine{
		DB: make(map[string]storage.DB),
	}
	for _, p := range path {
		fp, err := filepath.Abs(p)
		raw, err := ioutil.ReadFile(filepath.Join(fp, ".db.metainfo.json"))
		if err != nil {
			return nil, err
		}
		var i dbMetaInfo
		if err := json.Unmarshal(raw, &i); err != nil {
			return nil, err
		}

		e.DB[i.Name] = storage.DB{
			Path:   fp,
			Name:   i.Name,
			Tables: make(map[string]storage.Table),
		}

		for _, t := range i.Tables {
			p := t.Path
			if !filepath.IsAbs(t.Path) {
				p, err = filepath.Abs(filepath.Join(fp, t.Path))
				if err != nil {
					continue
				}
			}
			if err := e.readTableInfo(i.Name, p); err != nil {
				return nil, err
			}
		}
	}
	return e, nil
}

func (e *Engine) readTableInfo(db, path string) error {
	raw, err := ioutil.ReadFile(filepath.Join(path, ".table.metainfo.json"))
	if err != nil {
		return err
	}
	var i tableMetaInfo
	if err := json.Unmarshal(raw, &i); err != nil {
		return err
	}

	fullPath, err := filepath.Abs(filepath.Join(path, i.Path))
	if err != nil {
		return err
	}
	tbl := storage.Table{
		Name:    i.Name,
		Path:    fullPath,
		Columns: []storage.Column{},
	}

	for _, col := range i.Columns {
		c := storage.Column{
			Name: col.Name,
		}
		switch strings.ToUpper(col.Type) {
		case "INT":
			c.Type = value.INTEGER
		case "FLOAT":
			c.Type = value.FLOAT
		case "STRING":
			c.Type = value.STRING
		default:
			c.Type = value.UNKNOWN
		}
		e.parseDefinition(col.Definition)
		tbl.Columns = append(tbl.Columns, c)
	}
	e.DB[db].Tables[i.Name] = tbl

	return nil
}

func (e *Engine) parseDefinition(def string) {
	defs := strings.Split(strings.ToUpper(def), " ")
	for _, d := range defs {
		switch d {
		default:
			continue
		}
	}
}

func (e *Engine) GetColumnValues(db, table string) ([]map[string]value.Value, error) {
	ret := []map[string]value.Value{}

	fp, err := os.Open(e.DB[db].Tables[table].Path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = '\t'
	reader.Comment = '#'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		cols, err := e.parseColumn(db, table, record)
		if err != nil {
			return nil, err
		}

		lv := make(map[string]value.Value)
		for n, c := range cols {
			lv[e.DB[db].Tables[table].Columns[n].Name] = c
		}

		ret = append(ret, lv)
	}

	return ret, nil
}

func (e *Engine) parseColumn(db, table string, cols []string) ([]value.Value, error) {
	ret := []value.Value{}
	for n, col := range cols {
		val := value.Value{}
		switch e.DB[db].Tables[table].Columns[n].Type {
		case value.INTEGER:
			v, err := strconv.ParseInt(col, 10, 64)
			if err != nil {
				return nil, err
			}
			val.Int = v
			val.Type = value.INTEGER
		case value.FLOAT:
			v, err := strconv.ParseFloat(col, 64)
			if err != nil {
				return nil, err
			}
			val.Float = v
			val.Type = value.FLOAT
		case value.STRING:
			val.String = col
			val.Type = value.STRING
		default:
			return nil, fmt.Errorf("Not Impli")
		}
		ret = append(ret, val)
	}
	return ret, nil
}

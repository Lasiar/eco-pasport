package model

import (
	"eco-passport-back/base"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EpTree eptree
type EpTree struct {
	TreeItem []*nodeEpTree `xml:"TreeItem"`
}

type nodeEpTree struct {
	Name      string        `xml:"name,attr"`
	TableID   string        `xml:"table_id,attr" json:",omitempty"`
	TableName string        `xml:"table_name,attr"  json:",omitempty"`
	TreeItem  []*nodeEpTree `xml:"TreeItem"  json:",omitempty"`
}

// GetTablesInfo получение данных с базы
func (d *Database) GetTablesInfo() (map[int]TableInfo, error) {
	if d.err != nil {
		return nil, d.err
	}

	rows, err := d.db.Query(sqlGetTables)
	if err != nil {
		return nil, fmt.Errorf("[db] query %v", err)
	}

	t := make(map[int]TableInfo)

	for rows.Next() {
		row := struct {
			id      int
			dbName  string
			visName string
		}{}
		if err := rows.Scan(&row.id, &row.dbName, &row.visName); err != nil {
			return nil, fmt.Errorf("[db] scan %v", err)
		}
		t[row.id] = TableInfo{row.dbName, row.visName}
	}
	return t, nil
}

// GetTree get table tree
func GetTree() (EpTree, error) {
	res := struct {
		TablesMeta map[int]TableInfo
		EpTree
	}{}
	res.TablesMeta = make(map[int]TableInfo)
	res.EpTree.load("./Tree.xml")
	var err error
	res.TablesMeta, err = GetDatabase().GetTablesInfo()
	if err != nil {
		return EpTree{}, fmt.Errorf("[db] get table info: %v", err)
	}
	if err := changeName(res.TreeItem, res.TablesMeta); err != nil {
		return EpTree{}, fmt.Errorf("change name %v", err)
	}
	return res.EpTree, nil
}

func (e *EpTree) load(path string) {
	file, err := os.Open(path)
	if err != nil {
		base.GetConfig().Err.Fatalf("Can`t read tree file from %v err %v", path, err)
	}
	d := xml.NewDecoder(file)
	if err := d.Decode(&e); err != nil {
		base.GetConfig().Err.Fatalf("Can`t read tree file from %v err %v", path, err)
	}
}

func changeName(t []*nodeEpTree, table map[int]TableInfo) error {
	var sumError []string
	var hasError bool
	for _, node := range t {
		if node.Name == "" {
			id, err := strconv.Atoi(node.TableID)
			if err != nil {
				hasError = true
				sumError = append(sumError, fmt.Sprint(err))
			}
			if table, ok := table[id]; ok {
				node.Name = table.VisName
			}
		}
		sumError = append(sumError, fmt.Sprint(changeName(node.TreeItem, table)))
	}
	if hasError {
		return errors.New(strings.Join(sumError, "."))
	}
	return nil
}

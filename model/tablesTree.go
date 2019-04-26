package model

import (
	"EcoPasport/base"
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

	fmt.Println(res.TreeItem)

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

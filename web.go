package main

import (
	"fmt"
	"log"
	"strconv"
)

func WebGetTree() {
	t := GetEpTree()

	for _, node := range t.TreeItem {
		if len(node.TableID) != 4 {
			continue
		}
		id, err := strconv.Atoi(node.TableID)
		if err != nil {
			log.Println(err)
		}

		info := new(TablesInfo)

		if err := info.GetTables(); err != nil {
			log.Fatal(err)
		}



		for _, region := range *info {

			if id == region.ID {
				node.Name = region.VisName
			}

		}

		if len(t.TreeItem) != 0 {
			ChangeName(node, info)
		}
	}

}

func ChangeName(t *nodeEpTree, table *TablesInfo) {
	for _, node := range t.TreeItem {
		if len(node.TableID) != 4 {
			continue
		}
		id, err := strconv.Atoi(node.TableID)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, region := range *table {
			if id == region.ID {
				node.Name = region.VisName
				fmt.Println(id, region.VisName)
			}
		}
		if len(t.TreeItem) != 0 {
			ChangeName(node, table)
		}
	}

}

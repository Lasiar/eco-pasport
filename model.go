package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

const (
	SQLGetTables  = "SELECT ID_EP_Table, DBName, VisName FROM krasecology.eco.EP_Table"
	SQLGetRegions = "SELECT id, num_region, name, IsTown from krasecology.eco.EP_Region"
)

type database struct {
	*sql.DB
}

func (d *database) connect() (err error) {
	d.DB, err = sql.Open("mssql", GetConfig().ConnStr)
	if err != nil {
		return err
	}
	return nil
}

type region struct {
	ID        int
	NumRegion int
	Name      string
	IsTown    bool
}

type Regions []region

func (r *Regions) GetRegions() error {
	db := new(database)

	if err := db.connect(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(SQLGetRegions)
	if err != nil {
		return err
	}
	region := new(region)

	for rows.Next() {
		if err := rows.Scan(&region.ID, &region.NumRegion, &region.Name, &region.IsTown); err != nil {
			return err
		}
		*r = append(*r, *region)
	}
	return nil
}

type tableInfo struct {
	DBTable string
	VisName string
}

type TablesInfo map[int]tableInfo

func (t *TablesInfo) GetTables() error {
	db := new(database)

	if err := db.connect(); err != nil {
		return err
	}

	rows, err := db.Query(SQLGetTables)
	if err != nil {
		return err
	}

	*t = make(map[int]tableInfo)

	for rows.Next() {

		var id int
		var dbName, visName string

		if err := rows.Scan(&id, &dbName, &visName); err != nil {
			return err
		}
		(*t)[id] = tableInfo{dbName, visName}
	}

	return nil
}

package main

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

const (
	SQLGetTables  = `SELECT Table_ID, DB_Name, VisName FROM krasecology.eco_2018.Table_0_1_Tables;`
	SQLGetRegions = "SELECT id, num_region, name, CAST(IIF ( is_town = 1, 1, 0 ) AS BIT) AS is_town from krasecology.eco_2018.Table_0_0_Regions"
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

type Region struct {
	ID        int
	NumRegion int
	Name      string
	IsTown    bool
}

func GetRegions() *[]Region {
	db := new(database)

	if err := db.connect(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(SQLGetRegions)
	if err != nil {
		log.Fatal(err)
	}
	region := new(Region)
	regions := new([]Region)
	for rows.Next() {
		if err := rows.Scan(&region.ID, &region.NumRegion, &region.Name, &region.IsTown); err != nil {
			log.Fatal(err)
		}
		*regions = append(*regions, *region)
	}
	return regions
}

type tableInfo struct {
	ID      int
	Table   string
	VisName string
}

type TablesInfo []tableInfo

func (t *TablesInfo) GetTables() error {
	db := new(database)

	if err := db.connect(); err != nil {
		return err
	}

	rows, err := db.Query(SQLGetTables)
	if err != nil {
		return err
	}

	for rows.Next() {

		var id int
		var dbName, visName string

		if err := rows.Scan(&id, &dbName, &visName, ); err != nil {
			return err
		}

		*t = append(*t, tableInfo{id, dbName, visName})
	}
	return nil
}

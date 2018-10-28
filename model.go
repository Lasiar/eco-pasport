package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	SQLGetTables  = "SELECT ID_EP_Table, DBName, VisName FROM krasecology.eco.EP_Table"
	SQLGetRegions = "SELECT id, num_region, name, IsTown from krasecology.eco.EP_Region"
	SQLGetHeaders = "SELECT column_name, caption from krasecology.eco.Table_Column where ID_EP_Table=?"
	SQLGetSQL     = `
USE krasecology;

declare @SQL varchar(max) EXECUTE eco.sp_get_table ?,
?,
?,
@SQL output
EXECUTE (@sql)
`
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

func (r *Regions) FetchRegions() error {
	db := new(database)

	if err := db.connect(); err != nil {
		return fmt.Errorf("[DB] connect %v", err)
	}

	rows, err := db.Query(SQLGetRegions)
	if err != nil {
		return fmt.Errorf("[DB] query %v", err)
	}
	region := new(region)

	for rows.Next() {
		if err := rows.Scan(&region.ID, &region.NumRegion, &region.Name, &region.IsTown); err != nil {
			return fmt.Errorf("[DB] scan %v", err)
		}
		*r = append(*r, *region)
	}
	return nil
}

type RequestTableInfo struct {
	User     string `json:"user"`
	RegionID int    `json:"region_id"`
	TableID  int    `json:"table_id"`
}

type tableInfo struct {
	DBTable string
	VisName string
}

type TablesInfo map[int]tableInfo

func (t *TablesInfo) FetchTables() error {
	db := new(database)

	if err := db.connect(); err != nil {
		return fmt.Errorf("[DB] connect %v", err)
	}

	rows, err := db.Query(SQLGetTables)
	if err != nil {
		return fmt.Errorf("[DB] query %v", err)
	}

	*t = make(map[int]tableInfo)

	for rows.Next() {

		var id int
		var dbName, visName string

		if err := rows.Scan(&id, &dbName, &visName); err != nil {
			return fmt.Errorf("[DB] scan %v", err)
		}
		(*t)[id] = tableInfo{dbName, visName}
	}

	return nil
}

type Table struct {
	Header []string
	Value  [][]string
}

func (t *Table) FetchTableBySQL(info *RequestTableInfo) error {
	db := new(database)

	if err := db.connect(); err != nil {
		return fmt.Errorf("[DB] connect %v", err)
	}

	rows, err := db.Query(SQLGetSQL, info.User, info.TableID, info.RegionID)
	if err != nil {
		return fmt.Errorf("[DB] query %v", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[DB] column %v", err)
	}

	headers := new(headers)

	if err := headers.fetchHeaders(info.TableID); err != nil {
		return fmt.Errorf("[DB] fetch headers %v", err)
	}

	for _, column := range columns {
		t.Header = append(t.Header, (*headers)[column])
	}

	rawResult := make([][]byte, len(t.Header))
	result := make([]string, len(t.Header))

	dest := make([]interface{}, len(t.Header))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {

		err = rows.Scan(dest...)
		if err != nil {
			return fmt.Errorf("[DB] rows scan %v", err)
		}

		for j, raw := range rawResult {
			if raw == nil {
				result[j] = ""
			} else {
				result[j] = string(raw)
			}
		}
		t.Value = append(t.Value, result)
	}

	return nil
}

type headers map[string]string

func (h *headers) fetchHeaders(tableID int) error {
	db := new(database)
	if err := db.connect(); err != nil {
		return err
	}

	rows, err := db.Query(SQLGetHeaders, tableID)
	if err != nil {
		return err
	}

	var dbName, visName string

	*h = make(map[string]string)

	for rows.Next() {
		if err := rows.Scan(&dbName, &visName); err != nil {
			return err
		}
		(*h)[dbName] = visName
	}
	return nil
}

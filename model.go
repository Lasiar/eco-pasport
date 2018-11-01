package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	/*	SQLGetTables           = "SELECT ID_EP_Table, DBName, VisName FROM krasecology.eco.EP_Table"
		SQLGetRegions          = "SELECT id, num_region, name, IsTown from krasecology.eco.EP_Region"
		SQLGetHeaders          = "SELECT ID_EP_Table, column_name, caption from krasecology.eco.Table_Column"
		SQLGetEmptyText string = "SELECT Empty_text FROM krasecology.eco.l_EP_Tabe_EP_Region WHERE ID_EP_Table=? and ID_EP_Region=? "
		SQLGetSQL              = `
	USE krasecology;

	declare @SQL varchar(max) EXECUTE eco.sp_get_table 'babay@krasecology.ru',
	?,
	?,
	@SQL output
	EXECUTE (@sql)
	`
	)

	*/
	SQLGetTables    string = "SELECT Table_ID, DB_Name, VisName FROM krasecology.eco_2018.Table_0_1_Tables"
	SQLGetRegions   string = "SELECT id, num_region, name, cast(iif(is_town = 1,1,0) as BIT) from krasecology.eco_2018.Table_0_0_Regions"
	SQLGetHeaders   string = "SELECT Table_ID,column_name, caption from krasecology.eco_2018.Table_0_2_Columns"
	SQLGetEmptyText string = "SELECT Table_ID, Region_ID, Empty_text FROM krasecology.eco_2018.Table_0_3_Empty_text WHERE Table_ID=? and Region_ID=? "
	SQLGetSQL       string = `
USE krasecology;

declare @SQL varchar(max) EXECUTE eco_2018.sp_get_table ?,
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

type TablesMeta map[int]tableInfo

func (t *TablesMeta) Fetch() error {
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
	Header            []string
	Value             [][]string
	InfoForEmptyValue string
}

func (t *Table) Fetch(info *RequestTableInfo) error {
	db := new(database)

	if err := db.connect(); err != nil {
		return fmt.Errorf("[DB] connect %v", err)
	}

	rows, err := db.Query(SQLGetSQL, info.TableID, info.RegionID)
	if err != nil {
		return fmt.Errorf("[DB] query %v", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[DB] column %v", err)
	}

	headers := (*GetHeaders())[info.TableID]

	for _, column := range columns {
		t.Header = append(t.Header, headers[column])
	}

	rawResult := make([][]byte, len(t.Header))

	dest := make([]interface{}, len(t.Header))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {

		result := make([]string, len(t.Header))

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

	if len(t.Value) > 0 {
		return nil
	}

	rows, err = db.Query(SQLGetEmptyText, info.TableID, info.RegionID)
	if err != nil {
		return err
	}

	for rows.Next() {
		rows.Scan(&t.InfoForEmptyValue)
	}

	return nil
}

type Headers map[int]map[string]string

func (h *Headers) Fetch() error {
	db := new(database)
	if err := db.connect(); err != nil {
		return err
	}

	rows, err := db.Query(SQLGetHeaders)
	if err != nil {
		return err
	}

	var tableID int
	var dbName, visName string

	*h = make(map[int]map[string]string)

	for rows.Next() {
		if err := rows.Scan(&tableID, &dbName, &visName); err != nil {
			return err
		}
		if (*h)[tableID] == nil {
			(*h)[tableID] = make(map[string]string)
		}
		(*h)[tableID][dbName] = visName
	}
	return nil
}

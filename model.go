package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
)

const (
	/*sqlGetTables    string = "SELECT ID_EP_Table, DBName, VisName FROM krasecology.eco.EP_Table"
	sqlGetRegions   string = "SELECT id, num_region, name, IsTown from krasecology.eco.EP_Region"
	sqlGetHeaders   string = "SELECT ID_EP_Table, column_name, caption from krasecology.eco.Table_Column"
	sqlGetEmptyText string = "SELECT Empty_text FROM krasecology.eco.l_EP_Tabe_EP_Region"
	sqlGetSQL       string = `
	USE krasecology;

	declare @SQL varchar(max) EXECUTE eco.sp_get_table 'babay@krasecology.ru',
	?,
	?,
	@SQL output
	EXECUTE (@sql)
	`
) */

	// TODO: переделать на уровне базы этот шлак
	sqlGetTableSpecial string = `select
	p1.[Year],
	p1.Economic_activity,
	p1.Name,
	p1.License,
	p1.Document_validity,
	p2.Standard,
	p2.Hazard_class,
	p2.Beginning_of_the_year,
	p2.Waste_generation_for_the_year,
	p2.Waste_receipt_all,
	p2.Waste_receipt_import,
	p2.Processed_waste,
	p2.Recycled_waste_all,
	p2.Recycled_waste_of_them_processed,
	p2.Waste_transfer_processing,
	p2.Waste_transfer_utilization,
	p2.Waste_transfer_neutralization,
	p2.Waste_transfer_storage,
	p2.Waste_transfer_burial,
	p2.Waste_disposal_storage,
	p2.Waste_disposal_burial,
	p2.End_of_the_year
FROM
	eco_2018.Table_1_11_part_1 p1
INNER JOIN eco_2018.Table_1_11_part_2 p2 on
	p2.ID_p3 = p1.ID
	and p2.ID_Area = ?`
	sqlGetTables  string = "SELECT Table_ID, DB_Name, VisName FROM krasecology.eco_2018.Table_0_1_Tables"
	sqlGetRegions string = "SELECT id, num_region, name, cast(iif(is_town = 1,1,0) as BIT) from krasecology.eco_2018.Table_0_0_Regions"
	sqlGetHeaders string = `select
	table_id,
	'' as DB_Name,
	'' as VisName,
	header
from
	krasecology.eco_2018.Table_0_1_Tables
where
	header is not null
union SELECT
	Table_ID,
	column_name,
	caption,
	null as header
from
	krasecology.eco_2018.Table_0_2_Columns`
	sqlGetEmptyText string = "SELECT Table_ID, Region_ID, Empty_text FROM krasecology.eco_2018.Table_0_3_Empty_text"
	sqlGetSQL       string = `
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

func (d *database) connectMSSQL() (err error) {
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

//Regions структура всех регионов
type Regions []region

//Fetch получение данных с базы
func (r *Regions) Fetch() error {
	db := new(database)

	if err := db.connectMSSQL(); err != nil {
		return fmt.Errorf("[DB] connectMSSQL %v", err)
	}

	defer db.Close()

	rows, err := db.Query(sqlGetRegions)
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

//RequestTableInfo информация от пользвателя для выдачи таблицы
type RequestTableInfo struct {
	User     string `json:"user"`
	RegionID int    `json:"region_id"`
	TableID  int    `json:"table_id"`
}

type tableInfo struct {
	DBTable string
	VisName string
}

//TablesMeta информация про таблицы
type TablesMeta map[int]tableInfo

//Fetch получение данных с базы
func (t *TablesMeta) Fetch() error {
	db := new(database)

	if err := db.connectMSSQL(); err != nil {
		return fmt.Errorf("[DB] connectMSSQL %v", err)
	}

	defer db.Close()

	rows, err := db.Query(sqlGetTables)
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

//Table отдаваемая пользователю таблица
type Table struct {
	Header            []string `json:",omitempty"`
	HeaderAsHtml      string `json:",omitempty"`
	Value             [][]string
	InfoForEmptyValue string `json:",omitempty"`
}

//Fetch получение данных с базы
func (t *Table) Fetch(info *RequestTableInfo) error {
	db := new(database)

	if err := db.connectMSSQL(); err != nil {
		return fmt.Errorf("[DB] connectMSSQL %v", err)
	}

	defer db.Close()

	var (
		rows *sql.Rows
		err  error
	)

	switch info.TableID {
	case 1014:
		rows, err = db.Query(sqlGetTableSpecial, info.RegionID)
	default:
		rows, err = db.Query(sqlGetSQL, info.TableID, info.RegionID)
	}
	if err != nil {
		return fmt.Errorf("[DB] query %v", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("[DB] column %v", err)
	}

	headers := (*GetHeaders())[info.TableID]

	if headers.HTML != "" {
		t.HeaderAsHtml = headers.HTML
	} else {
		for _, column := range columns {
			t.Header = append(t.Header, headers.Columns[column])
		}
	}

	rawResult := make([][]byte, len(columns))

	dest := make([]interface{}, len(columns))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {

		result := make([]string, len(columns))

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

	if v, ok :=  (*GetEmptyText())[info.TableID][info.RegionID]; ok {
		t.InfoForEmptyValue = v
	}

	return nil
}

//Headers кеширование всех хейдоеров
type Headers map[int]struct {
	Columns map[string]string
	HTML    string
}

//Fetch получение данных с базы
func (h *Headers) Fetch() error {
	db := new(database)
	if err := db.connectMSSQL(); err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query(sqlGetHeaders)
	if err != nil {
		return err
	}

	var tableID int
	var dbName, visName string
	var htmlHeader sql.NullString

	*h = make(map[int]struct {
		Columns map[string]string
		HTML    string
	})

	for rows.Next() {
		if err := rows.Scan(&tableID, &dbName, &visName, &htmlHeader); err != nil {
			return err
		}

		if htmlHeader.Valid {
			(*h)[tableID] = struct {
				Columns map[string]string
				HTML    string
			}{Columns: nil, HTML: htmlHeader.String}
			continue
		}

		if (*h)[tableID].Columns == nil {
			(*h)[tableID] = struct {
				Columns map[string]string
				HTML    string
			}{Columns: make(map[string]string), HTML: ""}
		}
		(*h)[tableID].Columns[dbName] = visName
	}

	return nil
}

//EmptyText кеш инфомрации для пустых таблиц
type EmptyText map[int]map[int]string

//Fetch получение данных с базы
func (e *EmptyText) Fetch() error {
	db := new(database)
	if err := db.connectMSSQL(); err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query(sqlGetEmptyText)
	if err != nil {
		return err
	}

	var tableID, regionID int
	var emptyText string

	*e = make(map[int]map[int]string)

	for rows.Next() {
		rows.Scan(&tableID, &regionID, &emptyText)

		if (*e)[tableID] == nil {
			(*e)[tableID] = make(map[int]string)
		}
		(*e)[tableID][regionID] = emptyText
	}
	return nil
}

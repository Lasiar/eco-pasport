package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"strings"
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
	sqlTest string = `sELECT 
org.org_name, 
org.Adress,  
t19.Allotted_wastewater_total, 
t19.Water_object, 
t11.Waste_generation_for_the_year, 
t8.Into_the_atmosphere, 
org.lat, 
org.lng 
from 
eco_2018.Table_0_5_Org org 
LEFT join ( 
select 
p1.Name, 
p2.Waste_generation_for_the_year 
from 
eco_2018.Table_1_11_part_1 p1 
inner join eco_2018.Table_1_11_part_2 p2 on 
p2.ID_p3 = p1.ID 
and p2.Hazard_class = 'всего' ) t11 on 
t11.Name = org.Org_name 
left join ( 
select 
pd.Allotted_wastewater_total, 
pd.Water_object, 
pd.Name 
from 
eco_2018.Table_1_9_Pollutant_discharges as pd ) t19 on 
org.Org_name = t19.name 
left join( 
select p1.Name, 
p2.Into_the_atmosphere 
from 
eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p1 p1 
inner join eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p2 p2 on 
p2.ID_p1 = p1.ID and p2.Name_of_pollutant = 'всего') t8 on 
t8.name = org.Org_name
where org.ID_Area = ?
order by org.Org_name
`

	sqlSpectial18 string = `select 
p1.Year,
p1.Economic_activ,
p1.Name,
p1.Emission_permit,
p2.Name_of_pollutant, 
p2.Thrown_without_cleaning_all, 
p2.Thrown_without_cleaning_organized, 
p2.Received_pollution_treatment, 
p2.Caught_and_rendered_harmless_all,
p2.Caught_and_rendered_harmless_utilized, 
p2.Into_the_atmosphere, 
p2.Sources_of_pollution_all, 
p2.Sources_of_pollution_organized, 
p2.MPE, 
p2.TAR,

p1.Source


from eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p1 p1
inner join eco_2018.Table_1_8_Pollutants_into_the_atmosphere_p2 p2 on p1.ID = p2.ID_p1

where p1.ID_Area = ?`

	sqSpacial13 string = `SELECT
	[Year],
	Fee_total,
	Over_limit,
	From_stationary,
	From_mobile,
	Discharges,
	Waste_disposal,
	PNG,
	[Source]
FROM
	krasecology.eco_2018.Table_3_1_Fee_for_allowable_and_excess_emissions
where ID_Area = ?`

	sqlGetCenterArea string = `SELECT lat, lng from krasecology.eco_2018.Table_0_0_Regions where ID = ?`
	sqlGetPoints     string = `SELECT org_name, Adress, Waste_generation_for_the_year,Allotted_wastewater_total , lat, lng from krasecology.eco_2018.v_PopUp_Info where ID_Area=?`
	sqlGetInfoRegion string = `SELECT Admin_center, Creation_date, Population, Area, Gross_emissions, Withdrawn_water, Discharge_volume,Formed_waste  FROM eco_2018.Table_0_4_Regions_info WHERE Region_ID=?;`

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
	p2.Recycled_waste_of_them_recycling,
	p2.Recycled_waste_of_them_processed,
	p2.Neutralized_all,
	p2.Neutralized_processed,
	p2.Waste_transfer_processing,
	p2.Waste_transfer_utilization,
	p2.Waste_transfer_neutralization,
	p2.Waste_transfer_storage,
	p2.Waste_transfer_burial,
	p2.Waste_disposal_storage,
	p2.Waste_disposal_burial,
	p2.End_of_the_year,
	p1.[Source]
FROM
	eco_2018.Table_1_11_part_1 p1
INNER JOIN eco_2018.Table_1_11_part_2 p2 on
	p2.ID_p3 = p1.ID
	and p2.ID_Area = ?`

	sqlGetTables string = "SELECT Table_ID, DB_Name, VisName FROM krasecology.eco_2018.Table_0_1_Tables"

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

	sqlGetSQL string = `
USE krasecology;

declare @SQL varchar(max) EXECUTE eco_2018.sp_get_table ?,?,
?,
@SQL output
EXECUTE (@sql)
`
)

type Database struct {
	DB  *sql.DB
	err error
}

func (d *Database) Error() string {
	return d.err.Error()
}

func NewDatabase() *Database {
	db := new(Database)
	if err := db.connectMSSQL(); err != nil {
		db.err = fmt.Errorf("[DB CONNECT] %v", err)
		return db
	}
	return db
}

func (d *Database) connectMSSQL() (err error) {
	d.DB, err = sql.Open("mssql", GetConfig().ConnStr)
	if err != nil {
		return err
	}

	return d.DB.Ping()
}

func (d *Database) close() {
	if err := d.DB.Close(); err != nil {
		log.Printf("[DB CLOSE] %v", err)
	}
}

type Region struct {
	ID        int
	NumRegion int
	Name      string
	IsTown    bool
}

//Fetch получение данных с базы
func (d *Database) GetRegions() ([]Region, error) {
	if d.err != nil {
		return nil, d.err
	}
	defer d.close()

	rows, err := d.DB.Query(sqlGetRegions)
	if err != nil {
		return nil, fmt.Errorf("[DB] query %v", err)
	}

	regions := new([]Region)

	for rows.Next() {
		r := new(Region)

		if err := rows.Scan(&r.ID, &r.NumRegion, &r.Name, &r.IsTown); err != nil {
			return nil, fmt.Errorf("[DB] scan %v", err)
		}
		*regions = append(*regions, *r)
	}
	return *regions, nil
}

//RequestTableInfo информация от пользвателя для выдачи таблицы
type TableInfo struct {
	DBTable string
	VisName string
}

//Fetch получение данных с базы
func (d *Database) GetTablesInfo() (map[int]TableInfo, error) {
	if d.err != nil {
		return nil, d.err
	}
	defer d.close()

	rows, err := d.DB.Query(sqlGetTables)
	if err != nil {
		return nil, fmt.Errorf("[DB] query %v", err)
	}

	t := make(map[int]TableInfo)

	for rows.Next() {

		row := struct {
			id      int
			dbName  string
			visName string
		}{}

		if err := rows.Scan(&row.id, &row.dbName, &row.visName); err != nil {
			return nil, fmt.Errorf("[DB] scan %v", err)
		}
		t[row.id] = TableInfo{row.dbName, row.visName}
	}

	return t, nil
}

//Table отдаваемая пользователю таблица
type Table struct {
	Header            []string `json:",omitempty"`
	HeaderAsHtml      string   `json:",omitempty"`
	Value             [][]string
	InfoForEmptyValue string `json:",omitempty"`
}

//Fetch получение данных с базы
func (d *Database) GetTable(user string, regionID int, tableID int) (*Table, error) {
	if d.err != nil {
		return nil, d.err
	}

	defer d.close()

	var (
		rows *sql.Rows
		err  error
	)

	switch tableID {
	case 1014:
		rows, err = d.DB.Query(sqlGetTableSpecial, regionID)
	case 1027:
		rows, err = d.DB.Query(sqlSpectial18, regionID)
	case 1024:
		rows, err = d.DB.Query(sqSpacial13, regionID)
	default:
		rows, err = d.DB.Query(sqlGetSQL, user, tableID, regionID)
	}
	if err != nil {
		return nil, fmt.Errorf("[DB] query %v", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("[DB] column %v", err)
	}

	headers := (*GetHeaders())[tableID]

	t := new(Table)

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
			return nil, fmt.Errorf("[DB] rows scan %v", err)
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

	if v, ok := (*GetEmptyText())[tableID][regionID]; ok {
		t.InfoForEmptyValue = v
	}

	return t, nil
}

//Headers кеширование всех хейдоеров
type Headers map[int]struct {
	Columns map[string]string
	HTML    string
}

//Fetch получение данных с базы
func (d *Database) GetHeaders() (*Headers, error) {
	if d.err != nil {
		return nil, d.err
	}
	defer d.close()

	rows, err := d.DB.Query(sqlGetHeaders)
	if err != nil {
		return nil, err
	}

	headers := new(Headers)

	*headers = make(map[int]struct {
		Columns map[string]string
		HTML    string
	})

	for rows.Next() {

		row := struct {
			tableID    int
			dbName     string
			visName    string
			htmlHeader sql.NullString
		}{}

		if err := rows.Scan(&row.tableID, &row.dbName, &row.visName, &row.htmlHeader); err != nil {
			return nil, err
		}

		if row.htmlHeader.Valid {
			x := (*headers)[row.tableID]
			x.HTML = row.htmlHeader.String
			(*headers)[row.tableID] = x
			continue
		}

		x := (*headers)[row.tableID]
		if x.Columns == nil {
			x.Columns = make(map[string]string)
		}
		x.Columns[row.dbName] = row.visName
		(*headers)[row.tableID] = x
	}
	return headers, nil
}

//EmptyText кеш инфомрации для пустых таблиц

//Fetch получение данных с базы
func (d *Database) GetTextForEmptyTable() (map[int]map[int]string, error) {
	if d.err != nil {
		return nil, d.err
	}
	defer d.close()

	rows, err := d.DB.Query(sqlGetEmptyText)
	if err != nil {
		return nil, err
	}

	textForEmptyTable := make(map[int]map[int]string)

	for rows.Next() {

		row := struct {
			tableID  int
			regionID int
			text     string
		}{}

		if err := rows.Scan(&row.tableID, &row.regionID, &row.text); err != nil {
			return nil, err
		}

		if textForEmptyTable[row.tableID] == nil {
			textForEmptyTable[row.tableID] = make(map[int]string)
		}
		textForEmptyTable[row.tableID][row.regionID] = row.text
	}
	return textForEmptyTable, nil
}

type RegionInfo struct {
	GeneralInformation struct {
		AdminCenter  string
		CreationDate int
		Population   string
		Area         string
	}
	EnvironmentalAssessment struct {
		GrossEmissions  string
		WithdrawnWater  string
		DischargeVolume string
		FormedWaste     string
	}
}

func (d *Database) GetRegionInfo(id int) (*RegionInfo, bool, error) {
	if d.err != nil {
		return nil, false, d.err
	}
	defer d.close()

	regionInfo := new(RegionInfo)

	var (
		tmpArea sql.NullString
	)
	err := d.DB.QueryRow(sqlGetInfoRegion, id).Scan(&regionInfo.GeneralInformation.AdminCenter,
		&regionInfo.GeneralInformation.CreationDate,
		&regionInfo.GeneralInformation.Population,
		&tmpArea,
		&regionInfo.EnvironmentalAssessment.GrossEmissions,
		&regionInfo.EnvironmentalAssessment.WithdrawnWater,
		&regionInfo.EnvironmentalAssessment.DischargeVolume,
		&regionInfo.EnvironmentalAssessment.FormedWaste)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	if tmpArea.Valid {
		regionInfo.GeneralInformation.Area = tmpArea.String
	}

	return regionInfo, true, nil
}

type Point struct {
	Name                      string
	Address                   string
	WasteGenerationForTheYear string
	AllottedWastewaterTotal   string
	IntoTheAtmo               string
	Latitude                  float64
	Longitude                 float64
}

func (d *Database) GetMap(regionID int) (*[]float64, []Point, error) {
	if d.err != nil {
		return nil, nil, d.err
	}

	defer d.close()

	centerArea := new([]float64)

	center := struct {
		lat sql.NullFloat64
		lng sql.NullFloat64
	}{}

	err := d.DB.QueryRow(sqlGetCenterArea, regionID).Scan(&center.lat, &center.lng)
	if err != nil {
		return nil, nil, fmt.Errorf("[DB] get center %v", err)
	}

	if !center.lat.Valid || !center.lng.Valid {
		return nil, nil, sql.ErrNoRows
	}

	*centerArea = append(*centerArea, []float64{center.lat.Float64, center.lng.Float64}...)

	rows, err := d.DB.Query(sqlTest, regionID)
	if err != nil {
		return nil, nil, err
	}

	points := new([]Point)

	var tmpName string
	var tmpWater []string

	first := true

	for rows.Next() {
		var point Point
		var (
			tmpAllottedWastewaterTotal sql.NullString
			tmpPointWasteGenerator     sql.NullString
			tmpWaterObject             sql.NullString
			tmpIntoAmto                sql.NullString
		)

		if err := rows.Scan(&point.Name, &point.Address, &tmpAllottedWastewaterTotal, &tmpWaterObject, &tmpPointWasteGenerator, &tmpIntoAmto, &point.Latitude, &point.Longitude); err != nil {
			return nil, nil, fmt.Errorf("porint %v", err)
		}

		if first {
			tmpName = point.Name

			tmpWater = append(tmpWater, fmt.Sprintf("%v - %v", tmpWaterObject.String, tmpAllottedWastewaterTotal.String))

			first = false
		}

		if point.Name != tmpName {

			if tmpPointWasteGenerator.Valid {
				point.WasteGenerationForTheYear = tmpPointWasteGenerator.String
			}

			if tmpIntoAmto.Valid {
				point.IntoTheAtmo =tmpIntoAmto.String
			}

			point.AllottedWastewaterTotal = strings.Join(tmpWater, "; ")

			*points = append(*points, point)

			tmpWater = nil

			tmpName = point.Name

		}

		tmpWater = append(tmpWater, fmt.Sprintf("%v - %v", tmpWaterObject.String, tmpAllottedWastewaterTotal.String))

	}
	return centerArea, *points, nil
}

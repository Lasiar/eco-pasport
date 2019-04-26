package model

import (
	"EcoPasport/base"
	"database/sql"
	"fmt"
	"log"
	"sync"
	// driver mssql
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	_once sync.Once
	_db   *Database
)

// Database provide access to database
type Database struct {
	db  *sql.DB
	err error
}

func (d *Database) Error() string {
	return d.err.Error()
}

func (d *Database) newDatabase() {
	if err := d.connectMSSQL(); err != nil {
		d.err = fmt.Errorf("[db CONNECT] %v", err)
	}
}

// GetDatabase get connection
func GetDatabase() *Database {
	_once.Do(func() {
		_db = new(Database)
		_db.newDatabase()
		if _db.err != nil {
			log.Fatal(_db.err)
		}
	})
	return _db
}

func (d *Database) connectMSSQL() (err error) {
	d.db, err = sql.Open("mssql", base.GetConfig().ConnStr)
	if err != nil {
		return err
	}
	return d.db.Ping()
}

// TableInfo информация от пользвателя для выдачи таблицы
type TableInfo struct {
	DBTable string
	VisName string
}

// Table отдаваемая пользователю таблица
type Table struct {
	Header            []string `json:",omitempty"`
	HeaderAsHTML      string   `json:",omitempty"`
	Value             [][]string
	InfoForEmptyValue string `json:",omitempty"`
}

// Headers кеширование всех хейдоеров
type Headers struct {
	Columns map[string]string
	HTML    string
}

// RegionInfo info by region
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

// GetRegionInfo select info databases
func (d *Database) GetRegionInfo(id int) (*RegionInfo, bool, error) {
	if d.err != nil {
		return nil, false, d.err
	}
	regionInfo := new(RegionInfo)

	var (
		tmpArea sql.NullString
	)
	err := d.db.QueryRow(sqlGetInfoRegion, id).Scan(&regionInfo.GeneralInformation.AdminCenter,
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

// Point map point
type Point struct {
	Name                      string
	Address                   string
	WasteGenerationForTheYear string
	AllottedWastewaterTotal   string
	IntoTheAtmo               string
	Latitude                  float64
	Longitude                 float64
}

type nodeEpTree struct {
	Name      string        `xml:"name,attr"`
	TableID   string        `xml:"table_id,attr" json:",omitempty"`
	TableName string        `xml:"table_name,attr"  json:",omitempty"`
	TreeItem  []*nodeEpTree `xml:"TreeItem"  json:",omitempty"`
}

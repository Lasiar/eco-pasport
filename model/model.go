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

// SetDB set current database
func (d *Database) SetDB(db *sql.DB) {
	d.db = db
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
	d.db, err = sql.Open("sqlserver", base.GetConfig().ConnStr)
	if err != nil {
		return err
	}
	return d.db.Ping()
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

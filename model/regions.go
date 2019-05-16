package model

import (
	"database/sql"
	"fmt"
)

// Region save region
type Region struct {
	ID        int
	NumRegion int
	Name      string
	IsTown    bool
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

// SelectRegions get regions
func (d *Database) SelectRegions() ([]Region, error) {
	if d.err != nil {
		return nil, d.err
	}
	rows, err := d.db.Query(sqlGetRegions)
	if err != nil {
		return nil, fmt.Errorf("[db] query %v", err)
	}
	var regions []Region
	for rows.Next() {
		r := Region{}
		if err := rows.Scan(&r.ID, &r.NumRegion, &r.Name, &r.IsTown); err != nil {
			return nil, fmt.Errorf("[db] scan %v", err)
		}
		regions = append(regions, r)
	}
	return regions, nil
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

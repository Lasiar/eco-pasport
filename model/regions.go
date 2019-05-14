package model

import (
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

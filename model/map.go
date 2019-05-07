package model

import (
	"database/sql"
	"fmt"
)

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

// GetMap get map current region by region
func (d *Database) GetMap(regionID int) (cordsCentre *[2]float64, points *[]Point, err error) {
	if d.err != nil {
		return nil, nil, d.err
	}
	cordsCentre, err = d.SelectCentreMap(regionID)
	if err != nil {
		return nil, nil, err
	}
	points, err = d.SelectPointsMap(regionID)
	if err != nil {
		return nil, nil, err
	}
	return cordsCentre, points, nil
}

// SelectPointsMap select point map to current region
func (d *Database) SelectPointsMap(regionID int) (*[]Point, error) {
	rows, err := d.db.Query(sqlGetMapPoints, sql.Named("p1", regionID))
	if err != nil {
		return nil, err
	}
	points := make(map[string]*Point)
	for rows.Next() {
		point := new(Point)
		var (
			tmpAllottedWastewaterTotal sql.NullString
			tmpPointWasteGenerator     sql.NullString
			tmpWaterObject             sql.NullString
			tmpIntoAmto                sql.NullString
		)
		err := rows.Scan(
			&point.Name,
			&point.Address,
			&tmpAllottedWastewaterTotal,
			&tmpWaterObject,
			&tmpPointWasteGenerator,
			&tmpIntoAmto,
			&point.Latitude,
			&point.Longitude,
		)
		if err != nil {
			return nil, fmt.Errorf("porint %v", err)
		}
		if tmpPointWasteGenerator.Valid {
			point.WasteGenerationForTheYear = tmpPointWasteGenerator.String
		}
		if tmpIntoAmto.Valid {
			point.IntoTheAtmo = tmpIntoAmto.String
		}
		if tmpWaterObject.Valid || tmpAllottedWastewaterTotal.Valid {
			point.AllottedWastewaterTotal += fmt.Sprintf("%v - %v; ", tmpWaterObject.String, tmpAllottedWastewaterTotal.String)
		}
		if value, ok := points[point.Name]; ok {
			value.AllottedWastewaterTotal += point.AllottedWastewaterTotal
			continue
		}
		points[point.Name] = point
	}
	p := make([]Point, len(points))
	i := 0
	for _, value := range points {
		p[i] = *value
		i++
	}
	return &p, nil
}

// SelectCentreMap get cord centre map
func (d *Database) SelectCentreMap(regionID int) (*[2]float64, error) {
	centerArea := new([2]float64)
	center := struct {
		lat sql.NullFloat64
		lng sql.NullFloat64
	}{}

	err := d.db.QueryRow(sqlGetCenterArea, regionID).Scan(&center.lat, &center.lng)
	if err != nil {
		return nil, err
	}
	if !center.lat.Valid || !center.lng.Valid {
		return nil, sql.ErrNoRows
	}
	centerArea[0], centerArea[1] = center.lat.Float64, center.lng.Float64
	return centerArea, nil
}

package model

import (
	"database/sql"
	"fmt"
	"strings"
)

// GetMap получение данных с базы
func (d *Database) GetMap(regionID int) (*[2]float64, []Point, error) {
	if d.err != nil {
		return nil, nil, d.err
	}
	centerArea := new([2]float64)
	center := struct {
		lat sql.NullFloat64
		lng sql.NullFloat64
	}{}

	err := d.db.QueryRow(sqlGetCenterArea, regionID).Scan(&center.lat, &center.lng)
	if err != nil {
		return nil, nil, err
	}
	if !center.lat.Valid || !center.lng.Valid {
		return nil, nil, sql.ErrNoRows
	}
	centerArea[0] = center.lat.Float64
	centerArea[1] = center.lng.Float64
	rows, err := d.db.Query(sqlGetMapPoints, sql.Named("p1", regionID))
	if err != nil {
		return nil, nil, err
	}
	points := new([]Point)
	var currentName string
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
			return nil, nil, fmt.Errorf("porint %v", err)
		}
		if first {
			tmpWater = append(tmpWater, fmt.Sprintf("%v - %v", tmpWaterObject.String, tmpAllottedWastewaterTotal.String))
			if tmpPointWasteGenerator.Valid {
				point.WasteGenerationForTheYear = tmpPointWasteGenerator.String
			}
			if tmpIntoAmto.Valid {
				point.IntoTheAtmo = tmpIntoAmto.String
			}
			point.AllottedWastewaterTotal = strings.Join(tmpWater, "; ")
			*points = append(*points, point)
			tmpWater = nil
			first = false
		}
		if point.Name != currentName {
			if tmpPointWasteGenerator.Valid {
				point.WasteGenerationForTheYear = tmpPointWasteGenerator.String
			}
			if tmpIntoAmto.Valid {
				point.IntoTheAtmo = tmpIntoAmto.String
			}
			point.AllottedWastewaterTotal = strings.Join(tmpWater, "; ")
			*points = append(*points, point)
			tmpWater = nil
			currentName = point.Name
		}
		tmpWater = append(tmpWater, fmt.Sprintf("%v - %v", tmpWaterObject.String, tmpAllottedWastewaterTotal.String))
	}
	return centerArea, *points, nil
}

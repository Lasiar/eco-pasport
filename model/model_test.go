package model

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func BenchmarkDatabase_GetMap(b *testing.B) {
	db := GetDatabase()
	for i := 0; i < b.N; i++ {
		_, _, err := db.GetMap(45)
		if err != nil {
			b.Errorf("error was not expected while updating stats: %s", err)
		}
	}
}

func TestDatabase_GetMap(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db := new(Database)
	db.db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.db.Close()

	centre := &[2]float64{1.23, 2.1}
	rowsCentre := mock.NewRows([]string{"lat", "lng"}).AddRow(1.23, 2.1)

	dataPoint := []struct {
		Name                      string
		Address                   string
		WasteGenerationForTheYear string
		AllottedWastewaterTotal   string
		WaterObject               string
		IntoTheAtmo               string
		Latitude                  float64
		Longitude                 float64
	}{
		{
			Name:                      "test2",
			Address:                   "Los Santos",
			WasteGenerationForTheYear: "105888.635",
			AllottedWastewaterTotal:   "88009.12",
			WaterObject:               "Енисей",
			IntoTheAtmo:               "14052.45",
			Latitude:                  1,
			Longitude:                 13.5,
		},
		{
			Name:                      "test1",
			Address:                   "Los Santos",
			WasteGenerationForTheYear: "105888.635",
			AllottedWastewaterTotal:   "100",
			WaterObject:               "Кача",
			IntoTheAtmo:               "14052.45",
			Latitude:                  1,
			Longitude:                 13.5,
		},
		{
			Name:                      "test1",
			Address:                   "Los Santos",
			WasteGenerationForTheYear: "105888.635",
			AllottedWastewaterTotal:   "1001",
			WaterObject:               "Кача",
			IntoTheAtmo:               "14052.45",
			Latitude:                  1,
			Longitude:                 13.5,
		},
		{
			Name:                      "test3",
			Address:                   "Los Santos",
			WasteGenerationForTheYear: "105888.635",
			AllottedWastewaterTotal:   "88009.12",
			WaterObject:               "Азерот",
			IntoTheAtmo:               "14052.10",
			Latitude:                  1,
			Longitude:                 13.5,
		},
		{
			Name:                      "test5",
			Address:                   "Los Santos",
			WasteGenerationForTheYear: "105888.635",
			AllottedWastewaterTotal:   "88009.12",
			WaterObject:               "еуые",
			IntoTheAtmo:               "14052.10",
			Latitude:                  1,
			Longitude:                 13.5,
		},
	}

	rowsPoint := mock.NewRows([]string{
		"org.org_name",
		"org.Adress",
		"19.Allotted_wastewater_total",
		"t19.Water_object",
		"t11.Waste_generation_for_the_year",
		"t8.Into_the_atmosphere",
		"org.lat",
		"org.lng"})

	for _, data := range dataPoint {
		rowsPoint.AddRow(
			data.Name,
			data.Address,
			data.AllottedWastewaterTotal,
			data.WaterObject,
			data.WasteGenerationForTheYear,
			data.IntoTheAtmo,
			data.Latitude,
			data.Longitude,
		)
	}

	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rowsCentre)
	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rowsPoint)

	// now we execute our method
	c, point, err := db.GetMap(45)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	for _, p := range *point {
		t.Logf("%v: %v", p.Name, p.AllottedWastewaterTotal)
	}

	if !EqualFloat(*c, *centre) {
		t.Errorf("error wrong dataPoint centre, should: %v, get: %v", centre, c)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDatabase_GetMapError(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db := new(Database)
	db.db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(sqlmock.NewRows([]string{"lat", "long"}))

	_, _, err = db.GetMap(45)
	if err != sql.ErrNoRows {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDatabase_SelectRegions(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error

	db := new(Database)
	db.db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dataSet := []Region{
		{0, 0, "first", false},
		{1, 1, "second", true},
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"num_region",
		"name",
		"isTown",
	})
	for _, d := range dataSet {
		rows.AddRow(d.ID, d.NumRegion, d.Name, d.IsTown)
	}

	mock.ExpectQuery("[a-z]*").WillReturnRows(rows)

	regions, err := db.SelectRegions()
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	if !EqualRegions(regions, dataSet) {
		t.Errorf("error wrong data centre, should: %v, get: %v", regions, dataSet)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func EqualRegions(a, b []Region) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func EqualFloat(a, b [2]float64) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

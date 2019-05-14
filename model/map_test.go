package model

import (
	"errors"
	"reflect"
	"testing"
)

var (
	errorStub  = errors.New("i`m not nil")
	centreStub = [2]float64{1.23, 4.1}
	pointsStub = []struct {
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
)

func TestDatabase_GetMap(t *testing.T) {
	type args struct {
		regionID int
	}

	db, mock, err := Init()
	if err != nil {
		t.Error(err)
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

	for _, pointStub := range pointsStub {
		rowsPoint.AddRow(
			pointStub.Name,
			pointStub.Address,
			pointStub.AllottedWastewaterTotal,
			pointStub.WaterObject,
			pointStub.WasteGenerationForTheYear,
			pointStub.IntoTheAtmo,
			pointStub.Latitude,
			pointStub.Longitude,
		)
	}
	rows := mock.NewRows([]string{"lat", "lng"}).AddRow(centreStub[0], centreStub[1])
	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rows)
	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rowsPoint)
	tests := []struct {
		name            string
		d               *Database
		args            args
		wantCordsCentre *[2]float64
		wantPoints      *[]Point
		wantErr         error
	}{
		{name: "Test error", d: &Database{err: errorStub}, wantCordsCentre: nil, wantPoints: nil, wantErr: errorStub},
		{name: "Test error", d: db, args: args{45}, wantCordsCentre: &centreStub, wantPoints: nil, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCordsCentre, _, err := tt.d.GetMap(tt.args.regionID)
			if err != tt.wantErr {
				t.Errorf("Database.GetMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCordsCentre, tt.wantCordsCentre) {
				t.Errorf("Database.GetMap() gotCordsCentre = %v, want %v", gotCordsCentre, tt.wantCordsCentre)
			}
		})
	}
}

func TestDatabase_SelectPointsMap(t *testing.T) {
	type args struct {
		regionID int
	}
	db, mock, err := Init()
	if err != nil {
		t.Error(err)
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

	for _, pointStub := range pointsStub {
		rowsPoint.AddRow(
			pointStub.Name,
			pointStub.Address,
			pointStub.AllottedWastewaterTotal,
			pointStub.WaterObject,
			pointStub.WasteGenerationForTheYear,
			pointStub.IntoTheAtmo,
			pointStub.Latitude,
			pointStub.Longitude,
		)
	}
	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rowsPoint)
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *[]Point
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
		{name: "Test with mock", args: args{45}, d: db},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.d.SelectPointsMap(tt.args.regionID)
			if err != tt.wantErr {
				t.Errorf("Database.SelectPointsMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// fixme: selectPointsMap build slice from map, order not stability
			//  if !reflect.DeepEqual(got, tt.want) {
			//	  t.Errorf("Database.SelectPointsMap() = %v, want %v", got, tt.want)
			//  }
		})
	}
}

func TestDatabase_SelectMapCentre(t *testing.T) {
	type args struct {
		regionID int
	}
	// get mock database
	db, mock, err := Init()
	if err != nil {
		t.Error(err)
	}

	// add value in mock database
	rows := mock.NewRows([]string{"lat", "lng"}).AddRow(centreStub[0], centreStub[1])
	mock.ExpectQuery("[a-z]*").WithArgs(45).WillReturnRows(rows)

	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *[2]float64
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
		{name: "Test with mock", d: db, args: args{45}, want: &centreStub, wantErr: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.SelectMapCentre(tt.args.regionID)
			if err != tt.wantErr {
				t.Errorf("Database.SelectMapCentre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.SelectMapCentre() = %v, want %v", got, tt.want)
			}
		})
	}
}

package model

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabase_SelectRegions(t *testing.T) {
	db, mock, err := Init()
	if err != nil {
		t.Error(err)
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

	tests := []struct {
		name    string
		d       *Database
		want    []Region
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
		{name: "", d: &Database{db.db, db.err}, want: dataSet, wantErr: nil},
	}

	mock.ExpectQuery("[a-z]*").WillReturnRows(rows)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.SelectRegions()
			if err != tt.wantErr {
				t.Errorf("Database.SelectRegions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.SelectRegions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetRegionInfo(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *RegionInfo
		want1   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.d.GetRegionInfo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Database.GetRegionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.GetRegionInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Database.GetRegionInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

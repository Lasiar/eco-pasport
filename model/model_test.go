package model

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/denisenkom/go-mssqldb"
)

func Init() (*Database, sqlmock.Sqlmock, error) {
	var mock sqlmock.Sqlmock
	var err error

	db := new(Database)
	db.db, mock, err = sqlmock.New()
	return db, mock, err
}

func TestDatabase_SetDB(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name string
		d    *Database
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.SetDB(tt.args.db)
		})
	}
}

func TestDatabase_Error(t *testing.T) {
	tests := []struct {
		name string
		d    *Database
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Error(); got != tt.want {
				t.Errorf("Database.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_newDatabase(t *testing.T) {
	tests := []struct {
		name string
		d    *Database
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.newDatabase()
		})
	}
}

func TestGetDatabase(t *testing.T) {
	tests := []struct {
		name string
		want *Database
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_connectMSSQL(t *testing.T) {
	tests := []struct {
		name    string
		d       *Database
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.connectMSSQL(); (err != nil) != tt.wantErr {
				t.Errorf("Database.connectMSSQL() error = %v, wantErr %v", err, tt.wantErr)
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

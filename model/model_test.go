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
	db, _, err := Init()
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		d    *Database
		args args
		want *Database
	}{
		{name: "set current db", d: db, args: args{db.db}, want: db},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.SetDB(tt.args.db)
		})
		if !reflect.DeepEqual(tt.d, tt.want) {
			t.Errorf("Database.GetRegionInfo() got = %v, want %v", tt.d, tt.want)
		}
	}
}

func TestGetDatabase(t *testing.T) {
	db, _, err := Init()
	if err != nil {
		t.Error(err)
	}

	db.SetDB(db.db)

	tests := []struct {
		name string
		want *Database
	}{
		{name: "With set current db", want: db},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

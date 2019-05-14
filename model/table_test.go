package model

import (
	"reflect"
	"testing"
)

func TestDatabase_GetTable(t *testing.T) {
	type args struct {
		regionID int
		tableID  int
	}
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *Table
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.GetTable(tt.args.regionID, tt.args.tableID)
			if err != tt.wantErr {
				t.Errorf("Database.GetTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.GetTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetPrivilege(t *testing.T) {
	type args struct {
		emailUser string
		keyUser   string
		idTable   int
	}
//	dbDefaultTrue, mockDefultTrue, err := Init()
//	if err != nil {
//		t.Error(err)
//	}
//	mock.de

	tests := []struct {
		name    string
		d       *Database
		args    args
		want    bool
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: false, wantErr: errorStub},
	//	{name: "Empty email",}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.GetPrivilege(tt.args.emailUser, tt.args.keyUser, tt.args.idTable)
			if err != tt.wantErr {
				t.Errorf("Database.GetPrivilege() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Database.GetPrivilege() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetTablesInfo(t *testing.T) {
	tests := []struct {
		name    string
		d       *Database
		want    map[int]TableInfo
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.GetTablesInfo()
			if err != tt.wantErr {
				t.Errorf("Database.GetTablesInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.GetTablesInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetHeaders(t *testing.T) {
	type args struct {
		idTable int
	}
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *Headers
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.GetHeaders(tt.args.idTable)
			if err != tt.wantErr {
				t.Errorf("Database.GetHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Database.GetHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDatabase_GetTextForEmptyTable(t *testing.T) {
	type args struct {
		idRegion int
		idTable  int
	}
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    string
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: "", wantErr: errorStub},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.GetTextForEmptyTable(tt.args.idRegion, tt.args.idTable)
			if err != tt.wantErr {
				t.Errorf("Database.GetTextForEmptyTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Database.GetTextForEmptyTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

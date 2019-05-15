package model

import (
	"reflect"
	"testing"
)

func TestGetTree(t *testing.T) {
	tests := []struct {
		name    string
		want    EpTree
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTree()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTree() = %v, want %v", got, tt.want)
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

func TestEpTree_load(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		e    *EpTree
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.load(tt.args.path)
		})
	}
}

func Test_changeName(t *testing.T) {
	type args struct {
		t     []*nodeEpTree
		table map[int]TableInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := changeName(tt.args.t, tt.args.table); (err != nil) != tt.wantErr {
				t.Errorf("changeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

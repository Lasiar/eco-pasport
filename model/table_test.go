package model

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

	dbDefaultTrue, mockDefultTrue, err := Init()
	if err != nil {
		t.Error(err)
	}
	mockDefultTrue.ExpectQuery("[a-z]*").WillReturnRows(sqlmock.NewRows([]string{"isAccess"}).AddRow(true))

	dbDefaultFalse, mockDefultFalse, err := Init()
	if err != nil {
		t.Error(err)
	}
	mockDefultFalse.ExpectQuery("[a-z]*").WillReturnRows(sqlmock.NewRows([]string{"isAccess"}).AddRow(false))
	dbWithUser, mockWithUser, err := Init()
	if err != nil {
		t.Error(err)
	}
	mockWithUser.ExpectQuery("[a-z]*").
		WillReturnRows(sqlmock.NewRows([]string{"isAcccess"}).AddRow(false))
	dtTestRegUser := time.Date(2006, 5, 10, 10, 10, 10, 10, time.Local)
	mockWithUser.ExpectQuery("[a-z]*").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email"}).AddRow(100, dtTestRegUser),
		)
	mockWithUser.ExpectQuery("[a-z]*").
		WillReturnRows(
			sqlmock.NewRows([]string{"true"}).AddRow(true),
		)
	tests := []struct {
		name    string
		d       *Database
		args    args
		want    bool
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: false, wantErr: errorStub},
		{name: "Empty email", d: dbDefaultTrue, want: true, wantErr: nil},
		{name: "Non access", d: dbDefaultFalse, want: false, wantErr: nil},
		{
			name: "With email",
			d:    dbWithUser,
			args: args{
				emailUser: "test@mail.ru",
				idTable:   100,
				keyUser:   fmt.Sprintf("%X", md5.Sum([]byte(fmt.Sprint(dtTestRegUser.Unix()+int64(100)+int64(time.Now().Month()))))),
			},
			want:    true,
			wantErr: nil,
		},
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

func TestDatabase_GetHeaders(t *testing.T) {
	type args struct {
		idTable int
	}
	var idTable = 45

	dbNtvHeaders, mockNtvHeaders, err := Init()
	if err != nil {
		t.Error(err)
	}
	rowsNtvHeaders := sqlmock.NewRows([]string{"dbName", "VisName", "Header"}).
		AddRow("test", "тест", nil).
		AddRow("I`m DBNAME", "I`m VISNAME", nil)
	mockNtvHeaders.ExpectQuery("[a-z]*").WithArgs(idTable).WillReturnRows(rowsNtvHeaders)
	wantedNtvHeaders := new(Headers)
	wantedNtvHeaders.Columns = make(map[string]string)
	wantedNtvHeaders.Columns["test"] = "тест"
	wantedNtvHeaders.Columns["I`m DBNAME"] = "I`m VISNAME"

	dbHTMLHeaders, mockHTMLHeaders, err := Init()
	if err != nil {
		t.Log(err)
	}
	rowsHTMLHeaders := sqlmock.NewRows([]string{"dbName", "VisName", "Header"}).
		AddRow("foo", "bar", "<h1>I`m HTML header</h1>")
	mockHTMLHeaders.ExpectQuery("[a-z]*").WithArgs(idTable).WillReturnRows(rowsHTMLHeaders)
	wantedHTMLHeaders := new(Headers)
	wantedHTMLHeaders.HTML = "<h1>I`m HTML header</h1>"

	tests := []struct {
		name    string
		d       *Database
		args    args
		want    *Headers
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: nil, wantErr: errorStub},
		{name: "Test native header", d: dbNtvHeaders, args: args{idTable: idTable}, want: wantedNtvHeaders, wantErr: nil},
		{name: "Test HTML header", d: dbHTMLHeaders, args: args{idTable: idTable}, want: wantedHTMLHeaders, wantErr: nil},
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
	stubArg := args{idTable: 10, idRegion: 43}

	dbEmpty, mockEmpty, err := Init()
	if err != nil {
		t.Error(err)
	}
	mockEmpty.ExpectQuery("[a-z]*").WithArgs(stubArg.idTable, stubArg.idRegion).WillReturnRows(sqlmock.NewRows([]string{""}))

	dbNmpt, mockNmpt, err := Init()
	if err != nil {
		t.Error(err)
	}
	rowsNmpt := sqlmock.NewRows([]string{"empty_text"}).AddRow("FOO_BAR")
	mockNmpt.ExpectQuery("[a-z]*").WithArgs(stubArg.idTable, stubArg.idRegion).WillReturnRows(rowsNmpt)

	tests := []struct {
		name    string
		d       *Database
		args    args
		want    string
		wantErr error
	}{
		{name: "Test error", d: &Database{err: errorStub}, want: "", wantErr: errorStub},
		{name: "Test empty", d: dbEmpty, args: stubArg, want: "", wantErr: nil},
		{name: "Test not empty", d: dbNmpt, args: stubArg, want: "FOO_BAR", wantErr: nil},
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

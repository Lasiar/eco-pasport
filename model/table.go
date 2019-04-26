package model

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"time"
)

// GetTable получение данных с базы
func (d *Database) GetTable(user string, regionID, tableID int) (*Table, error) {
	if d.err != nil {
		return nil, d.err
	}

	var (
		rows *sql.Rows
		err  error
	)

	ctx := context.Background()
	switch tableID {
	case 1014:
		rows, err = d.db.Query(sqlGetTableSpecial, regionID)
	case 1027:
		rows, err = d.db.Query(sqlSpectial18, regionID)
	case 1024:
		rows, err = d.db.Query(sqSpacial13, regionID)
	default:
		rows, err = d.db.QueryContext(ctx, sqlGetSQL, user, tableID, regionID)
	}
	if err != nil {
		return nil, fmt.Errorf("[db] query %v", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("[db] column %v", err)
	}

	headers, err := GetDatabase().GetHeaders(tableID)
	if err != nil {
		return nil, fmt.Errorf("[db] получение заголовгов: %v", err)
	}

	t := new(Table)

	if headers.HTML != "" {
		t.HeaderAsHTML = headers.HTML
	} else {
		for _, column := range columns {
			t.Header = append(t.Header, headers.Columns[column])
		}
	}

	rawResult := make([][]byte, len(columns))

	dest := make([]interface{}, len(columns))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {

		result := make([]string, len(columns))

		err = rows.Scan(dest...)
		if err != nil {
			return nil, fmt.Errorf("[db] rows scan %v", err)
		}

		for j, raw := range rawResult {
			if raw == nil {
				result[j] = ""
			} else {
				result[j] = string(raw)
			}

		}

		t.Value = append(t.Value, result)
	}

	t.InfoForEmptyValue, err = d.GetTextForEmptyTable(regionID, tableID)
	if err != nil {
		return &Table{}, err
	}

	return t, nil
}

// GetPrivilege check user privilege to the table
func (d *Database) GetPrivilege(emailUser, keyUser string, idTable int) (bool, error) {
	if d.err != nil {
		return false, d.err
	}
	var (
		isAccess           bool
		idUser             int
		dateRegisteredUser time.Time
	)
	queryDefault := fmt.Sprintf(`select T_%v
			from eco_2018.Table_0_6_Access_Right
			where ID_Role = 1009`, idTable)

	if emailUser == "" {
		if err := d.db.QueryRow(queryDefault).Scan(&isAccess); err != nil {
			return false, err
		}
		return isAccess, nil
	}
	if err := d.db.QueryRow("select ID_USER_User, DateRegistered from USER_User where EMail = ?", emailUser).Scan(&idUser, &dateRegisteredUser); err != nil {
		return false, err
	}
	verificationSum := fmt.Sprint(dateRegisteredUser.Unix() + int64(idUser) + int64(time.Now().Month()))
	if fmt.Sprintf("%X", md5.Sum([]byte(verificationSum))) != keyUser {
		return false, nil
	}
	query := fmt.Sprintf(`SELECT top 1  T_%v
	FROM eco_2018.Table_0_8_Users eu
	inner join dbo.USER_User u on eu.User_ID = u.ID_USER_User
	inner join eco_2018.Table_0_6_Access_Right ac on ac.ID_Role = eu.Role_ID
	where u.EMail = '%v'	`, idTable, emailUser)
	fmt.Println(query)
	fmt.Println(emailUser)
	if err := d.db.QueryRow(query).Scan(&isAccess); err != nil {
		if err == sql.ErrNoRows {
			if err := d.db.QueryRow(queryDefault).Scan(&isAccess); err != nil {
				return false, err
			}
			return isAccess, nil
		}
		return false, err
	}
	fmt.Println(isAccess)
	return isAccess, nil
}

// GetTablesInfo получение данных с базы
func (d *Database) GetTablesInfo() (map[int]TableInfo, error) {
	if d.err != nil {
		return nil, d.err
	}

	rows, err := d.db.Query(sqlGetTables)
	if err != nil {
		return nil, fmt.Errorf("[db] query %v", err)
	}

	t := make(map[int]TableInfo)

	for rows.Next() {

		row := struct {
			id      int
			dbName  string
			visName string
		}{}

		if err := rows.Scan(&row.id, &row.dbName, &row.visName); err != nil {
			return nil, fmt.Errorf("[db] scan %v", err)
		}
		t[row.id] = TableInfo{row.dbName, row.visName}
	}

	return t, nil
}

package model

import (
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

	switch tableID {
	case 1014:
		rows, err = d.db.Query(sqlGetTableSpecial, regionID)
	case 1027:
		rows, err = d.db.Query(sqlSpectial18, regionID)
	case 1024:
		rows, err = d.db.Query(sqSpacial13, regionID)
	default:
		rows, err = d.db.Query("declare @SQL varchar(max) EXEC  krasecology.eco_2018.sp_get_table @User, @Table_id, @Region_id, @SQL output; EXECUTE (@sql)",
			sql.Named("User", user),
			sql.Named("Table_id", tableID),
			sql.Named("Region_id", regionID),
		)
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
	where u.EMail = @email`, idTable)
	if err := d.db.QueryRow(query, sql.Named("email", emailUser)).Scan(&isAccess); err != nil {
		if err == sql.ErrNoRows {
			if err := d.db.QueryRow(queryDefault).Scan(&isAccess); err != nil {
				return false, err
			}
			return isAccess, nil
		}
		return false, err
	}
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

// GetHeaders получение данных с базы
func (d *Database) GetHeaders(idTable int) (*Headers, error) {
	if d.err != nil {
		return nil, d.err
	}

	rows, err := d.db.Query(sqlGetHeaders, idTable)
	if err != nil {
		return nil, err
	}

	headers := new(Headers)

	headers.Columns = make(map[string]string)

	for rows.Next() {

		row := struct {
			dbName     string
			visName    string
			htmlHeader sql.NullString
		}{}

		if err := rows.Scan(&row.dbName, &row.visName, &row.htmlHeader); err != nil {
			return nil, err
		}

		if row.htmlHeader.Valid {
			headers.HTML = row.htmlHeader.String
			continue
		} else {
			headers.Columns[row.dbName] = row.visName
		}
	}

	return headers, nil
}

// GetTextForEmptyTable получение данных с базы
func (d *Database) GetTextForEmptyTable(idRegion, idTable int) (string, error) {
	if d.err != nil {
		return "", d.err
	}

	var text string

	err := d.db.QueryRow(sqlGetEmptyText, idTable, idRegion).Scan(&text)
	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("[db] quer row:  %v", err)
	}

	return text, nil
}

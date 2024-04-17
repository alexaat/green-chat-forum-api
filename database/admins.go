package database

import (
	types "green-chat-forum-api/types"
	"strings"
)

//     _________users__________________________________
//     |  id      |   email   |  password | session_id |
//     |  INTEGER |   TEXT    |  TEXT     | TEXT       |

func crerateAdminsTable() error {
	sql := "CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL, session_id TEXT)"

	statement, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func SaveAdmin() (*int64, error) {
	email := "admin@gmail.com"
	passwordEcripted := "$2a$10$RY6/ndU8o15ZBukOTUSaz.XgSwjrxec//SC.52Q9JXbZIJa5VoVtq"

	statement, err := db.Prepare("INSERT INTO admins (email, password, session_id) VALUES(?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Exec(email, passwordEcripted, "")
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func GetAdminBySessionId(sessionId string) (*types.User, error) {
	if strings.TrimSpace(sessionId) == "" {
		return nil, nil
	}
	rows, err := db.Query("SELECT id, email FROM admins WHERE session_id = ? LIMIT 1", sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user *types.User = nil
	for rows.Next() {
		user = &types.User{}
		err = rows.Scan(&(user.Id), &(user.Email))
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return user, nil
}

package database

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

package database

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func OpenDatabase() {
	dbLocal, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	db = dbLocal
	createTables()
	if err != nil {
		fmt.Println(err)
	}
}

func CloseDatabase() {
	db.Close()
}

func createTables() {

	err := crerateUsersTable()
	if err != nil {
		log.Fatal(err)
	}
	err = creratePostsTable()
	if err != nil {
		log.Fatal(err)
	}
	err = crerateMessagesTable()
	if err != nil {
		log.Fatal(err)
	}
	err = crerateCategoriesTable()
	if err != nil {
		log.Fatal(err)
	}
	_ = InsertCategories([]string{"green apple", "cucumber", "kivi", "green grapes", "avocado", "broccoli", "spinach"})

	err = crerateCommentsTable()
	if err != nil {
		log.Fatal(err)
	}

	err = crerateAdminsTable()
	if err != nil {
		log.Fatal(err)
	}
	SaveAdmin()
}

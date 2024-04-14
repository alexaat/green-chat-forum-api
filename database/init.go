package database

import (
	"database/sql"
	"log"
	"fmt"
)

var db *sql.DB

func CreateDatabase() {
	dbLocal, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	db = dbLocal
	defer db.Close()
	createTables()
	if err != nil {
		fmt.Println(err)
	}
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
	_ = InsertCategories([]string{"gereen apple", "cucumber", "kivi", "green grapes", "avocado", "broccoli", "spinach"})

	err = crerateCommentsTable()
	if err != nil {
		log.Fatal(err)
	}
}

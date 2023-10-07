package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jj-attaq/todo/utils"
)
// Check out GORM -> 'Object Relational Mapping tool'

func ConnDB() (db *sql.DB) {
	/*
	   db, err := sql.Open(
	       input("Please input dbDriver, mysql for example: "),
	       input("Enter Username: ") +
	       ":" +
	       input("Enter Password: ") +
	       "@" +
	       "tcp(127.0.0.1:3306)/todo_list")
	*/
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo_list")
	utils.HandleError(err)
	//    defer db.Close()
	return db
}

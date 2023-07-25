package main

import (
	"database/sql"
	"fmt"
	"os/user"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
    fmt.Println("hello world")

    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/todo_list")
    if err != nil {
        panic(err.Error())
    }
    defer db.Close()
}

package main

import (
	"bufio"
	"database/sql"
	"fmt"
//	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Entry struct {
    id int
    item string
    finished bool
}

func commands() []string {
    var commands []string
    commands = append(commands, "add", "delete", "quit", "?", "show", "update")
    return commands
}

func input(prompt string) string {
    fmt.Printf("%s", prompt)
    stdinScanner := bufio.NewScanner(os.Stdin)
    stdinScanner.Scan()
    output := stdinScanner.Text()
    return output
}
func execCommand(input string) string {
    commands()
    var output string
    if input == commands()[0] {
        output = commands()[0]
    } else if input == commands()[1] {
        output = commands()[1]
    } else if input == commands()[2] {
        output = commands()[2]
    } else if input == commands()[3] {
        output = commands()[3]
    } else if input == commands()[4] {
        output = commands()[4]
    } else if input == commands()[5] {
        output = commands()[5]
    } else {
        fmt.Printf("'%s' is not a valid command.\n", input)
    }
    return output
}
func connDB() (db *sql.DB) {
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
    if err != nil {
        panic(err.Error())
    }
//    defer db.Close()
    return db
}
func addTable() {
    removeTable, err := connDB().Exec("DROP TABLE IF EXISTS list")
    if err != nil {
        panic(err.Error())
    }
    removeTable.RowsAffected()
    table, err := connDB().Exec("CREATE TABLE list (id MEDIUMINT NOT NULL AUTO_INCREMENT, item VARCHAR(30), finished BOOL DEFAULT 0, PRIMARY KEY (id))")
    if err != nil {
        panic(err.Error())
    }
    table.RowsAffected()
}
func addTodo(todo string) {
    if len(todo) > 30 {
        fmt.Printf("Todo item is too long, must be below 30 characters.\n")
        return 
    }
    addTodo, err := connDB().Exec("INSERT INTO list (item) VALUES (?)", todo)
    if err != nil {
        panic(err.Error())
    }
    addTodo.RowsAffected()
}
func removeTodo(todoId string) {
    id, err := strconv.Atoi(todoId)
    remove, err := connDB().Exec("DELETE FROM list WHERE ? = id", id)
    if err != nil {
        panic(err.Error())
    }
    remove.RowsAffected()
}
func updateBool(status string, todoId string) {
    var isFinished int
    if status == "y" {
        isFinished = 1
    } else if status == "n" {
        isFinished = 0
    }
    id, err := strconv.Atoi(todoId)
    prep, err := connDB().Prepare("UPDATE list SET finished = ? WHERE id = ?")
    if err != nil {
        panic(err.Error())
    }
    defer prep.Close()
    update, err := prep.Exec(isFinished, id)
    if err != nil {
        panic(err.Error())
    }
    update.RowsAffected()
}
func showList() {
    var entry Entry
    show, err := connDB().Query("SELECT * FROM list")
    for show.Next() {
        err = show.Scan(&entry.id, &entry.item, &entry.finished)
        if err != nil {
            panic(err.Error())
        }
        fmt.Printf("%+15v\n", entry)
    }
}
func eventLoop() {
    for {
        enter := execCommand(input("Enter command (type ? for list of commands): "))
        if enter == "quit" {
            connDB().Close()
            break
        } else if enter == "delete" {
            removeTodo(input("Enter numerical id of item to be deleted: "))
        } else if enter == "add" {
            addTodo(input("Enter todo: "))
        } else if enter == "?" {
            fmt.Printf("%v\n", "Available commands: " + strings.Join(commands(), ", "))
        } else if enter == "show" {
            showList()
        } else if enter == "update" {
            whichId := input("Enter id of task to be updated: ")
            answer := input("Have you finished this task? Enter y/n: ")
            if answer == "y" {
                updateBool("y", whichId)
            } else {
                updateBool("n", whichId)
            }
        }
    }
}

func main() {
    addTable()
    eventLoop()
}

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Entry struct {
    id int
    item string
}

func commands() []string {
    var commands []string
    commands = append(commands, "add", "delete", "quit", "?", "show")
    return commands
}

func input(prompt string) string {
    fmt.Println(prompt)
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
    } else {
        log.Println("Oops, that is not a valid command.")
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
    table, err := connDB().Exec("CREATE TABLE list (id INT, item VARCHAR(30))")
    if err != nil {
        panic(err.Error())
    }
    table.RowsAffected()
}
func genId(idCount int) (n int) {
    var entry Entry
    count, err := connDB().Query("SELECT COUNT(id) FROM list WHERE id > 0")
    for count.Next() {
        err = count.Scan(&entry.id)
        if err != nil {
            panic(err.Error())
        }
    }
    fmt.Printf("%T\n", count)
    return entry.id + 1
}
func addTodo(todo string) {
    addTodo, err := connDB().Exec("INSERT INTO list VALUES (?, ?)", genId(idCount), todo)
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
func showList() {
    var entry Entry
    show, err := connDB().Query("SELECT * FROM list")
    for show.Next() {
        err = show.Scan(&entry.id, &entry.item)
        if err != nil {
            panic(err.Error())
        }
        fmt.Println(entry.id, entry.item)
    }
}
func eventLoop() {
    for {
        enter := execCommand(input("Enter command (type ? for list of commands): "))
        if enter == "quit" {
            break
        } else if enter == "delete" {
            removeTodo(input("Enter numerical id of item to be deleted: "))
        } else if enter == "add" {
            addTodo(input("Enter todo: "))
        } else if enter == "?" {
            fmt.Printf("%v\n", "Available commands: " + strings.Join(commands(), ", "))
        } else if enter == "show" {
            showList()
        } else {
            break
        }
    }
    /*
    for i := 0; i < todoNum; i++ {
        addTodo(input("Enter todo: "))
    }
    */
}

var idCount int
func main() {
    // ID GENERATION NEEDS TO BE FIXED!!!
    // DOESN'T UPDATE WHEN ITEMS ARE REMOVED FROM LIST
    addTable()
    eventLoop()
}

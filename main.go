package main

import (
//    "github.com/jj-attaq/todo/utils"
    "github.com/jj-attaq/todo/src/commands"
    "github.com/jj-attaq/todo/src/database"
	"fmt"
//	"log"
	"strings"

)

type Entry struct {
    id int
    item string
    finished bool
}

func eventLoop() {
    for {
        enter := commands.input("Enter command (type ? for list of commands): ")
        if enter == "quit" {
            break
        } else if enter == "delete" {
            commands.removeTodo(input("Enter numerical id of item to be deleted: "))
        } else if enter == "add" {
            commands.addTodo(input("Enter todo: "))
        } else if enter == "?" {
            fmt.Printf("%v\n", "Available commands: " + strings.Join(commands(), ", "))
        } else if enter == "show" {
            commands.showList()
        } else if enter == "update" {
            whichId := commands.input("Enter id of task to be updated: ")
            answer := commands.input("Have you finished this task? Enter y/n: ")
            if answer == "y" {
                commands.updateBool("y", whichId)
            } else {
                commands.updateBool("n", whichId)
            }
        }
    }
}

func main() {
    commands.addTable()
    eventLoop()
    database.connDB().Close()
}

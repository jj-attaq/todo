package main

import (
//    "github.com/jj-attaq/todo/utils"
    "github.com/jj-attaq/todo/src/commands"
    "github.com/jj-attaq/todo/src/database"
	"fmt"
    "strings"
//	"log"

)


func eventLoop() {
    input := commands.Input
    for {
        enter := input("Enter command (type ? for list of commands): ")
        if enter == "quit" {
            break
        } else if enter == "delete" {
            commands.RemoveTodo(input("Enter numerical id of item to be deleted: "))
        } else if enter == "add" {
            commands.AddTodo(input("Enter todo: "))
        } else if enter == "?" {
            fmt.Printf("%v\n", "Available commands: " + strings.Join(commands.Commands(), ", "))
        } else if enter == "show" {
            commands.ShowList()
        } else if enter == "update" {
            whichId := input("Enter id of task to be updated: ")
            answer := input("Have you finished this task? Enter y/n: ")
            if answer == "y" {
                commands.UpdateBool("y", whichId)
            } else {
                commands.UpdateBool("n", whichId)
            }
        }
    }
}

func main() {
    commands.AddTable()
    eventLoop()
    database.ConnDB().Close()
}

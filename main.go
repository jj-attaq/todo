package main

import (
	//    "github.com/jj-attaq/todo/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jj-attaq/todo/src/commands"
	"github.com/jj-attaq/todo/src/database"
	//	"log"
)

func eventLoop() {
	input := commands.Input
	for {
		enter := commands.ExecCommand(input("Enter command (type ? for list of commands): "))
		if enter == "quit" {
			break
		} else if enter == "delete" {
			commands.RemoveTodo(input("Enter numerical id of item to be deleted: "))
		} else if enter == "add" {
			commands.AddTodo(input("Enter todo: "))
		} else if enter == "?" {
			fmt.Printf("%v\n", "Available commands: "+strings.Join(commands.Commands(), ", "))
		} else if enter == "show" {
			commands.ShowList()
		} else if enter == "update" {
			whichId := input("Enter uuid of task to be updated: ")
			commands.UpdateBool(whichId)
			/*
				whichId := input("Enter id of task to be updated: ")
				answer := input("Have you finished this task? Enter y/n: ")
				if answer == "y" {
					commands.UpdateBool("y", whichId)
				} else {
					commands.UpdateBool("n", whichId)
				}
			*/
		}
	}
}

/*
 */

func main() {
	go func() {
		router := gin.Default()
		router.GET("/todo-list", commands.GetTodos)
		router.POST("/todo-list", commands.AddTask)
		router.GET("/todo-list/:uuid", commands.GetOneTodo)
		router.GET("/todo-list/remove/:uuid", commands.RemoveTask)
		router.GET("/todo-list/update/:uuid", commands.UpdateTask)

		router.Run("localhost:8080")
	}()
	//	commands.AddTable() // put back after testing
	eventLoop()
	database.ConnDB().Close()
}

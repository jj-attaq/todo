package main

import (
	//    "github.com/jj-attaq/todo/utils"
	"fmt"
	"net/http"
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

/*
func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", "Todo list: ")
	fmt.Fprintf(w, "%s", commands.ShowJSON())

	commands.ShowJSON()
}
*/

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, commands.ShowJSON())
}

func main() {

	go func() {
		router := gin.Default()
		router.GET("/show", getTodos)
		router.POST("/show", commands.AddTask)

		router.Run("localhost:8080")
		/*
			http.HandleFunc("/", greet)
			http.ListenAndServe(":8080", nil)
		*/
	}()
	commands.AddTable()
	eventLoop()
	database.ConnDB().Close()
}

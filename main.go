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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func main() {
	go func() {
		router := gin.Default()
		router.Use(CORSMiddleware())
		/*
			router.Use(cors.New(cors.Config{
				AllowOrigins:     []string{"http://localhost:5173/"},
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
				AllowHeaders:     []string{"Origin", "Access-Control-Allow-Origin"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: true,
				AllowOriginFunc: func(origin string) bool {
					return origin == "https://github.com"
				},
				MaxAge: 12 * time.Hour,
			}))
		*/
		/*
			corsConfig := cors.DefaultConfig()

			corsConfig.AllowOrigins = []string{"http://localhost:5173/"}
			corsConfig.AllowCredentials = true
			corsConfig.AddAllowMethods("OPTIONS")
			router.Use(cors.New(corsConfig))
		*/

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

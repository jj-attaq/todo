package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-attaq/todo/commands"
	"github.com/jj-attaq/todo/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnDB()
}

func main() {
	router := gin.Default()
	router.GET("/todos", commands.GetAllTodos)
	router.GET("/todos/:id", commands.GetTodo)
	router.POST("/addTodo", commands.AddTodo)
	router.PUT("/updateTodo/:id", commands.UpdateTodo)
	router.DELETE("/deleteTodo/:id", commands.DeleteTodo)
	router.Run()
}

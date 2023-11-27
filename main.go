package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jj-attaq/todo/commands"
	"github.com/jj-attaq/todo/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnDB()
}

func main() {
	fmt.Println("Starting server...")
    // commands.Test("this", "is", "a", "test")
    // commands.Test("1", "2", "3", "4")
    // commands.Test("Status Code", "")

	// Handlers
    router := gin.Default()
    // Task management handlers
	router.GET("/todos", commands.GetAllTodos) // might need to make POST because of user spec in json body, or use GET with :UserID in GET call
	router.GET("/todos/:id", commands.GetTodo)
	router.POST("/addTodo", commands.AddTodo) // UserID in json body
	router.PUT("/updateTodo", commands.UpdateTodo)
	router.PUT("/updateUser", commands.UpdateUser)
	router.DELETE("/deleteTodo", commands.DeleteTodo)
    // Auth handlers
    router.POST("/register", commands.Register)
    router.POST("/login", commands.Login)
    router.GET("/welcome", commands.Welcome)
    router.POST("/refresh", commands.Refresh)
    router.GET("/logout", commands.Logout)

	// Graceful shutdown
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", server.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}

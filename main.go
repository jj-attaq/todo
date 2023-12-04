package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func formatKeyAndMap(m map[string][]string) []string {
    var str  string
    var result []string
    for key, el := range m {
        resEl := strings.Join(el, " ")
        str  = key + ": " + resEl
        result = append(result, str )
    }
    return result
}

func logKeyAndMap(arr []string){
    log.Printf("Beginning of log entry: \n--------\n")
    for _, el  := range arr {
        fmt.Printf("%v\n", el)
    }
    fmt.Println("--------")
    log.Printf("End of log entry.\n")
}

func middleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            // log.Println(r.Cookies())
            // if logged in run commands.Refresh handler???
            next.ServeHTTP(w, r)
            
            /* log.Println(r.Method)
            log.Println(r.Host)
            logKeyAndMap(formatKeyAndMap(r.Header)) */
        },
    )
}

func main() {
    // Basic request logging
    gin.DisableConsoleColor()
    file, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	fmt.Println("Starting server...")
	// Handlers
    router := gin.Default()
    // configuredRouter := middleWare(router)

    // Task management handlers
	router.GET("/todos", commands.GetAllTodos) // might need to make POST because of user spec in json body, or use GET with :UserID in GET call
	router.GET("/todos/:id", commands.GetTodo)
	router.POST("/addTodo", commands.AddTodo) // UserID in json body
	router.PUT("/updateTodo", commands.UpdateTodo)
	router.DELETE("/deleteTodo", commands.DeleteTodo)
    // Auth handlers
    router.POST("/register", commands.Register)
    router.POST("/login", commands.Login)
    router.GET("/welcome", commands.Welcome)
    router.POST("/refresh", commands.Refresh)
    router.PUT("/updateUser", commands.UpdateUser)
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

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
	"text/tabwriter" // for formatting logs into columns: https://blog.el-chavez.me/2019/05/05/golang-tabwriter-aligned-text/
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jj-attaq/todo/commands"
	"github.com/jj-attaq/todo/initializers"
	"github.com/rs/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnDB()
}

func formatKeyAndMap(m map[string][]string) []string {
	var str string
	var result []string
	for key, el := range m {
		resEl := strings.Join(el, " ")
		str = key + ": " + resEl
		result = append(result, str)
	}
	return result
}

func logKeyAndMap(arr []string) {
	log.Printf("Beginning of log entry: \n--------\n")
	for _, el := range arr {
		fmt.Printf("%v\n", el)
	}
	fmt.Println("--------")
	log.Printf("End of log entry.\n")
}

// CORS middleware
func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			enCors := cors.Default().Handler(next)
			// next.ServeHTTP(w, r)
			enCors.ServeHTTP(w, r)
		},
	)
}
func logPrefix(str string, output any, isLog bool) {
	if isLog {
		log.Printf(str+": %v\n", output)
	} else {
		fmt.Printf(str+": %v\n", output)
	}
}

// https://www.enterpriseready.io/features/audit-log/
func loggingMiddleware(next http.Handler, file io.Writer) http.Handler {
	defaultWriter := io.MultiWriter(file, os.Stdout)
	writer := tabwriter.NewWriter(defaultWriter, 0, 8, 1, '\t', tabwriter.AlignRight)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			// Actor - username, uuid, api token
			// Group - organization, team, account for team admin history
			// Where - IP address, device ID, country
			time := time.Now()
			fmt.Fprintln(writer, "Time: ", time)
			fmt.Fprintln(writer, "Header: ")
			for key, el := range r.Header {
				fmt.Fprintln(writer, key+"\t"+strings.Join(el, ""))
				// fmt.Println("Key: ", key, " => ", el)
			}
			// When - NTP synced server time of the event
			// Target - object or resource being changed, the 'noun'
			// Action - the verb, how was the object changed
			fmt.Fprintf(writer, "\n")
			fmt.Fprintf(writer, "Method: %v\n", r.Method)
			// Action type - C, R, U, or D
			// Event Name
			// Description
			// --- Optional ---
			// Server
			// Version
			// Protocols
			fmt.Fprintf(writer, "Protocol: %v\n", r.Proto)
			fmt.Fprintf(writer, "URL: %v\n", r.URL)
			// Global Actor ID
			// log.Printf("Header: %v\n", r.Header)
			fmt.Fprintln(writer, "--------")
			writer.Flush()
		},
	)
}

func main() {
	// Basic request logging
	/* gin.DisableConsoleColor()
	   file, err := os.Create("gin.log")
	   if err != nil {
	       log.Panic(err)
	   }
	   gin.DefaultWriter = io.MultiWriter(file, os.Stdout) */

	file, err := os.OpenFile("middleWareLog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fmt.Println("Starting server...")
	// Handlers
	router := gin.Default()
	configuredRouter := loggingMiddleware(middleWare(router), file)

	// Task management handlers
	router.GET("/todos", commands.GetAllTodos) // might need to make POST because of user spec in json body, or use GET with :UserID in GET call
	router.GET("/todos/:id", commands.GetTodo)
	router.POST("/addTodo", commands.AddTodo)        // UserID in json body
	router.PATCH("/updateTodo", commands.UpdateTodo) // Maybe PATCH req?
	router.DELETE("/deleteTodo", commands.DeleteTodo)
	// Auth handlers
	router.POST("/register", commands.Register)
	router.POST("/login", commands.Login)
	router.GET("/welcome", commands.Welcome)
	router.POST("/refresh", commands.Refresh)
	router.PATCH("/updateUser", commands.UpdateUser) // Maybe PATCH req?
	router.GET("/logout", commands.Logout)
    // Team Handling
    router.POST("/createTeam", commands.CreateTeam)

	// Graceful shutdown
	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: configuredRouter, // Was router before
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

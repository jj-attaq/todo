package main

import (
	"github.com/jj-attaq/todo/initializers"
	"github.com/jj-attaq/todo/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnDB()
}

func main() {
    initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	initializers.DB.AutoMigrate(&models.Todo{})
    initializers.DB.AutoMigrate(&models.User{})
}

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
    if initializers.DB.Migrator().HasTable(&models.Todo{}) {
        initializers.DB.Migrator().DropTable(&models.Todo{})
    }
    if initializers.DB.Migrator().HasTable(&models.User{}) {
        initializers.DB.Migrator().DropTable(&models.User{})
    }
    initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	initializers.DB.AutoMigrate(&models.Todo{})
    initializers.DB.AutoMigrate(&models.User{})
}

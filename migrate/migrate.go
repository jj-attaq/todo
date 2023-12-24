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
	// RAW SQL FOR PARTIAL INDEXING // NEED FOR SOFT DELETION OF UNIQUE IDs
	sessionPartial := initializers.DB.Exec("CREATE UNIQUE INDEX \"sessions_user_id_unique\" ON sessions(user_id, deleted_at) WHERE deleted_at IS NULL;")
	if sessionPartial.Error != nil {
		panic(sessionPartial.Error)
	}
	// change for final version
	if initializers.DB.Migrator().HasTable(&models.Todo{}) {
		initializers.DB.Migrator().DropTable(&models.Todo{})
	}
	if initializers.DB.Migrator().HasTable(&models.TeamTodo{}) {
		initializers.DB.Migrator().DropTable(&models.TeamTodo{})
	}
	if initializers.DB.Migrator().HasTable(&models.User{}) {
		initializers.DB.Migrator().DropTable(&models.User{})
	}
	if initializers.DB.Migrator().HasTable(&models.Team{}) {
		initializers.DB.Migrator().DropTable(&models.Team{})
	}
	if initializers.DB.Migrator().HasTable(&models.Permissions{}) {
		initializers.DB.Migrator().DropTable(&models.Permissions{})
	}
	if initializers.DB.Migrator().HasTable(&models.Session{}) {
		initializers.DB.Migrator().DropTable(&models.Session{})
	}
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	initializers.DB.AutoMigrate(&models.Todo{})
	initializers.DB.AutoMigrate(&models.TeamTodo{})
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Team{})
	initializers.DB.AutoMigrate(&models.Permissions{})
	initializers.DB.AutoMigrate(&models.Session{})
}

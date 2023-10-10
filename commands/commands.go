package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/jj-attaq/todo/initializers"
	"github.com/jj-attaq/todo/models"
)

func AddTodo(c *gin.Context) {
	var body struct {
		Title string
		Body  string
	}
	c.Bind(&body)

	todo := models.Todo{Title: body.Title, Description: body.Body}
	result := initializers.DB.Create(&todo)
	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{
		"todo": todo,
	})
}

func GetAllTodos(c *gin.Context) {
	// Get todos
	var todos []models.Todo
	result := initializers.DB.Find(&todos)
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Respond
	c.JSON(200, gin.H{
		"todos": todos,
	})
}

func GetTodo(c *gin.Context) {
	// Get id off url
	id := c.Param("id")
	// Get todo
	var todo models.Todo
	result := initializers.DB.First(&todo, "id = ?", id)
	if result.Error != nil {
		c.Status(400)
		return
	}

	//Respond
	c.JSON(200, gin.H{
		"todo": todo,
	})
}

func UpdateTodo(c *gin.Context) {
	// status := c.Param("status")
	id := c.Param("id")
	var todo models.Todo

	result := initializers.DB.First(&todo, "id = ?", id)
	if result.Error != nil {
		c.Status(400)
		return
	} else {
		if todo.Status == false {
			initializers.DB.Model(&todo).Update("status", true)
		} else {
			initializers.DB.Model(&todo).Update("status", false)
		}

		c.JSON(200, gin.H{
			"todo": todo,
		})
	}
}

func DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	var todo models.Todo

	result := initializers.DB.Delete(&todo, "id = ?", id)
	if result.Error != nil {
		c.Status(400)
		return
	} else {
		c.Status(200)
	}
}

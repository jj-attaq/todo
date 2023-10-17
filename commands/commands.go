package commands

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jj-attaq/todo/initializers"
	"github.com/jj-attaq/todo/models"
	"golang.org/x/crypto/bcrypt"
)

func AddTodo(c *gin.Context) {
	var body struct {
        User uuid.UUID
		Title string
		Body  string
	}
	c.Bind(&body)

	todo := models.Todo{UserID: body.User,Title: body.Title, Description: body.Body}
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
	var body struct {
        User uuid.UUID
    }
    c.Bind(&body)
    userID := c.Param("UserID")
	// Get todos
	var todos []models.Todo
    fmt.Println(userID)
	result := initializers.DB.Find(&todos).Where("UserID = ?", userID)
    
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

func Register(c *gin.Context) {
    var body struct {
        Name string
        Email string 
        Password string
    }
    c.Bind(&body)
    pw := encrypt(body.Password)
    body.Password = pw
    user := models.User{Username: body.Name, Email: body.Email, Password: body.Password}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.Status(400)
		return
	}
    c.JSON(200, gin.H{
        "user": user,
    })
}

func Login(c *gin.Context) {
    // must be a better way to do this, check out gin.Accounts?
    var body struct {
        Email string 
        Password string
    }
    c.Bind(&body)

    email := body.Email
    // pw := c.Param("password")
    var user models.User
    result := initializers.DB.First(&user, "email = ?", email)
    /* fmt.Println(user.Password)
    fmt.Println(body.Password) */

    match := checkPasswordHash(body.Password, user.Password)
	if (result.Error != nil) || (!match) {
		c.Status(400)
		return
	}
    c.JSON(200, gin.H{
        "user": user,
    })
}

func encrypt(str string) string {
    hashed, err := bcrypt.GenerateFromPassword([]byte(str), 8)
    if err != nil {
        log.Fatal(err)
    }
    hashedPw := string(hashed[:])
    return hashedPw
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

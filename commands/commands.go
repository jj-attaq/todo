package commands

import (
	// "encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jj-attaq/todo/initializers"
	"github.com/jj-attaq/todo/models"
	"golang.org/x/crypto/bcrypt"
)

type reqBody struct {
    ID uuid.UUID // Todo
    Title string // Todo
    Body  string // Todo
    User uuid.UUID // User
    Name string // User
    Email string // User
    NewEmail string // User
    Password string // User
    NewPassword string // User
}

func AddTodo(c *gin.Context) {
    var body reqBody
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
    /* // Doesn't work if it's a GET request, requires POST
    var body struct {
        User uuid.UUID
    }
    c.Bind(&body)
    log.Printf("This is the userID input into Postman: %v\n", body.User) */
    userID := c.Param("UserID")
    // log.Printf("This is the userID input into Postman: %v\n", userID)
	// Get todos
	var todos []models.Todo
    // result := initializers.DB.Find(&todos) // ALL todos regardless of user
	result := initializers.DB.Where("user_id = ?", userID).Find(&todos)

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
    userID := c.Param("UserID")
	// Get todo
	var todo models.Todo
	result := initializers.DB.Where("user_id = ?", userID).First(&todo, "id = ?", id)
	if result.Error != nil {
		c.Status(400)
		return
	}

	//Respond
	c.JSON(200, gin.H{
		"todo": todo,
	})
}

func DeleteTodo(c *gin.Context) {
	// id := c.Param("id")
    var body reqBody
    c.Bind(&body)

	var todo models.Todo

	result := initializers.DB.Where("user_id = ?", body.User).Delete(&todo, "id = ?", body.ID)

    // current version returns 200 with empty request body, change later
	if result.Error != nil {
		c.Status(400)
		return
	} else {
		c.Status(200)
	}
}

func Register(c *gin.Context) {
    var body reqBody
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
    var body reqBody
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

func UpdateTodo(c *gin.Context) {
	// id := c.Param("id")
    var body reqBody
    c.Bind(&body)

	var todo models.Todo

	result := initializers.DB.Where("user_id = ?", body.User).First(&todo, "id = ?", body.ID)
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

func UpdateUser(c *gin.Context) {
    var body reqBody
    c.Bind(&body)
    pw := encrypt(body.NewPassword)

    var user models.User
    result := initializers.DB.First(&user, "id = ?", body.User)
    
    /* info, err := json.MarshalIndent(user, "\t", "")
    if err != nil {
        log.Println(err)
    }
    fmt.Println(string(info[:])) */

    if result.Error != nil { // NEVER FORGET .ERROR AFTER RESULT !!! I REPEAT, NEVER EVER FORGET .ERROR !!!
		c.Status(400)
		return
    } else {
        if len(body.NewPassword) > 0 {
            initializers.DB.Model(&user).Update("password", pw)
        }
        if len(body.NewEmail) > 0 {
            initializers.DB.Model(&user).Update("email", body.NewEmail)
        }
		c.JSON(200, gin.H{
			"user": user,
		})
    }
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

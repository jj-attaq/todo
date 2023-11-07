package commands

import (
	// "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jj-attaq/todo/initializers"
	"github.com/jj-attaq/todo/models"
	"golang.org/x/crypto/bcrypt"
)

// https://www.sohamkamani.com/golang/session-cookie-authentication/#overview

type session struct {
    Email string
    expiry time.Time
}

var sessions = map[string]session{}

func (s session) isExpired() bool {
    return s.expiry.Before(time.Now())
}


type todoReq struct {
	ID    uuid.UUID // Todo
	Title string    // Todo
	Body  string    // Todo
	User  uuid.UUID // User
}
type userReq struct {
	User        uuid.UUID // User
	Name        string    // User
	Email       string    // User
	NewEmail    string    // User
	Password    string    // User
	NewPassword string    // User
}

func AddTodo(c *gin.Context) {
	var todoBody todoReq
	c.Bind(&todoBody)
    /* var userBody userReq
    c.Bind(&userBody) */

	todo := models.Todo{UserID: todoBody.User, Title: todoBody.Title, Description: todoBody.Body}
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
	// result := initializers.DB.Find(&todos) // ALL todos regardless of user

    var user models.User
    initializers.DB.Where("user_id = ?", userID).Find(&user)
	result := initializers.DB.Where("user_id = ?", userID).Find(&user.Todos)
    /* func() {
        log.Println("User's todos:")
        for _, el := range user.Todos {
            log.Printf("%+v %+v\n", el.UserID, el.Title)
        }
    }() */

	if result.Error != nil {
		c.Status(400)
		return
	}
	// Respond
	c.JSON(200, gin.H{
        "user": user.Todos,
		// "todos": todos,
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
	var todoBody todoReq
	c.Bind(&todoBody)

	var todo models.Todo

	result := initializers.DB.Where("user_id = ?", todoBody.User).Delete(&todo, "id = ?", todoBody.ID)

	// current version returns 200 with empty request body, change later
	if result.Error != nil {
		c.Status(400)
		return
	} else {
		c.Status(200)
	}
}

func Register(c *gin.Context) {
	var userBody userReq
	c.Bind(&userBody)

	pw := encrypt(userBody.Password)
	userBody.Password = pw
	user := models.User{Username: userBody.Name, Email: userBody.Email, Password: userBody.Password}

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
    // https://www.sohamkamani.com/golang/password-authentication-and-storage/

	// must be a better way to do this, check out gin.Accounts?
	var userBody userReq
	c.Bind(&userBody)

	email := userBody.Email
	// pw := c.Param("password")
	var user models.User
	result := initializers.DB.First(&user, "email = ?", email)

	match := checkPasswordHash(userBody.Password, user.Password)
	if (result.Error != nil) || (!match) {
		c.Status(400)
		return
	}

    // NEW!
    sessionToken := uuid.NewString()
    expiresAt := time.Now().Add(120 * time.Second)

    sessions[sessionToken] = session{
        Email: userBody.Email,
        expiry: expiresAt,
    }

    http.SetCookie(c.Writer, &http.Cookie{
        Name: "session_token",
        Value: sessionToken,
        Expires: expiresAt,
    })
    log.Println(sessions)
    // END NEW!

	c.JSON(200, gin.H{
		"user": user,
        "token":sessionToken,
	})
}

func Welcome(c *gin.Context) {
    cookie, err := c.Request.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Status(401)
            return
        }
        c.Status(400)
        return
    }

    sessionToken := cookie.Value

    userSession, exists := sessions[sessionToken]
    if !exists {
        c.Status(401)
        return
    }
    if userSession.isExpired() {
        delete(sessions, sessionToken)
        c.Status(401)
        return
    }

    log.Println(sessions)
	c.JSON(200, gin.H{
        "welcome": userSession.Email,
	})
}

func Refresh(c *gin.Context) {
    cookie, err := c.Request.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Status(401)
            return
        }
        c.Status(400)
        return
    }

    sessionToken := cookie.Value

    userSession, exists := sessions[sessionToken]
    if !exists {
        c.Status(401)
        return
    }
    if userSession.isExpired() {
        delete(sessions, sessionToken)
        c.Status(401)
        return
    }

    newSessionToken := uuid.NewString()
    expiresAt := time.Now().Add(120 * time.Second)

    sessions[newSessionToken] = session{
        Email: userSession.Email,
        expiry: expiresAt,
    }
    delete(sessions, sessionToken)
    http.SetCookie(c.Writer, &http.Cookie{
        Name: "session_token",
        Value: newSessionToken,
        Expires: expiresAt,
    })
}

func Logout(c *gin.Context) {
    cookie, err := c.Request.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Status(401)
            return
        }
        c.Status(400)
        return
    }

    sessionToken := cookie.Value
    delete(sessions, sessionToken)
    http.SetCookie(c.Writer, &http.Cookie{
        Name: "session_token",
        Value: "",
        Expires: time.Now(),
    })
}

func UpdateTodo(c *gin.Context) {
	// id := c.Param("id")
	var todoBody todoReq
	c.Bind(&todoBody)

	var todo models.Todo

	result := initializers.DB.Where("user_id = ?", todoBody.User).First(&todo, "id = ?", todoBody.ID)
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
	var userBody userReq
	c.Bind(&userBody)
	pw := encrypt(userBody.NewPassword)

	var user models.User
	result := initializers.DB.First(&user, "id = ?", userBody.User)

	/* info, err := json.MarshalIndent(user, "\t", "")
	   if err != nil {
	       log.Println(err)
	   }
	   fmt.Println(string(info[:])) */
	match := checkPasswordHash(userBody.Password, user.Password)
	log.Println(match)

	if result.Error != nil { // NEVER FORGET .ERROR AFTER RESULT !!! I REPEAT, NEVER EVER FORGET .ERROR !!!
		c.Status(400)
		return
	} else {
		if len(userBody.NewPassword) > 0 {
			initializers.DB.Model(&user).Update("password", pw)
		}
		if len(userBody.NewEmail) > 0 {
			initializers.DB.Model(&user).Update("email", userBody.NewEmail)
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

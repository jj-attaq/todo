package commands

import (
	"errors"
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
type jsonMessage struct {
    m map[string]string
}

func AddTodo(c *gin.Context) {
    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        var todoBody todoReq
        c.Bind(&todoBody)
        // var userBody userReq
        // c.Bind(&userBody)

        uid, err := c.Cookie("session_token")
        if err != nil {
            panic(err)
        }

        log.Printf("UID: %+v", uid)

        // todo := models.Todo{UserID: todoBody.User, Title: todoBody.Title, Description: todoBody.Body}
        todo := models.Todo{UserID: user.ID, Title: todoBody.Title, Description: todoBody.Body}
        result := initializers.DB.Create(&todo)
        if result.Error != nil {
            c.Status(400)
            return
        }
        c.JSON(200, gin.H{
            "todo": todo,
        })
    }
}

func GetAllTodos(c *gin.Context) {
    // // Doesn't work if it's a GET request, requires POST
    //    var body struct {
    //        User uuid.UUID
    //    }
    //    c.Bind(&body)
    //    log.Printf("This is the userID input into Postman: %v\n", body.User)

    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        result := initializers.DB.Where("user_id = ?", user.ID).Find(&user.Todos)
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
}

func GetTodo(c *gin.Context) {
    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        id := c.Param("id")
        var todo models.Todo
        result := initializers.DB.Where("user_id = ?", user.ID).First(&todo, "id = ?", id)
        if result.Error != nil {
            c.Status(400)
            return
        }

        //Respond
        c.JSON(200, gin.H{
            "todo": todo,
        })
    }
}

func DeleteTodo(c *gin.Context) {
    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        // id := c.Param("id")
        var todoBody todoReq
        c.Bind(&todoBody)

        var todo models.Todo

        result := initializers.DB.Where("user_id = ?", user.ID).Delete(&todo, "id = ?", todoBody.ID)

        // current version returns 200 with empty request body, change later
        if result.Error != nil {
            c.Status(400)
            return
        } else {
            c.Status(200)
        }
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
func mkSession(user uuid.UUID, email string, expiresAt time.Time) models.Session {
    return models.Session{
        UserID: user,
        Email: email,
        Expiry: expiresAt,
    }
}

func RemoveExpiredCookies() { // helper function
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
	if (result.Error != nil) {
		c.Status(400) // maybe different code
		return
	}
    if !match {
        c.Status(401)
		return
    }

    // NEW!
    expiresAt := time.Now().Add(120 * time.Second)

    var session models.Session
    session = mkSession(user.ID, userBody.Email, expiresAt)

    userHasCookie := initializers.DB.First(&session, "user_id = ?", user.ID)
    log.Printf("Bind existing session in DB to session var: %+v\n", session.ID)

    mkCookie := func() {
        session = mkSession(user.ID, userBody.Email, expiresAt)
        sessionToken := uuid.New()
        session.ID = sessionToken
        // Sessions in DB
        initializers.DB.Create(&session)
        log.Printf("Session after creating new one in DB: %+v\n", session.ID)

        //check if already logged in with checkForCookie()

        http.SetCookie(c.Writer, &http.Cookie{
            Name: "session_token",
            Value: session.ID.String(),
            // User: session.UserID.String(),
            Expires: expiresAt,
        })
        // END NEW!
    }

    if userHasCookie.Error == nil && session.IsExpired() == false {
        // c.JSON(406, gin.H{
        //     "status": "You are already logged in!",
        //     "user": user,
        //     "token": session.ID,
        // })
        jsonFactory(c, 406, user, session.ID,"Hey there!", "You are already logged in!", "I'm hungry", "Me too!")
    } else if userHasCookie.Error == nil || session.IsExpired() == true {
        initializers.DB.Delete(&session, "user_id = ?", user.ID)
        mkCookie()
        log.Printf("Session after deleting cookie in DB: %+v\n", session.ID)
        // c.JSON(200, gin.H{
        //     "user": user,
        //     "token": session.ID,
        // })
        jsonFactory(c, 200, user, session.ID)
    } else if userHasCookie.Error != nil {
        mkCookie()
        log.Printf("No previous unexpired session: %+v\n", session.ID)
        // c.JSON(200, gin.H{
        //     // "message": "Thanks for joining!",
        //     "user": user,
        //     "token": session.ID,
        // })
        jsonFactory(c, 200, user, session.ID, "welcome", "Welcome " + user.Username + "!")
    }
}

func jsonFactory(c *gin.Context, code int, user models.User, token uuid.UUID, messages ...string) { // even messages (0 index) are the field name, odd are the actual message
    mess := new(jsonMessage)
    mess.m = make(map[string]string)
    var temp string
    for i, el := range messages {
        if i % 2 == 0 {
            temp = el
        } 
        if i % 2 != 0 {
            mess.m[temp] = el
        }
    }
    c.JSON(code, gin.H{
        "user": user,
        "token": token,
        "status code": code,
        "messages": mess.m,
    })
}

func Welcome(c *gin.Context) {
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
    }
    
    // SENDS 200 even when not logged in!!!
    if cookie.Email == "" {
        c.Status(400)
        return
    } else {
        c.JSON(200, gin.H{
            "welcome": cookie.Email,
        })
    }
}

func Refresh(c *gin.Context) {
    var oldSession models.Session
    var newSession models.Session
    user := getUser(c)
    
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
    }

    remove := initializers.DB.Delete(&oldSession, "user_id = ?", user.ID)
    if cookie.IsExpired() {
        if remove.Error != nil {
            log.Println(remove.Error)
        }
        c.Status(401)
        return
    }

    if cookie.IsExpired() == false {
        newSessionToken := uuid.New()
        expiresAt := time.Now().Add(120 * time.Second)
        newSession = mkSession(user.ID, user.Email, expiresAt)
        newSession.ID = newSessionToken
        // Sessions in DB
        initializers.DB.Create(&newSession)

        if remove.Error != nil {
            log.Println(remove.Error)
        }

        http.SetCookie(c.Writer, &http.Cookie{
            Name: "session_token",
            Value: newSessionToken.String(),
            Expires: expiresAt,
        })
    }
}

func Logout(c *gin.Context) {
    user := getUser(c)
    var session models.Session
    // cookie, err := checkForCookie(c)
    // if err != nil {
    //     log.Println(err)
    // }

    remove := initializers.DB.Delete(&session, "user_id = ?", user.ID)
    if remove.Error != nil {
        log.Println(remove.Error)
    }

    http.SetCookie(c.Writer, &http.Cookie{
        Name: "session_token",
        Value: "",
        Expires: time.Now(),
    })
}

func UpdateTodo(c *gin.Context) {
    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        // id := c.Param("id")
        var todoBody todoReq
        c.Bind(&todoBody)

        var todo models.Todo

        result := initializers.DB.Where("user_id = ?", user.ID).First(&todo, "id = ?", todoBody.ID)
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
}

func UpdateUser(c *gin.Context) {
    user := getUser(c)
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
        jsonFactory(c, 401, user, cookie.ID)
        return
    } else {
        var userBody userReq
        c.Bind(&userBody)
        pw := encrypt(userBody.NewPassword)

        match := checkPasswordHash(userBody.Password, user.Password)
        log.Println(match)

        if !match {
            c.Status(401)
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

func checkForCookie(c *gin.Context) (models.Session, error) {
    var session models.Session

    cookie, err := c.Request.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            c.Status(401)
            return session, errors.New("No session token.")
        }
        c.Status(400)
            return session, errors.New("Bad request.")
    }

    log.Printf("cookie: %v\n", cookie.Value)
    sessionToken := cookie.Value
    userSession := initializers.DB.First(&session, "id = ?", sessionToken)
    
    if userSession.Error != nil {
        c.Status(401)
        return session, errors.New("Unauthorized.")
    }
    
    if session.IsExpired() {
        initializers.DB.Delete(&session, "id = ?", session.ID)
        // maybe?
        session = models.Session{}
        c.Status(401)
        return session, errors.New("Session is expired.")
    }
    return session, nil
}

func getUser(c *gin.Context) models.User {
    // cookie, err := c.Request.Cookie("session_token")
    // if err != nil {
    //     if err == http.ErrNoCookie {
    //         c.Status(401)
    //     }
    //     c.Status(400)
    // }
    // var session models.Session
    // sessionToken := cookie.Value
    // userSession := initializers.DB.First(&session, "id = ?", sessionToken)
    // if userSession.Error != nil {
    //     c.Status(401)
    // }
    cookie, err := checkForCookie(c)
    if err != nil {
        log.Println(err)
    }

	var user models.User
	initializers.DB.First(&user, "email = ?", cookie.Email) // get rid of cookie.Email

    log.Printf("getUserId: %v\n", user.ID)
    return user
}


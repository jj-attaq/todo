package commands

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jj-attaq/todo/src/database"
	"github.com/jj-attaq/todo/utils"
)

var AllTodos []Entry

type Entry struct {
	Id       int            `json:"id"`
	UserID   sql.NullString `json:"userID"`
	UserPw   sql.NullString `json:"userPw"`
	Item     string         `json:"item"`
	Finished bool           `json:"finished"`
	UniqueID string         `json:"uuid"`
}

func contains(s []string, str string) bool {
	for _, el := range s {
		if el == str {
			return true
		}
	}
	return false
}

// change Input to work for api, so that a request can specify all fields when
// entering new entries outside terminal
func Input(prompt string) string {
	fmt.Printf("%s", prompt)
	stdinScanner := bufio.NewScanner(os.Stdin)
	stdinScanner.Scan()
	output := stdinScanner.Text()
	return output
}
func Commands() []string {
	var commands []string
	commands = append(commands, "add", "delete", "quit", "?", "show", "update")
	return commands
}
func ExecCommand(input string) string {
	var output string
	if contains(Commands(), input) {
		output = input
	} else {
		fmt.Printf("'%s' is not a valid command.\n", input)
	}
	return output
}
func AddTable() {
	removeTable, err := database.ConnDB().Exec("DROP TABLE IF EXISTS list")
	utils.HandleError(err)

	removeTable.RowsAffected()
	table, err := database.ConnDB().Exec("CREATE TABLE list (id MEDIUMINT NOT NULL AUTO_INCREMENT, userID VARCHAR(255) DEFAULT NULL, userPw VARCHAR(255) DEFAULT NULL, item VARCHAR(30), finished BOOL DEFAULT 0, PRIMARY KEY (id), uniqueID VARCHAR(255))")
	utils.HandleError(err)

	table.RowsAffected()
}
func AddTodo(todo string) {
	if len(todo) > 30 {
		fmt.Printf("Todo item is too long, must be below 30 characters.\n")
		return
	}
	if len(todo) <= 0 {
		fmt.Printf("Todo item is too short, must be at least 1 characters.\n")
		return
	}
	uniqueID := uuid.New()
	addTodo, err := database.ConnDB().Exec("INSERT INTO list (item, uniqueID) VALUES (?, ?)", todo, uniqueID)
	utils.HandleError(err)

	addTodo.RowsAffected()
}
func UpdateBool(todoId string) {
	var isFinished int
	var entry Entry
	status, err := database.ConnDB().Query("SELECT * FROM list WHERE uniqueID = ?", todoId)
	for status.Next() {
		err = status.Scan(&entry.Id, &entry.UserID, &entry.UserPw, &entry.Item, &entry.Finished, &entry.UniqueID)
		utils.HandleError(err)
		if entry.Finished == false {
			isFinished = 1
		} else if entry.Finished == true {
			isFinished = 0
		}
	}
	prep, err := database.ConnDB().Prepare("UPDATE list SET Finished = ? WHERE uniqueID = ?")
	utils.HandleError(err)

	defer prep.Close()
	update, err := prep.Exec(isFinished, todoId)

	utils.HandleError(err)
	update.RowsAffected()
}
func ShowList() {
	var entry Entry
	show, err := database.ConnDB().Query("SELECT * FROM list")
	for show.Next() {
		err = show.Scan(&entry.Id, &entry.UserID, &entry.UserPw, &entry.Item, &entry.Finished, &entry.UniqueID)
		utils.HandleError(err)

		res, err := json.MarshalIndent(entry, "", "    ")
		utils.HandleError(err)

		fmt.Printf("%+8v\n", string(res))
	}
}
func AddTask(c *gin.Context) {
	var newEntry Entry
	if err := c.BindJSON(&newEntry); err != nil {
		panic(err)
	}
	todo := newEntry.Item

	uniqueID := uuid.New()
	addTodo, err := database.ConnDB().Exec("INSERT INTO list (item, uniqueID) VALUES (?, ?)", todo, uniqueID)
	utils.HandleError(err)

	addTodo.RowsAffected()
}

func RemoveTodo(todoId string) {
	//	id, err := strconv.Atoi(todoId)
	remove, err := database.ConnDB().Exec("DELETE FROM list WHERE ? = uniqueID", todoId)
	utils.HandleError(err)

	remove.RowsAffected()
}

func UpdateTask(c *gin.Context) {
	uuid := c.Param("uuid")
	for _, el := range ShowJSON() {
		if el.UniqueID == uuid {
			UpdateBool(uuid)
		}
	}
}
func GetTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ShowJSON())
}
func GetOneTodo(c *gin.Context) {
	uuid := c.Param("uuid")
	for _, el := range ShowJSON() {
		if el.UniqueID == uuid {
			c.IndentedJSON(http.StatusOK, el)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "todo not found"})
}
func RemoveTask(c *gin.Context) {
	uuid := c.Param("uuid")
	for _, el := range ShowJSON() {
		if el.UniqueID == uuid {
			RemoveTodo(uuid)
		}
	}
}
func ShowJSON() []Entry {
	var entry Entry
	show, err := database.ConnDB().Query("SELECT * FROM list")
	result := AllTodos
	for show.Next() {
		err = show.Scan(&entry.Id, &entry.UserID, &entry.UserPw, &entry.Item, &entry.Finished, &entry.UniqueID)
		utils.HandleError(err)

		utils.HandleError(err)
		result = append(result, entry)
	}

	return result
}

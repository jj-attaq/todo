package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jj-attaq/todo/src/database"
	"github.com/jj-attaq/todo/utils"
)

//var commands = []string{"add", "delete", "quit", "?", "show", "update"}

type Entry struct {
	Id       int
	Item     string
	Finished bool
	UniqueID string
}

func contains(s []string, str string) bool {
	for _, el := range s {
		if el == str {
			return true
		}
	}
	return false
}
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
	//	if contains(commands, input) {
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
	table, err := database.ConnDB().Exec("CREATE TABLE list (id MEDIUMINT NOT NULL AUTO_INCREMENT, item VARCHAR(30), finished BOOL DEFAULT 0, PRIMARY KEY (id), uniqueID VARCHAR(255))")
	utils.HandleError(err)

	table.RowsAffected()
}
func AddTodo(todo string) {
	if len(todo) > 30 {
		fmt.Printf("Todo item is too long, must be below 30 characters.\n")
		return
	}
	uniqueID := uuid.New()
	addTodo, err := database.ConnDB().Exec("INSERT INTO list (item, uniqueID) VALUES (?, ?)", todo, uniqueID)
	utils.HandleError(err)

	addTodo.RowsAffected()
}
func RemoveTodo(todoId string) {
	id, err := strconv.Atoi(todoId)
	remove, err := database.ConnDB().Exec("DELETE FROM list WHERE ? = Id", id)
	utils.HandleError(err)

	remove.RowsAffected()
}
func UpdateBool(status string, todoId string) {
	var isFinished int
	if status == "y" {
		isFinished = 1
	} else if status == "n" {
		isFinished = 0
	}
	id, err := strconv.Atoi(todoId)
	prep, err := database.ConnDB().Prepare("UPDATE list SET Finished = ? WHERE Id = ?")
	utils.HandleError(err)

	defer prep.Close()
	update, err := prep.Exec(isFinished, id)
	utils.HandleError(err)

	update.RowsAffected()
}
func ShowList() {
	var entry Entry
	show, err := database.ConnDB().Query("SELECT * FROM list")
	for show.Next() {
		err = show.Scan(&entry.Id, &entry.Item, &entry.Finished, &entry.UniqueID)
		utils.HandleError(err)

		res, err := json.MarshalIndent(entry, "", "    ")
		utils.HandleError(err)

		fmt.Printf("%+8v\n", string(res))
	}
}

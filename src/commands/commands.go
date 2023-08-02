package commands

import (
    "fmt"
    "bufio"
    "os"
    "strconv"
    "github.com/jj-attaq/todo/utils"
    "github.com/jj-attaq/todo/src/database"
)

type Entry struct {
    id int
    item string
    finished bool
}

func contains(s []string, str string) bool {
    for _, el := range s {
        if el == str {
            return true
        }
    }
    return false
}
func commands() []string {
    var commands []string
    commands = append(commands, "add", "delete", "quit", "?", "show", "update")
    return commands
}
func input(prompt string) string {
    fmt.Printf("%s", prompt)
    stdinScanner := bufio.NewScanner(os.Stdin)
    stdinScanner.Scan()
    output := stdinScanner.Text()
    return output
}
func execCommand(input string) string {
    var output string
    if contains(commands(), input) {
        output = input
    } else {
        fmt.Printf("'%s' is not a valid command.\n", input)
    }
    return output
}
func addTable() {
    removeTable, err := database.connDB().Exec("DROP TABLE IF EXISTS list")
    utils.HandleError(err)

    removeTable.RowsAffected()
    table, err := database.connDB().Exec("CREATE TABLE list (id MEDIUMINT NOT NULL AUTO_INCREMENT, item VARCHAR(30), finished BOOL DEFAULT 0, PRIMARY KEY (id))")
    utils.HandleError(err)

    table.RowsAffected()
}
func addTodo(todo string) {
    if len(todo) > 30 {
        fmt.Printf("Todo item is too long, must be below 30 characters.\n")
        return 
    }
    addTodo, err := database.connDB().Exec("INSERT INTO list (item) VALUES (?)", todo)
    utils.HandleError(err)

    addTodo.RowsAffected()
}
func removeTodo(todoId string) {
    id, err := strconv.Atoi(todoId)
    remove, err := database.connDB().Exec("DELETE FROM list WHERE ? = id", id)
    utils.HandleError(err)

    remove.RowsAffected()
}
func updateBool(status string, todoId string) {
    var isFinished int
    if status == "y" {
        isFinished = 1
    } else if status == "n" {
        isFinished = 0
    }
    id, err := strconv.Atoi(todoId)
    prep, err := database.connDB().Prepare("UPDATE list SET finished = ? WHERE id = ?")
    utils.HandleError(err)

    defer prep.Close()
    update, err := prep.Exec(isFinished, id)
    utils.HandleError(err)

    update.RowsAffected()
}
func showList() {
    var entry Entry
    show, err := database.connDB().Query("SELECT * FROM list")
    for show.Next() {
        err = show.Scan(&entry.id, &entry.item, &entry.finished)
        utils.HandleError(err)

        fmt.Printf("%+15v\n", entry)
    }
}

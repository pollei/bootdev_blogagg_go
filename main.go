package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/pollei/bootdev_blogagg_go/internal/config"
	"github.com/pollei/bootdev_blogagg_go/internal/database"
)

type mainEvilGlobals struct {
	conf      *config.Config
	cmds      *cliCommands
	db        *sql.DB
	dbQueries *database.Queries
}

var mainGLOBS mainEvilGlobals

func (g *mainEvilGlobals) init() error {
	g.conf = config.Read()
	g.cmds = &cliCommands{}
	g.cmds.cmdList = make(map[string]cliCommand)
	g.cmds.register(
		cliCommand{
			name:        "exit",
			description: "Exit the Gator",
			callback:    commandExit,
		},
		cliCommand{
			name:        "login",
			description: "login the Gator",
			callback:    commandLogin,
		},
		cliCommand{
			name:        "register",
			description: "register new user into the Gator",
			callback:    commandRegister,
		},
		cliCommand{
			name:        "reset",
			description: "deletes all the users in the Gator",
			callback:    commandReset,
		},
		cliCommand{
			name:        "users",
			description: "list users in the Gator",
			callback:    commandUsers,
		},
		cliCommand{
			name:        "agg",
			description: "aggre with the Gator",
			callback:    commandNotImplementedYet,
		},
		cliCommand{
			name: "help", description: "Displays a help message",
			callback: commandHelp},
	)
	//var err error
	db, err := sql.Open("postgres", g.conf.Db_url)
	if err != nil {
		return err
	}
	g.db = db
	g.dbQueries = database.New(db)

	return nil
}
func main() {
	args := os.Args
	if len(args) < 2 {
		os.Exit(1)
	}
	err := mainGLOBS.init()
	if err != nil {
		fmt.Printf("init failed because of %v \n", err)
		os.Exit(1)
	}

	err = mainGLOBS.cmds.run(args[1:])
	if err != nil {
		os.Exit(1)
	}

	/* fmt.Println("start")
	conf := config.Read()
	conf.SetUser("sjp")

	conf2 := config.Read()

	fmt.Printf("user is %s\n", conf2.User_name)
	fmt.Printf("url is %s\n", conf2.Db_url) */

}

func dirtyNameRune(r rune) bool {
	if r >= '0' && r <= '9' {
		return false
	}
	if r >= 'a' && r <= 'z' {
		return false
	}
	if r >= 'A' && r <= 'Z' {
		return false
	}
	if r == '-' {
		return false
	}
	return true
}

func dirtyName(text string) bool {
	return strings.ContainsFunc(text, dirtyNameRune)
}

func cleanInput(text string) []string {
	lowStr := strings.ToLower(strings.TrimSpace(text))
	return strings.Fields(lowStr)
}

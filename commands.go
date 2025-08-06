package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	// "strings"
	"github.com/google/uuid"
	"github.com/pollei/bootdev_blogagg_go/internal/database"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

type cliCommands struct {
	cmdList map[string]cliCommand
}

func (c *cliCommands) run(args []string) error {
	//cmd, ok := mainGLOBS.cmds.cmdList[args[0]]
	cmd, ok := c.cmdList[args[0]]
	if ok {
		return cmd.callback(args)
	}
	return errors.New("unknown command")

}

func (c *cliCommands) register(cmds ...cliCommand) error {
	for _, cmd := range cmds {
		c.cmdList[cmd.name] = cmd
	}
	return nil
}

func commandExit([]string) error {
	fmt.Println("Closing the Gator... Goodbye!")
	os.Exit(0)
	return nil
}

func commandLogin(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if dirtyName(args[1]) {
		fmt.Println("not a proper user name")
		return errors.New("dirty name")
	}
	bgCtx := context.Background()
	usr, err := mainGLOBS.dbQueries.GetUserByName(bgCtx, args[1])
	if err != nil {
		return err
	}
	fmt.Printf("got usr %v\n", usr)
	mainGLOBS.conf.SetUser(args[1])
	fmt.Printf("logged in as %s to Gator... \n", args[1])
	return nil
}
func commandRegister(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if dirtyName(args[1]) {
		fmt.Println("not a proper user name")
		return errors.New("dirty name")
	}
	bgCtx := context.Background()
	oldUsr, err := mainGLOBS.dbQueries.GetUserByName(bgCtx, args[1])
	if err == nil {
		return errors.New("user with that name already exists")
	}
	fmt.Printf("got old usr %v %v \n", oldUsr, err)
	usrParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: args[1]}
	usr, err := mainGLOBS.dbQueries.CreateUser(bgCtx, usrParams)
	if err != nil {
		return err
	}
	fmt.Printf("created user %s -> %v \n", args[1], usr)
	mainGLOBS.conf.SetUser(args[1])
	return nil
}

func commandReset([]string) error {
	bgCtx := context.Background()
	err := mainGLOBS.dbQueries.DeleteAllUsers(bgCtx)
	return err
}

func commandUsers([]string) error {
	bgCtx := context.Background()
	users, err := mainGLOBS.dbQueries.GetUsers(bgCtx)
	for _, usr := range users {
		var suffix string
		if usr.Name == mainGLOBS.conf.User_name {
			suffix = " (current)"
		}
		fmt.Printf("* %s%s\n", usr.Name, suffix)

	}
	return err
}

func commandNotImplementedYet([]string) error {
	fmt.Println("command Not Implemented Yet")
	os.Exit(1)
	return nil
}

func commandHelp([]string) error {
	fmt.Print("Welcome to the Gator!\nUsage: \n\n")
	for _, cmdItm := range mainGLOBS.cmds.cmdList {
		fmt.Printf("%s: %s\n", cmdItm.name, cmdItm.description)
	}
	return nil
}

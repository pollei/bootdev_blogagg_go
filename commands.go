package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	//"html"

	// "strings"
	"github.com/google/uuid"
	"github.com/pollei/bootdev_blogagg_go/internal/database"
)

type cmdCallback func([]string) error

type cliCommand struct {
	name        string
	description string
	//callback    func([]string) error
	callback cmdCallback
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

func middlewareLoggedIn(hand cmdCallback) cmdCallback {
	return func(args []string) error {
		if mainGLOBS.currUser == nil {
			return errors.New("must be logged in")
		}
		return hand(args)
	}
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
	_, err := mainGLOBS.dbQueries.GetUserByName(bgCtx, args[1])
	if err == nil {
		return errors.New("user with that name already exists ")
	}
	//fmt.Printf("got old usr %v %v \n", oldUsr, err)
	now := time.Now().UTC()
	usrParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now, UpdatedAt: now, Name: args[1]}
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

func commandAgg(args []string) error {
	var durStr string
	if len(args) <= 1 {
		durStr = "2h"
	} else {
		durStr = args[1]
	}
	timeBetweenRequests, err := time.ParseDuration(durStr)
	if err != nil {
		fmt.Printf("duration was not understood: %v \n", err)
		return err
	}
	ac := aggController{interval: timeBetweenRequests}
	mainGLOBS.aggControl = &ac
	mainGLOBS.aggControl.tick = time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v \n", timeBetweenRequests)

	for ; ; <-mainGLOBS.aggControl.tick.C {
		scrapeFeeds()
	}

	/* bgCtx := context.Background()
	feed, err := fetchFeed(bgCtx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	// https://pkg.go.dev/html#UnescapeString
	fmt.Printf("%v \n", feed)
	inter, err := time.ParseDuration("1h")
	if err != nil {
		mainGLOBS.agg.interval = inter
	} */

	//return nil
}

func commandAddFeed(args []string) error {
	if len(args) <= 2 {
		return errors.New("not enough arguments")
	}
	if mainGLOBS.currUser == nil {
		return errors.New("must be logged in")
	}
	bgCtx := context.Background()
	now := time.Now().UTC()
	feedParams := database.CreateFeedParams{
		ID: uuid.New(), CreatedAt: now, UpdatedAt: now,
		Name: args[1], Url: args[2], UserID: mainGLOBS.currUser.ID}
	feed, err := mainGLOBS.dbQueries.CreateFeed(bgCtx, feedParams)
	if err != nil {
		return err
	}
	fmt.Printf("created feed \"%s\" \"%s\" -> %v \n", args[1], args[1], feed)
	ffParams := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: now, UpdatedAt: now,
		FeedID: feed.ID, UserID: mainGLOBS.currUser.ID,
	}
	ff, err := mainGLOBS.dbQueries.CreateFeedFollow(bgCtx, ffParams)
	if err != nil {
		fmt.Printf("feed follow err %v \n ", err)
		return err
	}
	fmt.Printf(" %s is now following %s \n", ff.UserName, ff.FeedName)
	return nil
}
func commandFeeds([]string) error {
	bgCtx := context.Background()
	feeds, err := mainGLOBS.dbQueries.GetFeedsSummary(bgCtx)
	if err != nil {
		fmt.Printf("feed err %v \n ", err)
		return err
	}
	fmt.Println("FEEDS: ")
	//fmt.Printf(" %d ", len(feeds))
	for _, fd := range feeds {
		fmt.Printf(" %s %s %s -> %v \n", fd.Name, fd.Url, fd.UserName, fd)
	}
	return nil
}
func commandFollow(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if mainGLOBS.currUser == nil {
		return errors.New("must be logged in")
	}
	bgCtx := context.Background()
	feed, err := mainGLOBS.dbQueries.GetFeedByUrl(bgCtx, args[1])
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	ffParams := database.CreateFeedFollowParams{
		ID: uuid.New(), CreatedAt: now, UpdatedAt: now,
		FeedID: feed.ID, UserID: mainGLOBS.currUser.ID,
	}
	ff, err := mainGLOBS.dbQueries.CreateFeedFollow(bgCtx, ffParams)
	if err != nil {
		fmt.Printf("feed follow err %v \n ", err)
		return err
	}
	fmt.Printf(" %s is now following %s \n", ff.UserName, ff.FeedName)
	return nil
}

func commandUnfollow(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if mainGLOBS.currUser == nil {
		return errors.New("must be logged in")
	}
	bgCtx := context.Background()
	feed, err := mainGLOBS.dbQueries.GetFeedByUrl(bgCtx, args[1])
	if err != nil {
		return err
	}
	delParams := database.DeleteFeedFollowByFeedIdParams{
		UserID: mainGLOBS.currUser.ID, FeedID: feed.ID}
	err = mainGLOBS.dbQueries.DeleteFeedFollowByFeedId(bgCtx, delParams)
	if err != nil {
		return err
	}
	return nil
}

func commandFollowing([]string) error {
	if mainGLOBS.currUser == nil {
		return errors.New("must be logged in")
	}
	bgCtx := context.Background()
	itFollows, err := mainGLOBS.dbQueries.GetFeedFollowsByUserId(bgCtx, mainGLOBS.currUser.ID)
	if err != nil {
		fmt.Printf("following err %v \n ", err)
		return err
	}
	for _, follow := range itFollows {
		fmt.Printf("%s \n", follow.FeedName)
	}
	return nil
}

func commandBrowse(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	fmt.Println("command Not Implemented Yet")
	return nil
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

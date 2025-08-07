package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pollei/bootdev_blogagg_go/internal/database"
)

type aggController struct {
	tick     *time.Ticker
	interval time.Duration
	done     chan bool
	mutx     sync.Mutex
}

func scrapeFeeds() {
	bgCtx := context.Background()
	feed, err := mainGLOBS.dbQueries.GetNextFeedToFetch(bgCtx)
	if err != nil {
		return
	}
	fetchParam := database.MarkFeedFetchedParams{
		ID: feed.ID, UpdatedAt: time.Now()}
	mainGLOBS.dbQueries.MarkFeedFetched(bgCtx, fetchParam)
	feedBuf, err := fetchFeed(bgCtx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return
	}
	// https://pkg.go.dev/html#UnescapeString
	fmt.Printf("%s %v \n", feedBuf.Channel.Title, feedBuf.Channel.Item)
	for _, itm := range feedBuf.Channel.Item {
		fmt.Printf("%s \n", itm.Title)
	}
}

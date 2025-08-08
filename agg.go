package main

import (
	"context"
	"database/sql"
	"fmt"

	//"sync"
	"time"

	"github.com/google/uuid"

	"github.com/pollei/bootdev_blogagg_go/internal/database"
)

type aggController struct {
	tick     *time.Ticker
	interval time.Duration
	//done     chan bool
	//mutx     sync.Mutex
}

func scrapeFeeds() {
	bgCtx := context.Background()
	fmt.Println("starting scrap ")
	feed, err := mainGLOBS.dbQueries.GetNextFeedToFetch(bgCtx)
	if err != nil {
		return
	}
	now := time.Now().UTC()
	if feed.LastFetchedAt.Valid {
		expireT := feed.LastFetchedAt.Time.Add(5 * time.Minute)
		fmt.Printf("last fetch is %v \n", feed.LastFetchedAt)
		fmt.Printf("expire is  %v \n", expireT)

		fmt.Printf("utc now is  %v \n", now)

		if expireT.After(now) {
			return
		}
		fmt.Println("feed expired")
	} else {
		fmt.Println("no last fetch ")
	}
	fetchParam := database.MarkFeedFetchedParams{
		ID: feed.ID, UpdatedAt: now}
	mainGLOBS.dbQueries.MarkFeedFetched(bgCtx, fetchParam)
	feedBuf, err := fetchFeed(bgCtx, feed.Url)
	if err != nil {
		return
	}
	// https://pkg.go.dev/html#UnescapeString
	fmt.Printf("feed channel title %s \n", feedBuf.Channel.Title)
	for _, itm := range feedBuf.Channel.Item {
		if len(itm.Title) <= 0 {
			continue
		}
		fmt.Printf("item title %s pubDate %s \n", itm.Title, itm.PubDate)
		pubT1123, err := time.Parse(time.RFC1123, itm.PubDate)
		if err == nil {
			//fmt.Printf("rfc3339 is %v \n", pubT1123)
			// sql.NullString( itm.Title)
			sqlTitle := sql.NullString{String: itm.Title, Valid: true}
			sqlDesc := sql.NullString{String: itm.Description}
			sqlDesc.Valid = len(itm.Description) > 0
			postParam := database.CreatePostParams{
				ID: uuid.New(), CreatedAt: now, UpdatedAt: now,
				FeedID: feed.ID, Title: sqlTitle, PublishedAt: pubT1123,
				Url: itm.Link, Description: sqlDesc,
			}
			mainGLOBS.dbQueries.CreatePost(bgCtx, postParam)
		}
	}
}

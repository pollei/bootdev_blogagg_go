package main

import (
	"context"
	"encoding/xml"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var retFeed RSSFeed
	buf, err := getBytesfromUrl(ctx, feedURL)
	if err != nil {
		return nil, err
	}
	xml.Unmarshal(buf, &retFeed)
	return &retFeed, err
}

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error Creating Request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error Requesting URL: %v", err)
	}
	defer res.Body.Close()

	var rssFeed RSSFeed
	decoder := xml.NewDecoder(res.Body)
	if err = decoder.Decode(&rssFeed); err != nil {
		return nil, fmt.Errorf("Error Decoding Data: %v", err)
	}

	return &rssFeed, nil
}

func handlerAgg(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Printf("%s (%s)\n", html.UnescapeString(rssFeed.Channel.Title), html.UnescapeString(rssFeed.Channel.Link))
	fmt.Printf("%s\n\n", html.UnescapeString(rssFeed.Channel.Description))

	for _, item := range(rssFeed.Channel.Item) {
		fmt.Printf("%s (%s) - %s\n", html.UnescapeString(item.Title), html.UnescapeString(item.Link), html.UnescapeString(item.PubDate))
		fmt.Printf("%s\n\n", html.UnescapeString(item.Description))
	}

	return nil
}

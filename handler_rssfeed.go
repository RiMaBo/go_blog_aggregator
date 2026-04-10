package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"time"

	"internal/database"

	"github.com/google/uuid"
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error Finding User: %v", err)
	}

	feedname := cmd.Args[0]
	feedUrl := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedname,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Creating Feed: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error Following Feed: %v", err)
	}

	fmt.Println("Feed Created Successfully:")
	fmt.Printf(" - ID:      %v\n", feed.ID)
	fmt.Printf(" - Name:    %v\n", feed.Name)
	fmt.Printf(" - URL:     %v\n", feed.Url)
	fmt.Printf(" - User ID: %v\n", feed.UserID)
	fmt.Println()
	fmt.Printf("User %s Now Following Feed '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error Listing Feeds: %v", err)
	}

	if len(feeds) < 1 {
		fmt.Println("No Feeds Found.")
		return nil
	}

	fmt.Printf("Found %d Feed(s):\n", len(feeds))
	for _, feed := range(feeds) {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("User Not Found: %v", err)
		}

		fmt.Printf("%s: \n", feed.Name)
		fmt.Printf(" - URL:        %s\n", feed.Url)
		fmt.Printf(" - Created By: %s\n\n", user.Name)
	}

	return nil
}

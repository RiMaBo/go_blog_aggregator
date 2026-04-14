package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
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

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error Finding Next Feed to Fetch: %v", err)
	}

	if err := s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("Error Marking Feed %s as Fetched: %v", feed.Name, err)
	}

	rssfeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Error Fetching Feed %s: %v", feed.Name, err)
	}

	for _, rssfeedItem := range rssfeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 Z0700", rssfeedItem.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		}


		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: html.UnescapeString(rssfeedItem.Title), Valid: true},
			Url:         rssfeedItem.Link,
			Description: sql.NullString{String: html.UnescapeString(rssfeedItem.Description), Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Unable to Create Post: %v", err)
			continue
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage: %s <time between requests>", cmd.Name)
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error Parsing Time Between Requests")
	}

	fmt.Printf("Collecting Feeds Every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
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

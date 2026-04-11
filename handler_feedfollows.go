package main

import (
	"context"
	"fmt"
	"time"

	"internal/database"

	"github.com/google/uuid"
)


func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}

	feedUrl := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Error Finding Feed %s: %v", feedUrl, err)
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

	fmt.Printf("User %s Now Following Feed '%s'\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollowsForUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error Finding Feed Follows for User %s: %v", user.Name, err)
	}

	if len(feedFollowsForUser) < 1 {
		fmt.Printf("User %s Not Following Any Feeds.\n", user.Name)
	}

	fmt.Printf("User %s Following %d Feed(s):\n", user.Name, len(feedFollowsForUser))
	for _, feedFollow := range feedFollowsForUser {
		fmt.Printf(" - %s\n", feedFollow.FeedName)
	}

	return nil
}

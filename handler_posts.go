package main

import (
	"context"
	"fmt"
	"strconv"

	"internal/database"
)

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.Args) > 0 {
		if givenLimit, err := strconv.Atoi(cmd.Args[0]); err != nil {
			return fmt.Errorf("Please provide a numeric value for limit.")
		} else {
			limit = givenLimit
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		fmt.Errorf("Error Getting Posts: %v", err)
	}

	for _, post := range(posts) {
		fmt.Printf("Feed:         %s\n", post.Name)
		fmt.Printf("Title:        %s\n", post.Title.String)
		fmt.Printf("Link:         %s\n", post.Url)
		fmt.Printf("Published At: %v\n", post.PublishedAt.Time)
		fmt.Printf("Description:  %s\n", post.Description.String)
		fmt.Println("")
	}

	return nil
}

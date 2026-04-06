package main

import (
	"context"
	"fmt"
	"time"

	"internal/database"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		return fmt.Errorf("Couldn't find user: %v", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("Couldn't set current user: %v", err)
	}

	fmt.Println("User switched successfully.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage %s <name>", cmd.Name)
	}

	username := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("Couldn't Create User: %v", err)
	}

	fmt.Println("User Created Successfully:")
	fmt.Printf(" - ID:    %v\n", user.ID)
	fmt.Printf(" - Name:  %v\n", user.Name)

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %v", err)
	}

	fmt.Println("User switched successfully.")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("Failed to reset users table: %v", err)
	}

	fmt.Println("Users table reset successfully.")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't list users: %v", err)
	}

	for _, usr := range users {
		if usr.Name == s.cfg.CurrentUserName {
			fmt.Printf(" - %s (current)\n", usr.Name)
		} else {
			fmt.Printf(" - %s\n", usr.Name)
		}
	}

	return nil
}

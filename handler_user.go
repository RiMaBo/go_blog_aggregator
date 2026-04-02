package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}

	username := cmd.Args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Println("The user has been set")
	return nil
}

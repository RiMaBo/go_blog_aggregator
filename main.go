package main

import (
	"internal/config"

	"log"
	"os"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error Reading Config File: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	err = cmds.run(programState, command{Name: cmd, Args: args})
	if err != nil {
		log.Fatal(err)
	}
}

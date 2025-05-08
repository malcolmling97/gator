package main

import (
	"fmt"
	"gator/internal/config"
	"log"
	"os"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := &state{
		config: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	// Check if enough arguments were provided
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Not enough arguments. Please provide a command.")
		os.Exit(1)
	}

	// Create a command struct
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Run the command
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}

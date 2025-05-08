package main

import (
	"fmt"
	"gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

// func signature of all command handlers
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("%s handler requires username", cmd.name)
	}

	username := cmd.args[0]

	// Set the user in the config
	err := s.config.SetUser(username)
	if err != nil {
		return err
	}

	// Print success message
	fmt.Printf("User set to: %s\n", username)

	return nil

}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)

}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f

}

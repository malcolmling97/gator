package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"
	"time"

	"github.com/google/uuid"
)

type state struct {
	db     *database.Queries
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

	// Step 1: Check if the user already exists
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		// User already exists — exit directly
		fmt.Fprintf(os.Stderr, "user with name '%s' already exists\n", username)
		os.Exit(1)
	}

	// Set the user in the config
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}

	// Print success message
	fmt.Printf("User set to: %s\n", username)

	return nil

}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("%s handler requires username", cmd.name)
	}

	username := cmd.args[0]

	// Step 1: Check if the user already exists
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		// User already exists — exit directly
		fmt.Fprintf(os.Stderr, "user with name '%s' already exists\n", username)
		os.Exit(1)
	}
	if err != sql.ErrNoRows {
		// Unexpected error
		return fmt.Errorf("error checking for user: %w", err)
	}

	// Step 2: User doesn't exist, create them
	userArgs := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	_, err = s.db.CreateUser(context.Background(), userArgs)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set the user in the config
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully!")
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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to delete users:", err)
		os.Exit(1) // explicitly required
	}
	fmt.Println("All users deleted successfully.")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	fmt.Println("Users:")
	for _, u := range users {
		if s.config.CurrentUserName == u.Name {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}

	}

	return nil
}

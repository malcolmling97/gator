package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: follow <feed-url>")
	}

	feedURL := cmd.args[0]

	// Get current username from config
	currentUsername := s.config.CurrentUserName
	if currentUsername == "" {
		return fmt.Errorf("no user is currently logged in")
	}

	// Look up user in DB
	user, err := s.db.GetUser(context.Background(), currentUsername)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Look up feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	// Create feed follow record
	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	result, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("%s followed %s\n", result.UserName, result.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	// Get current username from config
	currentUsername := s.config.CurrentUserName
	if currentUsername == "" {
		return fmt.Errorf("no user is currently logged in")
	}

	// Look up user in DB
	user, err := s.db.GetUser(context.Background(), currentUsername)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Get all feed follows for user
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", currentUsername)
	for _, f := range follows {
		fmt.Printf("- %s\n", f.FeedName)
	}

	return nil
}

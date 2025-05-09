package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	// HARD CODED FOR NOW, REMEMBER TO CHANGE WHEN WORKING
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	fmt.Println("Items:")
	for _, item := range feed.Channel.Items {
		fmt.Printf("- %s (%s)\n", item.Title, item.PubDate)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	// 1. Get current user from config
	currentUsername := s.config.CurrentUserName
	if currentUsername == "" {
		return fmt.Errorf("no user is currently logged in")
	}

	// 2. Get user record from DB
	user, err := s.db.GetUser(context.Background(), currentUsername)
	if err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}

	// 3. Create feed connected to that user
	feedArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedArgs)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	// 4. Print out the new feed fields
	fmt.Println("Feed successfully created:")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)
	fmt.Printf("Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated At: %s\n", feed.UpdatedAt.Format(time.RFC3339))

	return nil
}

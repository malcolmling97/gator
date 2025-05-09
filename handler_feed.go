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

	fmt.Printf("Feed created: %s (%s)\n", feed.Name, feed.Url)

	// 4. Automatically follow the feed for the user
	followArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	result, err := s.db.CreateFeedFollow(context.Background(), followArgs)
	if err != nil {
		return fmt.Errorf("failed to follow feed after creation: %w", err)
	}

	// 5. Confirm both actions
	fmt.Printf("%s is now following %s\n", result.UserName, result.FeedName)
	return nil
}

func handlerListAllFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeedsWithUsernames(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("Feeds:")
	for _, f := range feeds {
		fmt.Printf("- %s (%s) [added by: %s]\n", f.FeedName, f.FeedUrl, f.Username)
	}

	return nil
}

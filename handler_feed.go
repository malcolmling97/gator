package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

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

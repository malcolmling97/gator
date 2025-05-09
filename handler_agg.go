package main

import (
	"context"
	"fmt"
	"gator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	fmt.Println("Items:")
	for _, item := range feed.Channel.Items {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
	}

	return nil
}

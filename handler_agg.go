package main

import (
	"context"
	"fmt"
	"gator/internal/rss"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: agg <time_between_requests>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration string: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for {
		scrapeFeeds(s)
		<-ticker.C
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Printf("error getting next feed: %v", err)
		return
	}

	log.Printf("Fetching feed: %s (%s)", feed.Name, feed.Url)

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("error fetching feed from %s: %v", feed.Url, err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		log.Printf("Post: %s (%s)", item.Title, item.PubDate)
	}

	err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("error marking feed as fetched: %v", err)
	}
}

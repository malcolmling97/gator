package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		pubTime, _ := parseTime(item.PubDate) // helper to parse dates safely

		err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  pubTime,
				Valid: !pubTime.IsZero(),
			},
			FeedID: feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value") {
				continue // silently skip duplicates
			}
			log.Printf("error saving post: %v", err)
		}
	}

	err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("error marking feed as fetched: %v", err)
	}
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		t, err = time.Parse(time.RFC1123, s)
	}
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
	}
	return t, err
}

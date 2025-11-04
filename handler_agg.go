package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jkk290/gator/internal/database"
	"github.com/jkk290/gator/internal/feedapi"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error need duration argument")
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed: %w", err)
	}
	feedId := feedToFetch.ID
	updatedFeed, err := s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feedId,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}
	url := updatedFeed.Url

	fmt.Printf("fetching from %s...\n", url)
	RSSFeed, err := feedapi.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	for _, item := range RSSFeed.Channel.Item {
		fmt.Printf("Title: %s\n", item.Title)
	}
	return nil
}

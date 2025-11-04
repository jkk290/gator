package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		pubDate, err := timeParser(item.PubDate)
		if err != nil {
			return fmt.Errorf("error converting published Date to time.Time: %w", err)
		}
		_, error := s.db.AddPost(context.Background(), database.AddPostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feedId,
		})
		if error != nil {
			return fmt.Errorf("error saving post to database: %w", error)
		}
	}
	return nil
}

func timeParser(timeStamp string) (time.Time, error) {
	layout := "Nov 4, 2025 at 6:00am (HST)"
	return time.Parse(layout, timeStamp)
}

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
	for _, post := range RSSFeed.Channel.Item {
		pubDate, err := timeParser(post.PubDate)
		if err != nil {
			return fmt.Errorf("error converting published Date to time.Time: %w", err)
		}
		savedPost, err := s.db.AddPost(context.Background(), database.AddPostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       post.Title,
			Url:         post.Link,
			Description: sql.NullString{String: post.Description, Valid: true},
			PublishedAt: pubDate,
			FeedID:      feedId,
		})
		if err != nil {
			return fmt.Errorf("error saving post to database: %w", err)
		}
		fmt.Printf("successfully saved post to db, %s\n", savedPost.Title)
	}
	return nil
}

func timeParser(timeStamp string) (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	converted, err := time.Parse(layout, timeStamp)
	if err != nil {
		return time.Time{}, err
	}
	return converted, nil
}

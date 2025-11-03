package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("not enough arguments, need name and url")
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}
	userId := currentUser.ID
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	newFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    userId,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %w", err)
	}

	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
		FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Printf("feed successfully added: %v\n", newFeed)
	fmt.Printf("feed successfully followed: %s", newFeedFollow.FeedName)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all feeds")
	}
	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", feed.UserName)
	}
	return nil
}

func handlerAddFeedFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error not enough arguments, need url")
	}
	url := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user info")
	}
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed info")
	}
	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new feed follow: %w", err)
	}
	fmt.Printf("Following feed: %v", newFeedFollow)
	return nil
}

func handlerGetUserFeeds(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user info")
	}

	feeds, err := s.db.GetUserFeeds(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting user's feeds: %w", err)
	}
	fmt.Printf("%s's Feeds\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("Name: %s\n", feed.FeedName)
	}
	return nil
}

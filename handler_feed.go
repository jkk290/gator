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
	fmt.Printf("feed successfully added: %v\n", newFeed)
	return nil

}

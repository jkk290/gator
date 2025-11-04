package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jkk290/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.Args) == 1 {
		converted, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("error converting limit to int: %w", err)
		}
		limit = int32(converted)
	} else {
		limit = 2
	}
	userId := user.ID

	feeds, err := s.db.GetUserFeeds(context.Background(), userId)
	if err != nil {
		return fmt.Errorf("error getting user's feeds: %w", err)
	}

	feedIds := []uuid.UUID{}
	for _, feed := range feeds {
		feedIds = append(feedIds, feed.FeedID)
	}

	posts, err := s.db.GetPostsByFeeds(context.Background(), database.GetPostsByFeedsParams{
		Column1: feedIds,
		Limit:   limit,
	})
	if err != nil {
		return fmt.Errorf("error getting posts from feeds: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Published At: %v\n\n", post.PublishedAt)
	}
	return nil
}

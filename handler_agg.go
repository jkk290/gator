package main

import (
	"context"
	"fmt"

	"github.com/jkk290/gator/internal/feedapi"
)

func handlerAgg(s *state, cmd command) error {
	// url := cmd.Args[0]
	url := "https://www.wagslane.dev/index.xml"
	newFeed, err := feedapi.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	fmt.Printf("RSS Feed: %v\n", newFeed)
	return nil
}

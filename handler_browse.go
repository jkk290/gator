package main

import (
	"fmt"
	"strconv"
)

func handlerBrowse(s *state, cmd command) error {
	var limit int32
	if len(cmd.Args) == 1 {
		converted, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("error converting limit to int: %w", err)
		}
		limit = int32(converted)
	}
	posts, err := s.db.GetPostsByFeed()
}

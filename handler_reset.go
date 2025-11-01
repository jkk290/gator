package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("error resetting user table: %w", err)
	}
	fmt.Println("successfully reset user table")
	return nil
}

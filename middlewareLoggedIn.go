package main

import (
	"context"

	"github.com/jkk290/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		userName := s.cfg.CurrentUserName
		loggedInUser, err := s.db.GetUser(context.Background(), userName)
		if err != nil {
			return err
		}

		return handler(s, cmd, loggedInUser)
	}
}

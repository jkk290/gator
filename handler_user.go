package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jkk290/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no username provided")
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user does not exist")
	}

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("error setting current user: %w", err)
	}

	fmt.Printf("user %s successfully set\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no username provided")
	}
	name := cmd.Args[0]

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("error adding user to database: %w", err)
	}
	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("error setting current user: %w", err)
	}
	fmt.Printf("user successfully created:\n%v\n", newUser)
	return nil
}

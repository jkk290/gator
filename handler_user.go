package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("no username provided")
	}
	name := cmd.Args[0]

	if err := s.cfg.SetUser(name); err != nil {
		return fmt.Errorf("error setting current user: %w", err)
	}

	fmt.Printf("user %s successfully set\n", name)
	return nil
}

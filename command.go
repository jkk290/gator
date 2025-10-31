package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commands struct {
	addedCmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, passedCmd command) error {
	cmd, exists := c.addedCmds[passedCmd.Name]
	if !exists {
		return errors.New("command not found")
	}
	return cmd(s, passedCmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.addedCmds[name] = f
}

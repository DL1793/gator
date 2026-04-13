package main

import (
	"errors"

	"github.com/DL1793/gator/internal/config"
	"github.com/DL1793/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	callback map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.callback[cmd.name]
	if !ok {
		return errors.New("command not found: " + cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

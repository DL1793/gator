package main

import (
	"fmt"

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
	err := c.callback[cmd.name](s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: login <username>")
		return fmt.Errorf("invalid command arguments")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set %s\n", s.cfg.CurrentUserName)
	return nil
}

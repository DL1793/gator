package main

import (
	"fmt"

	"github.com/DL1793/gator/internal/config"
)

type state struct {
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

}

func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("invalid command arguments")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set %s\n", s.cfg.CurrentUserName)
}

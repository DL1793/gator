package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: login <username>")
		return fmt.Errorf("invalid command arguments")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("User not found")
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User set %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: register <username>")
	}

	params := database.CreateUserParams{
		uuid.New(),
		time.Now(),
		time.Now(),
		cmd.args[0]}

	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User %s created\n", cmd.args[0])
	fmt.Printf("User data:\n %v\n, %v\n, %v\n, %v\n", params.ID, params.CreatedAt, params.UpdatedAt, params.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("Usage: users")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("Usage: reset")
	}
	err := s.db.ClearUsers(context.Background())
	if err != nil {
		fmt.Printf("Error clearing users: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Users cleared")
	return nil
}

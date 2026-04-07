package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/DL1793/gator/internal/config"
	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		fmt.Println("Usage: add <name> <feed>")
		return fmt.Errorf("invalid command arguments")
	}

	feed_id := uuid.New()
	now := time.Now()

	params := database.CreateFeedParams{
		feed_id,
		now,
		now,
		cmd.args[0],
		cmd.args[1],
		user.ID,
	}

	feedFollow_params := database.CreateFeedFollowParams{
		uuid.New(),
		now,
		now,
		user.ID,
		feed_id,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed added: %v", feed)

	res, err := s.db.CreateFeedFollow(context.Background(), feedFollow_params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed followed: %v\n Current User: %v\n", res.FeedName, res.UserName)

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

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("Usage: feeds")
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: agg")
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: follow <url>")
	}
	url := cmd.args[0]
	now := time.Now()
	feed_id, err := s.db.GetFeedId(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		uuid.New(),
		now,
		now,
		user.ID,
		feed_id,
	}
	res, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed followed: %v\n Current User: %v\n", res.FeedName, res.UserName)
	return nil
}

func handlerGetFeedFollowsForUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		fmt.Println("Usage: following")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

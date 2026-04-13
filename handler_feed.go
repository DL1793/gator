package main

import (
	"context"
	"fmt"
	"time"

	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
)

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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: unfollow <url>")
	}
	url := cmd.args[0]
	feed_id, err := s.db.GetFeedId(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.DeleteFeedFollowParams{
		user.ID,
		feed_id,
	}
	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
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

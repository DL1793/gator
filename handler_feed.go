package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("usage: feeds")
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
		return errors.New("usage: add <name> <feed>")
	}

	feedId := uuid.New()
	now := time.Now()

	params := database.CreateFeedParams{
		ID:        feedId,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feedId,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Feed added: %v", feed)

	res, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("Feed followed: %v\n Current User: %v\n", res.FeedName, res.UserName)

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("usage: follow <url>")
	}
	url := cmd.args[0]
	now := time.Now()
	feedId, err := s.db.GetFeedId(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feedId,
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
		return errors.New("usage: unfollow <url>")
	}
	url := cmd.args[0]
	feedId, err := s.db.GetFeedId(context.Background(), url)
	if err != nil {
		return err
	}
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedId,
	}
	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func handlerGetFeedFollowsForUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return errors.New("usage: following")
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

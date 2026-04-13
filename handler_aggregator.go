package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/net/html"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("usage: agg <time_between_reqs>")
	}
	timeBetweenReqs := cmd.args[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		fmt.Println("Fetching feeds...")
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Failed to fetch feeds: %v\n", err)
		}
	}
	//return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return errors.New("usage: browse <limit>")
	}
	var limit int
	if len(cmd.args) == 0 {
		limit = 2
	} else {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	if len(posts) == 0 {
		return errors.New("no posts found")
	}
	for _, post := range posts {
		fmt.Printf("* TITLE: %s\n* LINK: %s\n",
			post.Title.String, post.Url)
		if post.Description.Valid {
			desc := stripHTML(post.Description.String)
			if desc != "" {
				fmt.Printf("* DESCRIPTION: %s\n", desc)
			}

		}
		fmt.Printf("* PUB DATE: %s\n\n", post.PublishedAt.Time.Format("Mon Jan 2"))

	}
	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return err
	}
	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, item := range rss.Channel.Items {
		var t time.Time
		now := time.Now()
		var pubDate sql.NullTime
		valid := true
		parsed := false
		for _, format := range formats {
			t, err = time.Parse(format, item.PubDate)
			if err == nil {
				parsed = true
				break
			}
		}
		if !parsed {
			valid = false
		}

		pubDate.Time = t
		pubDate.Valid = valid

		var validTitle bool

		if item.Title != "" {
			validTitle = true
		} else {
			validTitle = false
		}

		params := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title: sql.NullString{
				String: item.Title,
				Valid:  validTitle,
			},
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}
		fmt.Println("Creating", item.Title, "with", item.Description)
		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			if pqErr, ok := errors.AsType[*pq.Error](err); ok && pqErr.Code == "23505" {
				//post already exists, ignore
			} else {
				return err
			}
		}
	}
	return nil
}

func stripHTML(s string) string {
	doc, err := html.ParseFragment(strings.NewReader(s), nil)
	if err != nil {
		return s
	}
	var buf strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, node := range doc {
		walk(node)
	}
	return strings.TrimSpace(buf.String())
}

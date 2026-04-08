package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/DL1793/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("Usage: agg <time_between_reqs>")
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
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	params := database.MarkFeedFetchedParams{
		sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		feed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	}
	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range rss.Channel.Items {
		now := time.Now()
		var pubDate sql.NullTime
		valid := true
		if item.PubDate != "" {
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				t, err = time.Parse(time.RFC1123, item.PubDate)
				if err != nil {
					t, err = time.Parse(time.RFC822Z, item.PubDate)
					if err != nil {
						t, err = time.Parse(time.RFC822, item.PubDate)
						if err != nil {
							t, err = time.Parse(time.RFC3339, item.PubDate)
							if err != nil {
								fmt.Println(err)
								valid = false
							}
						}
					}
				}
			}
			pubDate.Time = t
			pubDate.Valid = valid
		} else {
			pubDate.Valid = false
		}
		var validTitle bool
		if item.Title != "" {
			validTitle = true
		} else {
			validTitle = false
		}

		params := database.CreatePostParams{
			uuid.New(),
			now,
			now,
			sql.NullString{
				item.Title,
				validTitle,
			},
			item.Link,
			sql.NullString{
				item.Description,
				true,
			},
			pubDate,
			feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				//post already exists, ignore
			} else {
				fmt.Println(err)
			}
		}
	}
}

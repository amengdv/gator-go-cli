package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/amengdv/gator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
    feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
    if err != nil {
        return err
    }

    log.Printf("Fetching %v\n", feedToFetch.Name)

    err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
    if err != nil {
        return err
    }

    feed, err := fetchFeed(feedToFetch.Url)
    if err != nil {
        return err
    }

    for _, item := range feed.Channel.Item {
        publishedTime := sql.NullTime{}
        if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
            publishedTime = sql.NullTime{
                Time: t,
                Valid: true,
            }
        }
        _, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
            ID: uuid.New(),
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Title: item.Title,
            Url: item.Link,
            Description: sql.NullString{
                String: item.Description,
                Valid: true,
            },
            PublishedAt: publishedTime,
            FeedID: feedToFetch.ID,
        })
        if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
    }

    log.Printf("Feed %s collected, %v posts found\n", feedToFetch.Name, len(feed.Channel.Item))

    return nil
}

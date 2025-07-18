package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gator/internal/database"
	"gator/rss"

	"github.com/google/uuid"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("%s: invalid time_between_reqs value: %w", cmd.name, err)
	}

	fmt.Printf("%s: collecting feeds every %s\n", cmd.name, timeBetweenRequests)
	fmt.Println()

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s.db)
		if err != nil {
			return fmt.Errorf("%s: error scraping feeds: %w", cmd.name, err)
		}
		fmt.Println()
	}
}

func scrapeFeeds(db *database.Queries) error {
	ctx := context.Background()

	nextFeed, err := db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = scrapeFeed(ctx, db, nextFeed)
	if err != nil {
		return err
	}

	return nil
}

func scrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) error {
	feedData, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	err = db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	for _, post := range feedData.Channel.Item {
		pubDate, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", post.PubDate)

		_, err := db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			Title:     post.Title,
			Url:       post.Link,
			Description: sql.NullString{
				String: post.Description,
				Valid:  len(post.Description) > 0,
			},
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		})

		if err != nil && !strings.Contains("duplicate key value violates unique constraint", err.Error()) {
			return err
		}
	}

	return nil
}

func handlerBrowsePosts(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("%s: invalid limit: %w", cmd.name, err)
		}
		limit = l
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("%s: error getting posts for user: %w", cmd.name, err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Format("Mon Jan 2 2006 15:04:05"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}

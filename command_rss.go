package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gator/internal/database"
	"gator/rss"
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

	fmt.Printf("%s\n", feedData.Channel.Title)
	fmt.Println("=======================")

	for _, item := range feedData.Channel.Item {
		fmt.Printf("%s\n", item.Title)
	}

	return nil
}

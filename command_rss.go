package main

import (
	"context"
	"fmt"
	"gator/rss"
)

func handlerAggregate(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	if len(cmd.args) == 1 {
		feedURL = cmd.args[0]
	}

	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("%s: error fetching feed %s: %w", cmd.name, feedURL, err)
	}

	fmt.Printf("%+v\n", *feed)
	return nil
}

package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	username := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		Name:      username,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: error creating feed: %w", cmd.name, err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: error creating feed follow: %w", cmd.name, err)
	}

	fmt.Printf("feed created: %+v\n", feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("%s: error getting feeds: %w", cmd.name, err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for i, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("%s: error getting user %s: %w", cmd.name, feed.UserID, err)
		}

		fmt.Println(feed.Name)
		fmt.Printf("- %s\n", feed.Url)
		fmt.Printf("- %s\n", user.Name)

		if i < len(feeds)-1 {
			fmt.Println()
		}
	}

	return nil
}

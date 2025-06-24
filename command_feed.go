package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("%s: error getting current user: %w", cmd.name, err)
	}

	username := cmd.args[0]
	url := cmd.args[1]
	now := time.Now().UTC()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: error creating feed: %w", cmd.name, err)
	}

	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
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

	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("%s: error getting feeds: %w", cmd.name, err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for i, feed := range feeds {

		user, err := s.db.GetUserById(ctx, feed.UserID)
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

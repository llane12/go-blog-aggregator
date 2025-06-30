package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("%s: error getting feed %s: %w", cmd.name, url, err)
	}

	now := time.Now().UTC()

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: error creating feed follow: %w", cmd.name, err)
	}

	fmt.Printf("user %s is now following feed %s\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("%s: error getting feed %s: %w", cmd.name, url, err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("%s: error deleting feed follow: %w", cmd.name, err)
	}

	fmt.Printf("%s unfollowed successfully!\n", feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	feed_follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("%s: error getting feed follows for user %s: %w", cmd.name, user.Name, err)
	}

	if len(feed_follows) == 0 {
		fmt.Printf("User %s is not following any feeds\n", user.Name)
		return nil
	}

	fmt.Printf("User %s is following these feeds:\n", user.Name)
	for _, follow := range feed_follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

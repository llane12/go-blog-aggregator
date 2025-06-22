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

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("%s: error getting current user: %w", cmd.name, err)
	}

	username := cmd.args[0]
	url := cmd.args[1]
	now := time.Now().UTC()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
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

	fmt.Printf("feed created: %+v\n", feed)
	return nil
}

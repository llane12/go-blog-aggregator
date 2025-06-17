package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	err := s.cfg.SetUser("")
	if err != nil {
		return fmt.Errorf("error clearing logged in user: %w", err)
	}

	err = s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}

	fmt.Printf("%s: reset successful\n", cmd.name)
	return nil
}

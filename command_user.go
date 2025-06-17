package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), sql.NullString{
		String: username,
		Valid:  true,
	})

	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return err
		} else {
			return fmt.Errorf("%s: user with name %s not registered", cmd.name, username)
		}
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("error setting current user: %w", err)
	}

	fmt.Printf("%s: user %s logged in successfully\n", cmd.name, username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	username := cmd.args[0]
	now := time.Now().UTC()

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name: sql.NullString{
			String: username,
			Valid:  true,
		},
	})

	if err != nil {
		return fmt.Errorf("%s: error creating user: %w", cmd.name, err)
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("%s: user %s registered successfully\n", cmd.name, user.Name.String)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("%s: error getting users: %w", cmd.name, err)
	}

	for _, user := range users {
		if !user.Name.Valid {
			continue
		}

		txt := fmt.Sprintf("* %s", user.Name.String)
		if user.Name.String == s.cfg.CurrentUserName {
			txt += " (current)"
		}
		fmt.Println(txt)
	}

	return nil
}

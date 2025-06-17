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
	if len(cmd.args) == 0 {
		return fmt.Errorf("%s: not enough arguments", cmd.name)
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("%s: too many arguments", cmd.name)
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
		return err
	}

	fmt.Printf("%s: username has been set to %s\n", cmd.name, username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("%s: not enough arguments", cmd.name)
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("%s: too many arguments", cmd.name)
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), sql.NullString{
		String: username,
		Valid:  true,
	})

	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return fmt.Errorf("%s: error opening database: %w", cmd.name, err)
	}

	if err == nil {
		return fmt.Errorf("%s: user with name %s already registered", cmd.name, username)
	}

	now := time.Now().UTC()

	user, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
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

	fmt.Printf("%s: user has been registered %s\n", cmd.name, user.Name.String)
	fmt.Printf("%s: user struct: %+v\n", cmd.name, user)

	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("%s: config username has been set to %s\n", cmd.name, username)

	return nil
}

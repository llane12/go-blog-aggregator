package main

import (
	"fmt"
	"gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("command %s not found", cmd.name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("%s: not enough arguments", cmd.name)
	}

	if len(cmd.args) > 1 {
		return fmt.Errorf("%s: too many arguments", cmd.name)
	}

	username := cmd.args[0]

	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("%s: username has been set to %s\n", cmd.name, username)
	return nil
}

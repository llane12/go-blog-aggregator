package main

import (
	"fmt"
	"gator/internal/config"
	"log"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	s := state{
		Config: &cfg,
	}

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("command name required")
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
}

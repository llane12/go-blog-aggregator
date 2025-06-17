package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", config)

	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	fmt.Printf("Opened database: %+v\n", config.DbUrl)

	s := state{
		cfg: &config,
		db:  database.New(db),
	}

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerListUsers)
	cmds.register("reset", handlerReset)

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

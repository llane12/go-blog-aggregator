package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	err = cfg.SetUser("luke")
	if err != nil {
		log.Fatalf("error setting current user: %v", err)
	}

	cfg2, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg2)
}

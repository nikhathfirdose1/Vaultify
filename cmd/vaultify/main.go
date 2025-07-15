package main

import (
	"fmt"
	"log"
	"github.com/nikhathfirdose1/vaultify/internal/config"
)

func main() {
	fmt.Println("Vaultify service starting...")

	// Load config from YAML
	cfg, err := config.LoadConfig("config/vaultify.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Loaded config. Server will run on port %d\n", cfg.Server.Port)
}

package main

import (
	"log"

	"github-analyzer/src/api"
	"github-analyzer/src/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := config.Load()

	server := api.NewServer(cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

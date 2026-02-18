package config

import (
	"log"
	"os"
)

type Config struct {
	Port             string
	OpenRouterAPIKey string
	GitHubToken      string
}

func Load() *Config {
	openRouterKey := os.Getenv("OPENROUTER_API_KEY")
	if openRouterKey == "" {
		log.Fatal("OPENROUTER_API_KEY is required in .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	return &Config{
		Port:             port,
		OpenRouterAPIKey: openRouterKey,
		GitHubToken:      githubToken,
	}
}

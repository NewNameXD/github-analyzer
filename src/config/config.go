package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	GroqAPIKey  string
	GitHubToken string
}

func Load() *Config {
	groqKey := os.Getenv("GROQ_API_KEY")
	if groqKey == "" {
		log.Fatal("GROQ_API_KEY is required in .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	githubToken := os.Getenv("GITHUB_TOKEN")

	return &Config{
		Port:        port,
		GroqAPIKey:  groqKey,
		GitHubToken: githubToken,
	}
}

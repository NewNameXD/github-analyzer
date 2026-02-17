package api

import (
	"net/http"

	"github-analyzer/src/config"
	"github-analyzer/src/handlers"
	"github-analyzer/src/services"
)

type Server struct {
	config  *config.Config
	handler *handlers.Handler
}

func NewServer(cfg *config.Config) *Server {
	githubService := services.NewGitHubService(cfg.GitHubToken)
	aiService := services.NewAIService(cfg.GroqAPIKey)

	handler := handlers.NewHandler(githubService, aiService)

	return &Server{
		config:  cfg,
		handler: handler,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/evaluate", s.handler.HandleEvaluate)
	mux.HandleFunc("/health", s.handler.HandleHealth)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/web/", http.StripPrefix("/web/", fs))

	mux.HandleFunc("/", s.handler.ServeHome)

	handler := enableCORS(mux)

	return http.ListenAndServe(":"+s.config.Port, handler)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

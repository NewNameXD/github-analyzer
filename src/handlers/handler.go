package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github-analyzer/src/models"
	"github-analyzer/src/services"
)

type Handler struct {
	githubService *services.GitHubService
	aiService     *services.AIService
}

func NewHandler(githubService *services.GitHubService, aiService *services.AIService) *Handler {
	return &Handler{
		githubService: githubService,
		aiService:     aiService,
	}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *Handler) ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "web/index.html")
}

func (h *Handler) HandleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.EvaluationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, models.EvaluationResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.Username == "" {
		respondJSON(w, http.StatusBadRequest, models.EvaluationResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	profile, err := h.githubService.FetchProfile(req.Username)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.EvaluationResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to fetch GitHub profile: %v", err),
		})
		return
	}

	evaluation, err := h.aiService.EvaluateProfile(profile, req.Language)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, models.EvaluationResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to generate evaluation: %v", err),
		})
		return
	}

	respondJSON(w, http.StatusOK, models.EvaluationResponse{
		Success:    true,
		Profile:    profile,
		Evaluation: evaluation,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github-analyzer/src/models"
)

type AIService struct {
	apiKey string
	client *http.Client
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message groqMessage `json:"message"`
	} `json:"choices"`
}

func NewAIService(apiKey string) *AIService {
	return &AIService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *AIService) EvaluateProfile(profile *models.GitHubProfile, language string) (string, error) {
	prompt := s.buildPrompt(profile, language)

	reqBody := groqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Groq API error %d: %s", resp.StatusCode, string(body))
	}

	var groqResp groqResponse
	json.NewDecoder(resp.Body).Decode(&groqResp)

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return groqResp.Choices[0].Message.Content, nil
}

func (s *AIService) buildPrompt(profile *models.GitHubProfile, language string) string {
	var sb strings.Builder

	languageMap := map[string]string{
		"en": "Respond in English.",
		"pt": "Responda em Português Brasileiro.",
		"es": "Responde en Español.",
		"fr": "Répondez en Français.",
		"de": "Antworten Sie auf Deutsch.",
		"ja": "日本語で回答してください。",
		"zh": "请用中文回答。",
	}

	languageInstruction := languageMap[language]
	if languageInstruction == "" {
		languageInstruction = languageMap["en"]
	}

	sb.WriteString(languageInstruction + " ")
	sb.WriteString("You are a senior software engineer and technical recruiter analyzing GitHub profiles. ")
	sb.WriteString("Provide a brutally honest, direct assessment without sugar-coating. ")
	sb.WriteString("Focus on facts, patterns, and actionable insights.\n\n")

	sb.WriteString(fmt.Sprintf("PROFILE: %s (@%s)\n", profile.Name, profile.Username))
	sb.WriteString(fmt.Sprintf("Bio: %s\n", profile.Bio))
	sb.WriteString(fmt.Sprintf("Location: %s | Company: %s\n", profile.Location, profile.Company))
	sb.WriteString(fmt.Sprintf("Followers: %d | Following: %d | Repos: %d\n\n", profile.Followers, profile.Following, profile.PublicRepos))

	sb.WriteString("REPOSITORIES:\n")
	for i, repo := range profile.Repositories {
		if i >= 20 {
			break
		}
		sb.WriteString(fmt.Sprintf("- %s (%s) | Stars: %d | Forks: %d | README: %v\n",
			repo.Name, repo.Language, repo.Stars, repo.Forks, repo.HasReadme))
		if repo.Description != "" {
			sb.WriteString(fmt.Sprintf("  Description: %s\n", repo.Description))
		}
		if len(repo.Topics) > 0 {
			sb.WriteString(fmt.Sprintf("  Topics: %v\n", repo.Topics))
		}
	}

	sb.WriteString("\n\nProvide a comprehensive analysis in the following structure:\n\n")
	sb.WriteString("## Technical Assessment\n")
	sb.WriteString("Evaluate code quality indicators, technology stack diversity, and project complexity.\n\n")
	sb.WriteString("## Profile Strengths\n")
	sb.WriteString("What genuinely stands out. Be specific with examples.\n\n")
	sb.WriteString("## Critical Weaknesses\n")
	sb.WriteString("What's missing or poorly executed. Don't hold back.\n\n")
	sb.WriteString("## Recommendations\n")
	sb.WriteString("Concrete, prioritized actions to improve profile credibility and technical presence.\n\n")
	sb.WriteString("## Market Perception\n")
	sb.WriteString("How this profile would be perceived by recruiters and technical hiring managers.\n\n")
	sb.WriteString("Be direct, factual, and constructive. Avoid generic advice.")

	return sb.String()
}

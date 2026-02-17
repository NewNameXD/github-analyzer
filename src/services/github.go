package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github-analyzer/src/models"
)

type GitHubService struct {
	client  *http.Client
	proxies []string
	retries int
	token   string
}

type githubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	PublicRepos int    `json:"public_repos"`
	AvatarURL   string `json:"avatar_url"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Blog        string `json:"blog"`
}

type githubRepo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Stars       int      `json:"stargazers_count"`
	Forks       int      `json:"forks_count"`
	Topics      []string `json:"topics"`
	UpdatedAt   string   `json:"updated_at"`
	Fork        bool     `json:"fork"`
	HTMLURL     string   `json:"html_url"`
}

func NewGitHubService(token string) *GitHubService {
	return &GitHubService{
		client: &http.Client{Timeout: 15 * time.Second},
		proxies: []string{
			"https://api.github.com",
		},
		retries: 3,
		token:   token,
	}
}

func (s *GitHubService) FetchProfile(username string) (*models.GitHubProfile, error) {
	user, err := s.fetchUser(username)
	if err != nil {
		return nil, err
	}

	profile := &models.GitHubProfile{
		Username:    username,
		Name:        user.Name,
		Bio:         user.Bio,
		Followers:   user.Followers,
		Following:   user.Following,
		PublicRepos: user.PublicRepos,
		AvatarURL:   user.AvatarURL,
		Company:     user.Company,
		Location:    user.Location,
		Blog:        user.Blog,
	}

	repos, err := s.fetchRepositories(username)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		hasReadme := s.checkReadme(username, repo.Name)

		profile.Repositories = append(profile.Repositories, models.Repository{
			Name:        repo.Name,
			Description: repo.Description,
			Language:    repo.Language,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			HasReadme:   hasReadme,
			Topics:      repo.Topics,
			UpdatedAt:   repo.UpdatedAt,
			HTMLURL:     repo.HTMLURL,
		})
	}

	return profile, nil
}

func (s *GitHubService) fetchUser(username string) (*githubUser, error) {
	var lastErr error

	for attempt := 0; attempt < s.retries; attempt++ {
		for _, baseURL := range s.proxies {
			url := fmt.Sprintf("%s/users/%s", baseURL, username)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Accept", "application/vnd.github.v3+json")
			req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")
			if s.token != "" {
				req.Header.Set("Authorization", "Bearer "+s.token)
			}

			resp, err := s.client.Do(req)
			if err != nil {
				lastErr = err
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusForbidden {
				lastErr = fmt.Errorf("rate limit exceeded")
				time.Sleep(time.Second * time.Duration(attempt+1))
				continue
			}

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("status %d", resp.StatusCode)
				continue
			}

			var user githubUser
			if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
				lastErr = err
				continue
			}

			return &user, nil
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %v", s.retries, lastErr)
}

func (s *GitHubService) fetchRepositories(username string) ([]githubRepo, error) {
	var lastErr error

	for attempt := 0; attempt < s.retries; attempt++ {
		for _, baseURL := range s.proxies {
			url := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", baseURL, username)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("Accept", "application/vnd.github.mercy-preview+json")
			req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")
			if s.token != "" {
				req.Header.Set("Authorization", "Bearer "+s.token)
			}

			resp, err := s.client.Do(req)
			if err != nil {
				lastErr = err
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusForbidden {
				lastErr = fmt.Errorf("rate limit exceeded")
				time.Sleep(time.Second * time.Duration(attempt+1))
				continue
			}

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("status %d", resp.StatusCode)
				continue
			}

			var repos []githubRepo
			if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
				lastErr = err
				continue
			}

			return repos, nil
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %v", s.retries, lastErr)
}

func (s *GitHubService) checkReadme(owner, repo string) bool {
	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/repos/%s/%s/readme", baseURL, owner, repo)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true
		}
	}

	return false
}

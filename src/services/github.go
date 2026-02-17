package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Login           string `json:"login"`
	Name            string `json:"name"`
	Bio             string `json:"bio"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	PublicRepos     int    `json:"public_repos"`
	AvatarURL       string `json:"avatar_url"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Blog            string `json:"blog"`
	TwitterUsername string `json:"twitter_username"`
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

type githubContent struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type githubReadme struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

type githubOrg struct {
	Login       string `json:"login"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
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
		Username:        username,
		Name:            user.Name,
		Bio:             user.Bio,
		Followers:       user.Followers,
		Following:       user.Following,
		PublicRepos:     user.PublicRepos,
		AvatarURL:       user.AvatarURL,
		Company:         user.Company,
		Location:        user.Location,
		Blog:            user.Blog,
		TwitterUsername: user.TwitterUsername,
	}

	repos, err := s.fetchRepositories(username)
	if err != nil {
		return nil, err
	}

	totalStars := 0
	totalForks := 0

	for _, repo := range repos {
		if repo.Fork {
			continue
		}

		totalStars += repo.Stars
		totalForks += repo.Forks

		readmeContent := s.fetchReadme(username, repo.Name)
		structure := s.fetchRepoStructure(username, repo.Name)

		profile.Repositories = append(profile.Repositories, models.Repository{
			Name:          repo.Name,
			Description:   repo.Description,
			Language:      repo.Language,
			Stars:         repo.Stars,
			Forks:         repo.Forks,
			HasReadme:     readmeContent != "",
			Topics:        repo.Topics,
			UpdatedAt:     repo.UpdatedAt,
			HTMLURL:       repo.HTMLURL,
			ReadmeContent: readmeContent,
			Structure:     structure,
		})
	}

	profile.TotalStars = totalStars
	profile.TotalForks = totalForks

	orgs, _ := s.fetchOrganizations(username)
	profile.Organizations = orgs

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

func (s *GitHubService) fetchReadme(owner, repo string) string {
	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/repos/%s/%s/readme", baseURL, owner, repo)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")
		if s.token != "" {
			req.Header.Set("Authorization", "Bearer "+s.token)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var readme githubReadme
			if err := json.NewDecoder(resp.Body).Decode(&readme); err != nil {
				continue
			}

			if readme.Encoding == "base64" {
				decoded, err := decodeBase64(readme.Content)
				if err == nil {
					return decoded
				}
			}
			return readme.Content
		}
	}

	return ""
}

func (s *GitHubService) fetchRepoStructure(owner, repo string) []string {
	var structure []string

	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/repos/%s/%s/contents", baseURL, owner, repo)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")
		if s.token != "" {
			req.Header.Set("Authorization", "Bearer "+s.token)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var contents []githubContent
			if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
				continue
			}

			for _, item := range contents {
				if item.Type == "dir" {
					structure = append(structure, item.Name+"/")
				} else {
					structure = append(structure, item.Name)
				}
			}
			break
		}
	}

	return structure
}

func decodeBase64(content string) (string, error) {
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")

	decoded := make([]byte, len(content))
	n := 0

	const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

	for i := 0; i < len(content); i += 4 {
		if i+3 >= len(content) {
			break
		}

		b1 := strings.IndexByte(base64Chars, content[i])
		b2 := strings.IndexByte(base64Chars, content[i+1])
		b3 := strings.IndexByte(base64Chars, content[i+2])
		b4 := strings.IndexByte(base64Chars, content[i+3])

		if b1 < 0 || b2 < 0 {
			continue
		}

		decoded[n] = byte((b1 << 2) | (b2 >> 4))
		n++

		if b3 >= 0 && content[i+2] != '=' {
			decoded[n] = byte((b2 << 4) | (b3 >> 2))
			n++
		}

		if b4 >= 0 && content[i+3] != '=' {
			decoded[n] = byte((b3 << 6) | b4)
			n++
		}
	}

	return string(decoded[:n]), nil
}

func (s *GitHubService) fetchOrganizations(username string) ([]models.Organization, error) {
	var orgs []models.Organization

	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/users/%s/orgs", baseURL, username)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "GitHub-Profile-Evaluator")
		if s.token != "" {
			req.Header.Set("Authorization", "Bearer "+s.token)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var githubOrgs []githubOrg
			if err := json.NewDecoder(resp.Body).Decode(&githubOrgs); err != nil {
				continue
			}

			for _, org := range githubOrgs {
				orgs = append(orgs, models.Organization{
					Login:       org.Login,
					AvatarURL:   org.AvatarURL,
					Description: org.Description,
				})
			}
			break
		}
	}

	return orgs, nil
}

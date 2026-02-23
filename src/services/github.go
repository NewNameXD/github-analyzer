package services

import (
	"encoding/json"
	"fmt"
	"io"
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
		client: &http.Client{Timeout: 10 * time.Second},
		proxies: []string{
			"https://api.github.com",
		},
		retries: 2,
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

	repoLimit := 10
	repoCount := 0

	type repoData struct {
		repo      githubRepo
		readme    string
		structure []string
		codeFiles map[string]string
	}

	repoChan := make(chan repoData, repoLimit)

	go func() {
		for _, repo := range repos {
			if repo.Fork {
				continue
			}
			if repoCount >= repoLimit {
				break
			}
			repoCount++

			data := repoData{repo: repo}

			if repoCount <= 5 {
				data.readme = s.fetchReadme(username, repo.Name)
				if data.readme == "" {
					data.structure = s.fetchRepoStructure(username, repo.Name)
				}
				data.codeFiles = s.fetchCodeFiles(username, repo.Name)
			}

			repoChan <- data
		}
		close(repoChan)
	}()

	for data := range repoChan {
		totalStars += data.repo.Stars
		totalForks += data.repo.Forks

		profile.Repositories = append(profile.Repositories, models.Repository{
			Name:          data.repo.Name,
			Description:   data.repo.Description,
			Language:      data.repo.Language,
			Stars:         data.repo.Stars,
			Forks:         data.repo.Forks,
			HasReadme:     data.readme != "",
			Topics:        data.repo.Topics,
			UpdatedAt:     data.repo.UpdatedAt,
			HTMLURL:       data.repo.HTMLURL,
			ReadmeContent: data.readme,
			Structure:     data.structure,
			CodeFiles:     data.codeFiles,
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

func (s *GitHubService) fetchCodeFiles(owner, repo string) map[string]string {
	codeFiles := make(map[string]string)
	filesToSkip := map[string]bool{
		"node_modules": true, "package-lock.json": true,
		"go.sum": true, "go.mod": true, ".env": true,
		"yarn.lock": true, "pnpm-lock.yaml": true, ".DS_Store": true,
		"dist": true, "build": true, ".git": true, ".github": true,
		"coverage": true, ".next": true, "out": true, "target": true,
		".venv": true, "venv": true, "__pycache__": true, ".pytest_cache": true,
		"vendor": true, ".bundle": true, "Gemfile.lock": true,
		"package.json": true, "Cargo.lock": true, "poetry.lock": true,
		".gitignore": true, ".env.example": true, ".env.local": true,
		"LICENSE": true, "CHANGELOG": true, "CHANGELOG.md": true,
	}

	s.fetchFilesRecursive(owner, repo, "", codeFiles, filesToSkip, 0)
	return codeFiles
}

func (s *GitHubService) fetchFilesRecursive(owner, repo, path string, codeFiles map[string]string, skip map[string]bool, depth int) {
	if depth > 2 || len(codeFiles) > 8 {
		return
	}

	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", baseURL, owner, repo, path)
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
				if skip[item.Name] {
					continue
				}

				if item.Type == "dir" {
					newPath := item.Path
					if newPath != "" {
						s.fetchFilesRecursive(owner, repo, newPath, codeFiles, skip, depth+1)
					}
				} else if isCodeFile(item.Name) && len(codeFiles) < 8 {
					content := s.fetchFileContent(owner, repo, item.Path)
					if content != "" && len(content) < 4000 {
						codeFiles[item.Path] = content
					}
				}
			}
			break
		}
	}
}

func isCodeFile(filename string) bool {
	skipExtensions := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".svg": true,
		".ico": true, ".webp": true, ".bmp": true, ".tiff": true, ".psd": true,
		".ai": true, ".eps": true, ".heic": true, ".raw": true,
		".mp3": true, ".mp4": true, ".wav": true, ".flac": true, ".aac": true,
		".m4a": true, ".ogg": true, ".wma": true, ".mov": true, ".avi": true,
		".mkv": true, ".webm": true, ".m3u8": true, ".m2ts": true,
		".zip": true, ".tar": true, ".gz": true, ".rar": true, ".7z": true,
		".bz2": true, ".xz": true, ".iso": true, ".dmg": true, ".apk": true,
		".exe": true, ".dll": true, ".so": true, ".dylib": true, ".o": true,
		".a": true, ".lib": true, ".class": true, ".pyc": true, ".pyo": true,
		".whl": true, ".jar": true, ".war": true, ".ear": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".odt": true, ".ods": true, ".odp": true,
		".rtf": true, ".pages": true, ".numbers": true, ".keynote": true,
		".db": true, ".sqlite": true, ".sqlite3": true, ".mdb": true,
		".accdb": true, ".dbf": true, ".ibd": true,
		".ttf": true, ".otf": true, ".woff": true, ".woff2": true, ".eot": true,
		".swp": true, ".swo": true, ".swn": true,
	}

	skipFiles := map[string]bool{
		"package-lock.json": true, "yarn.lock": true, "pnpm-lock.yaml": true,
		"Gemfile.lock": true, "Cargo.lock": true, "poetry.lock": true,
		"composer.lock": true, "pubspec.lock": true, "mix.lock": true,
		".DS_Store": true, "Thumbs.db": true, "desktop.ini": true,
		".gradle": true, ".m2": true, ".ivy2": true,
		"gradlew": true, "gradlew.bat": true,
		"mvnw": true, "mvnw.cmd": true,
		".npmrc": true, ".yarnrc": true, ".nvmrc": true,
		".ruby-version": true, ".python-version": true,
		"Procfile": true, "Procfile.dev": true,
		".editorconfig": true, ".prettierignore": true,
		".eslintignore": true, ".stylelintignore": true,
		"CNAME": true, "robots.txt": true, "sitemap.xml": true,
		"manifest.json": true, "browserconfig.xml": true,
		".htaccess": true, "web.config": true,
	}

	if skipFiles[filename] {
		return false
	}

	for ext := range skipExtensions {
		if strings.HasSuffix(filename, ext) {
			return false
		}
	}

	return true
}

func (s *GitHubService) fetchFileContent(owner, repo, path string) string {
	for _, baseURL := range s.proxies {
		url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", baseURL, owner, repo, path)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/vnd.github.v3.raw")
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
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				return string(body)
			}
		}
	}
	return ""
}

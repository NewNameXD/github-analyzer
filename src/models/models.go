package models

type GitHubProfile struct {
	Username        string         `json:"username"`
	Name            string         `json:"name"`
	Bio             string         `json:"bio"`
	Followers       int            `json:"followers"`
	Following       int            `json:"following"`
	PublicRepos     int            `json:"public_repos"`
	AvatarURL       string         `json:"avatar_url"`
	Company         string         `json:"company"`
	Location        string         `json:"location"`
	Blog            string         `json:"blog"`
	TwitterUsername string         `json:"twitter_username,omitempty"`
	Repositories    []Repository   `json:"repositories"`
	Organizations   []Organization `json:"organizations,omitempty"`
	TotalStars      int            `json:"total_stars"`
	TotalForks      int            `json:"total_forks"`
}

type Organization struct {
	Login       string `json:"login"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description,omitempty"`
}

type Repository struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Language      string            `json:"language"`
	Stars         int               `json:"stars"`
	Forks         int               `json:"forks"`
	HasReadme     bool              `json:"has_readme"`
	Topics        []string          `json:"topics"`
	UpdatedAt     string            `json:"updated_at"`
	HTMLURL       string            `json:"html_url"`
	ReadmeContent string            `json:"readme_content,omitempty"`
	Structure     []string          `json:"structure,omitempty"`
	CodeFiles     map[string]string `json:"code_files,omitempty"`
}

type EvaluationRequest struct {
	Username string `json:"username"`
	Language string `json:"language"`
}

type EvaluationResponse struct {
	Success    bool           `json:"success"`
	Profile    *GitHubProfile `json:"profile,omitempty"`
	Evaluation string         `json:"evaluation,omitempty"`
	Error      string         `json:"error,omitempty"`
}

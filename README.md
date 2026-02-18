# GitHub Profile Evaluator

> Honest, AI-powered analysis of GitHub profiles with brutally direct feedback

A professional web application that analyzes GitHub profiles using OpenRouter AI (Llama 3.3 70B) to provide data-driven insights, strengths, weaknesses, and actionable recommendations.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![OpenRouter AI](https://img.shields.io/badge/AI-OpenRouter-orange.svg)

## Screenshots

### Home Page
![Home Page](assets/home.png)

### Analysis Result
![Analysis Result](assets/analyze.png)

## Features

- **AI-Powered Analysis**: Uses OpenRouter with automatic model selection (Auto Free → Gemini 2.0 → Llama 3.3 70B → Llama 3.1 405B → Mistral → Qwen)
- **Smart Fallback System**: Automatically switches models if one is rate-limited
- **IP-Based Rate Limiting**: 1 request per minute per IP address to prevent abuse
- **Repository Deep Dive**: Analyzes folder structure and README content of each repository
- **Smart Prompt Splitting**: Automatically splits large profiles into two API calls if needed
- **Natural Language**: Responses sound like a real developer, not corporate AI
- **Multi-Language Support**: Interface and evaluation in 7 languages (EN, PT, ES, FR, DE, JA, ZH)
- **Comprehensive Evaluation**: Analyzes repositories, documentation, activity, and market perception
- **Modern UI**: Beautiful, responsive interface built with Tailwind CSS
- **Language Persistence**: Remembers your language preference via cookies
- **REST API**: Full API access for integration with other tools
- **Free Tier**: Uses OpenRouter's free models with intelligent fallback

## Quick Start

### Prerequisites

- Go 1.21 or higher
- OpenRouter API Key (free at [openrouter.ai/keys](https://openrouter.ai/keys))
- Optional: GitHub Token (to avoid rate limits)

### Installation

```bash
# Clone the repository
git clone https://github.com/pqpcara/github-analyzer.git
cd github-analyzer

# Install dependencies
go mod download

# Configure environment
cp .env.example .env
```

### Configuration

Edit `.env` file:

```env
# OpenRouter API Key (REQUIRED)
# Get it at: https://openrouter.ai/keys
OPENROUTER_API_KEY=your_openrouter_api_key

# GitHub Token (OPTIONAL - helps avoid rate limits)
# Get it at: https://github.com/settings/tokens
GITHUB_TOKEN=api_key

# Optional: Custom port (default: 8080)
PORT=8080
```

### Run

```bash
go run main.go
```

### Docker

```bash
docker build -t github-analyzer .
docker run -p 8080:8080 -e OPENROUTER_API_KEY=your_key -e GITHUB_TOKEN=github_token github-analyzer
```

Server will start on `http://localhost:8080`

## API Documentation

### Evaluate Profile

**Endpoint:** `POST /api/evaluate`

**Request:**
```json
{
  "username": "torvalds",
  "language": "en"
}
```

**Response:**
```json
{
  "success": true,
  "profile": {
    "username": "torvalds",
    "name": "Linus Torvalds",
    "bio": "...",
    "followers": 150000,
    "following": 0,
    "public_repos": 5,
    "avatar_url": "https://...",
    "repositories": [...]
  },
  "evaluation": "## Technical Assessment\n..."
}
```

**Language Codes:**
- `en` - English
- `pt` - Português
- `es` - Español
- `fr` - Français
- `de` - Deutsch
- `ja` - 日本語
- `zh` - 中文

### Health Check

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "healthy"
}
```

## Features in Detail

### AI Evaluation Sections

1. **Technical Assessment**: Code quality, tech stack diversity, project complexity
2. **Profile Strengths**: What stands out with specific examples
3. **Critical Weaknesses**: Areas needing improvement (no sugar-coating)
4. **Recommendations**: Concrete, prioritized actions
5. **Market Perception**: How recruiters/hiring managers view the profile

### Repository Analysis

For each repository, the system now analyzes:
- **Folder Structure**: Lists all files and directories in the root
- **README Content**: Reads and includes README.md content (up to 500 chars preview)
- **Technologies**: Identifies programming languages and frameworks
- **Topics/Tags**: Extracts repository topics for better categorization

This provides much deeper insights into code organization and project quality.

### Multi-Language Interface

The entire interface automatically translates when you change the language:
- All UI text and labels
- Placeholders and buttons
- Error messages
- Support card content

Your language preference is saved in cookies and restored on next visit.

## Tech Stack

- **Backend**: Go (net/http)
- **AI**: OpenRouter API with intelligent auto-selection (Auto Free → Gemini 2.0 → Llama 3.3 → Llama 3.1 405B → Mistral → Qwen)
- **Frontend**: Vanilla JavaScript, Tailwind CSS
- **Data Source**: GitHub REST API

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Star the repository
2. Fork the repository
3. Create your feature branch (`git checkout -b feature/AmazingFeature`)
4. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. Push to the branch (`git push origin feature/AmazingFeature`)
6. Open a Pull Request

## Contact

Created by [@pqpcara](https://github.com/pqpcara)

If you find this project helpful, consider:
- Starring the repository
- Reporting bugs
- Suggesting new features
- Contributing code

---

**Made with ❤️ by [@pqpcara](https://github.com/pqpcara)**

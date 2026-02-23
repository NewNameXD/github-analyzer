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

type openRouterRequest struct {
	Model    string              `json:"model"`
	Messages []openRouterMessage `json:"messages"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message openRouterMessage `json:"message"`
	} `json:"choices"`
}

func NewAIService(apiKey string) *AIService {
	return &AIService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 180 * time.Second},
	}
}

func (s *AIService) EvaluateProfile(profile *models.GitHubProfile, language string) (string, error) {
	prompt := s.buildPrompt(profile, language, len(profile.Repositories))

	result, err := s.callOpenRouterAPI(prompt)

	if err != nil && (strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "413") || strings.Contains(err.Error(), "too large")) {
		return s.evaluateInTwoParts(profile, language)
	}

	return result, err
}

func (s *AIService) evaluateInTwoParts(profile *models.GitHubProfile, language string) (string, error) {
	totalRepos := len(profile.Repositories)
	halfRepos := totalRepos / 2

	prompt1 := s.buildPrompt(profile, language, halfRepos)

	noteMap := map[string]string{
		"en": "\n\nNOTE: This is PART 1 of 2. Analyze these repos and wait for part 2 before final conclusion.",
		"pt": "\n\nNOTA: Esta é a PARTE 1 de 2. Analise estes repos e aguarde a parte 2 antes da conclusão final.",
		"es": "\n\nNOTA: Esta es la PARTE 1 de 2. Analiza estos repos y espera la parte 2 antes de la conclusión final.",
		"fr": "\n\nNOTE: Ceci est la PARTIE 1 de 2. Analysez ces repos et attendez la partie 2 avant la conclusion finale.",
		"de": "\n\nHINWEIS: Dies ist TEIL 1 von 2. Analysieren Sie diese Repos und warten Sie auf Teil 2 vor der endgültigen Schlussfolgerung.",
		"ja": "\n\n注意：これは2つのうちのパート1です。これらのリポジトリを分析し、最終結論の前にパート2を待ってください。",
		"zh": "\n\n注意：这是第1部分，共2部分。分析这些仓库并在最终结论前等待第2部分。",
	}

	note := noteMap[language]
	if note == "" {
		note = noteMap["en"]
	}

	prompt1 += note

	result1, err := s.callOpenRouterAPI(prompt1)
	if err != nil {
		return "", fmt.Errorf("error in part 1: %v", err)
	}

	profilePart2 := &models.GitHubProfile{
		Username:     profile.Username,
		Name:         profile.Name,
		Bio:          profile.Bio,
		Followers:    profile.Followers,
		Following:    profile.Following,
		PublicRepos:  profile.PublicRepos,
		AvatarURL:    profile.AvatarURL,
		Company:      profile.Company,
		Location:     profile.Location,
		Blog:         profile.Blog,
		Repositories: profile.Repositories[halfRepos:],
	}

	prompt2 := s.buildPrompt(profilePart2, language, len(profilePart2.Repositories))
	prompt2 = strings.Replace(prompt2, "REPOSITORIES:\n", "REPOSITORIES (CONTINUED):\n", 1)

	note2Map := map[string]string{
		"en": "\n\nNOTE: This is PART 2 of 2. Combine with previous analysis and provide the complete final evaluation.",
		"pt": "\n\nNOTA: Esta é a PARTE 2 de 2. Combine com a análise anterior e forneça a avaliação final completa.",
		"es": "\n\nNOTA: Esta es la PARTE 2 de 2. Combina con el análisis anterior y proporciona la evaluación final completa.",
		"fr": "\n\nNOTE: Ceci est la PARTIE 2 de 2. Combinez avec l'analyse précédente et fournissez l'évaluation finale complète.",
		"de": "\n\nHINWEIS: Dies ist TEIL 2 von 2. Kombinieren Sie mit der vorherigen Analyse und geben Sie die vollständige Endbewertung ab.",
		"ja": "\n\n注意：これは2つのうちのパート2です。前の分析と組み合わせて、完全な最終評価を提供してください。",
		"zh": "\n\n注意：这是第2部分，共2部分。与之前的分析结合，提供完整的最终评估。",
	}

	note2 := note2Map[language]
	if note2 == "" {
		note2 = note2Map["en"]
	}

	prompt2 += note2

	messages := []openRouterMessage{
		{Role: "user", Content: prompt1},
		{Role: "assistant", Content: result1},
		{Role: "user", Content: prompt2},
	}

	result2, err := s.callOpenRouterAPIWithMessages(messages)
	if err != nil {
		return "", fmt.Errorf("error in part 2: %v", err)
	}

	return result2, nil
}

func (s *AIService) callOpenRouterAPI(prompt string) (string, error) {
	return s.callOpenRouterAPIWithMessages([]openRouterMessage{
		{Role: "user", Content: prompt},
	})
}

func (s *AIService) callOpenRouterAPIWithMessages(messages []openRouterMessage) (string, error) {
	freeModels := []string{
		"openrouter/auto:free",
		"google/gemini-2.0-flash-exp:free",
		"meta-llama/llama-3.3-70b-instruct:free",
		"meta-llama/llama-3.1-405b-instruct:free",
		"mistralai/mistral-small-3.1:free",
		"qwen/qwen-2.5-72b-instruct:free",
	}

	var lastErr error
	for _, model := range freeModels {
		result, err := s.tryModelRequest(model, messages)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !strings.Contains(err.Error(), "429") &&
			!strings.Contains(err.Error(), "402") &&
			!strings.Contains(err.Error(), "404") &&
			!strings.Contains(err.Error(), "rate") {
			return "", err
		}
	}

	return "", fmt.Errorf("all free models are rate-limited or unavailable: %v", lastErr)
}

func (s *AIService) tryModelRequest(model string, messages []openRouterMessage) (string, error) {
	reqBody := openRouterRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/pqpcara/github-analyzer")
	req.Header.Set("X-Title", "GitHub Profile Evaluator")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter API error %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResp); err != nil {
		return "", err
	}

	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return openRouterResp.Choices[0].Message.Content, nil
}

func (s *AIService) buildPrompt(profile *models.GitHubProfile, language string, maxRepos int) string {
	var sb strings.Builder

	type promptText struct {
		instruction     string
		role            string
		structure       string
		technical       string
		strengths       string
		weaknesses      string
		recommendations string
		perception      string
		closing         string
	}

	languageMap := map[string]promptText{
		"en": {
			instruction:     "Respond in English.",
			role:            "Read the code. Say what's good, what's shit. No fluff. No corporate speak. Just facts. CRITICAL: Use ONLY ## (two hashtags) for titles. NEVER use # alone.",
			structure:       "",
			technical:       "## Code Quality\n\nWhat's the actual quality? Look at error handling, validation, architecture. Is it clean or messy? Does it handle edge cases? Be specific - point to actual problems you see in the code.",
			strengths:       "## What Works\n\nWhat's actually good in the code? Be specific. Don't say generic shit like 'good structure' - say WHAT is good and WHY.",
			weaknesses:      "## What's Broken\n\nWhat's bad? Missing error handling? No validation? Messy code? Security issues? Say it straight. Point to actual problems.",
			recommendations: "## Fix This\n\nWhat needs to change? Be specific. Don't say 'improve error handling' - say WHERE and HOW.",
			perception:      "## Level\n\nJunior, mid, or senior? Based on what you saw in the code. What can they handle? What can't they?",
			closing:         "",
		},
		"pt": {
			instruction:     "Responda em Português Brasileiro.",
			role:            "Leia o código. Fale o que é bom, o que é ruim. Sem enrolação. Sem papo corporativo. Só fatos. CRÍTICO: Use APENAS ## (dois hashtags) para títulos. NUNCA use # sozinho.",
			structure:       "",
			technical:       "## Qualidade do Código\n\nQual é a qualidade real? Olhe tratamento de erro, validação, arquitetura. Tá limpo ou bagunçado? Trata casos extremos? Seja específico - aponte problemas reais que você viu no código.",
			strengths:       "## O Que Funciona\n\nO que é realmente bom no código? Seja específico. Não fale genérico tipo 'boa estrutura' - fale O QUE é bom e POR QUÊ.",
			weaknesses:      "## O Que Tá Quebrado\n\nO que é ruim? Falta tratamento de erro? Sem validação? Código bagunçado? Problemas de segurança? Fale direto. Aponte problemas reais.",
			recommendations: "## Conserta Isso\n\nO que precisa mudar? Seja específico. Não fale 'melhore tratamento de erro' - fale ONDE e COMO.",
			perception:      "## Nível\n\nJunior, mid ou senior? Baseado no que você viu no código. O que ele consegue fazer? O que não consegue?",
			closing:         "",
		},
		"es": {
			instruction:     "Responde en Español.",
			role:            "Lee el código. Di qué es bueno, qué es malo. Sin rodeos. Sin habla corporativa. Solo hechos. CRÍTICO: Usa SOLO ## (dos hashtags) para títulos. NUNCA uses # solo.",
			structure:       "",
			technical:       "## Calidad del Código\n\n¿Cuál es la calidad real? Mira manejo de errores, validación, arquitectura. ¿Está limpio o desordenado? ¿Maneja casos extremos? Sé específico - señala problemas reales que viste en el código.",
			strengths:       "## Lo Que Funciona\n\n¿Qué es realmente bueno en el código? Sé específico. No digas cosas genéricas como 'buena estructura' - di QUÉ es bueno y POR QUÉ.",
			weaknesses:      "## Lo Que Está Roto\n\n¿Qué es malo? ¿Falta manejo de errores? ¿Sin validación? ¿Código desordenado? ¿Problemas de seguridad? Dilo directo. Señala problemas reales.",
			recommendations: "## Arregla Esto\n\n¿Qué necesita cambiar? Sé específico. No digas 'mejora el manejo de errores' - di DÓNDE y CÓMO.",
			perception:      "## Nivel\n\n¿Junior, mid o senior? Basado en lo que viste en el código. ¿Qué puede manejar? ¿Qué no puede?",
			closing:         "",
		},
		"fr": {
			instruction:     "Répondez en Français.",
			role:            "Lisez le code. Dites ce qui est bon, ce qui est mauvais. Sans détour. Sans parler corporate. Juste des faits. CRITIQUE: Utilisez SEULEMENT ## (deux hashtags) pour les titres. JAMAIS # seul.",
			structure:       "",
			technical:       "## Qualité du Code\n\nQuelle est la qualité réelle? Regardez la gestion des erreurs, la validation, l'architecture. C'est propre ou désordonné? Gère-t-il les cas limites? Soyez spécifique - pointez les vrais problèmes que vous voyez dans le code.",
			strengths:       "## Ce Qui Marche\n\nQu'est-ce qui est vraiment bon dans le code? Soyez spécifique. Ne dites pas des trucs génériques comme 'bonne structure' - dites CE QUI est bon et POURQUOI.",
			weaknesses:      "## Ce Qui Est Cassé\n\nQu'est-ce qui est mauvais? Manque de gestion d'erreurs? Pas de validation? Code désordonné? Problèmes de sécurité? Dites-le directement. Pointez les vrais problèmes.",
			recommendations: "## Réparez Ça\n\nQu'est-ce qui doit changer? Soyez spécifique. Ne dites pas 'améliorez la gestion des erreurs' - dites OÙ et COMMENT.",
			perception:      "## Niveau\n\nJunior, mid ou senior? Basé sur ce que vous avez vu dans le code. Qu'est-ce qu'il peut gérer? Qu'est-ce qu'il ne peut pas?",
			closing:         "",
		},
		"de": {
			instruction:     "Antworten Sie auf Deutsch.",
			role:            "Lesen Sie den Code. Sagen Sie, was gut ist, was schlecht ist. Ohne Umschweife. Ohne Unternehmenssprache. Nur Fakten. KRITISCH: Verwenden Sie NUR ## (zwei Hashtags) für Titel. NIEMALS # allein.",
			structure:       "",
			technical:       "## Code-Qualität\n\nWie ist die tatsächliche Qualität? Schauen Sie sich Fehlerbehandlung, Validierung, Architektur an. Ist es sauber oder unordentlich? Behandelt es Grenzfälle? Seien Sie spezifisch - zeigen Sie echte Probleme, die Sie im Code sehen.",
			strengths:       "## Was Funktioniert\n\nWas ist wirklich gut im Code? Seien Sie spezifisch. Sagen Sie nicht generische Dinge wie 'gute Struktur' - sagen Sie WAS gut ist und WARUM.",
			weaknesses:      "## Was Kaputt Ist\n\nWas ist schlecht? Fehlende Fehlerbehandlung? Keine Validierung? Unordentlicher Code? Sicherheitsprobleme? Sagen Sie es direkt. Zeigen Sie echte Probleme.",
			recommendations: "## Reparieren Sie Das\n\nWas muss sich ändern? Seien Sie spezifisch. Sagen Sie nicht 'verbessern Sie die Fehlerbehandlung' - sagen Sie WO und WIE.",
			perception:      "## Niveau\n\nJunior, mid oder senior? Basierend auf dem, was Sie im Code gesehen haben. Was kann er handhaben? Was nicht?",
			closing:         "",
		},
		"ja": {
			instruction:     "日本語で回答してください。",
			role:            "コードを読んでください。何が良くて何が悪いか言ってください。回りくどくなく。企業的な話し方なく。事実だけ。重要：タイトルには##（ハッシュタグ2つ）のみを使用。#単独は絶対に使わない。",
			structure:       "",
			technical:       "## コード品質\n\n実際の品質は？エラーハンドリング、検証、アーキテクチャを見てください。きれいですか、それとも乱雑ですか？エッジケースを処理していますか？具体的に - コードで見た実際の問題を指摘してください。",
			strengths:       "## 機能するもの\n\nコードで本当に良いものは何ですか？具体的に。「良い構造」のような一般的なことを言わないでください - 何が良くて、なぜかを言ってください。",
			weaknesses:      "## 壊れているもの\n\n何が悪いですか？エラーハンドリングがない？検証がない？乱雑なコード？セキュリティ問題？直接言ってください。実際の問題を指摘してください。",
			recommendations: "## これを修正\n\n何を変える必要がありますか？具体的に。「エラーハンドリングを改善」と言わないでください - どこでどのようにかを言ってください。",
			perception:      "## レベル\n\nジュニア、ミッド、またはシニア？コードで見たものに基づいて。何を扱えますか？何を扱えませんか？",
			closing:         "",
		},
		"zh": {
			instruction:     "请用中文回答。",
			role:            "阅读代码。说什么好，什么不好。不要绕圈子。不要企业话术。只要事实。关键：标题只使用##（两个井号）。绝不单独使用#。",
			structure:       "",
			technical:       "## 代码质量\n\n实际质量如何？看错误处理、验证、架构。干净还是混乱？处理边界情况吗？具体说明 - 指出您在代码中看到的实际问题。",
			strengths:       "## 有效的东西\n\n代码中真正好的是什么？具体说明。不要说\"良好的结构\"之类的泛泛之谈 - 说什么好以及为什么。",
			weaknesses:      "## 坏掉的东西\n\n什么不好？缺少错误处理？没有验证？混乱的代码？安全问题？直说。指出实际问题。",
			recommendations: "## 修复这个\n\n需要改变什么？具体说明。不要说\"改进错误处理\" - 说在哪里以及如何。",
			perception:      "## 水平\n\n初级、中级还是高级？基于您在代码中看到的。他能处理什么？不能处理什么？",
			closing:         "",
		},
	}

	lang := languageMap[language]
	if language == "" || lang.instruction == "" {
		lang = languageMap["en"]
	}

	fmt.Fprintf(&sb, "%s\n\n%s\n\n", lang.instruction, lang.role)

	fmt.Fprintf(&sb, "PROFILE: %s (@%s)\n", profile.Name, profile.Username)
	if profile.Bio != "" {
		fmt.Fprintf(&sb, "Bio: %s\n", profile.Bio)
	}
	fmt.Fprintf(&sb, "\n")

	sb.WriteString("CODE TO REVIEW:\n")
	repoCount := 0
	for _, repo := range profile.Repositories {
		if repoCount >= maxRepos {
			break
		}
		repoCount++

		fmt.Fprintf(&sb, "\n=== REPO %d: %s ===\n", repoCount, repo.Name)
		if repo.Language != "" {
			fmt.Fprintf(&sb, "Language: %s\n", repo.Language)
		}
		if repo.Description != "" {
			fmt.Fprintf(&sb, "Desc: %s\n", repo.Description)
		}

		if len(repo.CodeFiles) > 0 {
			fileCount := 0
			for path, content := range repo.CodeFiles {
				if fileCount >= 4 {
					break
				}
				fileCount++
				if len(content) > 1000 {
					content = content[:1000] + "..."
				}
				fmt.Fprintf(&sb, "\n%s:\n%s\n", path, content)
			}
		} else if repo.ReadmeContent != "" {
			readme := repo.ReadmeContent
			if len(readme) > 400 {
				readme = readme[:400] + "..."
			}
			fmt.Fprintf(&sb, "\nREADME: %s\n", readme)
		}
	}

	fmt.Fprintf(&sb, "\n%s\n%s\n%s\n%s\n%s", lang.technical, lang.strengths, lang.weaknesses, lang.recommendations, lang.perception)

	return sb.String()
}

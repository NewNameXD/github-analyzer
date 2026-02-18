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
		client: &http.Client{Timeout: 60 * time.Second},
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
			role:            "You're reviewing a GitHub profile. Talk like a real person - casual but knowledgeable. Skip the corporate buzzwords and AI-sounding phrases. Just give honest technical feedback like you would to a friend. CRITICAL: Use ONLY ## (two hashtags) for titles. NEVER use # alone. NEVER add # between sections.",
			structure:       "Break it down like this:",
			technical:       "## Technical Chops\n\nWhat's this person actually good at building? Look at the projects, tech choices, and how they organize code. Talk about the languages they know and what kind of stuff they build.",
			strengths:       "## The Good Stuff\n\nCall out 2-3 things that are legitimately impressive. Use specific examples from their repos. Be genuine.",
			weaknesses:      "## Where They Can Level Up\n\nWhat's missing or could be better? Be direct but constructive. Focus on stuff that actually matters.",
			recommendations: "## If I Were Them\n\nGive 3-4 specific things to work on next. Prioritize what'll make the biggest difference in their career. Be practical.",
			perception:      "## Real Talk\n\nWould you hire them? For what kind of role? What would they be great at? Be honest.",
			closing:         "Keep it real and direct. No fluff or generic advice. DO NOT add # alone anywhere.",
		},
		"pt": {
			instruction:     "Responda em Português Brasileiro.",
			role:            "Você está avaliando um perfil do GitHub. Fale como uma pessoa real - casual mas conhecedor. Sem jargões corporativos ou frases que soam como IA. Só dê feedback técnico honesto como você daria para um amigo. CRÍTICO: Use APENAS ## (dois hashtags) para títulos. NUNCA use # sozinho. NUNCA adicione # entre seções.",
			structure:       "Divida assim:",
			technical:       "## Bagagem Técnica\n\nNo que essa pessoa é boa de verdade? Analisa os projetos, escolhas de tech e como organiza o código. Fala sobre as linguagens que ela domina e o tipo de coisa que ela constrói.",
			strengths:       "## O Que Tá Bom\n\nDestaca 2-3 coisas que são realmente impressionantes. Usa exemplos específicos dos repos dela. Seja genuíno.",
			weaknesses:      "## Onde Pode Melhorar\n\nO que tá faltando ou poderia ser melhor? Seja direto mas construtivo. Foca em coisas que realmente importam.",
			recommendations: "## Se Fosse Eu\n\nDá 3-4 coisas específicas pra trabalhar agora. Prioriza o que vai fazer mais diferença na carreira. Seja prático.",
			perception:      "## Papo Reto\n\nVocê contrataria? Pra que tipo de vaga? O que essa pessoa mandaria bem? Seja honesto.",
			closing:         "Seja real e direto. Sem enrolação ou conselhos genéricos. NÃO adicione # sozinho em nenhum lugar.",
		},
		"es": {
			instruction:     "Responde en Español.",
			role:            "Estás evaluando un perfil de GitHub. Habla como una persona real - casual pero conocedor. Sin jerga corporativa o frases que suenan a IA. Solo da feedback técnico honesto como lo harías con un amigo. CRÍTICO: Usa SOLO ## (dos hashtags) para títulos. NUNCA uses # solo. NUNCA agregues # entre secciones.",
			structure:       "Divídelo así:",
			technical:       "## Habilidades Técnicas Reales\n\n¿En qué es realmente buena esta persona? Mira los proyectos, elecciones técnicas y cómo organiza el código. Habla de los lenguajes que domina y qué tipo de cosas construye.",
			strengths:       "## Lo Bueno\n\nDestaca 2-3 cosas que son genuinamente impresionantes. Usa ejemplos específicos de sus repos. Sé genuino.",
			weaknesses:      "## Donde Puede Mejorar\n\n¿Qué falta o podría ser mejor? Sé directo pero constructivo. Enfócate en cosas que realmente importan.",
			recommendations: "## Si Fuera Yo\n\nDa 3-4 cosas específicas en las que trabajar ahora. Prioriza lo que hará la mayor diferencia en su carrera. Sé práctico.",
			perception:      "## Hablando Claro\n\n¿Lo contratarías? ¿Para qué tipo de puesto? ¿En qué sería excelente? Sé honesto.",
			closing:         "Sé real y directo. Sin relleno o consejos genéricos. NO agregues # solo en ningún lugar.",
		},
		"fr": {
			instruction:     "Répondez en Français.",
			role:            "Tu évalues un profil GitHub. Parle comme une vraie personne - décontracté mais compétent. Sans jargon corporate ou phrases qui sonnent IA. Donne juste un feedback technique honnête comme tu le ferais avec un ami. CRITIQUE: Utilise SEULEMENT ## (deux hashtags) pour les titres. JAMAIS # seul. JAMAIS ajouter # entre les sections.",
			structure:       "Découpe comme ça:",
			technical:       "## Compétences Techniques Réelles\n\nDans quoi cette personne est vraiment bonne? Regarde les projets, choix tech et comment le code est organisé. Parle des langages qu'elle maîtrise et du type de trucs qu'elle construit.",
			strengths:       "## Ce Qui Est Bien\n\nSouligne 2-3 choses qui sont vraiment impressionnantes. Utilise des exemples spécifiques de leurs repos. Sois authentique.",
			weaknesses:      "## Où Progresser\n\nQu'est-ce qui manque ou pourrait être mieux? Sois direct mais constructif. Concentre-toi sur ce qui compte vraiment.",
			recommendations: "## Si C'était Moi\n\nDonne 3-4 trucs spécifiques sur lesquels bosser maintenant. Priorise ce qui fera la plus grande différence dans sa carrière. Sois pratique.",
			perception:      "## Franchement\n\nTu l'embaucherais? Pour quel genre de poste? Dans quoi serait-elle excellente? Sois honnête.",
			closing:         "Reste vrai et direct. Sans blabla ou conseils génériques. NE PAS ajouter # seul nulle part.",
		},
		"de": {
			instruction:     "Antworten Sie auf Deutsch.",
			role:            "Du bewertest ein GitHub-Profil. Rede wie ein echter Mensch - locker aber kompetent. Ohne Unternehmens-Jargon oder KI-klingende Phrasen. Gib einfach ehrliches technisches Feedback wie du es einem Freund geben würdest. KRITISCH: Verwende NUR ## (zwei Hashtags) für Titel. NIEMALS # allein. NIEMALS # zwischen Abschnitten hinzufügen.",
			structure:       "Teile es so auf:",
			technical:       "## Echte Technische Fähigkeiten\n\nWorin ist diese Person wirklich gut? Schau dir die Projekte, Tech-Entscheidungen und Code-Organisation an. Sprich über die Sprachen, die sie beherrscht, und was für Sachen sie baut.",
			strengths:       "## Das Gute Zeug\n\nHebe 2-3 Dinge hervor, die wirklich beeindruckend sind. Nutze spezifische Beispiele aus ihren Repos. Sei authentisch.",
			weaknesses:      "## Wo Verbesserung Möglich Ist\n\nWas fehlt oder könnte besser sein? Sei direkt aber konstruktiv. Konzentriere dich auf Dinge, die wirklich wichtig sind.",
			recommendations: "## Wenn Ich Es Wäre\n\nGib 3-4 spezifische Dinge, an denen jetzt gearbeitet werden sollte. Priorisiere, was den größten Unterschied in der Karriere machen wird. Sei praktisch.",
			perception:      "## Ehrlich Gesagt\n\nWürdest du sie einstellen? Für welche Art von Rolle? Worin wäre sie großartig? Sei ehrlich.",
			closing:         "Bleib echt und direkt. Kein Füllmaterial oder generische Ratschläge. Füge NIRGENDWO # allein hinzu.",
		},
		"ja": {
			instruction:     "日本語で回答してください。",
			role:            "GitHubプロフィールを評価しています。本物の人間のように話してください - カジュアルだけど知識がある。企業用語やAIっぽいフレーズは抜きで。友達に言うような正直な技術的フィードバックをください。重要：タイトルには##（ハッシュタグ2つ）のみを使用。#単独は絶対に使わない。セクション間に#を追加しない。",
			structure:       "こんな感じで分けて:",
			technical:       "## 実際の技術力\n\nこの人は何が本当に得意？プロジェクト、技術選択、コード構成を見て。得意な言語と作っているものの種類について話して。",
			strengths:       "## 良いところ\n\n本当に印象的な2-3のことを挙げて。リポジトリから具体的な例を使って。本物であること。",
			weaknesses:      "## 改善できるところ\n\n何が足りないか、何が良くなるか？率直に、でも建設的に。本当に重要なことに焦点を当てて。",
			recommendations: "## 自分だったら\n\n今取り組むべき3-4の具体的なこと。キャリアで最大の違いを生むことを優先して。実用的に。",
			perception:      "## 本音で\n\n採用する？どんな役割で？何が得意？正直に。",
			closing:         "リアルに、率直に。余計なことや一般的なアドバイスなしで。どこにも#単独を追加しない。",
		},
		"zh": {
			instruction:     "请用中文回答。",
			role:            "你在评估一个GitHub个人资料。像真人一样说话 - 随意但有见识。不要用企业术语或听起来像AI的短语。就像对朋友一样给出诚实的技术反馈。关键：只使用##（两个井号）作为标题。绝不单独使用#。绝不在章节之间添加#。",
			structure:       "这样分解:",
			technical:       "## 真实技术能力\n\n这个人真正擅长什么？看看项目、技术选择和代码组织。谈谈他们掌握的语言和构建的东西类型。",
			strengths:       "## 好的方面\n\n指出2-3个真正令人印象深刻的东西。用他们仓库的具体例子。要真诚。",
			weaknesses:      "## 可以改进的地方\n\n缺少什么或可以更好？直接但建设性。专注于真正重要的事情。",
			recommendations: "## 如果是我\n\n给出3-4个现在要做的具体事情。优先考虑对职业生涯影响最大的。要实用。",
			perception:      "## 实话实说\n\n你会雇用他们吗？担任什么角色？他们会擅长什么？诚实点。",
			closing:         "保持真实和直接。不要废话或通用建议。任何地方都不要单独添加#。",
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
	if profile.Location != "" || profile.Company != "" {
		fmt.Fprintf(&sb, "Location: %s | Company: %s\n", profile.Location, profile.Company)
	}
	fmt.Fprintf(&sb, "Followers: %d | Following: %d | Repos: %d\n\n",
		profile.Followers, profile.Following, profile.PublicRepos)

	sb.WriteString("REPOSITORIES:\n")
	repoCount := 0
	for _, repo := range profile.Repositories {
		if repoCount >= maxRepos {
			break
		}
		repoCount++

		fmt.Fprintf(&sb, "\n%d. %s", repoCount, repo.Name)
		if repo.Language != "" {
			fmt.Fprintf(&sb, " [%s]", repo.Language)
		}
		fmt.Fprintf(&sb, " - ⭐ %d | 🔀 %d\n", repo.Stars, repo.Forks)

		if repo.Description != "" {
			fmt.Fprintf(&sb, "   Description: %s\n", repo.Description)
		}

		if len(repo.Topics) > 0 {
			fmt.Fprintf(&sb, "   Topics: %v\n", repo.Topics)
		}

		if len(repo.Structure) > 0 {
			fmt.Fprintf(&sb, "   Structure: %v\n", repo.Structure)
		}

		if repo.ReadmeContent != "" {
			readmePreview := repo.ReadmeContent
			if len(readmePreview) > 400 {
				readmePreview = readmePreview[:400] + "..."
			}
			fmt.Fprintf(&sb, "   README: %s\n", readmePreview)
		}
	}

	fmt.Fprintf(&sb, "\n\n%s\n\n", lang.structure)
	fmt.Fprintf(&sb, "%s\n\n", lang.technical)
	fmt.Fprintf(&sb, "%s\n\n", lang.strengths)
	fmt.Fprintf(&sb, "%s\n\n", lang.weaknesses)
	fmt.Fprintf(&sb, "%s\n\n", lang.recommendations)
	fmt.Fprintf(&sb, "%s\n\n", lang.perception)
	fmt.Fprintf(&sb, "%s", lang.closing)

	return sb.String()
}

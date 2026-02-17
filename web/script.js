function setCookie(name, value, days = 365) {
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    const expires = "expires=" + date.toUTCString();
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
}

function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

const translations = {
    en: {
        title: 'AI-Powered GitHub Analysis',
        subtitle: 'Get honest, data-driven insights about any GitHub profile',
        usernamePlaceholder: 'Enter GitHub username',
        analyzeButton: 'Analyze Profile',
        languageLabel: 'Evaluation Language:',
        supportTitle: 'Enjoying this tool?',
        supportText: 'This project is open source and free to use. If you find it helpful, consider supporting:',
        followButton: 'Follow @pqpcara',
        starButton: 'Star this project',
        supportFooter: 'Your support helps keep this project alive and free for everyone!',
        loadingText: 'Analyzing profile...',
        footerPowered: 'Powered by Groq AI | Made with',
        footerBy: 'by',
        errorUsername: 'Please enter a GitHub username',
        followers: 'followers',
        following: 'following',
        repos: 'repos',
        organizations: 'Organizations',
        highlights: 'Highlights',
        totalStars: 'Total Stars',
        totalForks: 'Total Forks',
        repositories: 'Repositories',
        aiEvaluation: 'AI Evaluation'
    },
    pt: {
        title: 'Análise GitHub com IA',
        subtitle: 'Obtenha insights honestos e baseados em dados sobre qualquer perfil GitHub',
        usernamePlaceholder: 'Digite o nome de usuário do GitHub',
        analyzeButton: 'Analisar Perfil',
        languageLabel: 'Idioma da Avaliação:',
        supportTitle: 'Gostando desta ferramenta?',
        supportText: 'Este projeto é open source e gratuito. Se achar útil, considere apoiar:',
        followButton: 'Seguir @pqpcara',
        starButton: 'Dar estrela no projeto',
        supportFooter: 'Seu apoio ajuda a manter este projeto vivo e gratuito para todos!',
        loadingText: 'Analisando perfil...',
        footerPowered: 'Powered by Groq AI | Feito com',
        footerBy: 'por',
        errorUsername: 'Por favor, digite um nome de usuário do GitHub',
        followers: 'seguidores',
        following: 'seguindo',
        repos: 'repositórios',
        organizations: 'Organizações',
        highlights: 'Destaques',
        totalStars: 'Total de Estrelas',
        totalForks: 'Total de Forks',
        repositories: 'Repositórios',
        aiEvaluation: 'Avaliação da IA'
    },
    es: {
        title: 'Análisis GitHub con IA',
        subtitle: 'Obtén información honesta y basada en datos sobre cualquier perfil de GitHub',
        usernamePlaceholder: 'Ingrese nombre de usuario de GitHub',
        analyzeButton: 'Analizar Perfil',
        languageLabel: 'Idioma de Evaluación:',
        supportTitle: '¿Te gusta esta herramienta?',
        supportText: 'Este proyecto es de código abierto y gratuito. Si te resulta útil, considera apoyar:',
        followButton: 'Seguir @pqpcara',
        starButton: 'Dar estrella al proyecto',
        supportFooter: '¡Tu apoyo ayuda a mantener este proyecto vivo y gratuito para todos!',
        loadingText: 'Analizando perfil...',
        footerPowered: 'Powered by Groq AI | Hecho con',
        footerBy: 'por',
        errorUsername: 'Por favor, ingrese un nombre de usuario de GitHub',
        followers: 'seguidores',
        following: 'siguiendo',
        repos: 'repositorios',
        organizations: 'Organizaciones',
        highlights: 'Destacados',
        totalStars: 'Total de Estrellas',
        totalForks: 'Total de Forks',
        repositories: 'Repositorios',
        aiEvaluation: 'Evaluación de IA'
    },
    fr: {
        title: 'Analyse GitHub avec IA',
        subtitle: 'Obtenez des informations honnêtes et basées sur les données sur n\'importe quel profil GitHub',
        usernamePlaceholder: 'Entrez le nom d\'utilisateur GitHub',
        analyzeButton: 'Analyser le Profil',
        languageLabel: 'Langue d\'Évaluation:',
        supportTitle: 'Vous aimez cet outil?',
        supportText: 'Ce projet est open source et gratuit. Si vous le trouvez utile, envisagez de soutenir:',
        followButton: 'Suivre @pqpcara',
        starButton: 'Mettre une étoile au projet',
        supportFooter: 'Votre soutien aide à garder ce projet vivant et gratuit pour tous!',
        loadingText: 'Analyse du profil...',
        footerPowered: 'Powered by Groq AI | Fait avec',
        footerBy: 'par',
        errorUsername: 'Veuillez entrer un nom d\'utilisateur GitHub',
        followers: 'abonnés',
        following: 'abonnements',
        repos: 'dépôts',
        organizations: 'Organisations',
        highlights: 'Points forts',
        totalStars: 'Total d\'Étoiles',
        totalForks: 'Total de Forks',
        repositories: 'Dépôts',
        aiEvaluation: 'Évaluation IA'
    },
    de: {
        title: 'KI-gestützte GitHub-Analyse',
        subtitle: 'Erhalten Sie ehrliche, datengestützte Einblicke in jedes GitHub-Profil',
        usernamePlaceholder: 'GitHub-Benutzername eingeben',
        analyzeButton: 'Profil Analysieren',
        languageLabel: 'Bewertungssprache:',
        supportTitle: 'Gefällt Ihnen dieses Tool?',
        supportText: 'Dieses Projekt ist Open Source und kostenlos. Wenn Sie es nützlich finden, erwägen Sie Unterstützung:',
        followButton: '@pqpcara folgen',
        starButton: 'Projekt mit Stern markieren',
        supportFooter: 'Ihre Unterstützung hilft, dieses Projekt am Leben und kostenlos für alle zu halten!',
        loadingText: 'Profil analysieren...',
        footerPowered: 'Powered by Groq AI | Gemacht mit',
        footerBy: 'von',
        errorUsername: 'Bitte geben Sie einen GitHub-Benutzernamen ein',
        followers: 'Follower',
        following: 'folgt',
        repos: 'Repos',
        organizations: 'Organisationen',
        highlights: 'Highlights',
        totalStars: 'Gesamt Sterne',
        totalForks: 'Gesamt Forks',
        repositories: 'Repositories',
        aiEvaluation: 'KI-Bewertung'
    },
    ja: {
        title: 'AI搭載GitHub分析',
        subtitle: '任意のGitHubプロフィールについて正直でデータ駆動型の洞察を取得',
        usernamePlaceholder: 'GitHubユーザー名を入力',
        analyzeButton: 'プロフィールを分析',
        languageLabel: '評価言語:',
        supportTitle: 'このツールを気に入りましたか？',
        supportText: 'このプロジェクトはオープンソースで無料です。役立つと思ったら、サポートをご検討ください:',
        followButton: '@pqpcaraをフォロー',
        starButton: 'プロジェクトにスター',
        supportFooter: 'あなたのサポートがこのプロジェクトを維持し、誰でも無料で使えるようにします！',
        loadingText: 'プロフィールを分析中...',
        footerPowered: 'Powered by Groq AI | 作成者',
        footerBy: '',
        errorUsername: 'GitHubユーザー名を入力してください',
        followers: 'フォロワー',
        following: 'フォロー中',
        repos: 'リポジトリ',
        organizations: '組織',
        highlights: 'ハイライト',
        totalStars: '合計スター',
        totalForks: '合計フォーク',
        repositories: 'リポジトリ',
        aiEvaluation: 'AI評価'
    },
    zh: {
        title: 'AI驱动的GitHub分析',
        subtitle: '获取关于任何GitHub个人资料的诚实、数据驱动的见解',
        usernamePlaceholder: '输入GitHub用户名',
        analyzeButton: '分析个人资料',
        languageLabel: '评估语言:',
        supportTitle: '喜欢这个工具吗？',
        supportText: '这个项目是开源且免费的。如果您觉得有用，请考虑支持:',
        followButton: '关注 @pqpcara',
        starButton: '给项目加星',
        supportFooter: '您的支持帮助保持这个项目活跃并对所有人免费！',
        loadingText: '正在分析个人资料...',
        footerPowered: 'Powered by Groq AI | 制作者',
        footerBy: '',
        errorUsername: '请输入GitHub用户名',
        followers: '关注者',
        following: '正在关注',
        repos: '仓库',
        organizations: '组织',
        highlights: '亮点',
        totalStars: '总星标',
        totalForks: '总分支',
        repositories: '仓库',
        aiEvaluation: 'AI评估'
    }
};

const usernameInput = document.getElementById('username');
const languageSelect = document.getElementById('language');
const analyzeBtn = document.getElementById('analyzeBtn');
const loading = document.getElementById('loading');
const error = document.getElementById('error');
const errorText = document.getElementById('errorText');
const result = document.getElementById('result');
const profileHeader = document.getElementById('profileHeader');
const evaluation = document.getElementById('evaluation');

function updateLanguage(lang) {
    const t = translations[lang] || translations.en;
    
    document.querySelectorAll('[data-i18n]').forEach(element => {
        const key = element.getAttribute('data-i18n');
        if (t[key]) {
            element.textContent = t[key];
        }
    });
    
    document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
        const key = element.getAttribute('data-i18n-placeholder');
        if (t[key]) {
            element.placeholder = t[key];
        }
    });
}

function typeWriter(element, text, speed = 5) {
    return new Promise((resolve) => {
        let i = 0;
        element.innerHTML = '';
        element.classList.add('typing-cursor');
        
        function type() {
            if (i < text.length) {
                element.innerHTML += text.charAt(i);
                i++;
                setTimeout(type, speed);
            } else {
                element.classList.remove('typing-cursor');
                resolve();
            }
        }
        
        type();
    });
}

async function typeWriterHTML(element, html, speed = 1) {
    return new Promise((resolve) => {
        loading.classList.add('hidden');
        
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = html;
        
        let charIndex = 0;
        element.innerHTML = '';
        element.classList.add('typing-cursor');
        
        const htmlChunks = [];
        
        function extractText(node, chunks = []) {
            if (node.nodeType === Node.TEXT_NODE) {
                const text = node.textContent;
                for (let i = 0; i < text.length; i++) {
                    chunks.push({ char: text[i], html: text[i] });
                }
            } else if (node.nodeType === Node.ELEMENT_NODE) {
                const tagName = node.tagName.toLowerCase();
                const attrs = Array.from(node.attributes)
                    .map(a => `${a.name}="${a.value}"`)
                    .join(' ');
                
                chunks.push({ 
                    char: '', 
                    html: `<${tagName}${attrs ? ' ' + attrs : ''}>`,
                    isTag: true 
                });
                
                for (let child of node.childNodes) {
                    extractText(child, chunks);
                }
                
                chunks.push({ 
                    char: '', 
                    html: `</${tagName}>`,
                    isTag: true 
                });
            }
            return chunks;
        }
        
        const chunks = extractText(tempDiv);
        let currentHTML = '';
        let chunkIndex = 0;
        
        function type() {
            if (chunkIndex < chunks.length) {
                const chunk = chunks[chunkIndex];
                currentHTML += chunk.html;
                element.innerHTML = currentHTML;
                
                chunkIndex++;
                
                if (chunk.isTag) {
                    setTimeout(type, 0);
                } else {
                    setTimeout(type, speed);
                }
            } else {
                element.classList.remove('typing-cursor');
                resolve();
            }
        }
        
        type();
    });
}

async function evaluate() {
    const username = usernameInput.value.trim();
    const language = languageSelect.value;
    const t = translations[language] || translations.en;
    
    if (!username) {
        showError(t.errorUsername);
        return;
    }

    hideAll();
    loading.classList.remove('hidden');
    analyzeBtn.disabled = true;

    try {
        const apiUrl = window.API_URL || '/api/evaluate';
        const response = await fetch(apiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                username: username,
                language: language 
            }),
        });

        const data = await response.json();

        if (!data.success) {
            showError(data.error || 'Failed to evaluate profile');
            return;
        }

        await displayResult(data);
    } catch (err) {
        showError('Network error: ' + err.message);
    } finally {
        loading.classList.add('hidden');
        analyzeBtn.disabled = false;
    }
}

async function displayResult(data) {
    const profile = data.profile;
    
    const orgsHTML = profile.organizations && profile.organizations.length > 0 ? `
        <div class="mt-4 pt-4 border-t border-dark-border">
            <h4 class="text-sm font-semibold text-gray-400 mb-3 flex items-center gap-2">
                <i class="fas fa-building"></i> Organizations
            </h4>
            <div class="flex flex-wrap gap-2">
                ${profile.organizations.slice(0, 5).map(org => `
                    <a href="https://github.com/${org.login}" target="_blank" 
                       class="flex items-center gap-2 px-3 py-1.5 bg-dark-card border border-dark-border hover:border-cyan-500/30 rounded-lg transition-all group"
                       title="${org.description || org.login}">
                        <img src="${org.avatar_url}" alt="${org.login}" class="w-5 h-5 rounded">
                        <span class="text-sm group-hover:text-cyan-400 transition-colors">${org.login}</span>
                    </a>
                `).join('')}
            </div>
        </div>
    ` : '';

    const statsHTML = `
        <div class="mt-4 pt-4 border-t border-dark-border">
            <h4 class="text-sm font-semibold text-gray-400 mb-3 flex items-center gap-2">
                <i class="fas fa-trophy"></i> Stats
            </h4>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
                <div class="bg-dark-card border border-dark-border rounded-lg p-3 text-center">
                    <div class="text-2xl font-bold text-cyan-400">${profile.total_stars || 0}</div>
                    <div class="text-xs text-gray-500 mt-1">Total Stars</div>
                </div>
                <div class="bg-dark-card border border-dark-border rounded-lg p-3 text-center">
                    <div class="text-2xl font-bold text-cyan-400">${profile.total_forks || 0}</div>
                    <div class="text-xs text-gray-500 mt-1">Total Forks</div>
                </div>
                <div class="bg-dark-card border border-dark-border rounded-lg p-3 text-center">
                    <div class="text-2xl font-bold text-cyan-400">${profile.public_repos}</div>
                    <div class="text-xs text-gray-500 mt-1">Repositories</div>
                </div>
            </div>
        </div>
    `;
    
    profileHeader.innerHTML = `
        <div class="flex flex-col sm:flex-row items-center sm:items-start gap-6 w-full">
            <div class="relative group">
                <img 
                    src="${profile.avatar_url}" 
                    alt="${profile.name}" 
                    class="w-32 h-32 rounded-2xl border-4 border-cyan-500/30 shadow-xl glow-cyan"
                />
                <div class="absolute inset-0 bg-gradient-to-t from-cyan-500/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            </div>
            <div class="flex-1 text-center sm:text-left w-full">
                <div class="flex flex-col sm:flex-row sm:items-center gap-3 mb-3">
                    <h2 class="text-3xl font-bold">${profile.name || profile.username}</h2>
                    <a href="https://github.com/${profile.username}" target="_blank" 
                       class="inline-flex items-center gap-2 px-4 py-2 bg-dark-card border border-dark-border hover:border-cyan-500/30 rounded-lg text-sm transition-all group">
                        <i class="fab fa-github group-hover:text-cyan-400 transition-colors"></i>
                        <span class="group-hover:text-cyan-400 transition-colors">@${profile.username}</span>
                        <i class="fas fa-external-link-alt text-xs opacity-50 group-hover:opacity-100 transition-opacity"></i>
                    </a>
                </div>
                
                ${profile.bio ? `<p class="text-gray-400 mb-3 text-lg">${profile.bio}</p>` : ''}
                
                <div class="flex flex-wrap gap-3 text-sm text-gray-500 justify-center sm:justify-start mb-4">
                    ${profile.location ? `
                        <span class="flex items-center gap-1.5">
                            <i class="fas fa-map-marker-alt text-cyan-400"></i> ${profile.location}
                        </span>
                    ` : ''}
                    ${profile.company ? `
                        <span class="flex items-center gap-1.5">
                            <i class="fas fa-building text-cyan-400"></i> ${profile.company}
                        </span>
                    ` : ''}
                    ${profile.blog ? `
                        <a href="${profile.blog.startsWith('http') ? profile.blog : 'https://' + profile.blog}" target="_blank" 
                           class="flex items-center gap-1.5 hover:text-cyan-400 transition-colors">
                            <i class="fas fa-link text-cyan-400"></i> ${profile.blog.replace(/^https?:\/\//, '')}
                        </a>
                    ` : ''}
                    ${profile.twitter_username ? `
                        <a href="https://twitter.com/${profile.twitter_username}" target="_blank" 
                           class="flex items-center gap-1.5 hover:text-cyan-400 transition-colors">
                            <i class="fab fa-twitter text-cyan-400"></i> @${profile.twitter_username}
                        </a>
                    ` : ''}
                </div>
                
                <div class="flex flex-wrap gap-6 text-sm justify-center sm:justify-start">
                    <span class="flex items-center gap-2 text-gray-400">
                        <i class="fas fa-users text-cyan-400"></i>
                        <strong class="text-white">${profile.followers}</strong> followers
                    </span>
                    <span class="flex items-center gap-2 text-gray-400">
                        <i class="fas fa-user-friends text-cyan-400"></i>
                        <strong class="text-white">${profile.following}</strong> following
                    </span>
                </div>
                
                ${statsHTML}
                ${orgsHTML}
            </div>
        </div>
    `;

    const formattedEvaluation = data.evaluation
        .replace(/##\s+(.+)/g, '<h3 class="text-xl font-bold text-white mt-6 mb-3 first:mt-0 flex items-center gap-2"><i class="fas fa-chevron-right text-cyan-400 text-sm"></i>$1</h3>')
        .replace(/\*\*(.+?)\*\*/g, '<strong class="font-semibold text-cyan-400">$1</strong>')
        .replace(/\n\n/g, '</p><p class="mb-4">')
        .replace(/^(.+)$/gm, '<p class="mb-4 text-gray-400 leading-relaxed">$1</p>')
        .replace(/<p class="mb-4 text-gray-400 leading-relaxed"><h3/g, '<h3')
        .replace(/<\/h3><\/p>/g, '</h3>');

    result.classList.remove('hidden');
    
    result.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    
    await typeWriterHTML(evaluation, formattedEvaluation, 3);
}

function showError(message) {
    errorText.textContent = message;
    error.classList.remove('hidden');
    setTimeout(() => {
        error.classList.add('hidden');
    }, 5000);
}

function hideAll() {
    error.classList.add('hidden');
    result.classList.add('hidden');
}

analyzeBtn.addEventListener('click', evaluate);

usernameInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        evaluate();
    }
});

languageSelect.addEventListener('change', (e) => {
    const selectedLang = e.target.value;
    updateLanguage(selectedLang);
    setCookie('preferredLanguage', selectedLang, 365);
});

const savedLanguage = getCookie('preferredLanguage') || 'en';
if (savedLanguage !== 'en') {
    languageSelect.value = savedLanguage;
    updateLanguage(savedLanguage);
}

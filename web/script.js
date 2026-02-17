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
        title: 'GitHub Profile Evaluator',
        subtitle: 'Honest, data-driven analysis of GitHub profiles',
        usernamePlaceholder: 'Enter GitHub username',
        analyzeButton: 'Analyze Profile',
        languageLabel: 'Evaluation Language:',
        supportTitle: 'Enjoying this tool?',
        supportText: 'This project is open source and free to use. If you find it helpful, consider supporting:',
        followButton: 'Follow @pqpcara',
        starButton: 'Star this project',
        supportFooter: 'Your support helps keep this project alive and free for everyone!',
        loadingText: 'Analyzing profile and generating evaluation...',
        footerPowered: 'Powered by Groq AI | Made with',
        footerBy: 'by',
        errorUsername: 'Please enter a GitHub username'
    },
    pt: {
        title: 'Avaliador de Perfil GitHub',
        subtitle: 'Análise honesta e baseada em dados de perfis GitHub',
        usernamePlaceholder: 'Digite o nome de usuário do GitHub',
        analyzeButton: 'Analisar Perfil',
        languageLabel: 'Idioma da Avaliação:',
        supportTitle: 'Gostando desta ferramenta?',
        supportText: 'Este projeto é open source e gratuito. Se achar útil, considere apoiar:',
        followButton: 'Seguir @pqpcara',
        starButton: 'Dar estrela no projeto',
        supportFooter: 'Seu apoio ajuda a manter este projeto vivo e gratuito para todos!',
        loadingText: 'Analisando perfil e gerando avaliação...',
        footerPowered: 'Powered by Groq AI | Feito com',
        footerBy: 'por',
        errorUsername: 'Por favor, digite um nome de usuário do GitHub'
    },
    es: {
        title: 'Evaluador de Perfil GitHub',
        subtitle: 'Análisis honesto y basado en datos de perfiles GitHub',
        usernamePlaceholder: 'Ingrese nombre de usuario de GitHub',
        analyzeButton: 'Analizar Perfil',
        languageLabel: 'Idioma de Evaluación:',
        supportTitle: '¿Te gusta esta herramienta?',
        supportText: 'Este proyecto es de código abierto y gratuito. Si te resulta útil, considera apoyar:',
        followButton: 'Seguir @pqpcara',
        starButton: 'Dar estrella al proyecto',
        supportFooter: '¡Tu apoyo ayuda a mantener este proyecto vivo y gratuito para todos!',
        loadingText: 'Analizando perfil y generando evaluación...',
        footerPowered: 'Powered by Groq AI | Hecho con',
        footerBy: 'por',
        errorUsername: 'Por favor, ingrese un nombre de usuario de GitHub'
    },
    fr: {
        title: 'Évaluateur de Profil GitHub',
        subtitle: 'Analyse honnête et basée sur les données des profils GitHub',
        usernamePlaceholder: 'Entrez le nom d\'utilisateur GitHub',
        analyzeButton: 'Analyser le Profil',
        languageLabel: 'Langue d\'Évaluation:',
        supportTitle: 'Vous aimez cet outil?',
        supportText: 'Ce projet est open source et gratuit. Si vous le trouvez utile, envisagez de soutenir:',
        followButton: 'Suivre @pqpcara',
        starButton: 'Mettre une étoile au projet',
        supportFooter: 'Votre soutien aide à garder ce projet vivant et gratuit pour tous!',
        loadingText: 'Analyse du profil et génération de l\'évaluation...',
        footerPowered: 'Powered by Groq AI | Fait avec',
        footerBy: 'par',
        errorUsername: 'Veuillez entrer un nom d\'utilisateur GitHub'
    },
    de: {
        title: 'GitHub Profil Bewerter',
        subtitle: 'Ehrliche, datengestützte Analyse von GitHub-Profilen',
        usernamePlaceholder: 'GitHub-Benutzername eingeben',
        analyzeButton: 'Profil Analysieren',
        languageLabel: 'Bewertungssprache:',
        supportTitle: 'Gefällt Ihnen dieses Tool?',
        supportText: 'Dieses Projekt ist Open Source und kostenlos. Wenn Sie es nützlich finden, erwägen Sie Unterstützung:',
        followButton: '@pqpcara folgen',
        starButton: 'Projekt mit Stern markieren',
        supportFooter: 'Ihre Unterstützung hilft, dieses Projekt am Leben und kostenlos für alle zu halten!',
        loadingText: 'Profil analysieren und Bewertung generieren...',
        footerPowered: 'Powered by Groq AI | Gemacht mit',
        footerBy: 'von',
        errorUsername: 'Bitte geben Sie einen GitHub-Benutzernamen ein'
    },
    ja: {
        title: 'GitHub プロフィール評価ツール',
        subtitle: 'GitHubプロフィールの正直でデータ駆動型の分析',
        usernamePlaceholder: 'GitHubユーザー名を入力',
        analyzeButton: 'プロフィールを分析',
        languageLabel: '評価言語:',
        supportTitle: 'このツールを気に入りましたか？',
        supportText: 'このプロジェクトはオープンソースで無料です。役立つと思ったら、サポートをご検討ください:',
        followButton: '@pqpcaraをフォロー',
        starButton: 'プロジェクトにスター',
        supportFooter: 'あなたのサポートがこのプロジェクトを維持し、誰でも無料で使えるようにします！',
        loadingText: 'プロフィールを分析し、評価を生成中...',
        footerPowered: 'Powered by Groq AI | 作成者',
        footerBy: '',
        errorUsername: 'GitHubユーザー名を入力してください'
    },
    zh: {
        title: 'GitHub 个人资料评估器',
        subtitle: '基于数据的诚实GitHub个人资料分析',
        usernamePlaceholder: '输入GitHub用户名',
        analyzeButton: '分析个人资料',
        languageLabel: '评估语言:',
        supportTitle: '喜欢这个工具吗？',
        supportText: '这个项目是开源且免费的。如果您觉得有用，请考虑支持:',
        followButton: '关注 @pqpcara',
        starButton: '给项目加星',
        supportFooter: '您的支持帮助保持这个项目活跃并对所有人免费！',
        loadingText: '正在分析个人资料并生成评估...',
        footerPowered: 'Powered by Groq AI | 制作者',
        footerBy: '',
        errorUsername: '请输入GitHub用户名'
    }
};

const usernameInput = document.getElementById('username');
const languageSelect = document.getElementById('language');
const analyzeBtn = document.getElementById('analyzeBtn');
const loading = document.getElementById('loading');
const error = document.getElementById('error');
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

async function evaluate() {
    console.log('Evaluate function called');
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
        console.log('Fetching evaluation for:', username, 'in language:', language);
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

        console.log('Response status:', response.status);
        const data = await response.json();
        console.log('Response data:', data);

        if (!data.success) {
            showError(data.error || 'Failed to evaluate profile');
            return;
        }

        displayResult(data);
    } catch (err) {
        console.error('Error:', err);
        showError('Network error: ' + err.message);
    } finally {
        loading.classList.add('hidden');
        analyzeBtn.disabled = false;
    }
}

function displayResult(data) {
    console.log('Displaying result');
    const profile = data.profile;
    
    profileHeader.innerHTML = `
        <img 
            src="${profile.avatar_url}" 
            alt="${profile.name}" 
            class="w-20 h-20 rounded-full border-4 border-primary"
        />
        <div class="flex-1">
            <h2 class="text-2xl font-bold text-gray-800 mb-1">${profile.name || profile.username}</h2>
            <p class="text-gray-600 mb-1">${profile.bio || 'No bio available'}</p>
            <p class="text-gray-500 text-sm mb-2">
                ${profile.location || ''} 
                ${profile.company ? '| ' + profile.company : ''}
            </p>
            <div class="flex gap-5 text-sm">
                <span class="text-gray-600">
                    <strong class="text-primary">${profile.followers}</strong> followers
                </span>
                <span class="text-gray-600">
                    <strong class="text-primary">${profile.following}</strong> following
                </span>
                <span class="text-gray-600">
                    <strong class="text-primary">${profile.public_repos}</strong> repos
                </span>
            </div>
        </div>
    `;

    const formattedEvaluation = data.evaluation
        .replace(/##\s+(.+)/g, '<h3 class="text-xl font-bold text-gray-800 mt-6 mb-3 first:mt-0">$1</h3>')
        .replace(/\*\*(.+?)\*\*/g, '<strong class="font-semibold text-gray-900">$1</strong>')
        .replace(/\n\n/g, '</p><p class="mb-4">')
        .replace(/^(.+)$/gm, '<p class="mb-4 text-gray-700 leading-relaxed">$1</p>')
        .replace(/<p class="mb-4 text-gray-700 leading-relaxed"><h3/g, '<h3')
        .replace(/<\/h3><\/p>/g, '</h3>');

    evaluation.innerHTML = formattedEvaluation;
    result.classList.remove('hidden');
}

function showError(message) {
    console.log('Showing error:', message);
    error.textContent = message;
    error.classList.remove('hidden');
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
    updateLanguage(e.target.value);
});

console.log('GitHub Profile Evaluator loaded');

languageSelect.removeEventListener('change', () => {});
languageSelect.addEventListener('change', (e) => {
    const selectedLang = e.target.value;
    updateLanguage(selectedLang);
    setCookie('preferredLanguage', selectedLang, 365);
    console.log('Language saved to cookie:', selectedLang);
});

const savedLanguage = getCookie('preferredLanguage') || 'en';
if (savedLanguage !== 'en') {
    languageSelect.value = savedLanguage;
    updateLanguage(savedLanguage);
    console.log('Loaded saved language:', savedLanguage);
}

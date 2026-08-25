const SWEARS = new Set([
    // English
    'fuck', 'fucking', 'fucker', 'fucked',
    'cunt', 'cunts',
    'shit', 'shitty', 'bullshit',
    'asshole', 'arsehole',
    'bitch', 'bitches',
    'piss', 'pussy', 'pussies',
    'whore', 'slut', 'bastard',
    'motherfucker', 'motherfuckers', 'motherfucking',
    'twat', 'twats',
    'prick', 'pricks',
    'wanker', 'wankers',
    'bollocks',
    'bugger', 'buggers',
    'cocksucker', 'cocksuckers',
    'dickhead', 'dickheads',
    'asshat', 'asshats',
    'shithead', 'shitheads',
    'fucktard', 'fucktards',
    'pissed',
    'cum', 'jizz',

    // Russian
    'хуй', 'хуя', 'хуи', 'хуёв', 'хуйня',
    'пизда', 'пиздец', 'пизды',
    'ебать', 'ебал', 'ебаный', 'ёбаный', 'ебаная', 'ёбаная',
    'блядь', 'бля', 'бляди',
    'мудак', 'мудила',
    'сука', 'суки',
    'пидор', 'пидорас', 'пидоры',
    'заебал', 'заебала', 'заебали',
    'нахуй', 'похуй',
    'хуйло',
    'hui', 'khuy', 'huy',
    'pizda', 'pizdets', 'pizdec',
    'yebat', 'ebat', 'yobany', 'ebany',
    'blyad', 'blyat', 'blya',
    'mudak',
    'suka',
    'pidor', 'pidoras',
    'zaebal',
    'nahui', 'pohui',

    // Polish
    'kurwa', 'kurwy', 'kurwą',
    'chuj', 'chuja', 'chuje',
    'pizda', 'pizdy',
    'cipa', 'cipy',
    'jebać', 'jebany', 'jebana', 'jebane',
    'pierdolić', 'pierdolony', 'pierdolona',
    'dupa', 'dupy', 'dupek', 'dupki',
    'skurwysyn', 'skurwysyny',
    'gówno', 'gówna', 'gowno',
    'suka', 'suki',
    'spierdalaj',
]);

function containsSwear(value) {
    // 1. Normalize: lowercase + turn runs of separators into a single space
    const normalized = value
        .toLowerCase()
        .replace(/[\s'\-\.,]+/g, ' ')
        .trim();

    if (!normalized) return false;

    // 2. Split into tokens
    const tokens = normalized.split(' ');

    // 3. Exact token match (handles “fucking”, “cunt”, “fuck you”, etc.)
    for (const token of tokens) {
        if (SWEARS.has(token)) return true;
    }

    // 4. Optional light anti-obfuscation: check pairs of adjacent tokens
    //    (catches “fu ck”, “c unt”, “fu.ck” after normalization, etc.)
    for (let i = 0; i < tokens.length - 1; i++) {
        const joined = tokens[i] + tokens[i + 1];
        if (SWEARS.has(joined)) return true;
    }

    return false;
}

function isValidGenericField(value) {
    if (typeof value !== 'string' || value.length === 0) return false;

    // Existing bad-character / sinister / zero-width / control / private-use checks
    if (/[\p{C}\u2000-\u200F\u2028-\u202F\u205F-\u206F\u3000\uFEFF\u00AD]/u.test(value)) {
        return false;
    }

    // Zalgo / excessive marks
    if (/\p{M}{5,}/u.test(value)) return false;

    // Whitelist
    if (!/^[\p{L}\p{M}0-9\s'\-\.,]+$/u.test(value)) {
        return false;
    }

    // Swears
    if (containsSwear(value)) {
        return false;
    }

    return true;
}
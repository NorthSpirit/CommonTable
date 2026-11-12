const themes = [
    {key: 'default', id: 'themeCSS'},
    {key: 'alt', id: 'themeCSSAlt'},
    {key: 'altclean', id: 'themeCSSAltClean'}
];
const paletteMenuLink = document.getElementById('MColour');
const themeKey = 'GMSTheme';

function availableThemes() {
    return themes.filter(t => document.getElementById(t.id));
}

function applyTheme(key) {
    availableThemes().forEach(t => {
        document.getElementById(t.id).disabled = (t.key !== key);
    });
    localStorage.setItem(themeKey, key);
}

function applyThemeFromCache() {
    const cachedTheme = localStorage.getItem(themeKey);
    if (cachedTheme && availableThemes().some(t => t.key === cachedTheme)) {
        applyTheme(cachedTheme);
    }
}

document.addEventListener('DOMContentLoaded', applyThemeFromCache);

paletteMenuLink.addEventListener('click', (e) => {
    e.preventDefault();
    const available = availableThemes();
    // The theme that is on now is the one whose <link> is not disabled
    const currentIndex = available.findIndex(t => !document.getElementById(t.id).disabled);
    const next = available[(currentIndex + 1) % available.length];
    applyTheme(next.key);
});
const defaultCSS = document.getElementById('themeCSS');
const altCSS = document.getElementById('themeCSSAlt');
const paletteMenuLink = document.getElementById('MColour');
const themeKey = 'GMSTheme';

function applyTheme(isAlt) {
    if (isAlt) {
        defaultCSS.disabled = true;
        altCSS.disabled = false;
        localStorage.setItem(themeKey, 'alt');
    } else {
        defaultCSS.disabled = false;
        altCSS.disabled = true;
        localStorage.setItem(themeKey, 'default');
    }
}

function applyThemeFromCache() {
    const cachedTheme = localStorage.getItem(themeKey);
    if (cachedTheme === 'alt') {
        applyTheme(true);
    }
}

document.addEventListener('DOMContentLoaded', applyThemeFromCache);

paletteMenuLink.addEventListener('click', (e) => {
    e.preventDefault();
    const isCurrentlyDefault = altCSS.disabled;
    applyTheme(isCurrentlyDefault);
});
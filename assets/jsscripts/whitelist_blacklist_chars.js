function isValidGenericField(value) {
    if (typeof value !== 'string' || value.length === 0) return false;

    // Explicit bad characters / zero-widths / controls / private-use / odd spaces - sinisters
    if (/[\p{C}\u2000-\u200F\u2028-\u202F\u205F-\u206F\u3000\uFEFF\u00AD]/u.test(value)) {
        return false;
    }

    // Zalgo / excessive mark stacking
    if (/\p{M}{5,}/u.test(value)) return false;

    // Whitelist
    return /^[\p{L}\p{M}0-9\s'\-\.,]+$/u.test(value);
}
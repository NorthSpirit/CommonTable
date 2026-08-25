document.addEventListener('DOMContentLoaded', () => {
    const fields = document.querySelectorAll('[data-validate="generic"]');

    fields.forEach((field) => {
        const errorEl = document.createElement('p');
        errorEl.className = 'text-xs text-red-600 mt-1 hidden';
        errorEl.textContent = "That doesn't look right — please rephrase.";
        field.insertAdjacentElement('afterend', errorEl);

        const check = () => {
            const value = field.value.trim();
            // Empty is allowed here — required-ness is a separate concern
            // (the "required" attribute), this only judges content that IS present.
            const valid = value === '' || isValidGenericField(value);
            errorEl.classList.toggle('hidden', valid);
            field.setAttribute('aria-invalid', valid ? 'false' : 'true');
            return valid;
        };

        field.addEventListener('blur', check);
        field.addEventListener('input', () => {
            // Once an error is showing, re-check live so it clears the
            // moment the user fixes it, rather than waiting for another blur.
            if (!errorEl.classList.contains('hidden')) check();
        });

        const form = field.closest('form');
        if (form) {
            form.addEventListener('submit', (e) => {
                if (!check()) {
                    e.preventDefault();
                    field.focus();
                }
            });
        }
    });
});
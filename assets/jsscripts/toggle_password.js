const passwordField = document.getElementById('password-field');
const toggleButton = document.getElementById('toggle-password');

toggleButton.addEventListener('click', () => {
    const isHidden = passwordField.type === 'password';
    passwordField.type = isHidden ? 'text' : 'password';
    toggleButton.textContent = isHidden ? 'Hide' : 'Show';
});
// Wires up instant client-side preview for any <input type="file"> that has
// a matching data-preview-target attribute pointing at an <img> id.
// The preview is purely visual — nothing is uploaded until form submit.
document.querySelectorAll('input[type="file"][data-preview-target]').forEach((input) => {
    const targetId = input.getAttribute('data-preview-target');
    const img = document.getElementById(targetId);
    if (!img) return;

    input.addEventListener('change', () => {
        const file = input.files && input.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (event) => {
            img.src = event.target.result;
            img.classList.remove('hidden');
            const placeholder = document.getElementById(targetId + '-placeholder');
            if (placeholder) placeholder.classList.add('hidden');
        };
        reader.readAsDataURL(file);
    });
});
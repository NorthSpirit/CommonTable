function resizeChatBox() {
    const chatBox = document.getElementById('chat-messages');
    if (!chatBox) return;

    const rect = chatBox.getBoundingClientRect();
    const footer = document.querySelector('footer');
    const footerHeight = footer ? footer.getBoundingClientRect().height : 0;

    const available = window.innerHeight - rect.top - footerHeight - 24;
    chatBox.style.height = Math.max(available, 200) + 'px';
}

window.addEventListener('load', resizeChatBox);
window.addEventListener('resize', resizeChatBox);
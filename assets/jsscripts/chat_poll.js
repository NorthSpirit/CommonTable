const chatContainer = document.getElementById('chat-messages');
const chatForm = document.getElementById('chat-form');
const chatInput = document.getElementById('chat-input');

if (chatContainer) {
    const pollUrl = chatContainer.getAttribute('data-poll-url') || '/chef/chat/messages';

    // Mark as read: syncs the unread-badge's storage the moment this
    // page (the actual chat page) loads, since being here means
    // everything currently visible has been seen.
    const storageKeyMatch = pollUrl.match(/^\/(\w+)\//);
    if (storageKeyMatch) {
        localStorage.setItem('chat_last_read_' + storageKeyMatch[1], chatContainer.getAttribute('data-last-id') || '0');
    }

    function appendMessage(msg) {
        const div = document.createElement('div');
        div.className = 'text-sm gms-light-text';
        div.setAttribute('data-message-id', msg.id);
        div.innerHTML = '<span class="font-bold"></span> <span class="opacity-70"></span> ';
        div.querySelector('.font-bold').textContent = msg.sender_name;
        div.querySelector('.opacity-70').textContent = '(' + msg.sent_at + '):';
        div.appendChild(document.createTextNode(' ' + msg.message));
        chatContainer.appendChild(div);
    }

    async function pollMessages() {
        const lastID = chatContainer.getAttribute('data-last-id') || '0';
        try {
            const response = await fetch(pollUrl + '?after=' + lastID);
            if (!response.ok) return;
            const messages = await response.json();
            if (messages && messages.length > 0) {
                const wasAtBottom = chatContainer.scrollTop + chatContainer.clientHeight >= chatContainer.scrollHeight - 10;
                messages.forEach(appendMessage);
                chatContainer.setAttribute('data-last-id', messages[messages.length - 1].id);
                if (storageKeyMatch) {
                    localStorage.setItem('chat_last_read_' + storageKeyMatch[1], messages[messages.length - 1].id);
                }
                if (wasAtBottom) {
                    chatContainer.scrollTop = chatContainer.scrollHeight;
                }
            }
        } catch (e) {
            // Network hiccup — just try again on the next poll.
        }
    }

    chatContainer.scrollTop = chatContainer.scrollHeight;
    setInterval(pollMessages, 3000);
}

if (chatForm && chatInput) {
    chatForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const message = chatInput.value.trim();
        if (!message) return;

        chatInput.value = '';
        chatInput.focus();

        try {
            await fetch(chatForm.action, {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: 'message=' + encodeURIComponent(message),
            });

            if (chatContainer && typeof pollMessages === 'function') {
                await pollMessages();
            }
        } catch (e) {
            chatInput.value = message;
        }

        chatInput.focus();
    });
}
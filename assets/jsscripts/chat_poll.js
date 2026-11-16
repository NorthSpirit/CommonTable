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

    // Set on every chat page: lists the messages deleted so far, so a deletion shows up without reloading.
    const deletedUrl = chatContainer.getAttribute('data-deleted-url');

    // The message text sits in its own span, so it can be swapped for the "deleted" text later.
    function ensureTextSpan(div) {
        if (div.querySelector('.chat-text')) return;
        const nodes = Array.from(div.childNodes).reverse();
        const last = nodes.find((n) => n.nodeType === 3 && n.nodeValue.trim() !== '');
        if (!last) return;
        const span = document.createElement('span');
        span.className = 'chat-text';
        span.textContent = last.nodeValue;
        div.replaceChild(span, last);
    }

    function markDeleted(id, text) {
        const div = chatContainer.querySelector('[data-message-id="' + id + '"]');
        if (!div || div.hasAttribute('data-deleted')) return;
        ensureTextSpan(div);
        const span = div.querySelector('.chat-text');
        if (span) span.textContent = ' ' + text;
        div.setAttribute('data-deleted', '1');
        div.classList.remove('cursor-pointer', 'hover:opacity-80');
        div.classList.add('italic', 'opacity-60');
    }

    async function syncDeleted() {
        if (!deletedUrl) return;
        try {
            const response = await fetch(deletedUrl);
            if (!response.ok) return;
            const deleted = await response.json();
            deleted.forEach((d) => markDeleted(d.id, d.message));
        } catch (e) {
            // Try again on the next poll.
        }
    }

    // Admin and manager chat pages only: clicking a message opens a small popup to delete it.
    // The page gives the base URL in data-delete-url; staff and chef pages have none, so nothing happens there.
    const deleteBaseUrl = chatContainer.getAttribute('data-delete-url');
    let deletePopup = null;

    function closeDeletePopup() {
        if (deletePopup) {
            deletePopup.remove();
            deletePopup = null;
        }
    }

    function openDeletePopup(messageDiv, x, y) {
        closeDeletePopup();

        const id = messageDiv.getAttribute('data-message-id');

        const popup = document.createElement('div');
        popup.className = 'p-3 rounded-lg shadow-xl border gms-border-medium gms-content-bg gms-main-text text-sm';
        popup.style.position = 'fixed';
        popup.style.zIndex = '50';
        popup.style.width = '16rem';
        popup.style.left = '0px';
        popup.style.top = '0px';

        const title = document.createElement('p');
        title.className = 'font-bold mb-2';
        title.textContent = 'Delete this message?';

        const reason = document.createElement('input');
        reason.type = 'text';
        reason.maxLength = 255;
        reason.placeholder = 'Reason (optional)';
        reason.className = 'w-full p-1 mb-2 border gms-border-medium rounded-md gms-content-bg gms-main-text';

        const confirm = document.createElement('button');
        confirm.type = 'button';
        confirm.textContent = 'Delete';
        confirm.className = 'w-full px-3 py-1 font-bold rounded-lg shadow-md cursor-pointer gms-btn-submit border gms-border-medium';

        const error = document.createElement('p');
        error.className = 'text-xs mt-2 text-red-600';

        popup.append(title, reason, confirm, error);
        document.body.appendChild(popup);
        deletePopup = popup;

        // Place it next to the click, kept inside the window.
        const margin = 8;
        const left = Math.min(x + margin, window.innerWidth - popup.offsetWidth - margin);
        const top = Math.min(y + margin, window.innerHeight - popup.offsetHeight - margin);
        popup.style.left = Math.max(margin, left) + 'px';
        popup.style.top = Math.max(margin, top) + 'px';

        // A click inside the popup must not count as a click away.
        popup.addEventListener('click', (event) => event.stopPropagation());

        async function submit() {
            confirm.disabled = true;
            error.textContent = '';
            try {
                const response = await fetch(deleteBaseUrl + '/' + id + '/delete', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: 'reason=' + encodeURIComponent(reason.value.trim()),
                });
                if (response.ok || response.status === 404) {
                    // Reload so the page shows the message as deleted (404 = it was already deleted).
                    window.location.reload();
                    return;
                }
                error.textContent = (await response.text()).trim() || 'Could not delete.';
            } catch (e) {
                error.textContent = 'Could not delete.';
            }
            confirm.disabled = false;
        }

        confirm.addEventListener('click', submit);
        reason.addEventListener('keydown', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                submit();
            }
        });

        reason.focus();
    }

    // True only if the click landed on the message's actual text, not on empty space in its line.
    function clickedOnMessage(messageDiv, x, y) {
        const range = document.createRange();
        range.selectNodeContents(messageDiv);
        return Array.from(range.getClientRects()).some(
            (rect) => x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom
        );
    }

    if (deleteBaseUrl) {
        chatContainer.addEventListener('click', (event) => {
            const messageDiv = event.target.closest('[data-message-id]');
            if (!messageDiv || messageDiv.hasAttribute('data-deleted') || !clickedOnMessage(messageDiv, event.clientX, event.clientY)) return;
            event.stopPropagation();
            openDeletePopup(messageDiv, event.clientX, event.clientY);
        });
        document.addEventListener('click', closeDeletePopup);
        document.addEventListener('keydown', (event) => {
            if (event.key === 'Escape') closeDeletePopup();
        });
    }

    function appendMessage(msg) {
        const div = document.createElement('div');
        div.className = 'text-sm gms-light-text' + (deleteBaseUrl && !msg.deleted ? ' cursor-pointer hover:opacity-80 w-fit max-w-full' : '');
        if (msg.deleted) {
            div.setAttribute('data-deleted', '1');
            div.classList.add('italic', 'opacity-60');
        }
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
        await syncDeleted();
    }

    if (deleteBaseUrl) {
        chatContainer.querySelectorAll('[data-message-id]').forEach((div) => {
            if (div.hasAttribute('data-deleted')) return;
            div.classList.add('cursor-pointer', 'hover:opacity-80', 'w-fit', 'max-w-full');
        });
    }

    // Deleted messages are shown dimmed and cannot be clicked.
    chatContainer.querySelectorAll('[data-deleted]').forEach((div) => div.classList.add('italic', 'opacity-60'));

    syncDeleted();

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
                redirect: 'manual',
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
// Shows a small unread-dot on the nav "Chat" link when new messages
// exist that haven't been viewed yet. Runs on every page via the nav
// link itself — reads its own data attributes, so no page-specific
// wiring beyond the two attributes below.
(function () {
    const navLink = document.getElementById('nav-chat-link');
    if (!navLink) return;

    const pollUrl = navLink.getAttribute('data-poll-url');
    const storageKey = navLink.getAttribute('data-storage-key');
    if (!pollUrl || !storageKey) return;

    function getLastReadID() {
        return localStorage.getItem(storageKey) || '0';
    }

    function showBadge() {
        let dot = navLink.querySelector('.unread-dot');
        if (!dot) {
            dot = document.createElement('span');
            dot.className = 'unread-dot inline-block w-2 h-2 rounded-full bg-red-500 ml-1 align-middle';
            navLink.appendChild(dot);
        }
    }

    function hideBadge() {
        const dot = navLink.querySelector('.unread-dot');
        if (dot) dot.remove();
    }

    async function checkUnread() {
        try {
            const response = await fetch(pollUrl + '?after=' + getLastReadID());
            if (!response.ok) return;
            const messages = await response.json();
            if (messages && messages.length > 0) {
                showBadge();
            } else {
                hideBadge();
            }
        } catch (e) {
            // Network hiccup — leave badge state as-is, try again next tick.
        }
    }

    checkUnread();
    setInterval(checkUnread, 5000);
})();
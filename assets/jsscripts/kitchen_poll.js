async function pollKitchenBoard() {
    const board = document.getElementById('order-board');
    if (!board) return;

    try {
        const response = await fetch('/kitchen/orders/poll');
        if (!response.ok) return;
        const orders = await response.json();

        board.innerHTML = '';

        if (!orders || orders.length === 0) {
            const empty = document.createElement('p');
            empty.className = 'text-center italic gms-main-text col-span-full';
            empty.textContent = 'No active orders.';
            board.appendChild(empty);
            return;
        }

        orders.forEach((order) => {
            const card = document.createElement('div');
            card.className = 'border gms-border-medium rounded-lg shadow-lg overflow-hidden gms-content-bg';

            const header = document.createElement('div');
            header.className = 'gms-dark-bg gms-light-text px-3 py-2 flex justify-between items-center';
            const tableSpan = document.createElement('span');
            tableSpan.className = 'font-bold';
            tableSpan.textContent = 'Table ' + order.table_number;
            const statusSpan = document.createElement('span');
            statusSpan.className = 'text-xs uppercase opacity-80';
            statusSpan.textContent = order.status;
            header.appendChild(tableSpan);
            header.appendChild(statusSpan);

            const body = document.createElement('div');
            body.className = 'p-3';

            const meta = document.createElement('p');
            meta.className = 'text-xs opacity-60 mb-2';
            meta.textContent = 'Order #' + order.order_number + ' — ' + order.created_at;
            body.appendChild(meta);

            const list = document.createElement('ul');
            list.className = 'text-sm gms-main-text mb-3 space-y-1';
            (order.items || []).forEach((item) => {
                const li = document.createElement('li');
                const qty = document.createElement('span');
                qty.className = 'font-semibold';
                qty.textContent = item.quantity + 'x ';
                li.appendChild(qty);
                li.appendChild(document.createTextNode(item.name));
                if (item.notes) {
                    const noteSpan = document.createElement('span');
                    noteSpan.className = 'block text-xs opacity-70 pl-4';
                    noteSpan.textContent = item.notes;
                    li.appendChild(noteSpan);
                }
                list.appendChild(li);
            });
            body.appendChild(list);

            if (order.notes) {
                const noteP = document.createElement('p');
                noteP.className = 'text-xs italic opacity-70 mb-3';
                noteP.textContent = 'Note: ' + order.notes;
                body.appendChild(noteP);
            }

            if (order.next_status) {
                const form = document.createElement('form');
                form.method = 'POST';
                form.action = '/kitchen/orders/status';

                const idInput = document.createElement('input');
                idInput.type = 'hidden';
                idInput.name = 'order_id';
                idInput.value = order.id;
                form.appendChild(idInput);

                const statusInput = document.createElement('input');
                statusInput.type = 'hidden';
                statusInput.name = 'status';
                statusInput.value = order.next_status;
                form.appendChild(statusInput);

                const button = document.createElement('button');
                button.type = 'submit';
                button.className = 'w-full px-3 py-2 text-sm font-bold rounded-md shadow gms-btn-submit border gms-border-medium';
                button.textContent = order.next_status_label;
                form.appendChild(button);

                body.appendChild(form);
            }

            card.appendChild(header);
            card.appendChild(body);
            board.appendChild(card);
        });
    } catch (e) {
        // Network hiccup — try again next poll.
    }
}

setInterval(pollKitchenBoard, 4000);
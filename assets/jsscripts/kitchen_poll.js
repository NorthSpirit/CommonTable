async function pollKitchenBoard() {
    // Gets the order board the HTML object, if there is none, returns without breaking anything.
    const board = document.getElementById('order-board');
    if (!board) return;

    try {
        const response = await fetch('/kitchen/orders/poll');
        if (!response.ok) return;
        const orders = await response.json();

        board.innerHTML = '';

        // Obviously, if there are no orders (the list has 0 items)
        // !orders is only a safeguard, in case no list arrives at all
        if (!orders || orders.length === 0) {
            const empty = document.createElement('p');
            empty.className = 'text-center italic gms-main-text col-span-full';
            empty.textContent = 'No active orders.';
            board.appendChild(empty);
            return;
        }

        // For every order, reconstructs the display.
        orders.forEach((order) => {
            // Components have to match the ones on the templ templates
            // Otherwise first one upon opening the screen will differ than the ones later on
            // The order card box part, the one, contained all else
            const card = document.createElement('div');
            card.className = 'border gms-border-very-dark rounded-lg shadow-lg overflow-hidden gms-content-bg flex flex-col';

            // the "table name" + status part
            const header = document.createElement('div');
            header.className = 'gms-dark-bg gms-light-text px-3 py-2 flex justify-between items-center';
            const tableSpan = document.createElement('span');
            tableSpan.className = 'font-bold';
            // Gets the table number from the order
            tableSpan.textContent = 'Table ' + order.table_number;
            const statusSpan = document.createElement('span');
            // And it's status on the other side
            statusSpan.className = 'text-xs uppercase opacity-80';
            statusSpan.textContent = order.status;
            // Actually puts them inside the header part
            header.appendChild(tableSpan);
            header.appendChild(statusSpan);

            // Part below the header box thingy, the actual body
            const body = document.createElement('div');
            body.className = 'p-3 flex flex-col flex-1';

            // Part which tells which order it is (number) and it's date, for little extra prioritizing for the chef
            const meta = document.createElement('p');
            meta.className = 'text-xs opacity-60 mb-2';
            meta.textContent = 'Order #' + order.order_number + ' — ' + order.created_at;
            body.appendChild(meta);

            // Actual list of items in the order written as amount x unit.
            // Note, if the order was adjusted, it will repeat the item.
            // As in 4x Pink Soup, when line below it can be 2x Pink Soup, when line below it can be 3x Pink Soup
            const list = document.createElement('ul');
            list.className = 'text-sm gms-main-text mb-3 space-y-1';
            (order.items || []).forEach((item) => {
                const li = document.createElement('li');
                const qty = document.createElement('span');
                qty.className = 'font-semibold';
                qty.textContent = item.quantity + 'x ';
                // Actually appends the quantity to the box
                li.appendChild(qty);
                // And actually appends the name of the item to be box
                li.appendChild(document.createTextNode(item.name));
                // Finally, if the notes exists (FOR ITEM), it adds them, if not, there is nothing
                if (item.notes) {
                    const noteSpan = document.createElement('span');
                    noteSpan.className = 'block text-xs opacity-70 pl-4';
                    noteSpan.textContent = item.notes;
                    // Actually appends the notes there
                    li.appendChild(noteSpan);
                }
                list.appendChild(li);
            });
            body.appendChild(list);

            // Orders notes (FOR ORDER)
            if (order.notes) {
                const noteP = document.createElement('p');
                noteP.className = 'text-xs italic opacity-70 mb-3';
                noteP.textContent = 'Note: ' + order.notes;
                body.appendChild(noteP);
            }

            // Constructs the button (as in, full thing), once again, it has to match the one in the templ template
            // (internal/templates/kitchen/board_page.templ one)
            if (order.next_status) {
                const form = document.createElement('form');
                form.method = 'POST';
                form.action = '/kitchen/orders/status';
                form.className = 'mt-auto';

                // Order ID and status is not really meant to be inputted by hand
                // Instead it's gotten from the order itself
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

                // Constructs the button (as in the clickable component)
                const button = document.createElement('button');
                button.type = 'submit';
                button.className = 'w-full px-3 py-2 text-sm font-bold rounded-md shadow gms-btn-submit border gms-border-medium';
                button.textContent = order.next_status_label;
                form.appendChild(button);

                body.appendChild(form);
            }

            // Constructs the full card from its parts
            card.appendChild(header);
            card.appendChild(body);
            board.appendChild(card);
        });
    } catch (e) {
        // Network hiccup — try again next poll.
    }
}

// Fires the script every 4 seconds, realistically, could be set even longer
// Because humans are pretty rubbish event listeners
// Like, kitchen staff or a singular chef will probably either have his hands full and not react instantly
// Or on a slow day, will be half bored and also notice slowly - nobody have their eyes glued to the dashboard 24/7.
setInterval(pollKitchenBoard, 4000);
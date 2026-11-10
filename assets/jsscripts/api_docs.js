// Expand all / Collapse all buttons on the waiter API page
function setAllEndpoints(open) {
    document.querySelectorAll('details.api-endpoint').forEach(function (d) {
        d.open = open;
    });
}

document.getElementById('apiExpandAll').addEventListener('click', function () {
    setAllEndpoints(true);
});

document.getElementById('apiCollapseAll').addEventListener('click', function () {
    setAllEndpoints(false);
});

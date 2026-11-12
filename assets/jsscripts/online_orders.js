const approveDialog = document.getElementById('approve-dialog');
const approveForm = document.getElementById('approve-form');

document.querySelectorAll('[data-approve-action]').forEach(function (button) {
    button.addEventListener('click', function () {
        approveForm.action = button.dataset.approveAction;
        approveDialog.showModal();
    });
});

document.querySelectorAll('[data-dialog-close]').forEach(function (button) {
    button.addEventListener('click', function () {
        approveDialog.close();
    });
});

// Turns every [data-wysiwyg] element into a Quill editor that writes its
// HTML into a hidden input on submit.
(function () {
    if (typeof Quill === 'undefined') return;

    // Make Quill's default (light) look follow the page palette. The outer
    // wrapper in the template draws the single border, so Quill's own
    // borders are removed. !important is needed because Quill's stylesheet
    // is loaded after this one.
    var style = document.createElement('style');
    style.textContent =
        '.ql-toolbar.ql-snow,.ql-container.ql-snow{border:0 !important}' +
        '.ql-snow .ql-stroke{stroke:currentColor !important}' +
        '.ql-snow .ql-fill{fill:currentColor !important}' +
        '.ql-snow .ql-picker{color:inherit}' +
        '.ql-container.ql-snow{font-size:1rem}' +
        '.ql-editor{min-height:12rem}' +
        '.ql-toolbar.ql-snow button.wysiwyg-token{' +
            'width:auto !important;height:24px;padding:0 .6rem !important;margin-left:.25rem;' +
            'white-space:nowrap;font-size:.75rem;font-weight:600;color:inherit;' +
            'float:left}';
    document.body.appendChild(style);

    var editors = [];

    document.querySelectorAll('[data-wysiwyg]').forEach(function (box) {
        var input = document.getElementById(box.dataset.input);
        var toolbar = document.getElementById(box.dataset.toolbar);
        if (!input || !toolbar) return;

        var quill = new Quill(box, {
            theme: 'snow',
            modules: { toolbar: { container: toolbar } }
        });

        if (input.value) {
            quill.setContents(quill.clipboard.convert({ html: input.value }));
        }

        // Toolbar buttons that insert a placeholder token at the cursor.
        toolbar.querySelectorAll('[data-token]').forEach(function (btn) {
            btn.addEventListener('click', function () {
                var range = quill.getSelection(true);
                var text = '\n' + btn.dataset.token + '\n';
                quill.insertText(range.index, text, 'user');
                quill.setSelection(range.index + text.length);
            });
        });

        editors.push({
            quill: quill,
            input: input,
            required: box.dataset.required === 'true',
            label: box.dataset.label || 'This field'
        });
    });

    var form = document.querySelector('form[data-wysiwyg-form]');
    if (!form) return;

    form.addEventListener('submit', function (e) {
        for (var i = 0; i < editors.length; i++) {
            var ed = editors[i];
            var empty = ed.quill.getText().trim() === '';
            if (empty && ed.required) {
                e.preventDefault();
                alert(ed.label + ' is required.');
                return;
            }
            ed.input.value = empty ? '' : ed.quill.root.innerHTML;
        }
    });
})();

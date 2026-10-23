(function () {
  function noticeFor(el) {
    var next = el.nextElementSibling;
    if (next && next.hasAttribute && next.hasAttribute('data-limit-notice')) {
      return next;
    }
    var p = document.createElement('p');
    p.setAttribute('data-limit-notice', '');
    p.setAttribute('role', 'alert');
    p.className = 'text-xs mt-1';
    p.style.color = '#b91c1c';
    p.style.display = 'none';
    el.insertAdjacentElement('afterend', p);
    return p;
  }

  function show(el, text) {
    var p = noticeFor(el);
    p.textContent = text;
    p.style.display = 'block';
  }

  function hide(el) {
    var next = el.nextElementSibling;
    if (next && next.hasAttribute && next.hasAttribute('data-limit-notice')) {
      next.style.display = 'none';
    }
  }

  function checkText(el) {
    var max = parseInt(el.getAttribute('maxlength'), 10);
    if (!max) return;
    if (el.value.length >= max) {
      show(el, 'Maximum ' + max + ' characters reached.');
    } else {
      hide(el);
    }
  }

  function checkNumber(el) {
    if (el.value === '') { hide(el); return; }
    var v = parseFloat(el.value);
    var min = el.getAttribute('min');
    var max = el.getAttribute('max');
    if (max !== null && v > parseFloat(max)) {
      show(el, 'Maximum value is ' + max + '.');
    } else if (min !== null && v < parseFloat(min)) {
      show(el, 'Minimum value is ' + min + '.');
    } else {
      hide(el);
    }
  }

  document.addEventListener('input', function (e) {
    var el = e.target;
    if (!(el instanceof HTMLInputElement || el instanceof HTMLTextAreaElement)) return;
    if (el.type === 'number') { checkNumber(el); return; }
    if (el.hasAttribute('maxlength')) checkText(el);
  });

  document.addEventListener('beforeinput', function (e) {
    var el = e.target;
    if (!(el instanceof HTMLInputElement || el instanceof HTMLTextAreaElement)) return;
    if (!el.hasAttribute('maxlength') || el.type === 'number') return;
    if (e.inputType && e.inputType.indexOf('insert') === 0 &&
        el.selectionStart === el.selectionEnd &&
        el.value.length >= parseInt(el.getAttribute('maxlength'), 10)) {
      show(el, 'Maximum ' + el.getAttribute('maxlength') + ' characters reached.');
    }
  });
})();

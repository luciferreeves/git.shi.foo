(function () {
  var toggle = document.getElementById('user-menu-toggle');
  var panel = document.getElementById('user-menu-panel');
  if (!toggle || !panel) {
    return;
  }

  toggle.addEventListener('click', function (event) {
    event.stopPropagation();
    panel.classList.toggle('is-open');
  });

  document.addEventListener('click', function (event) {
    if (!panel.contains(event.target) && !toggle.contains(event.target)) {
      panel.classList.remove('is-open');
    }
  });
})();
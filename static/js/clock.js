(function () {
  var clock = document.getElementById('footer-clock');
  if (!clock) {
    return;
  }

  function pad(value) {
    return value < 10 ? '0' + value : '' + value;
  }

  function render() {
    var now = new Date();
    clock.textContent =
      now.getFullYear() + '-' + pad(now.getMonth() + 1) + '-' + pad(now.getDate()) + ' ' +
      pad(now.getHours()) + ':' + pad(now.getMinutes()) + ':' + pad(now.getSeconds());
  }

  render();
  setInterval(render, 1000);
})();
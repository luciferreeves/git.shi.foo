(function () {
  const buttons = document.querySelectorAll(".clone-copy");
  if (!buttons.length || !navigator.clipboard) {
    return;
  }

  buttons.forEach(function (button) {
    button.addEventListener("click", function () {
      const field = button.parentElement.querySelector(".clone-url");
      if (!field) {
        return;
      }

      navigator.clipboard.writeText(field.value).then(function () {
        const original = button.textContent;
        button.textContent = "copied";
        button.classList.add("is-copied");
        window.setTimeout(function () {
          button.textContent = original;
          button.classList.remove("is-copied");
        }, 1200);
      });
    });
  });
})();

(function () {
  const container = document.querySelector(".import-progress[data-events]");
  if (!container) {
    return;
  }

  const bar = container.querySelector('[data-role="bar"]');
  const phaseLabel = container.querySelector('[data-role="phase"]');
  const percentLabel = container.querySelector('[data-role="percent"]');
  const eventsUrl = container.getAttribute("data-events");

  const source = new EventSource(eventsUrl);

  source.addEventListener("message", function (message) {
    let progress;
    try {
      progress = JSON.parse(message.data);
    } catch (parseError) {
      return;
    }

    if (typeof progress.percent === "number" && bar) {
      bar.style.width = progress.percent + "%";
    }
    if (typeof progress.percent === "number" && percentLabel) {
      percentLabel.textContent = progress.percent;
    }
    if (progress.phase && phaseLabel) {
      phaseLabel.textContent = progress.phase;
    }

    if (progress.done) {
      source.close();
      window.location.reload();
    }
  });
})();
(function () {
  var root = document.documentElement;
  var button = document.getElementById("theme-toggle");
  if (!button) return;

  function label() {
    // The button names the theme you will switch to.
    button.textContent = root.getAttribute("data-theme") === "dark" ? "Light" : "Dark";
  }

  button.addEventListener("click", function () {
    var next = root.getAttribute("data-theme") === "dark" ? "light" : "dark";
    root.setAttribute("data-theme", next);
    try {
      localStorage.setItem("theme", next);
    } catch (e) {}
    label();
  });

  label();
})();

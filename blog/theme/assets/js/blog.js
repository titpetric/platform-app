// assets/js/blog.js
class BlogApp {
  constructor() {
  }
}

document.addEventListener('DOMContentLoaded', () => new DailyApp());

class ThemeMachine extends HTMLElement {
    constructor() {
      super();
      this.themeToggle;
      this.themeDisplay;
    }

    connectedCallback() {
      this.setupControl("appearance", document.documentElement);
      this.setupControl("theme", document.body);

      this.themeToggle = this.querySelector(".theme-display-toggle");
      this.themeDisplay = this.querySelector(".theme-display-wrapper");

      this.handleClickOutside();
      this.themeToggle.addEventListener("click", this.handleThemeToggleClick.bind(this));
    }

    handleChange(e, prop, el) {
      const attr = `data-${prop}`;
      const value = e.target.value;

      if (!value) {
        localStorage.removeItem(prop);
        el.removeAttribute(attr);
        return;
      }

      localStorage.setItem(prop, value);
      el.setAttribute(attr, value);
    }

    handleThemeToggleClick() {
      let expanded = this.themeToggle.getAttribute("aria-expanded") === "true" || false;

      this.themeToggle.setAttribute("aria-expanded", !expanded);
      this.themeDisplay.toggleAttribute("hidden", expanded);
    }

    handleClickOutside() {
      document.addEventListener(
        "click",
        (e) => {
          if (!e.target.closest("theme-machine")) {
            this.themeToggle.setAttribute("aria-expanded", false);
            this.themeDisplay.toggleAttribute("hidden", true);
          }
        },
        false
      );
    }

    setupControl(prop, el) {
      const initialValue = localStorage.getItem(prop) || "";
      const collection = this.querySelectorAll(`[name='${prop}']`);

      for (let item of collection) {
        item.checked = item.value === initialValue;
        item.addEventListener("change", (e) => this.handleChange(e, prop, el));
      }
    }
  }

  if ("customElements" in window) {
    window.customElements.define("theme-machine", ThemeMachine);
  }

// assets/js/daily.js
class DailyApp {
  constructor(formId = 'task-form', inputId = 'task-input', listId = 'task-list') {
    this.form = document.getElementById(formId);
    this.input = document.getElementById(inputId);
    this.list = document.getElementById(listId);

    if (!this.form || !this.input || !this.list) {
      throw new Error('DailyApp: required elements not found');
    }

    this.form.addEventListener('submit', (e) => this.onSubmit(e));
    this.list.addEventListener('change', (e) => this.onCheckboxChange(e));
  }

  async onSubmit(e) {
    e.preventDefault();
    const val = this.input.value && this.input.value.trim();
    if (!val) return;

    const payload = { title: val };

    try {
      const res = await fetch('/daily/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
        credentials: 'same-origin'
      });

      if (res.ok) {
        this.input.value = '';
        // reload to get server-side rendered list
        window.location.reload();
      } else {
        console.error('Save failed', res.status);
        alert('Could not save task (server returned ' + res.status + ')');
      }
    } catch (err) {
      console.error('Network error', err);
      alert('Network error while saving task');
    }
  }

  async onCheckboxChange(e) {
    const checkbox = e.target;
    if (!checkbox.classList.contains('task-checkbox')) return;

    const taskItem = checkbox.closest('.task-item');
    const taskId = taskItem?.dataset.taskId;
    if (!taskId) return;

    try {
      const res = await fetch(`/daily/complete/${encodeURIComponent(taskId)}`, {
        method: 'POST',
        credentials: 'same-origin'
      });

      if (res.ok) {
        // fade out for 0.3s
        taskItem.style.transition = 'opacity 0.3s';
        taskItem.style.opacity = '0';
        setTimeout(() => taskItem.remove(), 300);
      } else {
        console.error('Failed to complete task', res.status);
        checkbox.checked = false;
        alert('Could not complete task (server returned ' + res.status + ')');
      }
    } catch (err) {
      console.error('Network error', err);
      checkbox.checked = false;
      alert('Network error while completing task');
    }
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

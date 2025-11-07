// assets/js/daily.js
class DailyApp {
  constructor(formId = 'task-form', inputId = 'task-input', listId = 'task-list', themeBtnId = 'theme-toggle') {
    this.form = document.getElementById(formId);
    this.input = document.getElementById(inputId);
    this.list = document.getElementById(listId);
    this.themeBtn = document.getElementById(themeBtnId);

    if (!this.form || !this.input || !this.list) {
      throw new Error('DailyApp: required elements not found');
    }

    this.form.addEventListener('submit', (e) => this.onSubmit(e));
    this.list.addEventListener('change', (e) => this.onCheckboxChange(e));

    if (this.themeBtn) {
      this.themeBtn.addEventListener('click', () => this.toggleTheme());
    }
    this.loadTheme();
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

  toggleTheme() {
    const body = document.body;
    const isDark = body.classList.contains('dark');
    if (isDark) {
      body.classList.remove('dark');
      body.classList.add('light');
      localStorage.setItem('daily:theme', 'light');
      this.updateToggleLabel('Light');
    } else {
      body.classList.remove('light');
      body.classList.add('dark');
      localStorage.setItem('daily:theme', 'dark');
      this.updateToggleLabel('Dark');
    }
  }

  loadTheme() {
    const saved = localStorage.getItem('daily:theme');
    const body = document.body;
    if (saved === 'light') {
      body.classList.add('light');
      this.updateToggleLabel('Light');
    } else {
      body.classList.add('dark');
      this.updateToggleLabel('Dark');
    }
  }

  updateToggleLabel(text) {
    if (!this.themeBtn) return;
    this.themeBtn.textContent = text;
  }
}

document.addEventListener('DOMContentLoaded', () => new DailyApp());

package cmdui

const uiAssetsJSNodes = `    async function refreshNodes() {
      try {
        const res = await fetch('/api/ssh/nodes');
        const data = await res.json();
        const tbody = document.getElementById('nodes-body');
        tbody.innerHTML = '';
        if (!data || data.length === 0) {
          tbody.innerHTML = '<tr><td colspan="6" style="text-align:center">No nodes registered. Use "gitmap ssh join" to add nodes.</td></tr>';
          return;
        }
        data.forEach(n => {
          const row = document.createElement('tr');
          row.innerHTML = ` + "`" + `<td>${n.nodeId || '-'}</td><td><b>${n.alias}</b></td><td>${n.host}</td><td>${n.os || 'Linux'}</td><td><span class="badge ${n.isOnline ? 'badge-online' : 'badge-offline'}">${n.isOnline ? 'Online' : 'Offline'}</span></td><td><button class="btn" style="padding:0.25rem 0.5rem;font-size:0.75rem" onclick="testNode('${n.alias}')">Ping</button></td>` + "`" + `;
          tbody.appendChild(row);
        });
      } catch (e) { console.error(e); }
    }

    async function populateEditorNodes() {
      try {
        const res = await fetch('/api/ssh/nodes');
        const data = await res.json();
        const sel = document.getElementById('editor-node');
        if (sel) {
          sel.innerHTML = '<option value="local">local (current machine)</option>';
          if (data && data.length) {
            data.forEach(n => {
              const opt = document.createElement('option');
              opt.value = n.alias;
              opt.innerText = n.alias + ' (' + n.host + ')';
              sel.appendChild(opt);
            });
          }
        }
        const termSel = document.getElementById('term-node');
        if (termSel && data && data.length) {
          termSel.innerHTML = '<option value="local">local</option>';
          data.forEach(n => {
            const opt = document.createElement('option');
            opt.value = n.alias;
            opt.innerText = n.alias;
            termSel.appendChild(opt);
          });
        }
      } catch (e) {}
    }

    async function openRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      if (!path) { showToast('Please enter a file path', 'error'); return; }
      document.getElementById('editor-status').innerText = 'Loading ' + path + '...';
      try {
        const res = await fetch('/api/editor/read', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({nodeAlias: node, filePath: path})
        });
        const data = await res.json();
        if (data.isSuccess) {
          document.getElementById('editor-area').value = data.content;
          document.getElementById('editor-status').innerText = 'Loaded ' + path + ' from ' + node;
          document.getElementById('editor-lang').innerText = 'Language: ' + (data.language || 'plaintext');
        } else {
          showToast('Failed to read file: ' + (data.error || 'unknown error'), 'error');
        }
      } catch (e) { showToast('Error reading file: ' + e.message, 'error'); }
    }

    async function saveRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      const content = document.getElementById('editor-area').value;
      if (!path) { showToast('Please enter a file path', 'error'); return; }
      document.getElementById('editor-status').innerText = 'Saving to ' + node + '...';
      try {
        const res = await fetch('/api/editor/save', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({nodeAlias: node, filePath: path, content: content})
        });
        const data = await res.json();
        if (data.success) {
          document.getElementById('editor-status').innerText = 'Saved successfully at ' + new Date().toLocaleTimeString();
        } else {
          showToast('Failed to save file: ' + (data.error || 'unknown error'), 'error');
        }
      } catch (e) { showToast('Error saving file: ' + e.message, 'error'); }
    }

`

const uiAssetsJSSettings = `    function showSettingsSubTab(subId) {
      document.querySelectorAll('.sub-tab-btn').forEach(btn => btn.classList.remove('active'));
      const btn = document.getElementById('subnav-' + subId);
      if (btn) btn.classList.add('active');

      const panes = document.querySelectorAll('.settings-subtab-pane');
      if (subId === 'all') {
        panes.forEach(p => p.style.display = 'block');
        return;
      }
      panes.forEach(p => p.style.display = 'none');
      const target = document.getElementById('settings-pane-' + subId);
      if (target) target.style.display = 'block';
    }

    async function loadSettings() {
      try {
        const res = await fetch('/api/settings');
        const data = await res.json();
        if (data) {
          if (data.theme) applyTheme(data.theme);
          if (data.defaultRemote && document.getElementById('setting-remote')) document.getElementById('setting-remote').value = data.defaultRemote;
          if (data.clusterPort && document.getElementById('setting-port')) document.getElementById('setting-port').value = data.clusterPort;
          if (data.graphicsMode && document.getElementById('setting-graphics')) document.getElementById('setting-graphics').value = data.graphicsMode;
          if (data.autoOpenBrowser !== undefined && document.getElementById('setting-auto-open')) document.getElementById('setting-auto-open').value = String(data.autoOpenBrowser);
          if (data.commitInLayout && document.getElementById('setting-commitin-layout')) document.getElementById('setting-commitin-layout').value = data.commitInLayout;
          if (data.pullDirection && document.getElementById('setting-pull-direction')) document.getElementById('setting-pull-direction').value = data.pullDirection;
          if (data.prReplayMode && document.getElementById('setting-pr-mode')) document.getElementById('setting-pr-mode').value = data.prReplayMode;

          if (data.attributes) {
            document.querySelectorAll('[data-key]').forEach(el => {
              const key = el.getAttribute('data-key');
              if (data.attributes[key] !== undefined) {
                el.value = data.attributes[key];
              }
            });
          }
        }
      } catch (e) { console.error('Failed to load settings', e); }
    }

    async function saveSettings() {
      const attrs = {};
      document.querySelectorAll('[data-key]').forEach(el => {
        const key = el.getAttribute('data-key');
        if (key) {
          attrs[key] = el.value;
        }
      });

      const payload = {
        theme: document.getElementById('setting-theme') ? document.getElementById('setting-theme').value : 'dark',
        defaultRemote: document.getElementById('setting-remote') ? document.getElementById('setting-remote').value : 'origin',
        clusterPort: parseInt(document.getElementById('setting-port') ? document.getElementById('setting-port').value : 49152) || 49152,
        graphicsMode: document.getElementById('setting-graphics') ? document.getElementById('setting-graphics').value : 'high',
        autoOpenBrowser: document.getElementById('setting-auto-open') ? document.getElementById('setting-auto-open').value === 'true' : true,
        commitInLayout: document.getElementById('setting-commitin-layout') ? document.getElementById('setting-commitin-layout').value : 'split',
        pullDirection: document.getElementById('setting-pull-direction') ? document.getElementById('setting-pull-direction').value : 'pull-left',
        prReplayMode: document.getElementById('setting-pr-mode') ? document.getElementById('setting-pr-mode').value : 'merges',
        attributes: attrs
      };
      try {
        const res = await fetch('/api/settings', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload)
        });
        const data = await res.json();
        if (data.success) {
          showToast('Settings saved successfully!', 'success');
        } else {
          showToast('Failed to save settings: ' + (data.error || 'Unknown error'), 'error');
        }
      } catch (e) { showToast('Failed to save settings: ' + e.message, 'error'); }
    }

    let termHistory = [];
    let termHistoryIdx = -1;

    // ---- Theme system, icons, toasts, responsive shell (program 246, subtask 05) ----
    const ICONS = {
      settings: '<circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>',
      package: '<path d="M16.5 9.4 7.55 4.24"/><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><path d="M3.3 7 12 12l8.7-5"/><path d="M12 22V12"/>',
      server: '<rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><path d="M6 6h.01M6 18h.01"/>',
      'scroll-text': '<path d="M8 21h12a2 2 0 0 0 2-2v-2H10v2a2 2 0 1 1-4 0V5a2 2 0 1 0-4 0v3h4"/><path d="M19 17V5a2 2 0 0 0-2-2H4"/><path d="M15 8h-5"/><path d="M15 12h-5"/>',
      rocket: '<path d="M4.5 16.5c-1.5 1.26-2 5-2 5s3.74-.5 5-2c.71-.84.7-2.13-.09-2.91a2.18 2.18 0 0 0-2.91-.09z"/><path d="m12 15-3-3a22 22 0 0 1 2-3.95A12.88 12.88 0 0 1 22 2c0 2.72-.78 7.5-6 11a22.35 22.35 0 0 1-4 2z"/><path d="M9 12H4s.55-3.03 2-4c1.62-1.08 5 0 5 0"/><path d="M12 15v5s3.03-.55 4-2c1.08-1.62 0-5 0-5"/>',
      'message-square': '<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>',
      'arrow-left-right': '<path d="M8 3 4 7l4 4"/><path d="M4 7h16"/><path d="m16 21 4-4-4-4"/><path d="M20 17H4"/>',
      clock: '<circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/>',
      'pen-line': '<path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/>',
      'book-open': '<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>',
      check: '<path d="M20 6 9 17l-5-5"/>',
      'alert-triangle': '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><path d="M12 9v4"/><path d="M12 17h.01"/>',
      info: '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>',
      x: '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>'
    };
    function icon(name) { return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' + (ICONS[name] || '') + '</svg>'; }

    function showToast(message, kind) {
      kind = kind || 'info';
      const stack = document.getElementById('toast-stack');
      if (!stack) return;
      while (stack.children.length >= 4) stack.removeChild(stack.firstChild);
      const el = document.createElement('div');
      el.className = 'toast toast-' + kind;
      const ic = kind === 'success' ? 'check' : (kind === 'error' ? 'alert-triangle' : 'info');
      el.innerHTML = '<span class="toast-icon">' + icon(ic) + '</span><span class="toast-msg"></span>';
      el.querySelector('.toast-msg').textContent = message;
      el.addEventListener('click', function () { el.remove(); });
      stack.appendChild(el);
      setTimeout(function () { el.classList.add('toast-out'); setTimeout(function () { el.remove(); }, 250); }, 3500);
    }

    const NAV_ITEMS = [
      { id: 'settings', label: 'Settings', icon: 'settings' },
      { id: 'commitin', label: 'Commitin', icon: 'package' },
      { id: 'ssh', label: 'SSH Fleet', icon: 'server' },
      { id: 'macro', label: 'Macros', icon: 'scroll-text' },
      { id: 'installer', label: 'Installers', icon: 'rocket' },
      { id: 'prompts', label: 'Prompts & AI', icon: 'message-square' },
      { id: 'import-export', label: 'Import / Export', icon: 'arrow-left-right' },
      { id: 'schedules', label: 'Schedules', icon: 'clock' },
      { id: 'editor', label: 'Remote Editor', icon: 'pen-line' },
      { id: 'help', label: 'CLI & AGY Help', icon: 'book-open' }
    ];
    function renderNav() {
      const nav = document.getElementById('nav-list');
      if (!nav) return;
      nav.innerHTML = '';
      NAV_ITEMS.forEach(function (item, i) {
        const b = document.createElement('button');
        b.className = 'nav-btn' + (i === 0 ? ' active' : '');
        b.title = item.label;
        b.innerHTML = '<span class="nav-icon">' + icon(item.icon) + '</span><span class="nav-label">' + item.label + '</span>';
        b.addEventListener('click', function () { navGo(item.id, b); });
        nav.appendChild(b);
      });
    }
    function navGo(tabId, btn) {
      showTab(tabId);
      document.querySelectorAll('#nav-list .nav-btn').forEach(function (b) { b.classList.toggle('active', b === btn); });
      const sb = document.getElementById('sidebar');
      if (sb) sb.classList.remove('sidebar-open');
      const scrim = document.getElementById('sidebar-scrim');
      if (scrim) scrim.classList.remove('open');
    }
    function toggleSidebar() {
      const sb = document.getElementById('sidebar');
      if (!sb) return;
      sb.classList.toggle('sidebar-open');
      const scrim = document.getElementById('sidebar-scrim');
      if (scrim) scrim.classList.toggle('open', sb.classList.contains('sidebar-open'));
    }

    function applyTheme(name) {
      const theme = (name === 'light') ? 'light' : 'dark';
      document.documentElement.setAttribute('data-theme', theme);
      const sel = document.getElementById('setting-theme');
      if (sel && sel.value !== theme) sel.value = theme;
    }
    async function onThemeChange(e) {
      const sel = e.target;
      const next = sel.value === 'light' ? 'light' : 'dark';
      const prev = document.documentElement.getAttribute('data-theme') || 'dark';
      applyTheme(next);
      try {
        const res = await fetch('/api/settings', { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({ theme: next }) });
        const data = await res.json();
        if (!data || data.success === false) throw new Error((data && data.error) || 'save failed');
      } catch (err) {
        applyTheme(prev);
        showToast('Theme change failed: ' + err.message, 'error');
      }
    }
    function wireTheme() {
      const sel = document.getElementById('setting-theme');
      if (sel && !sel.dataset.themeWired) { sel.dataset.themeWired = '1'; sel.addEventListener('change', onThemeChange); }
    }

    (function syncDocTitle() {
      if (typeof showTab !== 'function') return;
      const orig = showTab;
      showTab = function (tabId) {
        orig(tabId);
        try {
          const item = NAV_ITEMS.find(function (n) { return n.id === tabId; });
          document.title = 'GitMap — ' + (item ? item.label : tabId);
        } catch (e) {}
      };
    })();

    window.addEventListener('load', function () {
      renderNav();
      wireTheme();
      document.title = 'GitMap — Settings';
    });

`

package cmdui

// IndexHTML provides the embedded single-page UI application.
const IndexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitMap Management & Fleet Dashboard</title>
  <style>
    :root {
      --bg: #0f172a;
      --card: #1e293b;
      --border: #334155;
      --text: #f8fafc;
      --muted: #94a3b8;
      --primary: #f59e0b;
      --primary-hover: #d97706;
      --success: #10b981;
      --error: #ef4444;
      --code-bg: #0b1120;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { background: var(--bg); color: var(--text); display: flex; height: 100vh; overflow: hidden; }
    #sidebar { width: 240px; background: #090d16; border-right: 1px solid var(--border); display: flex; flex-direction: column; }
    #sidebar .logo { padding: 1.25rem 1rem; font-size: 1.15rem; font-weight: 700; color: var(--primary); border-bottom: 1px solid var(--border); display: flex; align-items: center; gap: 8px; }
    #sidebar nav { flex: 1; padding: 0.75rem 0.5rem; display: flex; flex-direction: column; gap: 4px; overflow-y: auto; }
    .nav-btn { display: flex; align-items: center; gap: 10px; width: 100%; padding: 0.65rem 0.85rem; border: none; background: transparent; color: var(--muted); border-radius: 6px; cursor: pointer; text-align: left; font-size: 0.9rem; font-weight: 500; transition: all 0.15s; }
    .nav-btn:hover { background: var(--card); color: var(--text); }
    .nav-btn.active { background: var(--primary); color: #000; font-weight: 600; }
    #main { flex: 1; display: flex; flex-direction: column; overflow: hidden; background: var(--bg); }
    header { height: 56px; border-bottom: 1px solid var(--border); padding: 0 1.5rem; display: flex; align-items: center; justify-content: space-between; background: var(--card); }
    header h2 { font-size: 1.1rem; font-weight: 600; }
    .content-area { flex: 1; padding: 1.5rem; overflow-y: auto; display: none; }
    .content-area.active { display: block; }
    .card { background: var(--card); border: 1px solid var(--border); border-radius: 8px; padding: 1.25rem; margin-bottom: 1.25rem; }
    .card h3 { font-size: 1rem; margin-bottom: 0.75rem; color: var(--primary); }
    .form-group { margin-bottom: 1rem; }
    .form-group label { display: block; font-size: 0.85rem; color: var(--muted); margin-bottom: 0.35rem; }
    input, select, textarea { width: 100%; background: var(--code-bg); border: 1px solid var(--border); color: var(--text); padding: 0.6rem 0.75rem; border-radius: 6px; font-size: 0.9rem; }
    textarea { font-family: monospace; resize: vertical; }
    .btn { background: var(--primary); color: #000; font-weight: 600; padding: 0.6rem 1.2rem; border-radius: 6px; border: none; cursor: pointer; transition: 0.15s; }
    .btn:hover { background: var(--primary-hover); }
    .btn-secondary { background: #334155; color: var(--text); margin-left: 0.5rem; }
    .btn-secondary:hover { background: #475569; }
    table { width: 100%; border-collapse: collapse; margin-top: 0.5rem; }
    th, td { padding: 0.65rem 0.75rem; border: 1px solid var(--border); text-align: left; font-size: 0.85rem; }
    th { background: #131d2e; color: var(--muted); }
    .badge { display: inline-block; padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.75rem; font-weight: 600; }
    .badge-online { background: #064e3b; color: #34d399; }
    .badge-offline { background: #451a1a; color: #f87171; }
    #editor-container { display: flex; flex-direction: column; height: 500px; border: 1px solid var(--border); border-radius: 6px; overflow: hidden; }
    #editor-toolbar { background: #131d2e; padding: 0.5rem 0.75rem; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid var(--border); }
    #editor-area { flex: 1; width: 100%; height: 100%; background: var(--code-bg); color: #e2e8f0; font-family: "Cascadia Code", Consolas, monospace; padding: 0.75rem; border: none; outline: none; line-height: 1.5; font-size: 0.9rem; }
    .json-key { color: #f59e0b; }
    .json-str { color: #34d399; }
    .json-num { color: #60a5fa; }
    .json-bool { color: #f472b6; }
  </style>
</head>
<body>
  <div id="sidebar">
    <div class="logo">⚡ GitMap Fleet UI</div>
    <nav>
      <button class="nav-btn active" onclick="showTab('settings')">⚙️ Settings</button>
      <button class="nav-btn" onclick="showTab('commitin')">📦 Commitin</button>
      <button class="nav-btn" onclick="showTab('ssh')">🖥️ SSH Fleet</button>
      <button class="nav-btn" onclick="showTab('macro')">📜 Macros</button>
      <button class="nav-btn" onclick="showTab('installer')">🚀 Installers</button>
      <button class="nav-btn" onclick="showTab('prompts')">💬 Prompts & AI</button>
      <button class="nav-btn" onclick="showTab('import-export')">🔄 Import / Export</button>
      <button class="nav-btn" onclick="showTab('schedules')">⏱️ Schedules</button>
      <button class="nav-btn" onclick="showTab('editor')">📝 Remote Editor</button>
    </nav>
  </div>
  <div id="main">
    <header>
      <h2 id="header-title">Settings</h2>
      <div id="status-indicator" style="font-size: 0.85rem; color: var(--muted);">Connected: localhost</div>
    </header>

    <!-- SETTINGS TAB -->
    <div id="tab-settings" class="content-area active">
      <div class="card">
        <h3>General Settings</h3>
        <div class="form-group">
          <label>Theme</label>
          <select id="setting-theme"><option value="dark">Dark Theme (Standard)</option><option value="light">Light Theme</option></select>
        </div>
        <div class="form-group">
          <label>Default Remote Target</label>
          <input type="text" id="setting-remote" placeholder="origin or alias">
        </div>
        <div class="form-group">
          <label>Cluster REST Port</label>
          <input type="number" id="setting-port" value="49152">
        </div>
        <button class="btn" onclick="saveSettings()">Save Settings</button>
      </div>
    </div>

    <!-- COMMITIN TAB -->
    <div id="tab-commitin" class="content-area">
      <div class="card">
        <h3>Commitin Multi-Option Engine</h3>
        <div class="form-group">
          <label>Commit Message</label>
          <input type="text" id="commit-msg" placeholder="feat(scope): concise description">
        </div>
        <div class="form-group">
          <label>Commit Direction / Flow</label>
          <select id="commit-direction">
            <option value="standard">Standard Commit</option>
            <option value="right">Commit Right (Local -> Remote Stage)</option>
            <option value="left">Commit Left (Remote Rebase -> Local)</option>
          </select>
        </div>
        <div class="form-group">
          <label><input type="checkbox" id="commit-amend" style="width: auto;"> Amend Previous Commit</label>
        </div>
        <div class="form-group">
          <label><input type="checkbox" id="commit-push" checked style="width: auto;"> Auto Push Immediately</label>
        </div>
        <button class="btn" onclick="execCommitin()">Execute Commitin</button>
        <div id="commit-output" style="margin-top: 1rem; font-family: monospace; white-space: pre-wrap; font-size: 0.85rem; color: var(--muted);"></div>
      </div>
    </div>

    <!-- SSH TAB -->
    <div id="tab-ssh" class="content-area">
      <div class="card">
        <h3>SSH Cluster Nodes & Fleet Operations</h3>
        <button class="btn" onclick="refreshNodes()">Refresh Node Fleet</button>
        <button class="btn btn-secondary" onclick="deployMacroFleetModal()">Deploy Macros to Fleet</button>
        <table id="nodes-table" style="margin-top: 1rem;">
          <thead><tr><th>Node ID</th><th>Alias</th><th>Host / IP</th><th>OS</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody id="nodes-body"><tr><td colspan="6" style="text-align:center;">Loading nodes...</td></tr></tbody>
        </table>
      </div>
    </div>

    <!-- MACRO TAB -->
    <div id="tab-macro" class="content-area">
      <div class="card">
        <h3>Macro Automation Builder</h3>
        <button class="btn" onclick="loadMacros()">Reload Macros</button>
        <div id="macros-list" style="margin-top: 1rem;"></div>
      </div>
    </div>

    <!-- INSTALLER TAB -->
    <div id="tab-installer" class="content-area">
      <div class="card">
        <h3>Custom Multi-OS Installer Manager</h3>
        <div class="form-group"><label>Installer Name</label><input type="text" id="inst-name" placeholder="e.g. agy-tools"></div>
        <div class="form-group"><label>Target OS</label><select id="inst-os"><option value="windows">Windows</option><option value="ubuntu">Ubuntu / Debian</option><option value="centos">CentOS / RHEL</option><option value="unix">Universal Unix</option></select></div>
        <div class="form-group"><label>Install Command</label><input type="text" id="inst-cmd" placeholder="curl -fsSL ... | sh"></div>
        <div class="form-group"><label>Target Node for Immediate Test</label><select id="inst-node-select"><option value="local">Local Machine</option></select></div>
        <button class="btn" onclick="saveInstaller()">Save Installer</button>
        <button class="btn btn-secondary" onclick="testInstaller()">Test on Selected Node</button>
      </div>
    </div>

    <!-- PROMPTS TAB -->
    <div id="tab-prompts" class="content-area">
      <div class="card">
        <h3>Prompts & AI Instructions Manager</h3>
        <div class="form-group"><label>Prompt Template</label><textarea id="prompt-content" rows="10" placeholder="Paste or edit prompt markdown/JSON..."></textarea></div>
        <button class="btn" onclick="formatPromptJSON()">Format JSON</button>
        <button class="btn btn-secondary" onclick="importPrompt()">Import Prompt</button>
      </div>
    </div>

    <!-- IMPORT-EXPORT TAB -->
    <div id="tab-import-export" class="content-area">
      <div class="card">
        <h3>System Configuration Import / Export</h3>
        <button class="btn" onclick="exportFullConfig()">Export Full JSON Config</button>
        <button class="btn btn-secondary" onclick="importFullConfig()">Import JSON Config</button>
        <div style="margin-top: 1rem;"><textarea id="import-export-area" rows="12" placeholder="JSON snapshot data..."></textarea></div>
      </div>
    </div>

    <!-- SCHEDULES TAB -->
    <div id="tab-schedules" class="content-area">
      <div class="card">
        <h3>Task & Cron Schedules</h3>
        <div class="form-group"><label>Schedule Name</label><input type="text" id="sched-name" placeholder="e.g. night-sync"></div>
        <div class="form-group"><label>Cron Expression</label><input type="text" id="sched-cron" placeholder="0 2 * * *"></div>
        <button class="btn" onclick="addSchedule()">Add Schedule</button>
      </div>
    </div>

    <!-- REMOTE EDITOR TAB -->
    <div id="tab-editor" class="content-area">
      <div class="card" style="margin-bottom: 0.5rem;">
        <div style="display: flex; gap: 10px; align-items: flex-end;">
          <div style="flex: 1;"><label style="font-size:0.8rem;color:var(--muted)">Target Node</label><select id="editor-node"><option value="local">local (current machine)</option></select></div>
          <div style="flex: 3;"><label style="font-size:0.8rem;color:var(--muted)">Absolute / Relative File Path</label><input type="text" id="editor-path" placeholder="/etc/gitmap.conf or cli/main.go"></div>
          <button class="btn" onclick="openRemoteFile()">Load File</button>
          <button class="btn" style="background:var(--success)" onclick="saveRemoteFile()">Save to Host</button>
        </div>
      </div>
      <div id="editor-container">
        <div id="editor-toolbar"><span id="editor-status" style="font-size:0.85rem;color:var(--muted)">No file loaded</span><span id="editor-lang" style="font-size:0.8rem;color:var(--primary)">Language: Plaintext</span></div>
        <textarea id="editor-area" spellcheck="false" placeholder="Remote file content will appear here..."></textarea>
      </div>
    </div>
  </div>

  <script>
    function showTab(tabId) {
      document.querySelectorAll('.content-area').forEach(el => el.classList.remove('active'));
      document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));
      const target = document.getElementById('tab-' + tabId);
      if (target) target.classList.add('active');
      document.getElementById('header-title').innerText = tabId.charAt(0).toUpperCase() + tabId.slice(1);
      event?.target?.classList.add('active');
      if (tabId === 'ssh') refreshNodes();
      if (tabId === 'editor') populateEditorNodes();
    }

    async function refreshNodes() {
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
        sel.innerHTML = '<option value="local">local (current machine)</option>';
        if (data && data.length) {
          data.forEach(n => {
            const opt = document.createElement('option');
            opt.value = n.alias;
            opt.innerText = n.alias + ' (' + n.host + ')';
            sel.appendChild(opt);
          });
        }
      } catch (e) {}
    }

    async function openRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      if (!path) return alert('Please enter a file path');
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
          alert('Failed to read file: ' + (data.error || 'unknown error'));
        }
      } catch (e) { alert('Error reading file: ' + e.message); }
    }

    async function saveRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      const content = document.getElementById('editor-area').value;
      if (!path) return alert('Please enter a file path');
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
          alert('Failed to save file: ' + (data.error || 'unknown error'));
        }
      } catch (e) { alert('Error saving file: ' + e.message); }
    }

    async function execCommitin() {
      const msg = document.getElementById('commit-msg').value;
      const dir = document.getElementById('commit-direction').value;
      const isAmend = document.getElementById('commit-amend').checked;
      const isPush = document.getElementById('commit-push').checked;
      const out = document.getElementById('commit-output');
      out.innerText = 'Executing commitin...';
      try {
        const res = await fetch('/api/commitin/exec', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({message: msg, direction: dir, isAmend: isAmend, isPush: isPush})
        });
        const data = await res.json();
        out.innerText = data.output || data.error || 'Done';
      } catch (e) { out.innerText = 'Error: ' + e.message; }
    }

    // Auto-detect route on load
    window.addEventListener('load', () => {
      const path = window.location.pathname.replace(/^\//, '');
      if (path && document.getElementById('tab-' + path)) {
        showTab(path);
      }
    });
  </script>
</body>
</html>
`

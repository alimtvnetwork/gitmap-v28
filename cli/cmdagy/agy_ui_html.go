package cmdagy

const agyUIDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GitMap — Antigravity (AGY) Studio & Prompt Manager</title>
<style>
  :root {
    --bg: #0d1117;
    --surface: #161b22;
    --surface-hover: #1c2128;
    --border: #30363d;
    --text: #e6edf3;
    --dim: #8b949e;
    --cyan: #58a6ff;
    --green: #2ea043;
    --green-bg: rgba(46, 160, 67, 0.15);
    --yellow: #d29922;
    --yellow-bg: rgba(210, 153, 34, 0.15);
    --red: #f85149;
    --purple: #bc8cff;
  }
  * { box-sizing: border-box; }
  body { background: var(--bg); color: var(--text); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace; margin: 0; padding: 20px; font-size: 13px; line-height: 1.5; }
  .container { max-width: 1200px; margin: 0 auto; }
  header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); padding-bottom: 16px; margin-bottom: 20px; }
  .title-group { display: flex; align-items: center; gap: 12px; }
  h1 { font-size: 18px; font-weight: 700; color: var(--cyan); margin: 0; }
  .badge { font-size: 11px; font-weight: 600; padding: 3px 8px; border-radius: 12px; }
  .badge.online { background: var(--green-bg); color: var(--green); border: 1px solid var(--green); }
  .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
  @media (max-width: 900px) { .grid { grid-template-columns: 1fr; } }
  .card { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; padding: 16px; margin-bottom: 16px; }
  .card-title { font-size: 14px; font-weight: 600; color: var(--text); margin-top: 0; margin-bottom: 12px; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); padding-bottom: 8px; }
  button { background: var(--surface); color: var(--text); border: 1px solid var(--border); border-radius: 6px; padding: 6px 12px; font-size: 12px; font-weight: 600; cursor: pointer; transition: 0.15s; }
  button:hover { border-color: var(--cyan); background: var(--surface-hover); }
  button.primary { background: #238636; border-color: #2ea043; color: #fff; }
  button.primary:hover { background: #2ea043; }
  button.secondary { background: #1f6feb; border-color: #388bfd; color: #fff; }
  button.secondary:hover { background: #388bfd; }
  button.accent { background: #6e40c9; border-color: #8957e5; color: #fff; }
  button.accent:hover { background: #8957e5; }
  .form-group { margin-bottom: 12px; }
  label { display: block; font-size: 11px; text-transform: uppercase; color: var(--dim); margin-bottom: 4px; font-weight: 600; }
  input[type="text"], textarea, select { width: 100%; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; padding: 8px 10px; color: var(--text); font-family: monospace; font-size: 12px; }
  input[type="text"]:focus, textarea:focus, select:focus { outline: none; border-color: var(--cyan); }
  textarea { min-height: 120px; resize: vertical; }
  .btn-row { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 12px; }
  table { width: 100%; border-collapse: collapse; margin-top: 8px; }
  th, td { padding: 8px 10px; text-align: left; border-bottom: 1px solid var(--border); font-size: 12px; }
  th { color: var(--dim); font-size: 11px; text-transform: uppercase; }
  .status-tag { display: inline-block; font-size: 10px; font-weight: 700; padding: 2px 6px; border-radius: 4px; text-transform: uppercase; }
  .status-tag.running { background: var(--yellow-bg); color: var(--yellow); }
  .status-tag.queued { background: rgba(88, 166, 255, 0.15); color: var(--cyan); }
  .status-tag.saved { background: rgba(188, 140, 255, 0.15); color: var(--purple); }
  .snippet { font-family: monospace; color: var(--dim); max-width: 320px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .toast { position: fixed; bottom: 20px; right: 20px; background: #238636; color: #fff; padding: 10px 16px; border-radius: 6px; display: none; font-weight: 600; box-shadow: 0 4px 12px rgba(0,0,0,0.5); z-index: 100; }
</style>
</head>
<body>
<div class="container">
  <header>
    <div class="title-group">
      <h1>GitMap — Antigravity Studio</h1>
      <span class="badge online" id="nodeBadge">Node: Local</span>
      <span class="badge online">HTTP Server Active</span>
    </div>
    <div>
      <button onclick="fetchStatus()">↻ Refresh</button>
    </div>
  </header>

  <div class="grid">
    <div>
      <div class="card">
        <h2 class="card-title">Prompt Dispatcher & Lifecycle Manager</h2>
        <div class="form-group">
          <label>Target Project (Path or Directory)</label>
          <input type="text" id="targetProject" placeholder=". or D:\work\gitmap or /home/work/repo">
        </div>
        <div class="form-group">
          <label>Prompt Title</label>
          <input type="text" id="promptTitle" placeholder="Bugfix / Feature verification">
        </div>
        <div class="form-group">
          <label>Prompt Directive Text</label>
          <textarea id="promptText" placeholder="Enter comprehensive prompt instructions..."></textarea>
        </div>
        <div class="btn-row">
          <button class="primary" onclick="sendPromptImmediate()">▶ Dispatch Immediate</button>
          <button class="secondary" onclick="enqueuePrompt()">⏳ Enqueue to Project</button>
          <button class="accent" onclick="savePromptTemplate()">💾 Save as Template</button>
        </div>
      </div>

      <div class="card">
        <h2 class="card-title">Queued Prompts (<span id="queueCount">0</span>)</h2>
        <div style="max-height: 240px; overflow-y: auto;">
          <table>
            <thead><tr><th>ID</th><th>Title</th><th>Created</th><th>Status</th></tr></thead>
            <tbody id="queueTableBody"><tr><td colspan="4" style="color:var(--dim)">No queued prompts</td></tr></tbody>
          </table>
        </div>
      </div>
    </div>

    <div>
      <div class="card">
        <h2 class="card-title">Running Projects & Active Prompts (<span id="runningCount">0</span>)</h2>
        <div style="max-height: 280px; overflow-y: auto;">
          <table>
            <thead><tr><th>Project</th><th>Status</th><th>Active Prompt</th><th>Action</th></tr></thead>
            <tbody id="projectsTableBody"><tr><td colspan="4" style="color:var(--dim)">Scanning active projects...</td></tr></tbody>
          </table>
        </div>
      </div>

      <div class="card">
        <h2 class="card-title">Saved Templates & Recent History</h2>
        <div style="max-height: 220px; overflow-y: auto;">
          <table>
            <thead><tr><th>Title</th><th>Target</th><th>Type</th><th>Action</th></tr></thead>
            <tbody id="historyTableBody"><tr><td colspan="4" style="color:var(--dim)">No saved templates</td></tr></tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="toast" id="toastMsg">Action succeeded!</div>

<script>
  function showToast(msg) {
    const t = document.getElementById('toastMsg');
    t.innerText = msg;
    t.style.display = 'block';
    setTimeout(() => { t.style.display = 'none'; }, 2500);
  }

  function fetchStatus() {
    fetch('/api/status')
      .then(res => res.json())
      .then(data => {
        document.getElementById('nodeBadge').innerText = 'Node: ' + (data.nodeAlias || 'Local');
        renderRunningProjects(data.runningProjects || []);
        renderQueue(data.queuedPrompts || []);
        renderSavedPrompts(data.savedPrompts || []);
      })
      .catch(err => console.error('Status fetch failed:', err));
  }

  function renderRunningProjects(projects) {
    document.getElementById('runningCount').innerText = projects.length;
    const tbody = document.getElementById('projectsTableBody');
    if (!projects.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No active Antigravity projects</td></tr>';
      return;
    }
    tbody.innerHTML = projects.map(p => `
      <tr>
        <td><strong>${escapeHtml(p.projectName)}</strong></td>
        <td><span class="status-tag ${p.hasActivePrompt ? 'running' : 'saved'}">${p.status || 'IDLE'}</span></td>
        <td class="snippet" title="${escapeHtml(p.promptPreview || '')}">${escapeHtml(p.promptPreview || 'None')}</td>
        <td><button onclick="selectProject('${escapeHtml(p.projectPath)}')">Target</button></td>
      </tr>
    `).join('');
  }

  function renderQueue(queued) {
    document.getElementById('queueCount').innerText = queued.length;
    const tbody = document.getElementById('queueTableBody');
    if (!queued.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No prompts in queue</td></tr>';
      return;
    }
    tbody.innerHTML = queued.map(q => `
      <tr>
        <td>#${q.id}</td>
        <td>${escapeHtml(q.title || 'Untitled')}</td>
        <td>${escapeHtml(q.createdAt ? q.createdAt.split('T')[0] : '')}</td>
        <td><span class="status-tag queued">${q.status || 'QUEUED'}</span></td>
      </tr>
    `).join('');
  }

  function renderSavedPrompts(saved) {
    const tbody = document.getElementById('historyTableBody');
    if (!saved.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No saved templates</td></tr>';
      return;
    }
    tbody.innerHTML = saved.map(s => `
      <tr>
        <td><strong>${escapeHtml(s.title || 'Template')}</strong></td>
        <td>${escapeHtml(s.projectTarget || '.')}</td>
        <td><span class="status-tag saved">${s.status || 'SAVED'}</span></td>
        <td><button onclick="loadTemplate('${escapeHtml(s.projectTarget || '')}', '${escapeHtml(s.title || '')}', '${escapeHtml(s.promptText || '')}')">Load</button></td>
      </tr>
    `).join('');
  }

  function selectProject(path) {
    document.getElementById('targetProject').value = path;
  }

  function loadTemplate(target, title, text) {
    if (target) document.getElementById('targetProject').value = target;
    if (title) document.getElementById('promptTitle').value = title;
    if (text) document.getElementById('promptText').value = text;
    showToast('Loaded template: ' + title);
  }

  function sendPromptImmediate() {
    const payload = {
      projectTarget: document.getElementById('targetProject').value,
      title: document.getElementById('promptTitle').value,
      promptText: document.getElementById('promptText').value
    };
    if (!payload.promptText) { alert('Prompt Directive Text cannot be empty'); return; }
    fetch('/api/prompts/send', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(res => res.json()).then(res => {
      showToast('Prompt dispatched successfully!');
      fetchStatus();
    });
  }

  function enqueuePrompt() {
    const payload = {
      projectTarget: document.getElementById('targetProject').value,
      title: document.getElementById('promptTitle').value,
      promptText: document.getElementById('promptText').value,
      isEnqueue: true
    };
    if (!payload.promptText) { alert('Prompt Directive Text cannot be empty'); return; }
    fetch('/api/prompts/enqueue', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(res => res.json()).then(res => {
      showToast('Prompt enqueued to project!');
      fetchStatus();
    });
  }

  function savePromptTemplate() {
    const payload = {
      projectTarget: document.getElementById('targetProject').value,
      title: document.getElementById('promptTitle').value || 'Saved Template',
      promptText: document.getElementById('promptText').value
    };
    if (!payload.promptText) { alert('Prompt Directive Text cannot be empty'); return; }
    fetch('/api/prompts/save', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    }).then(res => res.json()).then(res => {
      showToast('Prompt template saved!');
      fetchStatus();
    });
  }

  function escapeHtml(str) {
    return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  fetchStatus();
  setInterval(fetchStatus, 5000);
</script>
</body>
</html>
`

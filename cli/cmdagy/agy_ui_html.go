package cmdagy

const agyUIDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GitMap — Antigravity (AGY) Studio & Fleet Manager</title>
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
    --red-bg: rgba(248, 81, 73, 0.15);
    --purple: #bc8cff;
  }
  * { box-sizing: border-box; }
  body { background: var(--bg); color: var(--text); font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace; margin: 0; padding: 20px; font-size: 13px; line-height: 1.5; }
  .container { max-width: 1280px; margin: 0 auto; }
  header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border); padding-bottom: 16px; margin-bottom: 16px; }
  .title-group { display: flex; align-items: center; gap: 12px; }
  h1 { font-size: 18px; font-weight: 700; color: var(--cyan); margin: 0; }
  .badge { font-size: 11px; font-weight: 600; padding: 3px 8px; border-radius: 12px; }
  .badge.online { background: var(--green-bg); color: var(--green); border: 1px solid var(--green); }
  .fleet-strip { display: flex; gap: 10px; overflow-x: auto; padding: 10px 14px; background: var(--surface); border: 1px solid var(--border); border-radius: 8px; margin-bottom: 16px; align-items: center; }
  .fleet-strip-label { font-size: 11px; font-weight: 700; text-transform: uppercase; color: var(--dim); margin-right: 8px; white-space: nowrap; }
  .fleet-node-card { display: flex; align-items: center; gap: 6px; padding: 4px 10px; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; font-size: 12px; white-space: nowrap; }
  .fleet-node-card.online { border-color: var(--green); }
  .fleet-node-card.offline { border-color: var(--border); opacity: 0.7; }
  .dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
  .dot.online { background: var(--green); }
  .dot.offline { background: var(--dim); }
  .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-bottom: 20px; }
  @media (max-width: 960px) { .grid { grid-template-columns: 1fr; } }
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
  textarea { min-height: 110px; resize: vertical; }
  .btn-row { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 12px; }
  table { width: 100%; border-collapse: collapse; margin-top: 8px; }
  th, td { padding: 8px 10px; text-align: left; border-bottom: 1px solid var(--border); font-size: 12px; }
  th { color: var(--dim); font-size: 11px; text-transform: uppercase; }
  .status-tag { display: inline-block; font-size: 10px; font-weight: 700; padding: 2px 6px; border-radius: 4px; text-transform: uppercase; }
  .status-tag.running { background: var(--yellow-bg); color: var(--yellow); }
  .status-tag.queued { background: rgba(88, 166, 255, 0.15); color: var(--cyan); }
  .status-tag.saved { background: rgba(188, 140, 255, 0.15); color: var(--purple); }
  .status-tag.completed { background: var(--green-bg); color: var(--green); }
  .snippet { font-family: monospace; color: var(--dim); max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .toast { position: fixed; bottom: 20px; right: 20px; background: #238636; color: #fff; padding: 10px 16px; border-radius: 6px; display: none; font-weight: 600; box-shadow: 0 4px 12px rgba(0,0,0,0.5); z-index: 100; }
</style>
</head>
<body>
<div class="container">
  <header>
    <div class="title-group">
      <h1>GitMap — Antigravity Studio</h1>
      <span class="badge online" id="nodeBadge">Node: Local</span>
      <span class="badge online">HTTP Active (:7430)</span>
    </div>
    <div style="display:flex; gap:8px; align-items:center;">
      <label style="margin:0; text-transform:none; font-size:12px; cursor:pointer;">
        <input type="checkbox" id="autoRefresh" checked style="vertical-align:middle;"> Auto-Refresh (5s)
      </label>
      <button onclick="fetchStatus()">↻ Refresh</button>
    </div>
  </header>

  <div class="fleet-strip" id="fleetStrip">
    <span class="fleet-strip-label">Fleet Nodes:</span>
    <span style="color:var(--dim); font-size:12px;">Discovering fleet nodes...</span>
  </div>

  <div class="grid">
    <div>
      <div class="card">
        <h2 class="card-title">Prompt Dispatcher & Lifecycle Manager</h2>
        <div style="display:grid; grid-template-columns: 1fr 1fr; gap:10px;">
          <div class="form-group">
            <label>Target Fleet Node</label>
            <select id="targetNode">
              <option value="local">local (Current Machine)</option>
            </select>
          </div>
          <div class="form-group">
            <label>Project Slug / Workspace</label>
            <input type="text" id="targetProject" placeholder=". or gitmap or repo-slug">
          </div>
        </div>
        <div class="form-group">
          <label>Prompt Title</label>
          <input type="text" id="promptTitle" placeholder="Code inspection / Feature verification">
        </div>
        <div class="form-group">
          <label>Prompt Directive Text</label>
          <textarea id="promptText" placeholder="Enter comprehensive prompt instructions..."></textarea>
        </div>
        <div class="btn-row">
          <button class="primary" onclick="sendPromptImmediate()">🚀 Dispatch Immediate</button>
          <button class="secondary" onclick="enqueuePrompt()">⏳ Enqueue to Project</button>
          <button class="accent" onclick="savePromptTemplate()">💾 Save as Template</button>
        </div>
      </div>

      <div class="card">
        <h2 class="card-title">Queued Prompts (<span id="queueCount">0</span>)</h2>
        <div style="max-height: 240px; overflow-y: auto;">
          <table>
            <thead><tr><th>ID</th><th>Title</th><th>Project</th><th>Status</th></tr></thead>
            <tbody id="queueTableBody"><tr><td colspan="4" style="color:var(--dim)">No queued prompts</td></tr></tbody>
          </table>
        </div>
      </div>
    </div>

    <div>
      <div class="card">
        <h2 class="card-title">
          <span>Active Projects & Prompts (<span id="runningCount">0</span>)</span>
          <button style="padding:2px 8px; font-size:11px;" onclick="loadPromptTree()">🌳 Tree View</button>
        </h2>
        <div style="max-height: 280px; overflow-y: auto;">
          <table>
            <thead><tr><th>Project</th><th>Status</th><th>Active Prompt</th><th>Action</th></tr></thead>
            <tbody id="projectsTableBody"><tr><td colspan="4" style="color:var(--dim)">Scanning active projects...</td></tr></tbody>
          </table>
        </div>
      </div>

      <div class="card">
        <h2 class="card-title">Saved Templates & Recent History</h2>
        <div style="max-height: 230px; overflow-y: auto;">
          <table>
            <thead><tr><th>Title</th><th>Target</th><th>Type</th><th>Actions</th></tr></thead>
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
        renderFleetNodes(data.fleetNodes || []);
        renderRunningProjects(data.runningProjects || []);
        renderQueue(data.queuedPrompts || []);
        renderSavedPrompts(data.savedPrompts || []);
      })
      .catch(err => console.error('Status fetch failed:', err));
  }

  function renderFleetNodes(nodes) {
    const strip = document.getElementById('fleetStrip');
    const select = document.getElementById('targetNode');
    const currentVal = select.value;

    if (!nodes.length) {
      strip.innerHTML = '<span class="fleet-strip-label">Fleet Nodes:</span><span style="color:var(--dim); font-size:12px;">Local machine only</span>';
      return;
    }

    strip.innerHTML = '<span class="fleet-strip-label">Fleet Nodes:</span>' + nodes.map(n => 
      '<div class="fleet-node-card ' + (n.isOnline ? 'online' : 'offline') + '">' +
        '<span class="dot ' + (n.isOnline ? 'online' : 'offline') + '"></span>' +
        '<strong>' + escapeHtml(n.alias || n.hostAlias) + '</strong>' +
        '<span style="color:var(--dim); font-size:11px;">(' + escapeHtml(n.os || 'os') + ')</span>' +
      '</div>'
    ).join('');

    select.innerHTML = nodes.map(n => 
      '<option value="' + escapeHtml(n.alias) + '">' + escapeHtml(n.alias) + ' (' + escapeHtml(n.hostAlias || 'host') + ')</option>'
    ).join('');

    if (currentVal) {
      select.value = currentVal;
    }
  }

  function renderRunningProjects(projects) {
    document.getElementById('runningCount').innerText = projects.length;
    const tbody = document.getElementById('projectsTableBody');
    if (!projects.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No active Antigravity projects</td></tr>';
      return;
    }
    tbody.innerHTML = projects.map(p => 
      '<tr>' +
        '<td><strong>' + escapeHtml(p.projectName) + '</strong></td>' +
        '<td><span class="status-tag ' + (p.hasActivePrompt ? 'running' : 'saved') + '">' + (p.status || 'IDLE') + '</span></td>' +
        '<td class="snippet" title="' + escapeHtml(p.promptPreview || '') + '">' + escapeHtml(p.promptPreview || 'None') + '</td>' +
        '<td><button onclick="selectProject(\'' + escapeHtml(p.projectName || p.projectPath) + '\')">Target</button></td>' +
      '</tr>'
    ).join('');
  }

  function renderQueue(queued) {
    document.getElementById('queueCount').innerText = queued.length;
    const tbody = document.getElementById('queueTableBody');
    if (!queued.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No prompts in queue</td></tr>';
      return;
    }
    tbody.innerHTML = queued.map(q => 
      '<tr>' +
        '<td>#' + q.id + '</td>' +
        '<td>' + escapeHtml(q.title || 'Untitled') + '</td>' +
        '<td>' + escapeHtml(q.projectName || '-') + '</td>' +
        '<td><span class="status-tag queued">' + (q.status || 'QUEUED') + '</span></td>' +
      '</tr>'
    ).join('');
  }

  function renderSavedPrompts(saved) {
    const tbody = document.getElementById('historyTableBody');
    if (!saved.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No saved templates</td></tr>';
      return;
    }
    tbody.innerHTML = saved.map(s => 
      '<tr>' +
        '<td><strong>' + escapeHtml(s.title || 'Template') + '</strong></td>' +
        '<td>' + escapeHtml(s.projectTarget || '.') + '</td>' +
        '<td><span class="status-tag saved">' + (s.status || 'SAVED') + '</span></td>' +
        '<td>' +
          '<button onclick="loadTemplate(\'' + escapeHtml(s.projectTarget || '') + '\', \'' + escapeHtml(s.title || '') + '\', \'' + escapeHtml(s.promptText || '') + '\')">Load</button> ' +
          '<button onclick="resendPrompt(\'' + escapeHtml(s.id) + '\')">↻ Resend</button>' +
        '</td>' +
      '</tr>'
    ).join('');
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
    }).then(res => res.json()).then(() => {
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
    }).then(res => res.json()).then(() => {
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
    }).then(res => res.json()).then(() => {
      showToast('Prompt template saved!');
      fetchStatus();
    });
  }

  function resendPrompt(id) {
    fetch('/api/prompts/resend', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: id })
    }).then(res => res.json()).then(() => {
      showToast('Prompt resent successfully!');
      fetchStatus();
    });
  }

  function loadPromptTree() {
    fetch('/api/prompts/tree')
      .then(res => res.json())
      .then(data => {
        let msg = 'Project Prompts Tree:\n\n';
        data.forEach(item => {
          msg += '• [' + item.nodeAlias + '] ' + item.projectName + ' (' + item.workspacePath + ')\n';
          (item.activePrompts || []).forEach(ap => {
            msg += '    - ACTIVE: ' + ap.promptSnippet + ' (' + ap.status + ')\n';
          });
          (item.queuedPrompts || []).forEach(qp => {
            msg += '    - QUEUED: ' + qp.title + ' (#' + qp.id + ')\n';
          });
        });
        alert(msg);
      });
  }

  function escapeHtml(str) {
    return String(str || '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  fetchStatus();
  setInterval(() => {
    if (document.getElementById('autoRefresh').checked) {
      fetchStatus();
    }
  }, 5000);
</script>
</body>
</html>
`

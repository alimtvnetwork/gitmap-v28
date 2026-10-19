package cmdagy

const agyUIDashboardScript = `<script>
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
    const sorted = projects.slice().sort(function(a, b) {
      function rank(p) {
        const st = ((p && p.status) || (p && p.hasActivePrompt ? 'running' : 'idle')).toLowerCase();
        if (st === 'running') return 0;
        if (st === 'queued' || st === 'enqueued') return 1;
        if (st === 'active') return 2;
        return 3;
      }
      return rank(a) - rank(b);
    });
    tbody.innerHTML = sorted.map(p => {
      const isRunning = p.hasActivePrompt || (p.status || '').toLowerCase() === 'running';
      const badgeClass = isRunning ? 'running' : 'saved';
      const badgeText = isRunning ? 'RUNNING' : (p.status || 'IDLE').toUpperCase();
      return '<tr>' +
        '<td><strong>' + escapeHtml(p.projectName) + '</strong></td>' +
        '<td><span class="status-tag ' + badgeClass + '">' + escapeHtml(badgeText) + '</span></td>' +
        '<td class="snippet" title="' + escapeHtml(p.promptPreview || '') + '">' + escapeHtml(p.promptPreview || 'None') + '</td>' +
        '<td><button onclick="selectProject(\'' + escapeHtml(p.projectName || p.projectPath) + '\')">Target</button></td>' +
      '</tr>';
    }).join('');
  }

  function renderQueue(queued) {
    document.getElementById('queueCount').innerText = queued.length;
    const tbody = document.getElementById('queueTableBody');
    if (!queued.length) {
      tbody.innerHTML = '<tr><td colspan="4" style="color:var(--dim)">No prompts in queue</td></tr>';
      return;
    }
    const sorted = queued.slice().sort(function(a, b) {
      function rank(q) {
        const st = ((q && q.status) || 'queued').toLowerCase();
        if (st === 'running') return 0;
        if (st === 'queued' || st === 'enqueued') return 1;
        if (st === 'active') return 2;
        return 3;
      }
      return rank(a) - rank(b);
    });
    tbody.innerHTML = sorted.map(q => {
      const st = (q.status || 'QUEUED').toUpperCase();
      return '<tr>' +
        '<td>#' + q.id + '</td>' +
        '<td>' + escapeHtml(q.title || 'Untitled') + '</td>' +
        '<td>' + escapeHtml(q.projectName || '-') + '</td>' +
        '<td><span class="status-tag queued">' + escapeHtml(st) + '</span></td>' +
      '</tr>';
    }).join('');
  }

  function renderPromptsList(items) {
    const list = (items || []).slice();
    list.sort(function(a, b) {
      function rank(s) {
        const st = ((s && s.status) || '').toLowerCase();
        if (st === 'running') return 0;
        if (st === 'queued' || st === 'enqueued') return 1;
        if (st === 'active') return 2;
        return 3;
      }
      return rank(a) - rank(b);
    });
    return list;
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
    if (!payload.promptText) {
      alert('Prompt Directive Text cannot be empty');
      return;
    }
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
    if (!payload.promptText) {
      alert('Prompt Directive Text cannot be empty');
      return;
    }
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
    if (!payload.promptText) {
      alert('Prompt Directive Text cannot be empty');
      return;
    }
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

  let selectedTreeProject = '';

  function selectTreeProject(name) {
    selectedTreeProject = name;
    document.querySelectorAll('.tree-node-block').forEach(function(el) {
      el.classList.toggle('selected', el.id === 'treeNode_' + name);
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
</script>`

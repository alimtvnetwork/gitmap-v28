package cmdagy

const agyUITreeExtensionHTML = `
<style>
  .modal-overlay { position: fixed; inset: 0; background: rgba(0, 0, 0, 0.75); display: none; align-items: center; justify-content: center; z-index: 1000; padding: 20px; }
  .modal-overlay.active { display: flex; }
  .modal-content { background: var(--surface); border: 1px solid var(--border); border-radius: 8px; width: 100%; max-width: 860px; max-height: 85vh; display: flex; flex-direction: column; box-shadow: 0 8px 24px rgba(0,0,0,0.6); }
  .modal-header { display: flex; justify-content: space-between; align-items: center; padding: 12px 16px; border-bottom: 1px solid var(--border); }
  .modal-header h3 { margin: 0; font-size: 15px; color: var(--cyan); display: flex; align-items: center; gap: 8px; }
  .modal-actions { display: flex; gap: 6px; align-items: center; flex-wrap: wrap; }
  .modal-body { padding: 16px; overflow-y: auto; flex: 1; font-family: monospace; font-size: 12px; }
  .tree-node-block { margin-bottom: 12px; background: var(--bg); border: 1px solid var(--border); border-radius: 6px; padding: 10px 14px; transition: border-color 0.15s; }
  .tree-node-block:hover { border-color: var(--cyan); }
  .tree-node-block.selected { border-color: var(--cyan); background: rgba(88, 166, 255, 0.05); }
  .tree-node-header { font-weight: 700; color: var(--text); display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; cursor: pointer; }
  .tree-node-line { font-family: monospace; padding: 3px 0; display: flex; align-items: center; gap: 6px; color: var(--text); }
</style>

<div class="modal-overlay" id="treeModal" onclick="handleTreeOverlayClick(event)">
  <div class="modal-content">
    <div class="modal-header">
      <h3>🌳 Interactive Prompt Tree</h3>
      <div class="modal-actions">
        <button class="primary" onclick="treeActionSend()">🚀 Send</button>
        <button onclick="treeActionBackup()">💾 Backup</button>
        <button onclick="treeActionRestore()">📥 Restore</button>
        <button onclick="treeActionReset()">↺ Reset</button>
        <button onclick="loadPromptTree()">↻ Refresh</button>
        <button onclick="closePromptTree()" style="padding:2px 8px;">✕</button>
      </div>
    </div>
    <div class="modal-body" id="treeModalBody">
      <div style="color:var(--dim); padding:20px; text-align:center;">Loading prompt tree...</div>
    </div>
  </div>
</div>

<script>
  let selectedTreeProject = '';

  function selectTreeProject(name) {
    selectedTreeProject = name;
    document.querySelectorAll('.tree-node-block').forEach(function(el) {
      el.classList.toggle('selected', el.id === 'treeNode_' + name);
    });
  }

  function closePromptTree() {
    const modal = document.getElementById('treeModal');
    if (modal) modal.classList.remove('active');
  }

  function handleTreeOverlayClick(e) {
    if (e.target && e.target.id === 'treeModal') {
      closePromptTree();
    }
  }

  function treeActionSend() {
    if (selectedTreeProject) {
      document.getElementById('targetProject').value = selectedTreeProject;
    }
    closePromptTree();
    const area = document.getElementById('promptText');
    if (area) {
      area.focus();
      area.scrollIntoView({ behavior: 'smooth' });
    }
    showToast('Targeted project: ' + (selectedTreeProject || 'default'));
  }

  function treeActionBackup() {
    fetch('/api/prompts/backup', { method: 'POST' })
      .then(res => res.json())
      .then(data => {
        showToast('Backup created! Batch: ' + (data.batchId || 'done') + ' (' + (data.count || 0) + ' prompts)');
        loadPromptTree();
        fetchStatus();
      })
      .catch(err => alert('Backup failed: ' + err));
  }

  function treeActionRestore() {
    fetch('/api/prompts/restore', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({})
    })
      .then(res => res.json())
      .then(data => {
        showToast('Restored ' + (data.restored || 0) + ' prompt(s)');
        loadPromptTree();
        fetchStatus();
      })
      .catch(err => alert('Restore failed: ' + err));
  }

  function treeActionReset() {
    fetch('/api/prompts/reset', { method: 'POST' })
      .then(res => res.json())
      .then(data => {
        showToast('Queue & transient state reset successfully');
        loadPromptTree();
        fetchStatus();
      })
      .catch(err => alert('Reset failed: ' + err));
  }

  function loadPromptTree() {
    const modal = document.getElementById('treeModal');
    if (modal) modal.classList.add('active');
    const body = document.getElementById('treeModalBody');
    if (body) body.innerHTML = '<div style="color:var(--dim); padding:20px; text-align:center;">Scanning prompt tree...</div>';
    fetch('/api/prompts/refresh', { method: 'GET' })
      .then(res => res.json())
      .then(data => {
        renderPromptTree(data);
      })
      .catch(err => {
        if (body) body.innerHTML = '<div style="color:var(--red); padding:20px;">Failed to load tree: ' + escapeHtml(err.message) + '</div>';
      });
  }

  function renderPromptTree(data) {
    const body = document.getElementById('treeModalBody');
    if (!data || !data.length) {
      body.innerHTML = '<div style="color:var(--dim); padding:20px; text-align:center;">No project prompt trees discovered.</div>';
      return;
    }
    let html = '';
    data.forEach(function(item, idx) {
      if (!selectedTreeProject && idx === 0) selectedTreeProject = item.projectName;
      const isSel = selectedTreeProject === item.projectName;
      html += '<div class="tree-node-block ' + (isSel ? 'selected' : '') + '" id="treeNode_' + escapeHtml(item.projectName) + '">';
      html += '  <div class="tree-node-header" onclick="selectTreeProject(\'' + escapeHtml(item.projectName) + '\')">';
      html += '    <span>• [' + escapeHtml(item.nodeAlias || 'local') + '] <strong>' + escapeHtml(item.projectName) + '</strong> <span style="color:var(--dim);font-weight:normal;">(' + escapeHtml(item.workspacePath) + ')</span></span>';
      html += '    <button style="padding:2px 8px;font-size:11px;" onclick="event.stopPropagation(); selectProject(\'' + escapeHtml(item.projectName) + '\'); closePromptTree(); document.getElementById(\'promptText\').focus();">Target</button>';
      html += '  </div>';

      const actives = item.activePrompts || [];
      const queueds = item.queuedPrompts || [];
      const total = actives.length + queueds.length;

      if (total === 0) {
        html += '  <div class="tree-node-line" style="color:var(--dim); padding-left:18px;">└── (no active or queued prompts)</div>';
      } else {
        actives.forEach(function(ap, i) {
          const isLast = (i === actives.length - 1 && queueds.length === 0);
          const branch = isLast ? '└──' : '├──';
          html += '  <div class="tree-node-line">';
          html += '    <span style="color:var(--dim);">' + branch + '</span>';
          html += '    <span class="status-tag running">RUNNING</span>';
          html += '    <span class="snippet" title="' + escapeHtml(ap.promptSnippet || '') + '">' + escapeHtml(ap.promptSnippet || 'Active Task') + '</span>';
          html += '  </div>';
        });
        queueds.forEach(function(qp, j) {
          const isLast = (j === queueds.length - 1);
          const branch = isLast ? '└──' : '├──';
          html += '  <div class="tree-node-line">';
          html += '    <span style="color:var(--dim);">' + branch + '</span>';
          html += '    <span class="status-tag queued">QUEUED</span>';
          html += '    <span>' + escapeHtml(qp.title || qp.prompt || 'Queued Instruction') + ' (#' + qp.id + ')</span>';
          html += '  </div>';
        });
      }
      html += '</div>';
    });
    body.innerHTML = html;
  }
</script>
`

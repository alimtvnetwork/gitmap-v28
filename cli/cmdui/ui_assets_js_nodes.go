package cmdui

const uiAssetsJSNodes = `    async function refreshNodes() {
      try {
        const res = await fetch('/api/ssh/nodes');
        const data = await res.json();
        const tbody = document.getElementById('nodes-body');
        tbody.innerHTML = '';
        if (!data || data.length === 0) {
          tbody.innerHTML = '<tr><td colspan="8" style="text-align:center">No nodes registered. Use "gitmap ssh join" to add nodes.</td></tr>';
          return;
        }
        data.forEach(n => {
          const row = document.createElement('tr');
          const statusLabel = n.status || (n.isOnline ? 'Online' : 'Offline');
          const statusClass = n.isOnline ? 'badge-online' : 'badge-offline';
          row.innerHTML = ` + "`" + `<td>${n.nodeId || '-'}</td><td><b>${n.alias}</b></td><td>${n.host}</td><td>${n.os || 'Linux'}</td><td>${n.account || '-'}</td><td>${n.instances || 0}</td><td><span class="badge ${statusClass}">${statusLabel}</span></td><td><button class="btn" style="padding:0.25rem 0.5rem;font-size:0.75rem" onclick="testNode('${n.alias}')">Ping</button></td>` + "`" + `;
          tbody.appendChild(row);
        });
      } catch (e) { console.error(e); }
    }
    async function testNode(alias) {
      showToast('Pinging node ' + alias + '...', 'info');
      try {
        const res = await fetch('/api/ssh/test?alias=' + encodeURIComponent(alias));
        const data = await res.json();
        if (data.success || data.isOnline) {
          showToast('Node ' + alias + ' is online (' + (data.latencyMs || 0) + 'ms)', 'success');
        } else {
          showToast('Node ' + alias + ' ping failed: ' + (data.error || 'unreachable'), 'error');
        }
        refreshNodes();
      } catch (e) {
        showToast('Ping error: ' + e.message, 'error');
      }
    }
    async function populateEditorNodes() {
      try {
        const data = await (await fetch('/api/ssh/nodes')).json();
        const sel = document.getElementById('editor-node');
        if (sel) {
          sel.innerHTML = '<option value="local">local (current machine)</option>';
          (data || []).forEach(function (n) {
            const opt = document.createElement('option');
            opt.value = n.alias; opt.innerText = n.alias + ' (' + n.host + ')';
            sel.appendChild(opt);
          });
        }
        const termSel = document.getElementById('term-node');
        if (termSel && data && data.length) {
          termSel.innerHTML = '<option value="local">local</option>';
          data.forEach(function (n) {
            const opt = document.createElement('option');
            opt.value = n.alias; opt.innerText = n.alias;
            termSel.appendChild(opt);
          });
        }
      } catch (e) {}
    }
    async function openRemoteFile() {
      const node = document.getElementById('editor-node').value, path = document.getElementById('editor-path').value.trim();
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
      const node = document.getElementById('editor-node').value, path = document.getElementById('editor-path').value.trim(), content = document.getElementById('editor-area').value;
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

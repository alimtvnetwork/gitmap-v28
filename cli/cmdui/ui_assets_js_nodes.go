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
          if (!data.some(n => n.alias === 'u1')) {
            const u1Opt = document.createElement('option');
            u1Opt.value = 'u1';
            u1Opt.innerText = 'u1';
            termSel.appendChild(u1Opt);
          }
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

`

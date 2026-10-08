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
          if (data.theme && document.getElementById('setting-theme')) document.getElementById('setting-theme').value = data.theme;
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
          alert('Settings saved successfully!');
        } else {
          alert('Failed to save settings: ' + (data.error || 'Unknown error'));
        }
      } catch (e) { alert('Failed to save settings: ' + e.message); }
    }

    let termHistory = [];
    let termHistoryIdx = -1;

`

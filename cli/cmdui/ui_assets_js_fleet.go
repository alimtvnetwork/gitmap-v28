package cmdui

const uiAssetsJSCore = `    function showTab(tabId) {
      document.querySelectorAll('.content-area').forEach(el => el.classList.remove('active'));
      document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));
      const target = document.getElementById('tab-' + tabId);
      if (target) target.classList.add('active');
      const prettyTitle = tabId.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase());
      document.getElementById('header-title').innerText = prettyTitle;
      document.title = 'GitMap — ' + prettyTitle;
      event?.target?.classList.add('active');
      if (tabId === 'ssh') refreshNodes();
      if (tabId === 'editor') populateEditorNodes();
    }

    // Fills every [data-icon] placeholder with the inline SVG from the icon()
    // helper (provided by the theme assets). Guarded so the page still boots
    // if the helper is unavailable; injected SVGs are normalized to 16px.
    function hydrateIcons() {
      if (typeof icon !== 'function') return;
      document.querySelectorAll('[data-icon]').forEach(function (el) {
        if (el.dataset.iconDone) return;
        el.dataset.iconDone = '1';
        el.innerHTML = icon(el.getAttribute('data-icon'));
        const svg = el.querySelector('svg');
        if (svg) {
          svg.setAttribute('width', '16');
          svg.setAttribute('height', '16');
          svg.style.verticalAlign = '-3px';
          svg.style.marginRight = '6px';
        }
      });
    }

`

const uiAssetsJSCommit = `    async function execCommitin() {
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

`

const uiAssetsJSFleet = `    async function deployKeysFleetUI() {
      const out = document.getElementById('ssh-keys-output') || document.getElementById('ssh-nodes-export-area');
      if (out) out.value = 'Deploying cluster SSH keys to all nodes...';
      try {
        const res = await fetch('/api/ssh/deploy-keys', { method: 'POST' });
        const data = await res.json();
        if (out) out.value = data.output || data.error || (data.success ? 'Keys deployed successfully!' : 'Deploy failed');
        showToast(data.success ? 'Public keys successfully deployed to fleet!' : 'Deployment finished with warnings', data.success ? 'success' : 'error');
      } catch (e) {
        if (out) out.value = 'Deploy error: ' + e.message;
        showToast('Error deploying keys: ' + e.message, 'error');
      }
    }

    async function exportSSHNodesUI() {
      const area = document.getElementById('ssh-nodes-export-area');
      try {
        const res = await fetch('/api/ssh/export');
        const data = await res.json();
        const jsonStr = JSON.stringify(data, null, 2);
        if (area) area.value = jsonStr;
        const blob = new Blob([jsonStr], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'gitmap-ssh-nodes.json';
        a.click();
        URL.revokeObjectURL(url);
      } catch (e) {
        showToast('Failed to export SSH nodes: ' + e.message, 'error');
      }
    }

    async function importSSHNodesUI() {
      const area = document.getElementById('ssh-nodes-export-area');
      if (!area || !area.value.trim()) {
        showToast('Please paste the fleet JSON export into the text area or use terminal: gitmap import-ssh <file.json>', 'info');
        return;
      }
      try {
        const payload = JSON.parse(area.value.trim());
        const res = await fetch('/api/ssh/import', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        const data = await res.json();
        if (data.success) {
          showToast('Successfully imported ' + (data.imported || 0) + ' fleet node(s)!', 'success');
          if (typeof loadSSHConnections === 'function') loadSSHConnections();
        } else {
          showToast('Import failed: ' + (data.error || 'Unknown error'), 'error');
        }
      } catch (e) {
        showToast('Invalid JSON or import error: ' + e.message, 'error');
      }
    }

    async function viewSSHPublicKeyUI() {
      const out = document.getElementById('ssh-quick-card-output');
      const txt = document.getElementById('ssh-quick-display');
      if (out) out.style.display = 'block';
      if (txt) txt.value = 'Retrieving public key...';
      try {
        const res = await fetch('/api/ssh/export');
        const data = await res.json();
        const payload = data.payload || data;
        let key = payload.publicKey || payload.PublicKey || '';
        if (!key && payload.connections && payload.connections.length > 0) {
          key = payload.connections[0].publicKey || payload.connections[0].PublicKey || '';
        }
        if (!key) {
          key = 'Local Public Key (run locally: gitmap ssh key show):\nUse "gitmap ssh copy-id <alias>" to deploy keys to remote nodes.';
        }
        if (txt) txt.value = key;
      } catch (e) {
        if (txt) txt.value = 'Run locally to view key: gitmap ssh key show\nError: ' + e.message;
      }
    }

    function managePortsAndFirewallUI() {
      const out = document.getElementById('ssh-quick-card-output');
      const txt = document.getElementById('ssh-quick-display');
      if (out) out.style.display = 'block';
      if (txt) txt.value = 'SSH Ports & Firewall Management Commands:\n' +
        '  • List open ports & firewall status:  gitmap ssh port ls\n' +
        '  • Open & enable port on target:       gitmap ssh enable --port 22\n' +
        '  • Configure listening SSH port:       gitmap ssh port set <port>\n' +
        '  • Diagnose connectivity & firewall:   gitmap ssh troubleshoot <ip>';
    }

    function openImportNodesModal() {
      showTab('import-export');
    }

    async function loadAgyInstances() {
      try {
        const res = await fetch('/api/instances');
        const data = await res.json();
        if (data && data.isSuccess && Array.isArray(data.instances)) {
          const sel = document.getElementById('setting-agy-instance');
          const list = document.getElementById('instances-list-container');
          if (sel) {
            const curVal = sel.value;
            sel.innerHTML = '';
            data.instances.forEach(inst => {
              const opt = document.createElement('option');
              opt.value = inst.instanceId;
              opt.textContent = inst.instanceId + ' — ' + inst.instanceName + (inst.isRunning ? ' (Running)' : ' (Idle)');
              sel.appendChild(opt);
            });
            if (curVal) sel.value = curVal;
          }
          if (list) {
            list.innerHTML = '';
            data.instances.forEach(inst => {
              const card = document.createElement('div');
              card.style.background = 'var(--code-bg)';
              card.style.border = '1px solid var(--border)';
              card.style.borderRadius = '6px';
              card.style.padding = '0.85rem';
              const badgeClass = inst.isRunning ? 'badge-online' : 'badge-offline';
              const statusText = inst.isRunning ? 'RUNNING' : 'IDLE';
              const bridgeInfo = inst.languageServer ? ' | Bridge: ' + inst.languageServer : '';
              card.innerHTML =
                '<div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:6px;">' +
                  '<strong style="color:var(--primary); font-size:0.9rem;">' + inst.instanceName + ' [' + inst.instanceId + ']</strong>' +
                  '<span class="badge ' + badgeClass + '">' + statusText + '</span>' +
                '</div>' +
                '<div style="font-size:0.8rem; color:var(--muted); line-height:1.5; font-family:monospace;">' +
                  'PID: ' + (inst.processId > 0 ? inst.processId : 'N/A') + bridgeInfo + '<br>' +
                  'Path: ' + (inst.configDir || 'N/A') +
                '</div>';
              list.appendChild(card);
            });
          }
        }
      } catch (e) {
        console.error('Failed to load AGY instances', e);
      }
    }

    // Auto-detect route on load
    window.addEventListener('load', () => {
      hydrateIcons();
      loadSettings();
      loadAgyInstances();
      initTerminalListeners();
      populateEditorNodes();
      const path = window.location.pathname.replace(/^\//, '');
      if (path && document.getElementById('tab-' + path)) {
        showTab(path);
      }
    });
  </script>
`

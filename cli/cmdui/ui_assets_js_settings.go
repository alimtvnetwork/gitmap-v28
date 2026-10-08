package cmdui

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

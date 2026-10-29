package cmdui

const uiAssetsJSHistory = `    window.runHistoryPurge = async function () {
      const body = {
        repo: document.getElementById('history-repo').value,
        file: document.getElementById('history-file').value,
        folder: document.getElementById('history-folder').value,
        commit: document.getElementById('history-commit').value,
        range: document.getElementById('history-range').value,
        isDryRun: document.getElementById('history-dry-run').checked,
        confirmText: document.getElementById('history-confirm').value
      };
      const out = document.getElementById('history-purge-result');
      out.textContent = 'Running...';
      try {
        const res = await fetch('/api/history/purge', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(body)
        });
        const data = await res.json();
        out.textContent = JSON.stringify(data, null, 2);
      } catch (err) {
        out.textContent = err.message;
      }
    };
`

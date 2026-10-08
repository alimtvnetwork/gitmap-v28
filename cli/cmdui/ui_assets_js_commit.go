package cmdui

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

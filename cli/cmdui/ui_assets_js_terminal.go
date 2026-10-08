package cmdui

const uiAssetsJSTerminal = `    function clearTerminalOutput() {
      const out = document.getElementById('term-output');
      if (out) out.innerHTML = '';
      const badge = document.getElementById('term-badge');
      if (badge) badge.style.display = 'none';
    }

    async function runTerminalCommand(cmdOverride) {
      const inputEl = document.getElementById('term-input');
      const cmd = cmdOverride ? cmdOverride.trim() : (inputEl ? inputEl.value.trim() : '');
      if (!cmd) return;

      if (!cmdOverride && inputEl) {
        termHistory.push(cmd);
        termHistoryIdx = termHistory.length;
        inputEl.value = '';
      }

      const node = document.getElementById('term-node') ? document.getElementById('term-node').value : 'local';
      const cwd = document.getElementById('term-cwd') ? document.getElementById('term-cwd').value.trim() : '';
      const out = document.getElementById('term-output');
      const badge = document.getElementById('term-badge');

      if (out) {
        const promptLine = document.createElement('div');
        promptLine.className = 'term-prompt';
        promptLine.textContent = '[' + node + (cwd ? ' ' + cwd : '') + ']$ ' + cmd;
        out.appendChild(promptLine);

        const runLine = document.createElement('div');
        runLine.className = 'term-info';
        runLine.textContent = '⏳ Executing...';
        out.appendChild(runLine);
        out.scrollTop = out.scrollHeight;

        try {
          const res = await fetch('/api/terminal/exec', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ command: cmd, nodeAlias: node, cwd: cwd, timeoutSec: 30 })
          });
          const data = await res.json();
          runLine.remove();

          if (data.stdout) {
            const stdoutEl = document.createElement('div');
            stdoutEl.className = 'term-stdout';
            stdoutEl.textContent = data.stdout;
            out.appendChild(stdoutEl);
          }

          if (data.stderr) {
            const stderrEl = document.createElement('div');
            stderrEl.className = 'term-stderr';
            stderrEl.textContent = data.stderr;
            out.appendChild(stderrEl);
          }

          if (data.error && !data.stderr) {
            const errEl = document.createElement('div');
            errEl.className = 'term-stderr';
            errEl.textContent = 'Error: ' + data.error;
            out.appendChild(errEl);
          }

          const metaEl = document.createElement('div');
          metaEl.style.fontSize = '0.75rem';
          metaEl.style.color = 'var(--muted)';
          metaEl.style.marginTop = '4px';
          metaEl.style.marginBottom = '8px';
          metaEl.textContent = 'Done in ' + (data.durationMs || 0) + 'ms (exit code ' + (data.exitCode !== undefined ? data.exitCode : (data.success ? 0 : 1)) + ')';
          out.appendChild(metaEl);

          if (badge) {
            badge.style.display = 'inline-block';
            if (data.success) {
              badge.className = 'badge badge-online';
              badge.textContent = 'Exit 0 (' + data.durationMs + 'ms)';
            } else {
              badge.className = 'badge badge-offline';
              badge.textContent = 'Exit ' + data.exitCode + ' (' + data.durationMs + 'ms)';
            }
          }
        } catch (e) {
          runLine.remove();
          const errEl = document.createElement('div');
          errEl.className = 'term-stderr';
          errEl.textContent = 'Network error: ' + e.message;
          out.appendChild(errEl);
        }
        out.scrollTop = out.scrollHeight;
      }
    }

    function initTerminalListeners() {
      const termInput = document.getElementById('term-input');
      if (termInput) {
        termInput.addEventListener('keydown', (e) => {
          if (e.key === 'Enter') {
            e.preventDefault();
            runTerminalCommand();
          } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            if (termHistory.length > 0 && termHistoryIdx > 0) {
              termHistoryIdx--;
              termInput.value = termHistory[termHistoryIdx];
            }
          } else if (e.key === 'ArrowDown') {
            e.preventDefault();
            if (termHistoryIdx < termHistory.length - 1) {
              termHistoryIdx++;
              termInput.value = termHistory[termHistoryIdx];
            } else {
              termHistoryIdx = termHistory.length;
              termInput.value = '';
            }
          }
        });
      }
    }

`

package cmdui

const uiAssetsJSSearch = `    // ---- Repository Search Functions ----
    let repoSearchDebounceTimer = null;
    let selectedSearchExt = '';

    function handleSearchInputKey(ev) {
      if (ev.key === 'Enter') {
        ev.preventDefault();
        triggerRepoSearch();
        return;
      }
      clearTimeout(repoSearchDebounceTimer);
      repoSearchDebounceTimer = setTimeout(triggerRepoSearch, 400);
    }

    function toggleSearchExtPill(btn) {
      document.querySelectorAll('#search-ext-pills .ext-pill').forEach(function (p) {
        p.classList.remove('active');
      });
      btn.classList.add('active');
      selectedSearchExt = btn.getAttribute('data-ext') || '';
      triggerRepoSearch();
    }

    async function triggerRepoSearch() {
      const qInput = document.getElementById('repo-search-input');
      const q = qInput ? qInput.value.trim() : '';
      if (!q) {
        const resCard = document.getElementById('search-results-card');
        if (resCard) resCard.style.display = 'none';
        return;
      }
      await performRepoSearch(q);
    }

    async function performRepoSearch(query) {
      const dirInput = document.getElementById('repo-search-dir');
      const dir = dirInput ? dirInput.value.trim() : '';
      const isRegex = document.getElementById('search-flag-regex') ? document.getElementById('search-flag-regex').checked : false;
      const isCase = document.getElementById('search-flag-case') ? document.getElementById('search-flag-case').checked : true;
      const badge = document.getElementById('search-status-badge');
      if (badge) badge.innerText = 'Searching...';
      const params = new URLSearchParams({ q: query, dir: dir, ext: selectedSearchExt, isRegex: isRegex, isCaseInsensitive: isCase });
      try {
        const res = await fetch('/api/search/repo?' + params.toString());
        const data = await res.json();
        renderRepoSearchResults(data);
        if (badge) badge.innerText = 'AUM Multi-Core';
      } catch (err) {
        showToast('Search error: ' + err.message, 'error');
        if (badge) badge.innerText = 'Error';
      }
    }

    function renderRepoSearchResults(data) {
      const card = document.getElementById('search-results-card');
      const container = document.getElementById('search-results-container');
      const summary = document.getElementById('search-results-summary');
      const dur = document.getElementById('search-duration-info');
      if (!card || !container) return;
      card.style.display = 'block';
      container.innerHTML = '';
      if (data.error) {
        summary.innerText = 'Search failed: ' + data.error;
        dur.innerText = '';
        return;
      }
      const matches = data.matches || [];
      const ms = Math.round((data.duration || 0) / 1000000);
      summary.innerText = matches.length + ' matches across ' + (data.totalFiles || 0) + ' files';
      dur.innerText = ms + 'ms';
      if (matches.length === 0) {
        container.innerHTML = '<div style="padding:var(--sp-3); color:var(--muted); font-size:var(--fs-sm);">No matches found.</div>';
        return;
      }
      matches.forEach(function (m) {
        container.appendChild(createSearchResultCard(m));
      });
    }

    function createSearchResultCard(match) {
      const row = document.createElement('div');
      row.className = 'search-match-card';
      row.style.background = 'var(--code-bg)';
      row.style.border = '1px solid var(--border)';
      row.style.borderRadius = 'var(--radius-sm)';
      row.style.padding = 'var(--sp-2) var(--sp-3)';
      row.style.position = 'relative';
      row.style.cursor = 'pointer';
      row.appendChild(buildSearchCardHeader(match));
      row.appendChild(buildSearchCardSnippet(match));
      attachSearchCardHoverTooltip(row, match);
      return row;
    }

    function buildSearchCardHeader(match) {
      const hdr = document.createElement('div');
      hdr.style.display = 'flex';
      hdr.style.justifyContent = 'space-between';
      hdr.style.alignItems = 'center';
      hdr.style.marginBottom = '4px';
      hdr.innerHTML = '<div style="display:flex; align-items:center; gap:var(--sp-2); overflow:hidden; text-overflow:ellipsis;">' +
        '<span style="font-weight:600; font-size:var(--fs-xs); color:var(--primary);">' + escapeHtmlSearch(match.path) + '</span>' +
        '<span class="badge badge-online" style="font-size:11px; padding:1px 6px;">Line ' + match.lineNumber + '</span>' +
      '</div>' +
      '<div style="display:flex; gap:6px;">' +
        '<button class="btn btn-xs btn-open-editor" title="Open in Remote Editor" style="font-size:11px; padding:2px 6px;">Editor</button>' +
        '<button class="btn btn-xs btn-open-agy" title="Open in Antigravity" style="font-size:11px; padding:2px 6px;">Antigravity</button>' +
      '</div>';
      hdr.querySelector('.btn-open-editor').addEventListener('click', function (e) {
        e.stopPropagation();
        openInEditor(match.path, match.lineNumber);
      });
      hdr.querySelector('.btn-open-agy').addEventListener('click', function (e) {
        e.stopPropagation();
        openInAntigravity(match.path);
      });
      return hdr;
    }

    function buildSearchCardSnippet(match) {
      const snip = document.createElement('div');
      snip.style.fontFamily = 'monospace';
      snip.style.fontSize = 'var(--fs-xs)';
      snip.style.color = 'var(--text)';
      snip.style.whiteSpace = 'pre';
      snip.style.overflow = 'hidden';
      snip.style.textOverflow = 'ellipsis';
      snip.textContent = match.lineContent || '';
      return snip;
    }

    function attachSearchCardHoverTooltip(row, match) {
      const tt = document.createElement('div');
      tt.className = 'search-tooltip';
      tt.style.display = 'none';
      tt.style.position = 'absolute';
      tt.style.bottom = '105%';
      tt.style.left = '16px';
      tt.style.zIndex = '100';
      tt.style.background = 'var(--raised)';
      tt.style.border = '1px solid var(--border)';
      tt.style.borderRadius = 'var(--radius-sm)';
      tt.style.padding = 'var(--sp-2)';
      tt.style.boxShadow = '0 4px 12px rgba(0,0,0,0.3)';
      tt.style.fontFamily = 'monospace';
      tt.style.fontSize = '11px';
      tt.style.maxWidth = '600px';
      tt.style.pointerEvents = 'none';
      row.appendChild(tt);
      row.addEventListener('mouseenter', async function () {
        await populateSearchTooltipContext(tt, match);
        tt.style.display = 'block';
      });
      row.addEventListener('mouseleave', function () {
        tt.style.display = 'none';
      });
      row.addEventListener('click', function () {
        openInEditor(match.path, match.lineNumber);
      });
    }

    async function populateSearchTooltipContext(tt, match) {
      if (tt.dataset.loaded) return;
      tt.dataset.loaded = '1';
      tt.innerText = 'Loading contextual preview...';
      try {
        const res = await fetch('/api/editor/read', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({nodeAlias: 'local', filePath: match.path})
        });
        const data = await res.json();
        if (data.isSuccess && data.content) {
          renderContextSnippetInTooltip(tt, data.content, match.lineNumber);
          return;
        }
        tt.innerText = match.lineContent || 'Preview unavailable';
      } catch (err) {
        tt.innerText = match.lineContent || 'Preview unavailable';
      }
    }

    function renderContextSnippetInTooltip(tt, fullContent, targetLine) {
      const lines = fullContent.split('\n');
      const start = Math.max(0, targetLine - 2);
      const end = Math.min(lines.length, targetLine + 2);
      let html = '';
      for (let i = start; i < end; i++) {
        const num = i + 1;
        const isMatch = num === targetLine;
        const lineText = escapeHtmlSearch(lines[i]);
        const bg = isMatch ? 'background:rgba(255,255,0,0.15); font-weight:600;' : '';
        html += '<div style="display:flex; gap:8px;' + bg + '">' +
          '<span style="color:var(--muted); min-width:30px; text-align:right;">' + num + '</span>' +
          '<span>' + lineText + '</span>' +
        '</div>';
      }
      tt.innerHTML = html;
    }

    function openInEditor(filePath, lineNumber) {
      const pathInput = document.getElementById('editor-path');
      if (pathInput) pathInput.value = filePath;
      showTab('editor');
      if (typeof openRemoteFile === 'function') {
        openRemoteFile().then(function () {
          scrollEditorToLine(lineNumber);
        });
        return;
      }
      scrollEditorToLine(lineNumber);
    }

    function scrollEditorToLine(lineNumber) {
      const area = document.getElementById('editor-area');
      if (!area || !lineNumber) return;
      const lines = area.value.split('\n');
      let offset = 0;
      for (let i = 0; i < lineNumber - 1 && i < lines.length; i++) {
        offset += lines[i].length + 1;
      }
      area.focus();
      area.setSelectionRange(offset, offset);
    }

    async function openInAntigravity(filePath) {
      try {
        const res = await fetch('/api/project/open-antigravity?path=' + encodeURIComponent(filePath));
        const data = await res.json();
        if (data.isSuccess) {
          showToast('Opened in Antigravity: ' + filePath, 'success');
          return;
        }
        showToast('Failed to open in Antigravity: ' + (data.error || 'unknown'), 'error');
      } catch (err) {
        showToast('Error opening Antigravity: ' + err.message, 'error');
      }
    }

    function escapeHtmlSearch(str) {
      if (!str) return '';
      return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }
`

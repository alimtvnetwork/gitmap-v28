package cmdlogviewer

// showLogMarkupScript contains client-side application logic.
const showLogMarkupScript = `<script>
let allCommits = [];
let activeHash = '';
let activePath = '';
let isCompareMode = false;
let compareBaseHash = '';
let compareTargetHash = '';

const params = new URLSearchParams(window.location.search);
activePath = params.get('path') || '.';
document.getElementById('path-indicator').textContent = activePath === '.' ? '[All repository files]' : activePath;
document.getElementById('modal-compare-path').value = activePath === '.' ? '' : activePath;

function showToast(msg) {
  const el = document.getElementById('toast');
  el.textContent = msg;
  setTimeout(() => { if (el.textContent === msg) el.textContent = ''; }, 5000);
}

async function loadBranches() {
  try {
    const res = await fetch('/api/log/branches');
    const data = await res.json();
    const sel = document.getElementById('branch-select');
    const searchInput = document.getElementById('branch-search');
    const datalist = document.getElementById('branch-list');
    const modalA = document.getElementById('modal-branch-a');
    const modalB = document.getElementById('modal-branch-b');
    sel.innerHTML = ''; modalA.innerHTML = ''; modalB.innerHTML = '';
    if (datalist) datalist.innerHTML = '';
    (data.branches || []).forEach(b => {
      const opt = new Option(b, b);
      if (b === data.activeBranch) opt.selected = true;
      sel.add(opt);
      if (datalist) {
        const dOpt = document.createElement('option');
        dOpt.value = b;
        datalist.appendChild(dOpt);
      }
      modalA.add(new Option(b, b));
      modalB.add(new Option(b, b));
    });
    if (searchInput) searchInput.value = data.activeBranch || '';
    if (modalB.options.length > 1) modalB.selectedIndex = 1;
  } catch (e) {
    console.error('Failed to load branches', e);
  }
}

const branchSearch = document.getElementById('branch-search');
if (branchSearch) {
  branchSearch.addEventListener('change', () => {
    const val = branchSearch.value.trim();
    const sel = document.getElementById('branch-select');
    if (sel && val) {
      for (let i = 0; i < sel.options.length; i++) {
        if (sel.options[i].value === val) {
          sel.selectedIndex = i;
          break;
        }
      }
    }
  });
}

const branchSel = document.getElementById('branch-select');
if (branchSel) {
  branchSel.addEventListener('change', () => {
    const val = branchSel.value;
    const searchInput = document.getElementById('branch-search');
    if (searchInput) searchInput.value = val;
  });
}

async function loadCommits() {
  try {
    const res = await fetch('/api/log/commits?path=' + encodeURIComponent(activePath) + '&limit=150');
    allCommits = await res.json();
    renderCommits(allCommits);
    if (allCommits.length > 0) {
      selectCommit(allCommits[0].hash);
    }
  } catch (e) {
    document.getElementById('commits-body').innerHTML = '<tr><td colspan="5" style="color:var(--red);text-align:center;">Failed to load commits</td></tr>';
  }
}

function renderCommits(commits) {
  const tbody = document.getElementById('commits-body');
  if (!commits || commits.length === 0) {
    tbody.innerHTML = '<tr><td colspan="5" style="text-align:center;color:var(--text-muted);padding:30px;">No commits found for path</td></tr>';
    return;
  }
  tbody.innerHTML = commits.map(c => {
    const selCls = c.hash === activeHash ? 'selected' : '';
    return '<tr data-hash="' + c.hash + '" class="' + selCls + '" onclick="handleCommitClick(event, \'' + c.hash + '\')" oncontextmenu="showCommitContextMenu(event, \'' + c.hash + '\', \'' + escapeQuotes(c.subject) + '\')">' +
      '<td style="color:var(--accent);">●</td>' +
      '<td class="mono hash-cell">' + c.shortHash + '</td>' +
      '<td class="date-cell">' + c.date + '</td>' +
      '<td class="author-cell" title="' + escapeQuotes(c.author) + '">' + c.author + '</td>' +
      '<td title="' + escapeQuotes(c.subject) + '">' + c.subject + '</td>' +
      '</tr>';
  }).join('');
}

function handleCommitClick(event, hash) {
  if (event && (event.ctrlKey || event.metaKey)) {
    if (activeHash && activeHash !== hash) {
      compareTwoCommits(activeHash, hash);
      return;
    }
  }
  isCompareMode = false;
  compareBaseHash = '';
  compareTargetHash = '';
  selectCommit(hash);
}

async function selectCommit(hash) {
  activeHash = hash;
  isCompareMode = false;
  document.querySelectorAll('#commits-body tr').forEach(tr => {
    tr.classList.toggle('selected', tr.getAttribute('data-hash') === hash);
  });
  try {
    const res = await fetch('/api/log/commit?hash=' + hash);
    const detail = await res.json();
    renderCommitDetail(detail);
    renderModifiedFiles(detail.files || []);
  } catch (e) {
    showToast('Failed to load commit details');
  }
}

function renderCommitDetail(detail) {
  const header = document.getElementById('detail-header');
  const body = document.getElementById('detail-body');
  const parents = (detail.parents || []).join(' ') || 'none (root)';
  header.innerHTML = '<div><b>Commit:</b> <span class="mono">' + detail.hash + '</span></div>' +
    '<div><b>Author:</b> ' + detail.author + ' &lt;' + (detail.authorEmail || '') + '&gt;</div>' +
    '<div><b>Date:</b> ' + detail.date + '</div>' +
    '<div><b>Parents:</b> <span class="mono">' + parents + '</span></div>';
  body.textContent = (detail.subject || '') + '\n\n' + (detail.body || '');
}

function renderModifiedFiles(files) {
  const tbody = document.getElementById('files-body');
  document.getElementById('files-summary').textContent = 'Modified Files (' + files.length + ')';
  if (!files || files.length === 0) {
    tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;color:var(--text-muted);padding:20px;">No files changed</td></tr>';
    return;
  }
  tbody.innerHTML = files.map(f => {
    const badgeCls = 'badge-' + (f.action || 'modified').toLowerCase();
    const stats = f.stats || '';
    return '<tr ondblclick="triggerDiff(\'' + f.path + '\')" oncontextmenu="showFileContextMenu(event, \'' + f.path + '\')" title="Double click or right click to diff">' +
      '<td><span class="badge-action ' + badgeCls + '">' + f.action + '</span></td>' +
      '<td class="mono" style="color:var(--text);">' + f.path + '</td>' +
      '<td class="mono" style="color:var(--yellow);">' + stats + '</td>' +
      '<td style="text-align:right;"><button onclick="triggerDiff(\'' + f.path + '\')">Diff</button></td>' +
      '</tr>';
  }).join('');
}

async function compareTwoCommits(hashA, hashB) {
  isCompareMode = true;
  compareBaseHash = hashA;
  compareTargetHash = hashB;
  activeHash = hashB;
  document.querySelectorAll('#commits-body tr').forEach(tr => {
    const h = tr.getAttribute('data-hash');
    tr.classList.toggle('selected', h === hashA || h === hashB);
  });
  showToast('Comparing ' + hashA.substring(0, 8) + ' with ' + hashB.substring(0, 8) + '...');
  try {
    const res = await fetch('/api/log/compare-commits?a=' + encodeURIComponent(hashA) + '&b=' + encodeURIComponent(hashB));
    const data = await res.json();
    renderCompareDetail(hashA, hashB, data.files || []);
    renderCompareFiles(hashA, hashB, data.files || []);
  } catch (e) {
    showToast('Failed to compare commits: ' + e.message);
  }
}

function renderCompareDetail(hashA, hashB, files) {
  const header = document.getElementById('detail-header');
  const body = document.getElementById('detail-body');
  header.innerHTML = '<div><b>⚖ Two-Commit Compare:</b> <span class="mono">' + hashA.substring(0, 8) + '</span> &harr; <span class="mono">' + hashB.substring(0, 8) + '</span></div>' +
    '<div><b>Changed Files:</b> ' + files.length + '</div>';
  body.textContent = 'Comparing changes between commit ' + hashA + '\nand commit ' + hashB + '.\n\nDouble-click any file below to launch external comparison (Beyond Compare / VS Code).';
}

function renderCompareFiles(hashA, hashB, files) {
  const tbody = document.getElementById('files-body');
  document.getElementById('files-summary').textContent = 'Comparison Diff Files (' + files.length + ')';
  if (!files || files.length === 0) {
    tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;color:var(--text-muted);padding:20px;">No differences between these two commits</td></tr>';
    return;
  }
  tbody.innerHTML = files.map(f => {
    const badgeCls = 'badge-' + (f.action || 'modified').toLowerCase();
    const stats = f.stats || '';
    return '<tr ondblclick="triggerTwoCommitDiff(\'' + hashA + '\', \'' + hashB + '\', \'' + f.path + '\')" oncontextmenu="showCompareFileContextMenu(event, \'' + hashA + '\', \'' + hashB + '\', \'' + f.path + '\')" title="Double click or right click to diff">' +
      '<td><span class="badge-action ' + badgeCls + '">' + f.action + '</span></td>' +
      '<td class="mono" style="color:var(--text);">' + f.path + '</td>' +
      '<td class="mono" style="color:var(--yellow);">' + stats + '</td>' +
      '<td style="text-align:right;"><button onclick="triggerTwoCommitDiff(\'' + hashA + '\', \'' + hashB + '\', \'' + f.path + '\')">Diff</button></td>' +
      '</tr>';
  }).join('');
}

async function triggerTwoCommitDiff(hashA, hashB, filePath) {
  const preferredTool = localStorage.getItem('gitmap_diff_tool') || '';
  showToast('Launching external compare tool for ' + filePath + '...');
  try {
    const res = await fetch('/api/log/diff-commits', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ hashA: hashA, hashB: hashB, path: filePath, tool: preferredTool })
    });
    const data = await res.json();
    if (data.isSuccess) {
      showToast('🚀 ' + (data.tool || 'Compare tool') + ' launched for ' + filePath);
    } else {
      showToast('Error: ' + (data.error || 'Failed to launch diff'));
    }
  } catch (e) {
    showToast('Failed to trigger diff: ' + e.message);
  }
}

function showCompareFileContextMenu(e, hashA, hashB, filePath) {
  positionContextMenu(e);
  contextMenu.innerHTML = '<div class="context-menu-item" onclick="navigator.clipboard.writeText(\'' + filePath + '\'); showToast(\'Copied path: ' + filePath + '\');">📋 Copy File Path</div>' +
    '<div class="context-menu-separator"></div>' +
    '<div class="context-menu-item" onclick="triggerTwoCommitDiff(\'' + hashA + '\', \'' + hashB + '\', \'' + filePath + '\')">⚖ Diff Commits (Beyond Compare)</div>';
}

async function triggerDiff(filePath) {
  if (!activeHash) return;
  const preferredTool = localStorage.getItem('gitmap_diff_tool') || '';
  showToast('Launching external compare tool for ' + filePath + '...');
  try {
    const res = await fetch('/api/log/diff', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ commit: activeHash, path: filePath, tool: preferredTool })
    });
    const data = await res.json();
    if (data.isSuccess) {
      showToast('🚀 ' + (data.tool || 'Compare tool') + ' launched for ' + filePath);
    } else {
      showToast('Error: ' + (data.error || 'Failed to launch diff'));
    }
  } catch (e) {
    showToast('Failed to trigger diff: ' + e.message);
  }
}

document.getElementById('commit-search').addEventListener('input', (e) => {
  const query = e.target.value.toLowerCase().trim();
  if (!query) {
    renderCommits(allCommits);
    return;
  }
  const filtered = allCommits.filter(c =>
    c.hash.toLowerCase().includes(query) ||
    c.author.toLowerCase().includes(query) ||
    c.subject.toLowerCase().includes(query)
  );
  renderCommits(filtered);
});

document.getElementById('btn-refresh').addEventListener('click', () => {
  showToast('Reloading commits...');
  loadCommits();
});

const modal = document.getElementById('compare-modal');
document.getElementById('btn-compare-branches').addEventListener('click', () => modal.classList.add('open'));
document.getElementById('modal-cancel').addEventListener('click', () => modal.classList.remove('open'));
document.getElementById('modal-submit').addEventListener('click', async () => {
  const bA = document.getElementById('modal-branch-a').value;
  const bB = document.getElementById('modal-branch-b').value;
  const p = document.getElementById('modal-compare-path').value.trim();
  modal.classList.remove('open');
  showToast('Launching comparison between ' + bA + ' and ' + bB + '...');
  try {
    const res = await fetch('/api/log/branch-compare', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ branchA: bA, branchB: bB, path: p })
    });
    const data = await res.json();
    if (data.isSuccess) {
      showToast('🚀 ' + (data.message || 'Comparison launched!'));
    } else {
      showToast('Error: ' + (data.error || 'Failed to compare branches'));
    }
  } catch (e) {
    showToast('Comparison failed: ' + e.message);
  }
});

function escapeQuotes(str) {
  return (str || '').replace(/'/g, "\\'").replace(/"/g, '&quot;');
}

const contextMenu = document.getElementById('context-menu');
document.addEventListener('click', () => { contextMenu.style.display = 'none'; });

function positionContextMenu(e) {
  e.preventDefault();
  contextMenu.style.display = 'block';
  contextMenu.style.left = Math.min(e.pageX, window.innerWidth - 200) + 'px';
  contextMenu.style.top = Math.min(e.pageY, window.innerHeight - 150) + 'px';
}

function showCommitContextMenu(e, hash, subject) {
  positionContextMenu(e);
  let compareItem = '';
  if (activeHash && activeHash !== hash) {
    compareItem = '<div class="context-menu-item" onclick="compareTwoCommits(\'' + activeHash + '\', \'' + hash + '\')">⚖ Compare with Selected (' + activeHash.substring(0, 7) + ')</div>';
  }
  contextMenu.innerHTML = '<div class="context-menu-item" onclick="navigator.clipboard.writeText(\'' + hash + '\'); showToast(\'Copied hash: ' + hash.substring(0, 8) + '\');">📋 Copy Commit Hash</div>' +
    '<div class="context-menu-item" onclick="navigator.clipboard.writeText(\'' + subject + '\'); showToast(\'Copied subject line\');">📋 Copy Subject Line</div>' +
    compareItem +
    '<div class="context-menu-separator"></div>' +
    '<div class="context-menu-item" onclick="triggerDiff(activePath)">🔍 Diff with Parent (Beyond Compare)</div>';
}

function showFileContextMenu(e, filePath) {
  positionContextMenu(e);
  contextMenu.innerHTML = '<div class="context-menu-item" onclick="navigator.clipboard.writeText(\'' + filePath + '\'); showToast(\'Copied path: ' + filePath + '\');">📋 Copy File Path</div>' +
    '<div class="context-menu-separator"></div>' +
    '<div class="context-menu-item" onclick="triggerDiff(\'' + filePath + '\')">🔍 Diff with Previous (Beyond Compare)</div>';
}

const btnFs = document.getElementById('btn-fullscreen');
if (btnFs) {
  btnFs.addEventListener('click', () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(() => {});
    } else {
      if (document.exitFullscreen) document.exitFullscreen().catch(() => {});
    }
  });
}

const settingsModal = document.getElementById('settings-modal');
const btnSettings = document.getElementById('btn-settings');
const settingDiffTool = document.getElementById('setting-diff-tool');
const settingAutoExit = document.getElementById('setting-auto-exit');

if (btnSettings && settingsModal) {
  btnSettings.addEventListener('click', () => {
    settingDiffTool.value = localStorage.getItem('gitmap_diff_tool') || 'auto';
    settingAutoExit.value = localStorage.getItem('gitmap_auto_exit') || '1';
    settingsModal.classList.add('open');
  });

  document.getElementById('modal-settings-cancel').addEventListener('click', () => {
    settingsModal.classList.remove('open');
  });

  document.getElementById('modal-settings-save').addEventListener('click', () => {
    localStorage.setItem('gitmap_diff_tool', settingDiffTool.value);
    localStorage.setItem('gitmap_auto_exit', settingAutoExit.value);
    settingsModal.classList.remove('open');
    showToast('⚙ Settings saved (diff tool: ' + settingDiffTool.value + ')');
  });
}

const btnRefresh = document.getElementById('btn-refresh');
if (btnRefresh) {
  btnRefresh.addEventListener('click', () => {
    showToast('Refreshing commits...');
    loadCommits();
  });
}

document.getElementById('btn-exit').addEventListener('click', async () => {
  showToast('Stopping server...');
  try {
    await fetch('/api/log/exit', { method: 'POST' });
  } catch (e) {
    // server shutdown terminates connection immediately
  }
  window.close();
});

document.addEventListener('keydown', (e) => {
  if (e.key !== 'ArrowDown' && e.key !== 'ArrowUp') return;
  if (!allCommits || allCommits.length === 0) return;
  const currentIndex = allCommits.findIndex(c => c.hash === activeHash);
  if (currentIndex === -1) return;
  const nextIndex = e.key === 'ArrowDown' ? Math.min(currentIndex + 1, allCommits.length - 1) : Math.max(currentIndex - 1, 0);
  if (nextIndex !== currentIndex) {
    e.preventDefault();
    selectCommit(allCommits[nextIndex].hash);
    const targetTr = document.querySelector('#commits-body tr[data-hash="' + allCommits[nextIndex].hash + '"]');
    if (targetTr) targetTr.scrollIntoView({ block: 'nearest' });
  }
});

loadBranches();
loadCommits();
if (params.get('compare') === '1') {
  modal.classList.add('open');
}
</script>
</body>
</html>`

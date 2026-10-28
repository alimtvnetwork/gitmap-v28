package cmdlogviewer

// showLogMarkupHead contains HTML structure and embedded styles.
const showLogMarkupHead = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GitMap — Show Log</title>
<link rel="manifest" href="/manifest.json">
<link rel="icon" type="image/x-icon" href="/favicon.ico">
<link rel="icon" type="image/svg+xml" href="/favicon.svg">
<meta name="theme-color" content="#181825">
<style>
  :root {
    --bg: #181825;
    --pane-bg: #1e1e2e;
    --card-bg: #252538;
    --border: #363a4f;
    --text: #cdd6f4;
    --text-muted: #a6adc8;
    --accent: #89b4fa;
    --accent-hover: #b4befe;
    --green: #a6e3a1;
    --red: #f38ba8;
    --yellow: #f9e2af;
    --blue: #89dceb;
    --sel-bg: #313244;
    --font-mono: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
    --font-ui: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  }
  ::-webkit-scrollbar { width: 8px; height: 8px; }
  ::-webkit-scrollbar-track { background: var(--pane-bg); }
  ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 4px; }
  ::-webkit-scrollbar-thumb:hover { background: var(--accent); }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body {
    background: var(--bg);
    color: var(--text);
    font-family: var(--font-ui);
    font-size: 13px;
    height: 100vh;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  header.toolbar {
    background: var(--pane-bg);
    border-bottom: 1px solid var(--border);
    padding: 8px 14px;
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
  }
  .brand {
    font-weight: 700;
    font-size: 14px;
    color: var(--accent);
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .path-badge {
    background: var(--card-bg);
    border: 1px solid var(--border);
    padding: 3px 10px;
    border-radius: 4px;
    color: var(--yellow);
    font-family: var(--font-mono);
    font-size: 11px;
    max-width: 320px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .search-box { flex: 1; max-width: 300px; }
  input, select, button {
    background: var(--card-bg);
    border: 1px solid var(--border);
    color: var(--text);
    padding: 5px 10px;
    border-radius: 4px;
    font-size: 12px;
    outline: none;
    transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }
  input:focus, select:focus { border-color: var(--accent); box-shadow: 0 0 0 2px rgba(137, 180, 250, 0.25); }
  button {
    cursor: pointer;
    background: #313244;
  }
  button:hover { background: var(--accent); color: #11111b; transform: translateY(-1px); box-shadow: 0 3px 10px rgba(137, 180, 250, 0.25); }
  button:active { transform: translateY(0); filter: brightness(0.95); }
  .btn-icon {
    width: 32px;
    height: 32px;
    padding: 6px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    position: relative;
    border-radius: 4px;
    flex-shrink: 0;
  }
  .btn-icon svg {
    width: 16px;
    height: 16px;
    fill: currentColor;
    pointer-events: none;
  }
  [data-tooltip] { position: relative; }
  [data-tooltip]::after {
    content: attr(data-tooltip);
    position: absolute;
    bottom: -28px;
    left: 50%;
    transform: translateX(-50%);
    background: #11111b;
    color: #cdd6f4;
    font-size: 11px;
    padding: 3px 8px;
    border-radius: 4px;
    border: 1px solid var(--border);
    white-space: nowrap;
    pointer-events: none;
    opacity: 0;
    transition: opacity 0.15s ease;
    z-index: 1000;
    box-shadow: 0 4px 12px rgba(0,0,0,0.5);
  }
  [data-tooltip]:hover::after { opacity: 1; }
  .btn-primary { background: var(--accent); color: #11111b; font-weight: 600; }
  .btn-primary:hover { background: var(--accent-hover); box-shadow: 0 3px 10px rgba(180, 190, 254, 0.3); }
  .btn-danger { background: #451a1a; color: #f38ba8; border-color: #f38ba8; }
  .btn-danger:hover { background: #f38ba8; color: #11111b; box-shadow: 0 3px 10px rgba(243, 139, 168, 0.3); }
  main.workspace { flex: 1; display: flex; flex-direction: column; overflow: hidden; animation: appEntrance 0.3s cubic-bezier(0.16, 1, 0.3, 1); }
  .pane-upper {
    flex: 4; min-height: 140px; display: flex; flex-direction: column;
    background: var(--pane-bg); border-bottom: 1px solid var(--border); overflow: hidden;
    transition: background 0.2s ease;
  }
  .table-scroll { flex: 1; overflow-y: auto; }
  table.data-table { width: 100%; border-collapse: collapse; font-size: 12px; }
  table.data-table th {
    position: sticky; top: 0; background: #252538; color: var(--text-muted);
    text-align: left; padding: 6px 10px; border-bottom: 1px solid var(--border); font-weight: 600;
  }
  table.data-table td {
    padding: 5px 10px; border-bottom: 1px solid rgba(255, 255, 255, 0.04);
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  }
  table.data-table tr { cursor: pointer; user-select: none; transition: background-color 0.15s cubic-bezier(0.16, 1, 0.3, 1), transform 0.15s cubic-bezier(0.16, 1, 0.3, 1); }
  table.data-table tr:hover { background: rgba(137, 180, 250, 0.08); }
  table.data-table tr.selected { background: var(--sel-bg); color: #fff; box-shadow: inset 3px 0 0 var(--accent); }
  .mono { font-family: var(--font-mono); }
  .hash-cell { color: var(--yellow); width: 90px; }
  .date-cell { color: var(--text-muted); width: 100px; }
  .author-cell { color: var(--green); width: 140px; }
  .pane-middle {
    flex: 3; min-height: 110px; background: var(--card-bg);
    border-bottom: 1px solid var(--border); display: flex; flex-direction: column; overflow: hidden;
    transition: background 0.2s ease;
  }
  .commit-header {
    padding: 8px 14px; background: rgba(0, 0, 0, 0.2); display: flex;
    flex-wrap: wrap; gap: 16px; border-bottom: 1px solid var(--border); font-size: 11px; color: var(--text-muted);
  }
  .commit-header b { color: var(--text); }
  .commit-body {
    flex: 1; padding: 10px 14px; font-family: var(--font-mono); font-size: 12px;
    white-space: pre-wrap; overflow-y: auto; color: var(--text); line-height: 1.45;
  }
  .pane-lower {
    flex: 3; min-height: 110px; background: var(--pane-bg);
    display: flex; flex-direction: column; overflow: hidden;
    transition: background 0.2s ease;
  }
  @keyframes appEntrance { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: none; } }
  .lower-header {
    padding: 6px 14px; background: #252538; border-bottom: 1px solid var(--border);
    display: flex; justify-content: space-between; align-items: center; color: var(--text-muted); font-size: 11px;
  }
  .badge-action { display: inline-block; padding: 2px 6px; border-radius: 3px; font-size: 10px; font-weight: 700; text-transform: uppercase; }
  .badge-added { background: rgba(166, 227, 161, 0.2); color: var(--green); }
  .badge-modified { background: rgba(137, 220, 235, 0.2); color: var(--blue); }
  .badge-deleted { background: rgba(243, 139, 168, 0.2); color: var(--red); }
  .badge-renamed { background: rgba(249, 226, 175, 0.2); color: var(--yellow); }
  footer.statusbar {
    background: #11111b; border-top: 1px solid var(--border); padding: 4px 14px;
    display: flex; justify-content: space-between; align-items: center; font-size: 11px; color: var(--text-muted); flex-shrink: 0;
  }
  .toast-message { color: var(--green); font-weight: 500; }
  .modal-overlay {
    position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.65);
    display: none; align-items: center; justify-content: center; z-index: 1000;
  }
  .modal-overlay.open { display: flex; }
  .modal-card {
    background: var(--pane-bg); border: 1px solid var(--border); border-radius: 8px;
    width: 440px; padding: 20px; box-shadow: 0 8px 30px rgba(0, 0, 0, 0.5);
  }
  .modal-card h3 { margin-bottom: 14px; color: var(--accent); }
  .form-row { margin-bottom: 12px; }
  .form-row label { display: block; margin-bottom: 4px; color: var(--text-muted); font-size: 11px; }
  .form-row select, .form-row input { width: 100%; }
  .modal-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 18px; }
  .btn-danger { background: var(--red); color: #11111b; font-weight: 600; }
  .btn-danger:hover { background: #eba0ac; }
  .context-menu {
    background: var(--pane-bg);
    border: 1px solid var(--border);
    border-radius: 6px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
    min-width: 180px;
    padding: 4px 0;
    z-index: 2000;
  }
  .context-menu-item {
    padding: 6px 14px;
    cursor: pointer;
    font-size: 12px;
    color: var(--text);
    display: flex;
    align-items: center;
    gap: 8px;
    user-select: none;
  }
  .context-menu-item:hover {
    background: var(--sel-bg);
    color: var(--accent);
  }
  .context-menu-separator {
    height: 1px;
    background: var(--border);
    margin: 4px 0;
  }
  @media (max-width: 900px) {
    body { height: auto; min-height: 100vh; overflow-y: auto; }
    main.workspace { overflow-y: auto; }
    header.toolbar { flex-wrap: wrap; }
    .pane-upper, .pane-middle, .pane-lower { min-height: 200px; flex: none; }
  }
</style>
</head>
<body>
<header class="toolbar">
  <div class="brand">
    <svg viewBox="0 0 32 32" width="22" height="22" style="vertical-align:middle; border-radius:6px; flex-shrink:0;">
      <defs>
        <linearGradient id="gmGrad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="#89b4fa"/>
          <stop offset="100%" stop-color="#b4befe"/>
        </linearGradient>
        <linearGradient id="gmBg" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="#1e1e2e"/>
          <stop offset="100%" stop-color="#11111b"/>
        </linearGradient>
        <filter id="gmShadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="1" stdDeviation="1.5" flood-color="#000000" flood-opacity="0.4"/>
        </filter>
      </defs>
      <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#gmBg)" stroke="#89b4fa" stroke-width="1.6" filter="url(#gmShadow)"/>
      <rect x="2.5" y="2.5" width="27" height="27" rx="5.5" fill="none" stroke="#313244" stroke-width="0.8" opacity="0.6"/>
      <path d="M11 7.5v17M11 14c3.5 0 10 2.5 10 9.5" fill="none" stroke="url(#gmGrad)" stroke-width="2.2" stroke-linecap="round"/>
      <circle cx="11" cy="7.5" r="3" fill="#89b4fa" stroke="#cdd6f4" stroke-width="0.8"/>
      <circle cx="11" cy="7.5" r="1.2" fill="#1e1e2e"/>
      <circle cx="11" cy="24.5" r="3" fill="#89b4fa" stroke="#cdd6f4" stroke-width="0.8"/>
      <circle cx="11" cy="24.5" r="1.2" fill="#1e1e2e"/>
      <circle cx="21" cy="23.5" r="3.2" fill="#a6e3a1" stroke="#ffffff" stroke-width="0.8"/>
      <circle cx="21" cy="23.5" r="1.3" fill="#1e1e2e"/>
    </svg>
    <span>GitMap Show Log</span>
  </div>
  <div class="path-badge" id="path-indicator" title="Active Path">Loading path...</div>
  <div class="branch-combo" style="display:flex; align-items:center;">
    <input type="text" id="branch-search" list="branch-list" placeholder="Filter branch..." style="width:130px; font-family:var(--font-mono); font-size:11px;" title="Filter or switch branch">
    <datalist id="branch-list"></datalist>
    <select id="branch-select" title="Active Branch" style="max-width:140px; margin-left:4px;"></select>
  </div>
  <input type="text" id="commit-search" class="search-box" placeholder="Filter commits (hash, author, message)...">
  <button id="btn-refresh" class="btn-icon" data-tooltip="Refresh commit list">
    <svg viewBox="0 0 24 24"><path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/></svg>
  </button>
  <button id="btn-compare-branches" class="btn-icon" data-tooltip="Branch comparison">
    <svg viewBox="0 0 24 24"><path d="M19 13v-2c-1.54.02-3.09-.75-4.07-1.83l-1.29-1.43c-.87-.97-2.11-1.55-3.41-1.61L10 6V4c0-1.1-.9-2-2-2S6 2.9 6 4v3.59c-.58.33-1 .95-1 1.66 0 1.07.87 1.94 1.94 1.94h.06c.45 1.5 1.62 2.68 3.12 3.12v1.88c-.58.33-1 .95-1 1.66 0 1.07.87 1.94 1.94 1.94s1.94-.87 1.94-1.94c0-.71-.42-1.33-1-1.66V13.8c1.37-.2 2.59-.96 3.32-2.07l1.05 1.16C16.33 13.87 17.64 14.48 19 14.5V17l4-4-4-4v2z"/></svg>
  </button>
  <button id="btn-settings" class="btn-icon" data-tooltip="Diff & viewer settings">
    <svg viewBox="0 0 24 24"><path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/></svg>
  </button>
  <button id="btn-fullscreen" class="btn-icon" data-tooltip="Toggle Fullscreen (F11)">
    <svg viewBox="0 0 24 24"><path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z"/></svg>
  </button>
  <button id="btn-exit" class="btn-icon btn-danger" data-tooltip="Stop server and close window">
    <svg viewBox="0 0 24 24"><path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/></svg>
  </button>
</header>
<main class="workspace">
  <section class="pane-upper">
    <div class="table-scroll">
      <table class="data-table" id="commits-table">
        <thead>
          <tr>
            <th style="width: 24px;"></th>
            <th class="hash-cell">Hash</th>
            <th class="date-cell">Date</th>
            <th class="author-cell">Author</th>
            <th>Subject</th>
          </tr>
        </thead>
        <tbody id="commits-body">
          <tr><td colspan="5" style="text-align: center; color: var(--text-muted); padding: 30px;">Loading commits...</td></tr>
        </tbody>
      </table>
    </div>
  </section>
  <section class="pane-middle">
    <div class="commit-header" id="detail-header">
      <div>Select a commit above to view full details</div>
    </div>
    <div class="commit-body" id="detail-body"></div>
  </section>
  <section class="pane-lower">
    <div class="lower-header">
      <span id="files-summary">Modified Files (0)</span>
      <span>💡 Double-click any file to launch Beyond Compare / external diff</span>
    </div>
    <div class="table-scroll">
      <table class="data-table" id="files-table">
        <thead>
          <tr>
            <th style="width: 100px;">Action</th>
            <th>File Path</th>
            <th style="width: 120px;">Changes</th>
            <th style="width: 80px; text-align: right;">Action</th>
          </tr>
        </thead>
        <tbody id="files-body">
          <tr><td colspan="4" style="text-align: center; color: var(--text-muted); padding: 20px;">No commit selected</td></tr>
        </tbody>
      </table>
    </div>
  </section>
</main>
<footer class="statusbar">
  <span id="status-left">Connected to GitMap local server</span>
  <span id="toast" class="toast-message"></span>
</footer>
<div class="modal-overlay" id="compare-modal">
  <div class="modal-card">
    <h3>Compare Branches</h3>
    <div class="form-row">
      <label>Base Branch (Left)</label>
      <select id="modal-branch-a"></select>
    </div>
    <div class="form-row">
      <label>Target Branch (Right)</label>
      <select id="modal-branch-b"></select>
    </div>
    <div class="form-row">
      <label>Path Filter (optional, leave blank for repository trees)</label>
      <input type="text" id="modal-compare-path" placeholder="e.g. cli/cmdgit/compare_launcher.go">
    </div>
    <div class="modal-actions">
      <button id="modal-cancel">Cancel</button>
      <button id="modal-submit" class="btn-primary">Launch Beyond Compare</button>
    </div>
  </div>
</div>
<div class="modal-overlay" id="settings-modal">
  <div class="modal-card">
    <h3>Settings & Compare Tool</h3>
    <div class="form-row">
      <label>Preferred External Diff Tool</label>
      <select id="setting-diff-tool">
        <option value="auto">Auto-detect (Beyond Compare → VS Code → Meld → WinMerge)</option>
        <option value="bcomp">Beyond Compare (bcomp / BCompare.exe)</option>
        <option value="vscode">Visual Studio Code (code -d)</option>
        <option value="meld">Meld</option>
        <option value="winmerge">WinMerge</option>
      </select>
    </div>
    <div class="form-row">
      <label>Auto-Close Server on Window Exit</label>
      <select id="setting-auto-exit">
        <option value="1">Enabled (Stop server immediately)</option>
        <option value="0">Disabled (Keep background daemon running)</option>
      </select>
    </div>
    <div class="modal-actions">
      <button id="modal-settings-cancel">Cancel</button>
      <button id="modal-settings-save" class="btn-primary">Save Settings</button>
    </div>
  </div>
</div>
<div id="context-menu" class="context-menu" style="display:none; position:fixed; z-index:2000;"></div>
`

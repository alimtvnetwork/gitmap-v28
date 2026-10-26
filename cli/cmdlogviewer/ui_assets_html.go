package cmdlogviewer

// showLogMarkupHead contains HTML structure and embedded styles.
const showLogMarkupHead = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GitMap — Show Log</title>
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
  }
  input:focus, select:focus { border-color: var(--accent); }
  button {
    cursor: pointer;
    background: #313244;
    transition: background 0.15s, border-color 0.15s;
  }
  button:hover { background: var(--accent); color: #11111b; }
  .btn-primary { background: var(--accent); color: #11111b; font-weight: 600; }
  .btn-primary:hover { background: var(--accent-hover); }
  main.workspace { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
  .pane-upper {
    flex: 4; min-height: 140px; display: flex; flex-direction: column;
    background: var(--pane-bg); border-bottom: 1px solid var(--border); overflow: hidden;
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
  table.data-table tr { cursor: pointer; user-select: none; }
  table.data-table tr:hover { background: rgba(255, 255, 255, 0.04); }
  table.data-table tr.selected { background: var(--sel-bg); color: #fff; }
  .mono { font-family: var(--font-mono); }
  .hash-cell { color: var(--yellow); width: 90px; }
  .date-cell { color: var(--text-muted); width: 100px; }
  .author-cell { color: var(--green); width: 140px; }
  .pane-middle {
    flex: 3; min-height: 110px; background: var(--card-bg);
    border-bottom: 1px solid var(--border); display: flex; flex-direction: column; overflow: hidden;
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
  }
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
  <div class="brand">⚡ GitMap Show Log</div>
  <div class="path-badge" id="path-indicator" title="Active Path">Loading path...</div>
  <select id="branch-select" title="Active Branch"></select>
  <input type="text" id="commit-search" class="search-box" placeholder="Filter commits (hash, author, message)...">
  <button id="btn-refresh" title="Refresh commit list">↻ Refresh</button>
  <button id="btn-compare-branches" title="Compare two branches via Beyond Compare">Branch Comparison...</button>
  <button id="btn-exit" class="btn-danger" title="Stop local server and close window">✕ Exit</button>
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
<div id="context-menu" class="context-menu" style="display:none; position:fixed; z-index:2000;"></div>
`

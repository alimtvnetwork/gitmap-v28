package cmddb

const visualizerHeadAndStyle = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GitMap Database Visualizer</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { background: #090d16; color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 13px; line-height: 1.4; display: flex; flex-direction: column; height: 100vh; overflow: hidden; }
header { background: #0f172a; border-bottom: 1px solid #1e293b; padding: 12px 20px; display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; }
.logo { display: flex; align-items: center; gap: 10px; font-weight: 700; font-size: 16px; color: #38bdf8; }
.db-select-box { display: flex; align-items: center; gap: 8px; }
.db-select-box select { background: #1e293b; border: 1px solid #334155; color: #f8fafc; border-radius: 4px; padding: 6px 12px; font-size: 13px; outline: none; }
.header-right { display: flex; align-items: center; gap: 16px; }
.status-pill { font-size: 11px; padding: 3px 8px; border-radius: 9999px; background: rgba(16,185,129,0.15); color: #10b981; border: 1px solid rgba(16,185,129,0.3); }
.header-link { color: #38bdf8; text-decoration: none; font-size: 12px; font-weight: 500; }
.header-link:hover { text-decoration: underline; }
.app-layout { display: flex; flex: 1; overflow: hidden; }
aside { width: 280px; background: #0c111d; border-right: 1px solid #1e293b; display: flex; flex-direction: column; flex-shrink: 0; }
.aside-search { padding: 12px; border-bottom: 1px solid #1e293b; }
.aside-search input { width: 100%; background: #1e293b; border: 1px solid #334155; color: #f8fafc; border-radius: 4px; padding: 6px 10px; font-size: 12px; }
.aside-list { flex: 1; overflow-y: auto; padding: 10px 0; }
.aside-section-title { font-size: 11px; font-weight: 700; color: #64748b; padding: 6px 16px; text-transform: uppercase; letter-spacing: 0.05em; display: flex; justify-content: space-between; align-items: center; }
.table-item { padding: 7px 16px; display: flex; justify-content: space-between; align-items: center; cursor: pointer; color: #cbd5e1; user-select: none; }
.table-item:hover { background: #1e293b; color: #f8fafc; }
.table-item.active { background: #1e293b; color: #38bdf8; font-weight: 600; border-left: 3px solid #38bdf8; }
.item-badge { font-size: 10px; background: #1e293b; color: #94a3b8; padding: 1px 6px; border-radius: 4px; }
.table-item.active .item-badge { background: #0284c7; color: #f0f9ff; }
.aside-actions { padding: 12px; border-top: 1px solid #1e293b; display: flex; gap: 8px; }
.btn-sm { flex: 1; background: #1e293b; border: 1px solid #334155; color: #cbd5e1; border-radius: 4px; padding: 6px; font-size: 11px; font-weight: 600; cursor: pointer; text-align: center; }
.btn-sm:hover { background: #334155; color: #fff; }
main.content-area { flex: 1; display: flex; flex-direction: column; background: #090d16; overflow: hidden; }
.tabs-nav { background: #0f172a; border-bottom: 1px solid #1e293b; display: flex; padding: 0 16px; flex-shrink: 0; }
.tab-btn { padding: 12px 16px; color: #94a3b8; font-size: 13px; font-weight: 600; border-bottom: 2px solid transparent; cursor: pointer; user-select: none; }
.tab-btn:hover { color: #f8fafc; }
.tab-btn.active { color: #38bdf8; border-bottom-color: #38bdf8; }
.tab-pane { flex: 1; display: none; flex-direction: column; overflow: hidden; padding: 16px 20px; }
.tab-pane.active { display: flex; }
.toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; gap: 12px; flex-shrink: 0; }
.filter-bar { display: flex; gap: 8px; align-items: center; flex: 1; }
.filter-bar select, .filter-bar input { background: #0f172a; border: 1px solid #1e293b; color: #f8fafc; border-radius: 4px; padding: 6px 10px; font-size: 12px; }
.btn-primary { background: #0284c7; color: #fff; border: none; border-radius: 4px; padding: 6px 14px; font-size: 12px; font-weight: 600; cursor: pointer; }
.btn-primary:hover { background: #0369a1; }
.btn-danger { background: #b91c1c; color: #fff; border: none; border-radius: 4px; padding: 4px 8px; font-size: 11px; cursor: pointer; }
.btn-danger:hover { background: #991b1b; }
.table-container { flex: 1; overflow: auto; border: 1px solid #1e293b; border-radius: 6px; background: #0f172a; }
table.data-table { width: 100%; border-collapse: collapse; text-align: left; font-size: 12px; }
table.data-table th { background: #1e293b; color: #94a3b8; font-weight: 600; padding: 8px 12px; border-bottom: 1px solid #334155; position: sticky; top: 0; z-index: 10; cursor: pointer; }
table.data-table th:hover { color: #f8fafc; }
table.data-table td { padding: 8px 12px; border-bottom: 1px solid #1e293b; color: #e2e8f0; font-family: monospace; }
table.data-table tr:hover td { background: #141f36; }
.pagination { display: flex; justify-content: space-between; align-items: center; margin-top: 12px; font-size: 12px; color: #94a3b8; flex-shrink: 0; }
.sql-editor { width: 100%; height: 140px; background: #0f172a; border: 1px solid #1e293b; border-radius: 6px; color: #f8fafc; font-family: monospace; font-size: 13px; padding: 12px; resize: none; margin-bottom: 10px; }
.telemetry { font-size: 12px; color: #38bdf8; font-weight: 600; display: inline-block; margin-left: 12px; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); display: none; align-items: center; justify-content: center; z-index: 100; }
.modal-card { background: #0f172a; border: 1px solid #334155; border-radius: 8px; width: 480px; max-width: 95%; padding: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.5); }
.modal-card h3 { font-size: 16px; margin-bottom: 14px; color: #38bdf8; }
.modal-form { display: flex; flex-direction: column; gap: 10px; }
.modal-form label { font-size: 11px; color: #94a3b8; font-weight: 600; }
.modal-form input, .modal-form select, .modal-form textarea { background: #1e293b; border: 1px solid #334155; color: #f8fafc; border-radius: 4px; padding: 6px 10px; font-size: 12px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 14px; }

.cell-badge { display: inline-flex; align-items: center; gap: 5px; padding: 2px 8px; border-radius: 4px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 11px; font-weight: 600; cursor: pointer; border: 1px solid transparent; user-select: none; transition: all 0.15s ease; }
.cell-badge:hover { filter: brightness(1.25); box-shadow: 0 0 8px rgba(56, 189, 248, 0.2); }
.badge-blob { background: rgba(168, 85, 247, 0.18); color: #c084fc; border-color: rgba(168, 85, 247, 0.4); }
.badge-json { background: rgba(56, 189, 248, 0.18); color: #38bdf8; border-color: rgba(56, 189, 248, 0.4); }
.badge-text { background: rgba(251, 191, 36, 0.18); color: #fbbf24; border-color: rgba(251, 191, 36, 0.4); }
.badge-bool { display: inline-flex; align-items: center; padding: 2px 6px; border-radius: 4px; font-size: 10px; font-weight: 700; letter-spacing: 0.04em; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; }
.badge-bool.true, .badge-bool-true { background: rgba(16, 185, 129, 0.2); color: #10b981; border: 1px solid rgba(16, 185, 129, 0.4); }
.badge-bool.false, .badge-bool-false { background: rgba(239, 68, 68, 0.16); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.38); }
.cell-truncated { max-width: 260px; display: inline-block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: middle; cursor: pointer; }
.cell-truncated:hover { color: #38bdf8; }
.cell-null { color: #64748b; font-style: italic; font-size: 11px; }

.inspector-modal-card { width: 840px; max-width: 96vw; max-height: 88vh; display: flex; flex-direction: column; }
.modal-header-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #1e293b; }
.modal-header-bar h3 { font-size: 15px; color: #38bdf8; display: flex; align-items: center; gap: 8px; }
.inspect-nav-tabs { display: flex; gap: 6px; border-bottom: 1px solid #1e293b; padding-bottom: 8px; margin-bottom: 12px; }
.inspect-tab-btn { background: #1e293b; border: 1px solid #334155; color: #94a3b8; border-radius: 4px; padding: 6px 14px; font-size: 12px; font-weight: 600; cursor: pointer; transition: all 0.15s ease; }
.inspect-tab-btn:hover { color: #f8fafc; background: #273549; }
.inspect-tab-btn.active { background: #0284c7; color: #fff; border-color: #0284c7; }
.inspect-tab-content { flex: 1; overflow: hidden; display: flex; flex-direction: column; min-height: 340px; }
.inspect-pane { flex: 1; display: none; overflow: auto; }
.inspect-pane.active { display: flex; flex-direction: column; }

.hex-table-wrap { max-height: 440px; overflow: auto; background: #060911; border: 1px solid #1e293b; border-radius: 6px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 12px; line-height: 1.5; padding: 8px 12px; }
.hex-grid-table { width: 100%; border-collapse: collapse; }
.hex-grid-table td { padding: 2px 8px; border: none; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; }
.hex-offset-col { color: #64748b; user-select: none; border-right: 1px solid #1e293b; width: 90px; }
.hex-bytes-col { color: #38bdf8; letter-spacing: 0.05em; padding-left: 12px; border-right: 1px solid #1e293b; width: 380px; }
.hex-ascii-col { color: #a5b4fc; padding-left: 12px; white-space: pre; }

.json-viewer-wrap { max-height: 440px; overflow: auto; background: #060911; border: 1px solid #1e293b; border-radius: 6px; padding: 12px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 12px; line-height: 1.5; }
.json-key { color: #38bdf8; }
.json-string { color: #34d399; }
.json-number { color: #fbbf24; }
.json-boolean { color: #f472b6; }
.json-null { color: #64748b; font-style: italic; }
.json-toggle { cursor: pointer; user-select: none; color: #94a3b8; display: inline-block; width: 14px; }
.json-toggle:hover { color: #f8fafc; }

.text-viewer-wrap { max-height: 440px; overflow: auto; background: #060911; border: 1px solid #1e293b; border-radius: 6px; padding: 12px; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 12px; line-height: 1.5; color: #f8fafc; white-space: pre-wrap; word-break: break-all; margin: 0; }
</style>
</head>
`

const visualizerBodyHtml = `<body>
<header>
  <div class="logo">GitMap DB Visualizer</div>
  <div class="db-select-box">
    <label for="db-select" style="color:#94a3b8; font-size:12px;">Database:</label>
    <select id="db-select"></select>
  </div>
  <div class="header-right">
    <span class="status-pill" id="conn-status">Online</span>
    <a href="/swagger" class="header-link" target="_blank">Swagger API &rarr;</a>
  </div>
</header>
<div class="app-layout">
  <aside>
    <div class="aside-search">
      <input type="text" id="aside-filter" placeholder="Filter tables or columns...">
    </div>
    <div class="aside-list" id="aside-items">
      <div class="aside-section-title"><span>Tables</span><span id="tab-count-badge">0</span></div>
      <div id="table-entries"></div>
      <div class="aside-section-title" style="margin-top:14px;"><span>Views</span><span id="view-count-badge">0</span></div>
      <div id="view-entries"></div>
    </div>
    <div class="aside-actions">
      <button class="btn-sm" id="btn-new-table">+ Table</button>
      <button class="btn-sm" id="btn-new-view">+ View</button>
    </div>
  </aside>
  <main class="content-area">
    <div class="tabs-nav">
      <div class="tab-btn active" data-tab="tab-browser">Table Browser</div>
      <div class="tab-btn" data-tab="tab-schema">Schema &amp; Columns</div>
      <div class="tab-btn" data-tab="tab-views">Views Builder</div>
      <div class="tab-btn" data-tab="tab-sql">SQL Console</div>
      <div class="tab-btn" data-tab="tab-migrate">Migrations &amp; Cloning</div>
    </div>
    <div class="tab-pane active" id="tab-browser">
      <div class="toolbar">
        <div class="filter-bar">
          <select id="filter-col"><option value="">-- Filter Column --</option></select>
          <select id="filter-op"><option value="=">=</option><option value="LIKE">LIKE</option><option value="!=">!=</option><option value=">">&gt;</option><option value="<">&lt;</option></select>
          <input type="text" id="filter-val" placeholder="Value...">
          <button class="btn-primary" id="btn-filter">Filter</button>
          <button class="btn-sm" id="btn-reset-filter" style="flex:none; padding:6px 10px;">Reset</button>
        </div>
        <button class="btn-primary" id="btn-insert-row">+ Insert Row</button>
      </div>
      <div class="table-container" id="grid-container">
        <table class="data-table" id="data-grid"><thead></thead><tbody></tbody></table>
      </div>
      <div class="pagination">
        <span id="page-info">Showing 0 records</span>
        <div style="display:flex; gap:8px; align-items:center;">
          <button class="btn-sm" id="btn-prev" style="flex:none; padding:4px 10px;">&lt; Prev</button>
          <span style="display:flex; align-items:center; gap:4px; font-size:12px; color:#94a3b8;">
            Page
            <input type="number" id="page-input" min="1" value="1" style="width:48px; background:#0f172a; border:1px solid #334155; color:#f8fafc; border-radius:4px; padding:3px 6px; text-align:center; font-size:12px;">
            <span id="page-max-info">/ 1</span>
          </span>
          <span id="page-num" style="display:none;">1</span>
          <button class="btn-sm" id="btn-next" style="flex:none; padding:4px 10px;">Next &gt;</button>
          <select id="page-size" style="background:#0f172a; border:1px solid #1e293b; color:#f8fafc; border-radius:4px; padding:4px 8px;">
            <option value="25">25 / page</option>
            <option value="50" selected>50 / page</option>
            <option value="100">100 / page</option>
          </select>
        </div>
      </div>
    </div>
    <div class="tab-pane" id="tab-schema">
      <div class="toolbar">
        <h3 id="schema-table-title" style="color:#38bdf8;">Schema</h3>
        <button class="btn-primary" id="btn-add-column">+ Add Column</button>
      </div>
      <div class="table-container" style="margin-bottom:14px;">
        <table class="data-table" id="schema-grid"><thead><tr><th>CID</th><th>Name</th><th>Data Type</th><th>Nullable</th><th>Default</th><th>Primary Key</th></tr></thead><tbody></tbody></table>
      </div>
      <h4 style="color:#94a3b8; margin-bottom:6px;">Table DDL</h4>
      <pre id="ddl-box" style="background:#0f172a; border:1px solid #1e293b; border-radius:6px; padding:12px; font-family:monospace; color:#38bdf8; overflow-x:auto;"></pre>
    </div>
    <div class="tab-pane" id="tab-views">
      <div class="toolbar">
        <h3 style="color:#38bdf8;">Dynamic Views</h3>
        <button class="btn-primary" id="btn-create-view-modal">+ Create View</button>
      </div>
      <div class="table-container">
        <table class="data-table" id="views-grid"><thead><tr><th>View Name</th><th>Source SQL</th><th>Actions</th></tr></thead><tbody></tbody></table>
      </div>
    </div>
    <div class="tab-pane" id="tab-sql">
      <div style="margin-bottom:10px;">
        <textarea class="sql-editor" id="sql-input" placeholder="SELECT * FROM Table LIMIT 10; (Ctrl+Enter to execute)"></textarea>
        <div>
          <button class="btn-primary" id="btn-exec-sql">Execute SQL</button>
          <button class="btn-sm" id="btn-export-csv" style="flex:none; padding:6px 12px; margin-left:8px;">Export CSV</button>
          <button class="btn-sm" id="btn-export-json" style="flex:none; padding:6px 12px; margin-left:4px;">Export JSON</button>
          <span class="telemetry" id="sql-telemetry"></span>
        </div>
      </div>
      <div class="table-container" id="sql-result-container">
        <table class="data-table" id="sql-grid"><thead></thead><tbody></tbody></table>
      </div>
    </div>
    <div class="tab-pane" id="tab-migrate">
      <div style="display:grid; grid-template-columns: 1fr 1fr; gap:20px; height:100%;">
        <div style="background:#0f172a; border:1px solid #1e293b; border-radius:6px; padding:16px;">
          <h3 style="color:#38bdf8; margin-bottom:12px;">Clone Database</h3>
          <div class="modal-form">
            <label>Destination File Path</label>
            <input type="text" id="clone-path" placeholder="./data/backup.db">
            <label>New Database Alias</label>
            <input type="text" id="clone-alias" placeholder="backup_db">
            <button class="btn-primary" id="btn-exec-clone" style="margin-top:10px;">Clone &amp; Mount</button>
          </div>
        </div>
        <div style="background:#0f172a; border:1px solid #1e293b; border-radius:6px; padding:16px; display:flex; flex-direction:column;">
          <h3 style="color:#38bdf8; margin-bottom:12px;">Run Migration</h3>
          <div class="modal-form" style="flex:1; display:flex; flex-direction:column;">
            <label>Version Tag</label>
            <input type="text" id="migrate-ver" placeholder="v1.0.1">
            <label>SQL Script</label>
            <textarea id="migrate-sql" style="flex:1; background:#1e293b; border:1px solid #334155; color:#f8fafc; font-family:monospace; padding:8px;" placeholder="ALTER TABLE ...;"></textarea>
            <button class="btn-primary" id="btn-exec-migrate" style="margin-top:10px;">Apply Migration</button>
          </div>
        </div>
      </div>
    </div>
  </main>
</div>
<div class="modal-overlay" id="generic-modal">
  <div class="modal-card">
    <h3 id="modal-title">Modal</h3>
    <div class="modal-form" id="modal-form-content"></div>
    <div class="modal-actions">
      <button class="btn-sm" id="btn-modal-cancel">Cancel</button>
      <button class="btn-primary" id="btn-modal-submit">Save</button>
    </div>
  </div>
</div>
<div class="modal-overlay" id="inspector-modal">
  <div class="modal-card inspector-modal-card">
    <div class="modal-header-bar">
      <div style="display:flex; align-items:center; gap:8px;">
        <h3 style="margin:0;">Cell Inspector</h3>
        <span id="inspect-modal-col" class="item-badge" style="background:#1e293b; color:#38bdf8;">COL</span>
        <span id="inspect-modal-badge" class="item-badge">[BLOB]</span>
      </div>
      <button class="btn-sm" id="btn-inspect-close" style="flex:none; padding:2px 8px;">&times;</button>
    </div>
    <div class="inspect-nav-tabs">
      <button class="inspect-tab-btn active" data-inspect-tab="inspect-tab-hex">Hex Dump</button>
      <button class="inspect-tab-btn" data-inspect-tab="inspect-tab-json">Formatted JSON</button>
      <button class="inspect-tab-btn" data-inspect-tab="inspect-tab-text">Decoded String</button>
    </div>
    <div class="inspect-tab-content">
      <div class="inspect-pane active" id="inspect-tab-hex">
        <div class="hex-table-wrap">
          <table class="hex-grid-table">
            <thead style="border-bottom:1px solid #1e293b; color:#64748b; font-size:11px;">
              <tr>
                <th style="padding:4px 8px; text-align:left; width:90px;">OFFSET</th>
                <th style="padding:4px 8px; text-align:left; width:380px;">HEX BYTES</th>
                <th style="padding:4px 8px; text-align:left;">ASCII</th>
              </tr>
            </thead>
            <tbody id="inspect-hex-body"></tbody>
          </table>
        </div>
      </div>
      <div class="inspect-pane" id="inspect-tab-json">
        <div class="json-viewer-wrap" id="inspect-json-body"></div>
      </div>
      <div class="inspect-pane" id="inspect-tab-text">
        <pre class="text-viewer-wrap" id="inspect-text-body"></pre>
      </div>
    </div>
    <div class="modal-actions" style="margin-top:12px; justify-content:space-between;">
      <span id="inspect-copy-status" style="color:#10b981; font-size:11px; align-self:center;"></span>
      <div style="display:flex; gap:8px;">
        <button class="btn-sm" id="btn-inspect-copy-hex">Copy Hex</button>
        <button class="btn-sm" id="btn-inspect-copy-json">Copy JSON</button>
        <button class="btn-sm" id="btn-inspect-copy-text">Copy String</button>
        <button class="btn-primary" id="btn-inspect-done">Done</button>
      </div>
    </div>
  </div>
</div>
`

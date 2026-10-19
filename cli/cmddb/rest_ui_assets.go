package cmddb

// SwaggerUIHtml is the standalone, zero-CDN embedded Swagger UI single-page application.
const SwaggerUIHtml = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>GitMap Database REST API - Swagger UI</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { background: #090d16; color: #f8fafc; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 14px; line-height: 1.5; }
header { background: #0f172a; border-bottom: 1px solid #1e293b; padding: 16px 24px; display: flex; align-items: center; justify-content: space-between; }
.brand { display: flex; align-items: center; gap: 12px; font-weight: 700; font-size: 18px; color: #38bdf8; }
.brand span { color: #94a3b8; font-size: 12px; font-weight: 400; padding: 2px 8px; border: 1px solid #334155; border-radius: 9999px; }
.links a { color: #38bdf8; text-decoration: none; font-size: 13px; margin-left: 16px; }
.links a:hover { text-decoration: underline; }
main { max-width: 1200px; margin: 24px auto; padding: 0 16px; }
.info-card { background: #0f172a; border: 1px solid #1e293b; border-radius: 8px; padding: 20px; margin-bottom: 24px; }
.info-card h1 { font-size: 20px; margin-bottom: 8px; color: #f8fafc; }
.info-card p { color: #94a3b8; margin-bottom: 12px; }
.tag-group { margin-bottom: 24px; }
.tag-title { font-size: 16px; font-weight: 700; color: #e2e8f0; margin-bottom: 12px; display: flex; align-items: center; gap: 8px; }
.op-card { background: #0f172a; border: 1px solid #1e293b; border-radius: 6px; margin-bottom: 10px; overflow: hidden; }
.op-summary { display: flex; align-items: center; gap: 12px; padding: 12px 16px; cursor: pointer; user-select: none; }
.op-summary:hover { background: #1e293b; }
.method-badge { font-weight: 700; font-size: 11px; padding: 4px 8px; border-radius: 4px; text-transform: uppercase; width: 68px; text-align: center; }
.method-get { background: #0369a1; color: #e0f2fe; }
.method-post { background: #15803d; color: #dcfce7; }
.method-put { background: #b45309; color: #fef3c7; }
.method-delete { background: #b91c1c; color: #fee2e2; }
.op-path { font-family: monospace; font-size: 13px; font-weight: 600; color: #f1f5f9; flex: 1; }
.op-desc { color: #94a3b8; font-size: 12px; }
.op-body { display: none; padding: 16px; border-top: 1px solid #1e293b; background: #090d16; }
.op-body.open { display: block; }
.form-group { margin-bottom: 12px; }
.form-group label { display: block; font-size: 12px; color: #94a3b8; margin-bottom: 4px; font-weight: 600; }
.form-group input, .form-group textarea { width: 100%; background: #0f172a; border: 1px solid #334155; border-radius: 4px; color: #f8fafc; padding: 8px 10px; font-family: monospace; font-size: 12px; }
.form-group textarea { height: 90px; resize: vertical; }
.btn-exec { background: #0284c7; color: #fff; border: none; border-radius: 4px; padding: 8px 16px; font-weight: 600; cursor: pointer; font-size: 12px; }
.btn-exec:hover { background: #0369a1; }
.resp-box { margin-top: 14px; background: #020617; border: 1px solid #1e293b; border-radius: 4px; padding: 12px; }
.resp-header { display: flex; justify-content: space-between; font-size: 12px; color: #94a3b8; margin-bottom: 6px; }
.resp-code { font-weight: 700; color: #38bdf8; }
.resp-code.err { color: #ef4444; }
pre.resp-content { font-family: monospace; font-size: 12px; color: #cbd5e1; overflow-x: auto; white-space: pre-wrap; max-height: 240px; }
</style>
</head>
<body>
<header>
  <div class="brand">GitMap REST <span>OpenAPI 3.0.3</span></div>
  <div class="links">
    <a href="/swagger/doc.json" target="_blank">Raw Spec (doc.json)</a>
    <a href="/ui">Database Visualizer UI</a>
  </div>
</header>
<main>
  <div class="info-card">
    <h1 id="spec-title">GitMap Database REST API</h1>
    <p id="spec-desc">Dynamic OpenAPI 3.0 specification generated on-the-fly from active SQLite databases.</p>
  </div>
  <div id="operations-container">Loading OpenAPI specification...</div>
</main>
<script>
(function() {
  fetch('/swagger/doc.json')
    .then(function(r) { return r.json(); })
    .then(function(spec) { renderSpec(spec); })
    .catch(function(err) {
      document.getElementById('operations-container').innerHTML = '<p style="color:#ef4444">Failed to load OpenAPI spec: ' + err + '</p>';
    });

  function renderSpec(spec) {
    if (spec.info) {
      if (spec.info.title) document.getElementById('spec-title').textContent = spec.info.title;
      if (spec.info.description) document.getElementById('spec-desc').textContent = spec.info.description;
    }
    var container = document.getElementById('operations-container');
    container.innerHTML = '';
    var groups = groupPathsByTag(spec);
    Object.keys(groups).forEach(function(tag) {
      var tagEl = document.createElement('div');
      tagEl.className = 'tag-group';
      tagEl.innerHTML = '<div class="tag-title">' + tag + '</div>';
      groups[tag].forEach(function(op) {
        tagEl.appendChild(createOpCard(op));
      });
      container.appendChild(tagEl);
    });
  }

  function groupPathsByTag(spec) {
    var groups = {};
    if (!spec.paths) return groups;
    Object.keys(spec.paths).forEach(function(path) {
      var item = spec.paths[path];
      ['get', 'post', 'put', 'delete', 'patch'].forEach(function(m) {
        if (!item[m]) return;
        var op = item[m];
        var tag = (op.tags && op.tags.length > 0) ? op.tags[0] : 'General';
        if (!groups[tag]) groups[tag] = [];
        groups[tag].push({ method: m, path: path, op: op });
      });
    });
    return groups;
  }

  function createOpCard(item) {
    var card = document.createElement('div');
    card.className = 'op-card';
    var summary = document.createElement('div');
    summary.className = 'op-summary';
    summary.innerHTML = '<span class="method-badge method-' + item.method + '">' + item.method + '</span>' +
      '<span class="op-path">' + item.path + '</span>' +
      '<span class="op-desc">' + (item.op.summary || '') + '</span>';
    var body = document.createElement('div');
    body.className = 'op-body';
    summary.onclick = function() { body.classList.toggle('open'); };

    var paramsHtml = '';
    var pathParams = [];
    var queryParams = [];
    if (item.op.parameters) {
      item.op.parameters.forEach(function(p) {
        if (p.in === 'path') pathParams.push(p);
        else queryParams.push(p);
      });
    }

    if (pathParams.length > 0 || queryParams.length > 0) {
      paramsHtml += '<div class="form-group"><label>Parameters</label>';
      pathParams.forEach(function(p) {
        paramsHtml += '<input placeholder="{' + p.name + '} (path param)" data-name="' + p.name + '" data-in="path" style="margin-bottom:6px;">';
      });
      queryParams.forEach(function(p) {
        paramsHtml += '<input placeholder="' + p.name + ' (query param)" data-name="' + p.name + '" data-in="query" style="margin-bottom:6px;">';
      });
      paramsHtml += '</div>';
    }

    var bodyHtml = '';
    if (item.method === 'post' || item.method === 'put' || item.method === 'patch') {
      bodyHtml = '<div class="form-group"><label>Request Body (JSON)</label><textarea placeholder="{\n  \"field\": \"value\"\n}"></textarea></div>';
    }

    body.innerHTML = paramsHtml + bodyHtml +
      '<button class="btn-exec">Execute</button>' +
      '<div class="resp-box" style="display:none;"><div class="resp-header"><span>Response</span><span class="resp-code"></span></div><pre class="resp-content"></pre></div>';

    var btn = body.querySelector('.btn-exec');
    btn.onclick = function() { executeRequest(item, body); };
    card.appendChild(summary);
    card.appendChild(body);
    return card;
  }

  function executeRequest(item, bodyEl) {
    var url = item.path;
    var queryParts = [];
    bodyEl.querySelectorAll('input').forEach(function(inp) {
      var val = inp.value.trim();
      var name = inp.getAttribute('data-name');
      var where = inp.getAttribute('data-in');
      if (where === 'path') {
        url = url.replace('{' + name + '}', encodeURIComponent(val));
      } else if (where === 'query' && val) {
        queryParts.push(encodeURIComponent(name) + '=' + encodeURIComponent(val));
      }
    });
    if (queryParts.length > 0) {
      url += (url.indexOf('?') >= 0 ? '&' : '?') + queryParts.join('&');
    }
    var options = { method: item.method.toUpperCase(), headers: {} };
    var txt = bodyEl.querySelector('textarea');
    if (txt && txt.value.trim()) {
      options.headers['Content-Type'] = 'application/json';
      options.body = txt.value.trim();
    }
    var respBox = bodyEl.querySelector('.resp-box');
    var codeEl = bodyEl.querySelector('.resp-code');
    var contentEl = bodyEl.querySelector('.resp-content');
    respBox.style.display = 'block';
    codeEl.textContent = 'Executing...';
    contentEl.textContent = '';
    var start = Date.now();
    fetch(url, options)
      .then(function(res) {
        var duration = Date.now() - start;
        codeEl.textContent = res.status + ' ' + res.statusText + ' (' + duration + 'ms)';
        codeEl.className = 'resp-code' + (res.status >= 400 ? ' err' : '');
        return res.text();
      })
      .then(function(text) {
        try {
          var j = JSON.parse(text);
          contentEl.textContent = JSON.stringify(j, null, 2);
        } catch(e) {
          contentEl.textContent = text;
        }
      })
      .catch(function(err) {
        codeEl.textContent = 'Error';
        codeEl.className = 'resp-code err';
        contentEl.textContent = String(err);
      });
  }
})();
</script>
</body>
</html>`

// VisualizerUIHtml is the embedded, standalone, zero-CDN custom Database Visualizer single-page application.
const VisualizerUIHtml = `<!DOCTYPE html>
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
</style>
</head>
<body>
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
          <span id="page-num" style="font-weight:600; color:#f8fafc;">1</span>
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
<script>
(function() {
  var state = { activeDb: '', activeTable: '', page: 1, limit: 50, tables: [], columns: [], lastSqlRows: [] };
  init();

  function init() {
    setupTabs();
    loadDatabases();
    setupEvents();
  }

  function setupTabs() {
    document.querySelectorAll('.tab-btn').forEach(function(btn) {
      btn.onclick = function() {
        document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
        document.querySelectorAll('.tab-pane').forEach(function(p) { p.classList.remove('active'); });
        btn.classList.add('active');
        var target = document.getElementById(btn.getAttribute('data-tab'));
        if (target) target.classList.add('active');
        if (btn.getAttribute('data-tab') === 'tab-schema') loadSchema();
        if (btn.getAttribute('data-tab') === 'tab-views') loadViews();
      };
    });
  }

  function loadDatabases() {
    fetch('/api/v1/databases')
      .then(function(r) { return r.json(); })
      .then(function(env) {
        var dbs = (env.Results && Array.isArray(env.Results)) ? env.Results : [];
        var select = document.getElementById('db-select');
        select.innerHTML = '';
        dbs.forEach(function(d) {
          var opt = document.createElement('option');
          opt.value = d.Alias;
          opt.textContent = d.Alias + ' (' + (d.TableCount || 0) + ' tables)';
          select.appendChild(opt);
        });
        if (dbs.length > 0) {
          state.activeDb = dbs[0].Alias;
          loadTables();
        }
      });
  }

  function loadTables() {
    if (!state.activeDb) return;
    fetch('/api/v1/db/' + state.activeDb + '/tables')
      .then(function(r) { return r.json(); })
      .then(function(env) {
        state.tables = env.Results || [];
        renderSidebar();
        if (state.tables.length > 0) {
          selectTable(state.tables[0].Name);
        }
      });
  }

  function renderSidebar() {
    var tabBox = document.getElementById('table-entries');
    var viewBox = document.getElementById('view-entries');
    tabBox.innerHTML = '';
    viewBox.innerHTML = '';
    var tabCount = 0; var viewCount = 0;
    state.tables.forEach(function(t) {
      var item = document.createElement('div');
      item.className = 'table-item' + (t.Name === state.activeTable ? ' active' : '');
      item.innerHTML = '<span>' + t.Name + '</span><span class="item-badge">' + (t.RowCount || 0) + '</span>';
      item.onclick = function() { selectTable(t.Name); };
      if (t.Type === 'view') {
        viewBox.appendChild(item);
        viewCount++;
      } else {
        tabBox.appendChild(item);
        tabCount++;
      }
    });
    document.getElementById('tab-count-badge').textContent = tabCount;
    document.getElementById('view-count-badge').textContent = viewCount;
  }

  function selectTable(name) {
    state.activeTable = name;
    state.page = 1;
    renderSidebar();
    loadColumns();
    loadGridData();
  }

  function loadColumns() {
    if (!state.activeDb || !state.activeTable) return;
    fetch('/api/v1/db/' + state.activeDb + '/schema?table=' + state.activeTable)
      .then(function(r) { return r.json(); })
      .then(function(env) {
        state.columns = [];
        if (env.Results && env.Results.length > 0) {
          state.columns = env.Results[0].Columns || [];
          document.getElementById('ddl-box').textContent = env.Results[0].Ddl || '';
        }
        var filterSel = document.getElementById('filter-col');
        filterSel.innerHTML = '<option value="">-- Filter Column --</option>';
        state.columns.forEach(function(c) {
          var o = document.createElement('option');
          o.value = c.Name;
          o.textContent = c.Name;
          filterSel.appendChild(o);
        });
      });
  }

  function loadGridData() {
    if (!state.activeDb || !state.activeTable) return;
    var offset = (state.page - 1) * state.limit;
    var url = '/api/v1/db/' + state.activeDb + '/' + state.activeTable + '?limit=' + state.limit + '&offset=' + offset;
    var filterCol = document.getElementById('filter-col').value;
    var filterVal = document.getElementById('filter-val').value;
    if (filterCol && filterVal) {
      url += '&' + encodeURIComponent(filterCol) + '=' + encodeURIComponent(filterVal);
    }
    fetch(url)
      .then(function(r) { return r.json(); })
      .then(function(env) {
        renderGrid(env);
      });
  }

  function renderGrid(env) {
    var rows = env.Results || [];
    var tableEl = document.getElementById('data-grid');
    var thead = tableEl.querySelector('thead');
    var tbody = tableEl.querySelector('tbody');
    thead.innerHTML = '';
    tbody.innerHTML = '';
    if (rows.length === 0) {
      tbody.innerHTML = '<tr><td style="text-align:center; padding:20px; color:#64748b;">No records found</td></tr>';
      document.getElementById('page-info').textContent = '0 records';
      return;
    }
    var cols = Object.keys(rows[0]);
    var trH = document.createElement('tr');
    cols.forEach(function(c) {
      var th = document.createElement('th');
      th.textContent = c;
      trH.appendChild(th);
    });
    var thAct = document.createElement('th');
    thAct.textContent = 'Actions';
    trH.appendChild(thAct);
    thead.appendChild(trH);

    rows.forEach(function(r) {
      var tr = document.createElement('tr');
      cols.forEach(function(c) {
        var td = document.createElement('td');
        td.textContent = r[c] !== null && r[c] !== undefined ? r[c] : 'NULL';
        tr.appendChild(td);
      });
      var tdA = document.createElement('td');
      var pkCol = state.columns.find(function(col) { return col.IsPrimaryKey; });
      var idVal = pkCol ? r[pkCol.Name] : (r['id'] || r['Id'] || r['ID']);
      tdA.innerHTML = '<button class="btn-sm" style="padding:2px 6px; margin-right:4px;" onclick="window.visualizerEditRow(' + idVal + ')">Edit</button>' +
        '<button class="btn-danger" onclick="window.visualizerDeleteRow(' + idVal + ')">Del</button>';
      tr.appendChild(tdA);
      tbody.appendChild(tr);
    });
    var total = (env.Attributes && env.Attributes.TotalRecords) ? env.Attributes.TotalRecords : rows.length;
    document.getElementById('page-info').textContent = 'Showing ' + rows.length + ' of ' + total + ' records';
    document.getElementById('page-num').textContent = state.page;
  }

  function loadSchema() {
    if (!state.activeDb || !state.activeTable) return;
    document.getElementById('schema-table-title').textContent = 'Schema: ' + state.activeTable;
    var tbody = document.querySelector('#schema-grid tbody');
    tbody.innerHTML = '';
    state.columns.forEach(function(c) {
      var tr = document.createElement('tr');
      tr.innerHTML = '<td>' + (c.ColumnId !== undefined ? c.ColumnId : '-') + '</td>' +
        '<td style="font-weight:600; color:#38bdf8;">' + c.Name + '</td>' +
        '<td>' + c.DataType + '</td>' +
        '<td>' + (c.IsNullable ? 'YES' : 'NO') + '</td>' +
        '<td>' + (c.DefaultValue !== null && c.DefaultValue !== undefined ? c.DefaultValue : 'NULL') + '</td>' +
        '<td>' + (c.IsPrimaryKey ? '<span style="color:#10b981;font-weight:700;">YES</span>' : '-') + '</td>';
      tbody.appendChild(tr);
    });
  }

  function loadViews() {
    if (!state.activeDb) return;
    fetch('/api/v1/db/' + state.activeDb + '/views')
      .then(function(r) { return r.json(); })
      .then(function(env) {
        var views = env.Results || [];
        var tbody = document.querySelector('#views-grid tbody');
        tbody.innerHTML = '';
        views.forEach(function(v) {
          var tr = document.createElement('tr');
          tr.innerHTML = '<td style="font-weight:600; color:#38bdf8;">' + (v.Name || v.ViewName) + '</td>' +
            '<td style="font-family:monospace; color:#cbd5e1;">' + (v.Sql || v.Query) + '</td>' +
            '<td><button class="btn-danger" onclick="window.visualizerDropView(\'' + (v.Name || v.ViewName) + '\')">Drop</button></td>';
          tbody.appendChild(tr);
        });
      });
  }

  function setupEvents() {
    document.getElementById('db-select').onchange = function(e) {
      state.activeDb = e.target.value;
      loadTables();
    };
    document.getElementById('btn-prev').onclick = function() {
      if (state.page > 1) { state.page--; loadGridData(); }
    };
    document.getElementById('btn-next').onclick = function() {
      state.page++; loadGridData();
    };
    document.getElementById('page-size').onchange = function(e) {
      state.limit = parseInt(e.target.value, 10);
      state.page = 1;
      loadGridData();
    };
    document.getElementById('btn-filter').onclick = function() {
      state.page = 1; loadGridData();
    };
    document.getElementById('btn-reset-filter').onclick = function() {
      document.getElementById('filter-col').value = '';
      document.getElementById('filter-val').value = '';
      state.page = 1; loadGridData();
    };
    document.getElementById('btn-exec-sql').onclick = executeSql;
    document.getElementById('sql-input').onkeydown = function(e) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        executeSql();
      }
    };
    document.getElementById('btn-insert-row').onclick = openInsertRowModal;
    document.getElementById('btn-add-column').onclick = openAddColumnModal;
    document.getElementById('btn-create-view-modal').onclick = openCreateViewModal;
    document.getElementById('btn-exec-clone').onclick = executeClone;
    document.getElementById('btn-exec-migrate').onclick = executeMigrate;
    document.getElementById('btn-modal-cancel').onclick = closeModal;
    document.getElementById('btn-export-csv').onclick = exportSqlCsv;
    document.getElementById('btn-export-json').onclick = exportSqlJson;
  }

  function executeSql() {
    var sql = document.getElementById('sql-input').value.trim();
    if (!sql || !state.activeDb) return;
    var start = Date.now();
    var isSelect = /^SELECT|^PRAGMA|^EXPLAIN/i.test(sql);
    var p;
    if (isSelect) {
      p = fetch('/api/v1/db/' + state.activeDb + '/query?sql=' + encodeURIComponent(sql)).then(function(r) { return r.json(); });
    } else {
      p = fetch('/api/v1/db/' + state.activeDb + '/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Sql: sql })
      }).then(function(r) { return r.json(); });
    }
    p.then(function(env) {
      var ms = Date.now() - start;
      document.getElementById('sql-telemetry').textContent = 'Executed in ' + ms + 'ms';
      renderSqlResult(env.Results);
    }).catch(function(err) {
      document.getElementById('sql-telemetry').textContent = 'Error: ' + err;
    });
  }

  function renderSqlResult(results) {
    state.lastSqlRows = results;
    var thead = document.querySelector('#sql-grid thead');
    var tbody = document.querySelector('#sql-grid tbody');
    thead.innerHTML = '';
    tbody.innerHTML = '';
    if (!results || !Array.isArray(results) || results.length === 0) {
      tbody.innerHTML = '<tr><td style="padding:16px; color:#94a3b8;">' + JSON.stringify(results) + '</td></tr>';
      return;
    }
    var cols = Object.keys(results[0]);
    var trH = document.createElement('tr');
    cols.forEach(function(c) {
      var th = document.createElement('th');
      th.textContent = c;
      trH.appendChild(th);
    });
    thead.appendChild(trH);
    results.forEach(function(r) {
      var tr = document.createElement('tr');
      cols.forEach(function(c) {
        var td = document.createElement('td');
        td.textContent = r[c] !== null && r[c] !== undefined ? r[c] : 'NULL';
        tr.appendChild(td);
      });
      tbody.appendChild(tr);
    });
  }

  function exportSqlCsv() {
    if (!state.lastSqlRows || !Array.isArray(state.lastSqlRows) || state.lastSqlRows.length === 0) return;
    var cols = Object.keys(state.lastSqlRows[0]);
    var lines = [cols.join(',')];
    state.lastSqlRows.forEach(function(r) {
      var line = cols.map(function(c) { return '"' + String(r[c] || '').replace(/"/g, '""') + '"'; }).join(',');
      lines.push(line);
    });
    downloadBlob(lines.join('\n'), 'query_results.csv', 'text/csv');
  }

  function exportSqlJson() {
    if (!state.lastSqlRows) return;
    downloadBlob(JSON.stringify(state.lastSqlRows, null, 2), 'query_results.json', 'application/json');
  }

  function downloadBlob(content, filename, type) {
    var blob = new Blob([content], { type: type });
    var a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = filename;
    a.click();
  }

  function openInsertRowModal() {
    document.getElementById('modal-title').textContent = 'Insert Row in ' + state.activeTable;
    var form = document.getElementById('modal-form-content');
    form.innerHTML = '';
    state.columns.forEach(function(c) {
      var wrap = document.createElement('div');
      wrap.innerHTML = '<label>' + c.Name + ' (' + c.DataType + ')' + (c.IsPrimaryKey ? ' [PK]' : '') + '</label>' +
        '<input type="text" data-col="' + c.Name + '" placeholder="' + (c.IsPrimaryKey ? 'auto' : '') + '">';
      form.appendChild(wrap);
    });
    document.getElementById('btn-modal-submit').onclick = function() {
      var body = {};
      form.querySelectorAll('input').forEach(function(inp) {
        var col = inp.getAttribute('data-col');
        var val = inp.value.trim();
        if (val) body[col] = val;
      });
      fetch('/api/v1/db/' + state.activeDb + '/' + state.activeTable, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      }).then(function() {
        closeModal();
        loadGridData();
      });
    };
    document.getElementById('generic-modal').style.display = 'flex';
  }

  function openAddColumnModal() {
    document.getElementById('modal-title').textContent = 'Add Column to ' + state.activeTable;
    var form = document.getElementById('modal-form-content');
    form.innerHTML = '<label>Column Name</label><input type="text" id="new-col-name">' +
      '<label>Data Type</label><select id="new-col-type"><option>TEXT</option><option>INTEGER</option><option>REAL</option><option>BLOB</option><option>BOOLEAN</option></select>' +
      '<label>Default Value</label><input type="text" id="new-col-dflt">';
    document.getElementById('btn-modal-submit').onclick = function() {
      var colName = document.getElementById('new-col-name').value.trim();
      var colType = document.getElementById('new-col-type').value;
      var dflt = document.getElementById('new-col-dflt').value.trim();
      fetch('/api/v1/db/' + state.activeDb + '/alter', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          TableName: state.activeTable,
          Operation: 'add_column',
          ColumnName: colName,
          ColumnType: colType,
          DefaultValue: dflt
        })
      }).then(function() {
        closeModal();
        loadColumns();
        loadSchema();
      });
    };
    document.getElementById('generic-modal').style.display = 'flex';
  }

  function openCreateViewModal() {
    document.getElementById('modal-title').textContent = 'Create Dynamic View';
    var form = document.getElementById('modal-form-content');
    form.innerHTML = '<label>View Name</label><input type="text" id="new-view-name" placeholder="v_summary">' +
      '<label>SQL Query</label><textarea id="new-view-sql" style="height:90px;" placeholder="SELECT * FROM ..."></textarea>';
    document.getElementById('btn-modal-submit').onclick = function() {
      var vName = document.getElementById('new-view-name').value.trim();
      var vSql = document.getElementById('new-view-sql').value.trim();
      fetch('/api/v1/db/' + state.activeDb + '/views', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ViewName: vName, Query: vSql })
      }).then(function() {
        closeModal();
        loadTables();
        loadViews();
      });
    };
    document.getElementById('generic-modal').style.display = 'flex';
  }

  function executeClone() {
    var p = document.getElementById('clone-path').value.trim();
    var a = document.getElementById('clone-alias').value.trim();
    if (!p || !a || !state.activeDb) return;
    fetch('/api/v1/db/' + state.activeDb + '/clone', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ DestinationPath: p, NewAlias: a, AutoMount: true })
    }).then(function(r) { return r.json(); })
      .then(function() {
        alert('Cloned successfully to ' + p + ' under alias ' + a);
        loadDatabases();
      });
  }

  function executeMigrate() {
    var ver = document.getElementById('migrate-ver').value.trim();
    var sql = document.getElementById('migrate-sql').value.trim();
    if (!ver || !sql || !state.activeDb) return;
    fetch('/api/v1/db/' + state.activeDb + '/migrate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ MigrationName: ver, UpSql: sql })
    }).then(function(r) { return r.json(); })
      .then(function() {
        alert('Migration applied successfully');
        loadTables();
      });
  }

  function closeModal() {
    document.getElementById('generic-modal').style.display = 'none';
  }

  window.visualizerDeleteRow = function(id) {
    if (!confirm('Are you sure you want to delete record #' + id + '?')) return;
    fetch('/api/v1/db/' + state.activeDb + '/' + state.activeTable + '/' + id, { method: 'DELETE' })
      .then(function() { loadGridData(); });
  };

  window.visualizerDropView = function(name) {
    if (!confirm('Drop view ' + name + '?')) return;
    fetch('/api/v1/db/' + state.activeDb + '/views/' + name, { method: 'DELETE' })
      .then(function() { loadTables(); loadViews(); });
  };

  window.visualizerEditRow = function(id) {
    fetch('/api/v1/db/' + state.activeDb + '/' + state.activeTable + '/' + id)
      .then(function(r) { return r.json(); })
      .then(function(env) {
        var rec = (env.Results && Array.isArray(env.Results)) ? env.Results[0] : env.Results;
        if (!rec) return;
        document.getElementById('modal-title').textContent = 'Edit Record #' + id;
        var form = document.getElementById('modal-form-content');
        form.innerHTML = '';
        state.columns.forEach(function(c) {
          if (c.IsPrimaryKey) return;
          var wrap = document.createElement('div');
          wrap.innerHTML = '<label>' + c.Name + '</label><input type="text" data-col="' + c.Name + '" value="' + (rec[c.Name] || '') + '">';
          form.appendChild(wrap);
        });
        document.getElementById('btn-modal-submit').onclick = function() {
          var body = {};
          form.querySelectorAll('input').forEach(function(inp) {
            body[inp.getAttribute('data-col')] = inp.value.trim();
          });
          fetch('/api/v1/db/' + state.activeDb + '/' + state.activeTable + '/' + id, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
          }).then(function() {
            closeModal();
            loadGridData();
          });
        };
        document.getElementById('generic-modal').style.display = 'flex';
      });
  };
})();
</script>
</body>
</html>`

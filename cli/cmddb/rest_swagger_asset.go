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
    <div style="font-size: 12px; color: #64748b;">Servers: <code id="server-url">/</code></div>
  </div>
  <div id="endpoints-container">Loading endpoints...</div>
</main>
<script>
fetch('/swagger/doc.json')
  .then(function(r) { return r.json(); })
  .then(function(spec) {
    document.getElementById('spec-title').textContent = spec.info.title || 'GitMap REST API';
    document.getElementById('spec-desc').textContent = spec.info.description || '';
    renderEndpoints(spec);
  })
  .catch(function(err) {
    document.getElementById('endpoints-container').innerHTML = '<div style="color:#ef4444;">Failed to load OpenAPI spec: ' + err + '</div>';
  });

function renderEndpoints(spec) {
  var container = document.getElementById('endpoints-container');
  container.innerHTML = '';
  var tags = {};
  var paths = spec.paths || {};

  Object.keys(paths).forEach(function(path) {
    var methods = paths[path];
    Object.keys(methods).forEach(function(m) {
      var op = methods[m];
      var tag = (op.tags && op.tags[0]) || 'General';
      if (!tags[tag]) tags[tag] = [];
      tags[tag].push({ method: m.toLowerCase(), path: path, op: op });
    });
  });

  Object.keys(tags).forEach(function(tag) {
    var grp = document.createElement('div');
    grp.className = 'tag-group';
    grp.innerHTML = '<div class="tag-title">' + tag + '</div>';
    tags[tag].forEach(function(item) {
      grp.appendChild(createOpCard(item.method, item.path, item.op));
    });
    container.appendChild(grp);
  });
}

function createOpCard(method, path, op) {
  var card = document.createElement('div');
  card.className = 'op-card';
  var summary = document.createElement('div');
  summary.className = 'op-summary';
  summary.innerHTML = '<span class="method-badge method-' + method + '">' + method + '</span>' +
    '<span class="op-path">' + path + '</span>' +
    '<span class="op-desc">' + (op.summary || '') + '</span>';

  var body = document.createElement('div');
  body.className = 'op-body';
  summary.onclick = function() { body.classList.toggle('open'); };

  var formHtml = '';
  var params = op.parameters || [];
  if (params.length > 0) {
    params.forEach(function(p) {
      formHtml += '<div class="form-group"><label>' + p.name + ' (' + p.in + ')' + (p.required ? ' *' : '') + '</label>' +
        '<input type="text" data-in="' + p.in + '" data-name="' + p.name + '" placeholder="' + (p.description || '') + '"></div>';
    });
  }

  if (op.requestBody) {
    formHtml += '<div class="form-group"><label>Request Body (JSON)</label>' +
      '<textarea data-type="body" placeholder="{}"></textarea></div>';
  }

  formHtml += '<button class="btn-exec">Try it out</button>' +
    '<div class="resp-box" style="display:none;"><div class="resp-header"><span>Response</span><span class="resp-code"></span></div>' +
    '<pre class="resp-content"></pre></div>';

  body.innerHTML = formHtml;
  var btn = body.querySelector('.btn-exec');
  btn.onclick = function() { executeRequest(method, path, body); };

  card.appendChild(summary);
  card.appendChild(body);
  return card;
}

function executeRequest(method, pathTemplate, bodyEl) {
  var resolvedPath = pathTemplate;
  var queryParams = [];
  var inputs = bodyEl.querySelectorAll('input');
  inputs.forEach(function(inp) {
    var where = inp.getAttribute('data-in');
    var name = inp.getAttribute('data-name');
    var val = inp.value.trim();
    if (where === 'path') {
      resolvedPath = resolvedPath.replace('{' + name + '}', encodeURIComponent(val));
    } else if (where === 'query' && val) {
      queryParams.push(encodeURIComponent(name) + '=' + encodeURIComponent(val));
    }
  });

  if (queryParams.length > 0) {
    resolvedPath += (resolvedPath.indexOf('?') === -1 ? '?' : '&') + queryParams.join('&');
  }

  var textarea = bodyEl.querySelector('textarea[data-type="body"]');
  var bodyData = null;
  var headers = {};
  if (textarea && textarea.value.trim()) {
    bodyData = textarea.value.trim();
    headers['Content-Type'] = 'application/json';
  }

  var respBox = bodyEl.querySelector('.resp-box');
  var respCode = bodyEl.querySelector('.resp-code');
  var respContent = bodyEl.querySelector('.resp-content');
  respBox.style.display = 'block';
  respCode.textContent = 'Sending...';
  respContent.textContent = '';

  var opts = { method: method.toUpperCase(), headers: headers };
  if (bodyData && method !== 'get' && method !== 'head') {
    opts.body = bodyData;
  }

  fetch(resolvedPath, opts)
    .then(function(r) {
      respCode.textContent = r.status + ' ' + r.statusText;
      respCode.className = 'resp-code' + (r.ok ? '' : ' err');
      return r.text();
    })
    .then(function(text) {
      try {
        var parsed = JSON.parse(text);
        respContent.textContent = JSON.stringify(parsed, null, 2);
      } catch(e) {
        respContent.textContent = text;
      }
    })
    .catch(function(err) {
      respCode.textContent = 'Network Error';
      respCode.className = 'resp-code err';
      respContent.textContent = String(err);
    });
}
</script>
</body>
</html>`

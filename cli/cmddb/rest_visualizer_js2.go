package cmddb

const visualizerScriptJsPart2 = `  function executeSql() {
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
`

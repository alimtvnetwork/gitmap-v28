package cmddb

const visualizerScriptJsPart1 = `(function() {
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
`

package cmddb

const visualizerScriptJsPart1 = `(function() {
  var state = {
    activeDb: '',
    activeTable: '',
    page: 1,
    limit: 50,
    totalPages: 1,
    tables: [],
    columns: [],
    lastSqlRows: [],
    cellStore: {},
    cellCounter: 0,
    blobStore: {},
    blobCounter: 0,
    jsonStore: {},
    jsonCounter: 0,
    protoStore: {},
    protoCounter: 0
  };
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

  function escapeHtml(s) {
    if (s === null || s === undefined) return '';
    return String(s)
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#39;');
  }

  function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }

  function isBinaryString(str) {
    if (typeof str !== 'string' || str.length === 0) return false;
    if (str.indexOf('\uFFFD') !== -1 || str.indexOf('\x00') !== -1) return true;
    var checkLen = Math.min(str.length, 256);
    for (var i = 0; i < checkLen; i++) {
      var code = str.charCodeAt(i);
      if (code < 32 && code !== 9 && code !== 10 && code !== 13) return true;
    }
    return false;
  }

  function renderCellContent(val, colName) {
    if (val === null || val === undefined) {
      return '<span class="cell-null">NULL</span>';
    }
    if (typeof val === 'boolean') {
      var isTrue = val === true;
      return '<span class="badge-bool ' + (isTrue ? 'true' : 'false') + '">[BOOL ' + (isTrue ? 'true' : 'false') + ']</span>';
    }
    if (typeof val === 'object' && (val.__type === 'binary' || val.__blob || val.hex || val.base64Data)) {
      var bId = 'cell_' + (++state.cellCounter);
      var sz = val.byteLength || val.size || (val.hex ? Math.floor(val.hex.length / 2) : 0);
      state.cellStore[bId] = { col: colName, type: 'blob', data: val, size: sz, hex: val.hex, base64: val.base64Data };
      return '<button type="button" class="badge-blob cell-badge" onclick="window.visualizerOpenInspect(\'' + bId + '\', \'hex\')">[BLOB ' + formatBytes(sz) + ']</button>';
    }
    if (typeof val === 'object') {
      var jId = 'cell_' + (++state.cellCounter);
      state.cellStore[jId] = { col: colName, type: 'json', data: val, text: JSON.stringify(val, null, 2) };
      var count = Array.isArray(val) ? val.length + ' items' : Object.keys(val).length + ' keys';
      return '<button type="button" class="badge-json cell-badge" onclick="window.visualizerOpenInspect(\'' + jId + '\', \'json\')">[JSON ' + count + ']</button>';
    }
    if (typeof val === 'string') {
      if (isBinaryString(val)) {
        var pId = 'cell_' + (++state.cellCounter);
        state.cellStore[pId] = { col: colName, type: 'blob', text: val, size: val.length };
        return '<button type="button" class="badge-blob cell-badge" onclick="window.visualizerOpenInspect(\'' + pId + '\', \'hex\')">[BLOB ' + formatBytes(val.length) + ']</button>';
      }
      var trimmed = val.trim();
      if ((trimmed.startsWith('{') && trimmed.endsWith('}')) || (trimmed.startsWith('[') && trimmed.endsWith(']'))) {
        try {
          var parsed = JSON.parse(trimmed);
          if (typeof parsed === 'object' && parsed !== null) {
            var jId2 = 'cell_' + (++state.cellCounter);
            state.cellStore[jId2] = { col: colName, type: 'json', data: parsed, text: JSON.stringify(parsed, null, 2) };
            var jCount = Array.isArray(parsed) ? parsed.length + ' items' : Object.keys(parsed).length + ' keys';
            return '<button type="button" class="badge-json cell-badge" onclick="window.visualizerOpenInspect(\'' + jId2 + '\', \'json\')">[JSON ' + jCount + ']</button>';
          }
        } catch (e) {}
      }
      if (val.indexOf('\n') !== -1 || val.indexOf('\r') !== -1) {
        var lines = val.split(/\r?\n/);
        var tId = 'cell_' + (++state.cellCounter);
        state.cellStore[tId] = { col: colName, type: 'text', text: val };
        var firstLine = lines[0].slice(0, 35);
        return '<button type="button" class="badge-text cell-badge" onclick="window.visualizerOpenInspect(\'' + tId + '\', \'text\')">' + escapeHtml(firstLine) + ' &#8629; [' + lines.length + ' lines]</button>';
      }
      if (val.length > 70) {
        var lId = 'cell_' + (++state.cellCounter);
        state.cellStore[lId] = { col: colName, type: 'text', text: val };
        return '<button type="button" class="cell-truncated cell-badge" onclick="window.visualizerOpenInspect(\'' + lId + '\', \'text\')" title="' + escapeHtml(val) + '">' + escapeHtml(val.slice(0, 68)) + '&hellip;</button>';
      }
    }
    return escapeHtml(String(val));
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
      updatePaginationControls(0, 0);
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
        td.innerHTML = renderCellContent(r[c], c);
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
    updatePaginationControls(rows.length, total);
  }

  function updatePaginationControls(rowCount, total) {
    var totalPages = Math.max(1, Math.ceil(total / state.limit));
    state.totalPages = totalPages;
    document.getElementById('page-info').textContent = 'Showing ' + rowCount + ' of ' + total + ' records (Page ' + state.page + ' of ' + totalPages + ')';
    var pageNum = document.getElementById('page-num');
    if (pageNum) pageNum.textContent = state.page;
    var pageInput = document.getElementById('page-input');
    if (pageInput) {
      pageInput.value = state.page;
      pageInput.max = totalPages;
    }
    var pageMaxInfo = document.getElementById('page-max-info');
    if (pageMaxInfo) pageMaxInfo.textContent = '/ ' + totalPages;

    var btnPrev = document.getElementById('btn-prev');
    var isPrevDisabled = state.page <= 1;
    btnPrev.disabled = isPrevDisabled;
    btnPrev.style.opacity = isPrevDisabled ? '0.4' : '1';
    btnPrev.style.pointerEvents = isPrevDisabled ? 'none' : 'auto';

    var btnNext = document.getElementById('btn-next');
    var isNextDisabled = state.page >= totalPages || rowCount < state.limit;
    btnNext.disabled = isNextDisabled;
    btnNext.style.opacity = isNextDisabled ? '0.4' : '1';
    btnNext.style.pointerEvents = isNextDisabled ? 'none' : 'auto';
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
      if (state.page > 1) {
        state.page--;
        loadGridData();
      }
    };
    document.getElementById('btn-next').onclick = function() {
      if (state.page < state.totalPages) {
        state.page++;
        loadGridData();
      }
    };
    var pageInput = document.getElementById('page-input');
    if (pageInput) {
      pageInput.onchange = function(e) {
        var p = parseInt(e.target.value, 10);
        if (isNaN(p) || p < 1) p = 1;
        if (state.totalPages && p > state.totalPages) p = state.totalPages;
        state.page = p;
        e.target.value = p;
        loadGridData();
      };
    }
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

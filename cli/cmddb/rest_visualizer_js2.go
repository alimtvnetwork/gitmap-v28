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
        td.innerHTML = renderCellContent(r[c], c);
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

  function stringToHex(str) {
    if (!str) return '';
    var out = '';
    for (var i = 0; i < str.length; i++) {
      var h = (str.charCodeAt(i) & 0xff).toString(16);
      if (h.length < 2) h = '0' + h;
      out += h;
    }
    return out;
  }

  function hexToAscii(hexStr) {
    if (!hexStr) return '';
    var out = '';
    for (var i = 0; i < hexStr.length; i += 2) {
      var code = parseInt(hexStr.slice(i, i + 2), 16);
      out += (code >= 32 && code <= 126) ? String.fromCharCode(code) : '.';
    }
    return out;
  }

  function renderHexLinesInto(hexStr, tbody) {
    tbody.innerHTML = '';
    if (!hexStr) {
      tbody.innerHTML = '<tr><td colspan="3" style="color:#64748b; padding:12px;">(empty data)</td></tr>';
      return;
    }
    var len = hexStr.length;
    var offset = 0;
    while (offset < len) {
      var chunkHex = hexStr.slice(offset, offset + 32);
      var bytes = [];
      var ascii = '';
      for (var i = 0; i < chunkHex.length; i += 2) {
        var b = chunkHex.slice(i, i + 2);
        bytes.push(b.toUpperCase());
        var code = parseInt(b, 16);
        ascii += (code >= 32 && code <= 126) ? String.fromCharCode(code) : '.';
      }
      var byteStr1 = bytes.slice(0, 8).join(' ');
      var byteStr2 = bytes.slice(8).join(' ');
      var fullByteStr = byteStr1 + (bytes.length > 8 ? '  ' + byteStr2 : '');
      var offStr = ('00000000' + Math.floor(offset / 2).toString(16).toUpperCase()).slice(-8);

      var tr = document.createElement('tr');
      tr.innerHTML = '<td class="hex-offset-col">' + offStr + '</td>' +
        '<td class="hex-bytes-col">' + fullByteStr + '</td>' +
        '<td class="hex-ascii-col">|' + ascii + '|</td>';
      tbody.appendChild(tr);
      offset += 32;
    }
  }

  function renderJsonTree(obj, container) {
    container.innerHTML = '';
    container.appendChild(createJsonNode(obj, 'root', true));
  }

  function createJsonNode(val, key, isRoot) {
    var wrap = document.createElement('div');
    wrap.style.marginLeft = isRoot ? '0' : '16px';
    if (val === null) {
      wrap.innerHTML = (key !== 'root' ? '<span class="json-key">\"' + key + '\": </span>' : '') + '<span class="json-null">null</span>';
      return wrap;
    }
    var t = typeof val;
    if (t === 'string') {
      wrap.innerHTML = (key !== 'root' ? '<span class="json-key">\"' + key + '\": </span>' : '') + '<span class="json-string">\"' + escapeHtml(val) + '\"</span>';
      return wrap;
    }
    if (t === 'number') {
      wrap.innerHTML = (key !== 'root' ? '<span class="json-key">\"' + key + '\": </span>' : '') + '<span class="json-number">' + val + '</span>';
      return wrap;
    }
    if (t === 'boolean') {
      wrap.innerHTML = (key !== 'root' ? '<span class="json-key">\"' + key + '\": </span>' : '') + '<span class="json-boolean">' + (val ? 'true' : 'false') + '</span>';
      return wrap;
    }
    if (t === 'object') {
      var isArr = Array.isArray(val);
      var openChar = isArr ? '[' : '{';
      var closeChar = isArr ? ']' : '}';
      var keys = Object.keys(val);
      var header = document.createElement('div');
      var toggle = document.createElement('span');
      toggle.className = 'json-toggle';
      toggle.textContent = '\u25BC ';
      header.appendChild(toggle);
      if (key !== 'root') {
        var kSpan = document.createElement('span');
        kSpan.className = 'json-key';
        kSpan.textContent = '\"' + key + '\": ';
        header.appendChild(kSpan);
      }
      var oSpan = document.createElement('span');
      oSpan.textContent = openChar + ' (' + keys.length + ')';
      header.appendChild(oSpan);
      wrap.appendChild(header);

      var childBox = document.createElement('div');
      keys.forEach(function(k) {
        childBox.appendChild(createJsonNode(val[k], k, false));
      });
      wrap.appendChild(childBox);

      var footer = document.createElement('div');
      footer.textContent = closeChar;
      wrap.appendChild(footer);

      toggle.onclick = function() {
        var isHidden = childBox.style.display === 'none';
        childBox.style.display = isHidden ? 'block' : 'none';
        footer.style.display = isHidden ? 'block' : 'none';
        toggle.textContent = isHidden ? '\u25BC ' : '\u25B6 ';
        oSpan.textContent = isHidden ? openChar + ' (' + keys.length + ')' : openChar + ' ... ' + closeChar;
      };
      return wrap;
    }
    wrap.textContent = String(val);
    return wrap;
  }

  window.visualizerOpenInspect = function(cellId, initialTab) {
    var entry = state.cellStore[cellId];
    if (!entry) return;
    state.activeInspectEntry = entry;

    document.getElementById('inspect-modal-col').textContent = entry.col || 'Cell';
    var badgeLabel = entry.type === 'blob' ? '[BLOB ' + formatBytes(entry.size || 0) + ']' :
                     entry.type === 'json' ? '[JSON]' : '[TEXT]';
    document.getElementById('inspect-modal-badge').textContent = badgeLabel;

    // 1. Prepare Hex Dump
    var hexStr = '';
    if (entry.hex) {
      hexStr = entry.hex;
    } else if (entry.data && entry.data.hex) {
      hexStr = entry.data.hex;
    } else if (entry.text) {
      hexStr = stringToHex(entry.text);
    } else if (entry.data && typeof entry.data === 'object') {
      hexStr = stringToHex(JSON.stringify(entry.data));
    }
    state.currentHexStr = hexStr;
    renderHexLinesInto(hexStr, document.getElementById('inspect-hex-body'));

    // 2. Prepare JSON view
    var jsonContainer = document.getElementById('inspect-json-body');
    var dataObj = entry.data || null;
    if (!dataObj && entry.type === 'json' && entry.text) {
      try { dataObj = JSON.parse(entry.text); } catch(e) {}
    }
    if (dataObj && typeof dataObj === 'object') {
      state.currentJsonData = dataObj;
      renderJsonTree(dataObj, jsonContainer);
    } else {
      state.currentJsonData = null;
      jsonContainer.innerHTML = '<div style="color:#64748b; padding:12px; font-style:italic;">Data is not formatted JSON. View Hex Dump or Decoded String tabs.</div>';
    }

    // 3. Prepare Decoded String view
    var rawText = entry.text || (dataObj ? JSON.stringify(dataObj, null, 2) : hexToAscii(hexStr));
    state.currentRawText = rawText;
    document.getElementById('inspect-text-body').textContent = rawText;

    // Switch to active tab
    var targetPaneId = initialTab === 'json' ? 'inspect-tab-json' :
                       initialTab === 'text' ? 'inspect-tab-text' : 'inspect-tab-hex';
    switchInspectTab(targetPaneId);

    document.getElementById('inspector-modal').style.display = 'flex';
  };

  function switchInspectTab(targetPaneId) {
    document.querySelectorAll('.inspect-tab-btn').forEach(function(b) {
      b.classList.remove('active');
      if (b.getAttribute('data-inspect-tab') === targetPaneId) b.classList.add('active');
    });
    document.querySelectorAll('.inspect-pane').forEach(function(p) {
      p.classList.remove('active');
      if (p.id === targetPaneId) p.classList.add('active');
    });
  }

  function closeInspectorModal() {
    document.getElementById('inspector-modal').style.display = 'none';
  }

  window.visualizerOpenHex = function(bId) {
    window.visualizerOpenInspect(bId, 'hex');
  };

  window.visualizerOpenProto = function(pId) {
    window.visualizerOpenInspect(pId, 'hex');
  };

  window.visualizerOpenJson = function(jId) {
    window.visualizerOpenInspect(jId, 'json');
  };

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

  function initVisualizerModals() {
    var btnClose = document.getElementById('btn-inspect-close');
    if (btnClose) btnClose.onclick = closeInspectorModal;
    var btnDone = document.getElementById('btn-inspect-done');
    if (btnDone) btnDone.onclick = closeInspectorModal;

    document.querySelectorAll('.inspect-tab-btn').forEach(function(btn) {
      btn.onclick = function() {
        var targetPane = btn.getAttribute('data-inspect-tab');
        switchInspectTab(targetPane);
      };
    });

    var btnCopyHex = document.getElementById('btn-inspect-copy-hex');
    if (btnCopyHex) {
      btnCopyHex.onclick = function() {
        if (!state.currentHexStr) return;
        navigator.clipboard.writeText(state.currentHexStr).then(function() {
          var st = document.getElementById('inspect-copy-status');
          st.textContent = 'Hex copied!';
          setTimeout(function() { st.textContent = ''; }, 2000);
        });
      };
    }

    var btnCopyJson = document.getElementById('btn-inspect-copy-json');
    if (btnCopyJson) {
      btnCopyJson.onclick = function() {
        var jsonText = state.currentJsonData ? JSON.stringify(state.currentJsonData, null, 2) : (state.currentRawText || '');
        navigator.clipboard.writeText(jsonText).then(function() {
          var st = document.getElementById('inspect-copy-status');
          st.textContent = 'JSON copied!';
          setTimeout(function() { st.textContent = ''; }, 2000);
        });
      };
    }

    var btnCopyText = document.getElementById('btn-inspect-copy-text');
    if (btnCopyText) {
      btnCopyText.onclick = function() {
        var rawText = state.currentRawText || '';
        navigator.clipboard.writeText(rawText).then(function() {
          var st = document.getElementById('inspect-copy-status');
          st.textContent = 'String copied!';
          setTimeout(function() { st.textContent = ''; }, 2000);
        });
      };
    }
  }

  initVisualizerModals();
})();
`

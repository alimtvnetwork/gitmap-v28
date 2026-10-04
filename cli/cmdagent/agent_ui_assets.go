package cmdagent

import (
	"net/http"
)

// ServeDashboardHTML renders the embedded single-page visualizer dashboard.
func ServeDashboardHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitMap Agent Fleet Visualizer</title>
  <style>
    :root {
      --bg-body: #090d16;
      --bg-header: #0d1322;
      --bg-panel: #111827;
      --bg-card: #182235;
      --bg-card-hover: #1f2c44;
      --border-subtle: #202e48;
      --border-strong: #2f4266;
      --text-main: #f8fafc;
      --text-muted: #94a3b8;
      --text-faint: #64748b;
      --color-cyan: #38bdf8;
      --color-green: #10b981;
      --color-amber: #f59e0b;
      --color-red: #ef4444;
      --color-purple: #a855f7;
      --color-indigo: #6366f1;
      --font-mono: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', Consolas, monospace;
      --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-body);
      color: var(--text-main);
      font-family: var(--font-sans);
      font-size: 13px;
      line-height: 1.5;
      overflow-x: hidden;
    }
    header {
      background-color: var(--bg-header);
      border-bottom: 1px solid var(--border-subtle);
      padding: 12px 24px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      position: sticky;
      top: 0;
      z-index: 40;
    }
    .brand-group { display: flex; align-items: center; gap: 12px; }
    .brand-icon {
      width: 28px;
      height: 28px;
      background: linear-gradient(135deg, var(--color-cyan), var(--color-indigo));
      border-radius: 6px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 800;
      font-size: 14px;
      color: #fff;
    }
    .brand-title { font-size: 16px; font-weight: 700; letter-spacing: -0.3px; }
    .header-actions { display: flex; align-items: center; gap: 12px; }
    .status-pill {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 4px 10px;
      border-radius: 9999px;
      font-size: 11px;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .status-pill.RUNNING { background: rgba(16, 185, 129, 0.15); color: var(--color-green); border: 1px solid rgba(16, 185, 129, 0.4); }
    .status-pill.ACTIVE { background: rgba(56, 189, 248, 0.15); color: var(--color-cyan); border: 1px solid rgba(56, 189, 248, 0.4); }
    .status-pill.CRASH_DETECTED { background: rgba(239, 68, 68, 0.2); color: var(--color-red); border: 1px solid rgba(239, 68, 68, 0.6); animation: pulseRed 1.8s infinite; }
    .status-pill.COMPLETED { background: rgba(16, 185, 129, 0.15); color: var(--color-green); border: 1px solid rgba(16, 185, 129, 0.4); }
    .status-pill.IDLE { background: rgba(100, 116, 139, 0.15); color: var(--text-faint); border: 1px solid rgba(100, 116, 139, 0.3); }
    @keyframes pulseRed { 0%, 100% { opacity: 1; } 50% { opacity: 0.5; } }
    .btn {
      background: var(--bg-card);
      border: 1px solid var(--border-strong);
      color: var(--text-main);
      padding: 6px 12px;
      border-radius: 6px;
      cursor: pointer;
      font-size: 12px;
      font-weight: 600;
      display: inline-flex;
      align-items: center;
      gap: 6px;
      transition: all 0.15s ease;
    }
    .btn:hover { background: var(--bg-card-hover); border-color: var(--color-cyan); }
    .btn-danger { background: rgba(239, 68, 68, 0.15); border-color: rgba(239, 68, 68, 0.4); color: var(--color-red); }
    .btn-danger:hover { background: rgba(239, 68, 68, 0.25); border-color: var(--color-red); }
    .container { padding: 20px 24px; max-width: 1600px; margin: 0 auto; }
    .cards-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 16px;
      margin-bottom: 20px;
    }
    .card {
      background: var(--bg-card);
      border: 1px solid var(--border-subtle);
      border-radius: 8px;
      padding: 16px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    }
    .card-title {
      font-size: 11px;
      text-transform: uppercase;
      letter-spacing: 0.8px;
      color: var(--text-muted);
      margin-bottom: 8px;
      font-weight: 600;
    }
    .card-value {
      font-size: 20px;
      font-weight: 700;
      color: var(--text-main);
      margin-bottom: 6px;
      word-break: break-all;
    }
    .card-subtext { font-size: 12px; color: var(--text-muted); }
    .progress-bar-bg {
      width: 100%;
      height: 6px;
      background: var(--border-subtle);
      border-radius: 3px;
      overflow: hidden;
      margin-top: 8px;
    }
    .progress-bar-fill {
      height: 100%;
      background: linear-gradient(90deg, var(--color-cyan), var(--color-green));
      border-radius: 3px;
      transition: width 0.3s ease;
    }
    .main-grid {
      display: grid;
      grid-template-columns: 360px 1fr;
      gap: 20px;
    }
    @media (max-width: 1024px) {
      .main-grid { grid-template-columns: 1fr; }
    }
    .panel {
      background: var(--bg-panel);
      border: 1px solid var(--border-subtle);
      border-radius: 8px;
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }
    .panel-header {
      background: var(--bg-card);
      border-bottom: 1px solid var(--border-subtle);
      padding: 12px 16px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 700;
      font-size: 13px;
    }
    .panel-body { padding: 16px; overflow-y: auto; max-height: calc(100vh - 280px); }
    .tree-node { margin-bottom: 8px; }
    .tree-item {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 6px 10px;
      border-radius: 6px;
      cursor: pointer;
      user-select: none;
      transition: background 0.15s ease;
      font-family: var(--font-mono);
      font-size: 12px;
    }
    .tree-item:hover { background: var(--bg-card); }
    .tree-item.selected { background: rgba(56, 189, 248, 0.15); border: 1px solid rgba(56, 189, 248, 0.3); }
    .tree-children { margin-left: 20px; border-left: 1px dashed var(--border-subtle); padding-left: 10px; margin-top: 4px; }
    .badge {
      display: inline-block;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 10px;
      font-weight: 700;
      text-transform: uppercase;
      font-family: var(--font-mono);
    }
    .badge-DONE { background: rgba(16, 185, 129, 0.2); color: var(--color-green); }
    .badge-IN_PROGRESS { background: rgba(245, 158, 11, 0.2); color: var(--color-amber); }
    .badge-PENDING { background: rgba(100, 116, 139, 0.2); color: var(--text-faint); }
    .badge-FAILED { background: rgba(239, 68, 68, 0.2); color: var(--color-red); }
    .action-badge {
      display: inline-block;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 10px;
      font-weight: 700;
      font-family: var(--font-mono);
    }
    .act-WRITE { background: rgba(56, 189, 248, 0.2); color: var(--color-cyan); }
    .act-READ { background: rgba(168, 85, 247, 0.2); color: var(--color-purple); }
    .act-EXEC { background: rgba(99, 102, 241, 0.2); color: var(--color-indigo); }
    .act-COMPLETE { background: rgba(16, 185, 129, 0.2); color: var(--color-green); }
    .act-START, .act-CLAIM { background: rgba(245, 158, 11, 0.2); color: var(--color-amber); }
    .act-FAIL, .act-CRASH { background: rgba(239, 68, 68, 0.25); color: var(--color-red); }
    .act-SEARCH, .act-CHECK, .act-LINT { background: rgba(100, 116, 139, 0.25); color: var(--text-muted); }
    .filter-bar {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      padding: 12px 16px;
      background: var(--bg-card);
      border-bottom: 1px solid var(--border-subtle);
    }
    .search-input, .select-input {
      background: var(--bg-panel);
      border: 1px solid var(--border-strong);
      color: var(--text-main);
      padding: 6px 10px;
      border-radius: 6px;
      font-size: 12px;
      outline: none;
    }
    .search-input { flex: 1; min-width: 180px; }
    .search-input:focus, .select-input:focus { border-color: var(--color-cyan); }
    .logs-table {
      width: 100%;
      border-collapse: collapse;
      font-size: 12px;
    }
    .logs-table th {
      text-align: left;
      padding: 10px 12px;
      background: rgba(13, 19, 34, 0.6);
      color: var(--text-muted);
      font-weight: 600;
      border-bottom: 1px solid var(--border-subtle);
    }
    .logs-table td {
      padding: 8px 12px;
      border-bottom: 1px solid var(--border-subtle);
      font-family: var(--font-mono);
      vertical-align: middle;
    }
    .logs-table tr:hover td { background: var(--bg-card-hover); cursor: pointer; }
    .modal-backdrop {
      position: fixed;
      top: 0; left: 0; width: 100%; height: 100%;
      background: rgba(0, 0, 0, 0.7);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 100;
      visibility: hidden;
      opacity: 0;
      transition: all 0.2s ease;
    }
    .modal-backdrop.open { visibility: visible; opacity: 1; }
    .modal-window {
      background: var(--bg-panel);
      border: 1px solid var(--border-strong);
      border-radius: 10px;
      width: 90%;
      max-width: 720px;
      max-height: 85vh;
      overflow-y: auto;
      box-shadow: 0 16px 40px rgba(0, 0, 0, 0.6);
    }
    .modal-header-danger {
      background: linear-gradient(135deg, rgba(239, 68, 68, 0.2), rgba(239, 68, 68, 0.05));
      border-bottom: 1px solid rgba(239, 68, 68, 0.4);
      padding: 16px 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .modal-title-danger { font-size: 15px; font-weight: 700; color: var(--color-red); display: flex; align-items: center; gap: 8px; }
    .modal-body { padding: 20px; }
    .autopsy-card {
      background: var(--bg-card);
      border: 1px solid rgba(239, 68, 68, 0.3);
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 14px;
    }
    .autopsy-field { display: flex; margin-bottom: 6px; font-size: 12px; }
    .autopsy-label { width: 140px; color: var(--text-muted); font-weight: 600; font-family: var(--font-sans); }
    .autopsy-value { flex: 1; font-family: var(--font-mono); color: var(--text-main); word-break: break-all; }
    .diagnosis-box {
      margin-top: 10px;
      padding: 12px;
      background: rgba(239, 68, 68, 0.1);
      border-left: 3px solid var(--color-red);
      border-radius: 4px;
      color: #fca5a5;
      font-size: 12px;
      line-height: 1.4;
    }
    .modal-footer {
      padding: 12px 20px;
      background: var(--bg-card);
      border-top: 1px solid var(--border-subtle);
      display: flex;
      justify-content: flex-end;
      gap: 10px;
    }
  </style>
</head>
<body>
  <header>
    <div class="brand-group">
      <div class="brand-icon">GM</div>
      <div>
        <div class="brand-title">GitMap AI Agent Visualizer</div>
        <div style="font-size: 11px; color: var(--text-muted);">Multi-Tiered Split-DB Task Orchestrator</div>
      </div>
    </div>
    <div class="header-actions">
      <span id="headerStatusPill" class="status-pill IDLE">IDLE</span>
      <label style="display: flex; align-items: center; gap: 6px; font-size: 12px; cursor: pointer; color: var(--text-muted);">
        <input type="checkbox" id="autoRefreshCheck" checked> Auto-refresh (2.5s)
      </label>
      <button class="btn" id="refreshBtn" onclick="manualRefresh()">
        <span>↻</span> Refresh
      </button>
      <button class="btn btn-danger" onclick="confirmClearTask()">
        <span>🗑</span> Clear Task
      </button>
    </div>
  </header>

  <div class="container">
    <div class="cards-grid">
      <div class="card">
        <div class="card-title">Active Parent Task</div>
        <div class="card-value" id="cardTaskName" style="font-size: 15px;">Scanning...</div>
        <div class="card-subtext" id="cardTaskSlug">-</div>
      </div>
      <div class="card">
        <div class="card-title">Steps Budget</div>
        <div class="card-value" id="cardStepsCount">0 / 300</div>
        <div class="card-subtext" id="cardStepsPercent">0% Completed</div>
        <div class="progress-bar-bg">
          <div class="progress-bar-fill" id="cardStepsBar" style="width: 0%;"></div>
        </div>
      </div>
      <div class="card">
        <div class="card-title">In-Flight Agent Fleet</div>
        <div class="card-value" id="cardFleetCount">0 Active</div>
        <div class="card-subtext" id="cardFleetRoles">No agents claimed</div>
      </div>
      <div class="card" id="cardCrashContainer">
        <div class="card-title">Crashes Detected</div>
        <div class="card-value" id="cardCrashCount" style="color: var(--color-green);">0 Clean</div>
        <div class="card-subtext" id="cardCrashSubtext">Fleet status healthy</div>
      </div>
    </div>

    <div class="main-grid">
      <div class="panel">
        <div class="panel-header">
          <span>Split-DB Hierarchy</span>
          <span style="font-size: 11px; color: var(--text-muted); font-family: var(--font-mono);" id="dbHierarchyTag">Tier 1..3</span>
        </div>
        <div class="panel-body" id="treeViewContainer">
          <div style="color: var(--text-muted); padding: 10px;">Loading database tree...</div>
        </div>
      </div>

      <div class="panel">
        <div class="panel-header">
          <span>Audit Trail Explorer</span>
          <span style="font-size: 11px; color: var(--text-muted);" id="logsCountTag">0 records</span>
        </div>
        <div class="filter-bar">
          <input type="text" id="searchInput" class="search-input" placeholder="Search target file, details, or command..." oninput="handleFilterChange()">
          <select id="agentFilter" class="select-input" onchange="handleFilterChange()">
            <option value="ALL">All Agents</option>
          </select>
          <select id="actionFilter" class="select-input" onchange="handleFilterChange()">
            <option value="ALL">All Actions</option>
            <option value="SEARCH">SEARCH</option>
            <option value="READ">READ</option>
            <option value="WRITE">WRITE</option>
            <option value="EXEC">EXEC</option>
            <option value="LINT">LINT</option>
            <option value="CHECK">CHECK</option>
            <option value="CLAIM">CLAIM</option>
            <option value="START">START</option>
            <option value="COMPLETE">COMPLETE</option>
            <option value="FAIL">FAIL</option>
            <option value="CRASH">CRASH</option>
          </select>
          <select id="subtaskFilter" class="select-input" onchange="handleFilterChange()">
            <option value="ALL">All Subtasks</option>
          </select>
        </div>
        <div style="overflow-x: auto; flex: 1;">
          <table class="logs-table">
            <thead>
              <tr>
                <th style="width: 140px;">Time</th>
                <th style="width: 90px;">Agent</th>
                <th style="width: 70px;">Subtask</th>
                <th style="width: 80px;">Action</th>
                <th>Target File / Command Details</th>
                <th style="width: 80px;">Status</th>
              </tr>
            </thead>
            <tbody id="logsTableBody">
              <tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 20px;">Scanning agent action telemetry...</td></tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <div class="modal-backdrop" id="crashModal">
    <div class="modal-window">
      <div class="modal-header-danger">
        <div class="modal-title-danger">
          <span>⚠</span> CRASH AUTOPSY & FORENSIC DIAGNOSTICS
        </div>
        <button class="btn" onclick="closeCrashModal()">✕</button>
      </div>
      <div class="modal-body" id="crashModalBody">
        <div style="color: var(--text-muted);">No crash autopsies detected.</div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-danger" onclick="triggerClearCurrentTask()">Clear Task Run</button>
        <button class="btn" onclick="closeCrashModal()">Close</button>
      </div>
    </div>
  </div>

  <div class="modal-backdrop" id="detailModal">
    <div class="modal-window">
      <div class="panel-header">
        <span>Action Log Forensic Record</span>
        <button class="btn" onclick="closeDetailModal()">✕</button>
      </div>
      <div class="modal-body" id="detailModalBody"></div>
      <div class="modal-footer">
        <button class="btn" onclick="closeDetailModal()">Close</button>
      </div>
    </div>
  </div>

  <script>
    var appState = {
      selectedTaskId: '',
      autoRefresh: true,
      lastStatusData: null,
      lastTasksData: null,
      lastCrashesData: null,
      timer: null
    };

    function init() {
      fetchFullDashboard();
      appState.timer = setInterval(function() {
        var isChecked = document.getElementById('autoRefreshCheck').checked;
        if (isChecked) {
          fetchFullDashboard();
        }
      }, 2500);
    }

    function manualRefresh() {
      fetchFullDashboard();
    }

    function fetchFullDashboard() {
      var taskQuery = appState.selectedTaskId ? '?taskId=' + encodeURIComponent(appState.selectedTaskId) : '';
      Promise.all([
        fetch('/api/agent/status' + taskQuery).then(function(r) { return r.json(); }).catch(function() { return null; }),
        fetch('/api/agent/tasks' + taskQuery).then(function(r) { return r.json(); }).catch(function() { return null; }),
        fetch('/api/agent/crashes' + taskQuery).then(function(r) { return r.json(); }).catch(function() { return null; }),
        fetchActionLogs()
      ]).then(function(results) {
        var statusData = results[0];
        var tasksData = results[1];
        var crashesData = results[2];
        if (statusData) renderStatusCards(statusData);
        if (tasksData) renderTasksAndTree(tasksData);
        if (crashesData) renderCrashesCard(crashesData);
      });
    }

    function fetchActionLogs() {
      var searchVal = document.getElementById('searchInput').value.trim();
      var agentVal = document.getElementById('agentFilter').value;
      var actionVal = document.getElementById('actionFilter').value;
      var subVal = document.getElementById('subtaskFilter').value;
      var params = [];
      if (appState.selectedTaskId) params.push('taskId=' + encodeURIComponent(appState.selectedTaskId));
      if (searchVal) params.push('q=' + encodeURIComponent(searchVal));
      if (agentVal && agentVal !== 'ALL') params.push('agent=' + encodeURIComponent(agentVal));
      if (actionVal && actionVal !== 'ALL') params.push('actionType=' + encodeURIComponent(actionVal));
      if (subVal && subVal !== 'ALL') params.push('subtask=' + encodeURIComponent(subVal));
      params.push('limit=120');

      var url = '/api/agent/logs?' + params.join('&');
      return fetch(url).then(function(r) { return r.json(); }).then(function(data) {
        renderLogsTable(data.logs || []);
      }).catch(function() {
        renderLogsTable([]);
      });
    }

    function renderStatusCards(s) {
      appState.lastStatusData = s;
      var pill = document.getElementById('headerStatusPill');
      pill.textContent = s.status || 'IDLE';
      pill.className = 'status-pill ' + (s.status || 'IDLE');

      var task = s.activeParentTask || {};
      var taskName = task.taskName || task.taskSlug || 'No active task';
      document.getElementById('cardTaskName').textContent = taskName;
      document.getElementById('cardTaskSlug').textContent = task.runDirectory || s.tempDirectory || '-';

      var totalBudget = task.totalStepsBudget || 300;
      var completedSteps = task.completedSteps || 0;
      var pct = Math.min(100, Math.round((completedSteps / totalBudget) * 100));
      document.getElementById('cardStepsCount').textContent = completedSteps + ' / ' + totalBudget;
      document.getElementById('cardStepsPercent').textContent = pct + '% Completed';
      document.getElementById('cardStepsBar').style.width = pct + '%';

      var metrics = s.globalMetrics || {};
      var agentCount = metrics.total_agents || '2';
      document.getElementById('cardFleetCount').textContent = agentCount + ' Total Agents';
      document.getElementById('cardFleetRoles').textContent = (s.status === 'RUNNING' ? 'Fleet actively executing' : 'Fleet idle / claimed');
    }

    function renderCrashesCard(c) {
      appState.lastCrashesData = c;
      var hasCrashes = c.hasCrashesDetected || (c.crashedAgents && c.crashedAgents.length > 0);
      var crashCountEl = document.getElementById('cardCrashCount');
      var crashSubtextEl = document.getElementById('cardCrashSubtext');
      var crashContainer = document.getElementById('cardCrashContainer');

      if (hasCrashes) {
        var count = c.crashedAgents ? c.crashedAgents.length : 1;
        crashCountEl.textContent = count + ' Crash(es)';
        crashCountEl.style.color = 'var(--color-red)';
        crashSubtextEl.innerHTML = '<span style="color: var(--color-red); font-weight:700; cursor:pointer;" onclick="openCrashModal()">⚠ Click to View Autopsy</span>';
        crashContainer.style.borderColor = 'rgba(239, 68, 68, 0.5)';
      } else {
        crashCountEl.textContent = '0 Clean';
        crashCountEl.style.color = 'var(--color-green)';
        crashSubtextEl.textContent = 'Fleet status healthy';
        crashContainer.style.borderColor = 'var(--border-subtle)';
      }
    }

    function renderTasksAndTree(t) {
      appState.lastTasksData = t;
      var parentTasks = t.parentTasks || [];
      var activeTask = t.activeTask || {};
      var subtasks = t.subtasks || [];
      var agents = t.agents || [];

      if (!appState.selectedTaskId && activeTask.taskSlug) {
        appState.selectedTaskId = activeTask.taskSlug;
      }

      updateAgentAndSubtaskDropdowns(subtasks, agents);
      var html = '<div class="tree-node">';
      html += '<div class="tree-item"><span style="color:var(--color-cyan);">🗄</span> <strong>ai_agents.db</strong> <span style="font-size:10px; color:var(--text-faint);">(Master Registry)</span></div>';
      html += '<div class="tree-children">';

      if (parentTasks.length === 0 && activeTask.taskSlug) {
        parentTasks = [activeTask];
      }

      for (var i = 0; i < parentTasks.length; i++) {
        var p = parentTasks[i];
        var isSelected = (p.taskSlug === appState.selectedTaskId || p.parentTaskId === appState.selectedTaskId);
        html += '<div class="tree-node">';
        html += '<div class="tree-item ' + (isSelected ? 'selected' : '') + '" onclick="selectTask(\'' + (p.taskSlug || p.parentTaskId) + '\')">';
        html += '<span>📁</span> <strong>' + (p.taskSlug || p.taskName || 'Task') + '</strong>';
        html += '<span class="badge badge-' + (p.status || 'ACTIVE') + '" style="margin-left:auto;">' + (p.status || 'ACTIVE') + '</span>';
        html += '</div>';

        if (isSelected) {
          html += '<div class="tree-children">';
          html += '<div style="font-size:11px; color:var(--text-muted); margin:4px 0 2px 4px; font-weight:600;">📑 Subtasks (' + subtasks.length + ')</div>';
          for (var s = 0; s < subtasks.length; s++) {
            var sub = subtasks[s];
            html += '<div class="tree-item" style="padding:4px 6px;">';
            html += '<span class="badge badge-' + sub.status + '">' + sub.status + '</span>';
            html += '<span style="font-size:11px;">' + sub.taskCode + ': ' + sub.title.substring(0, 24) + '...</span>';
            html += '</div>';
          }

          html += '<div style="font-size:11px; color:var(--text-muted); margin:6px 0 2px 4px; font-weight:600;">🗂 Agent Split DBs (' + agents.length + ')</div>';
          for (var a = 0; a < agents.length; a++) {
            var ag = agents[a];
            html += '<div class="tree-item" style="padding:4px 6px; color:var(--color-purple);">';
            html += '<span>⚡</span> <span>' + ag.agentRole + '</span>';
            html += '<span style="font-size:10px; color:var(--text-faint); margin-left:auto;">' + (ag.agentSlug || 'db') + '.db</span>';
            html += '</div>';
          }
          html += '</div>';
        }
        html += '</div>';
      }
      html += '</div></div>';
      document.getElementById('treeViewContainer').innerHTML = html;
    }

    function selectTask(slug) {
      appState.selectedTaskId = slug;
      fetchFullDashboard();
    }

    function updateAgentAndSubtaskDropdowns(subtasks, agents) {
      var agentSelect = document.getElementById('agentFilter');
      var subSelect = document.getElementById('subtaskFilter');
      var currAgent = agentSelect.value;
      var currSub = subSelect.value;

      var agentHtml = '<option value="ALL">All Agents</option>';
      var seenRoles = {};
      for (var a = 0; a < agents.length; a++) {
        var r = agents[a].agentRole;
        if (!seenRoles[r]) {
          agentHtml += '<option value="' + r + '">' + r + '</option>';
          seenRoles[r] = true;
        }
      }
      agentSelect.innerHTML = agentHtml;
      agentSelect.value = currAgent || 'ALL';

      var subHtml = '<option value="ALL">All Subtasks</option>';
      for (var s = 0; s < subtasks.length; s++) {
        var code = subtasks[s].taskCode || subtasks[s].subtaskId;
        subHtml += '<option value="' + subtasks[s].subtaskId + '">' + code + ': ' + subtasks[s].title.substring(0, 20) + '</option>';
      }
      subSelect.innerHTML = subHtml;
      subSelect.value = currSub || 'ALL';
    }

    function renderLogsTable(logs) {
      var tbody = document.getElementById('logsTableBody');
      document.getElementById('logsCountTag').textContent = logs.length + ' records';
      if (!logs || logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; color: var(--text-muted); padding: 24px;">No action logs matching filter criteria</td></tr>';
        return;
      }

      var html = '';
      for (var i = 0; i < logs.length; i++) {
        var log = logs[i];
        var timeStr = log.createdAt ? log.createdAt.substring(11, 19) : '--:--:--';
        var actionClass = 'act-' + (log.actionType || 'EXEC');
        var targetInfo = log.targetFile ? log.targetFile : (log.actionDetails || '-');
        html += '<tr onclick="openDetailModal(' + i + ')">';
        html += '<td>' + timeStr + '</td>';
        html += '<td><strong>' + (log.agentRole || 'Worker') + '</strong></td>';
        html += '<td>' + (log.subtaskId || '-') + '</td>';
        html += '<td><span class="action-badge ' + actionClass + '">' + (log.actionType || 'EXEC') + '</span></td>';
        html += '<td style="color:var(--text-main); word-break:break-all;">' + escapeHtml(targetInfo) + '</td>';
        html += '<td><span class="badge badge-' + (log.status || 'DONE') + '">' + (log.status || 'DONE') + '</span></td>';
        html += '</tr>';
      }
      tbody.innerHTML = html;
      window._currentLogs = logs;
    }

    function handleFilterChange() {
      fetchActionLogs();
    }

    function openCrashModal() {
      var modal = document.getElementById('crashModal');
      var body = document.getElementById('crashModalBody');
      var c = appState.lastCrashesData;
      if (!c || !c.crashedAgents || c.crashedAgents.length === 0) {
        body.innerHTML = '<div style="color:var(--text-muted); padding:20px;">No crashed agents or incomplete subtasks detected in current run.</div>';
      } else {
        var html = '';
        for (var i = 0; i < c.crashedAgents.length; i++) {
          var a = c.crashedAgents[i];
          html += '<div class="autopsy-card">';
          html += '<div class="autopsy-field"><div class="autopsy-label">Subtask ID & Code:</div><div class="autopsy-value">#' + a.subtaskId + ' (' + (a.taskCode || '-') + ')</div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Task Title:</div><div class="autopsy-value">' + (a.title || '-') + '</div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Assigned Role:</div><div class="autopsy-value" style="color:var(--color-amber);">' + (a.assignedAgentRole || '-') + '</div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Last Action Type:</div><div class="autopsy-value"><span class="action-badge act-' + a.lastActionType + '">' + a.lastActionType + '</span></div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Target File:</div><div class="autopsy-value" style="color:var(--color-cyan);">' + (a.lastTargetFile || '-') + '</div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Action Details:</div><div class="autopsy-value">' + (a.lastActionDetails || '-') + '</div></div>';
          html += '<div class="autopsy-field"><div class="autopsy-label">Last Timestamp:</div><div class="autopsy-value">' + (a.lastTimestamp || '-') + '</div></div>';
          html += '<div class="diagnosis-box"><strong>Forensic Diagnosis:</strong><br>' + (a.diagnosis || 'Agent crashed or timed out before recording completion evidence.') + '</div>';
          html += '</div>';
        }
        body.innerHTML = html;
      }
      modal.classList.add('open');
    }

    function closeCrashModal() {
      document.getElementById('crashModal').classList.remove('open');
    }

    function openDetailModal(idx) {
      var log = window._currentLogs && window._currentLogs[idx];
      if (!log) return;
      var body = document.getElementById('detailModalBody');
      var html = '<div class="autopsy-card" style="border-color:var(--border-strong);">';
      html += '<div class="autopsy-field"><div class="autopsy-label">Action Log ID:</div><div class="autopsy-value">#' + log.actionLogId + '</div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Subtask ID:</div><div class="autopsy-value">' + log.subtaskId + '</div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Agent Role:</div><div class="autopsy-value" style="color:var(--color-amber);">' + log.agentRole + '</div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Action Type:</div><div class="autopsy-value"><span class="action-badge act-' + log.actionType + '">' + log.actionType + '</span></div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Target File:</div><div class="autopsy-value" style="color:var(--color-cyan);">' + (log.targetFile || '-') + '</div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Action Details:</div><div class="autopsy-value">' + escapeHtml(log.actionDetails || '-') + '</div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Status:</div><div class="autopsy-value"><span class="badge badge-' + log.status + '">' + log.status + '</span></div></div>';
      html += '<div class="autopsy-field"><div class="autopsy-label">Recorded At:</div><div class="autopsy-value">' + log.createdAt + '</div></div>';
      html += '</div>';
      body.innerHTML = html;
      document.getElementById('detailModal').classList.add('open');
    }

    function closeDetailModal() {
      document.getElementById('detailModal').classList.remove('open');
    }

    function confirmClearTask() {
      var taskId = appState.selectedTaskId || '';
      var confirmed = window.confirm('Are you sure you want to clear task lifecycle data? This will purge temporary databases for this run.');
      if (confirmed) {
        fetch('/api/agent/clear', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ taskId: taskId })
        }).then(function(r) { return r.json(); }).then(function(res) {
          alert('Task cleared: ' + (res.status || 'OK'));
          fetchFullDashboard();
        });
      }
    }

    function triggerClearCurrentTask() {
      closeCrashModal();
      confirmClearTask();
    }

    function escapeHtml(str) {
      if (!str) return '';
      return String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
    }

    window.onload = init;
  </script>
</body>
</html>`

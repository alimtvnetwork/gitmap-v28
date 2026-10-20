package cmdschedule

const scheduleAppHTMLHeader = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>GitMap Schedule & Notifier Settings</title>`

const scheduleAppHTMLBody = `
<div class="dashboard-container">
	<div class="header-bar">
		<div class="header-title">
			<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
			GitMap Schedule & Notifier Settings
		</div>
		<label class="switch-label">
			<input type="checkbox" id="maskToggle" checked>
			<span>Privacy Masking (u***r@d***n.com)</span>
		</label>
	</div>

	<div class="section-card">
		<div class="section-header">
			<span class="section-title">Active Cron & Interval Tasks</span>
			<button type="button" class="btn btn-primary" id="btnRefreshTasks">Refresh</button>
		</div>
		<table class="data-table">
			<thead><tr><th>Name</th><th>Interval</th><th>Command</th><th>Last Run</th><th>Status</th><th>Toggle</th></tr></thead>
			<tbody id="tasksTableBody"><tr><td colspan="6" style="text-align:center; color:var(--text-muted);">Loading tasks...</td></tr></tbody>
		</table>
	</div>

	<div class="section-card">
		<div class="section-header"><span class="section-title">Folder Watcher Triggers</span></div>
		<table class="data-table">
			<thead><tr><th>Target Path</th><th>Events</th><th>Debounce</th><th>Action Command</th><th>Status</th></tr></thead>
			<tbody id="watchersTableBody"><tr><td colspan="5" style="text-align:center; color:var(--text-muted);">Loading watchers...</td></tr></tbody>
		</table>
		<div class="form-grid">
			<div class="form-group"><label>Folder Path</label><input type="text" id="watcherPath" placeholder="/repo/subfolder"></div>
			<div class="form-group"><label>Events (comma-separated)</label><input type="text" id="watcherEvents" value="modify,create" placeholder="modify, create"></div>
			<div class="form-group"><label>Debounce (s)</label><input type="number" id="watcherDebounce" value="3" min="1" max="60"></div>
			<div class="form-group"><label>Action Command</label><input type="text" id="watcherAction" placeholder="gitmap pe -t"></div>
		</div>
		<div style="margin-top: 12px; display: flex; justify-content: flex-end;">
			<button type="button" class="btn btn-primary" id="btnAddWatcher">Add Folder Watcher</button>
		</div>
	</div>

	<div class="section-card">
		<div class="section-header"><span class="section-title">Notifier Recipient Manager</span></div>
		<table class="data-table">
			<thead><tr><th>Type</th><th>Destination Target</th><th>Trigger Events</th><th>Status</th></tr></thead>
			<tbody id="notifiersTableBody"><tr><td colspan="4" style="text-align:center; color:var(--text-muted);">Loading notifiers...</td></tr></tbody>
		</table>
	</div>
</div>`

const scheduleAppHTMLFooter = `</body></html>`

// ScheduleUIJavaScript provides client-side JavaScript for the schedule dashboard.
var ScheduleUIJavaScript = `
(function() {
	var isMasked = true, tasks = [], watchers = [], notifiers = [];
	document.addEventListener('DOMContentLoaded', function() {
		initPrivacyMask(); initActions(); loadTasks(); loadWatchers(); loadNotifiers();
	});
	function initPrivacyMask() {
		var toggle = document.getElementById('maskToggle');
		toggle.checked = isMasked;
		toggle.addEventListener('change', function(e) { isMasked = e.target.checked; renderAllTables(); });
	}
	function maskText(str) {
		if (!str || !isMasked) return str;
		if (str.indexOf('@') > 0) {
			var p = str.split('@'), loc = p[0], dom = p[1];
			var mLoc = loc.length <= 2 ? loc[0] + '***' : loc[0] + '***' + loc[loc.length - 1];
			var dIdx = dom.lastIndexOf('.'), mDom = dom;
			if (dIdx > 0) {
				var nm = dom.substring(0, dIdx), tld = dom.substring(dIdx);
				mDom = (nm.length <= 2 ? nm[0] + '***' : nm[0] + '***' + nm[nm.length - 1]) + tld;
			}
			return mLoc + '@' + mDom;
		}
		if (str.indexOf('/') >= 0 || str.indexOf('\\') >= 0) {
			var segs = str.split(/[/\\]/);
			for (var i = 1; i < segs.length; i++) {
				if (segs[i].length > 2) segs[i] = segs[i][0] + '***' + segs[i][segs[i].length - 1];
			}
			return segs.join('/');
		}
		return str;
	}
	function renderAllTables() { renderTasks(); renderWatchers(); renderNotifiers(); }
	function loadTasks() {
		fetch('/api/schedule/tasks').then(function(r) { return r.json(); }).then(function(data) {
			tasks = data || []; renderTasks();
		});
	}
	function renderTasks() {
		var tb = document.getElementById('tasksTableBody'); tb.innerHTML = '';
		if (tasks.length === 0) { tb.innerHTML = '<tr><td colspan="6" style="text-align:center;">No scheduled tasks registered.</td></tr>'; return; }
		tasks.forEach(function(t) {
			var tr = document.createElement('tr');
			var badge = t.isEnabled ? '<span class="badge badge-active">Active</span>' : '<span class="badge badge-inactive">Paused</span>';
			var chk = t.isEnabled ? 'checked' : '';
			tr.innerHTML = '<td><strong>' + escapeHTML(t.name) + '</strong></td>' +
				'<td><span class="code-pill">' + escapeHTML(t.interval) + '</span></td>' +
				'<td><span class="code-pill">' + escapeHTML(maskText(t.commandLine)) + '</span></td>' +
				'<td>' + (t.lastRunAt || 'Never') + '</td>' +
				'<td>' + badge + '</td>' +
				'<td><label class="toggle-switch"><input type="checkbox" data-name="' + escapeHTML(t.name) + '" ' + chk + ' class="task-toggle"><span class="slider"></span></label></td>';
			tb.appendChild(tr);
		});
		tb.querySelectorAll('.task-toggle').forEach(function(el) {
			el.addEventListener('change', function(e) { toggleTask(e.target.getAttribute('data-name'), e.target.checked); });
		});
	}
	function toggleTask(name, isEnabled) {
		fetch('/api/schedule/toggle', {
			method: 'POST', headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ name: name, isEnabled: isEnabled })
		}).then(function(r) { return r.json(); }).then(function() { loadTasks(); });
	}
	function loadWatchers() {
		fetch('/api/schedule/watchers').then(function(r) { return r.json(); }).then(function(data) {
			watchers = data || []; renderWatchers();
		});
	}
	function renderWatchers() {
		var tb = document.getElementById('watchersTableBody'); tb.innerHTML = '';
		if (watchers.length === 0) { tb.innerHTML = '<tr><td colspan="5" style="text-align:center;">No folder watchers active.</td></tr>'; return; }
		watchers.forEach(function(w) {
			var tr = document.createElement('tr');
			var evts = (w.events || []).join(', ');
			tr.innerHTML = '<td><span class="code-pill">' + escapeHTML(maskText(w.path)) + '</span></td>' +
				'<td>' + escapeHTML(evts) + '</td><td>' + w.debounceSeconds + 's</td>' +
				'<td><span class="code-pill">' + escapeHTML(maskText(w.actionCommand)) + '</span></td>' +
				'<td><span class="badge badge-active">Watching</span></td>';
			tb.appendChild(tr);
		});
	}
	function loadNotifiers() {
		fetch('/api/schedule/notifiers').then(function(r) { return r.json(); }).then(function(data) {
			notifiers = data || []; renderNotifiers();
		});
	}
	function renderNotifiers() {
		var tb = document.getElementById('notifiersTableBody'); tb.innerHTML = '';
		if (notifiers.length === 0) { tb.innerHTML = '<tr><td colspan="4" style="text-align:center;">No notifiers configured.</td></tr>'; return; }
		notifiers.forEach(function(n) {
			var tr = document.createElement('tr');
			var badge = n.isEnabled ? '<span class="badge badge-active">Enabled</span>' : '<span class="badge badge-inactive">Muted</span>';
			tr.innerHTML = '<td><strong>' + escapeHTML(n.type.toUpperCase()) + '</strong></td>' +
				'<td><span class="code-pill">' + escapeHTML(maskText(n.target)) + '</span></td>' +
				'<td>' + escapeHTML(n.events) + '</td><td>' + badge + '</td>';
			tb.appendChild(tr);
		});
	}
	function initActions() {
		document.getElementById('btnRefreshTasks').addEventListener('click', loadTasks);
		document.getElementById('btnAddWatcher').addEventListener('click', createWatcher);
	}
	function createWatcher() {
		var path = document.getElementById('watcherPath').value.trim();
		if (!path) { alert('Folder path is required'); return; }
		var events = document.getElementById('watcherEvents').value.split(',').map(function(s) { return s.trim(); }).filter(Boolean);
		var debounce = parseInt(document.getElementById('watcherDebounce').value, 10) || 3;
		var action = document.getElementById('watcherAction').value.trim();
		var payload = { path: path, events: events, debounceSeconds: debounce, actionCommand: action, isEnabled: true };
		fetch('/api/schedule/watchers', {
			method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload)
		}).then(function(r) { return r.json(); }).then(function(res) {
			if (res.isSuccess) { document.getElementById('watcherPath').value = ''; loadWatchers(); }
			else { alert(res.error || 'Failed to add watcher'); }
		});
	}
	function escapeHTML(s) { return (s || '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;'); }
})();
`

package cmdmail

const mailAppHTMLHeader = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>GitMap Mail Composer</title>`

const mailAppHTMLBody = `
<div class="composer-container">
	<div class="header-bar">
		<div class="header-title">
			<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#89b4fa" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
			GitMap Mail Composer & Mailbox Studio
		</div>
		<div class="header-actions">
			<label class="switch-label"><input type="checkbox" id="maskToggle" checked><span>Privacy Mask</span></label>
			<button type="button" class="btn btn-sm btn-secondary" id="btnRefreshMailboxes" title="Refresh Mailbox Status">
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>
				Check Mailboxes
			</button>
		</div>
	</div>
	<div class="mailbox-panel" id="mailboxPanel">
		<div class="panel-header">
			<span class="panel-title">Configured Mailboxes</span>
			<span class="unread-badge" id="totalUnreadBadge">0 unread</span>
		</div>
		<div class="mailbox-cards" id="mailboxCards">
			<div class="mailbox-loading">Inspecting mailbox health & unread messages...</div>
		</div>
	</div>
	<div class="row">
		<div class="col"><label for="accountSelect">Sender Account</label><select id="accountSelect"><option value="">Loading accounts...</option></select></div>
		<div class="col"><label for="templateSelect">Apply Template</label><select id="templateSelect"><option value="">Select a template...</option></select></div>
	</div>
	<div class="row">
		<div class="col"><label for="toInput">To Recipients (Press Enter or comma to add)</label><div class="chips-container" id="toChipsContainer"><input type="text" class="chips-input" id="toInput" placeholder="recipient@domain.com"></div></div>
	</div>
	<div class="row">
		<div class="col"><label for="ccInput">CC Recipients</label><div class="chips-container" id="ccChipsContainer"><input type="text" class="chips-input" id="ccInput" placeholder="cc@domain.com"></div></div>
	</div>
	<div class="row">
		<div class="col"><label for="subjectInput">Subject</label><input type="text" id="subjectInput" placeholder="Subject headline..."></div>
	</div>
	<div class="tabs">
		<button type="button" class="tab-btn active" id="tabCompose">Compose (Markdown)</button>
		<button type="button" class="tab-btn" id="tabPreview">HTML Preview</button>
	</div>
	<div id="composeContainer"><textarea id="mailBody" placeholder="Compose message in Markdown format..."></textarea></div>
	<div id="previewContainer" style="display: none;"><div id="previewPane"></div></div>
	<div class="action-bar">
		<button type="button" class="btn btn-secondary" id="btnClear">Clear</button>
		<button type="button" class="btn btn-primary" id="btnSend">
			<svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
			Send Email
		</button>
	</div>
</div>
<div class="modal" id="statusModal">
	<div class="modal-box">
		<div class="modal-title" id="modalTitle">Sending Dispatch</div>
		<div class="modal-msg" id="modalMsg">Connecting to mail agent...</div>
		<button type="button" class="btn btn-secondary" id="modalCloseBtn">Close</button>
	</div>
</div>`

const mailAppHTMLFooter = `</body></html>`

// MailUIJavaScript provides the client-side JavaScript for offline mail composer operations.
var MailUIJavaScript = `
(function() {
	var toList = [], ccList = [], isMasked = true, templates = [];
	document.addEventListener('DOMContentLoaded', function() {
		initMaskToggle(); initChipInputs(); initTabs(); initActions(); fetchAccounts(); fetchTemplates(); fetchMailboxStatus();
	});
	function initMaskToggle() {
		var toggle = document.getElementById('maskToggle');
		toggle.checked = isMasked;
		toggle.addEventListener('change', function(e) { isMasked = e.target.checked; renderAllChips(); });
	}
	function maskEmail(email) {
		if (!email || email.indexOf('@') === -1) return email;
		var p = email.split('@'), local = p[0], domain = p[1];
		var maskedLocal = local.length <= 2 ? local[0] + '***' : local[0] + '***' + local[local.length - 1];
		var dotIdx = domain.lastIndexOf('.'), maskedDomain = domain;
		if (dotIdx > 0) {
			var name = domain.substring(0, dotIdx), tld = domain.substring(dotIdx);
			var dMask = name.length <= 2 ? name[0] + '***' : name[0] + '***' + name[name.length - 1];
			maskedDomain = dMask + tld;
		}
		return maskedLocal + '@' + maskedDomain;
	}
	function initChipInputs() {
		setupChipField('toInput', 'toChipsContainer', toList);
		setupChipField('ccInput', 'ccChipsContainer', ccList);
	}
	function setupChipField(inputId, containerId, list) {
		var input = document.getElementById(inputId);
		input.addEventListener('keydown', function(e) {
			if (e.key === 'Enter' || e.key === ',') { e.preventDefault(); addRecipientFromInput(input, list, containerId); }
		});
	}
	function addRecipientFromInput(input, list, containerId) {
		var val = input.value.trim().replace(/,/g, '');
		if (val.length > 0 && list.indexOf(val) === -1) { list.push(val); input.value = ''; renderChips(containerId, list); }
	}
	function renderAllChips() { renderChips('toChipsContainer', toList); renderChips('ccChipsContainer', ccList); }
	function renderChips(containerId, list) {
		var container = document.getElementById(containerId), input = container.querySelector('.chips-input');
		var existing = container.querySelectorAll('.chip');
		for (var i = 0; i < existing.length; i++) existing[i].remove();
		list.forEach(function(email, idx) {
			var chip = document.createElement('span'); chip.className = 'chip';
			chip.innerHTML = (isMasked ? maskEmail(email) : email) + ' <span class="chip-remove" data-idx="' + idx + '">&times;</span>';
			chip.querySelector('.chip-remove').addEventListener('click', function() { list.splice(idx, 1); renderChips(containerId, list); });
			container.insertBefore(chip, input);
		});
	}
	function initTabs() {
		var tC = document.getElementById('tabCompose'), tP = document.getElementById('tabPreview');
		var cB = document.getElementById('composeContainer'), pB = document.getElementById('previewContainer');
		var pP = document.getElementById('previewPane'), txt = document.getElementById('mailBody');
		tC.addEventListener('click', function() { tC.classList.add('active'); tP.classList.remove('active'); cB.style.display = 'block'; pB.style.display = 'none'; });
		tP.addEventListener('click', function() { tP.classList.add('active'); tC.classList.remove('active'); cB.style.display = 'none'; pB.style.display = 'block'; pP.innerHTML = renderMarkdown(txt.value); });
	}
	function renderMarkdown(md) {
		if (!md) return '<p style="color:var(--text-muted)">Empty email content.</p>';
		var lines = md.split('\n'), html = [];
		for (var i = 0; i < lines.length; i++) {
			var line = lines[i];
			if (line.indexOf('### ') === 0) html.push('<h3>' + escapeHTML(line.slice(4)) + '</h3>');
			else if (line.indexOf('## ') === 0) html.push('<h2>' + escapeHTML(line.slice(3)) + '</h2>');
			else if (line.indexOf('# ') === 0) html.push('<h1>' + escapeHTML(line.slice(2)) + '</h1>');
			else if (line.indexOf('- ') === 0) html.push('<li>' + formatInline(line.slice(2)) + '</li>');
			else if (line.trim().length > 0) html.push('<p>' + formatInline(line) + '</p>');
		}
		return html.join('');
	}
	function formatInline(text) {
		var safe = escapeHTML(text);
		safe = safe.replace(/\\*\\*(.*?)\\*\\*/g, '<strong>$1</strong>');
		safe = safe.replace(/\\*(.*?)\\*/g, '<em>$1</em>');
		return safe.replace(new RegExp('\\x60([^\\x60]+)\\x60', 'g'), '<code>$1</code>');
	}
	function escapeHTML(str) { return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;'); }
	function initActions() {
		document.getElementById('btnSend').addEventListener('click', sendMail);
		document.getElementById('btnClear').addEventListener('click', clearForm);
		document.getElementById('modalCloseBtn').addEventListener('click', hideModal);
		document.getElementById('templateSelect').addEventListener('change', onTemplateChange);
		document.getElementById('btnRefreshMailboxes').addEventListener('click', fetchMailboxStatus);
	}
	function fetchMailboxStatus() {
		var container = document.getElementById('mailboxCards');
		container.innerHTML = '<div class="mailbox-loading">Updating mailbox telemetry...</div>';
		fetch('/api/mail/check').then(function(r) { return r.json(); }).then(function(summary) {
			renderMailboxSummary(summary);
		}).catch(function() {
			container.innerHTML = '<div class="mailbox-card"><span class="account-name">Local Mailer</span><span class="account-status status-active">Active</span></div>';
		});
	}
	function renderMailboxSummary(summary) {
		var container = document.getElementById('mailboxCards');
		container.innerHTML = '';
		var unreadBadge = document.getElementById('totalUnreadBadge');
		unreadBadge.textContent = (summary.totalUnread || 0) + ' unread';
		if (!summary.items || summary.items.length === 0) {
			container.innerHTML = '<div class="mailbox-card"><span class="account-name">Default Mail Gateway</span><span class="account-status status-active">Ready</span></div>';
			return;
		}
		summary.items.forEach(function(item) {
			var card = document.createElement('div');
			card.className = 'mailbox-card' + (item.hasUnread ? ' card-unread' : '');
			var addr = isMasked ? item.maskedAddress : item.emailAddress;
			card.innerHTML = '<div class="card-left"><span class="account-name">' + escapeHTML(item.accountName) + '</span><span class="account-email">' + escapeHTML(addr) + '</span></div><div class="card-right"><span class="unread-count">' + (item.unreadCount || 0) + ' unread</span><span class="account-status ' + (item.isActive ? 'status-active' : 'status-inactive') + '">' + escapeHTML(item.status || 'Active') + '</span></div>';
			container.appendChild(card);
		});
	}
	function fetchAccounts() {
		fetch('/api/mail/accounts').then(function(r) { return r.json(); }).then(function(accounts) {
			var sel = document.getElementById('accountSelect'); sel.innerHTML = '';
			accounts.forEach(function(acc) {
				var opt = document.createElement('option'); opt.value = acc.name;
				opt.textContent = acc.name + ' (' + acc.email + ')'; if (acc.isDefault) opt.selected = true;
				sel.appendChild(opt);
			});
		});
	}
	function fetchTemplates() {
		fetch('/api/mail/templates').then(function(r) { return r.json(); }).then(function(data) {
			templates = data; var sel = document.getElementById('templateSelect');
			sel.innerHTML = '<option value="">-- Choose Template --</option>';
			templates.forEach(function(tmpl) {
				var opt = document.createElement('option'); opt.value = tmpl.id; opt.textContent = tmpl.name;
				sel.appendChild(opt);
			});
		});
	}
	function onTemplateChange(e) {
		var id = e.target.value;
		for (var i = 0; i < templates.length; i++) {
			if (templates[i].id === id) {
				document.getElementById('subjectInput').value = templates[i].subject;
				document.getElementById('mailBody').value = templates[i].bodyMarkdown; break;
			}
		}
	}
	function sendMail() {
		if (toList.length === 0) { showModal('Validation Error', 'Please enter at least one recipient email.'); return; }
		var payload = { account: document.getElementById('accountSelect').value, to: toList, cc: ccList, subject: document.getElementById('subjectInput').value, bodyMarkdown: document.getElementById('mailBody').value, isHTML: true };
		showModal('Sending...', 'Dispatching payload to mail gateway...');
		fetch('/api/mail/send', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) })
		.then(function(r) { return r.json(); }).then(function(res) {
			showModal(res.isSuccess ? 'Sent Successfully' : 'Dispatch Failed', res.isSuccess ? res.message : (res.error || 'Server error'));
		}).catch(function(err) { showModal('Network Error', err.message); });
	}
	function clearForm() {
		toList = []; ccList = []; renderAllChips();
		document.getElementById('subjectInput').value = ''; document.getElementById('mailBody').value = '';
		document.getElementById('templateSelect').value = '';
	}
	function showModal(title, msg) {
		document.getElementById('modalTitle').textContent = title;
		document.getElementById('modalMsg').textContent = msg;
		document.getElementById('statusModal').classList.add('visible');
	}
	function hideModal() { document.getElementById('statusModal').classList.remove('visible'); }
})();
`

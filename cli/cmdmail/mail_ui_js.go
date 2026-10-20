//go:build !js || js

package cmdmail

// MailUIJavaScript contains the client-side logic for the Mail Composer.
const MailUIJavaScript = `
(function() {
	var toList = [];
	var ccList = [];
	var isMasked = true;
	var templates = [];

	document.addEventListener('DOMContentLoaded', function() {
		initMaskToggle();
		initChipInputs();
		initTabs();
		initActions();
		fetchAccounts();
		fetchTemplates();
	});

	function initMaskToggle() {
		var toggle = document.getElementById('maskToggle');
		toggle.checked = isMasked;
		toggle.addEventListener('change', function(e) {
			isMasked = e.target.checked;
			renderAllChips();
		});
	}

	function maskEmail(email) {
		if (!email || email.indexOf('@') === -1) return email;
		var parts = email.split('@');
		var local = parts[0];
		var domain = parts[1];
		var maskedLocal = local.length <= 2 ? local[0] + '***' : local[0] + '***' + local[local.length - 1];
		var dotIdx = domain.lastIndexOf('.');
		var maskedDomain = domain;
		if (dotIdx > 0) {
			var name = domain.substring(0, dotIdx);
			var tld = domain.substring(dotIdx);
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
			if (e.key === 'Enter' || e.key === ',') {
				e.preventDefault();
				addRecipientFromInput(input, list, containerId);
			}
		});
	}

	function addRecipientFromInput(input, list, containerId) {
		var val = input.value.trim().replace(/,/g, '');
		if (val.length > 0 && list.indexOf(val) === -1) {
			list.push(val);
			input.value = '';
			renderChips(containerId, list);
		}
	}

	function renderAllChips() {
		renderChips('toChipsContainer', toList);
		renderChips('ccChipsContainer', ccList);
	}

	function renderChips(containerId, list) {
		var container = document.getElementById(containerId);
		var input = container.querySelector('.chips-input');
		var existing = container.querySelectorAll('.chip');
		for (var i = 0; i < existing.length; i++) existing[i].remove();

		list.forEach(function(email, idx) {
			var chip = document.createElement('span');
			chip.className = 'chip';
			var displayText = isMasked ? maskEmail(email) : email;
			chip.innerHTML = displayText + ' <span class="chip-remove" data-idx="' + idx + '">&times;</span>';
			chip.querySelector('.chip-remove').addEventListener('click', function() {
				list.splice(idx, 1);
				renderChips(containerId, list);
			});
			container.insertBefore(chip, input);
		});
	}

	function initTabs() {
		var tabCompose = document.getElementById('tabCompose');
		var tabPreview = document.getElementById('tabPreview');
		var composeBox = document.getElementById('composeContainer');
		var previewBox = document.getElementById('previewContainer');
		var previewPane = document.getElementById('previewPane');
		var textarea = document.getElementById('mailBody');

		tabCompose.addEventListener('click', function() {
			tabCompose.classList.add('active');
			tabPreview.classList.remove('active');
			composeBox.style.display = 'block';
			previewBox.style.display = 'none';
		});

		tabPreview.addEventListener('click', function() {
			tabPreview.classList.add('active');
			tabCompose.classList.remove('active');
			composeBox.style.display = 'none';
			previewBox.style.display = 'block';
			previewPane.innerHTML = renderMarkdown(textarea.value);
		});
	}

	function renderMarkdown(md) {
		if (!md) return '<p style="color:var(--text-muted)">Empty email content.</p>';
		var lines = md.split('\n');
		var html = [];
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
		safe = safe.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
		safe = safe.replace(/\*(.*?)\*/g, '<em>$1</em>');
		safe = safe.replace(new RegExp('\\x60([^\\x60]+)\\x60', 'g'), '<code>$1</code>');
		return safe;
	}

	function escapeHTML(str) {
		return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
	}

	function initActions() {
		document.getElementById('btnSend').addEventListener('click', sendMail);
		document.getElementById('btnClear').addEventListener('click', clearForm);
		document.getElementById('modalCloseBtn').addEventListener('click', hideModal);
		document.getElementById('templateSelect').addEventListener('change', onTemplateChange);
	}

	function fetchAccounts() {
		fetch('/api/mail/accounts')
			.then(function(res) { return res.json(); })
			.then(function(accounts) {
				var select = document.getElementById('accountSelect');
				select.innerHTML = '';
				accounts.forEach(function(acc) {
					var opt = document.createElement('option');
					opt.value = acc.name;
					opt.textContent = acc.name + ' (' + acc.email + ')';
					if (acc.isDefault) opt.selected = true;
					select.appendChild(opt);
				});
			});
	}

	function fetchTemplates() {
		fetch('/api/mail/templates')
			.then(function(res) { return res.json(); })
			.then(function(data) {
				templates = data;
				var select = document.getElementById('templateSelect');
				select.innerHTML = '<option value="">-- Choose Template --</option>';
				templates.forEach(function(tmpl) {
					var opt = document.createElement('option');
					opt.value = tmpl.id;
					opt.textContent = tmpl.name;
					select.appendChild(opt);
				});
			});
	}

	function onTemplateChange(e) {
		var tmplId = e.target.value;
		for (var i = 0; i < templates.length; i++) {
			if (templates[i].id === tmplId) {
				document.getElementById('subjectInput').value = templates[i].subject;
				document.getElementById('mailBody').value = templates[i].bodyMarkdown;
				break;
			}
		}
	}

	function sendMail() {
		if (toList.length === 0) {
			showModal('Validation Error', 'Please enter at least one recipient email.');
			return;
		}
		var payload = {
			account: document.getElementById('accountSelect').value,
			to: toList,
			cc: ccList,
			subject: document.getElementById('subjectInput').value,
			bodyMarkdown: document.getElementById('mailBody').value,
			isHTML: true
		};

		showModal('Sending...', 'Dispatching payload to mail gateway...');
		fetch('/api/mail/send', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(payload)
		})
		.then(function(res) { return res.json(); })
		.then(function(result) {
			if (result.isSuccess) {
				showModal('Sent Successfully', result.message);
			} else {
				showModal('Dispatch Failed', result.error || 'Server error occurred.');
			}
		})
		.catch(function(err) {
			showModal('Network Error', err.message);
		});
	}

	function clearForm() {
		toList = [];
		ccList = [];
		renderAllChips();
		document.getElementById('subjectInput').value = '';
		document.getElementById('mailBody').value = '';
		document.getElementById('templateSelect').value = '';
	}

	function showModal(title, msg) {
		document.getElementById('modalTitle').textContent = title;
		document.getElementById('modalMsg').textContent = msg;
		document.getElementById('statusModal').classList.add('visible');
	}

	function hideModal() {
		document.getElementById('statusModal').classList.remove('visible');
	}
})();
`

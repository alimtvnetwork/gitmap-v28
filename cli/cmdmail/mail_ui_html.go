package cmdmail

const mailAppHTMLHeader = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>GitMap Mail Composer</title>
`

const mailAppHTMLBody = `
<div class="composer-container">
	<div class="header-bar">
		<div class="header-title">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
				<polyline points="22,6 12,13 2,6"/>
			</svg>
			GitMap Mail Composer
		</div>
		<label class="switch-label">
			<input type="checkbox" id="maskToggle" checked>
			<span>Mask View (u***r@d***n.com)</span>
		</label>
	</div>

	<div class="row">
		<div class="col">
			<label for="accountSelect">Sender Account</label>
			<select id="accountSelect">
				<option value="">Loading accounts...</option>
			</select>
		</div>
		<div class="col">
			<label for="templateSelect">Apply Template</label>
			<select id="templateSelect">
				<option value="">Select a template...</option>
			</select>
		</div>
	</div>

	<div class="row">
		<div class="col">
			<label for="toInput">To Recipients (Press Enter or comma to add)</label>
			<div class="chips-container" id="toChipsContainer">
				<input type="text" class="chips-input" id="toInput" placeholder="recipient@domain.com">
			</div>
		</div>
	</div>

	<div class="row">
		<div class="col">
			<label for="ccInput">CC Recipients</label>
			<div class="chips-container" id="ccChipsContainer">
				<input type="text" class="chips-input" id="ccInput" placeholder="cc@domain.com">
			</div>
		</div>
	</div>

	<div class="row">
		<div class="col">
			<label for="subjectInput">Subject</label>
			<input type="text" id="subjectInput" placeholder="Subject headline...">
		</div>
	</div>

	<div class="tabs">
		<button type="button" class="tab-btn active" id="tabCompose">Compose (Markdown)</button>
		<button type="button" class="tab-btn" id="tabPreview">HTML Preview</button>
	</div>

	<div id="composeContainer">
		<textarea id="mailBody" placeholder="Compose message in Markdown format..."></textarea>
	</div>

	<div id="previewContainer" style="display: none;">
		<div id="previewPane"></div>
	</div>

	<div class="action-bar">
		<button type="button" class="btn btn-secondary" id="btnClear">Clear</button>
		<button type="button" class="btn btn-primary" id="btnSend">Send Email</button>
	</div>
</div>

<div class="modal" id="statusModal">
	<div class="modal-box">
		<div class="modal-title" id="modalTitle">Sending Dispatch</div>
		<div class="modal-msg" id="modalMsg">Connecting to mail agent...</div>
		<button type="button" class="btn btn-secondary" id="modalCloseBtn">Close</button>
	</div>
</div>
`

const mailAppHTMLFooter = `
</body>
</html>
`

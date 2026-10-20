package cmdmail

import (
	"strings"
)

// MailUICustomCSS defines the embedded stylesheet using GitMap's dark palette.
const MailUICustomCSS = `
:root {
	--bg-canvas: #090d16;
	--bg-card: #182235;
	--bg-input: #0d1524;
	--border-color: #1e293b;
	--accent-cyan: #38bdf8;
	--accent-hover: #0284c7;
	--text-primary: #f1f5f9;
	--text-muted: #94a3b8;
	--success-color: #34d399;
	--danger-color: #f87171;
}

* { box-sizing: border-box; margin: 0; padding: 0; }
body {
	background-color: var(--bg-canvas);
	color: var(--text-primary);
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
	font-size: 14px;
	line-height: 1.5;
	padding: 24px;
}

.composer-container {
	max-width: 900px;
	margin: 0 auto;
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 24px;
	box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
}

.header-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	border-bottom: 1px solid var(--border-color);
	padding-bottom: 16px;
	margin-bottom: 20px;
}

.header-title {
	display: flex;
	align-items: center;
	gap: 10px;
	font-size: 18px;
	font-weight: 700;
	color: var(--accent-cyan);
}

.switch-label {
	display: flex;
	align-items: center;
	gap: 8px;
	cursor: pointer;
	font-size: 13px;
	color: var(--text-muted);
}

.row {
	display: flex;
	gap: 16px;
	margin-bottom: 14px;
}

.col { flex: 1; }
label {
	display: block;
	font-size: 12px;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: var(--text-muted);
	margin-bottom: 6px;
}

input[type="text"], select, textarea {
	width: 100%;
	background: var(--bg-input);
	border: 1px solid var(--border-color);
	color: var(--text-primary);
	border-radius: 6px;
	padding: 8px 12px;
	font-size: 14px;
	outline: none;
	transition: border-color 0.2s;
}

input[type="text"]:focus, select:focus, textarea:focus {
	border-color: var(--accent-cyan);
}

.chips-container {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
	align-items: center;
	background: var(--bg-input);
	border: 1px solid var(--border-color);
	border-radius: 6px;
	padding: 6px 10px;
	min-height: 40px;
}

.chip {
	display: inline-flex;
	align-items: center;
	gap: 6px;
	background: rgba(56, 189, 248, 0.15);
	border: 1px solid rgba(56, 189, 248, 0.3);
	color: var(--accent-cyan);
	padding: 2px 8px;
	border-radius: 4px;
	font-size: 12px;
}

.chip-remove {
	cursor: pointer;
	font-weight: bold;
	color: var(--text-muted);
}
.chip-remove:hover { color: var(--danger-color); }

.chips-input {
	flex: 1;
	min-width: 140px;
	border: none;
	background: transparent;
	color: var(--text-primary);
	outline: none;
	padding: 2px 4px;
}

.tabs {
	display: flex;
	gap: 8px;
	border-bottom: 1px solid var(--border-color);
	margin-bottom: 12px;
}

.tab-btn {
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	color: var(--text-muted);
	padding: 8px 16px;
	cursor: pointer;
	font-weight: 600;
}

.tab-btn.active {
	color: var(--accent-cyan);
	border-bottom-color: var(--accent-cyan);
}

textarea#mailBody {
	min-height: 220px;
	resize: vertical;
	font-family: monospace;
}

#previewPane {
	min-height: 220px;
	background: var(--bg-input);
	border: 1px solid var(--border-color);
	border-radius: 6px;
	padding: 16px;
	overflow-y: auto;
}

#previewPane h1, #previewPane h2, #previewPane h3 {
	color: var(--accent-cyan);
	margin-bottom: 8px;
}

#previewPane p { margin-bottom: 8px; }
#previewPane code {
	background: #1e293b;
	padding: 2px 6px;
	border-radius: 4px;
	font-size: 13px;
}

.action-bar {
	display: flex;
	justify-content: flex-end;
	gap: 12px;
	margin-top: 18px;
}

.btn {
	border: none;
	border-radius: 6px;
	padding: 10px 20px;
	font-weight: 600;
	cursor: pointer;
	transition: background-color 0.2s;
}

.btn-primary {
	background: var(--accent-cyan);
	color: #090d16;
}
.btn-primary:hover { background: var(--accent-hover); }

.btn-secondary {
	background: transparent;
	border: 1px solid var(--border-color);
	color: var(--text-muted);
}
.btn-secondary:hover { color: var(--text-primary); border-color: var(--text-muted); }

.modal {
	display: none;
	position: fixed;
	top: 0; left: 0; right: 0; bottom: 0;
	background: rgba(9, 13, 22, 0.85);
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.modal.visible { display: flex; }
.modal-box {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 24px;
	max-width: 440px;
	width: 90%;
	text-align: center;
}
.modal-title { font-size: 16px; font-weight: 700; margin-bottom: 12px; }
.modal-msg { color: var(--text-muted); margin-bottom: 18px; }
`

// BuildMailComposerHTML stitches together the zero-CDN single page HTML document.
func BuildMailComposerHTML() string {
	var sb strings.Builder
	sb.WriteString(mailAppHTMLHeader)
	sb.WriteString("<style>\n")
	sb.WriteString(MailUICustomCSS)
	sb.WriteString("\n</style>\n</head>\n<body>\n")
	sb.WriteString(mailAppHTMLBody)
	sb.WriteString("\n<script>\n")
	sb.WriteString(MailUIJavaScript)
	sb.WriteString("\n</script>\n")
	sb.WriteString(mailAppHTMLFooter)

	return sb.String()
}

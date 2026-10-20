package cmdschedule

import (
	"strings"
)

// ScheduleUICustomCSS defines embedded stylesheet with Dark GitMap styling.
const ScheduleUICustomCSS = `
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

.dashboard-container {
	max-width: 1040px;
	margin: 0 auto;
}

.header-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 16px 24px;
	margin-bottom: 20px;
}

.header-title {
	display: flex;
	align-items: center;
	gap: 12px;
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

.section-card {
	background: var(--bg-card);
	border: 1px solid var(--border-color);
	border-radius: 8px;
	padding: 20px;
	margin-bottom: 20px;
}

.section-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
	padding-bottom: 8px;
	border-bottom: 1px solid var(--border-color);
}

.section-title {
	font-size: 16px;
	font-weight: 600;
	color: var(--text-primary);
}

.data-table {
	width: 100%;
	border-collapse: collapse;
	margin-top: 8px;
}

.data-table th {
	text-align: left;
	font-size: 12px;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: var(--text-muted);
	padding: 10px 12px;
	border-bottom: 1px solid var(--border-color);
}

.data-table td {
	padding: 12px;
	border-bottom: 1px solid rgba(30, 41, 59, 0.5);
	vertical-align: middle;
}

.data-table tr:hover td {
	background-color: rgba(56, 189, 248, 0.04);
}

.badge {
	display: inline-block;
	padding: 2px 8px;
	border-radius: 12px;
	font-size: 11px;
	font-weight: 600;
}

.badge-active { background: rgba(52, 211, 153, 0.15); color: var(--success-color); border: 1px solid rgba(52, 211, 153, 0.3); }
.badge-inactive { background: rgba(148, 163, 184, 0.15); color: var(--text-muted); border: 1px solid rgba(148, 163, 184, 0.3); }

.toggle-switch {
	position: relative;
	display: inline-block;
	width: 36px;
	height: 20px;
}
.toggle-switch input { opacity: 0; width: 0; height: 0; }
.slider {
	position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0;
	background-color: #334155; transition: .2s; border-radius: 20px;
}
.slider:before {
	position: absolute; content: ""; height: 14px; width: 14px; left: 3px; bottom: 3px;
	background-color: white; transition: .2s; border-radius: 50%;
}
input:checked + .slider { background-color: var(--accent-cyan); }
input:checked + .slider:before { transform: translateX(16px); }

.form-grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
	gap: 12px;
	margin-top: 14px;
}

.form-group { display: flex; flex-direction: column; gap: 4px; }
.form-group label { font-size: 12px; font-weight: 600; color: var(--text-muted); }
.form-group input, .form-group select {
	background: var(--bg-input); border: 1px solid var(--border-color);
	color: var(--text-primary); border-radius: 6px; padding: 8px 10px; font-size: 13px; outline: none;
}
.form-group input:focus { border-color: var(--accent-cyan); }

.btn {
	padding: 8px 16px; border-radius: 6px; font-weight: 600; cursor: pointer;
	border: none; font-size: 13px; transition: background-color 0.2s;
}
.btn-primary { background: var(--accent-cyan); color: #090d16; }
.btn-primary:hover { background: var(--accent-hover); }

.code-pill {
	background: #0d1524; padding: 2px 6px; border-radius: 4px; font-family: monospace; font-size: 12px; border: 1px solid var(--border-color);
}
`

// BuildScheduleDashboardHTML stitches together the complete dashboard SPA.
func BuildScheduleDashboardHTML() string {
	var sb strings.Builder
	sb.WriteString(scheduleAppHTMLHeader)
	sb.WriteString("<style>\n")
	sb.WriteString(ScheduleUICustomCSS)
	sb.WriteString("\n</style>\n</head>\n<body>\n")
	sb.WriteString(scheduleAppHTMLBody)
	sb.WriteString("\n<script>\n")
	sb.WriteString(ScheduleUIJavaScript)
	sb.WriteString("\n</script>\n")
	sb.WriteString(scheduleAppHTMLFooter)

	return sb.String()
}

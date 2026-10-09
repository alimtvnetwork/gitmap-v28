package cmdui

// uiAssetsStyle holds the full <style> block for the embedded UI.
// Moved byte-identical from ui_assets_markup_misc.go (program 246, subtask 05)
// and extended: refined dark tokens + a real light theme, toast CSS,
// responsive shell CSS, and the new design-system token scales.
// NOTE: terminal/editor surfaces (--code-bg, --code-text, --json-*) stay dark
// in BOTH themes by design — a deliberate premium choice (spec §4.1).

const uiAssetsStyle = `  <style>
    :root {
      /* 4-Plane Neutral Depth Hierarchy — dark theme (default, refined) */
      --bg: hsl(230, 25%, 8%);
      --raised: hsl(230, 18%, 18%);
      --card: hsl(230, 20%, 12%);
      --popover: hsl(230, 20%, 16%);
      --border: hsl(230, 18%, 20%);
      --text: #f8fafc;
      --muted: #94a3b8;
      --primary: #f59e0b;
      --primary-hover: #d97706;
      --primary-ink: #000000;
      --accent: #fbbf24;
      --success: #10b981;
      --warning: #fbbf24;
      --info: #38bdf8;
      --error: #ef4444;
      --code-bg: hsl(230, 25%, 6%);
      --code-text: #e2e8f0;
      --term-prompt: #f59e0b;
      --json-key: #f59e0b;
      --json-str: #34d399;
      --json-num: #60a5fa;
      --json-bool: #f472b6;
      --sidebar-bg: hsl(230, 28%, 6%);
      --badge-ok-bg: #064e3b;
      --badge-ok-fg: #34d399;
      --badge-err-bg: #451a1a;
      --badge-err-fg: #f87171;
      /* Spacing scale (4pt) */
      --sp-1: 4px; --sp-2: 8px; --sp-3: 12px; --sp-4: 16px;
      --sp-5: 20px; --sp-6: 24px; --sp-8: 32px;
      /* Type scale */
      --fs-xs: 0.75rem; --fs-sm: 0.85rem; --fs-md: 0.95rem;
      --fs-lg: 1.1rem; --fs-xl: 1.35rem; --fs-2xl: 1.6rem;
      /* Radii, fonts, shadow */
      --radius-sm: 6px; --radius-md: 8px; --radius-lg: 12px;
      --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      --font-mono: "Cascadia Code", Consolas, monospace;
      --shadow-card: 0 8px 24px rgba(0, 0, 0, 0.35);
    }
    :root[data-theme="light"] {
      --bg: #eef2f7;
      --raised: #ffffff;
      --card: #ffffff;
      --popover: #ffffff;
      --border: #e2e8f0;
      --text: #0f172a;
      --muted: #64748b;
      --primary: #b45309;
      --primary-hover: #92400e;
      --primary-ink: #ffffff;
      --accent: #b45309;
      --success: #047857;
      --warning: #b45309;
      --info: #0284c7;
      --error: #dc2626;
      --code-bg: hsl(230, 25%, 8%);
      --sidebar-bg: #ffffff;
      --shadow-card: 0 8px 24px rgba(15, 23, 42, 0.12);
    }
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: var(--font-sans); }
    :focus-visible { outline: 2px solid var(--primary); outline-offset: 2px; }
    body { background: var(--bg); color: var(--text); display: flex; height: 100vh; overflow: hidden; transition: background 0.15s ease, color 0.15s ease; }
    #sidebar { width: 240px; background: var(--sidebar-bg); border-right: 1px solid var(--border); display: flex; flex-direction: column; }
    #sidebar .logo { padding: var(--sp-5) var(--sp-4); font-size: var(--fs-lg); font-weight: 700; color: var(--primary); border-bottom: 1px solid var(--border); display: flex; align-items: center; gap: var(--sp-2); }
    #sidebar nav { flex: 1; padding: var(--sp-3) var(--sp-2); display: flex; flex-direction: column; gap: var(--sp-1); overflow-y: auto; }
    .nav-btn { display: flex; align-items: center; gap: var(--sp-2); width: 100%; padding: var(--sp-2) var(--sp-3); border: none; background: transparent; color: var(--muted); border-radius: var(--radius-sm); cursor: pointer; text-align: left; font-size: var(--fs-sm); font-weight: 500; transition: all 0.15s ease; }
    .nav-btn:hover { background: var(--card); color: var(--text); }
    .nav-btn.active { background: var(--primary); color: var(--primary-ink); font-weight: 600; }
    .nav-btn .nav-icon { display: inline-flex; width: 18px; height: 18px; flex: none; }
    .nav-btn .nav-icon svg { width: 100%; height: 100%; }
    #main { flex: 1; display: flex; flex-direction: column; overflow: hidden; background: var(--bg); min-width: 0; }
    header { min-height: 56px; border-bottom: 1px solid var(--border); padding: var(--sp-2) var(--sp-6); display: flex; align-items: center; justify-content: space-between; gap: var(--sp-3); background: var(--card); }
    header h2 { font-size: var(--fs-lg); font-weight: 600; }
    .hamburger { display: none; background: transparent; border: none; color: var(--text); cursor: pointer; padding: var(--sp-2); border-radius: var(--radius-sm); align-items: center; }
    .hamburger:hover { background: var(--raised); }
    .sidebar-scrim { display: none; }
    #status-indicator { display: flex; align-items: center; font-size: var(--fs-sm); color: var(--muted); }
    .status-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--success); margin-right: var(--sp-2); animation: pulse-dot 2s infinite; }
    @keyframes pulse-dot { 0%, 100% { opacity: 1; } 50% { opacity: 0.4; } }
    .content-area { flex: 1; padding: var(--sp-6); overflow-y: auto; display: none; }
    .content-area.active { display: block; }
    .card { background: var(--card); border: 1px solid var(--border); box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05), 0 2px 4px rgba(0, 0, 0, 0.25); border-radius: var(--radius-md); padding: var(--sp-5); margin-bottom: var(--sp-5); }
    .card h3 { font-size: var(--fs-md); margin-bottom: var(--sp-3); color: var(--primary); }
    .card-hero { background: var(--popover); box-shadow: var(--shadow-card); border-top: 2px solid var(--primary); }
    .card-flat { background: transparent; box-shadow: none; }
    .form-group { margin-bottom: var(--sp-4); }
    .form-group label { display: block; font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-1); }
    input, select, textarea { width: 100%; background: var(--code-bg); border: 1px solid var(--border); color: var(--text); padding: var(--sp-2) var(--sp-3); border-radius: var(--radius-sm); font-size: var(--fs-sm); }
    input, select { color-scheme: dark; }
    :root[data-theme="light"] input, :root[data-theme="light"] select { color-scheme: light; }
    textarea { font-family: var(--font-mono); resize: vertical; }
    .btn { background: var(--primary); color: var(--primary-ink); font-weight: 600; padding: var(--sp-2) var(--sp-4); border-radius: var(--radius-sm); border: none; cursor: pointer; transition: all 0.15s ease; font-size: var(--fs-sm); }
    .btn:hover { background: var(--primary-hover); filter: brightness(1.04); }
    .btn-secondary { background: var(--raised); color: var(--text); margin-left: var(--sp-2); }
    .btn-secondary:hover { background: var(--popover); }
    table { width: 100%; border-collapse: collapse; margin-top: var(--sp-2); font-variant-numeric: tabular-nums; }
    th, td { padding: var(--sp-2) var(--sp-3); border: 1px solid var(--border); text-align: left; font-size: var(--fs-sm); }
    th { background: var(--raised); color: var(--muted); }
    .table-wrap { overflow-x: auto; }
    .badge { display: inline-block; padding: var(--sp-1) var(--sp-2); border-radius: var(--radius-sm); font-size: var(--fs-xs); font-weight: 600; }
    .badge-online { background: var(--badge-ok-bg); color: var(--badge-ok-fg); }
    .badge-offline { background: var(--badge-err-bg); color: var(--badge-err-fg); }
    #editor-container { display: flex; flex-direction: column; height: 500px; border: 1px solid var(--border); border-radius: var(--radius-sm); overflow: hidden; }
    #editor-toolbar { background: var(--raised); padding: var(--sp-2) var(--sp-3); display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid var(--border); }
    #editor-area { flex: 1; width: 100%; height: 100%; background: var(--code-bg); color: var(--code-text); font-family: var(--font-mono); padding: var(--sp-3); border: none; outline: none; line-height: 1.5; font-size: var(--fs-sm); }
    .json-key { color: var(--json-key); }
    .json-str { color: var(--json-str); }
    .json-num { color: var(--json-num); }
    .json-bool { color: var(--json-bool); }
    .sub-tab-btn { display: inline-flex; align-items: center; gap: var(--sp-2); background: transparent; border: 1px solid var(--border); color: var(--muted); padding: var(--sp-2) var(--sp-3); border-radius: var(--radius-sm); cursor: pointer; font-size: var(--fs-sm); font-weight: 500; transition: all 0.15s ease; }
    .sub-tab-btn:hover { background: var(--card); color: var(--text); }
    .sub-tab-btn.active { background: var(--primary); color: var(--primary-ink); font-weight: 600; border-color: var(--primary); }
    .sub-tab-btn .nav-icon { display: inline-flex; width: 16px; height: 16px; }
    .sub-tab-btn .nav-icon svg { width: 100%; height: 100%; }
    .settings-savebar { position: sticky; bottom: 0; background: var(--card); border-top: 1px solid var(--border); padding: var(--sp-3) var(--sp-4); display: flex; justify-content: flex-end; z-index: 10; }
    .chip { background: var(--raised); border: 1px solid var(--border); color: var(--muted); padding: var(--sp-1) var(--sp-3); border-radius: 9999px; font-size: var(--fs-xs); cursor: pointer; transition: all 0.15s ease; font-family: var(--font-mono); }
    .chip:hover { border-color: var(--primary); color: var(--text); background: var(--popover); }
    .term-stdout { color: var(--code-text); }
    .term-stderr { color: var(--error); }
    .term-info { color: var(--info); }
    .term-prompt { color: var(--term-prompt); font-weight: 600; }
    #toast-stack { position: fixed; right: var(--sp-4); bottom: var(--sp-4); z-index: 200; display: flex; flex-direction: column; gap: var(--sp-2); max-width: 360px; }
    .toast { display: flex; align-items: flex-start; gap: var(--sp-2); background: var(--popover); color: var(--text); border: 1px solid var(--border); border-left: 3px solid var(--info); border-radius: var(--radius-md); box-shadow: var(--shadow-card); padding: var(--sp-3) var(--sp-4); font-size: var(--fs-sm); cursor: pointer; animation: toast-in 0.25s ease; }
    .toast-success { border-left-color: var(--success); }
    .toast-error { border-left-color: var(--error); }
    .toast-info { border-left-color: var(--info); }
    .toast-icon { display: inline-flex; width: 18px; height: 18px; flex: none; }
    .toast-icon svg { width: 100%; height: 100%; }
    .toast-success .toast-icon { color: var(--success); }
    .toast-error .toast-icon { color: var(--error); }
    .toast-info .toast-icon { color: var(--info); }
    .toast-msg { flex: 1; overflow-wrap: anywhere; }
    .toast-out { opacity: 0; transform: translateX(16px); transition: opacity 0.25s ease, transform 0.25s ease; }
    @keyframes toast-in { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: none; } }
    @media (max-width: 1200px) {
      #sidebar { width: 68px; }
      #sidebar .logo { justify-content: center; padding: var(--sp-4) var(--sp-2); }
      .nav-label { display: none; }
      .nav-btn { justify-content: center; padding: var(--sp-2); }
    }
    @media (max-width: 900px) {
      .hamburger { display: inline-flex; }
      #sidebar { position: fixed; left: 0; top: 0; bottom: 0; width: 240px; z-index: 100; transform: translateX(-105%); transition: transform 0.2s ease; }
      #sidebar.sidebar-open { transform: none; }
      #sidebar.sidebar-open .nav-label { display: inline; }
      #sidebar.sidebar-open .nav-btn { justify-content: flex-start; padding: var(--sp-2) var(--sp-3); }
      .sidebar-scrim.open { display: block; position: fixed; inset: 0; background: rgba(0, 0, 0, 0.45); z-index: 90; }
      .content-area { padding: var(--sp-4); }
      header { padding: var(--sp-2) var(--sp-4); }
    }
  </style>
`

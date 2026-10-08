package cmdui

const uiAssetsDocHead = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitMap Management & Fleet Dashboard</title>
`

const uiAssetsStyle = `  <style>
    :root {
      /* 4-Plane Neutral Depth Hierarchy */
      --bg: hsl(230, 25%, 8%);            /* Plane 0: Base Canvas */
      --raised: hsl(230, 18%, 18%);       /* Plane 1: Raised Wells */
      --card: hsl(230, 20%, 12%);         /* Plane 2: Surface Cards */
      --popover: hsl(230, 20%, 16%);      /* Plane 3: Elevated Flyouts */
      --border: hsl(230, 18%, 20%);       /* Hairline border */
      --text: #f8fafc;
      --muted: #94a3b8;
      --primary: #f59e0b;
      --primary-hover: #d97706;
      --success: #10b981;
      --error: #ef4444;
      --code-bg: hsl(230, 25%, 6%);
    }
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { background: var(--bg); color: var(--text); display: flex; height: 100vh; overflow: hidden; }
    #sidebar { width: 240px; background: hsl(230, 28%, 6%); border-right: 1px solid var(--border); display: flex; flex-direction: column; }
    #sidebar .logo { padding: 1.25rem 1rem; font-size: 1.15rem; font-weight: 700; color: var(--primary); border-bottom: 1px solid var(--border); display: flex; align-items: center; gap: 8px; }
    #sidebar nav { flex: 1; padding: 0.75rem 0.5rem; display: flex; flex-direction: column; gap: 4px; overflow-y: auto; }
    .nav-btn { display: flex; align-items: center; gap: 10px; width: 100%; padding: 0.65rem 0.85rem; border: none; background: transparent; color: var(--muted); border-radius: 6px; cursor: pointer; text-align: left; font-size: 0.9rem; font-weight: 500; transition: all 0.15s; }
    .nav-btn:hover { background: var(--card); color: var(--text); }
    .nav-btn.active { background: var(--primary); color: #000; font-weight: 600; }
    #main { flex: 1; display: flex; flex-direction: column; overflow: hidden; background: var(--bg); }
    header { height: 56px; border-bottom: 1px solid var(--border); padding: 0 1.5rem; display: flex; align-items: center; justify-content: space-between; background: var(--card); }
    header h2 { font-size: 1.1rem; font-weight: 600; }
    .content-area { flex: 1; padding: 1.5rem; overflow-y: auto; display: none; }
    .content-area.active { display: block; }
    .card { background: var(--card); border: 1px solid var(--border); box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.05), 0 2px 4px rgba(0, 0, 0, 0.25); border-radius: 8px; padding: 1.25rem; margin-bottom: 1.25rem; }
    .card h3 { font-size: 1rem; margin-bottom: 0.75rem; color: var(--primary); }
    .form-group { margin-bottom: 1rem; }
    .form-group label { display: block; font-size: 0.85rem; color: var(--muted); margin-bottom: 0.35rem; }
    input, select, textarea { width: 100%; background: var(--code-bg); border: 1px solid var(--border); color: var(--text); padding: 0.6rem 0.75rem; border-radius: 6px; font-size: 0.9rem; }
    textarea { font-family: monospace; resize: vertical; }
    .btn { background: var(--primary); color: #000; font-weight: 600; padding: 0.6rem 1.2rem; border-radius: 6px; border: none; cursor: pointer; transition: 0.15s; }
    .btn:hover { background: var(--primary-hover); filter: brightness(1.04); }
    .btn-secondary { background: var(--raised); color: var(--text); margin-left: 0.5rem; }
    .btn-secondary:hover { background: var(--popover); }
    table { width: 100%; border-collapse: collapse; margin-top: 0.5rem; }
    th, td { padding: 0.65rem 0.75rem; border: 1px solid var(--border); text-align: left; font-size: 0.85rem; }
    th { background: #131d2e; color: var(--muted); }
    .badge { display: inline-block; padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.75rem; font-weight: 600; }
    .badge-online { background: #064e3b; color: #34d399; }
    .badge-offline { background: #451a1a; color: #f87171; }
    #editor-container { display: flex; flex-direction: column; height: 500px; border: 1px solid var(--border); border-radius: 6px; overflow: hidden; }
    #editor-toolbar { background: #131d2e; padding: 0.5rem 0.75rem; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid var(--border); }
    #editor-area { flex: 1; width: 100%; height: 100%; background: var(--code-bg); color: #e2e8f0; font-family: "Cascadia Code", Consolas, monospace; padding: 0.75rem; border: none; outline: none; line-height: 1.5; font-size: 0.9rem; }
    .json-key { color: #f59e0b; }
    .json-str { color: #34d399; }
    .json-num { color: #60a5fa; }
    .json-bool { color: #f472b6; }
    .sub-tab-btn { background: transparent; border: 1px solid var(--border); color: var(--muted); padding: 0.45rem 0.9rem; border-radius: 6px; cursor: pointer; font-size: 0.85rem; font-weight: 500; transition: all 0.15s; }
    .sub-tab-btn:hover { background: var(--card); color: var(--text); }
    .sub-tab-btn.active { background: var(--primary); color: #000; font-weight: 600; border-color: var(--primary); }
    .chip { background: #131d2e; border: 1px solid var(--border); color: #cbd5e1; padding: 0.3rem 0.75rem; border-radius: 9999px; font-size: 0.8rem; cursor: pointer; transition: 0.15s; font-family: monospace; }
    .chip:hover { border-color: var(--primary); color: #fff; background: #1e293b; }
    .term-stdout { color: #f8fafc; }
    .term-stderr { color: #f87171; }
    .term-info { color: #38bdf8; }
    .term-prompt { color: #f59e0b; font-weight: 600; }
  </style>
`

const uiAssetsHeadClose = `</head>
`

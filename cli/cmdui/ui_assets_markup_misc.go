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

const uiAssetsTail = `</body>
</html>
`

const uiAssetsMarkupShell = `<body>
  <div id="sidebar">
    <div class="logo">⚡ GitMap Fleet UI</div>
    <nav>
      <button class="nav-btn active" onclick="showTab('settings')">⚙️ Settings</button>
      <button class="nav-btn" onclick="showTab('commitin')">📦 Commitin</button>
      <button class="nav-btn" onclick="showTab('ssh')">🖥️ SSH Fleet</button>
      <button class="nav-btn" onclick="showTab('macro')">📜 Macros</button>
      <button class="nav-btn" onclick="showTab('installer')">🚀 Installers</button>
      <button class="nav-btn" onclick="showTab('prompts')">💬 Prompts & AI</button>
      <button class="nav-btn" onclick="showTab('import-export')">🔄 Import / Export</button>
      <button class="nav-btn" onclick="showTab('schedules')">⏱️ Schedules</button>
      <button class="nav-btn" onclick="showTab('editor')">📝 Remote Editor</button>
      <button class="nav-btn" onclick="showTab('help')">📖 CLI & AGY Help</button>
    </nav>
  </div>
  <div id="main">
    <header>
      <h2 id="header-title">Settings</h2>
      <div id="status-indicator" style="font-size: 0.85rem; color: var(--muted);">Connected: localhost</div>
    </header>

    <!-- SETTINGS TAB -->
`

const uiAssetsMarkupMisc = `    <div id="tab-editor" class="content-area">
      <div class="card" style="margin-bottom: 0.5rem;">
        <div style="display: flex; gap: 10px; align-items: flex-end;">
          <div style="flex: 1;"><label style="font-size:0.8rem;color:var(--muted)">Target Node</label><select id="editor-node"><option value="local">local (current machine)</option></select></div>
          <div style="flex: 3;"><label style="font-size:0.8rem;color:var(--muted)">Absolute / Relative File Path</label><input type="text" id="editor-path" placeholder="/etc/gitmap.conf or cli/main.go"></div>
          <button class="btn" onclick="openRemoteFile()">Load File</button>
          <button class="btn" style="background:var(--success)" onclick="saveRemoteFile()">Save to Host</button>
        </div>
      </div>
      <div id="editor-container">
        <div id="editor-toolbar"><span id="editor-status" style="font-size:0.85rem;color:var(--muted)">No file loaded</span><span id="editor-lang" style="font-size:0.8rem;color:var(--primary)">Language: Plaintext</span></div>
        <textarea id="editor-area" spellcheck="false" placeholder="Remote file content will appear here..."></textarea>
      </div>
    </div>

    <!-- CLI & AGY HELP TAB -->
    <div id="tab-help" class="content-area">
      <div class="card">
        <h3>📖 CLI &amp; AGY Help — Speed Commands, LAP, Rerun-With-ID, Machine Identity &amp; Bot Reference</h3>
        <table>
          <thead>
            <tr>
              <th>Command &amp; Syntax</th>
              <th>Aliases</th>
              <th>Description &amp; Flags</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>gitmap agy add .</code> / <code>gitmap agy add &lt;path&gt;</code></td>
              <td><code>agy add</code></td>
              <td>Register current repository (<code>.</code>) or explicit <code>&lt;path&gt;</code> in Antigravity workspace registry and 24h SQLite sequence cache.</td>
            </tr>
            <tr>
              <td><code>gitmap agy add-read .</code> / <code>gitmap agy add-read &lt;path&gt;</code></td>
              <td><code>ar</code></td>
              <td>Register current repository (<code>.</code>) or <code>&lt;path&gt;</code> and immediately dispatch Read Memory onboarding prompt.</td>
            </tr>
            <tr>
              <td><code>gitmap agy rp ls</code> &amp; <code>gitmap agy rp prompts ls [--wc 200]</code></td>
              <td><code>running-projects</code></td>
              <td>List running projects with 24h sequence ID (<code>#1</code>), Project ID, Alias, Path, and <code>[convID]</code>, or render active prompt tree truncated to 200 words (<code>--wc T</code>, <code>--json</code>, <code>--file</code>).</td>
            </tr>
            <tr>
              <td><code>gitmap agy last-active-projects (lap) N [ls/help] [--limit Y] [--offset/skip Z] [--page P] [--wc T] [--json] [--file &lt;path&gt;]</code></td>
              <td><code>lap</code></td>
              <td>Inspect projects with activity in the last <code>N</code> hours (default <code>24</code>, <code>--limit 10</code>, <code>--offset/--skip</code>, <code>--page</code>, <code>--wc 200</code>, <code>--json</code>, <code>--file &lt;path&gt;</code>).</td>
            </tr>
            <tr>
              <td><code>gitmap agy rerun-with-id (rwi) &lt;id|alias|seq|path&gt; &lt;convid&gt; "&lt;prompt|file&gt;" [-p &lt;name&gt;]</code></td>
              <td><code>rwi</code></td>
              <td>Resolve target project by 24h sequence (<code>1</code>), ID, alias, or path and conversation (<code>&lt;convid&gt;</code> or <code>P1</code>), then inject prompt text or markdown file with optional prefix template (<code>-p</code>).</td>
            </tr>
            <tr>
              <td><code>gitmap agy rerun-with-convid (rwc / rwp) &lt;convid|P1&gt; "&lt;prompt|file&gt;" [-p &lt;name&gt;]</code></td>
              <td><code>rwc</code>, <code>rwp</code>, <code>rerun-with-prompt-id</code></td>
              <td>Auto-resolve project + conversation from <code>&lt;convid|P1&gt;</code> in the 24h SQLite sequence cache and inject prompt or file.</td>
            </tr>
            <tr>
              <td><code>gitmap agy fast-forward (ff / account-switch / asw)</code></td>
              <td><code>ff</code>, <code>account-switch</code>, <code>asw</code></td>
              <td>Fast-forward switch Antigravity account automatically when usage reaches threshold (default 15% remaining prod / 98% E2E test).</td>
            </tr>
            <tr>
              <td><code>gitmap machine ls/change/set/revert/help [--ssh]</code> &amp; <code>gitmap alias ls/change/set/revert/help [--ssh]</code></td>
              <td><code>gitmap os machine</code>, <code>gitmap os alias</code></td>
              <td>Inspect, set, change, or revert local and SSH fleet machine hostname and network alias (auto-defaults to Local IPv4; accepts <code>y</code> confirmation).</td>
            </tr>
            <tr>
              <td><code>gitmap telegram</code>, <code>gitmap email</code>, <code>gitmap settings</code>, <code>gitmap os help</code></td>
              <td><code>agy telegram</code>, <code>agy email</code>, <code>agy settings</code></td>
              <td>Configure two-way Telegram bot, speed SMTP email alerts, unified speed settings (<code>lap.default_hours=24</code>), and OS diagnostics help.</td>
            </tr>
            <tr>
              <td><code>gitmap deploy keys all</code> / <code>gitmap deploy-keys-all</code></td>
              <td><code>deploy-keys</code>, <code>deploy-key-all</code></td>
              <td>Distribute and trust all cluster SSH public keys across all fleet machines passwordlessly (no <code>ssh</code> subcommand prefix required).</td>
            </tr>
            <tr>
              <td><code>gitmap ssh export-json [file]</code> / <code>gitmap export-ssh</code></td>
              <td><code>ssh-export</code>, <code>nodes-export-json</code></td>
              <td>Export full fleet node configurations, IPs, hostnames, and public keys into JSON for instant transfer to another machine.</td>
            </tr>
            <tr>
              <td><code>gitmap ssh import-json [file]</code> / <code>gitmap import-ssh</code></td>
              <td><code>ssh-import</code>, <code>nodes-import-json</code></td>
              <td>Import fleet node configurations and public keys from JSON snapshot, automatically merging connections and authorized keys.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="card">
        <h3>Special Repositories (repo-secrets / rs &amp; repo-cache / rc)</h3>
        <table>
          <thead>
            <tr>
              <th>Command &amp; Syntax</th>
              <th>Aliases</th>
              <th>Description &amp; Sequenced Storage</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>gitmap cd rs</code> &amp; <code>gitmap cd rc</code></td>
              <td><code>cd repo-secrets</code>, <code>cd repo-cache</code></td>
              <td>Navigate directly into the configured <code>repo-secrets</code> (<code>rs</code>) or <code>repo-cache</code> (<code>rc</code>) repository directory.</td>
            </tr>
            <tr>
              <td><code>gitmap rs file &lt;path&gt;</code>, <code>gitmap rs folder &lt;path&gt;</code>, <code>gitmap rs text "&lt;secret&gt;"</code></td>
              <td><code>repo-secrets</code>, <code>rs</code></td>
              <td>Store files, folders, or inline secret strings into <code>repo-secrets</code> with auto-sequenced <code>XX-&lt;repo&gt;/01-&lt;slug&gt;.ext</code> hierarchy + auto-commit &amp; push.</td>
            </tr>
            <tr>
              <td><code>gitmap rc file &lt;script.ps1&gt;</code>, <code>gitmap rc folder &lt;dir&gt;</code>, <code>gitmap rc text "&lt;script&gt;" --ext .ps1</code></td>
              <td><code>repo-cache</code>, <code>repo-storage</code>, <code>rc</code></td>
              <td>Archive reusable PowerShell (<code>.ps1</code>) scripts, test harnesses, and folders into <code>repo-cache</code> under <code>XX-&lt;repo&gt;/01-&lt;slug&gt;.ext</code> + auto-commit &amp; push.</td>
            </tr>
            <tr>
              <td><code>gitmap settings set special_repos.secrets_name repo-secrets</code> &amp; <code>gitmap settings set special_repos.cache_name repo-cache</code></td>
              <td><code>secrets_repo</code>, <code>cache_repo</code></td>
              <td>Customize the default special repository directory names for <code>repo-secrets</code> and <code>repo-cache</code> across local and fleet workflows.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <script>
`

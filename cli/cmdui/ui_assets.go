package cmdui

// IndexHTML provides the embedded single-page UI application.
const IndexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitMap Management & Fleet Dashboard</title>
  <style>
    :root {
      --bg: #0f172a;
      --card: #1e293b;
      --border: #334155;
      --text: #f8fafc;
      --muted: #94a3b8;
      --primary: #f59e0b;
      --primary-hover: #d97706;
      --success: #10b981;
      --error: #ef4444;
      --code-bg: #0b1120;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; }
    body { background: var(--bg); color: var(--text); display: flex; height: 100vh; overflow: hidden; }
    #sidebar { width: 240px; background: #090d16; border-right: 1px solid var(--border); display: flex; flex-direction: column; }
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
    .card { background: var(--card); border: 1px solid var(--border); border-radius: 8px; padding: 1.25rem; margin-bottom: 1.25rem; }
    .card h3 { font-size: 1rem; margin-bottom: 0.75rem; color: var(--primary); }
    .form-group { margin-bottom: 1rem; }
    .form-group label { display: block; font-size: 0.85rem; color: var(--muted); margin-bottom: 0.35rem; }
    input, select, textarea { width: 100%; background: var(--code-bg); border: 1px solid var(--border); color: var(--text); padding: 0.6rem 0.75rem; border-radius: 6px; font-size: 0.9rem; }
    textarea { font-family: monospace; resize: vertical; }
    .btn { background: var(--primary); color: #000; font-weight: 600; padding: 0.6rem 1.2rem; border-radius: 6px; border: none; cursor: pointer; transition: 0.15s; }
    .btn:hover { background: var(--primary-hover); }
    .btn-secondary { background: #334155; color: var(--text); margin-left: 0.5rem; }
    .btn-secondary:hover { background: #475569; }
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
  </style>
</head>
<body>
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
    <div id="tab-settings" class="content-area active">
      <div class="card">
        <h3>General & UI Layout Settings</h3>
        <div class="form-group">
          <label>Theme</label>
          <select id="setting-theme"><option value="dark">Dark Theme (Standard)</option><option value="light">Light Theme</option></select>
        </div>
        <div class="form-group">
          <label>Graphics Rendering & Hardware Acceleration</label>
          <select id="setting-graphics"><option value="high">High Performance / Hardware Acceleration</option><option value="fast">Fast Standard</option><option value="plain">Plain / Low Graphics</option></select>
        </div>
        <div class="form-group">
          <label>Auto-Open Browser on UI Launch</label>
          <select id="setting-auto-open"><option value="true">Enabled (Launch default browser automatically)</option><option value="false">Disabled (Emit URL to console only)</option></select>
        </div>
        <div class="form-group">
          <label>Commit-In Migration Studio Layout</label>
          <select id="setting-commitin-layout"><option value="split">Split Dual-Pane (Config & Output)</option><option value="left">Left Focus (Config dominant)</option><option value="right">Right Focus (Output dominant)</option><option value="stacked">Stacked Vertical</option></select>
        </div>
        <div class="form-group">
          <label>Pull Direction Orientation</label>
          <select id="setting-pull-direction"><option value="pull-left">Pull-Left (Replay from target down to left)</option><option value="pull-right">Pull-Right (Replay from left source into right target)</option><option value="bidirectional">Bidirectional Sync</option></select>
        </div>
        <div class="form-group">
          <label>Commit-Pull PR Replay Mode</label>
          <select id="setting-pr-mode"><option value="merges">merges (Feature branches & release PRs)</option><option value="feature-per-commit">feature-per-commit</option><option value="direct">direct mainline</option></select>
        </div>
        <div class="form-group">
          <label>Default Remote Target</label>
          <input type="text" id="setting-remote" placeholder="origin or alias">
        </div>
        <div class="form-group">
          <label>Cluster REST Port</label>
          <input type="number" id="setting-port" value="49152">
        </div>
        <button class="btn" onclick="saveSettings()">Save Settings</button>
      </div>

      <div class="card">
        <h3>Speed Settings (AGY, LAP, Account Switch, Telegram Bot, Email &amp; Machine Identity)</h3>
        <div class="form-group">
          <label>LAP Default Lookback Hours (N = 24) (<code>lap.default_hours</code>)</label>
          <input type="number" id="setting-lap-hours" name="lap.default_hours" data-key="lap.default_hours" value="24" placeholder="24">
        </div>
        <div class="form-group">
          <label>Account Switch Fast-Forward Threshold (%) (<code>account_switch.threshold</code> — default 15% prod / 98% E2E test)</label>
          <input type="number" id="setting-account-switch-threshold" name="account_switch.threshold" data-key="account_switch.threshold" value="15" placeholder="15 (prod) or 98 (E2E test)">
        </div>
        <div class="form-group">
          <label>Telegram Two-Way Bot Setup — Bot Token (<code>telegram.bot_token</code>)</label>
          <input type="text" id="setting-telegram-token" name="telegram.bot_token" data-key="telegram.bot_token" placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11">
        </div>
        <div class="form-group">
          <label>Telegram Two-Way Bot Setup — Chat ID (<code>telegram.chat_id</code>)</label>
          <input type="text" id="setting-telegram-chat" name="telegram.chat_id" data-key="telegram.chat_id" placeholder="-1001234567890">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup — SMTP Host (<code>email.smtp_host</code>)</label>
          <input type="text" id="setting-email-smtp" name="email.smtp_host" data-key="email.smtp_host" placeholder="smtp.gmail.com:587">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup — From Address (<code>email.from</code>)</label>
          <input type="email" id="setting-email-from" name="email.from" data-key="email.from" placeholder="bot@example.com">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup — To Address (<code>email.to</code>)</label>
          <input type="email" id="setting-email-to" name="email.to" data-key="email.to" placeholder="dev@example.com">
        </div>
        <div class="form-group">
          <label>Machine Identity &amp; Network Alias — Machine Name (<code>machine.name</code>)</label>
          <input type="text" id="setting-machine-name" name="machine.name" data-key="machine.name" placeholder="dev-win-01 / ubuntu-node-02">
        </div>
        <div class="form-group">
          <label>Machine Identity &amp; Network Alias — Network Alias (<code>machine.alias</code> — auto-defaults to Local IPv4)</label>
          <input type="text" id="setting-machine-alias" name="machine.alias" data-key="machine.alias" placeholder="Auto-defaults to Local IPv4 if unset">
        </div>
        <div class="form-group">
          <label>Special Secrets Repository Name (<code>special_repos.secrets_name</code> — shortcut: <code>gitmap rs</code> / <code>gitmap cd rs</code>)</label>
          <input type="text" id="setting-special-secrets-name" name="special_repos.secrets_name" data-key="special_repos.secrets_name" value="repo-secrets" placeholder="repo-secrets">
        </div>
        <div class="form-group">
          <label>Special Cache / Storage Repository Name (<code>special_repos.cache_name</code> — shortcut: <code>gitmap rc</code> / <code>gitmap cd rc</code>)</label>
          <input type="text" id="setting-special-cache-name" name="special_repos.cache_name" data-key="special_repos.cache_name" value="repo-cache" placeholder="repo-cache">
        </div>
        <button class="btn" onclick="saveSettings()">Save Speed Settings</button>
      </div>
    </div>

    <!-- COMMITIN TAB -->
    <div id="tab-commitin" class="content-area">
      <div class="card">
        <h3>Commitin Multi-Option Engine</h3>
        <div class="form-group">
          <label>Commit Message</label>
          <input type="text" id="commit-msg" placeholder="feat(scope): concise description">
        </div>
        <div class="form-group">
          <label>Commit Direction / Flow</label>
          <select id="commit-direction">
            <option value="standard">Standard Commit</option>
            <option value="right">Commit Right (Local -> Remote Stage)</option>
            <option value="left">Commit Left (Remote Rebase -> Local)</option>
          </select>
        </div>
        <div class="form-group">
          <label><input type="checkbox" id="commit-amend" style="width: auto;"> Amend Previous Commit</label>
        </div>
        <div class="form-group">
          <label><input type="checkbox" id="commit-push" checked style="width: auto;"> Auto Push Immediately</label>
        </div>
        <button class="btn" onclick="execCommitin()">Execute Commitin</button>
        <div id="commit-output" style="margin-top: 1rem; font-family: monospace; white-space: pre-wrap; font-size: 0.85rem; color: var(--muted);"></div>
      </div>
    </div>

    <!-- SSH TAB -->
    <div id="tab-ssh" class="content-area">
      <div class="card">
        <h3>SSH Cluster Nodes & Fleet Operations</h3>
        <button class="btn" onclick="refreshNodes()">Refresh Node Fleet</button>
        <button class="btn btn-secondary" onclick="deployMacroFleetModal()">Deploy Macros to Fleet</button>
        <table id="nodes-table" style="margin-top: 1rem;">
          <thead><tr><th>Node ID</th><th>Alias</th><th>Host / IP</th><th>OS</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody id="nodes-body"><tr><td colspan="6" style="text-align:center;">Loading nodes...</td></tr></tbody>
        </table>
      </div>
    </div>

    <!-- MACRO TAB -->
    <div id="tab-macro" class="content-area">
      <div class="card">
        <h3>Macro Automation Builder</h3>
        <button class="btn" onclick="loadMacros()">Reload Macros</button>
        <div id="macros-list" style="margin-top: 1rem;"></div>
      </div>
    </div>

    <!-- INSTALLER TAB -->
    <div id="tab-installer" class="content-area">
      <div class="card">
        <h3>Custom Multi-OS Installer Manager</h3>
        <div class="form-group"><label>Installer Name</label><input type="text" id="inst-name" placeholder="e.g. agy-tools"></div>
        <div class="form-group"><label>Target OS</label><select id="inst-os"><option value="windows">Windows</option><option value="ubuntu">Ubuntu / Debian</option><option value="centos">CentOS / RHEL</option><option value="unix">Universal Unix</option></select></div>
        <div class="form-group"><label>Install Command</label><input type="text" id="inst-cmd" placeholder="curl -fsSL ... | sh"></div>
        <div class="form-group"><label>Target Node for Immediate Test</label><select id="inst-node-select"><option value="local">Local Machine</option></select></div>
        <button class="btn" onclick="saveInstaller()">Save Installer</button>
        <button class="btn btn-secondary" onclick="testInstaller()">Test on Selected Node</button>
      </div>
    </div>

    <!-- PROMPTS TAB -->
    <div id="tab-prompts" class="content-area">
      <div class="card">
        <h3>Prompts & AI Instructions Manager</h3>
        <div class="form-group"><label>Prompt Template</label><textarea id="prompt-content" rows="10" placeholder="Paste or edit prompt markdown/JSON..."></textarea></div>
        <button class="btn" onclick="formatPromptJSON()">Format JSON</button>
        <button class="btn btn-secondary" onclick="importPrompt()">Import Prompt</button>
      </div>
    </div>

    <!-- IMPORT-EXPORT TAB -->
    <div id="tab-import-export" class="content-area">
      <div class="card">
        <h3>System Configuration Import / Export</h3>
        <button class="btn" onclick="exportFullConfig()">Export Full JSON Config</button>
        <button class="btn btn-secondary" onclick="importFullConfig()">Import JSON Config</button>
        <div style="margin-top: 1rem;"><textarea id="import-export-area" rows="12" placeholder="JSON snapshot data..."></textarea></div>
      </div>
    </div>

    <!-- SCHEDULES TAB -->
    <div id="tab-schedules" class="content-area">
      <div class="card">
        <h3>Task & Cron Schedules</h3>
        <div class="form-group"><label>Schedule Name</label><input type="text" id="sched-name" placeholder="e.g. night-sync"></div>
        <div class="form-group"><label>Cron Expression</label><input type="text" id="sched-cron" placeholder="0 2 * * *"></div>
        <button class="btn" onclick="addSchedule()">Add Schedule</button>
      </div>
    </div>

    <!-- REMOTE EDITOR TAB -->
    <div id="tab-editor" class="content-area">
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
    function showTab(tabId) {
      document.querySelectorAll('.content-area').forEach(el => el.classList.remove('active'));
      document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));
      const target = document.getElementById('tab-' + tabId);
      if (target) target.classList.add('active');
      document.getElementById('header-title').innerText = tabId.charAt(0).toUpperCase() + tabId.slice(1);
      event?.target?.classList.add('active');
      if (tabId === 'ssh') refreshNodes();
      if (tabId === 'editor') populateEditorNodes();
    }

    async function refreshNodes() {
      try {
        const res = await fetch('/api/ssh/nodes');
        const data = await res.json();
        const tbody = document.getElementById('nodes-body');
        tbody.innerHTML = '';
        if (!data || data.length === 0) {
          tbody.innerHTML = '<tr><td colspan="6" style="text-align:center">No nodes registered. Use "gitmap ssh join" to add nodes.</td></tr>';
          return;
        }
        data.forEach(n => {
          const row = document.createElement('tr');
          row.innerHTML = ` + "`" + `<td>${n.nodeId || '-'}</td><td><b>${n.alias}</b></td><td>${n.host}</td><td>${n.os || 'Linux'}</td><td><span class="badge ${n.isOnline ? 'badge-online' : 'badge-offline'}">${n.isOnline ? 'Online' : 'Offline'}</span></td><td><button class="btn" style="padding:0.25rem 0.5rem;font-size:0.75rem" onclick="testNode('${n.alias}')">Ping</button></td>` + "`" + `;
          tbody.appendChild(row);
        });
      } catch (e) { console.error(e); }
    }

    async function populateEditorNodes() {
      try {
        const res = await fetch('/api/ssh/nodes');
        const data = await res.json();
        const sel = document.getElementById('editor-node');
        sel.innerHTML = '<option value="local">local (current machine)</option>';
        if (data && data.length) {
          data.forEach(n => {
            const opt = document.createElement('option');
            opt.value = n.alias;
            opt.innerText = n.alias + ' (' + n.host + ')';
            sel.appendChild(opt);
          });
        }
      } catch (e) {}
    }

    async function openRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      if (!path) return alert('Please enter a file path');
      document.getElementById('editor-status').innerText = 'Loading ' + path + '...';
      try {
        const res = await fetch('/api/editor/read', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({nodeAlias: node, filePath: path})
        });
        const data = await res.json();
        if (data.isSuccess) {
          document.getElementById('editor-area').value = data.content;
          document.getElementById('editor-status').innerText = 'Loaded ' + path + ' from ' + node;
          document.getElementById('editor-lang').innerText = 'Language: ' + (data.language || 'plaintext');
        } else {
          alert('Failed to read file: ' + (data.error || 'unknown error'));
        }
      } catch (e) { alert('Error reading file: ' + e.message); }
    }

    async function saveRemoteFile() {
      const node = document.getElementById('editor-node').value;
      const path = document.getElementById('editor-path').value.trim();
      const content = document.getElementById('editor-area').value;
      if (!path) return alert('Please enter a file path');
      document.getElementById('editor-status').innerText = 'Saving to ' + node + '...';
      try {
        const res = await fetch('/api/editor/save', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({nodeAlias: node, filePath: path, content: content})
        });
        const data = await res.json();
        if (data.success) {
          document.getElementById('editor-status').innerText = 'Saved successfully at ' + new Date().toLocaleTimeString();
        } else {
          alert('Failed to save file: ' + (data.error || 'unknown error'));
        }
      } catch (e) { alert('Error saving file: ' + e.message); }
    }

    async function execCommitin() {
      const msg = document.getElementById('commit-msg').value;
      const dir = document.getElementById('commit-direction').value;
      const isAmend = document.getElementById('commit-amend').checked;
      const isPush = document.getElementById('commit-push').checked;
      const out = document.getElementById('commit-output');
      out.innerText = 'Executing commitin...';
      try {
        const res = await fetch('/api/commitin/exec', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify({message: msg, direction: dir, isAmend: isAmend, isPush: isPush})
        });
        const data = await res.json();
        out.innerText = data.output || data.error || 'Done';
      } catch (e) { out.innerText = 'Error: ' + e.message; }
    }

    async function loadSettings() {
      try {
        const res = await fetch('/api/settings');
        const data = await res.json();
        if (data) {
          if (data.theme && document.getElementById('setting-theme')) document.getElementById('setting-theme').value = data.theme;
          if (data.defaultRemote && document.getElementById('setting-remote')) document.getElementById('setting-remote').value = data.defaultRemote;
          if (data.clusterPort && document.getElementById('setting-port')) document.getElementById('setting-port').value = data.clusterPort;
          if (data.graphicsMode && document.getElementById('setting-graphics')) document.getElementById('setting-graphics').value = data.graphicsMode;
          if (data.autoOpenBrowser !== undefined && document.getElementById('setting-auto-open')) document.getElementById('setting-auto-open').value = String(data.autoOpenBrowser);
          if (data.commitInLayout && document.getElementById('setting-commitin-layout')) document.getElementById('setting-commitin-layout').value = data.commitInLayout;
          if (data.pullDirection && document.getElementById('setting-pull-direction')) document.getElementById('setting-pull-direction').value = data.pullDirection;
          if (data.prReplayMode && document.getElementById('setting-pr-mode')) document.getElementById('setting-pr-mode').value = data.prReplayMode;
        }
      } catch (e) { console.error('Failed to load settings', e); }
    }

    async function saveSettings() {
      const payload = {
        theme: document.getElementById('setting-theme').value,
        defaultRemote: document.getElementById('setting-remote').value,
        clusterPort: parseInt(document.getElementById('setting-port').value) || 49152,
        graphicsMode: document.getElementById('setting-graphics') ? document.getElementById('setting-graphics').value : 'high',
        autoOpenBrowser: document.getElementById('setting-auto-open') ? document.getElementById('setting-auto-open').value === 'true' : true,
        commitInLayout: document.getElementById('setting-commitin-layout') ? document.getElementById('setting-commitin-layout').value : 'split',
        pullDirection: document.getElementById('setting-pull-direction') ? document.getElementById('setting-pull-direction').value : 'pull-left',
        prReplayMode: document.getElementById('setting-pr-mode') ? document.getElementById('setting-pr-mode').value : 'merges'
      };
      try {
        const res = await fetch('/api/settings', {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: JSON.stringify(payload)
        });
        const data = await res.json();
        if (data.success) {
          alert('Settings saved successfully!');
        }
      } catch (e) { alert('Failed to save settings: ' + e.message); }
    }

    // Auto-detect route on load
    window.addEventListener('load', () => {
      loadSettings();
      const path = window.location.pathname.replace(/^\//, '');
      if (path && document.getElementById('tab-' + path)) {
        showTab(path);
      }
    });
  </script>
</body>
</html>
`

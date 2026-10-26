package cmdui

const uiAssetsDocHead = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>GitMap Management & Fleet Dashboard</title>
`

// uiAssetsStyle moved to ui_assets_theme.go (program 246) — extended there with
// the real light theme, toast CSS, and responsive shell styles.

const uiAssetsHeadClose = `</head>
`

const uiAssetsTail = `</body>
</html>
`

const uiAssetsMarkupShell = `<body>
  <a class="skip-link" href="#main-content">Skip to main content</a>
  <div id="sidebar">
    <div class="logo"><svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg><span class="nav-label">GitMap Fleet UI</span></div>
    <nav id="nav-list" role="tablist" aria-label="Primary"></nav>
  </div>
  <div id="sidebar-scrim" class="sidebar-scrim" tabindex="0" role="button" aria-label="Close navigation" onclick="toggleSidebar()" onkeydown="if(event.key==='Enter'||event.key==='Escape'){event.preventDefault();toggleSidebar();}"></div>
  <div id="main-content" tabindex="-1">
    <header>
      <div style="display:flex;align-items:center;gap:var(--sp-3);">
        <button class="hamburger" onclick="toggleSidebar()" aria-label="Toggle navigation" aria-expanded="false"><svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><path d="M3 6h18M3 12h18M3 18h18"/></svg></button>
        <h2 id="header-title">Settings</h2>
      </div>
      <div style="display:flex;align-items:center;gap:var(--sp-2);">
        <div id="status-indicator"><span class="status-dot" aria-hidden="true"></span><span class="badge badge-online">Connected</span></div>
        <button id="btn-fullscreen" class="btn" style="padding:4px 10px;font-size:var(--fs-xs);background:var(--raised);border:1px solid var(--border);" title="Toggle Fullscreen (F11)">⛶ Fullscreen</button>
        <button id="btn-exit" class="btn" style="padding:4px 10px;font-size:var(--fs-xs);background:var(--error);color:#fff;" title="Close application window and stop server">✕ Close Window</button>
      </div>
    </header>

    <!-- SETTINGS TAB -->
`

const uiAssetsMarkupMisc = `    <div id="tab-editor" class="content-area">
      <div class="card" style="margin-bottom: 0.5rem;">
        <div style="display: flex; gap: 10px; align-items: flex-end;">
          <div style="flex: 1;"><label for="editor-node" style="font-size:0.8rem;color:var(--muted)">Target Node</label><select id="editor-node"><option value="local">local (current machine)</option></select></div>
          <div style="flex: 3;"><label for="editor-path" style="font-size:0.8rem;color:var(--muted)">Absolute / Relative File Path</label><input type="text" id="editor-path" placeholder="/etc/gitmap.conf or cli/main.go"></div>
          <button class="btn" onclick="openRemoteFile()">Load File</button>
          <button class="btn" style="background:var(--success)" onclick="saveRemoteFile()">Save to Host</button>
        </div>
      </div>
      <div id="editor-container">
        <div id="editor-toolbar"><span id="editor-status" style="font-size:0.85rem;color:var(--muted)">No file loaded</span><span id="editor-lang" style="font-size:0.8rem;color:var(--primary)">Language: Plaintext</span></div>
        <label for="editor-area" class="visually-hidden">Remote file content</label><textarea id="editor-area" spellcheck="false" placeholder="Remote file content will appear here..."></textarea>
      </div>
    </div>

    <!-- CLI & AGY HELP TAB -->
    <div id="tab-help" class="content-area">
      <div class="card card-hero">
        <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-2); flex-wrap:wrap; gap:var(--sp-2);">
          <h3 style="margin-bottom:0;"><span data-icon="book-open"></span> Embedded GitMap Manuals &amp; CLI Reference</h3>
          <span class="badge badge-online" id="help-topics-count">100+ Topics</span>
        </div>
        <p style="font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-3);">
          Search and view official markdown manuals and command documentation embedded directly in the binary.
        </p>
        <div class="form-group" style="margin-bottom: var(--sp-3);">
          <input type="text" id="help-search" placeholder="Search 100+ embedded GitMap manuals and help topics..." oninput="filterHelpTopicsUI(this.value)">
        </div>
        <div style="display: grid; grid-template-columns: 280px 1fr; gap: var(--sp-3); min-height: 460px;" id="help-viewer-layout">
          <div style="display: flex; flex-direction: column; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--raised); overflow: hidden;">
            <div style="padding: var(--sp-2) var(--sp-3); border-bottom: 1px solid var(--border); font-size: var(--fs-xs); font-weight: 600; color: var(--muted); text-transform: uppercase;">
              Topics (<span id="help-visible-count">...</span>)
            </div>
            <div id="help-topics-list" style="flex: 1; max-height: 500px; overflow-y: auto; padding: var(--sp-1);">
              <div style="padding: var(--sp-2); font-size: var(--fs-sm); color: var(--muted);">Loading topics...</div>
            </div>
          </div>
          <div style="display: flex; flex-direction: column; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--card); overflow: hidden;">
            <div style="padding: var(--sp-2) var(--sp-3); border-bottom: 1px solid var(--border); display: flex; justify-content: space-between; align-items: center; background: var(--raised);">
              <span id="help-active-topic-title" style="font-size: var(--fs-sm); font-weight: 600; color: var(--primary);">Select a topic</span>
              <button class="btn btn-secondary" onclick="copyCurrentHelpDoc()" style="padding: 2px 8px; font-size: var(--fs-xs);"><span data-icon="check"></span>Copy Markdown</button>
            </div>
            <div id="help-doc-viewer" tabindex="0" aria-label="Help document content" style="flex: 1; padding: var(--sp-4); overflow-y: auto; max-height: 500px; font-family: var(--font-mono); font-size: 0.85rem; line-height: 1.6; white-space: pre-wrap; color: var(--text); background: var(--code-bg);">Select a command topic from the list to view its official manual.</div>
          </div>
        </div>
      </div>

      <div class="card">
        <h3><span data-icon="book-open"></span> CLI &amp; AGY Help — Speed Commands, LAP, Rerun-With-ID, Machine Identity &amp; Bot Reference</h3>
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

  <div id="toast-stack" aria-live="polite"></div>

  <script>
`

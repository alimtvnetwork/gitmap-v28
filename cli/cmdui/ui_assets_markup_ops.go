package cmdui

const uiAssetsMarkupOps = `    <div id="tab-commitin" class="content-area">
      <div class="card card-hero">
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
      <div class="card card-hero">
        <h3>Public Key &amp; Firewall Management</h3>
        <p style="font-size: 0.85rem; color: var(--muted); margin-bottom: 0.75rem;">
          Quick actions for node authentication keys, listening ports, and firewall rules.
        </p>
        <div style="display: flex; gap: 8px; flex-wrap: wrap;">
          <button class="btn" onclick="viewSSHPublicKeyUI()">View SSH Public Key</button>
          <button class="btn btn-secondary" onclick="managePortsAndFirewallUI()">Manage Ports &amp; Firewall</button>
        </div>
        <div id="ssh-quick-card-output" style="margin-top: 0.75rem; display: none;">
          <textarea id="ssh-quick-display" rows="5" readonly style="font-family: monospace; font-size: 0.85rem;"></textarea>
        </div>
      </div>

      <div class="card">
        <h3>SSH Cluster Nodes & Fleet Operations</h3>
        <div style="display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 1rem;">
          <button class="btn" onclick="refreshNodes()">Refresh Node Fleet</button>
          <button class="btn btn-secondary" onclick="deployKeysFleetUI()"><span data-icon="rocket"></span>Deploy Keys (All)</button>
          <button class="btn btn-secondary" onclick="exportSSHNodesUI()"><span data-icon="arrow-left-right"></span>Export Nodes JSON</button>
          <button class="btn btn-secondary" onclick="openImportNodesModal()"><span data-icon="arrow-left-right"></span>Import Nodes JSON</button>
          <button class="btn btn-secondary" onclick="deployMacroFleetModal()">Deploy Macros to Fleet</button>
        </div>
        <div id="ssh-action-status" style="font-size: 0.85rem; margin-bottom: 0.5rem; display: none;"></div>
        <table id="nodes-table" style="margin-top: 0.5rem;">
          <thead><tr><th>Node ID</th><th>Alias</th><th>Host / IP</th><th>OS</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody id="nodes-body"><tr><td colspan="6" style="text-align:center;">Loading nodes...</td></tr></tbody>
        </table>
        <div style="margin-top: 1.25rem; padding: 1rem; background: var(--bg); border: 1px solid var(--border); border-radius: 6px;">
          <h4 style="margin-bottom: 0.5rem; color: var(--primary);"><span data-icon="info"></span> Fleet CLI Commands &amp; Cross-Machine Sharing</h4>
          <p style="font-size: 0.85rem; color: var(--muted); line-height: 1.6; margin-bottom: 0.5rem;">
            • <b>Deploy Keys Across Fleet:</b> <code>gitmap deploy keys all</code> or <code>gitmap deploy-keys-all</code> (alias: <code>gitmap ssh deploy keys all</code>)<br>
            • <b>Export Nodes to JSON:</b> <code>gitmap ssh export-json [file]</code> or <code>gitmap export-ssh</code><br>
            • <b>Import Nodes from JSON:</b> <code>gitmap ssh import-json [file]</code> or <code>gitmap deploy import [file]</code><br>
            • <b>Single-Line Compact Share:</b> <code>gitmap ssh export-oneliner</code> (copies base64 command to import on any machine)<br>
            • <b>Deploy Cluster Topology:</b> <code>gitmap deploy config ssh [all]</code> (syncs all nodes and IPs across the fleet)
          </p>
          <p style="font-size: 0.8rem; color: var(--accent);">
            Tip: You can also manage full configuration backups in the <b><span data-icon="arrow-left-right"></span> Import / Export</b> tab.
          </p>
        </div>
      </div>
    </div>

    <!-- MACRO TAB -->
    <div id="tab-macro" class="content-area">
      <div class="card card-hero">
        <h3>Macro Automation Builder</h3>
        <button class="btn" onclick="loadMacros()">Reload Macros</button>
        <div id="macros-list" style="margin-top: 1rem;"></div>
      </div>
    </div>

    <!-- INSTALLER TAB -->
    <div id="tab-installer" class="content-area">
      <div class="card card-hero">
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
      <div class="card card-hero">
        <h3>Prompts & AI Instructions Manager</h3>
        <div class="form-group"><label>Prompt Template</label><textarea id="prompt-content" rows="10" placeholder="Paste or edit prompt markdown/JSON..."></textarea></div>
        <button class="btn" onclick="formatPromptJSON()">Format JSON</button>
        <button class="btn btn-secondary" onclick="importPrompt()">Import Prompt</button>
      </div>
    </div>

    <!-- IMPORT-EXPORT TAB -->
    <div id="tab-import-export" class="content-area">
      <div class="card card-hero">
        <h3>System Configuration Import / Export</h3>
        <button class="btn" onclick="exportFullConfig()">Export Full JSON Config</button>
        <button class="btn btn-secondary" onclick="importFullConfig()">Import JSON Config</button>
        <div style="margin-top: 1rem;"><textarea id="import-export-area" rows="12" placeholder="JSON snapshot data..."></textarea></div>
      </div>
      <div class="card" style="margin-top: 1rem;">
        <h3>SSH Fleet Nodes &amp; Public Keys Export / Import</h3>
        <p style="font-size:0.85rem;color:var(--muted);margin-bottom:0.75rem;">
          Export all registered cluster nodes, aliases, IPs, and public keys to share across machines, or import nodes directly from another node's JSON export.
        </p>
        <div style="display:flex;gap:0.5rem;flex-wrap:wrap;margin-bottom:0.75rem;">
          <button class="btn" onclick="exportSSHNodesUI()"><span data-icon="arrow-left-right"></span>Export Fleet Nodes JSON</button>
          <button class="btn btn-secondary" onclick="importSSHNodesUI()"><span data-icon="arrow-left-right"></span>Import Fleet Nodes JSON</button>
          <button class="btn" style="background:#0284c7" onclick="deployKeysFleetUI()"><span data-icon="rocket"></span>Deploy Keys (All Nodes)</button>
        </div>
        <textarea id="ssh-nodes-export-area" rows="8" placeholder="Fleet nodes JSON export/import data will appear here..."></textarea>
      </div>
    </div>

    <!-- SCHEDULES TAB -->
    <div id="tab-schedules" class="content-area">
      <div class="card card-hero">
        <h3>Task & Cron Schedules</h3>
        <div class="form-group"><label>Schedule Name</label><input type="text" id="sched-name" placeholder="e.g. night-sync"></div>
        <div class="form-group"><label>Cron Expression</label><input type="text" id="sched-cron" placeholder="0 2 * * *"></div>
        <button class="btn" onclick="addSchedule()">Add Schedule</button>
      </div>
    </div>

    <!-- REMOTE EDITOR TAB -->
`

const uiAssetsMarkupSettingsA = `    <div id="tab-settings" class="content-area active">
      <div class="sub-tab-bar" style="display: flex; gap: 8px; margin-bottom: 1.25rem; border-bottom: 1px solid var(--border); padding-bottom: 0.75rem; flex-wrap: wrap;">
        <button type="button" class="sub-tab-btn active" onclick="showSettingsSubTab('all')" id="subnav-all"><span data-icon="settings"></span>All Settings</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('general')" id="subnav-general"><span data-icon="pen-line"></span>General &amp; UI</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('cursor')" id="subnav-cursor"><span data-icon="rocket"></span>Cursor Fleet &amp; Sync</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('speed')" id="subnav-speed"><span data-icon="clock"></span>Speed &amp; Automation</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('alerts')" id="subnav-alerts"><span data-icon="alert-triangle"></span>Alerts &amp; Notifications</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('instances')" id="subnav-instances"><span data-icon="server"></span>Antigravity Instances</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('catalog')" id="subnav-catalog"><span data-icon="book-open"></span>UI/UX Repo Index</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('terminal')" id="subnav-terminal"><span data-icon="scroll-text"></span>Embedded Terminal</button>
      </div>

      <!-- 1. GENERAL & UI APPEARANCE -->
      <div class="card card-hero settings-subtab-pane" id="settings-pane-general">
        <h3><span data-icon="pen-line"></span> General &amp; UI Appearance</h3>
        <div class="form-group">
          <label>Theme</label>
          <select id="setting-theme">
            <option value="dark">Dark</option>
            <option value="light">Light</option>
          </select>
        </div>
        <div class="form-group">
          <label>Graphics Rendering &amp; Hardware Acceleration</label>
          <select id="setting-graphics">
            <option value="high">High Performance / Hardware Acceleration</option>
            <option value="fast">Fast Standard</option>
            <option value="plain">Plain / Low Graphics</option>
          </select>
        </div>
        <div class="form-group">
          <label>Auto-Open Browser on UI Launch</label>
          <select id="setting-auto-open">
            <option value="true">Enabled (Launch default browser automatically)</option>
            <option value="false">Disabled (Emit URL to console only)</option>
          </select>
        </div>
        <div class="form-group">
          <label>Cluster REST Port</label>
          <input type="number" id="setting-port" value="49152">
        </div>
        <div class="form-group">
          <label>Default Remote Target</label>
          <input type="text" id="setting-remote" placeholder="origin or alias">
        </div>
        <div class="form-group">
          <label>Pull Direction Orientation</label>
          <select id="setting-pull-direction">
            <option value="pull-left">Pull-Left (Replay from target down to left)</option>
            <option value="pull-right">Pull-Right (Replay from left source into right target)</option>
            <option value="bidirectional">Bidirectional Sync</option>
          </select>
        </div>
        <div class="form-group">
          <label>Commit-Pull PR Replay Mode</label>
          <select id="setting-pr-mode">
            <option value="merges">merges (Feature branches &amp; release PRs)</option>
            <option value="feature-per-commit">feature-per-commit</option>
            <option value="direct">direct mainline</option>
          </select>
        </div>
        <div class="form-group">
          <label>Commit-In Migration Studio Layout</label>
          <select id="setting-commitin-layout">
            <option value="split">Split Dual-Pane (Config &amp; Output)</option>
            <option value="left">Left Focus (Config dominant)</option>
            <option value="right">Right Focus (Output dominant)</option>
            <option value="stacked">Stacked Vertical</option>
          </select>
        </div>
      </div>

      <!-- 2. CURSOR FLEET & SYNC -->
      <div class="card settings-subtab-pane" id="settings-pane-cursor">
        <h3><span data-icon="rocket"></span> Cursor Fleet &amp; Sync</h3>
        <div class="form-group">
          <label>Machine Name (<code>machine.name</code>)</label>
          <input type="text" id="setting-machine-name" name="machine.name" data-key="machine.name" placeholder="dev-win-01 / ubuntu-node-02">
        </div>
        <div class="form-group">
          <label>Network Alias (<code>machine.alias</code> &mdash; auto-defaults to Local IPv4)</label>
          <input type="text" id="setting-machine-alias" name="machine.alias" data-key="machine.alias" placeholder="Auto-defaults to Local IPv4 if unset">
        </div>
        <div class="form-group">
          <label>Special Secrets Repository Name (<code>special_repos.secrets_name</code> &mdash; shortcut: <code>gitmap rs</code> / <code>gitmap cd rs</code>)</label>
          <input type="text" id="setting-special-secrets-name" name="special_repos.secrets_name" data-key="special_repos.secrets_name" value="repo-secrets" placeholder="repo-secrets">
        </div>
        <div class="form-group">
          <label>Special Cache / Storage Repository Name (<code>special_repos.cache_name</code> &mdash; shortcut: <code>gitmap rc</code> / <code>gitmap cd rc</code>)</label>
          <input type="text" id="setting-special-cache-name" name="special_repos.cache_name" data-key="special_repos.cache_name" value="repo-cache" placeholder="repo-cache">
        </div>
      </div>

      <!-- 3. SPEED & AUTOMATION -->
      <div class="card settings-subtab-pane" id="settings-pane-speed">
        <h3><span data-icon="clock"></span> Speed &amp; Automation</h3>
        <div class="form-group">
          <label>LAP Default Lookback Hours (N = 24) (<code>lap.default_hours</code>)</label>
          <input type="number" id="setting-lap-hours" name="lap.default_hours" data-key="lap.default_hours" value="24" placeholder="24">
        </div>
        <div class="form-group">
          <label>Account Switch Fast-Forward Threshold (%) (<code>account_switch.threshold</code> &mdash; default 15% prod / 98% E2E test)</label>
          <input type="number" id="setting-account-switch-threshold" name="account_switch.threshold" data-key="account_switch.threshold" value="15" placeholder="15 (prod) or 98 (E2E test)">
        </div>
      </div>

      <!-- 4. ALERTS & NOTIFICATIONS -->
      <div class="card settings-subtab-pane" id="settings-pane-alerts">
        <h3><span data-icon="alert-triangle"></span> Alerts &amp; Notifications</h3>
        <div class="form-group">
          <label>Telegram Two-Way Bot Setup &mdash; Bot Token (<code>telegram.bot_token</code>)</label>
          <input type="text" id="setting-telegram-token" name="telegram.bot_token" data-key="telegram.bot_token" placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11">
        </div>
        <div class="form-group">
          <label>Telegram Two-Way Bot Setup &mdash; Chat ID (<code>telegram.chat_id</code>)</label>
          <input type="text" id="setting-telegram-chat" name="telegram.chat_id" data-key="telegram.chat_id" placeholder="-1001234567890">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup &mdash; SMTP Host (<code>email.smtp_host</code>)</label>
          <input type="text" id="setting-email-smtp" name="email.smtp_host" data-key="email.smtp_host" placeholder="smtp.gmail.com:587">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup &mdash; From Address (<code>email.from</code>)</label>
          <input type="email" id="setting-email-from" name="email.from" data-key="email.from" placeholder="bot@example.com">
        </div>
        <div class="form-group">
          <label>Email Speed SMTP Setup &mdash; To Address (<code>email.to</code>)</label>
          <input type="email" id="setting-email-to" name="email.to" data-key="email.to" placeholder="dev@example.com">
        </div>
      </div>

`

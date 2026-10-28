package cmdui

const uiAssetsMarkupOps = `    <div id="tab-commitin" class="content-area">
      <div class="card card-hero">
        <h3>Commitin Multi-Option Engine</h3>
        <div class="form-group">
          <label for="commit-msg">Commit Message</label>
          <input type="text" id="commit-msg" placeholder="feat(scope): concise description">
        </div>
        <div class="form-group">
          <label for="commit-direction">Commit Direction / Flow</label>
          <select id="commit-direction">
            <option value="standard">Standard Commit</option>
            <option value="right">Commit Right (Local -> Remote Stage)</option>
            <option value="left">Commit Left (Remote Rebase -> Local)</option>
          </select>
        </div>
        <div class="form-group">
          <label for="commit-amend"><input type="checkbox" id="commit-amend" style="width: auto;"> Amend Previous Commit</label>
        </div>
        <div class="form-group">
          <label for="commit-push"><input type="checkbox" id="commit-push" checked style="width: auto;"> Auto Push Immediately</label>
        </div>
        <button class="btn" onclick="execCommitin()">Execute Commitin</button>
        <div id="commit-output" style="margin-top: var(--sp-4); font-family: monospace; white-space: pre-wrap; font-size: var(--fs-sm); color: var(--muted);"></div>
      </div>
    </div>

    <!-- SSH TAB -->
    <div id="tab-ssh" class="content-area">
      <div class="card card-hero">
        <h3>Public Key &amp; Firewall Management</h3>
        <p style="font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-3);">
          Quick actions for node authentication keys, listening ports, and firewall rules.
        </p>
        <div style="display: flex; gap: var(--sp-2); flex-wrap: wrap;">
          <button class="btn" onclick="viewSSHPublicKeyUI()">View SSH Public Key</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="managePortsAndFirewallUI()">Manage Ports &amp; Firewall</button>
        </div>
        <div id="ssh-quick-card-output" style="margin-top: var(--sp-3); display: none;">
          <textarea id="ssh-quick-display" rows="5" readonly aria-label="SSH key output" style="font-family: monospace; font-size: 0.85rem;"></textarea>
        </div>
      </div>

      <div class="card">
        <h3>SSH Cluster Nodes & Fleet Operations</h3>
        <div style="display: flex; gap: var(--sp-2); flex-wrap: wrap; margin-bottom: var(--sp-4);">
          <button class="btn" onclick="refreshNodes()">Refresh Node Fleet</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="deployKeysFleetUI()"><span data-icon="rocket"></span>Deploy Keys (All)</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="exportSSHNodesUI()"><span data-icon="arrow-left-right"></span>Export Nodes JSON</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="openImportNodesModal()"><span data-icon="arrow-left-right"></span>Import Nodes JSON</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="deployMacroFleetModal()">Deploy Macros to Fleet</button>
        </div>
        <div id="ssh-action-status" style="font-size: var(--fs-sm); margin-bottom: var(--sp-2); display: none;"></div>
        <table id="nodes-table" style="margin-top: var(--sp-2);">
          <thead><tr><th>Node ID</th><th>Alias</th><th>Host / IP</th><th>OS</th><th>Account</th><th>Instances</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody id="nodes-body"><tr><td colspan="8" style="text-align:center;">Loading nodes...</td></tr></tbody>
        </table>
        <div style="margin-top: var(--sp-4); padding: var(--sp-4); background: var(--bg); border: 1px solid var(--border); border-radius: var(--radius-sm);">
          <h4 style="margin-bottom: var(--sp-2); color: var(--primary); font-size: var(--fs-sm);"><span data-icon="info"></span> Fleet CLI Commands &amp; Cross-Machine Sharing</h4>
          <p style="font-size: var(--fs-sm); color: var(--muted); line-height: 1.6; margin-bottom: var(--sp-2);">
            • <b>Fleet Pull All:</b> <code>gitmap pa --ssh</code> (Enqueue pull-all across SSH fleet)<br>
            • <b>Inspect Cluster Nodes:</b> <code>gitmap nodes</code> (Fleet node discovery, status &amp; latency)<br>
            • <b>Cluster Ping Test:</b> <code>gitmap nodes ping</code> (Dual-stack ICMP/TCP ping probe)<br>
            • <b>Pipeline Diagnostics:</b> <code>gitmap pe</code> (Pipeline error analyzer &amp; runner log extraction)<br>
            • <b>Deploy Keys Across Fleet:</b> <code>gitmap deploy keys all</code> or <code>gitmap deploy-keys-all</code> (alias: <code>gitmap ssh deploy keys all</code>)<br>
            • <b>Export Nodes to JSON:</b> <code>gitmap ssh export-json [file]</code> or <code>gitmap export-ssh</code><br>
            • <b>Import Nodes from JSON:</b> <code>gitmap ssh import-json [file]</code> or <code>gitmap deploy import [file]</code><br>
            • <b>Single-Line Compact Share:</b> <code>gitmap ssh export-oneliner</code> (copies base64 command to import on any machine)<br>
            • <b>Deploy Cluster Topology:</b> <code>gitmap deploy config ssh [all]</code> (syncs all nodes and IPs across the fleet)
          </p>
          <p style="font-size: var(--fs-xs); color: var(--accent);">
            Tip: You can also manage full configuration backups in the <b><span data-icon="arrow-left-right"></span> Import / Export</b> tab.
          </p>
        </div>
      </div>
    </div>

    <!-- MACRO TAB -->
    <div id="tab-macro" class="content-area">
      <div class="card card-hero">
        <h3><span data-icon="scroll-text"></span> Macro Automation Builder</h3>
        <p style="font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-3);">
          Inspect, create, test, and replay saved CLI workflow macros across all workspaces.
        </p>
        <div style="display: flex; gap: var(--sp-2); margin-bottom: var(--sp-3);">
          <button class="btn" onclick="loadMacros()"><span data-icon="clock"></span> Reload Macros</button>
          <button class="btn btn-secondary" onclick="openNewMacroDialog()"><span data-icon="plus"></span> New Macro</button>
        </div>
        <div id="macros-list" style="margin-top: 1rem; display: flex; flex-direction: column; gap: var(--sp-2);"></div>
      </div>
    </div>

    <!-- INSTALLER TAB -->
    <div id="tab-installer" class="content-area">
      <div class="card card-hero">
        <h3><span data-icon="rocket"></span> Installer &amp; Workstation Profile Tree</h3>
        <p style="font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-3);">
          Explore workstation profiles, install/uninstall profiles, run single tool installations, and inspect dependency hierarchies.
        </p>
        <div style="display: flex; gap: var(--sp-2); margin-bottom: var(--sp-3);">
          <button class="btn" onclick="loadInstallerProfiles()"><span data-icon="clock"></span> Reload Profiles</button>
        </div>
        <div id="profiles-tree-list" style="margin-top: 1rem; display: flex; flex-direction: column; gap: var(--sp-3);"></div>
      </div>
      <div class="card" style="margin-top: var(--sp-4);">
        <h3>Custom Multi-OS Installer Definition</h3>
        <div class="form-group"><label for="inst-name">Installer Name / Slug</label><input type="text" id="inst-name" placeholder="e.g. agy-tools"></div>
        <div class="form-group"><label for="inst-os">Target OS</label><select id="inst-os"><option value="windows">Windows</option><option value="ubuntu">Ubuntu / Debian</option><option value="centos">CentOS / RHEL</option><option value="unix">Universal Unix</option></select></div>
        <div class="form-group"><label for="inst-cmd">Install Command</label><input type="text" id="inst-cmd" placeholder="curl -fsSL ... | sh"></div>
        <div class="form-group"><label for="inst-node-select">Target Node for Immediate Test</label><select id="inst-node-select"><option value="local">Local Machine</option></select></div>
        <div style="display: flex; gap: var(--sp-2); margin-top: var(--sp-2);">
          <button class="btn" onclick="saveInstaller()">Save Installer</button>
          <button class="btn btn-secondary" onclick="testInstaller()">Test on Selected Node</button>
        </div>
      </div>
    </div>

    <!-- PROMPTS TAB -->
    <div id="tab-prompts" class="content-area">
      <div class="card card-hero">
        <h3><span data-icon="message-square"></span> Prompts &amp; AI Instructions Manager</h3>
        <p style="font-size: var(--fs-sm); color: var(--muted); margin-bottom: var(--sp-3);">
          Search, filter, and inspect canonical prompts, agent workflows, and queued instructions across Antigravity instances.
        </p>
        <div class="form-group" style="margin-bottom: var(--sp-3);">
          <input type="text" id="prompt-search" placeholder="Search prompts, templates, and agent instructions..." oninput="filterPromptsUI(this.value)">
        </div>
        <div class="chip-container" id="prompt-category-pills" style="display: flex; gap: var(--sp-2); flex-wrap: wrap; margin-bottom: var(--sp-3);">
          <button type="button" class="chip active" onclick="setPromptCategory('all', this)">All</button>
          <button type="button" class="chip" onclick="setPromptCategory('active', this)">Active / Queued</button>
          <button type="button" class="chip" onclick="setPromptCategory('completed', this)">Completed</button>
          <button type="button" class="chip" onclick="loadPromptsList()" style="margin-left: auto; background: var(--raised);"><span data-icon="clock"></span>Refresh</button>
        </div>
        <div id="prompts-list" style="margin-top: var(--sp-2); display: flex; flex-direction: column; gap: var(--sp-2);">
          <div style="color: var(--muted); font-size: var(--fs-sm);">Loading prompts from Antigravity instances...</div>
        </div>
      </div>
      <div class="card">
        <h3>Prompt Template Scratchpad</h3>
        <div class="form-group"><label for="prompt-content">Prompt Template</label><textarea id="prompt-content" rows="8" placeholder="Paste or edit prompt markdown/JSON..."></textarea></div>
        <div style="display: flex; gap: var(--sp-2); flex-wrap: wrap;">
          <button class="btn" onclick="formatPromptJSON()">Format JSON</button>
          <button class="btn btn-secondary" onclick="importPrompt()">Import Prompt</button>
          <button class="btn btn-secondary" onclick="copyPromptToClipboard()"><span data-icon="check"></span>Copy to Clipboard</button>
        </div>
      </div>
    </div>

    <!-- IMPORT-EXPORT TAB -->
    <div id="tab-import-export" class="content-area">
      <div class="card card-hero">
        <h3>System Configuration Import / Export</h3>
        <button class="btn" onclick="exportFullConfig()">Export Full JSON Config</button>
        <button class="btn btn-secondary" onclick="importFullConfig()">Import JSON Config</button>
        <div style="margin-top: 1rem;"><textarea id="import-export-area" rows="12" aria-label="Full config JSON data" placeholder="JSON snapshot data..."></textarea></div>
      </div>
      <div class="card" style="margin-top: var(--sp-4);">
        <h3>SSH Fleet Nodes &amp; Public Keys Export / Import</h3>
        <p style="font-size:var(--fs-sm);color:var(--muted);margin-bottom:var(--sp-3);">
          Export all registered cluster nodes, aliases, IPs, and public keys to share across machines, or import nodes directly from another node's JSON export.
        </p>
        <div style="display:flex;gap:var(--sp-2);flex-wrap:wrap;margin-bottom:var(--sp-3);">
          <button class="btn" onclick="exportSSHNodesUI()"><span data-icon="arrow-left-right"></span>Export Fleet Nodes JSON</button>
          <button class="btn btn-secondary" style="margin-left: 0;" onclick="importSSHNodesUI()"><span data-icon="arrow-left-right"></span>Import Fleet Nodes JSON</button>
          <button class="btn" style="background:var(--primary)" onclick="deployKeysFleetUI()"><span data-icon="rocket"></span>Deploy Keys (All Nodes)</button>
        </div>
        <textarea id="ssh-nodes-export-area" rows="8" aria-label="Fleet nodes JSON export/import data" placeholder="Fleet nodes JSON export/import data will appear here..."></textarea>
      </div>
    </div>

    <!-- SCHEDULES TAB -->
    <div id="tab-schedules" class="content-area">
      <div class="card card-hero">
        <h3>Task & Cron Schedules</h3>
        <div class="form-group"><label for="sched-name">Schedule Name</label><input type="text" id="sched-name" placeholder="e.g. night-sync"></div>
        <div class="form-group"><label for="sched-cron">Cron Expression</label><input type="text" id="sched-cron" placeholder="0 2 * * *"></div>
        <button class="btn" onclick="addSchedule()">Add Schedule</button>
      </div>
    </div>

    <!-- TERMINAL TAB (Top-Level) -->
    <div id="tab-terminal" class="content-area">
      <div class="card card-hero">
        <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:var(--sp-3); flex-wrap:wrap; gap:var(--sp-2);">
          <div style="display:flex; align-items:center; gap:8px;">
            <span style="display:inline-flex; gap:5px; font-size:0.85rem;">
              <span style="color:#ef4444;">●</span>
              <span style="color:#f59e0b;">●</span>
              <span style="color:#10b981;">●</span>
            </span>
            <h3 style="margin-bottom:0; color:var(--text);"><span data-icon="scroll-text"></span> GitMap Live Terminal Runner</h3>
          </div>
          <div style="display:flex; align-items:center; gap:8px; flex:1; max-width:480px; justify-content:flex-end;">
            <select id="term-node" aria-label="Terminal target node" style="width:auto; min-width:110px; padding:0.4rem 0.6rem;">
              <option value="local">local</option>
            </select>
            <input type="text" id="term-cwd" aria-label="Terminal working directory" placeholder="cwd (optional)" style="flex:1; min-width:140px; padding:0.4rem 0.6rem;">
            <span id="term-badge" class="badge" style="display:none; white-space:nowrap;"></span>
          </div>
        </div>

        <!-- Quick command chips -->
        <div class="chip-container" style="display:flex; gap:var(--sp-2); flex-wrap:wrap; margin-bottom:var(--sp-3); align-items:center;">
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap status')">gitmap status</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap pe')">gitmap pe</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap nodes ls')">gitmap nodes ls</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap cursor settings view')">gitmap cursor settings view</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap lap 24')">gitmap lap 24</button>
          <button type="button" class="chip" onclick="clearTerminalOutput()" style="margin-left:auto; background:var(--raised); color:var(--muted);"><span data-icon="x"></span>Clear</button>
        </div>

        <!-- Terminal Console Output -->
        <div id="term-output" tabindex="0" aria-label="Terminal output" style="background:var(--terminal-bg, var(--code-bg)); border:1px solid var(--terminal-border, var(--border)); border-radius:var(--radius-sm); min-height:280px; max-height:480px; overflow-y:auto; padding:var(--sp-3); font-family:var(--font-mono); font-size:0.85rem; line-height:1.5; white-space:pre-wrap; color:var(--terminal-text, var(--code-text)); margin-bottom:0.75rem;">GitMap Fleet Terminal ready. Type a command or click a quick action above.</div>

        <!-- Command Prompt Input -->
        <div style="display:flex; gap:var(--sp-2); align-items:center;">
          <span style="font-family:monospace; color:var(--primary); font-weight:700; font-size:1.1rem;">$</span>
          <input type="text" id="term-input" placeholder="Enter command... (Press Enter or Run)" style="flex:1;">
          <button type="button" class="btn" id="term-run-btn" onclick="runTerminalCommand()"><span data-icon="rocket"></span>Run</button>
        </div>
      </div>
    </div>

    <!-- REMOTE EDITOR TAB -->

`

package cmdui

const uiAssetsMarkupOps = `    <div id="tab-commitin" class="content-area">
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
          <button class="btn btn-secondary" onclick="deployKeysFleetUI()">🔑 Deploy Keys (All)</button>
          <button class="btn btn-secondary" onclick="exportSSHNodesUI()">📤 Export Nodes JSON</button>
          <button class="btn btn-secondary" onclick="openImportNodesModal()">📥 Import Nodes JSON</button>
          <button class="btn btn-secondary" onclick="deployMacroFleetModal()">Deploy Macros to Fleet</button>
        </div>
        <div id="ssh-action-status" style="font-size: 0.85rem; margin-bottom: 0.5rem; display: none;"></div>
        <table id="nodes-table" style="margin-top: 0.5rem;">
          <thead><tr><th>Node ID</th><th>Alias</th><th>Host / IP</th><th>OS</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody id="nodes-body"><tr><td colspan="6" style="text-align:center;">Loading nodes...</td></tr></tbody>
        </table>
        <div style="margin-top: 1.25rem; padding: 1rem; background: var(--bg); border: 1px solid var(--border); border-radius: 6px;">
          <h4 style="margin-bottom: 0.5rem; color: var(--primary);">💡 Fleet CLI Commands & Cross-Machine Sharing</h4>
          <p style="font-size: 0.85rem; color: var(--muted); line-height: 1.6; margin-bottom: 0.5rem;">
            • <b>Deploy Keys Across Fleet:</b> <code>gitmap deploy keys all</code> or <code>gitmap deploy-keys-all</code> (alias: <code>gitmap ssh deploy keys all</code>)<br>
            • <b>Export Nodes to JSON:</b> <code>gitmap ssh export-json [file]</code> or <code>gitmap export-ssh</code><br>
            • <b>Import Nodes from JSON:</b> <code>gitmap ssh import-json [file]</code> or <code>gitmap deploy import [file]</code><br>
            • <b>Single-Line Compact Share:</b> <code>gitmap ssh export-oneliner</code> (copies base64 command to import on any machine)<br>
            • <b>Deploy Cluster Topology:</b> <code>gitmap deploy config ssh [all]</code> (syncs all nodes and IPs across the fleet)
          </p>
          <p style="font-size: 0.8rem; color: var(--accent);">
            Tip: You can also manage full configuration backups in the <b>🔄 Import / Export</b> tab.
          </p>
        </div>
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
      <div class="card" style="margin-top: 1rem;">
        <h3>SSH Fleet Nodes &amp; Public Keys Export / Import</h3>
        <p style="font-size:0.85rem;color:var(--muted);margin-bottom:0.75rem;">
          Export all registered cluster nodes, aliases, IPs, and public keys to share across machines, or import nodes directly from another node's JSON export.
        </p>
        <div style="display:flex;gap:0.5rem;flex-wrap:wrap;margin-bottom:0.75rem;">
          <button class="btn" onclick="exportSSHNodesUI()">📤 Export Fleet Nodes JSON</button>
          <button class="btn btn-secondary" onclick="importSSHNodesUI()">📥 Import Fleet Nodes JSON</button>
          <button class="btn" style="background:#0284c7" onclick="deployKeysFleetUI()">🔑 Deploy Keys (All Nodes)</button>
        </div>
        <textarea id="ssh-nodes-export-area" rows="8" placeholder="Fleet nodes JSON export/import data will appear here..."></textarea>
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
`

package cmdui

const uiAssetsMarkupSettingsA = `    <div id="tab-settings" class="content-area active">
      <div class="sub-tab-bar" style="display: flex; gap: 8px; margin-bottom: 1.25rem; border-bottom: 1px solid var(--border); padding-bottom: 0.75rem; flex-wrap: wrap;">
        <button type="button" class="sub-tab-btn active" onclick="showSettingsSubTab('all')" id="subnav-all">📑 All Settings</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('general')" id="subnav-general">🎨 General &amp; UI</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('cursor')" id="subnav-cursor">🛸 Cursor Fleet &amp; Sync</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('speed')" id="subnav-speed">⚡ Speed &amp; Automation</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('alerts')" id="subnav-alerts">🔔 Alerts &amp; Notifications</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('instances')" id="subnav-instances">🤖 Antigravity Instances</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('catalog')" id="subnav-catalog">✨ UI/UX Repo Index</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('terminal')" id="subnav-terminal">💻 Embedded Terminal</button>
      </div>

      <!-- 1. GENERAL & UI APPEARANCE -->
      <div class="card settings-subtab-pane" id="settings-pane-general">
        <h3>🎨 General &amp; UI Appearance</h3>
        <div class="form-group">
          <label>Theme</label>
          <select id="setting-theme">
            <option value="dark">Dark Theme (Standard)</option>
            <option value="light">Light Theme</option>
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
        <button class="btn" onclick="saveSettings()">Save General Settings</button>
      </div>

      <!-- 2. CURSOR FLEET & SYNC -->
      <div class="card settings-subtab-pane" id="settings-pane-cursor">
        <h3>🛸 Cursor Fleet &amp; Sync</h3>
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
        <button class="btn" onclick="saveSettings()">Save Cursor Fleet Settings</button>
      </div>

      <!-- 3. SPEED & AUTOMATION -->
      <div class="card settings-subtab-pane" id="settings-pane-speed">
        <h3>⚡ Speed &amp; Automation</h3>
        <div class="form-group">
          <label>LAP Default Lookback Hours (N = 24) (<code>lap.default_hours</code>)</label>
          <input type="number" id="setting-lap-hours" name="lap.default_hours" data-key="lap.default_hours" value="24" placeholder="24">
        </div>
        <div class="form-group">
          <label>Account Switch Fast-Forward Threshold (%) (<code>account_switch.threshold</code> &mdash; default 15% prod / 98% E2E test)</label>
          <input type="number" id="setting-account-switch-threshold" name="account_switch.threshold" data-key="account_switch.threshold" value="15" placeholder="15 (prod) or 98 (E2E test)">
        </div>
        <button class="btn" onclick="saveSettings()">Save Speed Settings</button>
      </div>

      <!-- 4. ALERTS & NOTIFICATIONS -->
      <div class="card settings-subtab-pane" id="settings-pane-alerts">
        <h3>🔔 Alerts &amp; Notifications</h3>
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
        <button class="btn" onclick="saveSettings()">Save Notification Settings</button>
      </div>

`

package cmdui

const uiAssetsMarkupSettingsA = `    <div id="tab-settings" class="content-area active">
      <div class="sub-tab-bar" style="display: flex; gap: var(--sp-2); margin-bottom: var(--sp-4); border-bottom: 1px solid var(--border); flex-wrap: wrap;">
        <button type="button" class="sub-tab-btn active" onclick="showSettingsSubTab('all')" id="subnav-all"><span data-icon="settings"></span>All Settings</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('general')" id="subnav-general"><span data-icon="pen-line"></span>General &amp; UI</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('cursor')" id="subnav-cursor"><span data-icon="rocket"></span>Cursor Fleet &amp; Sync</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('speed')" id="subnav-speed"><span data-icon="clock"></span>Speed &amp; Automation</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('alerts')" id="subnav-alerts"><span data-icon="alert-triangle"></span>Alerts &amp; Notifications</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('instances')" id="subnav-instances"><span data-icon="server"></span>Antigravity Instances</button>
        <button type="button" class="sub-tab-btn" onclick="showSettingsSubTab('catalog')" id="subnav-catalog"><span data-icon="book-open"></span>UI/UX Repo Index</button>
        <button type="button" class="sub-tab-btn" onclick="showTab('terminal')" id="subnav-terminal"><span data-icon="scroll-text"></span>Embedded Terminal</button>
      </div>

      <!-- 1. GENERAL & UI APPEARANCE -->
      <div class="card card-hero settings-subtab-pane" id="settings-pane-general">
        <h3><span data-icon="pen-line"></span> General &amp; UI Appearance</h3>
        <div class="form-group">
          <label for="setting-theme">Theme</label>
          <select id="setting-theme">
            <option value="dark">Dark</option>
            <option value="light">Light</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-graphics">Graphics Rendering &amp; Hardware Acceleration</label>
          <select id="setting-graphics">
            <option value="high">High Performance / Hardware Acceleration</option>
            <option value="fast">Fast Standard</option>
            <option value="plain">Plain / Low Graphics</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-auto-open">Auto-Open Browser on UI Launch</label>
          <select id="setting-auto-open">
            <option value="true">Enabled (Launch default browser automatically)</option>
            <option value="false">Disabled (Emit URL to console only)</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-port">Cluster REST Port</label>
          <input type="number" id="setting-port" value="49152">
        </div>
        <div class="form-group">
          <label for="setting-remote">Default Remote Target</label>
          <input type="text" id="setting-remote" placeholder="origin or alias">
        </div>
        <div class="form-group">
          <label for="setting-pull-direction">Pull Direction Orientation</label>
          <select id="setting-pull-direction">
            <option value="pull-left">Pull-Left (Replay from target down to left)</option>
            <option value="pull-right">Pull-Right (Replay from left source into right target)</option>
            <option value="bidirectional">Bidirectional Sync</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-pr-mode">Commit-Pull PR Replay Mode</label>
          <select id="setting-pr-mode">
            <option value="merges">merges (Feature branches &amp; release PRs)</option>
            <option value="feature-per-commit">feature-per-commit</option>
            <option value="direct">direct mainline</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-commitin-layout">Commit-In Migration Studio Layout</label>
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
          <label for="setting-machine-name">Machine Name (<code>machine.name</code>)</label>
          <input type="text" id="setting-machine-name" name="machine.name" data-key="machine.name" placeholder="dev-win-01 / ubuntu-node-02">
        </div>
        <div class="form-group">
          <label for="setting-machine-alias">Network Alias (<code>machine.alias</code> &mdash; auto-defaults to Local IPv4)</label>
          <input type="text" id="setting-machine-alias" name="machine.alias" data-key="machine.alias" placeholder="Auto-defaults to Local IPv4 if unset">
        </div>
        <div class="form-group">
          <label for="setting-special-secrets-name">Special Secrets Repository Name (<code>special_repos.secrets_name</code> &mdash; shortcut: <code>gitmap rs</code> / <code>gitmap cd rs</code>)</label>
          <input type="text" id="setting-special-secrets-name" name="special_repos.secrets_name" data-key="special_repos.secrets_name" value="repo-secrets" placeholder="repo-secrets">
        </div>
        <div class="form-group">
          <label for="setting-special-cache-name">Special Cache / Storage Repository Name (<code>special_repos.cache_name</code> &mdash; shortcut: <code>gitmap rc</code> / <code>gitmap cd rc</code>)</label>
          <input type="text" id="setting-special-cache-name" name="special_repos.cache_name" data-key="special_repos.cache_name" value="repo-cache" placeholder="repo-cache">
        </div>
        <div style="margin-top:var(--sp-4); padding-top:var(--sp-3); border-top:1px solid var(--border);">
          <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:var(--sp-2); flex-wrap:wrap; gap:var(--sp-2);">
            <h4 style="margin:0; font-size:var(--fs-base); color:var(--text);"><span data-icon="cpu"></span> AI IDE Setup (Cursor Alternatives: Cline, Continue, Zed)</h4>
            <div style="display:flex; gap:var(--sp-2);">
              <button class="btn btn-secondary btn-sm" type="button" onclick="loadIDESetupStatus()"><span data-icon="refresh-cw"></span> Refresh Status</button>
              <button class="btn btn-secondary btn-sm" type="button" onclick="testIDEConnection()"><span data-icon="activity"></span> Test Connection</button>
              <button class="btn btn-primary btn-sm" type="button" onclick="applyIDESetup()"><span data-icon="check"></span> Apply Setup</button>
            </div>
          </div>
          <p style="font-size:var(--fs-xs); color:var(--muted); margin-bottom:var(--sp-3);">
            Configure VS Code + Cline (Composer multi-file agent), Continue.dev (Tab-autocomplete &amp; inline edit), and Zed AI with Moonshot Kimi K3.1 Pro or custom providers across Windows, macOS, and Linux.
          </p>
          <div style="display:grid; grid-template-columns:repeat(auto-fit, minmax(200px, 1fr)); gap:var(--sp-2); margin-bottom:var(--sp-3);">
            <div class="form-group" style="margin:0;">
              <label for="ide-target-select" style="font-size:var(--fs-xs);">Target IDE / Extension</label>
              <select id="ide-target-select">
                <option value="all">All Targets (Cline, Continue, Zed)</option>
                <option value="cline">VS Code + Cline (Composer Agent)</option>
                <option value="continue">VS Code + Continue (Inline &amp; Tab-Autocomplete)</option>
                <option value="zed">Zed AI Editor (Native BYOK)</option>
              </select>
            </div>
            <div class="form-group" style="margin:0;">
              <label for="ide-provider-select" style="font-size:var(--fs-xs);">AI Provider</label>
              <select id="ide-provider-select" onchange="onIDEProviderChange()">
                <option value="kimi">Moonshot / Kimi (kimi-k3.1-pro)</option>
                <option value="deepseek">DeepSeek (deepseek-chat / reasoner)</option>
                <option value="openrouter">OpenRouter (Multi-model gateway)</option>
                <option value="anthropic">Anthropic Claude</option>
                <option value="ollama">Ollama (Local)</option>
                <option value="custom">Custom Endpoint</option>
              </select>
            </div>
            <div class="form-group" style="margin:0;">
              <label for="ide-api-key" style="font-size:var(--fs-xs);">API Key</label>
              <input type="password" id="ide-api-key" placeholder="Enter API Key (sk-...)">
            </div>
          </div>
          <div style="display:grid; grid-template-columns:1fr 1fr; gap:var(--sp-2); margin-bottom:var(--sp-3);">
            <div class="form-group" style="margin:0;">
              <label for="ide-base-url" style="font-size:var(--fs-xs);">Base URL</label>
              <input type="text" id="ide-base-url" value="https://api.moonshot.cn/v1">
            </div>
            <div class="form-group" style="margin:0;">
              <label for="ide-model-name" style="font-size:var(--fs-xs);">Model ID</label>
              <input type="text" id="ide-model-name" value="kimi-k3.1-pro">
            </div>
          </div>
          <div id="ide-setup-status-box" style="font-size:var(--fs-xs); color:var(--muted); background:var(--code-bg); padding:var(--sp-2); border-radius:var(--radius-sm); border:1px solid var(--border);">
            Click <strong>Refresh Status</strong> to detect installed IDEs and verify configuration paths.
          </div>
        </div>
      </div>

      <!-- 3. SPEED & AUTOMATION -->
      <div class="card settings-subtab-pane" id="settings-pane-speed">
        <h3><span data-icon="clock"></span> Speed &amp; Automation</h3>
        <div class="form-group">
          <label for="setting-lap-hours">LAP Default Lookback Hours (N = 24) (<code>lap.default_hours</code>)</label>
          <input type="number" id="setting-lap-hours" name="lap.default_hours" data-key="lap.default_hours" value="24" placeholder="24">
        </div>
        <div class="form-group">
          <label for="setting-account-switch-threshold">Account Switch Fast-Forward Threshold (%) (<code>account_switch.threshold</code> &mdash; default 15% prod / 98% E2E test)</label>
          <input type="number" id="setting-account-switch-threshold" name="account_switch.threshold" data-key="account_switch.threshold" value="15" placeholder="15 (prod) or 98 (E2E test)">
        </div>
        <div class="form-group">
          <label for="setting-ui-refresh-interval">UI Background Auto-Refresh Interval (ms) (<code>ui.refresh_interval_ms</code> &mdash; default 300000ms / 5m)</label>
          <input type="number" id="setting-ui-refresh-interval" name="ui.refresh_interval_ms" data-key="ui.refresh_interval_ms" value="300000" placeholder="300000">
        </div>
        <div class="form-group">
          <label for="setting-commit-timezone">Commit Schedule Timezone (<code>commit.timezone</code>)</label>
          <input type="text" id="setting-commit-timezone" name="commit.timezone" data-key="commit.timezone" value="Local" placeholder="Local or America/New_York">
        </div>
        <div class="form-group">
          <label for="setting-commit-schedule-active">Commit Schedule Enforcement Active (<code>commit.schedule_active</code>)</label>
          <select id="setting-commit-schedule-active" name="commit.schedule_active" data-key="commit.schedule_active">
            <option value="false">Disabled (Unconstrained commits anytime)</option>
            <option value="true">Enabled (Enforce weekday cutoff before 08:00 AM)</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-commit-weekday-cutoff">Commit Weekday Cutoff Time (<code>commit.weekday_cutoff</code>)</label>
          <input type="text" id="setting-commit-weekday-cutoff" name="commit.weekday_cutoff" data-key="commit.weekday_cutoff" value="08:00 AM" placeholder="08:00 AM">
          <p style="font-size:var(--fs-xs);color:var(--muted);margin-top:var(--sp-1);">Terminal check: <code>gitmap restrict status</code>. Add or edit with <code>--name</code>, <code>--blocked-after</code>, <code>--email</code>, <code>--tz</code>. Weekends: <code>--weekend-open</code> or <code>--weekend-closed</code>. Stay and commit after the cutoff with <code>--permit-stay</code> and <code>--auto-window</code>. <code>--commit-only</code> keeps the restriction and skips push.</p>
        </div>
        <div class="form-group" style="margin-top:var(--sp-3);">
          <label>Prompt Preview Header Metadata (Default Instance Template)</label>
          <div class="prompt-preview-header-wrapper" style="border:1px solid var(--border);border-radius:var(--radius-sm);background:var(--code-bg);padding:var(--sp-2);">
            <div style="display:flex;align-items:center;justify-content:space-between;">
              <div style="display:flex;align-items:center;gap:var(--sp-2);">
                <strong style="font-size:var(--fs-sm);color:var(--text);">Prompt Execution Preview</strong>
                <span class="badge badge-online">ACTIVE</span>
              </div>
              <details class="prompt-meta-details" style="font-size:var(--fs-xs);color:var(--muted);">
                <summary class="btn-xs" style="cursor:pointer;display:inline-flex;align-items:center;padding:2px 8px;background:var(--raised);border:1px solid var(--border);border-radius:var(--radius-sm);user-select:none;">Details</summary>
                <div style="margin-top:6px;padding:8px;background:var(--bg);border:1px solid var(--border);border-radius:var(--radius-sm);display:flex;flex-direction:column;gap:4px;">
                  <div><b>Instance ID:</b> <code>default-copy-814</code></div>
                  <div><b>Process PID:</b> <code>Live Subprocess</code></div>
                  <div><b>Error Backoff:</b> <code>5m Suppression Window Active</code></div>
                  <div><b>Credit Sync:</b> <code>&ge; 30s Throttle Active</code></div>
                  <div style="margin-top:4px;"><b>Prompt Context:</b><pre style="max-height:80px;overflow:auto;margin:4px 0 0 0;font-size:11px;padding:4px;background:var(--code-bg);">Execution environment: default-copy-814</pre></div>
                </div>
              </details>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. ALERTS & NOTIFICATIONS -->
      <div class="card settings-subtab-pane" id="settings-pane-alerts">
        <h3><span data-icon="alert-triangle"></span> Alerts &amp; Notifications</h3>
        <div class="form-group">
          <label for="setting-telegram-token">Telegram Two-Way Bot Setup &mdash; Bot Token (<code>telegram.bot_token</code>)</label>
          <input type="text" id="setting-telegram-token" name="telegram.bot_token" data-key="telegram.bot_token" placeholder="123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11">
        </div>
        <div class="form-group">
          <label for="setting-telegram-chat">Telegram Two-Way Bot Setup &mdash; Chat ID (<code>telegram.chat_id</code>)</label>
          <input type="text" id="setting-telegram-chat" name="telegram.chat_id" data-key="telegram.chat_id" placeholder="-1001234567890">
        </div>
        <div class="form-group" style="margin-top:var(--sp-2); margin-bottom:var(--sp-3);">
          <button type="button" class="btn btn-secondary" id="btn-test-telegram" onclick="testTelegramDispatch()"><span data-icon="rocket"></span>Test Telegram Dispatch</button>
        </div>
        <div class="form-group">
          <label for="setting-email-smtp">Email Speed SMTP Setup &mdash; SMTP Host (<code>email.smtp_host</code>)</label>
          <input type="text" id="setting-email-smtp" name="email.smtp_host" data-key="email.smtp_host" placeholder="smtp.gmail.com:587">
        </div>
        <div class="form-group">
          <label for="setting-email-from">Email Speed SMTP Setup &mdash; From Address (<code>email.from</code>)</label>
          <input type="email" id="setting-email-from" name="email.from" data-key="email.from" placeholder="bot@example.com">
        </div>
        <div class="form-group">
          <label for="setting-email-to">Email Speed SMTP Setup &mdash; To Address (<code>email.to</code>)</label>
          <input type="email" id="setting-email-to" name="email.to" data-key="email.to" placeholder="dev@example.com">
        </div>
      </div>
`

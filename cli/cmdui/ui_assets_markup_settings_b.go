package cmdui

const uiAssetsMarkupSettingsB = `      <!-- 5. ANTIGRAVITY MULTI-INSTANCE PANEL -->
      <div class="card settings-subtab-pane" id="settings-pane-instances">
        <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:0.75rem; flex-wrap:wrap; gap:8px;">
          <h3>🤖 Antigravity Multi-Instance Engine</h3>
          <button class="btn btn-secondary" type="button" onclick="loadAgyInstances()">🔄 Scan Instances</button>
        </div>
        <p style="font-size:0.85rem; color:var(--muted); margin-bottom:1rem;">
          Configure active Antigravity IDE instance, monitor language server host bridges, and view prompt queue partitions.
        </p>
        <div class="form-group">
          <label>Target Active Instance (<code>antigravity.instance_id</code>)</label>
          <select id="setting-agy-instance" data-key="antigravity.instance_id">
            <option value="primary">primary (Default Workspace: ~/.gemini/antigravity)</option>
          </select>
        </div>
        <div class="form-group">
          <label>Language Server Auto-Discovery</label>
          <select id="setting-agy-autodiscover" data-key="antigravity.auto_discover">
            <option value="true">Enabled (Auto-detect running language_server ports)</option>
            <option value="false">Disabled (Static instances only)</option>
          </select>
        </div>
        <div style="margin-top:1.25rem;">
          <h4 style="font-size:0.9rem; color:var(--text); margin-bottom:0.6rem;">Detected Instances &amp; Status</h4>
          <div id="instances-list-container" style="display:grid; grid-template-columns:repeat(auto-fit, minmax(280px, 1fr)); gap:12px;">
            <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:0.85rem;">
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:6px;">
                <strong style="color:var(--primary); font-size:0.9rem;">primary</strong>
                <span class="badge badge-online">RUNNING</span>
              </div>
              <div style="font-size:0.8rem; color:var(--muted); line-height:1.5; font-family:monospace;">
                PID: 4757 | Port: 127.0.0.1:33419<br>
                Path: ~/.gemini/antigravity
              </div>
            </div>
          </div>
        </div>
        <div style="margin-top:1.25rem;">
          <button class="btn" onclick="saveSettings()">Save Instance Settings</button>
        </div>
      </div>

      <!-- 6. AUTHORITATIVE UI/UX & CSS3 REPO INDEX -->
      <div class="card settings-subtab-pane" id="settings-pane-catalog">
        <h3>✨ Curated UI/UX &amp; CSS3 Repository Index</h3>
        <p style="font-size:0.85rem; color:var(--muted); margin-bottom:1.25rem;">
          High-craft architectural benchmarks for layout, accessibility primitives, HSL tokens, and micro-interactions.
        </p>
        <div style="display:grid; grid-template-columns:repeat(auto-fit, minmax(280px, 1fr)); gap:14px;">
          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">shadcn/ui</strong>
                <span class="badge badge-online" style="background:#1e3a5f; color:#60a5fa;">Radix + Tailwind</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">Accessible Primitives &amp; Token Theming</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Copy-paste accessible components built on Radix UI and Tailwind CSS with first-class HSL token architecture.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">shadcn-ui/ui</code>
              <a href="https://github.com/shadcn-ui/ui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Radix Primitives</strong>
                <span class="badge badge-online" style="background:#1e3a5f; color:#60a5fa;">Headless ARIA</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">WAI-ARIA Compliant Primitives</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Unstyled accessible components (dialogs, tooltips, popovers, dropdowns) providing headless behavioral foundations.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">radix-ui/primitives</code>
              <a href="https://github.com/radix-ui/primitives" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Aceternity UI</strong>
                <span class="badge badge-online" style="background:#3b1f5f; color:#c084fc;">Canvas Effects</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">Visual Craft &amp; Canvas Polish</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Modern interactive components, border beams, glowing cards, and canvas animations for high-impact experiences.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">mannupaaji/aceternity-ui</code>
              <a href="https://github.com/mannupaaji/aceternity-ui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Magic UI</strong>
                <span class="badge badge-online" style="background:#3b1f5f; color:#c084fc;">Landing Polish</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">Micro-Interactions &amp; Grids</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">20+ copy-paste components: retro grids, animated border rays, particle backgrounds, and tactile buttons.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">magicuidesign/magicui</code>
              <a href="https://github.com/magicuidesign/magicui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Tremor</strong>
                <span class="badge badge-online" style="background:#0f402c; color:#34d399;">Analytics UI</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">Data Telemetry &amp; Charts</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">React dashboard and telemetry components; clean dark mode metrics cards, KPI callouts, and area charts.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">tremorlabs/tremor</code>
              <a href="https://github.com/tremorlabs/tremor" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Lucide Icons</strong>
                <span class="badge badge-online" style="background:#422006; color:#fbbf24;">Vector Icons</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">Consistent Iconography</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Pixel-precise SVG icon system maintaining uniform 24px viewports with zero third-party script bloat.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">lucide-icons/lucide</code>
              <a href="https://github.com/lucide-icons/lucide" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; padding:1rem; display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:4px;">
                <strong style="color:var(--text); font-size:0.95rem;">Modern CSS Solutions</strong>
                <span class="badge badge-online" style="background:#1e3a5f; color:#60a5fa;">CSS3 Patterns</span>
              </div>
              <div style="font-size:0.75rem; color:var(--primary); font-weight:600; margin-bottom:6px;">CSS3 Architecture Recipes</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Modern CSS3 layout architectures: CSS grid, subgrid, container queries, and resilient responsive layouts.</p>
            </div>
            <div style="margin-top:10px; padding-top:8px; border-top:1px solid var(--border); font-size:0.75rem; display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">moderncss.dev</code>
              <a href="https://moderncss.dev" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>
        </div>
      </div>

      <!-- 7. INTERACTIVE EMBEDDED TERMINAL -->
      <div class="card settings-subtab-pane" id="settings-pane-terminal">
        <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:0.75rem; flex-wrap:wrap; gap:8px;">
          <div style="display:flex; align-items:center; gap:8px;">
            <span style="display:inline-flex; gap:5px; font-size:0.85rem;">
              <span style="color:#ef4444;">●</span>
              <span style="color:#f59e0b;">●</span>
              <span style="color:#10b981;">●</span>
            </span>
            <h3 style="margin-bottom:0; color:var(--text);">GitMap Live Terminal</h3>
          </div>
          <div style="display:flex; align-items:center; gap:8px; flex:1; max-width:480px; justify-content:flex-end;">
            <select id="term-node" style="width:auto; min-width:110px; padding:0.4rem 0.6rem;">
              <option value="local">local</option>
              <option value="u1">u1</option>
            </select>
            <input type="text" id="term-cwd" placeholder="cwd (optional)" style="flex:1; min-width:140px; padding:0.4rem 0.6rem;">
            <span id="term-badge" class="badge" style="display:none; white-space:nowrap;"></span>
          </div>
        </div>

        <!-- Quick command chips -->
        <div class="chip-container" style="display:flex; gap:6px; flex-wrap:wrap; margin-bottom:0.75rem; align-items:center;">
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap status')">gitmap status</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap pe')">gitmap pe</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap nodes ls')">gitmap nodes ls</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap cursor settings view')">gitmap cursor settings view</button>
          <button type="button" class="chip" onclick="runTerminalCommand('gitmap lap 24')">gitmap lap 24</button>
          <button type="button" class="chip" onclick="clearTerminalOutput()" style="margin-left:auto; background:#1e293b; color:var(--muted);">🧹 Clear</button>
        </div>

        <!-- Terminal Console Output -->
        <div id="term-output" style="background:var(--code-bg); border:1px solid var(--border); border-radius:6px; min-height:260px; max-height:420px; overflow-y:auto; padding:0.85rem; font-family:'Cascadia Code', Consolas, monospace; font-size:0.85rem; line-height:1.5; white-space:pre-wrap; color:#e2e8f0; margin-bottom:0.75rem;">GitMap Fleet Terminal ready. Type a command or click a quick action above.</div>

        <!-- Command Prompt Input -->
        <div style="display:flex; gap:8px; align-items:center;">
          <span style="font-family:monospace; color:var(--primary); font-weight:700; font-size:1.1rem;">$</span>
          <input type="text" id="term-input" placeholder="Enter command... (Press Enter or Run)" style="flex:1;">
          <button type="button" class="btn" id="term-run-btn" onclick="runTerminalCommand()">⚡ Run</button>
        </div>
      </div>
    </div>

    <!-- COMMITIN TAB -->
`

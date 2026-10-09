package cmdui

const uiAssetsMarkupSettingsB = `      <!-- 5. ANTIGRAVITY MULTI-INSTANCE PANEL -->
      <div class="card settings-subtab-pane" id="settings-pane-instances">
        <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:var(--sp-3); flex-wrap:wrap; gap:var(--sp-2);">
          <h3><span data-icon="server"></span> Antigravity Multi-Instance Engine</h3>
          <button class="btn btn-secondary" type="button" onclick="loadAgyInstances()"><span data-icon="check"></span>Scan Instances</button>
        </div>
        <p style="font-size:var(--fs-sm); color:var(--muted); margin-bottom:var(--sp-4);">
          Configure active Antigravity IDE instance, monitor language server host bridges, and view prompt queue partitions.
        </p>
        <div class="form-group">
          <label for="setting-agy-instance">Target Active Instance (<code>antigravity.instance_id</code>)</label>
          <select id="setting-agy-instance" data-key="antigravity.instance_id">
            <option value="primary">primary (Default Workspace: ~/.gemini/antigravity)</option>
          </select>
        </div>
        <div class="form-group">
          <label for="setting-agy-autodiscover">Language Server Auto-Discovery</label>
          <select id="setting-agy-autodiscover" data-key="antigravity.auto_discover">
            <option value="true">Enabled (Auto-detect running language_server ports)</option>
            <option value="false">Disabled (Static instances only)</option>
          </select>
        </div>
        <div style="margin-top:1.25rem;">
          <h4 style="font-size:var(--fs-sm); color:var(--text); margin-bottom:var(--sp-2);">Detected Instances &amp; Status</h4>
          <div id="instances-list-container" style="display:grid; grid-template-columns:repeat(auto-fit, minmax(280px, 1fr)); gap:var(--sp-3);">
            <p style="color:var(--muted); font-size:var(--fs-sm); margin:0;">No instances detected — click Scan Instances.</p>
          </div>
        </div>
      </div>

      <!-- 6. AUTHORITATIVE UI/UX & CSS3 REPO INDEX -->
      <div class="card settings-subtab-pane" id="settings-pane-catalog">
        <h3><span data-icon="book-open"></span> Curated UI/UX &amp; CSS3 Repository Index</h3>
        <p style="font-size:var(--fs-sm); color:var(--muted); margin-bottom:var(--sp-4);">
          High-craft architectural benchmarks for layout, accessibility primitives, HSL tokens, and micro-interactions.
        </p>
        <div style="display:grid; grid-template-columns:repeat(auto-fit, minmax(280px, 1fr)); gap:var(--sp-3);">
          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">shadcn/ui</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Radix + Tailwind</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">Accessible Primitives &amp; Token Theming</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Copy-paste accessible components built on Radix UI and Tailwind CSS with first-class HSL token architecture.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">shadcn-ui/ui</code>
              <a href="https://github.com/shadcn-ui/ui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Radix Primitives</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Headless ARIA</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">WAI-ARIA Compliant Primitives</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Unstyled accessible components (dialogs, tooltips, popovers, dropdowns) providing headless behavioral foundations.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">radix-ui/primitives</code>
              <a href="https://github.com/radix-ui/primitives" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Aceternity UI</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Canvas Effects</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">Visual Craft &amp; Canvas Polish</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Modern interactive components, border beams, glowing cards, and canvas animations for high-impact experiences.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">mannupaaji/aceternity-ui</code>
              <a href="https://github.com/mannupaaji/aceternity-ui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Magic UI</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Landing Polish</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">Micro-Interactions &amp; Grids</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">20+ copy-paste components: retro grids, animated border rays, particle backgrounds, and tactile buttons.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">magicuidesign/magicui</code>
              <a href="https://github.com/magicuidesign/magicui" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Tremor</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Analytics UI</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">Data Telemetry &amp; Charts</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">React dashboard and telemetry components; clean dark mode metrics cards, KPI callouts, and area charts.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">tremorlabs/tremor</code>
              <a href="https://github.com/tremorlabs/tremor" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Lucide Icons</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">Vector Icons</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">Consistent Iconography</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Pixel-precise SVG icon system maintaining uniform 24px viewports with zero third-party script bloat.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">lucide-icons/lucide</code>
              <a href="https://github.com/lucide-icons/lucide" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>

          <div style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); padding:var(--sp-3); display:flex; flex-direction:column; justify-content:space-between;">
            <div>
              <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:var(--sp-1);">
                <strong style="color:var(--text); font-size:var(--fs-sm);">Modern CSS Solutions</strong>
                <span class="badge badge-online" style="background:var(--badge-ok-bg); color:var(--badge-ok-fg);">CSS3 Patterns</span>
              </div>
              <div style="font-size:var(--fs-xs); color:var(--primary); font-weight:600; margin-bottom:var(--sp-2);">CSS3 Architecture Recipes</div>
              <p style="font-size:0.8rem; color:var(--muted); line-height:1.4;">Modern CSS3 layout architectures: CSS grid, subgrid, container queries, and resilient responsive layouts.</p>
            </div>
            <div style="margin-top:var(--sp-2); padding-top:var(--sp-2); border-top:1px solid var(--border); font-size:var(--fs-xs); display:flex; justify-content:space-between;">
              <code style="color:var(--muted);">moderncss.dev</code>
              <a href="https://moderncss.dev" target="_blank" style="color:var(--primary); text-decoration:none; font-weight:600;">Visit &rarr;</a>
            </div>
          </div>
        </div>
      </div>

      <!-- 7. INTERACTIVE EMBEDDED TERMINAL -->
      <div class="card settings-subtab-pane" id="settings-pane-terminal">
        <div style="display:flex; align-items:center; justify-content:space-between; margin-bottom:var(--sp-3); flex-wrap:wrap; gap:var(--sp-2);">
          <div style="display:flex; align-items:center; gap:8px;">
            <span style="display:inline-flex; gap:5px; font-size:0.85rem;">
              <span style="color:#ef4444;">●</span>
              <span style="color:#f59e0b;">●</span>
              <span style="color:#10b981;">●</span>
            </span>
            <h3 style="margin-bottom:0; color:var(--text);">GitMap Live Terminal</h3>
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
        <div id="term-output" tabindex="0" aria-label="Terminal output" style="background:var(--code-bg); border:1px solid var(--border); border-radius:var(--radius-sm); min-height:180px; max-height:420px; overflow-y:auto; padding:var(--sp-3); font-family:'Cascadia Code', Consolas, monospace; font-size:0.85rem; line-height:1.5; white-space:pre-wrap; color:var(--code-text); margin-bottom:0.75rem;">GitMap Fleet Terminal ready. Type a command or click a quick action above.</div>

        <!-- Command Prompt Input -->
        <div style="display:flex; gap:var(--sp-2); align-items:center;">
          <span style="font-family:monospace; color:var(--primary); font-weight:700; font-size:1.1rem;">$</span>
          <input type="text" id="term-input" placeholder="Enter command... (Press Enter or Run)" style="flex:1;">
          <button type="button" class="btn" id="term-run-btn" onclick="runTerminalCommand()"><span data-icon="rocket"></span>Run</button>
        </div>
      </div>

      <!-- ONE sticky save bar for all settings panes (replaces the per-pane Save buttons) -->
      <div class="settings-savebar">
        <span style="font-size: 0.8rem; color: var(--muted);">Saves every settings pane above.</span>
        <button class="btn" type="button" onclick="saveSettings()"><span data-icon="check"></span>Save All Settings</button>
      </div>
    </div>

    <!-- COMMITIN TAB -->
`

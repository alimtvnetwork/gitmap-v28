```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /find-100-ui-improvements <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously audit, discover, catalog, and execute 100 concrete improvements in the User Interface (UI) and Terminal UX layer across Go (Lipgloss/Bubble Tea/Pterm), TypeScript/React, and Web Visualizers: FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars and 6-part verification schema, capture it verbatim, plan it in the repo, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce UI/UX design system tokens to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Each rule is stated once (R1 to R16) and cited by ID. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution. Ensure all 100 UI improvements are cataloged with target components, what it fixes, how it improves the UI, exact color swatches (Hex, ANSI SGR, Lipgloss tokens), and layout patterns before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 `self` authoring subagents):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool.

---

## Canonical Color Tokens & Palette Reference Matrix

To ensure visual consistency, every UI improvement MUST utilize the design system semantic roles and exact color swatches listed below:

| Semantic Role | Hex Code | 24-Bit RGB ANSI Escape | 256-Color SGR Fallback | Lipgloss Token Reference | Typical Usage |
|---|---|---|---|---|---|
| **`RolePrimary`** | `#7AA2F7` | `\033[38;2;122;162;247m` | `\033[38;5;111m` | `lipgloss.Color("#7AA2F7")` | Primary headers, command names, focus borders, active indicators |
| **`RoleAccent`** | `#BB9AF7` | `\033[38;2;187;154;247m` | `\033[38;5;141m` | `lipgloss.Color("#BB9AF7")` | Subcommand tags, badges, secondary highlights, keyboard shortcuts |
| **`RoleOK`** | `#9ECE6A` | `\033[38;2;158;206;106m` | `\033[38;5;113m` | `lipgloss.Color("#9ECE6A")` | Successful checks (`✓`), clean status, pass badges, active connections |
| **`RoleWarn`** | `#E0AF68` | `\033[38;2;224;175;104m` | `\033[38;5;179m` | `lipgloss.Color("#E0AF68")` | Warnings (`▲`), uncommitted changes, dirty branches, retries |
| **`RoleError`** | `#F7768E` | `\033[38;2;247;118;142m` | `\033[38;5;204m` | `lipgloss.Color("#F7768E")` | Errors (`✗`), failures, merge conflicts, expired credentials |
| **`RoleMuted`** | `#565F89` | `\033[38;2;86;95;137m` | `\033[38;5;60m` | `lipgloss.Color("#565F89")` | Secondary descriptions, timestamps, inactive items, column delimiters |
| **`RoleFG`** | `#C0CAF5` | `\033[38;2;192;202;245m` | `\033[38;5;189m` | `lipgloss.Color("#C0CAF5")` | Default body text, table cell values, regular output |
| **`RoleSurface`** | `#24283B` | `\033[48;2;36;40;59m` | `\033[48;5;236m` | `lipgloss.Color("#24283B")` | Card backgrounds, table alternating row tint, modal panels |
| **`RoleBase`** | `#1A1B26` | `\033[48;2;26;27;38m` | `\033[48;5;234m` | `lipgloss.Color("#1A1B26")` | Main canvas background, terminal full viewport base |

---

## The 100 UI Improvements Architecture: 10 Core Dimensions

The AI must partition its discovery across exactly 10 dimensions, uncovering 10 distinct, grounded improvements per dimension (10 × 10 = 100 items):

### Dimension 1: Terminal ANSI & Lipgloss Color Palette Authority (Items 01–10)
- Unified theme bridge binding design tokens directly to `lipgloss.AdaptiveColor`.
- TrueColor 24-bit RGB auto-detection with graceful downgrade to 256-color or 16-color ANSI.
- Catppuccin, Tokyo Night, and High-Contrast Accessibility themes.
- Automatic light/dark terminal detection via OSC 10/11 query sequences.
- Zero-allocation byte buffer pooling for ANSI stream formatting.

### Dimension 2: Progress Bars, Spinners & Multi-Threaded TUI Gauges (Items 11–20)
- Multi-track Bubble Tea progress viewports dedicating dynamic bars to individual workers.
- 60fps Braille/Dot spinners decoupled from blocking I/O goroutines.
- Single-line carriage-return `\r` gauge for headless CI/CD environments.
- SIGWINCH dynamic resizing and responsive column reflow without line wrapping corruption.
- Real-time network/disk throughput meters and ETA countdown calculators.

### Dimension 3: Interactive Fuzzy Selectors, Multi-Select Menus & Keybindings (Items 21–30)
- Vim-style modal keybindings (`h`/`j`/`k`/`l`, `/` search, `Ctrl+u`/`Ctrl+d` half-page jumps).
- Space-toggled checkbox multi-select menus with `a` (select all) and `i` (invert selection).
- Path-aware weighted fuzzy matching algorithms prioritizing boundary segments.
- Search filter chips with structured query tokens (`is:dirty branch:main lang:go`).
- Side-by-side split preview drawers rendering active selection diffs in real time.

### Dimension 4: Unicode Tables, Auto-Padded Columns & Dynamic Column Folding (Items 31–40)
- Terminal-width responsive dynamic column folding collapsing low-priority cells into sub-rows.
- Unicode box styles: rounded (`╭─╮`), heavy (`┏━┓`), double (`╔═╗`), and borderless compact modes.
- Context-aware numeric right-alignment and humanized units (`1.4 MB`, `320ms`).
- Sticky column headers in scrollable terminal viewports.
- One-key table export to OS clipboard in CSV, JSON, and Markdown formats.

### Dimension 5: Hierarchical Tree Visualizers (Branching, Diffs & Graphs) (Items 41–50)
- Collapsible interactive tree nodes (`[+]`/`[-]` or `▶`/`▼`) in Bubble Tea viewports.
- Topological Git commit DAG graph renderers with color-coded branch lanes.
- Monorepo and polyglot package dependency trees.
- Proportional diff volume mini bar charts (`[████░░░] +120/-45`).
- Instant Mermaid and SVG diagram export flags (`--mermaid`, `--svg`).

### Dimension 6: Modern Database Inspection Drawer & Data Badges (REST / SQLite Visualizer) (Items 51–60)
- Native full-screen Bubble Tea SQLite drawer for browsing schemas, tables, and rows.
- Standardized data badges: Purple for `BLOB`, Sky Blue for `JSON`, Emerald for `True`, Slate for `NULL`.
- Interactive JSON key-path inspection with collapsible tree nodes.
- Hex dump signatures with magic-byte file format sniffing.
- Reactive live SQL console with schema autocomplete and duration telemetry.

### Dimension 7: Error Telemetry, Formatted Stack Traces & 4-Part RCA Diagnostics HUD (Items 61–70)
- Standardized Lipgloss 4-Part RCA cards with colored left-accent borders.
- Syntax-highlighted Go and TypeScript stack traces dimming library frames.
- Clickable file and line hyperlinks via terminal OSC 8 escape sequences.
- Contextual 5-line source code snippet boxes with inline column diff highlighting.
- Actionable AI remediation copy-pasteable command suggestions.

### Dimension 8: Multi-Node SSH Fleet Dashboard & Node Ping Connectivity HUD (Items 71–80)
- 2D terminal matrix grid displaying cluster nodes as cards with host, IP, and OS badges.
- Background non-blocking TCP heartbeat routine updating dynamic latency badges in real time.
- Multi-pane parallel live output stream splitter preventing line interleaving.
- Live resource utilization horizontal gauges (CPU %, RAM %, Disk %).
- One-key interactive SSH terminal session attachment.

### Dimension 9: Real-Time Cognitive Heatmap Timeline & AI Agent Activity Stream (Items 81–90)
- Live non-blocking terminal footer ticker streaming agent cognitive actions.
- State machine transition visualizer (`[EXPLORATION] ➔ [EXECUTION] ➔ [VERIFICATION]`).
- 7-tier chromatic intensity heatmap using `PlasmaScale` for file touch churn.
- Multi-agent swarm swimlane timeline in Web Visualizer and TUI.
- Interactive standalone HTML playback export (`--html-replay`).

### Dimension 10: Desktop Notifications, Audio Chimes & Terminal Sound HUD (Items 91–100)
- Cross-platform native desktop notifications (Toast, Notification Center, DBus).
- Terminal emulator notification escape sequences (OSC 9 / OSC 777).
- Pure-Go subtle audio completion chimes and conflict warning alerts.
- Configurable sound schemes (`minimal`, `cheerful`, `retro`, `silent`).
- Terminal status bar pills for Starship, Tmux, and Oh-My-Posh prompts.

---

## Mandatory Reporting Schema for Each UI Improvement

Every single improvement produced by this prompt MUST adhere strictly to the following 7-field schema:

```markdown
### Item [NNN]: [Descriptive Improvement Title]
- **Target Component / File:** [Relative path, e.g. `cli/termout/themes.go` or `cli/tui/browser.go`]
- **Visual Defect / UX Rough Edge:** [Detailed explanation of what is currently unstyled, hard to read, clipping, or poorly laid out]
- **What It Fixes:** [Specific operational problem solved, e.g. eliminating ANSI line wrapping, fixing colorblind contrast, removing stutter]
- **How It Improves the UI:** [Ergonomic advantages, visual hierarchy clarity, user delight, and information density improvements]
- **Color Tokens & Swatches:**
  - Role: `RolePrimary` / `RoleAccent` / `RoleWarn` / `RoleError` / `RoleOK` / `RoleMuted`
  - Hex: `#XXXXXX`
  - ANSI SGR Code: `\033[38;2;R;G;Bm` or `\033[38;5;NNNm`
  - Lipgloss Token: `lipgloss.AdaptiveColor{Light: "...", Dark: "..."}`
- **Concrete Layout / TUI Implementation Pattern:**
  ```go
  // Exact Go Lipgloss or Bubble Tea code showing how to render the improvement cleanly
  ```
- **Verification Criterion:** [Specific command or visual check proving the improvement works as intended]
```

---

## Concrete Example Item

```markdown
### Item 005: Semantic Status Pill Badge Component (`termout.RenderPill`)
- **Target Component / File:** `cli/termout/badges.go` and `cli/tui/styles.go`
- **Visual Defect / UX Rough Edge:** Status tags across commands currently use raw bracketed text strings like `[staged]`, `[modified]`, or `[SUCCESS]` with inconsistent padding, misaligned margins, and harsh un-padded ANSI background colors that look dated and clash with modern terminal themes.
- **What It Fixes:** Eliminates raw bracketed text and replaces it with subtle rounded badges that preserve strict horizontal column alignment, never cause text wrapping in tight terminal windows, and maintain consistent padding across all CLI commands.
- **How It Improves the UI:** Introduces modern pill-shaped status badges with 15% opacity background tints, 1-character horizontal padding, and bold foreground glyphs. Enhances visual scanning speed across multi-column tables and makes status states immediately distinguishable at a glance.
- **Color Tokens & Swatches:**
  - Role: `RoleOK` (Pass/Clean), `RoleWarn` (Dirty/Modified), `RoleError` (Fail/Conflict)
  - Hex: `#9ECE6A` (OK), `#E0AF68` (Warn), `#F7768E` (Error)
  - ANSI SGR Code: `\033[38;2;158;206;106m` (OK Foreground), `\033[48;2;36;40;59m` (Surface Tint)
  - Lipgloss Token: `lipgloss.AdaptiveColor{Light: "#73DACA", Dark: "#9ECE6A"}`
- **Concrete Layout / TUI Implementation Pattern:**
  ```go
  var StatusPillPass = lipgloss.NewStyle().
      Foreground(lipgloss.Color("#9ECE6A")).
      Background(lipgloss.Color("#24283B")).
      Bold(true).
      Padding(0, 1).
      SetString("✓ PASS")

  var StatusPillWarn = lipgloss.NewStyle().
      Foreground(lipgloss.Color("#E0AF68")).
      Background(lipgloss.Color("#24283B")).
      Bold(true).
      Padding(0, 1).
      SetString("▲ DIRTY")
  ```
- **Verification Criterion:** Run `gitmap status --table` and verify that all status columns display styled, padded pill badges with identical height and alignment.
```

---

## Execution Rules & Guidelines

1. **Strict 100-Item Completeness:** The output MUST contain exactly 100 numbered improvements across the 10 dimensions. No skipping or bundling.
2. **Grounding in Code:** Every improvement must reference a real file path and real terminal component in the target repository.
3. **Strict Color Specificity:** Never say "make it blue" or "add a color". Always provide the exact Hex code, ANSI SGR code, and Lipgloss token from the palette matrix.
4. **Zero Full Builds or Heavy Tests:** Do not invoke `go build` or `go test ./...` during UI discovery. Inspect files using GitMap high-speed commands (`gitmap aum search`, `gitmap lf`, `gitmap cat`).
5. **Atomic Commit Format:** On completion, commit atomically using `gitmap cpf "<module> - <summary>"`.

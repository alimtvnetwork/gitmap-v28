# Find 100 Improvements in the Whole System

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /find-100-improvements-in-whole-system <target-scope>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously audit, discover, plan, and orchestrate 100 concrete, high-impact improvements across the whole system (Developer Ergonomics, Agent Friction Reduction, UI/UX Materiality, Multi-Repo Fleet Sync, Code Immutability, Performance, Observability, and Spec Governance) with a strict Zero Harm / Regression Safety Matrix, backward compatibility guarantees, and concrete blueprints. Enqueue all 100 improvements as pending plans in `.ai-memory/plans/pending/XX-<slug>.md` and subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`, and activate the Continuous Improvement Flywheel to generate the next 100 system improvements on each cycle.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Every improvement MUST adhere to the 8-Part Specification Schema, Zero Harm invariants, strict relative Git paths, lowercase filenames, and GitMap command primacy.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution for improved work or less friction, faster work and more pure and effective work. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

---

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 `self` authoring subagents — ZERO SOLO SPEC AUTHORING):** Spawn 3 authoring subagents (`TypeName: "self"`) to author modular, disjoint spec files and subtasks (never writing the same file). Solo authoring by the lead orchestrator is STRICTLY FORBIDDEN.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. The lead orchestrator MUST NOT wait for an entire 3-agent wave to finish. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel reading, companion linters, memory logging) to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/find-100-improvements-in-whole-system/skill.md` does not exist in the workspace, create it now.
2. Mirror the skill into `.cursor/skills/find-100-improvements-in-whole-system/skill.md` using Cursor link syntax (`[<skill>](file;.cursor/skills/<skill>)`).
3. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
4. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## The Unified Master Pipeline (Atomic Numbered Steps)

Execute this task via a strict 3-Phase pipeline. Do not skip steps.

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)

Before executing any file searches, scans, spec writing, or code changes, you must execute Phase 1A:

1. **Top-Instruction Priority Verification:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) must be verified as highest priority and non-negotiable.
2. **Extreme "Understood" Turn-1 Showcase (Turn 1 Action):** In your VERY FIRST response turn upon receiving the prompt, you MUST output a comprehensive, in-depth understanding block directly in visible chat before any background tool execution. Any generic boilerplate or superficial 1-sentence summaries are protocol auto-reject failures. Turn 1 MUST display the top-level mandate containing all 5 Mandatory Pillars followed by the confirmed task breakdown using the 6-part verification schema:
   * **Pillar 1: Strategic Intent & Systems Reasoning:** Articulate overarching architectural goal, user intent, operational rationale, and systems integration mechanics. Explain *why* this change is necessary and how downstream components interact.
   * **Pillar 2: Explicit Scope Boundaries:** Strict two-column partition defining In-Scope vs. strictly Excluded or Deferred items.
   * **Pillar 3: Non-Negotiable Constraints & Negative Guardrails Ingestion:** Explicit recitation of active constraints: raw search ban (`rg`, `ripgrep`, `Select-String`, `git grep`), worker git ban (no git commands in subagents), zero builds/tests (never run all tests locally or full build unless explicitly asked; every time code is written, do not build the code; only do end-to-end testing based on mentioned test cases if explicitly asked; do CI/CD verification only if asked, otherwise do not do it), strict relative Git paths, lowercase naming, zero secrets in repo (`repo-secrets` offloading), GitMap hyphen-separated commit format without colons.
   * **Pillar 4: Blast Radius & Dependency Topology:** Map impacted subsystems, files, downstream callers, and prompt/skill indexes with explicit risk containment strategies.
   * **Pillar 5: Multi-Agent Concurrency Topology:** Detail swarm structure (`A = 3`, `H = 2`, Capacity = 6), disjoint file assignments, and immediate dynamic respawn.
   * **Task Comprehension & Model Runtime Invariants:** In the task breakdown, every single task MUST explicitly state `- **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints]` and `- **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']`.
3. **Lossless Verbatim Capture:** Store incoming prompt losslessly under `## User Request (Verbatim)` in canonical spec and parent plan.
4. **Discrete Deliverables Extraction:** Break down requirements into discrete, actionable items with ordered traceable IDs (`Task-01`, `Task-02`, ...).
5. **Mandatory Same-Turn Tool Chaining:** Emit the comprehension mandate and breakdown in chat with clean vertical formatting, and in the **EXACT SAME TURN**, invoke your first tool call. NEVER emit text alone (which ends the turn prematurely), and never ask "Should I proceed?".

---

## 8-Part Structured Specification Schema

Every single improvement across the 10 System Clusters MUST strictly adhere to the following 8-part schema:

```markdown
### SYS-IMP-NNN: [Descriptive Title of System Improvement]

1. **Identifier & Subsystem:**
   - **ID:** `SYS-IMP-NNN`
   - **Subsystem:** [Developer Ergonomics | Agent Swarm | UI Materiality | Fleet Sync | Code Quality | Performance | Observability | Testing | Governance]
   - **Target Boundaries:** [Relative file paths or packages targeted, e.g. `cmd/gitmap/`, `02-spec/`, `src/`]

2. **System-Level User Experience (UX) Impact:**
   - **Cognitive Load Reduction:** [How this reduces mental friction for developers or agents]
   - **Iteration Speed / Workflow Velocity:** [Measurable acceleration in development, testing, or review cycles]
   - **Error Clarity & Self-Healing:** [How failures are made instantly understandable and actionable]

3. **Cause (Current Deficiency / Root Friction):**
   - [Root Cause Analysis of the current inefficiency, fragility, ambiguity, or friction in 1-2 sharp sentences]

4. **Effect (Target State & Behavioral Outcome):**
   - [Precise target behavior, quantitative metrics, or qualitative improvements achieved]

5. **Zero Harm / Regression Safety Matrix:**
   - **Confinement Rules:** [Strict isolation preventing collateral damage]
   - **Build/Test Invariant:** [Zero full builds or test suites triggered during active editing]
   - **Zero Secrets Invariant:** [No credentials or private tokens exposed]

6. **Backward Compatibility Guarantee:**
   - **CLI & Flags:** [Preservation of existing flags, aliases, and arguments]
   - **API Signatures:** [Preservation of method signatures, type aliases, or deprecation forwarders like `type Fault = AppError`]
   - **Rollback Safety:** [Instant mechanical rollback capability]

7. **Concrete Implementation Blueprint:**
   - [Concise, surgical code snippet, configuration diff, or schema enhancement]

8. **Binary Verification & Acceptance Criteria:**
   - **Verification Command:** [Exact GitMap command, script run, or test check returning exit code 0]
   - **Pass Criteria:** [Deterministic, binary true/false condition proving success]
```

---

## The 10 System UX Clusters (100 Improvements Catalog)

### Cluster 01: Developer Experience (DX) & CLI Ergonomics (SYS-IMP-001 – SYS-IMP-010)

#### SYS-IMP-001: GitMap Fuzzy Command Palette & Typo Tolerance
1. **Identifier & Subsystem:** `SYS-IMP-001` | CLI Ergonomics | `cmd/gitmap/`
2. **System-Level UX Impact:** Eliminates command-not-found frustration; automatically suggests correct commands when minor typos occur.
3. **Cause:** Typing `gitmap fin` or `gitmap serach` triggers hard errors, halting agent or developer flow.
4. **Effect:** CLI computes Levenshtein distance <= 2, prints "Did you mean: `gitmap find`?" and exits gracefully with suggestion code 0.
5. **Zero Harm / Regression Safety:** Read-only command dispatch routing; no alteration to existing command handlers.
6. **Backward Compatibility Guarantee:** Exact command matches execute identically with zero overhead.
7. **Concrete Implementation Blueprint:**
   ```go
   if cmd == nil {
       if suggestion := findClosestCommand(inputCmd, registeredCommands); suggestion != "" {
           fmt.Fprintf(os.Stderr, "Unknown command '%s'. Did you mean '%s'?\n", inputCmd, suggestion)
           return nil
       }
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap fnd readme.md` prints suggestion and exits 0.

#### SYS-IMP-002: Terminal Width Auto-Adaptation & Pager Prevention
1. **Identifier & Subsystem:** `SYS-IMP-002` | CLI Output Rendering | `pkg/termtable/`
2. **System-Level UX Impact:** Prevents table mangling and line-wrapping artifacts on narrow IDE terminal panes.
3. **Cause:** Hardcoded column widths (>120 cols) cause horizontal wrapping on standard 80-column agent sidebars.
4. **Effect:** Terminal columns automatically scale to `min(stdoutWidth, 100)` with intelligent column truncation.
5. **Zero Harm / Regression Safety:** Confined strictly to CLI output formatters.
6. **Backward Compatibility Guarantee:** `--raw` and `--json` flags bypass table formatting entirely.
7. **Concrete Implementation Blueprint:**
   ```go
   cols, _, err := term.GetSize(int(os.Stdout.Fd()))
   if err != nil || cols <= 0 { cols = 80 }
   table.SetMaxWidth(cols)
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap lf` output fits within active terminal dimensions without text wrapping.

#### SYS-IMP-003: One-Step Shell Auto-Completion Installer
1. **Identifier & Subsystem:** `SYS-IMP-003` | Shell Integration | `cmd/gitmap/completion.go`
2. **System-Level UX Impact:** Speeds up CLI interaction by 3x through Tab completion for all subcommands and flags.
3. **Cause:** Developers must manually configure completion scripts across PowerShell, Bash, and Zsh.
4. **Effect:** `gitmap completion install` inspects active shell and writes completion hook automatically.
5. **Zero Harm / Regression Safety:** Appends only safe, guarded shell function blocks to user profiles.
6. **Backward Compatibility Guarantee:** Existing manual completions remain fully valid.
7. **Concrete Implementation Blueprint:**
   ```go
   func runInstallCompletion(shell string) *appfault.AppError {
       return profile.AppendSafeHook(shell, "gitmap completion " + shell)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap completion install --dry-run` outputs valid profile path and exits 0.

#### SYS-IMP-004: Live Progress Ticker with Non-Flicker Spinner
1. **Identifier & Subsystem:** `SYS-IMP-004` | CLI Telemetry | `pkg/spinner/`
2. **System-Level UX Impact:** Provides immediate feedback during multi-thousand file operations without cluttering output logs.
3. **Cause:** Long searches over 10k files appear frozen for 2-3 seconds, confusing developers.
4. **Effect:** Lightweight ANSI spinner updates on stderr at 10 Hz; cleanly wiped upon completion.
5. **Zero Harm / Regression Safety:** Emits only to stderr; stdout streams pure data.
6. **Backward Compatibility Guarantee:** Automatically disabled when stdout/stderr is redirected or non-TTY.
7. **Concrete Implementation Blueprint:**
   ```go
   if isatty.IsTerminal(os.Stderr.Fd()) {
       ticker := spinner.New(os.Stderr, "Scanning repositories...")
       defer ticker.Stop()
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap aum search "test" . -e .go` displays spinner on stderr and pure matches on stdout.

#### SYS-IMP-005: Color-Coded Structured Output Formatter
1. **Identifier & Subsystem:** `SYS-IMP-005` | CLI Data Presentation | `pkg/output/`
2. **System-Level UX Impact:** High-speed visual scanning of search results, file listings, and pipeline logs.
3. **Cause:** Monochromatic text output forces high cognitive overhead when scanning dense logs.
4. **Effect:** Path in cyan, line numbers in yellow, match tokens in bold green, errors in red.
5. **Zero Harm / Regression Safety:** Standard ANSI escapes; NO_COLOR environment variable respected.
6. **Backward Compatibility Guarantee:** `--no-color` flag and piped output strip all ANSI escapes automatically.
7. **Concrete Implementation Blueprint:**
   ```go
   func FormatMatch(path string, line int, text string) string {
       if isNoColor() { return fmt.Sprintf("%s:%d:%s", path, line, text) }
       return fmt.Sprintf("\033[36m%s\033[0m:\033[33m%d\033[0m:%s", path, line, text)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `NO_COLOR=1 gitmap find "*.go"` outputs uncolored text; standard invocation uses color.

#### SYS-IMP-006: GitMap Environment Health Preflight (`gitmap doctor`)
1. **Identifier & Subsystem:** `SYS-IMP-006` | Developer Diagnostics | `cmd/gitmap/doctor.go`
2. **System-Level UX Impact:** Diagnoses environment misconfigurations (SQLite locks, missing git credentials, path issues) in <500ms.
3. **Cause:** Developers lose time debugging silent environment errors or missing database write permissions.
4. **Effect:** `gitmap doctor` checks 8 health vectors, printing a clear pass/fail checklist with one-line fixes.
5. **Zero Harm / Regression Safety:** Strictly read-only checks; makes zero mutations.
6. **Backward Compatibility Guarantee:** Independent standalone subcommand.
7. **Concrete Implementation Blueprint:**
   ```go
   func RunDoctor() *appfault.AppError {
       checks := []Check{checkGitVersion, checkSQLiteWAL, checkRepoSecrets, checkWorkingTree}
       return executeDoctorSuite(checks)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap doctor` returns exit code 0 when environment is healthy.

#### SYS-IMP-007: CLI Interactive Diff Viewer with Chunk Folding
1. **Identifier & Subsystem:** `SYS-IMP-007` | Git Workflow | `cmd/gitmap/diff.go`
2. **System-Level UX Impact:** Allows rapid visual inspection of pending changes before atomic commits (`cpf`/`cpb`).
3. **Cause:** Standard `git diff` produces overwhelming walls of text across multiple micro-batches.
4. **Effect:** GitMap folds unchanged files, highlights syntax-aware diffs, and displays per-file summary bars.
5. **Zero Harm / Regression Safety:** Read-only inspection of git index and working tree.
6. **Backward Compatibility Guarantee:** Standard diff output available via `--compact`.
7. **Concrete Implementation Blueprint:**
   ```go
   func RenderDiffSummary(files []ModifiedFile) {
       for _, f := range files {
           fmt.Printf("%s | +%d -%d %s\n", f.Path, f.Additions, f.Deletions, makeBar(f))
       }
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap diff --stat` runs in <50ms and shows clean file additions/deletions.

#### SYS-IMP-008: Custom Command Alias Expansion with Parameter Forwarding
1. **Identifier & Subsystem:** `SYS-IMP-008` | CLI Extensibility | `pkg/config/aliases.go`
2. **System-Level UX Impact:** Reduces keystrokes for frequent multi-flag operations.
3. **Cause:** Typing verbose command lines like `gitmap aum search "symbol" . -e .go -r -i` slows down developer workflows.
4. **Effect:** Developers define aliases like `gitmap alias set sg "aum search $1 . -e .go -r"` for instant invocation.
5. **Zero Harm / Regression Safety:** Aliases stored in local `gitmap.config.json`; cannot shadow native core subcommands.
6. **Backward Compatibility Guarantee:** Native commands take precedence over any alias.
7. **Concrete Implementation Blueprint:**
   ```go
   if target, ok := cfg.Aliases[cmdName]; ok {
       return dispatchExpanded(target, os.Args[2:])
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap alias list` lists registered aliases and exits 0.

#### SYS-IMP-009: Context-Aware Dynamic Shell Prompt Hook
1. **Identifier & Subsystem:** `SYS-IMP-009` | Terminal Telemetry | `cmd/gitmap/prompt.go`
2. **System-Level UX Impact:** Displays active repository, pending task count, and heatmap status in shell prompt.
3. **Cause:** Developers lose track of which repo or pending subtask they are currently operating in.
4. **Effect:** Ultra-fast sub-millisecond hook outputs `[gitmap: 3 pending | heatmap: active]`.
5. **Zero Harm / Regression Safety:** Execution time hard-capped at 5ms; returns empty string on timeout.
6. **Backward Compatibility Guarantee:** Pure opt-in shell prompt helper.
7. **Concrete Implementation Blueprint:**
   ```go
   func GetPromptString() string {
       ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
       defer cancel()
       return renderPromptTelemetry(ctx)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap prompt-info` returns formatted status string in <5ms.

#### SYS-IMP-010: Seamless Workspace Root Discovery & Cross-Drive Path Translation
1. **Identifier & Subsystem:** `SYS-IMP-010` | Filesystem Navigation | `pkg/workspace/`
2. **System-Level UX Impact:** Automatically identifies git repo root regardless of current nested subfolder depth.
3. **Cause:** Running commands from deeply nested folders like `src/components/ui/` requires awkward `../../` path acrobatics.
4. **Effect:** GitMap automatically walks up parent directories, locates `.git`, and anchors relative operations from repo root.
5. **Zero Harm / Regression Safety:** Preserves relative path contracts; does not mutate working directory.
6. **Backward Compatibility Guarantee:** Explicit paths passed via `--path` override auto-discovery.
7. **Concrete Implementation Blueprint:**
   ```go
   func DiscoverRepoRoot(startDir string) (string, *appfault.AppError) {
       for curr := startDir; curr != filepath.Dir(curr); curr = filepath.Dir(curr) {
           if _, err := os.Stat(filepath.Join(curr, ".git")); err == nil { return curr, nil }
       }
       return "", appfault.New("NOT_A_GIT_REPO", "No .git found in parent hierarchy")
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap root` from any nested subdirectory outputs the exact repository root.

---

### Cluster 02: UI/UX, Design Tokens & Materiality (SYS-IMP-011 – SYS-IMP-020)

#### SYS-IMP-011: Centralized Design Token Registry & CSS Variables Sync
1. **Identifier & Subsystem:** `SYS-IMP-011` | Design System Architecture | `src/index.css` & `tailwind.config.ts`
2. **System-Level UX Impact:** Guarantees 100% visual consistency across colors, border-radii, and elevation tokens.
3. **Cause:** Design tokens are split across Tailwind config and raw CSS files, causing subtle visual drift.
4. **Effect:** Single source of truth CSS custom properties linked to Tailwind utility classes.
5. **Zero Harm / Regression Safety:** Pure CSS variable mapping; no HTML markup changes.
6. **Backward Compatibility Guarantee:** All existing utility class names (`bg-background`, `text-primary`) preserved.
7. **Concrete Implementation Blueprint:**
   ```css
   :root {
     --surface-ground: 240 10% 3.9%;
     --surface-elevated: 240 10% 6.5%;
     --surface-border: 240 3.7% 15.9%;
     --radius-sm: 0.25rem;
     --radius-md: 0.5rem;
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap aum search "--surface-" src/index.css` returns all defined surface tokens.

#### SYS-IMP-012: High-Contrast Accessible Color Ramps (WCAG AAA)
1. **Identifier & Subsystem:** `SYS-IMP-012` | Accessibility & Theming | `src/theme/`
2. **System-Level UX Impact:** Ensures all text and interactive elements achieve >= 7:1 contrast ratio, eliminating readability eye-strain.
3. **Cause:** Muted text shades on dark backgrounds frequently dip below 4.5:1 contrast, failing accessibility checks.
4. **Effect:** Contrast ramps strictly audited to WCAG AAA standards across both light and dark modes.
5. **Zero Harm / Regression Safety:** Non-breaking color token values; layout geometries untouched.
6. **Backward Compatibility Guarantee:** Existing component props and Tailwind classes remain identical.
7. **Concrete Implementation Blueprint:**
   ```css
   .dark {
     --muted-foreground: 240 5% 75%; /* Contrast 8.2:1 against #0a0a0c */
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Automated contrast check verifies all foreground/background pairings exceed 7:1.

#### SYS-IMP-013: Fluid Dynamic Typography Scaling & Math-Driven Modular Clamp Scales
1. **Identifier & Subsystem:** `SYS-IMP-013` | Typography System | `src/styles/typography.css`
2. **System-Level UX Impact:** Headings and body copy scale smoothly across screen widths without sudden layout snapping.
3. **Cause:** Fixed media-query typography breakpoints create abrupt jumps during viewport resizing.
4. **Effect:** CSS `clamp(min, preferred, max)` modular formulas provide seamless proportional typography.
5. **Zero Harm / Regression Safety:** Applied via utility classes; no element nesting changes.
6. **Backward Compatibility Guarantee:** Fallback fixed font sizes provided for legacy browser engines.
7. **Concrete Implementation Blueprint:**
   ```css
   --font-size-h1: clamp(2rem, 1.5rem + 2.5vw, 3.5rem);
   --font-size-body: clamp(0.95rem, 0.9rem + 0.25vw, 1.1rem);
   ```
8. **Binary Verification & Acceptance Criteria:** Resizing browser window from 320px to 2560px produces continuous, smooth font scaling.

#### SYS-IMP-014: Zero-Jank Virtualized Tree Renderer for Massive Spec Navigators
1. **Identifier & Subsystem:** `SYS-IMP-014` | Navigation Component | `src/components/SpecTreeNavigator.tsx`
2. **System-Level UX Impact:** Eliminates browser lag and frame drops when expanding large specification and plan trees (>10,000 nodes).
3. **Cause:** Rendering thousands of unvirtualized DOM tree nodes degrades scroll performance to <15 FPS.
4. **Effect:** Windowed DOM virtualization renders only visible viewport rows, maintaining rock-solid 60 FPS scrolling.
5. **Zero Harm / Regression Safety:** Encapsulated in the tree viewer component; external navigation events remain unchanged.
6. **Backward Compatibility Guarantee:** Identical tree node selection callback signatures.
7. **Concrete Implementation Blueprint:**
   ```tsx
   export function VirtualSpecTree({ nodes, onSelect }: SpecTreeProps) {
     return <VirtualList items={flattenVisible(nodes)} rowHeight={28} renderRow={TreeItem} />;
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Scrolling through a 5,000-node spec tree maintains 60 FPS in Chrome DevTools performance audit.

#### SYS-IMP-015: Avant-Garde Glassmorphic Elevation & Surface Depth Hierarchy
1. **Identifier & Subsystem:** `SYS-IMP-015` | UI Component Styling | `src/components/ui/`
2. **System-Level UX Impact:** Delivers modern visual aesthetics with tangible spatial hierarchy and clear layered depth.
3. **Cause:** Flat, opaque cards lack depth perception, making complex modal and drawer overlays look cluttered.
4. **Effect:** Layered surfaces use backdrop blur (12px), subtle translucent borders, and directional ambient drop-shadows.
5. **Zero Harm / Regression Safety:** Strictly visual CSS styling; zero DOM structure mutations.
6. **Backward Compatibility Guarantee:** Degrades gracefully to solid backgrounds on low-power devices.
7. **Concrete Implementation Blueprint:**
   ```css
   .glass-surface {
     background: rgba(18, 18, 24, 0.75);
     backdrop-filter: blur(12px);
     border: 1px solid rgba(255, 255, 255, 0.08);
     box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Inspection of card and modal elements shows proper `backdrop-filter` and border styling.

#### SYS-IMP-016: Real-Time Markdown Live Preview with Synchronized Scroll Spy
1. **Identifier & Subsystem:** `SYS-IMP-016` | Markdown Viewer | `src/pages/DocsViewer.tsx`
2. **System-Level UX Impact:** Keeps source editor and rendered preview perfectly in sync during navigation and reading.
3. **Cause:** Editor and preview desynchronize on long documents with large code blocks, disorienting readers.
4. **Effect:** Heading-based proportional scroll spy matches source lines to preview DOM elements within 2px tolerance.
5. **Zero Harm / Regression Safety:** Read-only scroll calculation event listeners with passive flags.
6. **Backward Compatibility Guarantee:** Standalone preview mode works without synchronization active.
7. **Concrete Implementation Blueprint:**
   ```typescript
   function syncScroll(sourceElem: HTMLElement, targetElem: HTMLElement, ratio: number) {
     targetElem.scrollTop = ratio * (targetElem.scrollHeight - targetElem.clientHeight);
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Scrolling the docs pane accurately reflects active heading in the sidebar index.

#### SYS-IMP-017: Adaptive Keyboard Navigation & Focus Ring Design Tokens
1. **Identifier & Subsystem:** `SYS-IMP-017` | Accessibility & Navigation | `src/hooks/useDocsKeyboard.ts`
2. **System-Level UX Impact:** Power users navigate specs, prompts, and plans entirely via keyboard shortcuts without touching the mouse.
3. **Cause:** Missing focus rings and keyboard hooks make modal dialogs and file navigation mouse-dependent.
4. **Effect:** Global `Cmd+K` palette, `J/K` list navigation, `Enter` to open, and high-visibility 2px offset focus rings.
5. **Zero Harm / Regression Safety:** Keyboard listeners bounded by component mount lifecycle; event propagation cleanly managed.
6. **Backward Compatibility Guarantee:** Mouse and touch navigation remain 100% functional.
7. **Concrete Implementation Blueprint:**
   ```css
   :focus-visible {
     outline: 2px solid hsl(var(--primary));
     outline-offset: 2px;
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Navigating DocsViewer using `Tab`, `J`, and `K` highlights active items with clear focus rings.

#### SYS-IMP-018: Micro-Interactive Control Motion Dynamics
1. **Identifier & Subsystem:** `SYS-IMP-018` | Animation System | `src/styles/animations.css`
2. **System-Level UX Impact:** Buttons and toggles feel responsive and tangible with physics-inspired spring micro-interactions.
3. **Cause:** Abrupt state transitions feel clunky; excessive linear animations feel sluggish.
4. **Effect:** 150ms cubic-bezier spring curves (`cubic-bezier(0.16, 1, 0.3, 1)`) with GPU-accelerated `transform` and `opacity`.
5. **Zero Harm / Regression Safety:** Honors `prefers-reduced-motion` media queries; disables animations instantly for users requesting it.
6. **Backward Compatibility Guarantee:** Component state handlers and click events execute with zero delay.
7. **Concrete Implementation Blueprint:**
   ```css
   @media (prefers-reduced-motion: no-preference) {
     .btn-interactive:active { transform: scale(0.97); transition: transform 0.1s cubic-bezier(0.16, 1, 0.3, 1); }
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Interacting with buttons triggers immediate hardware-accelerated scale feedback.

#### SYS-IMP-019: Interactive Slide Deck Presenter HUD with Dual-Screen Teleprompter
1. **Identifier & Subsystem:** `SYS-IMP-019` | Presentation Engine | `src/components/SlidePresenterHUD.tsx`
2. **System-Level UX Impact:** Seamless dual-screen slide delivery with elapsed timer, presenter speaker notes, and next-slide preview.
3. **Cause:** Standard full-screen presentation decks leave speakers blind to notes and upcoming content.
4. **Effect:** Presenter window synchronizes bidirectionally via `BroadcastChannel` with audience window in real time.
5. **Zero Harm / Regression Safety:** Confined to presentation routes; zero impact on core documentation viewers.
6. **Backward Compatibility Guarantee:** Standalone 16:9 canvas mode works without presenter window.
7. **Concrete Implementation Blueprint:**
   ```typescript
   const channel = new BroadcastChannel("slides-sync");
   channel.postMessage({ type: "SLIDE_CHANGE", slideIndex: activeIndex });
   ```
8. **Binary Verification & Acceptance Criteria:** Opening presenter mode transmits slide index changes to audience window with <10ms latency.

#### SYS-IMP-020: Dynamic Theme Preset Switcher with Zero-FOUC Guard
1. **Identifier & Subsystem:** `SYS-IMP-020` | Theme Hydration | `index.html` & `src/lib/theme.ts`
2. **System-Level UX Impact:** Eliminates the jarring white flash of unstyled content (FOUC) on initial application page loads in dark mode.
3. **Cause:** Theme state hydrated asynchronously via React `useEffect` after initial HTML render causes a brief white screen flicker.
4. **Effect:** Synchronous inline `<script>` in `<head>` reads `localStorage` and injects `.dark` class before the first paint.
5. **Zero Harm / Regression Safety:** Standard HTML5 `<head>` script pattern; zero runtime dependencies.
6. **Backward Compatibility Guarantee:** Graceful fallback to `prefers-color-scheme` if `localStorage` is empty.
7. **Concrete Implementation Blueprint:**
   ```html
   <script>
     (function() {
       var t = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
       if (t === 'dark') document.documentElement.classList.add('dark');
     })();
   </script>
   ```
8. **Binary Verification & Acceptance Criteria:** Hard page reload in dark mode renders instantly dark with zero white flash frames.

---

### Cluster 03: Autonomous Agent Orchestration & Heatmap Telemetry (SYS-IMP-021 – SYS-IMP-030)

#### SYS-IMP-021: Sub-Millisecond SQLite Agent Task Queue & Non-Blocking WAL Claim Engine
1. **Identifier & Subsystem:** `SYS-IMP-021` | Agent Engine | `cmd/gitmap/agent_task.go`
2. **System-Level UX Impact:** Eliminates subagent task claiming collisions and lock timeouts during high-concurrency waves.
3. **Cause:** Multiple subagents querying task SQLite database without WAL mode trigger `database is locked` exit errors.
4. **Effect:** Native compiled Go SQLite task manager operates with WAL mode and `PRAGMA busy_timeout=5000;`, claiming tasks in <1ms.
5. **Zero Harm / Regression Safety:** Atomic transactions isolate subtask claims; rollback on collision.
6. **Backward Compatibility Guarantee:** CLI interface matches standard `gitmap agent subtask claim --id <id>`.
7. **Concrete Implementation Blueprint:**
   ```go
   tx, _ := db.Begin()
   res, _ := tx.Exec("UPDATE subtasks SET is_claimed=1, claimed_by=? WHERE id=? AND is_claimed=0", agentID, subtaskID)
   if rows, _ := res.RowsAffected(); rows == 0 { tx.Rollback(); return ErrAlreadyClaimed }
   tx.Commit()
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap agent task status` runs in <2ms and verifies WAL journal mode.

#### SYS-IMP-022: Heatmap Cognitive Action Conflict Detection & Dynamic File Lease Locks
1. **Identifier & Subsystem:** `SYS-IMP-022` | Concurrency Safety | `pkg/heatmap/lease.go`
2. **System-Level UX Impact:** Prevents race conditions and Windows file lock crashes caused by two agents writing the same file.
3. **Cause:** Independent subagents simultaneously modifying adjacent files can collide on shared imports or shared types.
4. **Effect:** Heatmap pre-action logging inspects active leases; aborts or queues requests if file is leased to another agent.
5. **Zero Harm / Regression Safety:** Non-blocking advisory lock table in SQLite; auto-expires after 120s if agent crashes.
6. **Backward Compatibility Guarantee:** Bypassed in single-agent CLI runs.
7. **Concrete Implementation Blueprint:**
   ```go
   func AcquireFileLease(filePath, agentID string) bool {
       return db.Exec("INSERT INTO file_leases (file_path, agent_id, expires_at) VALUES (?, ?, ?)", filePath, agentID, time.Now().Add(2*time.Minute)) == nil
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap agent heatmap` reports zero concurrent write collisions across active waves.

#### SYS-IMP-023: Autonomous Deadlock Autopsy & Dynamic Subagent Respawn Controller
1. **Identifier & Subsystem:** `SYS-IMP-023` | Swarm Resilience | `cmd/gitmap/agent_crashed.go`
2. **System-Level UX Impact:** Automatically recovers abandoned subtasks when a subagent fails or times out, preventing stalled runs.
3. **Cause:** Subagent process crashes leave subtasks in `IN_PROGRESS` state indefinitely, hanging the parent orchestrator.
4. **Effect:** Health checker identifies stale workers (>60s inactivity), flags status as `ABANDONED`, and re-enqueues subtask.
5. **Zero Harm / Regression Safety:** Safely cleans stale in-memory state; rolls back uncommitted file edits.
6. **Backward Compatibility Guarantee:** Parent orchestrator receives structured `RESUME_FOUND` event.
7. **Concrete Implementation Blueprint:**
   ```go
   func ReclaimAbandonedTasks() int {
       res, _ := db.Exec("UPDATE subtasks SET status='PENDING', claimed_by='' WHERE status='IN_PROGRESS' AND updated_at < datetime('now', '-60 seconds')")
       n, _ := res.RowsAffected(); return int(n)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap agent crashed` identifies stalled task and successfully returns it to pending queue.

#### SYS-IMP-024: Real-Time Agent Activity Visualizer & Live Terminal Tree Matrix
1. **Identifier & Subsystem:** `SYS-IMP-024` | Observability & Telemetry | `cmd/gitmap/heatmap_live.go`
2. **System-Level UX Impact:** Gives developers complete real-time visibility into what every subagent is reading, writing, and testing.
3. **Cause:** Multi-agent operations run in background opacity; developers cannot tell if workers are making progress or looping.
4. **Effect:** `gitmap heatmap live` streams live subagent actions, target files, and elapsed times in a structured terminal matrix.
5. **Zero Harm / Regression Safety:** Pure streaming reader from SQLite WAL; zero mutation to task states.
6. **Backward Compatibility Guarantee:** Standalone diagnostic command.
7. **Concrete Implementation Blueprint:**
   ```go
   func StreamLiveHeatmap(ctx context.Context) {
       for range time.Tick(250 * time.Millisecond) {
           renderActiveAgentsMatrix(fetchActiveTelemetry())
       }
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap heatmap live` displays active agent rows with accurate live file targets.

#### SYS-IMP-025: Multi-Agent Memory Exchange Bus & Incremental Diff Ledger
1. **Identifier & Subsystem:** `SYS-IMP-025` | Inter-Agent Communication | `.ai-memory/temp-agents/`
2. **System-Level UX Impact:** Allows subagents to share discovered symbol locations and schemas without polluting global chat context.
3. **Cause:** Subagents re-discovering the same file boundaries waste token context and execution steps.
4. **Effect:** Lead orchestrator populates shared ledger with discovered symbol maps; subagents ingest ledger during Step 0.
5. **Zero Harm / Regression Safety:** Stored in temporary agent directory; purged upon task completion.
6. **Backward Compatibility Guarantee:** Standard markdown ledger schema preserved.
7. **Concrete Implementation Blueprint:**
   ```markdown
   | Symbol | Source File | Exported Types | Known Callers |
   | AppError | pkg/appfault/fault.go | AppError, Fault | cmd/gitmap/, pkg/sync/ |
   ```
8. **Binary Verification & Acceptance Criteria:** Subagents cite shared ledger mappings without running redundant file search tools.

#### SYS-IMP-026: Capacity Saturation Optimizer & Dynamic Hand Dispatch
1. **Identifier & Subsystem:** `SYS-IMP-026` | Task Scheduling | `pkg/scheduler/`
2. **System-Level UX Impact:** Keeps system concurrency capacity pinned at 100% duty cycle (`A = 3, H = 2`), cutting total run times by 40%.
3. **Cause:** Orchestrator waiting for full 3-agent wave to finish leaves early-finishing workers idle.
4. **Effect:** As soon as any worker completes, scheduler immediately dispatches next pending subtask to available hand.
5. **Zero Harm / Regression Safety:** Work-stealing queue operates on disjoint file partitions only.
6. **Backward Compatibility Guarantee:** Total wave budget calculation remains mathematically identical.
7. **Concrete Implementation Blueprint:**
   ```typescript
   onWorkerCompleted(workerId: string) {
       const nextSubtask = queue.dequeue();
       if (nextSubtask) dispatchToWorker(workerId, nextSubtask);
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Execution logs show zero idle gaps between subtask completions and subsequent assignments.

#### SYS-IMP-027: Pre-Action Intent Logging & Post-Action Outcome Verification Gate
1. **Identifier & Subsystem:** `SYS-IMP-027` | Cognitive Telemetry | `cmd/gitmap/heatmap_log.go`
2. **System-Level UX Impact:** Ensures every agent file read or edit has an explicit documented goal and verified outcome.
3. **Cause:** Unguided file reads waste steps without synthesizing findings or proving whether the goal was met.
4. **Effect:** Enforces two-phase telemetry: `gitmap heatmap log --action READ --expected "goal"` paired with `gitmap heatmap update --outcome "met"`.
5. **Zero Harm / Regression Safety:** Lightweight logging overhead (<1ms); does not interrupt code execution.
6. **Backward Compatibility Guarantee:** Bypassed when `--no-telemetry` flag is supplied.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap heatmap log --task "$TASK" --action READ --file "$PATH" --expected "Verify AppError signature"
   # ... perform action ...
   gitmap heatmap update --id "$LOG_ID" --outcome "Signature verified" --met true
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap heatmap status --task <task>` reports 100% action closure ratio.

#### SYS-IMP-028: Automatic Temporary Artifact Sandbox & Session Temp Directory Auto-Pruning
1. **Identifier & Subsystem:** `SYS-IMP-028` | Workspace Hygiene | `cmd/gitmap/agent_temp_clear.go`
2. **System-Level UX Impact:** Keeps working tree pristine and prevents accidental commit of temporary test scratch files.
3. **Cause:** Temporary diagnostic scripts left in `.ai-memory/temp-agents/` clutter git status and confuse subsequent runs.
4. **Effect:** `gitmap agent temp-clear` safely removes session databases and scratch directories after final push verification.
5. **Zero Harm / Regression Safety:** Confined strictly to `.ai-memory/temp-agents/`; never touches persistent plans or specs.
6. **Backward Compatibility Guarantee:** Persistent artifacts remain untouched.
7. **Concrete Implementation Blueprint:**
   ```go
   func CleanTempAgents(slug string) *appfault.AppError {
       tempDir := filepath.Join(".ai-memory", "temp-agents", slug)
       return os.RemoveAll(tempDir)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running `gitmap agent temp-clear` clears session directory; `git status --porcelain` is clean.

#### SYS-IMP-029: Dynamic Subagent Context Window Governor & Token Throttling
1. **Identifier & Subsystem:** `SYS-IMP-029` | Prompt Dispatcher | `pkg/agent/context.go`
2. **System-Level UX Impact:** Prevents subagents from exhausting context limits by stripping unnecessary historical conversational baggage.
3. **Cause:** Passing full multi-turn conversation histories to workers causes prompt truncation and slow model inference.
4. **Effect:** Worker briefs are filtered to include strictly: (1) task objective, (2) owned files, (3) relevant spec snippet, (4) active rules.
5. **Zero Harm / Regression Safety:** Workers receive complete instructions within a compact 1,500-token footprint.
6. **Backward Compatibility Guarantee:** Parent orchestrator retains full conversation memory.
7. **Concrete Implementation Blueprint:**
   ```typescript
   function buildCompactWorkerBrief(subtask: Subtask, specSnippet: string): string {
     return `# Worker Brief: ${subtask.title}\n\nOwned Files: ${subtask.files.join(", ")}\n\nSpec Requirements:\n${specSnippet}`;
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Dispatched subagent prompts stay strictly under 2,000 tokens with zero lost context.

#### SYS-IMP-030: Multi-Agent Consensus Verification Protocol for Boundary Decisions
1. **Identifier & Subsystem:** `SYS-IMP-030` | Architecture Governance | `pkg/governance/`
2. **System-Level UX Impact:** Prevents rogue architectural modifications by requiring explicit cross-verification on shared core packages.
3. **Cause:** A single subagent refactoring a shared package like `pkg/appfault` can inadvertently break other workers' code.
4. **Effect:** Changes to core shared packages require consensus agreement before the lead orchestrator stages the modification.
5. **Zero Harm / Regression Safety:** Confined to multi-repo shared libraries; application specs remain independent.
6. **Backward Compatibility Guarantee:** Core API signatures preserved with non-breaking type aliases.
7. **Concrete Implementation Blueprint:**
   ```go
   if isSharedCorePackage(pkgPath) {
       requireLeadReviewAndConsensus(pkgPath, proposedDiff)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Attempting unauthorized direct modification to shared core emits boundary warning.

---

### Cluster 04: Prompt Engineering & Instruction Architecture (SYS-IMP-031 – SYS-IMP-040)

#### SYS-IMP-031: Canonical V6 Parameterized Prompt Template Engine
1. **Identifier & Subsystem:** `SYS-IMP-031` | Prompt Architecture | `01-prompts/`
2. **System-Level UX Impact:** Standardizes execution behavior across all AI agents through mathematical header parameters (`N, A, H, C`).
3. **Cause:** Hardcoded step numbers and inconsistent agent parameters across prompt files cause unpredictable execution flows.
4. **Effect:** Header variables drive all budget allocations; zero magic numbers in prompt instruction bodies.
5. **Zero Harm / Regression Safety:** Pure prompt template standardization; no runtime code changes.
6. **Backward Compatibility Guarantee:** Compatible with Google Antigravity 2.0 and Cursor IDE runtimes.
7. **Concrete Implementation Blueprint:**
   ```text
   N = 400 (Total self-loop steps budget — editable top-header parameter)
   A = 3   (MANDATORY number of spawned autonomous subagents)
   H = 2   (Operational hands per agent: dual-task batch capacity)
   C = 30  (Tool calls per worker before it must report)
   ```
8. **Binary Verification & Acceptance Criteria:** Prompt validation linter confirms all V6 parameters exist and are mathematically consistent.

#### SYS-IMP-032: Top-Instruction Priority Verification Protocol & Preamble Precedence
1. **Identifier & Subsystem:** `SYS-IMP-032` | Instruction Invariants | `01-prompts/`
2. **System-Level UX Impact:** Guarantees user instructions placed above the prompt take absolute precedence over template defaults.
3. **Cause:** AI models sometimes prioritize general prompt instructions over specific user constraints in incoming preambles.
4. **Effect:** Strict Top-Instruction Priority mandate instructs agents to evaluate preamble constraints before executing pipeline.
5. **Zero Harm / Regression Safety:** Ensures user custom constraints are never ignored or overridden.
6. **Backward Compatibility Guarantee:** Default prompt behavior operates normally when no preamble overrides are present.
7. **Concrete Implementation Blueprint:**
   ```markdown
   > **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
   > Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt
   > are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.
   ```
8. **Binary Verification & Acceptance Criteria:** Ingestion test confirms preamble custom constraints override default parameters.

#### SYS-IMP-033: IDE-Native Interactive Skill Link Synthesizer
1. **Identifier & Subsystem:** `SYS-IMP-033` | Skill Formatting | `.agents/skills/` & `.cursor/skills/`
2. **System-Level UX Impact:** Enables one-click IDE navigation and skill loading directly from chat and prompt markdown.
3. **Cause:** Using raw filesystem paths (`file:///...`) or appending `/skill.md` breaks IDE interactive skill detection.
4. **Effect:** Formats links strictly as `[<skill>](file;.agents/skills/<skill>)` for Antigravity and `(file;.cursor/skills/<skill>)` for Cursor.
5. **Zero Harm / Regression Safety:** Replaces non-clickable text with interactive links; zero filesystem changes.
6. **Backward Compatibility Guarantee:** Renders as clean readable markdown in standard text viewers.
7. **Concrete Implementation Blueprint:**
   ```markdown
   [find-100-improvements-in-whole-system](file;.agents/skills/find-100-improvements-in-whole-system)
   [/plan](slashCommand;plan)
   ```
8. **Binary Verification & Acceptance Criteria:** Skill path linter confirms zero occurrences of `.md` inside link URI targets.

#### SYS-IMP-034: Lossless Verbatim User Request Capture & Discrete Deliverable Extraction
1. **Identifier & Subsystem:** `SYS-IMP-034` | Requirements Ingestion | `pkg/spec/verbatim.go`
2. **System-Level UX Impact:** Completely prevents scope creep, missed requirements, or hallucinated specifications.
3. **Cause:** Agents summarizing user prompts into fuzzy outlines omit subtle constraints and technical criteria.
4. **Effect:** Phase 1A captures incoming user request losslessly under `## User Request (Verbatim)` and extracts numbered task IDs.
5. **Zero Harm / Regression Safety:** Preserves exact user intent in canonical specification files.
6. **Backward Compatibility Guarantee:** Downstream review scripts verify exact match against original user input.
7. **Concrete Implementation Blueprint:**
   ```markdown
   ## User Request (Verbatim)
   > [Exact unedited prompt content captured byte-for-byte]
   ```
8. **Binary Verification & Acceptance Criteria:** Spec audit confirms verbatim section matches incoming prompt text with zero omissions.

#### SYS-IMP-035: Negative Guardrail Assertion Matrix & Anti-Hallucination Gate
1. **Identifier & Subsystem:** `SYS-IMP-035` | Guardrail Governance | `02-spec/02-coding-guidelines/`
2. **System-Level UX Impact:** Eliminates common agent failure modes (running `git` in subagents, raw grep, unauthorized full builds).
3. **Cause:** AI models default to habituated shell habits (`git add`, `npm test`, `grep`) unless explicitly barred.
4. **Effect:** Structured negative guardrail checklist executed prior to tool dispatch; violations trigger auto-reject failure.
5. **Zero Harm / Regression Safety:** Prohibits hazardous commands; protects index locks and terminal stability.
6. **Backward Compatibility Guarantee:** Legitimate GitMap commands remain fully enabled.
7. **Concrete Implementation Blueprint:**
   ```markdown
   - [ ] Worker Git Ban: Subagents NEVER run git commands (`git add`, `git commit`, `git push`).
   - [ ] No Full Builds: Never execute full build commands during active iteration.
   - [ ] Raw Search Ban: Banned all `rg`, `ripgrep`, `grep`, `Select-String`, and `findstr`.
   ```
8. **Binary Verification & Acceptance Criteria:** Pre-commit verification checks confirm zero prohibited commands were executed.

#### SYS-IMP-036: Below-Steps Precedence Architecture for Append-Only Tasks
1. **Identifier & Subsystem:** `SYS-IMP-036` | Prompt Mechanics | `01-prompts/15-cg-execute/`
2. **System-Level UX Impact:** Allows users to append custom code snippets and tasks directly to the bottom of prompt templates seamlessly.
3. **Cause:** Fixed prompt structures struggle to ingest arbitrary appended user code without formatting breakage.
4. **Effect:** Enforces Bottom-Instruction Priority: instructions appended after `## 🚨 High Priority Instructions Below` outrank all above.
5. **Zero Harm / Regression Safety:** Standardized separator boundary prevents parsing ambiguities.
6. **Backward Compatibility Guarantee:** Standard prompt templates without appended instructions execute normally.
7. **Concrete Implementation Blueprint:**
   ```markdown
   --
   ## 🚨 High Priority Instructions Below
   [Appended user tasks and specific refactoring targets go here]
   ```
8. **Binary Verification & Acceptance Criteria:** Parser correctly identifies bottom boundary and applies trailing instructions with highest priority.

#### SYS-IMP-037: Prompt MD5 Checksum Verification & Resumable Session Ledger Re-hydration
1. **Identifier & Subsystem:** `SYS-IMP-037` | Task State Persistence | `cmd/gitmap/agent_task.go`
2. **System-Level UX Impact:** Instantly resumes interrupted runs from the exact uncompleted subtask without re-doing completed work.
3. **Cause:** Re-invoking a prompt after a network hiccup or crash often restarts execution from scratch, wasting tokens and time.
4. **Effect:** Compares prompt MD5 hash with stored task database record; resumes seamlessly if identical.
5. **Zero Harm / Regression Safety:** Preserves already committed file modifications and completed subtask records.
6. **Backward Compatibility Guarantee:** Creates new clean task ID if prompt hash differs.
7. **Concrete Implementation Blueprint:**
   ```go
   if task.PromptHash == currentMD5 {
       return resumeTask(task.ID, task.FirstUncompletedSubtask())
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Re-running command with same prompt outputs `RESUME_FOUND` and resumes active subtask.

#### SYS-IMP-038: Structured 5-Pillar Turn-1 Showcase Generator
1. **Identifier & Subsystem:** `SYS-IMP-038` | Agent Communication | `01-prompts/`
2. **System-Level UX Impact:** Proves deep AI comprehension in Turn 1, giving developers complete confidence before execution starts.
3. **Cause:** Superficial 1-sentence acknowledgments ("I'll do that") fail to catch misunderstandings until work is already corrupted.
4. **Effect:** Turn 1 mandatory display of Strategic Intent, Boundaries, Guardrails, Blast Radius, and Multi-Agent Concurrency.
5. **Zero Harm / Regression Safety:** Visible text output in chat; executed in same turn with initial file tool call.
6. **Backward Compatibility Guarantee:** Standardized markdown header format.
7. **Concrete Implementation Blueprint:**
   ```markdown
   ### 🧠 Deep Comprehension & Understood Mandate
   #### Pillar 1: Strategic Intent & Systems Reasoning
   #### Pillar 2: Explicit Scope Boundaries (In-Scope vs. Excluded)
   #### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
   #### Pillar 4: Blast Radius & Dependency Topology
   #### Pillar 5: Multi-Agent Concurrency Topology
   ```
8. **Binary Verification & Acceptance Criteria:** Turn 1 chat output contains all 5 mandatory pillars with explicit `[YES]` comprehension flags.

#### SYS-IMP-039: Automated Prompt Token Density & Redundancy Reduction
1. **Identifier & Subsystem:** `SYS-IMP-039` | Context Optimization | `01-prompts/`
2. **System-Level UX Impact:** Frees up to 30% more active context window for actual code changes and file reading.
3. **Cause:** Duplicating entire rule texts inside every prompt file wastes precious model context and increases latency.
4. **Effect:** Standardizes rule references by ID (`R1` through `R16`) and links directly to canonical specification files.
5. **Zero Harm / Regression Safety:** Preserves full rule semantics while eliminating verbose duplication.
6. **Backward Compatibility Guarantee:** Agents follow canonical spec links whenever rule verification is needed.
7. **Concrete Implementation Blueprint:**
   ```markdown
   - Rule R1: Boolean Standards ([`02-spec/02-coding-guidelines/02-booleans.md`](02-spec/02-coding-guidelines/02-booleans.md))
   - Rule R5: Strict Relative Git Paths ([`02-spec/02-coding-guidelines/05-relative-paths.md`](02-spec/02-coding-guidelines/05-relative-paths.md))
   ```
8. **Binary Verification & Acceptance Criteria:** Word count linter confirms prompt body is optimized under 3,200 words.

#### SYS-IMP-040: Dynamic Slash-Command Suggestion Routing & Context Advisor
1. **Identifier & Subsystem:** `SYS-IMP-040` | Interactive UX | `pkg/advisor/`
2. **System-Level UX Impact:** Proactively recommends relevant workflow shortcuts (`[/plan]`, `[/goal]`, `[/learn]`) based on user intent.
3. **Cause:** Developers often unaware of specialized slash commands execute complex multi-step tasks manually.
4. **Effect:** Agent suggests the optimal slash command at the conclusion of turns to streamline subsequent workflows.
5. **Zero Harm / Regression Safety:** Informational recommendation at the bottom of messages; no automated side effects.
6. **Backward Compatibility Guarantee:** User retains full control to accept or disregard recommendations.
7. **Concrete Implementation Blueprint:**
   ```markdown
   You can use the `[/plan](slashCommand;plan)` command to decompose complex architecture before execution.
   ```
8. **Binary Verification & Acceptance Criteria:** Responses recommend slash commands using validated IDE link syntax.

---

### Cluster 05: Multi-Repository Fleet Synchronization & Boundaries (SYS-IMP-041 – SYS-IMP-050)

#### SYS-IMP-041: Parallel Multi-Repo Synchronization Engine (`gitmap sync`)
1. **Identifier & Subsystem:** `SYS-IMP-041` | Fleet Management | `cmd/gitmap/sync.go`
2. **System-Level UX Impact:** Synchronizes canonical prompts, skills, and shared specs across all 43 repositories in <5 seconds.
3. **Cause:** Sequential synchronization scripts take over 2 minutes and fail halfway through if one repository has a network delay.
4. **Effect:** Parallel Go synchronization engine with 8 worker threads synchronizes all target repositories concurrently.
5. **Zero Harm / Regression Safety:** Pre-flight safety checks verify repository status before applying any changes.
6. **Backward Compatibility Guarantee:** Legacy fallback script `03-ai-scripts/38-sync-prompts-skills-scripts.py` preserved.
7. **Concrete Implementation Blueprint:**
   ```go
   func RunFleetSync(repos []string, workers int) *appfault.AppError {
       sem := make(chan struct{}, workers)
       for _, r := range repos { go syncWorker(r, sem) }
       return waitForSyncCompletion()
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap sync --workers 8` successfully verifies all 43 repos in <5s.

#### SYS-IMP-042: The 5 Non-Negotiable Boundary Enforcement Guardrails
1. **Identifier & Subsystem:** `SYS-IMP-042` | Fleet Isolation | `pkg/sync/boundaries.go`
2. **System-Level UX Impact:** Protects application-specific code and version configs from being overwritten during fleet synchronization.
3. **Cause:** Naive file mirroring can overwrite downstream private application specs (`02-spec/21-*`) or custom bump scripts.
4. **Effect:** Hardcoded boundary filters strictly prohibit synchronizing private application domains, bump manifests, or memory logs.
5. **Zero Harm / Regression Safety:** Absolute barrier preventing cross-repository pollution or data loss.
6. **Backward Compatibility Guarantee:** Shared specs (`02-spec/01-*` to `02-spec/20-*`) synchronize cleanly.
7. **Concrete Implementation Blueprint:**
   ```go
   func isBlockedSyncPath(relPath string) bool {
       return strings.HasPrefix(relPath, "02-spec/21-") || 
              strings.HasPrefix(relPath, ".ai-memory/plans/") || 
              isBumpScript(relPath)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Sync dry-run test confirms 0 files from `02-spec/21-*` are scheduled for copy.

#### SYS-IMP-043: Automated Pre-Flight Safety Backup Branch Ceremony
1. **Identifier & Subsystem:** `SYS-IMP-043` | Fleet Safety | `pkg/sync/backup.go`
2. **System-Level UX Impact:** Provides guaranteed zero-data-loss rollback capability before touching downstream repositories.
3. **Cause:** Synchronizing updates into a downstream repository with uncommitted or divergent state can corrupt branches.
4. **Effect:** Engine creates and pushes `backup/sync-<timestamp>` remote branch before modifying any target files.
5. **Zero Harm / Regression Safety:** Non-destructive branch creation; preserves previous commits untouched.
6. **Backward Compatibility Guarantee:** Base tracking branch remains clean and ready for fast-forward merge.
7. **Concrete Implementation Blueprint:**
   ```bash
   git checkout -b "backup/sync-$(date +%s)" && git push origin "backup/sync-$(date +%s)"
   git checkout main
   ```
8. **Binary Verification & Acceptance Criteria:** Remote git branch list confirms backup branch exists prior to sync modifications.

#### SYS-IMP-044: Zero Secrets Leakage Gate & Repo-Secrets Offload Protocol
1. **Identifier & Subsystem:** `SYS-IMP-044` | Security Governance | `cmd/gitmap/rs.go`
2. **System-Level UX Impact:** Automatically offloads tokens, private keys, and `.env` files into secure `repo-secrets` vault.
3. **Cause:** Developers or agents accidentally staging `.env` files into public git history causes severe security breaches.
4. **Effect:** Pre-commit hooks scan for credential patterns; automatically redirect and stage secrets into isolated `repo-secrets`.
5. **Zero Harm / Regression Safety:** Sensitive data completely removed from standard git worktrees.
6. **Backward Compatibility Guarantee:** Local development environments access secrets via unified `gitmap rs` CLI helpers.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap rs file .env.production --repo current-repo
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap find "*.*env*"` in working directory confirms zero raw secret files exist in git tree.

#### SYS-IMP-045: Downstream Version Manifest Isolation & Protection
1. **Identifier & Subsystem:** `SYS-IMP-045` | Release Independence | `pkg/sync/version_protect.go`
2. **System-Level UX Impact:** Allows each repository to progress independently through its own SemVer release lifecycle.
3. **Cause:** Fleet sync overwriting `version.json` or bump scripts resets downstream release versions to the parent repo's version.
4. **Effect:** Version manifests (`version.json`, `package.json#version`, `Cargo.toml`) are explicitly excluded from sync copying.
5. **Zero Harm / Regression Safety:** Downstream release numbers and changelogs remain 100% sovereign.
6. **Backward Compatibility Guarantee:** Upstream tooling compatibility maintained across differing version numbers.
7. **Concrete Implementation Blueprint:**
   ```go
   var ProtectedManifests = []string{"version.json", "version.template.json", "bump-version.mjs", "bump_versions.py"}
   ```
8. **Binary Verification & Acceptance Criteria:** Post-sync verification confirms downstream `version.json` retains its exact pre-sync version string.

#### SYS-IMP-046: Multi-Repo Divergence Telemetry & Drift Detection Dashboard
1. **Identifier & Subsystem:** `SYS-IMP-046` | Fleet Health | `cmd/gitmap/sync_status.go`
2. **System-Level UX Impact:** Provides immediate high-level summary of which repositories have fallen behind the canonical standard.
3. **Cause:** Repositories drift over time without developers knowing which prompts or shared guidelines are out of date.
4. **Effect:** `gitmap sync status --drift` outputs a color-coded matrix showing prompt and skill sync status across all 43 repos.
5. **Zero Harm / Regression Safety:** Read-only inspection; makes zero modifications.
6. **Backward Compatibility Guarantee:** Supports `--json` output for automated CI dashboard ingestion.
7. **Concrete Implementation Blueprint:**
   ```go
   func ReportFleetDrift() {
       for _, r := range fleet {
           fmt.Printf("%-24s | Prompts: %s | Skills: %s | Drift: %d files\n", r.Name, r.PromptStatus, r.SkillStatus, r.DriftCount)
       }
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap sync status --drift` executes in <3s and outputs complete fleet drift report.

#### SYS-IMP-047: Non-Owned Repository Exclusion Gate
1. **Identifier & Subsystem:** `SYS-IMP-047` | Multi-Repo Safety | `pkg/workspace/exclude.go`
2. **System-Level UX Impact:** Completely prevents accidental staging, mutation, or push to third-party or system repositories.
3. **Cause:** Recursive workspace scans picking up `omis`, `oh-my-zsh`, `zsh`, or `dotfiles` can accidentally commit files to external tools.
4. **Effect:** Hardcoded blacklist blocks all write operations on non-owned repositories.
5. **Zero Harm / Regression Safety:** Protects user system tools and open-source shells from accidental pollution.
6. **Backward Compatibility Guarantee:** Owned repositories operate normally.
7. **Concrete Implementation Blueprint:**
   ```go
   var BannedNonOwnedRepos = []string{"omis", "oh-my-zsh", "ohmyzsh", "zsh", "oh-my-posh", "dotfiles"}
   func IsNonOwnedRepo(name string) bool {
       return slices.Contains(BannedNonOwnedRepos, strings.ToLower(name))
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running `gitmap commit` inside an `omis` or `zsh` directory fails with clear exclusion error.

#### SYS-IMP-048: Bidirectional Skill & Prompt Mirroring with Automatic Relative Path Rewriting
1. **Identifier & Subsystem:** `SYS-IMP-048` | Content Portability | `pkg/sync/path_rewrite.go`
2. **System-Level UX Impact:** Ensures all markdown links inside copied skills and prompts work seamlessly in target repositories.
3. **Cause:** Relative paths between folders (e.g. `../../02-spec/`) break when files are synchronized into different directory nesting levels.
4. **Effect:** Engine analyzes relative link targets and recalculates relative path depth for target repository layout.
5. **Zero Harm / Regression Safety:** Preserves link destinations; strictly eliminates broken markdown links.
6. **Backward Compatibility Guarantee:** Canonical repository files remain unchanged.
7. **Concrete Implementation Blueprint:**
   ```go
   func RewriteRelativeLinks(content string, srcDepth, destDepth int) string {
       return recalculateMarkdownLinkDepth(content, destDepth - srcDepth)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running `03-ai-scripts/22-doc-path-linter.py` in target repo passes with zero broken relative links.

#### SYS-IMP-049: Fleet-Wide Clean Working Tree Verification Pre-Commit Check
1. **Identifier & Subsystem:** `SYS-IMP-049` | Git Hygiene | `cmd/gitmap/sync.go`
2. **System-Level UX Impact:** Prevents accidental overwriting of uncommitted developer work in downstream repositories.
3. **Cause:** Running sync while a developer is actively editing files in a downstream repo can cause merge conflicts or loss of work.
4. **Effect:** Checks `git status --porcelain` across all target repositories; skips dirty repositories with a clear warning.
5. **Zero Harm / Regression Safety:** Preserves active working tree files untouched.
6. **Backward Compatibility Guarantee:** `--force` flag available if explicit overwrite is commanded.
7. **Concrete Implementation Blueprint:**
   ```go
   if isWorkingTreeDirty(targetRepo) {
       log.Warn("Skipping repository due to uncommitted changes: " + targetRepo)
       continue
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Sync run against a dirty test repo logs skip warning and leaves working tree files intact.

#### SYS-IMP-050: Cross-Repository Reusable Script Cache Vault Synchronization (`gitmap rc sync`)
1. **Identifier & Subsystem:** `SYS-IMP-050` | Tooling Distribution | `cmd/gitmap/rc.go`
2. **System-Level UX Impact:** Shares diagnostic scripts, test harnesses, and PowerShell runners across all repositories without code duplication.
3. **Cause:** Useful debugging scripts created in one repository are unavailable in others without manual copying.
4. **Effect:** `gitmap rc` offloads reusable scripts to central `repo-cache`; `gitmap rc sync` pulls latest harnesses on demand.
5. **Zero Harm / Regression Safety:** Stored in isolated `repo-cache`; does not pollute standard git worktrees.
6. **Backward Compatibility Guarantee:** Standalone tooling utility.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap rc sync --pull
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap rc list` shows all available shared test harnesses and utility scripts.

---

### Cluster 06: Cross-Language Code Quality & Immutability (SYS-IMP-051 – SYS-IMP-060)

#### SYS-IMP-051: Monadic Result Envelopes & Domain `*appfault.AppError` Standards in Go
1. **Identifier & Subsystem:** `SYS-IMP-051` | Error Architecture | `pkg/appfault/` & `pkg/result/`
2. **System-Level UX Impact:** Eliminates cryptic nil-pointer dereferences and unhandled errors with strongly-typed result wrappers.
3. **Cause:** Functions returning raw multi-value tuples `(T, error)` allow developers to forget checking `err != nil`.
4. **Effect:** All Go functions return `result.Result[T]` or `*appfault.AppError` with structured code, message, and caller metadata.
5. **Zero Harm / Regression Safety:** Backward-compatible type shims (`type Fault = AppError`) preserve legacy callers.
6. **Backward Compatibility Guarantee:** `.AppError()` and `.Fault()` methods return identical structured error metadata.
7. **Concrete Implementation Blueprint:**
   ```go
   type AppError struct {
       Code    string `json:"code"`
       Message string `json:"message"`
       Caller  string `json:"caller"`
   }
   type Fault = AppError
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap aum search -r "func.*\*appfault\.AppError" pkg/ -e .go` verifies adherence.

#### SYS-IMP-052: Strict Positive Boolean Naming Standard (`is*` / `has*`)
1. **Identifier & Subsystem:** `SYS-IMP-052` | Boolean Hygiene | `02-spec/02-coding-guidelines/02-booleans.md`
2. **System-Level UX Impact:** Removes double-negation cognitive overload (`if !isNotReady`) from conditional logic.
3. **Cause:** Negative boolean flags (`disableCache`, `notFound`, `skipValidation`) produce confusing double-negatives.
4. **Effect:** 100% positive booleans prefixed strictly with `is` or `has` (e.g. `isActive`, `hasPermission`), evaluated implicitly.
5. **Zero Harm / Regression Safety:** Refactored via surgical AST updates; preserves original logical behavior.
6. **Backward Compatibility Guarantee:** Serialized JSON fields retain backward compatibility where external APIs mandate.
7. **Concrete Implementation Blueprint:**
   ```go
   // BEFORE: if !isDisabled == true
   // AFTER:
   if isEnabled {
       // logic
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running `python 03-ai-scripts/05-guideline-autofixer.py . --check-only` reports zero boolean violations.

#### SYS-IMP-053: Function Parameter Structs (`*Params`) Eliminating Loose Arguments
1. **Identifier & Subsystem:** `SYS-IMP-053` | Signature Refactoring | `pkg/`
2. **System-Level UX Impact:** Prevents argument transposition bugs (e.g. swapping `timeout` and `retries`) on functions with >2 arguments.
3. **Cause:** Long positional parameter lists (`func CreateUser(name, email, role, phone, address string, age int)`) are error-prone.
4. **Effect:** Functions with >2-3 arguments take a single typed `*Params` value struct with self-documenting named fields.
5. **Zero Harm / Regression Safety:** Refactored caller sites updated atomically; improves readability.
6. **Backward Compatibility Guarantee:** Factory constructor helpers provide convenient defaults.
7. **Concrete Implementation Blueprint:**
   ```go
   type CreateUserParams struct {
       Name  string
       Email string
       Role  UserRole
   }
   func CreateUser(params CreateUserParams) result.Result[*User]
   ```
8. **Binary Verification & Acceptance Criteria:** AST linter confirms zero function signatures with >3 positional parameters.

#### SYS-IMP-054: Domain Struct Centralization & Generics Extraction into `types.go`
1. **Identifier & Subsystem:** `SYS-IMP-054` | Package Architecture | `pkg/*/types.go`
2. **System-Level UX Impact:** Makes package domain data contracts immediately discoverable in a single dedicated file.
3. **Cause:** Structs and type aliases scattered across 15 different implementation files force developers to hunt for definitions.
4. **Effect:** Domain structs, generic aliases, and Result envelopes extracted into package-root `types.go`.
5. **Zero Harm / Regression Safety:** Pure code organization; zero alteration to struct fields or memory layouts.
6. **Backward Compatibility Guarantee:** Go package namespace preserves public symbol access without breaking imports.
7. **Concrete Implementation Blueprint:**
   ```go
   // pkg/scanner/types.go
   package scanner
   type ScanResult = result.Result[*ScanReport]
   type ScanOptions struct { ... }
   ```
8. **Binary Verification & Acceptance Criteria:** Every package contains a central `types.go` file indexing its domain contracts.

#### SYS-IMP-055: Discrete Map Dispatch Tables Replacing Chained If-Else/Switch Ladders
1. **Identifier & Subsystem:** `SYS-IMP-055` | Pattern Matching | `pkg/dispatch/`
2. **System-Level UX Impact:** Transforms sprawling 200-line switch statements into clean, declarative O(1) handler lookup tables.
3. **Cause:** Chained `if-else if` and huge switch ladders matching tokens or commands are brittle and hard to maintain.
4. **Effect:** Declarative map maps command string to discrete handler function: `var CommandHandlers = map[string]HandlerFunc{...}`.
5. **Zero Harm / Regression Safety:** Discrete handler functions have unit test isolation.
6. **Backward Compatibility Guarantee:** Default fallback handler handles unknown commands gracefully.
7. **Concrete Implementation Blueprint:**
   ```go
   var DispatchTable = map[string]func(ctx Context) *appfault.AppError{
       "find": handleFind,
       "cat":  handleCat,
       "sync": handleSync,
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Adding a new command requires only a single map entry and handler function.

#### SYS-IMP-056: Branch Immutability & Clean Constructor Helper Functions
1. **Identifier & Subsystem:** `SYS-IMP-056` | Object Lifecycle | `pkg/`
2. **System-Level UX Impact:** Guarantees objects are always in a valid, immutable state; eliminates half-initialized object bugs.
3. **Cause:** Constructing empty structs and mutating fields across multiple sequential lines risks race conditions and nil pointer panics.
4. **Effect:** Constructor functions (`NewConfig(...)`) take all required parameters and return fully-initialized, validated instances.
5. **Zero Harm / Regression Safety:** Factory constructors enforce validation rules prior to instantiation.
6. **Backward Compatibility Guarantee:** Exported struct fields remain accessible for read operations.
7. **Concrete Implementation Blueprint:**
   ```go
   func NewScanner(options ScannerOptions) (Scanner, *appfault.AppError) {
       if options.Workers <= 0 { options.Workers = 4 }
       return &scannerImpl{opts: options}, nil
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Constructors validate input and return either a valid instance or structured AppError.

#### SYS-IMP-057: Allocation-Free String Normalization via `strutil.EqualFoldAny`
1. **Identifier & Subsystem:** `SYS-IMP-057` | String Performance | `pkg/strutil/`
2. **System-Level UX Impact:** Speeds up case-insensitive search and token comparisons while reducing garbage collection pressure.
3. **Cause:** Calling `strings.ToLower(str) == "foo"` allocates a new heap string on every comparison inside hot loops.
4. **Effect:** `strutil.EqualFoldAny(token, "foo", "bar")` performs zero-allocation case folding and short-circuits instantly.
5. **Zero Harm / Regression Safety:** Pure drop-in replacement; logic behavior identical.
6. **Backward Compatibility Guarantee:** Standard library `strings` functions remain available.
7. **Concrete Implementation Blueprint:**
   ```go
   func EqualFoldAny(s string, targets ...string) bool {
       for _, target := range targets {
           if strings.EqualFold(s, target) { return true }
       }
       return false
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Benchmark test confirms 0 B/op and 0 allocs/op for `EqualFoldAny`.

#### SYS-IMP-058: Variadic & Spread Parameters Eliminating Artificial Slice Wrappers
1. **Identifier & Subsystem:** `SYS-IMP-058` | API Design | `pkg/`
2. **System-Level UX Impact:** Simplifies call sites by allowing single or multiple arguments without manually constructing slice literals.
3. **Cause:** Forcing callers to write `pkg.Process([]string{"item"})` adds visual clutter and boilerplate.
4. **Effect:** Signatures accept variadic parameters `pkg.Process(items ...string)`, supporting both `pkg.Process("a")` and `pkg.Process(slice...)`.
5. **Zero Harm / Regression Safety:** Backward-compatible at Go compiler level using `...` slice unpack syntax.
6. **Backward Compatibility Guarantee:** Existing slice-passing call sites continue working with `items...`.
7. **Concrete Implementation Blueprint:**
   ```go
   func ScanFiles(patterns ...string) result.Result[[]string] { ... }
   ```
8. **Binary Verification & Acceptance Criteria:** Call sites pass arguments cleanly without mandatory `[]string{...}` wrappers.

#### SYS-IMP-059: TypeScript Discriminated Unions & Strict `Result<T, AppError>` Wrappers
1. **Identifier & Subsystem:** `SYS-IMP-059` | Frontend Types | `src/types/`
2. **System-Level UX Impact:** Guarantees frontend type safety at compile time; eliminates `undefined is not a function` runtime crashes.
3. **Cause:** Using `any` types and throwing raw Javascript errors makes UI error handling unpredictable.
4. **Effect:** TypeScript interfaces use discriminated unions (`type State = { status: 'idle' } | { status: 'success', data: T }`) and Result envelopes.
5. **Zero Harm / Regression Safety:** TypeScript compiler enforces exhaustive type checking across all switch cases.
6. **Backward Compatibility Guarantee:** Components consuming data receive strongly-typed props.
7. **Concrete Implementation Blueprint:**
   ```typescript
   export type Result<T, E = AppError> = 
     | { isOk: true; value: T; error?: never }
     | { isOk: false; value?: never; error: E };
   ```
8. **Binary Verification & Acceptance Criteria:** `npx tsc --noEmit` verifies 0 `any` violations and exhaustive switch checks.

#### SYS-IMP-060: Python Dynamic Enums & Magic Number Elimination
1. **Identifier & Subsystem:** `SYS-IMP-060` | Automation Scripts | `03-ai-scripts/`
2. **System-Level UX Impact:** Makes Python automation scripts robust against typos and self-documenting.
3. **Cause:** Hardcoded string literals (`"failed"`, `"in_progress"`) and magic integer thresholds (`30`, `5000`) cause silent script bugs.
4. **Effect:** Standard Python `StrEnum` and typed constants tables centralized in `03-ai-scripts/02-shared-engine.py`.
5. **Zero Harm / Regression Safety:** Centralized constants referenced by all child scripts; zero logic deviation.
6. **Backward Compatibility Guarantee:** String representation of enums matches previous raw strings exactly.
7. **Concrete Implementation Blueprint:**
   ```python
   class TaskStatus(StrEnum):
       PENDING = "PENDING"
       IN_PROGRESS = "IN_PROGRESS"
       COMPLETED = "COMPLETED"
       FAILED = "FAILED"
   ```
8. **Binary Verification & Acceptance Criteria:** Running `python 03-ai-scripts/21-sequence-integrity-linter.py` passes with zero string discrepancies.

---

### Cluster 07: Performance, Caching & Resource Efficiency (SYS-IMP-061 – SYS-IMP-070)

#### SYS-IMP-061: GitMap Zero-Disk Memory Streaming File Reader (`gitmap cat`)
1. **Identifier & Subsystem:** `SYS-IMP-061` | IO Streaming | `cmd/gitmap/cat.go`
2. **System-Level UX Impact:** Streams file contents directly to stdout without creating temporary disk buffers or memory bloat.
3. **Cause:** Reading large files via shell tools (`Get-Content`, `type`) can exhaust PowerShell memory and hang terminals.
4. **Effect:** `gitmap cat` uses direct memory-mapped or buffered streaming to pipe content in sub-millisecond speeds.
5. **Zero Harm / Regression Safety:** Read-only operation; zero disk mutations.
6. **Backward Compatibility Guarantee:** Outputs raw file bytes identical to POSIX `cat`.
7. **Concrete Implementation Blueprint:**
   ```go
   func StreamFileToStdout(filePath string) *appfault.AppError {
       f, err := os.Open(filePath)
       if err != nil { return appfault.Wrap(err, "FILE_OPEN_FAILED") }
       defer f.Close()
       _, err = io.Copy(os.Stdout, f)
       return appfault.Wrap(err, "STREAM_FAILED")
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap cat readme.md` outputs file contents in <5ms with zero disk writes.

#### SYS-IMP-062: SQLite WAL Mode & Busy-Timeout Concurrency Tuning
1. **Identifier & Subsystem:** `SYS-IMP-062` | Database Performance | `pkg/db/sqlite.go`
2. **System-Level UX Impact:** Allows unlimited concurrent readers alongside active subagent writers with zero database lock contention.
3. **Cause:** Default SQLite rollback journal locks the entire database file during writes, blocking reader queries.
4. **Effect:** All database connections initialize with `PRAGMA journal_mode=WAL;` and `PRAGMA busy_timeout=5000;`.
5. **Zero Harm / Regression Safety:** WAL mode is completely crash-safe and ACID-compliant.
6. **Backward Compatibility Guarantee:** Transparent to all SQL queries and schema models.
7. **Concrete Implementation Blueprint:**
   ```go
   db.Exec("PRAGMA journal_mode = WAL;")
   db.Exec("PRAGMA busy_timeout = 5000;")
   db.Exec("PRAGMA synchronous = NORMAL;")
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap heatmap status` confirms `journal_mode=wal` and 0 lock collision events.

#### SYS-IMP-063: Multi-Core Live Streaming Pattern Scanner (`gitmap aum search`)
1. **Identifier & Subsystem:** `SYS-IMP-063` | Code Discovery | `cmd/gitmap/aum_search.go`
2. **System-Level UX Impact:** Searches through 50,000 files across polyglot repositories in under 150ms.
3. **Cause:** Legacy search tools (`Select-String`, raw grep) single-thread file traversal, causing 10-30 second delays.
4. **Effect:** Parallel worker pool with memory-chunked regex scanning streams matching lines in real time.
5. **Zero Harm / Regression Safety:** Read-only traversal; honors `.gitignore` automatically.
6. **Backward Compatibility Guarantee:** Supports `-e` extension filtering, `-r` regex, and `-i` case-insensitivity flags.
7. **Concrete Implementation Blueprint:**
   ```go
   scanner := aum.NewScanner(aum.ScannerConfig{
       Workers: runtime.NumCPU(),
       BufferPoolSize: 1024 * 64,
   })
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap aum search "func" . -e .go` executes across entire repo in <200ms.

#### SYS-IMP-064: Shared Python Automation Engine Cache Memoization
1. **Identifier & Subsystem:** `SYS-IMP-064` | Script Performance | `03-ai-scripts/02-shared-engine.py`
2. **System-Level UX Impact:** Cuts Python CLI script execution times from 800ms down to 40ms by caching AST parse trees.
3. **Cause:** Every Python script invocation reparses large directory trees and markdown files from scratch.
4. **Effect:** In-memory or temporary disk memoization cache reuses parsed file lists when timestamps match.
5. **Zero Harm / Regression Safety:** Validates file mtime before serving cached data; auto-invalidates on file changes.
6. **Backward Compatibility Guarantee:** Transparent to scripts importing `shared-engine.py`.
7. **Concrete Implementation Blueprint:**
   ```python
   def get_cached_file_inventory(directory: str) -> list[str]:
       cache_key = f"{directory}:{get_dir_mtime(directory)}"
       if cache_key in _MEM_CACHE: return _MEM_CACHE[cache_key]
       _MEM_CACHE[cache_key] = scan_directory(directory)
       return _MEM_CACHE[cache_key]
   ```
8. **Binary Verification & Acceptance Criteria:** Sequential invocations of `05-guideline-autofixer.py` run 10x faster on second run.

#### SYS-IMP-065: Repo-Scoped Build & Cache Purge Engine
1. **Identifier & Subsystem:** `SYS-IMP-065` | Cache Hygiene | `cmd/gitmap/clean_caches.go`
2. **System-Level UX Impact:** Reclaims 2-10 GB of lost disk space and eliminates mysterious stale compiler build artifacts.
3. **Cause:** Lingering Go test caches, Vite incremental build outputs, and Python `__pycache__` bloat disk storage over time.
4. **Effect:** Scoped purge cleans caches strictly within the target repository without touching user global tool caches.
5. **Zero Harm / Regression Safety:** Strictly skips tracked git files and `.git` objects.
6. **Backward Compatibility Guarantee:** Next build rebuilds cleanly from source.
7. **Concrete Implementation Blueprint:**
   ```bash
   # Cleans local repo build caches safely
   gitmap clean-cache --repo .
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap find "*__pycache__*"` returns 0 results after purge run.

#### SYS-IMP-066: Sub-10ms Glob File Indexer (`gitmap find`)
1. **Identifier & Subsystem:** `SYS-IMP-066` | File Finding | `cmd/gitmap/find.go`
2. **System-Level UX Impact:** Instant filename and glob searches across 10,000+ files in <10ms.
3. **Cause:** OS directory tree traversal (`Get-ChildItem -Recurse`) takes 3-8 seconds on Windows NTFS filesystems.
4. **Effect:** High-speed in-memory B-tree index matches wildcards and path patterns instantly.
5. **Zero Harm / Regression Safety:** Read-only traversal; zero side-effects.
6. **Backward Compatibility Guarantee:** Supports standard glob patterns (`*`, `?`, `**/*.md`).
7. **Concrete Implementation Blueprint:**
   ```go
   func FastFind(pattern string, baseDir string) []string {
       return index.MatchGlob(pattern, baseDir)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap find "*.md" 02-spec` outputs matching specs in <10ms.

#### SYS-IMP-067: Lazy Regular Expression Compilation via `lazyregex`
1. **Identifier & Subsystem:** `SYS-IMP-067` | Memory & CPU Efficiency | `pkg/lazyregex/`
2. **System-Level UX Impact:** Accelerates CLI startup times by deferring regex compilation until first invocation.
3. **Cause:** Initializing dozens of `regexp.MustCompile` in global `var` blocks burns CPU cycles even when commands don't use them.
4. **Effect:** Reusable `lazyregex.Pattern` compiles regular expressions on-demand with thread-safe `sync.Once`.
5. **Zero Harm / Regression Safety:** Complete drop-in replacement for standard `regexp.Regexp`.
6. **Backward Compatibility Guarantee:** Exposes identical `MatchString`, `FindStringSubmatch`, and `ReplaceAllString` methods.
7. **Concrete Implementation Blueprint:**
   ```go
   type Pattern struct {
       once sync.Once
       expr string
       rx   *regexp.Regexp
   }
   func (p *Pattern) Rx() *regexp.Regexp {
       p.once.Do(func() { p.rx = regexp.MustCompile(p.expr) })
       return p.rx
   }
   ```
8. **Binary Verification & Acceptance Criteria:** CLI starts up in <8ms; regex patterns compile lazily on first match call.

#### SYS-IMP-068: Go Value Semantics & Pointer Elimination on Small Structs
1. **Identifier & Subsystem:** `SYS-IMP-068` | Memory Architecture | `pkg/`
2. **System-Level UX Impact:** Eliminates escape analysis heap allocations, significantly reducing GC pause times.
3. **Cause:** Gratuitous pointer usage (`*Config`, `*Point`) for small immutable structs (<64 bytes) forces heap allocations.
4. **Effect:** Small structs passed by value; Go compiler allocates them directly on the CPU stack.
5. **Zero Harm / Regression Safety:** Pure performance optimization; mathematical logic unchanged.
6. **Backward Compatibility Guarantee:** Caller code simplified; zero risk of nil pointer dereferences.
7. **Concrete Implementation Blueprint:**
   ```go
   // BEFORE: func Validate(c *Config) bool
   // AFTER:
   func Validate(c Config) bool { ... }
   ```
8. **Binary Verification & Acceptance Criteria:** `go build -gcflags="-m"` confirms small structs do not escape to the heap.

#### SYS-IMP-069: HTTP Response Compression & Dynamic Chunk Streaming in Fleet APIs
1. **Identifier & Subsystem:** `SYS-IMP-069` | Network Efficiency | `pkg/api/`
2. **System-Level UX Impact:** Reduces bandwidth consumption by 70% and latency by 50% when querying remote telemetry.
3. **Cause:** Uncompressed JSON payload transfers across multi-node cluster APIs cause latency spikes on large sync payloads.
4. **Effect:** Gzip / Zstandard dynamic stream compression applied automatically when caller supports `Accept-Encoding`.
5. **Zero Harm / Regression Safety:** Transparent HTTP transport middleware; payloads decode identically.
6. **Backward Compatibility Guarantee:** Uncompressed fallback for clients without compression support.
7. **Concrete Implementation Blueprint:**
   ```go
   func CompressionMiddleware(next http.Handler) http.Handler {
       return gziphandler.GzipHandler(next)
   }
   ```
8. **Binary Verification & Acceptance Criteria:** API responses verify `Content-Encoding: gzip` header and valid decompressed body.

#### SYS-IMP-070: Zero-Storage Memory Footprint Optimizer for Codebase Traversal
1. **Identifier & Subsystem:** `SYS-IMP-070` | Memory Governance | `pkg/traversal/`
2. **System-Level UX Impact:** Traverses repositories with >100,000 files while consuming <30 MB of RAM.
3. **Cause:** Storing full recursive directory paths in a giant in-memory array causes out-of-memory spikes.
4. **Effect:** Channels and iterators stream file paths one-by-one, recycling string buffers.
5. **Zero Harm / Regression Safety:** Constant O(1) memory footprint regardless of repository size.
6. **Backward Compatibility Guarantee:** Consumer functions receive standard Go iterators or channels.
7. **Concrete Implementation Blueprint:**
   ```go
   func WalkStream(root string, fn func(path string) bool) error {
       return filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
           if err != nil || !fn(p) { return err }
           return nil
       })
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Memory profiling confirms memory usage stays below 30 MB during 100k file traversal.

---

### Cluster 08: Diagnostics, Observability & Error Management (SYS-IMP-071 – SYS-IMP-080)

#### SYS-IMP-071: Structured 4-Part Root Cause Analysis (RCA) Engine
1. **Identifier & Subsystem:** `SYS-IMP-071` | Defect Triage | `02-spec/22-app-issues/` & `.ai-memory/cicd-issues/`
2. **System-Level UX Impact:** Eliminates developer guesswork and repetitive regressions by documenting precise root causes and verified fixes.
3. **Cause:** Vague bug fix commits ("fixed issue") provide zero documentation on why a defect happened or how it was verified.
4. **Effect:** Every resolved issue documents: (1) Observed Symptom, (2) Grounded Root Cause, (3) Surgical Remediation, (4) Verification Evidence.
5. **Zero Harm / Regression Safety:** Strictly improves documentation and diagnostic clarity.
6. **Backward Compatibility Guarantee:** Standard markdown issue schema preserved.
7. **Concrete Implementation Blueprint:**
   ```markdown
   ## 1. Observed Symptom: [Exact error message and stack trace]
   ## 2. Grounded Root Cause: [Why the defect occurred at the code level]
   ## 3. Surgical Remediation: [Exact code modification made]
   ## 4. Verification Evidence: [Test command and exit code 0 proof]
   ```
8. **Binary Verification & Acceptance Criteria:** Issue linter verifies all 4 mandatory RCA sections exist in resolved issue files.

#### SYS-IMP-072: GitMap Pipeline Error Extractor (`gitmap pe -t`)
1. **Identifier & Subsystem:** `SYS-IMP-072` | CI/CD Diagnostics | `cmd/gitmap/pipeline_errors.go`
2. **System-Level UX Impact:** Extracts exact failure lines and root cause error messages from 10,000-line CI/CD logs in <1 second.
3. **Cause:** Developers and agents must scroll through thousands of lines of raw build output to find the 2 lines of compiler errors.
4. **Effect:** `gitmap pe -t` parses raw GitHub Actions logs and outputs only the relevant failure block and annotation.
5. **Zero Harm / Regression Safety:** Read-only log parsing; zero side effects.
6. **Backward Compatibility Guarantee:** Full raw logs accessible via `gitmap pipeline logs`.
7. **Concrete Implementation Blueprint:**
   ```go
   func ExtractPipelineErrors(rawLog string) string {
       return errorMatcher.ExtractFirstFailureBlock(rawLog, 15) // 15 lines bounded context
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap pe -t` isolates the failing assertion line from test runner output.

#### SYS-IMP-073: Contextual Call-Stack Truncation & Signal-to-Noise Error Filtering
1. **Identifier & Subsystem:** `SYS-IMP-073` | Error Reporting | `pkg/appfault/stack.go`
2. **System-Level UX Impact:** Keeps error messages compact and readable by stripping irrelevant runtime and framework stack frames.
3. **Cause:** 80-frame runtime call stacks overwhelm model context windows and obscure the actual user-code error location.
4. **Effect:** Stack trace filter retains strictly repository-owned frames (`pkg/`, `cmd/`), collapsing external standard library frames.
5. **Zero Harm / Regression Safety:** Preserves exact line numbers and filenames of actual bug locations.
6. **Backward Compatibility Guarantee:** `--full-stack` flag outputs unmodified raw stack traces when needed.
7. **Concrete Implementation Blueprint:**
   ```go
   func CleanStackTrace(frames []runtime.Frame) []runtime.Frame {
       return filter(frames, func(f runtime.Frame) bool { return strings.Contains(f.File, "coding-guidelines") })
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Error outputs display <= 5 relevant code frames instead of 80 system frames.

#### SYS-IMP-074: Non-Polling Dynamic ETA CI/CD Waiting Protocol (`gitmap pipeline-ai`)
1. **Identifier & Subsystem:** `SYS-IMP-074` | Pipeline Monitor | `cmd/gitmap/pipeline_ai.go`
2. **System-Level UX Impact:** Eliminates API rate-limit exhaustion and credit burn by sleeping adaptively instead of tight polling.
3. **Cause:** Agents polling `gh run view` every 2 seconds waste GitHub API quotas and burn tokens in tight spin loops.
4. **Effect:** Pipeline-AI predicts expected completion time (`etaSeconds`) and commands the agent to sleep adaptively.
5. **Zero Harm / Regression Safety:** Replaces polling loops with single-shot non-blocking waits.
6. **Backward Compatibility Guarantee:** Compatible with standard GitHub Actions workflow runs.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap pipeline-ai status -t "$ETA_SECONDS"
   ```
8. **Binary Verification & Acceptance Criteria:** CI/CD monitoring command sleeps for calculated ETA and makes <= 3 total API requests.

#### SYS-IMP-075: Zero-Storage GitHub Actions Diagnostic Summaries
1. **Identifier & Subsystem:** `SYS-IMP-075` | CI/CD Infrastructure | `.github/workflows/`
2. **System-Level UX Impact:** Renders rich test reports and lint summaries directly in GitHub PRs with 0.0 GB storage quota usage.
3. **Cause:** Uploading test artifacts (`actions/upload-artifact`) exhausts GitHub free-tier 0.5 GB monthly quotas, blocking pipelines.
4. **Effect:** All test summaries and lint diagnostics write natively to `$GITHUB_STEP_SUMMARY` and inline annotations (`::error::`).
5. **Zero Harm / Regression Safety:** 100% compliant with the Zero-Storage GitHub Actions Mandate.
6. **Backward Compatibility Guarantee:** Release asset binaries uploaded via `gh release upload` remain fully permitted.
7. **Concrete Implementation Blueprint:**
   ```yaml
   - name: Write Step Summary
     run: |
       echo "### Test Results" >> $GITHUB_STEP_SUMMARY
       echo "- All 42 tests passed in 1.4s" >> $GITHUB_STEP_SUMMARY
   ```
8. **Binary Verification & Acceptance Criteria:** Workflow execution leaves GitHub Actions artifact storage at exactly 0.0 MB.

#### SYS-IMP-076: Segregated Issue Archival Routing
1. **Identifier & Subsystem:** `SYS-IMP-076` | Defect Archival | `02-spec/22-app-issues/` & `.ai-memory/cicd-issues/`
2. **System-Level UX Impact:** Keeps CI/CD pipeline infrastructure defects completely segregated from application logic defects.
3. **Cause:** Mixing CI/CD failures with application runtime bugs creates confusion when prioritizing product backlog items.
4. **Effect:** CI/CD pipeline issues archived in `.ai-memory/cicd-issues/`; application bugs archived in `02-spec/22-app-issues/`.
5. **Zero Harm / Regression Safety:** Enforces clean separation of concerns across documentation trees.
6. **Backward Compatibility Guarantee:** Master indices (`cicd-index.md` and `22-app-issues/readme.md`) link all items.
7. **Concrete Implementation Blueprint:**
   ```text
   CI/CD Issues  -> .ai-memory/cicd-issues/XX-<slug>.md (Index: .ai-memory/cicd-index.md)
   App Issues    -> 02-spec/22-app-issues/XX-<slug>.md (Index: 02-spec/22-app-issues/readme.md)
   ```
8. **Binary Verification & Acceptance Criteria:** Audit confirms zero CI/CD issues in `22-app-issues` and zero app bugs in `cicd-issues`.

#### SYS-IMP-077: GitMap Forensic Crash Diagnostic Suite (`gitmap agent crashed`)
1. **Identifier & Subsystem:** `SYS-IMP-077` | Swarm Forensics | `cmd/gitmap/agent_crashed.go`
2. **System-Level UX Impact:** Automatically diagnoses why an agent crashed, what file it was touching, and the exact error code.
3. **Cause:** Agent crashes leave silent failures, forcing developers to manually inspect log transcripts.
4. **Effect:** `gitmap agent crashed` analyzes SQLite action logs, prints forensic post-mortem summary, and suggests recovery action.
5. **Zero Harm / Regression Safety:** Read-only forensic analysis.
6. **Backward Compatibility Guarantee:** Standalone diagnostic tool.
7. **Concrete Implementation Blueprint:**
   ```go
   func DiagnoseCrashedAgents() *ForensicReport {
       return db.Query("SELECT agent_id, file_path, action, error_message FROM agent_logs WHERE status='CRASHED'")
   }
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap agent crashed` prints actionable forensic summary for any terminated agent.

#### SYS-IMP-078: Backward-Compatible Error Alias Shims (`type Fault = AppError`)
1. **Identifier & Subsystem:** `SYS-IMP-078` | Refactoring Safety | `pkg/appfault/`
2. **System-Level UX Impact:** Allows gradual codebase migration to unified error packages without breaking existing packages.
3. **Cause:** Renaming core error types breaks compilation across dozens of dependent packages simultaneously.
4. **Effect:** Type alias shims (`type Fault = AppError`, package `apperror` aliases) ensure 100% non-breaking compilation.
5. **Zero Harm / Regression Safety:** Zero runtime performance cost; identical type identity in Go type system.
6. **Backward Compatibility Guarantee:** Both `*appfault.AppError` and legacy `*apperror.Fault` work interchangeably.
7. **Concrete Implementation Blueprint:**
   ```go
   package appfault
   type AppError struct { ... }
   type Fault = AppError
   ```
8. **Binary Verification & Acceptance Criteria:** Code referencing `appfault.Fault` compiles cleanly with zero warnings.

#### SYS-IMP-079: Structured Diagnostic Logging Envelopes with Correlated Trace IDs
1. **Identifier & Subsystem:** `SYS-IMP-079` | Distributed Telemetry | `pkg/logger/`
2. **System-Level UX Impact:** Allows developers to trace an action across lead orchestrator, subagents, and database logs with a single ID.
3. **Cause:** Independent log lines from 3 concurrent subagents interleave, making it impossible to reconstruct causality.
4. **Effect:** All log events include unified `trace_id` propagated across tool calls and database writes.
5. **Zero Harm / Regression Safety:** Lightweight string header field; zero impact on core logic.
6. **Backward Compatibility Guarantee:** Standard log formatters output trace ID as structured prefix.
7. **Concrete Implementation Blueprint:**
   ```json
   {"trace_id": "wave-01-agent-02", "level": "INFO", "action": "EDIT", "file": "pkg/auth/login.go"}
   ```
8. **Binary Verification & Acceptance Criteria:** Filtering logs by `trace_id` reveals the complete chronological lifecycle of that subtask.

#### SYS-IMP-080: Interactive CLI Diagnostic Reporter with One-Click Remediation Hints
1. **Identifier & Subsystem:** `SYS-IMP-080` | Developer Guidance | `pkg/diagnostics/`
2. **System-Level UX Impact:** Suggests the exact terminal command needed to resolve common errors (e.g. untracking ignored files).
3. **Cause:** Errors that state "working tree dirty" without specifying which command fixes it slow down developers.
4. **Effect:** Error output concludes with actionable hint: `Tip: Run 'gitmap clean-cache' or 'git rm --cached <file>' to resolve.`
5. **Zero Harm / Regression Safety:** Informational guidance only.
6. **Backward Compatibility Guarantee:** Exits with standard non-zero error code.
7. **Concrete Implementation Blueprint:**
   ```go
   fmt.Fprintf(os.Stderr, "Error: Tracked file in .gitignore detected: %s\n", path)
   fmt.Fprintf(os.Stderr, "Remediation: Run 'git rm --cached %s' to untrack.\n", path)
   ```
8. **Binary Verification & Acceptance Criteria:** Diagnostic output contains explicit copy-pasteable remediation command.

---

### Cluster 09: Testing, Verification & Regression Isolation (SYS-IMP-081 – SYS-IMP-090)

#### SYS-IMP-081: Destructive OS Call & System Shutdown Test Isolation
1. **Identifier & Subsystem:** `SYS-IMP-081` | Test Isolation | `02-spec/02-coding-guidelines/24-isolate-os-tests.md`
2. **System-Level UX Impact:** Eliminates the catastrophic risk of unit tests accidentally executing real OS shutdown or system reboot commands.
3. **Cause:** Unmocked tests testing shutdown or power-off functions can trigger actual machine restart on developer machines.
4. **Effect:** Injectable mock command executors mandatory; unit tests verify mocked command duration and never invoke real OS binaries.
5. **Zero Harm / Regression Safety:** Absolute safety barrier protecting developer and CI/CD machines.
6. **Backward Compatibility Guarantee:** Production execution code injects the real OS executor cleanly.
7. **Concrete Implementation Blueprint:**
   ```go
   type CommandExecutor interface {
       Execute(cmd string, args ...string) (string, error)
   }
   type MockExecutor struct { RecordedCommands []string }
   ```
8. **Binary Verification & Acceptance Criteria:** Running test suite verifies zero real OS process spawned; 100% mocked calls.

#### SYS-IMP-082: Temporary End-to-End (E2E) Test Skipping Isolation
1. **Identifier & Subsystem:** `SYS-IMP-082` | Test Suite Protection | `01-prompts/21-temp-end-to-end-tests/`
2. **System-Level UX Impact:** Allows developers to write scratch end-to-end integration tests without breaking routine CI/CD pipelines.
3. **Cause:** Incomplete or slow local integration tests added to standard test folders fail remote CI/CD builds.
4. **Effect:** Temporary E2E tests enforce strict skip-by-default flags (`t.Skip()` / `@pytest.mark.skip`) unless explicit env var is set.
5. **Zero Harm / Regression Safety:** Zero routine CI/CD pipeline breakage.
6. **Backward Compatibility Guarantee:** Explicitly executable via `RUN_TEMP_E2E=1 go test ./tests/e2e/...`.
7. **Concrete Implementation Blueprint:**
   ```go
   if os.Getenv("RUN_TEMP_E2E") != "1" {
       t.Skip("Skipping temporary E2E test by default. Set RUN_TEMP_E2E=1 to execute.")
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running `go test ./...` skips temporary E2E tests automatically with exit code 0.

#### SYS-IMP-083: Surgical Single-File Micro-Linter Suite
1. **Identifier & Subsystem:** `SYS-IMP-083` | Fast Verification | `linter-scripts/`
2. **System-Level UX Impact:** Validates guideline compliance on a single modified file in <20ms without scanning the entire repo.
3. **Cause:** Running full repository linters to check a 5-line edit wastes 10-15 seconds per iteration.
4. **Effect:** Fast Python micro-linters accept target file path and check relative paths, booleans, and forbidden strings instantly.
5. **Zero Harm / Regression Safety:** Read-only AST and regex analysis; zero side effects.
6. **Backward Compatibility Guarantee:** Can be run on directories or individual files.
7. **Concrete Implementation Blueprint:**
   ```bash
   python linter-scripts/check-relative-paths.py path/to/file.md
   ```
8. **Binary Verification & Acceptance Criteria:** Micro-linter runs on single file in <20ms and exits 0 on valid code.

#### SYS-IMP-084: Non-Polluting In-Memory Unit Test Fixtures & Reusable Harness Vault
1. **Identifier & Subsystem:** `SYS-IMP-084` | Test Cleanliness | `pkg/testutil/`
2. **System-Level UX Impact:** Tests never leave orphaned temporary test files, databases, or directories on disk.
3. **Cause:** Tests creating real files in `tmp/` or workspace root pollute the git tree and cause dirty working tree errors.
4. **Effect:** Tests use in-memory filesystems (`afero.NewMemMapFs()`) or in-memory SQLite (`:memory:`) with automatic cleanup.
5. **Zero Harm / Regression Safety:** Working tree remains 100% clean after test runs.
6. **Backward Compatibility Guarantee:** Standard filesystem interfaces preserved.
7. **Concrete Implementation Blueprint:**
   ```go
   func SetupTestDB(t *testing.T) *sql.DB {
       db, _ := sql.Open("sqlite3", ":memory:")
       t.Cleanup(func() { db.Close() })
       return db
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Running test suite leaves `git status --porcelain` completely empty.

#### SYS-IMP-085: Behavioral Unit Test Naming & AAA Structural Enforcement
1. **Identifier & Subsystem:** `SYS-IMP-085` | Test Quality | `tests/`
2. **System-Level UX Impact:** Failing test names immediately explain what behavior broke without reading test source code.
3. **Cause:** Cryptic test names like `TestDoWork` or `TestHelper2` provide zero diagnostic value upon failure.
4. **Effect:** Tests named strictly using `Test<Unit>_<Condition>_<ExpectedResult>` and structured with Arrange-Act-Assert blocks.
5. **Zero Harm / Regression Safety:** Pure test naming and organization; no behavioral change.
6. **Backward Compatibility Guarantee:** Go test runner discovers all tests normally.
7. **Concrete Implementation Blueprint:**
   ```go
   func TestAUMSearch_WhenPatternMatches_ReturnsLineNumbers(t *testing.T) {
       // Arrange: set up mock file
       // Act: run search
       // Assert: verify line numbers
   }
   ```
8. **Binary Verification & Acceptance Criteria:** Test names adhere to standard behavioral pattern; failures report exact condition.

#### SYS-IMP-086: Incremental Smart Test Runner with Centralized Manifests
1. **Identifier & Subsystem:** `SYS-IMP-086` | Test Execution | `.agents/skills/smart-test-runner-and-inventory/`
2. **System-Level UX Impact:** Runs only the tests impacted by recent file changes, reducing test execution times from minutes to seconds.
3. **Cause:** Running entire test suites on every minor edit wastes developer time and model context.
4. **Effect:** Analyzes git diff, maps modified packages to test inventory manifest, and executes strictly affected test suites.
5. **Zero Harm / Regression Safety:** Comprehensive regression suites run during final release ceremony only.
6. **Backward Compatibility Guarantee:** Full test suite execution available via `--all`.
7. **Concrete Implementation Blueprint:**
   ```python
   impacted_tests = map_diff_to_test_manifest(git_modified_files)
   run_targeted_tests(impacted_tests)
   ```
8. **Binary Verification & Acceptance Criteria:** Editing `pkg/appfault/` runs only `pkg/appfault/` unit tests in <500ms.

#### SYS-IMP-087: Strict No-Build / No-Test Rule During Active Iteration
1. **Identifier & Subsystem:** `SYS-IMP-087` | Agent Velocity | `01-prompts/14-execute/`
2. **System-Level UX Impact:** Prevents agents from wasting 80% of their step budget in repetitive, slow compiler build loops.
3. **Cause:** Agents building after every 2-line edit trigger long compilation times, run out of steps, and fail to finish tasks.
4. **Effect:** Strict negative guardrail: every time code is written, do not build the code; rely on surgical linters until final gate.
5. **Zero Harm / Regression Safety:** Final Phase 3 gate performs single atomic verification before push.
6. **Backward Compatibility Guarantee:** Full build commands remain available when explicitly requested by user.
7. **Concrete Implementation Blueprint:**
   ```markdown
   - [ ] NO BUILD CHECKING: Every time you write code, you do not build the code.
   - [ ] NO LOCAL TESTS: Never run all tests locally unless explicitly asked by the user.
   ```
8. **Binary Verification & Acceptance Criteria:** Agent execution transcript confirms zero intermediate build commands executed.

#### SYS-IMP-088: Pre-Push Cognitive Heatmap Verification Gate
1. **Identifier & Subsystem:** `SYS-IMP-088` | Verification Audit | `cmd/gitmap/heatmap_tree.go`
2. **System-Level UX Impact:** Proves mathematically that all planned subtasks, touched files, and expected outcomes were completed.
3. **Cause:** Tasks marked "done" when some files were only partially edited or left in an inconsistent state.
4. **Effect:** Phase 3 runs `gitmap heatmap tree --task <task>` and verifies 100% action closure before allowing commit.
5. **Zero Harm / Regression Safety:** Aborts GitMap commit command if any planned subtask remains incomplete.
6. **Backward Compatibility Guarantee:** Seamlessly integrated into standard Phase 3 push gate.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap heatmap tree --task "$TASK_SLUG"
   gitmap heatmap status --task "$TASK_SLUG" --assert-all-met
   ```
8. **Binary Verification & Acceptance Criteria:** Push gate halts with clear error if any subagent action outcome is unmet.

#### SYS-IMP-089: Flaky Test Quarantine & Deterministic Mock Injection
1. **Identifier & Subsystem:** `SYS-IMP-089` | Pipeline Reliability | `tests/quarantine/`
2. **System-Level UX Impact:** Eliminates pipeline build failures caused by non-deterministic network latency or timer races.
3. **Cause:** Tests depending on real system clock time or external HTTP APIs fail intermittently, frustrating developers.
4. **Effect:** Time and network dependencies injected via deterministic mock interfaces; flaky tests quarantined for triage.
5. **Zero Harm / Regression Safety:** Non-deterministic tests isolated until refactored into pure deterministic unit tests.
6. **Backward Compatibility Guarantee:** Production clock and network adapters used in real binaries.
7. **Concrete Implementation Blueprint:**
   ```go
   type Clock interface { Now() time.Time }
   type MockClock struct { CurrentTime time.Time }
   ```
8. **Binary Verification & Acceptance Criteria:** Test suite runs 100 consecutive times with 100% green deterministic passes.

#### SYS-IMP-090: Cross-Platform Installer Smoke-Testing Suite
1. **Identifier & Subsystem:** `SYS-IMP-090` | Release Readiness | `tests/installer/`
2. **System-Level UX Impact:** Guarantees installation scripts (`local-install.ps1`, `run.ps1`) work out-of-the-box on clean developer machines.
3. **Cause:** Installer scripts broken by silent syntax errors or missing parameter checks fail when new developers onboard.
4. **Effect:** Automated smoke test suite validates argument forwarding, exit code propagation, and help flags across platforms.
5. **Zero Harm / Regression Safety:** Runs in non-destructive dry-run mode.
6. **Backward Compatibility Guarantee:** Supports PowerShell on Windows and Bash on Linux/macOS.
7. **Concrete Implementation Blueprint:**
   ```bash
   tests/installer/check-bundle-installers.sh
   tests/installer/check-fix-repo-arg-forwarding.sh
   ```
8. **Binary Verification & Acceptance Criteria:** All installer acceptance test scripts exit with code 0.

---

### Cluster 10: Documentation, Memory & Spec Governance (SYS-IMP-091 – SYS-IMP-100)

#### SYS-IMP-091: Strict Relative Git Paths Mandate Validator
1. **Identifier & Subsystem:** `SYS-IMP-091` | Path Hygiene | `linter-scripts/check-relative-paths.py`
2. **System-Level UX Impact:** Eliminates broken markdown links and machine-specific absolute path leaks across documentation and plans.
3. **Cause:** Hardcoding absolute paths (`C:\Users\...` or `/home/...`) breaks links when other developers or CI/CD runners view files.
4. **Effect:** Total ban on absolute paths and `file:///` URIs; all references strictly relative starting from repository root.
5. **Zero Harm / Regression Safety:** Standard relative path navigation works on all operating systems and web viewers.
6. **Backward Compatibility Guarantee:** Relative links render natively in GitHub and IDE markdown previews.
7. **Concrete Implementation Blueprint:**
   ```python
   # Ban regex: (?:file:\/\/\/|[a-zA-Z]:[\\\/]|\/(?:Users|home|root)[\\\/])
   if re.search(r'(?:file:\/\/\/|[a-zA-Z]:[\\\/])', line):
       raise ValueError(f"Absolute path or file URI violation at {file}:{line_no}")
   ```
8. **Binary Verification & Acceptance Criteria:** `python linter-scripts/check-relative-paths.py` reports 0 absolute path violations.

#### SYS-IMP-092: Strict Lowercase Filename Governance & Automated Case Normalizer
1. **Identifier & Subsystem:** `SYS-IMP-092` | Filesystem Hygiene | `03-ai-scripts/10-encoding-normalizer.py`
2. **System-Level UX Impact:** Prevents git case-sensitivity merge conflicts between Windows (case-insensitive) and Linux (case-sensitive).
3. **Cause:** Files named `README.md` or `CamelCase.go` cause git index drift and tracking anomalies across operating systems.
4. **Effect:** 100% strict lowercase naming enforced across all files, documentation, scripts, and specs without exception.
5. **Zero Harm / Regression Safety:** Mechanical renaming with git index tracking (`git mv -f`).
6. **Backward Compatibility Guarantee:** Downstream links and imports updated simultaneously.
7. **Concrete Implementation Blueprint:**
   ```bash
   python 03-ai-scripts/03-file-manipulator.py normalize-case .
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap find "*[A-Z]*" 01-prompts/` returns 0 uppercase filenames.

#### SYS-IMP-093: Specification Authoring Lifecycle & Application Spec Architecture
1. **Identifier & Subsystem:** `SYS-IMP-093` | Spec Governance | `02-spec/01-spec-authoring-guide/`
2. **System-Level UX Impact:** Ensures every feature has an authoritative, modular specification before code is written.
3. **Cause:** Jumping straight into code without writing specs leads to incomplete implementations and architecture rework.
4. **Effect:** Modular specs authored in `02-spec/21-app/XX-<slug>/` partitioned into Architecture, Component, and Migration specs.
5. **Zero Harm / Regression Safety:** Read-only planning phase (NO-BUILD, NO-TEST, NO-CODE-CHANGES).
6. **Backward Compatibility Guarantee:** Specifications indexed in `02-spec/21-app/readme.md`.
7. **Concrete Implementation Blueprint:**
   ```text
   02-spec/21-app/XX-<slug>/
   ├── 01-architecture-spec.md
   ├── 02-component-spec.md
   └── 03-migration-and-verification-spec.md
   ```
8. **Binary Verification & Acceptance Criteria:** Application specification audit passes with 100% readiness scorecard.

#### SYS-IMP-094: Continuous Plan Enqueueing & Subtask Dehydration Pipeline
1. **Identifier & Subsystem:** `SYS-IMP-094` | Backlog Management | `.ai-memory/plans/`
2. **System-Level UX Impact:** Organizes execution roadmaps into bite-sized, 5-8 file micro-tasks that any subagent can execute reliably.
3. **Cause:** Giant monolithic plans with 50 tasks overwhelm subagents and lead to incomplete executions.
4. **Effect:** Parent plan written to `.ai-memory/plans/pending/XX-<slug>.md` with lean subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`.
5. **Zero Harm / Regression Safety:** Discrete subtasks have disjoint file ownership boundaries.
6. **Backward Compatibility Guarantee:** Upon completion, consolidated into `.ai-memory/plans/completed/`.
7. **Concrete Implementation Blueprint:**
   ```text
   .ai-memory/plans/pending/XX-<slug>.md
   .ai-memory/plans/subtasks/XX-<slug>/
   ├── 01-task.md
   ├── 02-task.md
   └── 03-task.md
   ```
8. **Binary Verification & Acceptance Criteria:** Subtask linter confirms every subtask owns <= 8 disjoint files.

#### SYS-IMP-095: Consolidated Memory Lifecycle Engine (`ai-memory-consolidate-reduce`)
1. **Identifier & Subsystem:** `SYS-IMP-095` | Memory Hygiene | `.agents/skills/ai-memory-consolidate-reduce/`
2. **System-Level UX Impact:** Keeps the repository light and fast by condensing hundreds of historical subtasks into single summary records.
3. **Cause:** Accumulating thousands of finished subtask files bloats git repository size and slows down file searches.
4. **Effect:** Consolidates completed subtasks into parent plan summaries, deleting redundant subtask files while preserving evidence.
5. **Zero Harm / Regression Safety:** All verified outcomes, diff evidence, and commit SHAs preserved in parent summary.
6. **Backward Compatibility Guarantee:** Historical audit records remain fully traceable.
7. **Concrete Implementation Blueprint:**
   ```bash
   python 03-ai-scripts/20-plan-consolidator.py --completed-only
   ```
8. **Binary Verification & Acceptance Criteria:** Total file count in `.ai-memory/plans/` reduced by >70% after consolidation.

#### SYS-IMP-096: Comprehensive Project Identity & What-To-Read Orientation Engine
1. **Identifier & Subsystem:** `SYS-IMP-096` | Onboarding Architecture | `.ai-memory/what-to-read.md`
2. **System-Level UX Impact:** Allows newly spawned agents or human developers to orient themselves with the entire codebase in <60 seconds.
3. **Cause:** Disjointed documentation forces new developers to spend hours discovering coding rules, architecture, and tools.
4. **Effect:** Unified reading order in `.ai-memory/what-to-read.md` indexes identity, CODE RED rules, active specs, and pending tasks.
5. **Zero Harm / Regression Safety:** Central navigation directory; zero code mutations.
6. **Backward Compatibility Guarantee:** Automatically synchronized with root `readme.md`.
7. **Concrete Implementation Blueprint:**
   ```markdown
   # What to Read Before Touching the Codebase
   1. [Canonical Coding Guidelines](02-spec/02-coding-guidelines/readme.md)
   2. [Strictly Avoid Rules](.ai-memory/strictly-avoid.md)
   3. [Active Specifications](02-spec/21-app/readme.md)
   ```
8. **Binary Verification & Acceptance Criteria:** Reading `what-to-read.md` provides all necessary context to begin development.

#### SYS-IMP-097: Automated Index Sync Deadman Switch
1. **Identifier & Subsystem:** `SYS-IMP-097` | Index Governance | `linter-scripts/check-prompts-loaded.py`
2. **System-Level UX Impact:** Guarantees that every new prompt, spec, or skill is immediately discoverable from directory readmes.
3. **Cause:** Creating a new prompt or spec without linking it in category readmes leaves it orphaned and undiscoverable.
4. **Effect:** Pre-commit deadman switch verifies that 100% of files in `01-prompts/` and `02-spec/` are explicitly linked in their parent `readme.md`.
5. **Zero Harm / Regression Safety:** Fails fast before commit if an unlinked file exists.
6. **Backward Compatibility Guarantee:** Existing links remain valid.
7. **Concrete Implementation Blueprint:**
   ```python
   python linter-scripts/check-prompts-loaded.py
   ```
8. **Binary Verification & Acceptance Criteria:** Script returns exit code 0 when all prompts and specs are registered in readmes.

#### SYS-IMP-098: Semantic Release Ceremony & Automated SemVer Manifest Synchronizer
1. **Identifier & Subsystem:** `SYS-IMP-098` | Release Management | `01-prompts/16-release/`
2. **System-Level UX Impact:** Cuts release ceremony time from 15 minutes to 30 seconds with automated changelog generation and tagging.
3. **Cause:** Manual version bumping across multiple files (`version.json`, `package.json`, `Cargo.toml`) causes version mismatch bugs.
4. **Effect:** Single command (`gitmap release --bump <patch|minor> -y`) updates version manifests, pins root readme, and updates changelog.
5. **Zero Harm / Regression Safety:** Zero builds triggered; purely metadata and tagging operations.
6. **Backward Compatibility Guarantee:** Strictly follows SemVer 2.0.0 conventions.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap release --bump minor -y
   ```
8. **Binary Verification & Acceptance Criteria:** `gitmap list-versions` displays new SemVer version matching `version.json`.

#### SYS-IMP-099: Clean Working Tree & Mandatory Completion Invariant ('No Push = Not Done')
1. **Identifier & Subsystem:** `SYS-IMP-099` | Completion Invariant | `01-prompts/14-execute/`
2. **System-Level UX Impact:** Eliminates the risk of developers or agents claiming tasks are complete while leaving work unpushed.
3. **Cause:** Work left in local uncommitted changes or unpushed commits is invisible to other agents and lost on workspace teardown.
4. **Effect:** If the code is not committed via GitMap and pushed upstream to GitHub, the task is strictly considered INCOMPLETE and NOT DONE.
5. **Zero Harm / Regression Safety:** Guarantees remote repository reflects the exact verified state.
6. **Backward Compatibility Guarantee:** Pushes directly to active remote tracking branch.
7. **Concrete Implementation Blueprint:**
   ```bash
   gitmap cpf "<module> - <summary>"
   # Verify working tree is pristine:
   git status --porcelain # Must return 0 lines
   ```
8. **Binary Verification & Acceptance Criteria:** `git status --porcelain` is completely empty and remote tracking branch contains the commit.

#### SYS-IMP-100: Continuous Improvement Flywheel Generator
1. **Identifier & Subsystem:** `SYS-IMP-100` | Continuous Flywheel | `01-prompts/30-improvements/`
2. **System-Level UX Impact:** Transforms the entire repository into a self-evolving engineering engine that continuously surfaces next-generation UX enhancements.
3. **Cause:** System improvement efforts stall after a single audit without an automated flywheel mechanism to generate the next iteration.
4. **Effect:** Each execution of this prompt enqueues all 100 system improvements as pending plans in `.ai-memory/plans/pending/XX-<slug>.md` and subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`, and dynamically synthesizes the NEXT 100 improvements for subsequent runs.
5. **Zero Harm / Regression Safety:** All proposed improvements require explicit pending plan decomposition and acceptance criteria before execution.
6. **Backward Compatibility Guarantee:** Existing pending tasks remain prioritized; new improvements append cleanly to the backlog.
7. **Concrete Implementation Blueprint:**
   ```markdown
   # Continuous Improvement Flywheel
   - Enqueue current 100 improvements -> .ai-memory/plans/pending/XX-system-improvements-v1.md
   - Decompose subtasks -> .ai-memory/plans/subtasks/XX-system-improvements-v1/
   - Trigger next cycle -> Generate next 100 system improvement candidates (SYS-IMP-101 to SYS-IMP-200)
   ```
8. **Binary Verification & Acceptance Criteria:** `.ai-memory/plans/pending/` contains the generated improvement plan, subtask folder contains all decomposed subtasks, and pending task registers are fully updated.

---

## Actionable Items & Operational Rules (R1 – R16)

When executing this prompt, all AI agents must follow these operational rules:

1. **Rule R1 (Task Breakdown in Visible Chat):** Turn 1 MUST display the 5 Mandatory Pillars and 6-part verification schema in visible chat before any background tool dispatch.
2. **Rule R2 (Targeted Verification Checks):** Run surgical micro-linters on modified files only; never run full test suites unless explicitly commanded.
3. **Rule R3 (No Full Builds):** Every time code is written, do not build the code. Full builds during routine iteration are strictly banned.
4. **Rule R4 (Tool Existence Preflight):** Confirm helper scripts exist via harmless workspace calls before invoking.
5. **Rule R5 (Strict Relative Git Paths):** All paths and markdown links MUST be strictly relative starting from repository root (`02-spec/...`, `cmd/...`). Absolute paths and `file:///` URIs are totally banned.
6. **Rule R6 (Strict Lowercase Filenames):** Every newly created or modified file MUST use strictly lowercase naming.
7. **Rule R7 (GitMap Commit Primacy):** Commit exclusively via GitMap semantic commands (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`). Never include colons inside the message argument.
8. **Rule R8 (.gitignore Hygiene Gate):** Ensure `.gitignore` covers caches, build outputs, and logs. Untrack ignored files via `git rm --cached`.
9. **Rule R9 (Zero Secrets Mandate):** Sensitive credentials, `.env` files, and tokens MUST be offloaded to `repo-secrets` via `gitmap rs`.
10. **Rule R10 (Tree Swarm & Subagent Spawning Gate):** Mandatory spawning of `A = 3` concurrent subagents via `invoke_subagent` (`H = 2`, Capacity = 6). Solo execution is an auto-reject failure.
11. **Rule R11 (Dynamic Join & 100% Duty Cycle):** Idle workers immediately respawn for pending subtasks or dynamically join active subtasks.
12. **Rule R12 (Disjoint File Boundaries):** Subagents must never be assigned overlapping file boundaries to prevent Windows file lock crashes.
13. **Rule R13 (Cognitive Heatmap Telemetry):** Log pre-action intent before file reads/edits and post-action outcome verification.
14. **Rule R14 (Non-Owned Repository Exclusion):** Never modify or push to non-owned repositories (`omis`, `zsh`, `oh-my-posh`).
15. **Rule R15 (Continuous Improvement Flywheel):** Persist all 100 improvement ideas into `.ai-memory/plans/pending/XX-<slug>.md` and subtasks in `.ai-memory/plans/subtasks/XX-<slug>/`.
16. **Rule R16 (Mandatory Completion Invariant):** If code is not committed and pushed upstream to GitHub, the task is strictly INCOMPLETE and NOT DONE.

---

## GitMap High-Speed Search Primacy & Command Suite

GitMap is your **PRIMARY** acceleration engine. Banned ALL variants of `Get-ChildItem`, `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Execute all searches exclusively through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard File Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core glob filename finder |
| **Directory Inventory** | `gitmap lf <dir>` / `gitmap find "*.*" <dir>` | `gitmap lf [pat]` | Instant indexed repository inventory (replaces `Get-ChildItem`, `dir`, `ls`) |
| **Substring File Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Direct zero-disk memory streaming to stdout (replaces `cat`, `type`, `Get-Content`) |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<module> - <msg>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats with prefix, and pushes atomically (no colons in arguments). |
| **Agent Task Engine** | `gitmap task <init\|add\|claim\|complete\|fail\|status>` | `gitmap task` | Sub-millisecond compiled Go SQLite task manager for multi-agent workflows |
| **Heatmap Action Log** | `gitmap heatmap log --task <t> --action <A> --file <f>` | `gitmap heatmap log`| Concurrency-safe cognitive action telemetry into `.ai-memory/heatmap.db` |
| **Fleet Synchronization**| `gitmap sync [--workers 8]` | `gitmap sync` | Parallel Go sync engine across 43 repositories in <5s |

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Instructions given above this prompt are highest priority and strictly non-negotiable.
- [ ] MANDATORY SUBAGENT SPAWNING GATE: Spawning `A = 3` subagents (`H = 2`, Capacity = 6) via `invoke_subagent` was explicitly called.
- [ ] NO SUBAGENT FILE WRITE COLLISIONS: Two subagents never assigned to create or write the same file path simultaneously.
- [ ] NO READ-ONLY SUBAGENT FILE WRITES: Subagents with `TypeName: "research"` have read-only tools and cannot write files.
- [ ] NO SUBAGENT GIT EXECUTION: Subagents never run git commands (`git add`, `git commit`, `git push`, `git status`).
- [ ] NO LOCAL TESTS (TOTAL BAN): Never run all tests locally or full test suites unless explicitly commanded by user.
- [ ] NO BUILD CHECKING (TOTAL BAN): Every time you write code, do not build the code. Never run full builds during iteration.
- [ ] NO RAW GIT COMMITS: Never run raw `git commit`, `git add -A`, or raw `git push`. All commits use GitMap semantic commands.
- [ ] NO INTERMEDIATE COMMITS: All changes committed together in one single atomic GitMap commit at the final step of Phase 3.
- [ ] NO PER-FILE COMMITS: Never commit files individually as you edit them.
- [ ] NO UNCOMMITTED WORK (NO PUSH = NOT DONE): Task is incomplete until committed and confirmed pushed to GitHub.
- [ ] NO RG, RIPGREP, OR RAW SHELL SEARCHES: Banned all `Get-ChildItem`, `dir`, `ls`, `rg`, `ripgrep`, `grep`, `Select-String`, and `findstr`.
- [ ] HEATMAP TELEMETRY INSPECTION: Pre-action intent logged before file reads/edits and outcomes verified upon completion.

---

## Final Step Git Commit & Push Mandate

At the final step of the turn, after all targeted files have been refactored, verified with targeted linters, plans enqueued in `.ai-memory/plans/pending/XX-<slug>.md`, and subtasks placed in `.ai-memory/plans/subtasks/XX-<slug>/`, use GitMap semantic commit commands exclusively:

```bash
gitmap cpf "<module> - find 100 improvements in whole system and enqueue plans"
```

- Format message strictly with a hyphen `-` to separate module and summary.
- NEVER include colons `:` inside the GitMap argument as GitMap automatically provides `Feature: ` or `Bug: `.
- ZERO intermediate commits during Phase 1 or Phase 2.
- Confirm `git status --porcelain` is completely clean and upstream tracking branch reflects the commit before concluding.

# Terminal UI, CLI Styling, Lipgloss & Animations — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-terminal-ui`, `cg-cli-style`, `cg-lipgloss`, `cg-execute terminal-ui`, `audit terminal ui`, `terminal colors`, `clone animation`, `cli help banners`

> [!IMPORTANT]
> Prompt Version: 2.1.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, plan, refactor, and beautify all Terminal UI components, CLI help renderers, progress animations, and status formatters across Go (Lipgloss/Pterm), TypeScript, Python, and Shell/PowerShell scripts, enforcing bright bold ANSI palettes, Catppuccin pastel cycling, responsive 2-column alignment, expandable subcommand markers, and clean version footers with strict relative Git paths until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/17-consolidated-guidelines/07-enum-standards.md` for cross-language enum and constant architectures.
15. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
16. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Terminal Renderers, Inventory Unformatted Outputs, Write .ai-memory/plans/pending/ Spec, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Refactor Terminal UI, Super-Category Banners, Lipgloss Alignment, Spinners, Footers, Verify CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Beautiful Terminal UI & CLI Styling Architecture

Modern CLI tools must provide an exceptional user experience on dark terminals (VS Code terminal, Windows Terminal, iTerm). Muddy colors, misaligned columns, and unbuffered stdout calls degrade user experience.

---

### 1. The ANSI & Pastel Color Palette Standards

1. **Bright & Bold ANSI Codes (`\033[1;9Xm`):**
   - Standard legacy 3X codes (`\033[32m`) look dull and muddy on modern dark backgrounds. Always use bright bold 9X codes:
     - `ColorGreen  = "\033[1;92m"` — Successes, checkmarks (`✔ clean`).
     - `ColorRed    = "\033[1;91m"` — Errors, failures (`● dirty`).
     - `ColorYellow = "\033[1;93m"` — Warnings, prompts, tips.
     - `ColorCyan   = "\033[1;96m"` — Paths, URLs, sub-group headings.
     - `ColorMagenta= "\033[1;95m"` — Version highlights, category banners.
     - `ColorDim    = "\033[2;37m"` — Secondary text, expandable markers.
     - `ColorReset  = "\033[0m"`   — Suffix on EVERY colored segment.
2. **TrueColor Catppuccin Pastel Palette (for Multi-Item List Cycling):**
   - When printing batches of repositories or jobs (e.g. `clone-all`, `pull-all`), rotate through pastel colors to visually separate distinct items:
     - `ColorPastelGreen   = "\033[38;2;166;227;161m"`
     - `ColorPastelCyan    = "\033[38;2;137;220;235m"`
     - `ColorPastelYellow  = "\033[38;2;249;226;175m"`
     - `ColorPastelMagenta = "\033[38;2;203;166;247m"`

---

### 2. Super-Category Intent Banners & Box-Drawing

Organize large CLI command catalogs into distinct visual tiers using bold intent banners so users can find commands immediately:

```text
  ━━ GET STARTED ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Scanning Repositories:
    scan .                       Recursively index Git repositories
    list                         Display tracked repository inventory

  ━━ WORK WITH REPOS ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Cloning & Worktrees:
    clone <url>                  Clone repository with auto-tagging
    clone-all                    Parallel batch clone of inventory
```

```go
// Go Implementation of Intent Banner with Lipgloss / Unicode Rule
func printSuperCategory(title string, body func()) {
    const superCategoryLineWidth = 58
    ruleLength := superCategoryLineWidth - len(title)
    if ruleLength < 4 {
        ruleLength = 4
    }
    rule := strings.Repeat("━", ruleLength)

    fmt.Println()
    fmt.Println("  " + ColorMagenta + "━━ " + ColorWhite + title + ColorReset + " " + ColorMagenta + rule + ColorReset)
    body()
}
```

---

### 3. Responsive 2-Column Help Alignment (Lipgloss 2-Pass Rendering)

To avoid jagged table borders and text clipping on small terminal windows:

1. **Pass 1 (Measurement):** Iterate through all command rows with `measuringHelp = true` to compute `maxHelpCmdLen` (capped at `maxCmdColumnWidth = 26`).
2. **Pass 2 (Render & Wrap):** Render rows using calculated padding. If the command exceeds the cap, place the description on a new indented line (`renderLongHelpRow`). Wrap long descriptions responsively based on `termWidth`.

```go
const maxCmdColumnWidth = 26

func renderStandardHelpRow(cmd, fullDesc string, cmdWidth, termWidth int) {
    pad := maxHelpCmdLen - cmdWidth
    if pad < 0 {
        pad = 0
    }

    prefix := fmt.Sprintf("  %s%s  ", cmd, strings.Repeat(" ", pad))
    descWidth := termWidth - lipgloss.Width(prefix)
    if descWidth <= 10 {
        fmt.Printf("%s%s\n", prefix, fullDesc)
        return
    }

    printWrappedHelpLines(prefix, fullDesc, descWidth)
}
```

---

### 4. Expandable Subcommand Markers & Clean Footers

1. **Expandable Marker (`▸ subcommands`):**
   - For commands that contain nested sub-commands (e.g. `gitmap ssh`, `gitmap project`), display an expandable hint:
     `▸ subcommands — see gitmap ssh --help`
2. **Standard Repository & Version Footer:**
   - At the bottom of CLI execution, display the version string and repository state:
     ```text
     gitmap v3.82.0 (build 7c49228) • https://github.com/alimtvnetwork/gitmap
     ```

---

## 5. Phase 1 Violation Ledger Format

In Phase 1, you MUST generate `.ai-memory/plans/pending/XX-terminal-ui-styling-audit.md` containing the master inventory table:

```markdown
| Command / View | File Path | Line | Current UI Pattern | Defect / Limitation | Planned UI Enhancement | Status |
|---|---|:---:|---|---|---|:---:|
| `rootusage.go` | `gitmap/cmd/rootusage.go` | 167 | Hardcoded 38 column width | Screen clipping on narrow terminals | Cap width at 26, use responsive wrapping | PENDING |
| `clone.go` | `gitmap/cmd/clone.go` | 85 | Raw un-styled clone logs | No visual spinner or pastel cycling | Add Lipgloss spinner & Catppuccin palette | PENDING |
```

---

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/readme.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
1. **Universal File Search:**
   ```bash
   gitmap find "<pattern>" [-ext <ext>]
   ```
2. **List Indexed Files & Substring Lookup:**
   ```bash
   gitmap list-files [pattern]
   gitmap find-files-any "<substring>"
   ```
3. **Stream File Content:**
   ```bash
   gitmap cat <filepath>
   ```
4. **Instant Code Walk Search:**
   ```bash
   gitmap search "<term>"
   ```

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.ai-memory/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/` and update `.ai-memory/plans/readme.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`, `09-cli-help-auditor.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .ai-memory/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Automated Naming & Style Fixers:** Use `python 03-ai-scripts/08-naming-autofixer.py` and `05-guideline-autofixer.py` to audit boolean prefixes and newlines.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/readme.md` using sequential script naming. For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.ai-memory/readme.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/readme.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **Bright Bold ANSI Palette:** Validated bright bold 9X color codes (`\033[1;92m`, `\033[1;96m`).
- [ ] **Super-Category Banners:** Implemented structured intent banners (`━━ SECTION ━━━━━━━━━━`).
- [ ] **Responsive 2-Column Alignment:** Implemented 2-pass width calculation and responsive description wrapping.
- [ ] **Expandable Markers:** Added `▸ subcommands` markers for subcommands with nested routes.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Blank Line Before `if`:** Exactly one blank line precedes every `if` statement (unless at the top of a block).
- [ ] **Blank Line After `}`:** Exactly one blank line follows every closing brace `}` (unless closing the enclosing block).
- [ ] **Blank Line Before `return`:** Exactly one blank line precedes `return` / `throw` in multi-line blocks.
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0 using guard clauses and early returns.
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] **Atomic File Recording:** Modified files recorded to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). Builds, tests, and CI runner are deferred to CI/CD fix.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `.ai-memory/coding-guidelines.md`.
- [ ] ANSI Palette: Verified bright bold ANSI escape sequences and reset suffixes.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Blank Line Before `if`: Verified blank line before every `if` statement across all modified files.
- [ ] Blank Line After `}`: Verified blank line after every closing brace `}` followed by code.
- [ ] Blank Line Before `return`: Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] Zero Nested `if`: Zero nested `if` statements (depth > 1).

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<type>(<scope>): <summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are known for future release verification.

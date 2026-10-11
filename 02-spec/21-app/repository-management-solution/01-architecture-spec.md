# Architecture Specification: Subtree Remediation Hierarchy & Unified Fleet Control

- **Feature Slug:** `repository-management-solution`
- **Module:** `cli/cmdpending`, `cli/constants`
- **Specification Version:** 1.0.0
- **Status:** APPROVED
- **Author:** Spec Agent 01
- **Target Release:** Minor Bump

---

## 1. Architectural Overview & Design Goals

The GitMap pending commits subsystem (`cli/cmdpending`) scans repositories across local workspaces and SSH cluster nodes, detecting uncommitted files and unpushed commits. When repositories are in a dirty state, the CLI outputs remediation options inside an 80-column bordered table.

This architecture accomplishes three primary objectives:
1. **Hierarchical Subtree Rendering:** Replaces flat inline suggestions with structured two-tier subtrees where each option node hosts its executable command as an indented child solution node.
2. **Unified Fleet Multi-Repo Control:** Integrates Option 3 (`All Repos (Commit & Push)`) targeting `all` repositories via `gitmap sends cp all "wip: save changes"`, empowering operators to remediate the entire workspace with a single command.
3. **Box Width & Layout Invariants:** Preserves the 78-rune inner table box width using rune-aware string truncation and padding, preventing table border misalignment across polyglot terminal emulators.

---

## 2. Subtree Topology & Mathematical Connector Calculation

### 2.1 Two-Tier Subtree Model

Each remediation entry is broken into two distinct rendering lines:
- **Parent Node (Option Definition):** Identifies the option number and descriptive label.
- **Child Node (Solution Invocation):** Indents the concrete GitMap CLI command under the option.

```text
│   ├── Option 1: Commit & Push WIP                                            │
│   │   └── gitmap sends cp <repo> "wip: save changes"                         │
│   ├── Option 2: Stash WIP                                                    │
│   │   └── gitmap stash <repo>                                                │
│   └── Option 3: All Repos (Commit & Push)                                    │
│       └── gitmap sends cp all "wip: save changes"                            │
```

### 2.2 Connector Calculation Mechanics

Let $N$ represent the total count of remediation options for a dirty repository ($N = 3$). For each option at index $i \in [0, N-1]$:

1. **Option Node Connectors:**
   - When $i < N - 1$ (non-terminal option):
     - Connector: `├──`
     - Formatted string: `   ├── Option %d: %s`
   - When $i == N - 1$ (terminal option):
     - Connector: `└──`
     - Formatted string: `   └── Option %d: %s`

2. **Child Solution Node Connectors:**
   - Under a non-terminal option ($i < N - 1$):
     - The vertical trunk continues downward to subsequent options.
     - Continuation prefix: `│   `
     - Formatted string: `   │   └── %s`
   - Under the terminal option ($i == N - 1$):
     - The trunk terminates at Option 3; no vertical line continues downward.
     - Terminal prefix: `    ` (four spaces)
     - Formatted string: `       └── %s` (seven leading spaces + `└── `)

This mathematical distinction (`│   └──` vs `    └──`) guarantees standard hierarchical tree aesthetics in monospaced terminal fonts.

---

## 3. Data Contracts & Model Evolution

### 3.1 Model Definition in `cli/cmdpending/pending_commits_types.go`

The existing `RemediationOption` struct captures both the label and executable command:

```go
type RemediationOption struct {
	OptionNumber int    `json:"optionNumber"`
	Label        string `json:"label"`
	Command      string `json:"command"`
}
```

- `OptionNumber`: 1-based sequential integer (`1`, `2`, `3`).
- `Label`: Human-readable remediation title (`Commit & Push WIP`, `Stash WIP`, `All Repos (Commit & Push)`).
- `Command`: Concrete executable GitMap invocation.

### 3.2 Option Triplet Construction in `cli/cmdpending/pending_commits_cmd.go`

`buildRepoRemediationOptions` constructs three distinct options for each dirty repository:

| Option # | Label | Target Scope | Command |
|---|---|---|---|
| **1** | `Commit & Push WIP` | Single Repo | `gitmap sends cp <target> "wip: save changes"` |
| **2** | `Stash WIP` | Single Repo | `gitmap stash <target>` |
| **3** | `All Repos (Commit & Push)` | Entire Workspace | `gitmap sends cp all "wip: save changes"` |

Single-repository targets are quoted when path names contain whitespace via `quoteRemediationTarget`.

---

## 4. Unified Fleet Multi-Repo Control Engine

### 4.1 Semantics of `gitmap sends cp all`

Option 3 provides workspace-wide remediation through `gitmap sends cp all "wip: save changes"`:
1. **Target Discovery:** The execution engine identifies all registered repositories in the active workspace.
2. **Dirty Filtering:** Repositories are inspected for uncommitted changes (staged, unstaged, untracked) using GitMap telemetry.
3. **Atomic Commit & Push:** For each dirty repository:
   - Untracked and modified files are staged.
   - A commit is authored with the standardized message `wip: save changes`.
   - The commit is pushed upstream to the tracked remote branch.
4. **Clean Skip Guarantee:** Repositories without pending changes are marked as `clean-skipped` without executing redundant network operations.
5. **Idempotence:** Re-running the command after resolution produces a clean pass across all nodes.

---

## 5. Box Width & Terminal Invariants (78 Runes)

### 5.1 The 78-Rune Invariant

The outer table boundary is fixed at 80 terminal columns:
- Left border: `  │` (2 leading spaces, 1 vertical border line)
- Content space: Exactly 78 runes
- Right border: `│` (1 vertical border line)
- Total line width: $2 + 1 + 78 + 1 = 82$ print positions (enclosing 80-column inner frame)

### 5.2 UTF-8 Rune-Aware Calculation

Terminal box drawing characters (`├──`, `└──`, `│`, `…`) consume multiple UTF-8 bytes (e.g., `├` is 3 bytes: `\xe2\x94\x9c`) but occupy exactly 1 column width in terminal displays.

The string formatter `formatTreeLine` enforces:
1. **Rune Conversion:** Converts raw strings into rune slices `[]rune(raw)` to ensure character-level indexing rather than byte indexing.
2. **Truncation with Ellipsis:** If `len(runes) > 78`, the slice is truncated to index 77 and appended with `…` (1 rune), maintaining exactly 78 visual columns.
3. **Space Padding:** If `len(runes) < 78`, trailing spaces are appended: `raw + strings.Repeat(" ", 78 - len(runes))`.

---

## 6. Component Architecture & File Responsibilities

| Subsystem File | Architectural Responsibility |
|---|---|
| `cli/cmdpending/pending_commits_types.go` | Telemetry models, `RemediationOption`, `RepoPendingCommitRecord`. |
| `cli/cmdpending/pending_commits_cmd.go` | Option generation, `buildRepoRemediationOptions`, single vs `all` scope resolution. |
| `cli/cmdpending/pending_commits_render.go` | Tree connector calculation, two-tier subtree loop, 78-rune box width enforcement. |
| `cli/cmdpending/pending_commits_render_test.go` | Unit testing subtree hierarchy, connector prefixes, and rune count invariants. |
| `cli/cmdpending/pending_commits_repo_test.go` | Unit testing Option 1, Option 2, Option 3 generation across dirty repositories. |

---

## 7. Implementation Guidelines & Quality Constraints

1. **Function Size:** Strict limit of 15 lines per function. Decompose rendering logic into modular helper functions:
   - `resolveOptionPrefix(isTerminal bool) string`
   - `resolveSolutionPrefix(isTerminal bool) string`
   - `renderOptionSubtreeRow(opt RemediationOption, isTerminal bool, border, reset string)`
2. **File Size:** All Go and markdown files must remain strictly $\le 300$ lines.
3. **Positive Booleans:** Boolean identifiers must use affirmative prefixes (`isTerminal`, `isDirty`, `hasUnpushed`). Comparison against `true` or negative polarity is prohibited.
4. **Relative Path Hygiene:** Only repository-relative paths (`cli/...`, `02-spec/...`, `.ai-memory/...`) are permitted in documentation and code comments.
5. **No Direct Git Commands:** All operations and documentation must reference GitMap CLI equivalents.

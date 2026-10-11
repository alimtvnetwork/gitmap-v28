# Master Specification: Repository Management Solution & Subtree Remediation Hierarchy

- **Feature Slug:** `repository-management-solution`
- **Module:** `cli/cmdpending`, `cli/constants`
- **Specification Version:** 1.0.0
- **Status:** APPROVED
- **Author:** Spec Agent 01
- **Target Release:** Minor Bump

---

## User Request (Verbatim)

```text
All option points, the option solution needs to be as a subtree as well. It's not like option one and it's going to show everything there. Don't do that. And with all these options, try to have one version that will control all these repositories. One option that we can do that will control all this or fix all these repositories, make sure of that. I like the table format. That's nice. I appreciate that.
```

---

## 1. Problem Statement & System Context

During workspace repository status inspection (`gitmap pc`, `gitmap pending-commits`), dirty repositories displaying uncommitted files or unpushed commits present remediation hints within the 80-column summary table.

Prior implementations suffered from two primary user experience and architectural deficiencies:

1. **Flat Single-Line Option Formatting:**
   - Previous versions rendered remediation suggestions as flat lines combining the option number, action title, and command into a single crowded text string:
     ```text
     │   ├── Option 1: gitmap sends cp repo-a "wip: save changes"                 │
     │   └── Option 2: gitmap stash repo-a                                        │
     ```
   - This flat layout cluttered terminal readability, failed to visually separate the conceptual remediation strategy from the raw shell invocation, and violated hierarchical subtree principles.

2. **Absence of Unified Multi-Repo Control:**
   - Operators inspecting dozens of dirty repositories had no single command option to remediate the entire workspace fleet in one action.
   - Remediation options only offered per-repository targets (`<repo>`), forcing operators to manually repeat commands across every individual repository.

---

## 2. Core Functional Requirements

### FR-01: Two-Tier Hierarchical Subtree Formatting
Remediation hints for each dirty repository must be structured as a two-tier hierarchical subtree inside the 80-column table box:
- **Tier 1 (Option Node):** Describes the remediation strategy with an option number and title.
  - Non-terminal branch: `├── Option 1: Commit & Push WIP`
  - Intermediate branch: `├── Option 2: Stash WIP`
  - Terminal branch: `└── Option 3: All Repos (Commit & Push)`
- **Tier 2 (Child Solution Node):** Emits the exact executable GitMap command indented underneath its corresponding Option Node:
  - Non-terminal child: `│   └── gitmap sends cp <repo> "wip: save changes"`
  - Intermediate child: `│   └── gitmap stash <repo>`
  - Terminal child: `    └── gitmap sends cp all "wip: save changes"`

### FR-02: Unified Fleet Multi-Repo Control (Option 3)
- Every dirty repository record must include Option 3: `All Repos (Commit & Push)`.
- The child solution command is `gitmap sends cp all "wip: save changes"`.
- This provides operators with a single, immediate command to commit and push WIP across all dirty repositories in the workspace in one shot.

### FR-03: Dynamic Connector Calculation
The rendering engine must dynamically calculate tree connectors based on terminal status:
- For an Option Node at index `i` out of `N` options:
  - Non-terminal (`i < N-1`): uses `├── `
  - Terminal (`i == N-1`): uses `└── `
- For a Child Solution Node at index `i`:
  - Non-terminal parent (`i < N-1`): trunk continuation prefix `│   └── `
  - Terminal parent (`i == N-1`): empty trunk prefix `    └── `

### FR-04: Strict 78-Rune Box Width Invariant
- The terminal table features an 80-column width (78 inner runes enclosed by `│` borders).
- All formatted tree lines must be measured using UTF-8 rune counts via `formatTreeLine`.
- Overlong lines exceeding 78 runes must be truncated to 77 runes and terminated with an ellipsis (`…`).
- Lines shorter than 78 runes must be padded with trailing spaces to maintain perfect border alignment.

---

## 3. Terminal Visual Topology

```text
  ┌──────────────────────────────────────────────────────────────────────────────┐
  │                    GITMAP PENDING COMMITS SUMMARY                            │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ Scanned: 43 repos │ Dirty: 2 repos │ Uncommitted: 5 files │ Unpushed: 1      │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ REPOSITORY              VER/BRANCH        UNCOMMITTED    UNPUSHED  STATUS    │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ gitmap                  main                        3           0  ● PEND    │
  │   ├── Option 1: Commit & Push WIP                                            │
  │   │   └── gitmap sends cp gitmap "wip: save changes"                         │
  │   ├── Option 2: Stash WIP                                                    │
  │   │   └── gitmap stash gitmap                                                │
  │   └── Option 3: All Repos (Commit & Push)                                    │
  │       └── gitmap sends cp all "wip: save changes"                            │
  │ docs-site               main                        2           1  ● PEND    │
  │   ├── Option 1: Commit & Push WIP                                            │
  │   │   └── gitmap sends cp docs-site "wip: save changes"                      │
  │   ├── Option 2: Stash WIP                                                    │
  │   │   └── gitmap stash docs-site                                             │
  │   └── Option 3: All Repos (Commit & Push)                                    │
  │       └── gitmap sends cp all "wip: save changes"                            │
  ├──────────────────────────────────────────────────────────────────────────────┤
  │ ○ 41 clean repositories omitted (use --all to inspect)                       │
  └──────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Architectural Invariants & Quality Standards

1. **File Size Limit:** No source or documentation file may exceed 300 lines.
2. **Function Length Limit:** All Go functions must not exceed 15 lines.
3. **Boolean Naming Conventions:** Positive polarity booleans only (`isTerminal`, `hasChanges`, `isDirty`). Explicit comparisons (`== true`) and mixed polarity are strictly forbidden.
4. **Strict Relative Paths:** All file links and paths must use repository-relative forward slashes (`02-spec/...`, `.ai-memory/...`, `cli/...`). Absolute paths and URI schemes (`file:///`) are forbidden.
5. **Direct Git Command Ban:** Direct invocations of `git add`, `git commit`, `git status`, or `git push` are strictly forbidden. All operations must route through GitMap CLI commands (`gitmap sends cp`, `gitmap stash`).
6. **GitMap Tooling:** Code searches must use `gitmap aum search` and `gitmap find`.

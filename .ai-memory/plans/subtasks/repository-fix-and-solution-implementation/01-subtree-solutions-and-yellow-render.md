# Engineering Subtask Plan: Subtree Solutions Rendering & Vibrant Yellow Terminal Palette

**Subtask ID:** `01-subtree-solutions-and-yellow-render`  
**Subtask Code:** `Task-01-Subtree-Solutions`  
**Parent Plan:** `repository-fix-and-solution-implementation` ([repository-fix-and-solution-implementation.md](../../repository-fix-and-solution-implementation.md))  
**Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/repository-fix-and-solution-implementation/01-architecture-spec.md)  
**Companion Subtask:** [02-non-repo-smart-clone-and-fix-engine.md](02-non-repo-smart-clone-and-fix-engine.md)  
**Assigned Agent Role:** Spec Subagent 01 / Implementation Worker 01  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory

| Component Area | Primary Implementation Files | Test & Documentation Files |
| :--- | :--- | :--- |
| **Pull Results Renderer** | `cli/cmdpull/pull_efficient_render.go` | `cli/cmdpull/pull_efficient_render_test.go` |
| **Remediation Hint Resolver** | `cli/cmdpull/pull_remediation_hint.go` | `cli/cmdpull/pull_remediation_forward_slash_test.go` |
| **Terminal Constants** | `cli/constants/constants_terminal.go` | `cli/cmdpull/pull_efficient_render_test.go` |

---

## 2. Disjoint File Ownership & Strict Boundaries

To ensure safe multi-agent execution and prevent file modification collisions:

- **Owned Implementation Files (for Worker 01):**
  * `cli/cmdpull/pull_efficient_render.go`
  * `cli/cmdpull/pull_remediation_hint.go`
  * `cli/cmdpull/pull_efficient_render_test.go`
- **Downstream Owned Files (Reserved exclusively for Worker 02 / Subtask 02):**
  * `cli/cloner/pulldiag.go`
  * `cli/cmdfix/fix_cmd.go`
  * `cli/cmdautofix/fix.go`
- **Strictly Prohibited Actions:**
  * **TOTAL BAN ON GIT COMMANDS:** Subagents MUST NOT execute `git add`, `git commit`, `git status`, `git push`, or `git diff`.
  * **Strictly Relative Paths Only:** All file references across documentation, code comments, and test fixtures must use forward-slash paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute filesystem paths or `file:///` URIs.
  * **Search Exclusively via GitMap:** Use `gitmap aum search`, `gitmap find`, and `gitmap cat`. Strict ban on `grep`, `ripgrep`, `rg`, `Select-String`.
  * **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.

---

## 3. Step-by-Step Implementation Tasks

### Step 1: Deprecate and Remove Flat `Next Step:` Colon Line in `cli/cmdpull/pull_efficient_render.go`
* **Objective:** Completely eradicate the unstructured `Next Step: <cmd1> or <cmd2>` line with inline "or".
* **Exact Modifications:**
  1. In `renderSingleFailedItem`:
     - Remove the call to `ResolvePullRemediationHint(s)`:
       ```go
       // REMOVE lines:
       remHint := ResolvePullRemediationHint(s)
       if remHint == "" {
           remHint = fmt.Sprintf("gitmap status %s or gitmap fix %s", s.RepoName, s.RepoName)
       }
       fmt.Fprintf(w, "    %sNext Step: %s%s%s\n", treeBranch, constants.ColorCyan, remHint, constants.ColorReset)
       ```
     - Directly route rendering into `renderStructuredSolutionsSubtree(w, s)`.
* **Acceptance Criteria:**
  - Pull output for failed repositories never contains the token `"Next Step:"`.

---

### Step 2: Implement `renderStructuredSolutionsSubtree` in `cli/cmdpull/pull_efficient_render.go`
* **Objective:** Format `Solutions:` as a nested subtree with all option labels, titles, and commands highlighted in vibrant yellow (`constants.ColorYellow`).
* **Exact Modifications:**
  1. Replace `renderStructuredOptions` with `renderStructuredSolutionsSubtree(w io.Writer, s *PullRepoState)`.
  2. Emit the subtree root branch:
     ```go
     fmt.Fprintf(w, "    %sSolutions:\n", treeBranch)
     ```
  3. Resolve structured remediation using `ResolveStructuredRemediation(s)`.
  4. If `len(structured.Options) == 0`, generate synthetic fallback dual options:
     - Option 1 (Auto-Fix): `gitmap fix ` + repoName
     - Option 2 (Inspect Status): `gitmap status ` + repoName
  5. Iterate over options:
     - If option index is the final option (`i == len(opts)-1`), connector is `treeTerminal` (`└── `), otherwise `treeBranch` (`├── `).
     - Format each option line using `treeContinuation` (`│   `) and bright yellow styling:
       ```go
       fmt.Fprintf(w, "    %s%s%sOption %d (%s): %s%s\n",
           treeContinuation, connector,
           constants.ColorYellow, opt.OptionNumber, opt.Title, opt.Command, constants.ColorReset)
       ```
* **Acceptance Criteria:**
  - `Solutions:` appears with `├── Solutions:`.
  - Nested sub-items indent with `│   ├── Option 1 (...)` and `│   └── Option 2 (...)`.
  - ANSI yellow code `\033[1;93m` wraps each option line.

---

### Step 3: Implement Unified Batch Command Footer in `renderFailedGroup`
* **Objective:** Provide a single prominent batch resolution command at the conclusion of failed repositories list to resolve all failed repositories simultaneously.
* **Exact Modifications:**
  1. In `renderFailedGroup`, after the loop over `deduped` failed repository states:
     ```go
     if len(deduped) > 0 {
         fmt.Fprintf(w, "\n  %s💡 Unified Resolution:%s %sgitmap fix --all%s %s(or: %sgitmap pull-fix --all%s)%s\n",
             constants.ColorBold, constants.ColorReset,
             constants.ColorYellow, constants.ColorReset,
             constants.ColorDim, constants.ColorYellow, constants.ColorDim, constants.ColorReset)
     }
     ```
* **Acceptance Criteria:**
  - When 1 or more failed repositories are rendered, the bottom of the section displays `💡 Unified Resolution: gitmap fix --all (or: gitmap pull-fix --all)` in yellow.
  - When 0 failed repositories exist, the footer is not emitted.

---

### Step 4: Refactor Option Hints in `cli/cmdpull/pull_remediation_hint.go`
* **Objective:** Eliminate inline "or" strings and guarantee clean dual options across all failure classifications.
* **Exact Modifications:**
  1. Review all dual hint resolvers (`resolveConflictDualHints`, `resolveDivergedDualHints`, `resolveUntrackedDualHints`, `resolveDirtyTreeDualHints`, `resolveAuthDualHints`, `resolveMissingRepoDualHints`, `resolveFallbackDualHints`).
  2. Guarantee that each resolver returns two distinct, executable options without inline "or" disjunctions in the title or command strings.
  3. Verify that all repository paths within commands are normalized with `filepath.ToSlash` to avoid double-backslash escapes on Windows.
* **Acceptance Criteria:**
  - Zero option titles or commands contain the substring `" or "`.
  - All file and directory paths in commands use forward slashes `/`.

---

### Step 5: Update & Expand Test Suite in `cli/cmdpull/pull_efficient_render_test.go`
* **Objective:** Prevent regressions and ensure automated verification of the subtree structure, ANSI yellow coloring, and unified batch command.
* **Exact Modifications:**
  1. Update existing `TestRenderConciseActiveResultsTo`:
     - Remove assertion checking for `"Next Step:"`.
     - Add assertion verifying `"Solutions:"` is present.
     - Add assertion verifying `constants.ColorYellow` is present.
     - Add assertion verifying `💡 Unified Resolution: gitmap fix --all` is present.
  2. Implement new targeted test: `TestRenderFailedGroup_SubtreeAndYellow`:
     - Set up mock `PullRepoState` with failure details.
     - Execute `RenderConciseActiveResultsTo` into a buffer.
     - Verify exact tree lines:
       * `├── Solutions:`
       * `│   ├── Option 1`
       * `│   └── Option 2`
       * `└── Diagnostic:`
     * Verify presence of `\033[1;93m` (yellow).
* **Acceptance Criteria:**
  - Tests pass cleanly with zero failures or race conditions.

---

## 4. Verification & Quality Gates

| Gate ID | Verification Item | Success Criteria |
| :--- | :--- | :--- |
| **VG-01** | Subtree Solution Layout | `Solutions:` rendered as nested subtree; zero flat colon lines with inline "or". |
| **VG-02** | Yellow Option Palette | `Option 1` and `Option 2` titles and commands wrapped in `constants.ColorYellow`. |
| **VG-03** | Batch Resolution Command | `💡 Unified Resolution: gitmap fix --all (or: gitmap pull-fix --all)` emitted at summary bottom. |
| **VG-04** | Total Removal of `Next Step:` | Zero occurrences of `"Next Step:"` in pull output. |
| **VG-05** | Relative Path Hygiene | 100% relative Git paths across all documentation, comments, and tests. |
| **VG-06** | Test Suite Parity | Unit tests verify subtree layout, ANSI colors, and batch command presence. |

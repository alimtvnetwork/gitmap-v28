# Subtask 04: Failure Tree Subtrees & Dual Options Remediation

> **Parent Plan:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`  
> **Status:** READY  
> **Target Files:**  
> - `cli/cmdpull/pull_efficient_render.go`  
> - `cli/cmdpull/pull_remediation_hint.go`  
> - `cli/cmdpull/pull_efficient_render_test.go`  
> - `cli/cmdpull/pull_remediation_hint_test.go`  

---

## 1. Objectives

1. **Tree Connector Refactoring:**
   Replace flat repetitive `↳` indentation in `renderSingleFailedItem` and `renderStructuredOptions` within `cli/cmdpull/pull_efficient_render.go` with hierarchical Unicode box connectors (`├──`, `└──`, `│   `) matching the pipeline failure tree design from `cli/cmdpipeline/pipeline_failure_tree.go`.
2. **Missing Repository Dual Options Logic:**
   Add `resolveMissingRepoDualHints(repoDir, repoName string)` in `cli/cmdpull/pull_remediation_hint.go`. Ensure missing repositories return:
   - **Option 1:** `gitmap clone <repo>` ("Clone Missing Repository")
   - **Option 2:** `gitmap rm --db-only <repo>` ("Deregister from Database")
   Eliminate the erroneous fallback to `gitmap status` and `gitmap pull` when repository folders are missing on disk.
3. **Comprehensive Pattern Matching:**
   Enhance `isMissingRepoFailure` to detect missing directory failures across Linux and Windows: `"missing repository directory"`, `"directory does not exist"`, `"no such file or directory"`, and `"cannot change to"`.
4. **Unit Test Verification:**
   Verify tree connector rendering and missing repo dual hint resolution via unit tests in `pull_efficient_render_test.go` and `pull_remediation_hint_test.go`.

---

## 2. Implementation Steps

### Step 2.1: Implement Dual Options for Missing Repositories (`cli/cmdpull/pull_remediation_hint.go`)

1. Define `resolveMissingRepoDualHints`:
   ```go
   func resolveMissingRepoDualHints(repoDir, repoName string) (string, string, string, string) {
       cloneCmd := fmt.Sprintf("gitmap clone %s", repoName)
       removeCmd := fmt.Sprintf("gitmap rm --db-only %s", repoName)
       return "Clone Missing Repository", cloneCmd, "Deregister from Database", removeCmd
   }
   ```
2. Update `ResolveDualPullRemediationHints`:
   Prioritize `isMissingRepoFailure(msg)` before diverged, untracked, dirty, and auth checks:
   ```go
   if isMissingRepoFailure(msg) {
       return resolveMissingRepoDualHints(repoDir, repoName)
   }
   ```
3. Update `isMissingRepoFailure`:
   ```go
   func isMissingRepoFailure(msg string) bool {
       lower := strings.ToLower(msg)
       return strings.Contains(lower, "missing repository directory") ||
           strings.Contains(lower, "directory does not exist") ||
           strings.Contains(lower, "no such file or directory") ||
           strings.Contains(lower, "cannot change to")
   }
   ```

### Step 2.2: Refactor Failure Tree Rendering (`cli/cmdpull/pull_efficient_render.go`)

1. Define tree connector constants in `cmdpull`:
   - `treeBranch = "├── "`
   - `treeTerminal = "└── "`
   - `treeContinuation = "│   "`
   - `treeIndent = "    "`
2. Refactor `renderSingleFailedItem`:
   - Line 1: `FormatConciseActiveResultLine(colWidth, displayName, "failed")`
   - Line 2 (Reason): `fmt.Fprintf(w, "    %s%sReason:%s %s\n", treeBranch, constants.ColorDim, constants.ColorReset, errDetails)`
   - Line 3 (Next Step): `fmt.Fprintf(w, "    %s%sNext Step:%s %s%s%s\n", treeBranch, constants.ColorCyan, constants.ColorReset, constants.ColorDim, remHint, constants.ColorReset)`
   - Lines 4+ (Options Subtree):
     - Subtree Root: `fmt.Fprintf(w, "    %sOptions:\n", treeBranch)`
     - For Option `i` of `N`:
       - If `i < N-1`: `fmt.Fprintf(w, "    %s%sOption %d (%s):%s %s%s%s\n", treeContinuation, treeBranch, opt.OptionNumber, opt.Title, ...)`
       - If `i == N-1`: `fmt.Fprintf(w, "    %s%sOption %d (%s):%s %s%s%s\n", treeContinuation, treeTerminal, opt.OptionNumber, opt.Title, ...)`
   - Terminal Line (Diagnostic hint):
     - `fmt.Fprintf(w, "    %sTo inspect stack trace: %sgitmap pull-error %s%s (or: %sgitmap pe%s)\n", treeTerminal, constants.ColorYellow, s.RepoName, constants.ColorReset, constants.ColorDim, constants.ColorReset)`

### Step 2.3: Unit Test Suite Verification

1. Add tests in `cli/cmdpull/pull_remediation_hint_test.go`:
   - Test `ResolveDualPullRemediationHints` with `"missing repository directory"` -> assert Option 1 is `gitmap clone` and Option 2 is `gitmap rm --db-only`.
   - Test `ResolveDualPullRemediationHints` with `"fatal: cannot change to '/work/repo': No such file or directory"` -> assert matches missing repo hints.
2. Add tests in `cli/cmdpull/pull_efficient_render_test.go`:
   - Create failed `PullRepoState` and execute `RenderConciseActiveResultsTo`.
   - Assert buffer contains `├── Reason:`, `├── Next Step:`, `├── Options:`, `│   ├── Option 1`, `│   └── Option 2`, and `└── To inspect stack trace:`.

---

## 3. Acceptance Criteria

- [ ] `resolveMissingRepoDualHints` implemented and exported/tested.
- [ ] `ResolveDualPullRemediationHints` branches into `resolveMissingRepoDualHints` on missing directory errors.
- [ ] Failure rendering produces box tree connectors (`├──`, `└──`, `│   `) with no loose flat `↳` prefixes for failed repos.
- [ ] Diagnostic hint `└── To inspect stack trace: gitmap pull-error <repo> (or: gitmap pe)` renders as the tree terminus.
- [ ] `gofmt -w` cleanly applied with 0 lint violations.
- [ ] All unit tests pass with `go test ./cli/cmdpull/...`.

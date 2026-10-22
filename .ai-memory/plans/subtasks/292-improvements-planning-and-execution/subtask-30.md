# Subtask Plan 30: Operational Handler Extraction & Route Decoupling (`cli/cmd/see_ops.go`)

> **Subtask ID:** `subtask-30`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmd/see_ops.go`, `cli/cmd/cmd_see.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Extract operational handlers and execution functions from `cli/cmd/cmd_see.go` into a new companion file `cli/cmd/see_ops.go` within package `cmd`.
2. Move functions:
   - `runSeeCommitPending() *apperror.AppError`
   - `printPendingCommits(dirty []cmdcpar.DirtyRepoSummary)`
   - `runSeeHistory(args []string) *apperror.AppError`
   - `runSeeHistorySSH(args []string) *apperror.AppError`
   - `resolveHistorySSHArgs(args []string) []string`
   - `RunSeeErrorsSSH(args []string) *apperror.AppError`
   - `runSeeErrorsSSH(args []string) *apperror.AppError`
   - `resolveErrorsSSHArgs(args []string) []string`
   - `RunRepoManageUI() *apperror.AppError`
   - `PrintSeeHelp() *apperror.AppError`
3. Decouple routing dispatch (`cmd_see.go`) from operational execution (`see_ops.go`), reducing `cmd_see.go` from 204 lines to ~100 lines.
4. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
`cli/cmd/cmd_see.go` currently mixes two separate architectural concerns:
- Subcommand routing dispatch (parsing CLI argument tokens, matching command aliases).
- Subprocess execution and UI orchestration (running `gitmap interactive`, executing SSH fleet PAS commands, formatting pending commit tables).

By extracting execution logic into `cli/cmd/see_ops.go`:
1. `cmd_see.go` becomes a pure routing table of ~100 lines.
2. `see_ops.go` encapsulates the subprocess and fleet execution handlers (~110 lines).
3. Both files remain cleanly bounded under modular sizing targets, with zero circular dependency since both reside in `package cmd`.

---

## 3. Concrete Implementation Details

### 3.1 New File `cli/cmd/see_ops.go`
Create `cli/cmd/see_ops.go` in package `cmd`:

```go
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcpar"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmderrors"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runSeeCommitPending() *apperror.AppError {
	dirty := cmdcpar.CollectDirtyRepositories()
	isClean := len(dirty) == 0
	if isClean {
		fmt.Printf("%s✓ No repositories with pending commits.%s\n", constants.ColorGreen, constants.ColorReset)
		return nil
	}

	printPendingCommits(dirty)

	return nil
}

func printPendingCommits(dirty []cmdcpar.DirtyRepoSummary) {
	fmt.Printf("\n%s  === PENDING COMMITS (%d repositories) ===%s\n", constants.ColorYellow, len(dirty), constants.ColorReset)
	for _, d := range dirty {
		fmt.Printf("  • %-30s | modified: %d | untracked: %d | staged: %d\n",
			d.RepoName, d.Diagnosis.ModifiedCount, d.Diagnosis.UntrackedCount, d.Diagnosis.StagedCount)
	}

	fmt.Println()
}

func runSeeHistory(args []string) *apperror.AppError {
	hasHook := HistoryRunnerFn != nil
	if hasHook {
		return HistoryRunnerFn(args)
	}

	err := exec.Command("gitmap", append([]string{"history"}, args...)...).Run()

	return apperror.WrapSimple(err, "history")
}

func runSeeHistorySSH(args []string) *apperror.AppError {
	subArgs := resolveHistorySSHArgs(args)
	err := cmdssh.RunFleetPASCommand("history", "gitmap history --limit 10", func() error {
		hasHook := HistoryRunnerFn != nil
		if hasHook {
			return HistoryRunnerFn(subArgs)
		}

		return exec.Command("gitmap", append([]string{"history"}, subArgs...)...).Run()
	})

	return apperror.WrapSimple(err, "history-ssh")
}

func resolveHistorySSHArgs(args []string) []string {
	isLonger := len(args) > 2
	if isLonger {
		return args[2:]
	}

	return []string{}
}

// RunSeeErrorsSSH executes see errors across local host and remote fleet nodes.
func RunSeeErrorsSSH(args []string) *apperror.AppError {
	err := cmdssh.RunFleetPASCommand("see errors", "gitmap errors --limit 10", func() error {
		return cmderrors.RunErrorsCLI(args)
	})

	return apperror.WrapSimple(err, "errors-ssh")
}

func runSeeErrorsSSH(args []string) *apperror.AppError {
	subArgs := resolveErrorsSSHArgs(args)

	return RunSeeErrorsSSH(subArgs)
}

func resolveErrorsSSHArgs(args []string) []string {
	isLonger := len(args) > 2
	if isLonger {
		return args[2:]
	}

	return []string{}
}

// RunRepoManageUI opens the repository management interactive terminal dashboard.
func RunRepoManageUI() *apperror.AppError {
	fmt.Println("Launching GitMap Repository Management Dashboard...")
	cmd := exec.Command("gitmap", "interactive")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	hasErr := err != nil
	if hasErr {
		cmdWeb := exec.Command("gitmap", "ui")
		_ = cmdWeb.Start()
	}

	return nil
}

// PrintSeeHelp displays help for gitmap see suite.
func PrintSeeHelp() *apperror.AppError {
	fmt.Printf(`
  GitMap Inspection Suite (see / c)
    • gitmap see commit pending          - List repositories with uncommitted changes
    • gitmap see git-ignore issues       - List repositories with ignore/duplicate issues
    • gitmap see errors                  - Display recent execution errors
    • gitmap see history                 - Display command history
    • gitmap see errors ssh (ses)        - Display errors across SSH fleet (PAS Formula)
    • gitmap history ssh                 - Display history across SSH fleet (PAS Formula)
    • gitmap nodes history               - Display history across SSH fleet nodes
    • gitmap repo-manage ui              - Launch interactive terminal dashboard
`)

	return nil
}
```

### 3.2 Update `cli/cmd/cmd_see.go`
Remove lines 104 through 204 from `cli/cmd/cmd_see.go`. All calls from `cmd_see.go` resolve transparently to `see_ops.go` within `package cmd`.

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All functions return `*apperror.AppError` directly. No `cliexit.HandleError` or `os.Exit` calls in leaf functions.
- **Rule 2 (Function Limit):** Every function <= 15 lines (e.g. `RunSeeErrorsSSH` is 7 lines, `RunRepoManageUI` is 14 lines).
- **Rule 3 (Zero Nested Ifs):** Linear guard clauses with early returns (`isClean`, `hasHook`, `isLonger`, `hasErr`).
- **Rule 4 (Positive Booleans):** Positive naming used throughout (`isClean`, `hasHook`, `isLonger`, `hasErr`).
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [ ] File `cli/cmd/see_ops.go` created in `package cmd`.
- [ ] 10 operational functions and handlers moved into `see_ops.go`.
- [ ] Extracted functions removed from `cli/cmd/cmd_see.go`.
- [ ] `cmd_see.go` line count reduced to ~100 lines.
- [ ] All functions satisfy <= 15 lines and zero nested ifs.
- [ ] Zero build or test commands executed.

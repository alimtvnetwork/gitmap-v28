# Subtask Plan 29: Rerun Function Decomposition (`cli/cmdpe/pe_rerun.go`)

> **Subtask ID:** `subtask-29`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpe/pe_rerun.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Decompose `RunPeRerun` in `cli/cmdpe/pe_rerun.go` (currently 55 lines, violating Rule 2) into discrete helper functions.
2. Ensure every single function in `cli/cmdpe/pe_rerun.go` is strictly <= 15 lines (target <= 10 lines).
3. Preserve identical CLI flag parsing, validation, dry-run display, and process execution behavior.
4. Adhere strictly to Error Return Sovereignty (Rule 1) and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
The current implementation of `RunPeRerun` combines:
- Argument pre-scanning for help flags.
- `flag.FlagSet` instantiation and parsing.
- Argument slice length and run ID validation.
- GitHub CLI path availability checks.
- Command argument array construction.
- Dry-run branching and formatted output.
- Subprocess execution and output printing.

This monolithic 55-line structure violates the project's strict 15-line function limit and makes testing individual phases impossible without full subprocess execution.

---

## 3. Concrete Implementation Details

### 3.1 Refactoring in `cli/cmdpe/pe_rerun.go`

Replace the monolithic `RunPeRerun` with modular helper functions:

```go
type peRerunConfig struct {
	repoFlag string
	dryRun   bool
	yes      bool
	runID    string
}

// RunPeRerun implements `gitmap pe rerun <run-id>` (alias: perr).
func RunPeRerun(args []string) *appfault.AppError {
	if isPeRerunHelp(args) {
		fmt.Print(peRerunHelp)
		return nil
	}

	cfg, parseErr := parsePeRerunFlags(args)
	if parseErr != nil {
		return parseErr
	}

	return executePeRerun(cfg)
}

func isPeRerunHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}

	return false
}

func parsePeRerunFlags(args []string) (peRerunConfig, *appfault.AppError) {
	fs := flag.NewFlagSet("pe rerun", flag.ContinueOnError)
	repoFlag := fs.String("repo", "", "owner/repo target")
	dryRun := fs.Bool("dry-run", true, "print the plan without executing")
	yes := fs.Bool("yes", false, "actually request the re-run")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(peRerunHelp)
		return peRerunConfig{}, apperror.WrapSimple(parseErr, "pe.rerun.flags")
	}

	return validatePeRerunArgs(fs.Args(), *repoFlag, *dryRun, *yes)
}

func validatePeRerunArgs(rest []string, repo string, dryRun, yes bool) (peRerunConfig, *appfault.AppError) {
	if len(rest) == 0 || strings.TrimSpace(rest[0]) == "" {
		fmt.Print(peRerunHelp)
		return peRerunConfig{}, appfault.NewValidationError("missing required argument: run-id")
	}

	return peRerunConfig{
		repoFlag: strings.TrimSpace(repo),
		dryRun:   dryRun,
		yes:      yes,
		runID:    strings.TrimSpace(rest[0]),
	}, nil
}

func executePeRerun(cfg peRerunConfig) *appfault.AppError {
	if !hasGh() {
		return appfault.NewValidationError("pe rerun requires the GitHub CLI (gh) on PATH")
	}

	ghArgs := buildPeRerunGhArgs(cfg)
	if cfg.dryRun && !cfg.yes {
		printPeRerunPlan(ghArgs, cfg.runID)
		return nil
	}

	return executePeRerunCmd(ghArgs, cfg.runID)
}

func buildPeRerunGhArgs(cfg peRerunConfig) []string {
	ghArgs := []string{"run", "rerun", cfg.runID, "--failed"}
	if cfg.repoFlag != "" {
		ghArgs = append(ghArgs, "-R", cfg.repoFlag)
	}

	return ghArgs
}

func printPeRerunPlan(ghArgs []string, runID string) {
	planLine := "gh " + strings.Join(ghArgs, " ")
	fmt.Println("pe rerun plan (dry-run, no changes made):")
	fmt.Printf("  would execute: %s\n", planLine)
	fmt.Printf("  scope: failed jobs only, run %s\n", runID)
	fmt.Println("  re-run with --yes to execute.")
}

func executePeRerunCmd(ghArgs []string, runID string) *appfault.AppError {
	planLine := "gh " + strings.Join(ghArgs, " ")
	fmt.Printf("requesting re-run: %s\n", planLine)
	cmd := exec.Command("gh", ghArgs...)
	out, runErr := cmd.CombinedOutput()
	fmt.Print(string(out))
	if runErr != nil {
		return apperror.WrapSimple(runErr, "pe.rerun.exec")
	}

	fmt.Printf("re-run requested for run %s (failed jobs only)\n", runID)

	return nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All validation and execution errors bubble directly via `*appfault.AppError`. No calls to `cliexit.HandleError`.
- **Rule 2 (Function Limit):** All 8 functions are strictly <= 15 lines (e.g. `RunPeRerun` is 10 lines).
- **Rule 3 (Zero Nested Ifs):** Guard clauses on help, flag parse errors, missing arguments, and dry-run execution.
- **Rule 4 (Positive Booleans):** Clean positive polarity across all branches.
- **Rule 6 (Vertical Whitespace):** Blank line before `return` and after `}` throughout.

---

## 5. Verification Criteria
- [x] `RunPeRerun(args []string) *appfault.AppError` is <= 15 lines.
- [x] Every decomposed helper function is <= 15 lines.
- [x] Dry-run plan formatting and live execution flags (`--yes`) remain identical.
- [x] All file paths are strictly relative lowercase Git paths.

# Subtask Plan 30: Diagnostics Triage Decomposition & Loop Flattening (`cli/cmdpe/pe_triage.go`)

> **Subtask ID:** `subtask-30`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpe/pe_triage.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Flatten nested loops in `worstFailureDetail` in `cli/cmdpe/pe_triage.go` by decomposing job and step extraction into discrete helper functions.
2. Decompose oversized `RunPeTriage` function (currently 100 lines) into small, focused functions adhering to the <= 15 line rule (Rule 2).
3. Ensure zero nested conditionals (Rule 3) and proper error return sovereignty (Rule 1).
4. Maintain 100% output parity for `gitmap pe triage` and alias `gitmap pet`.

---

## 2. Targeted Problem & Architectural Rationale
In `cli/cmdpe/pe_triage.go`:
1. **Nested Loops Violation:** `worstFailureDetail` iterates over `view.Jobs` and nests a secondary loop iterating over `j.Steps`:
   ```go
   for _, j := range view.Jobs {
       if !strings.EqualFold(j.Conclusion, "failure") {
           continue
       }
       for _, s := range j.Steps {
           if strings.EqualFold(s.Conclusion, "failure") {
               return j.Name, s.Name
           }
       }
       return j.Name, "unknown"
   }
   ```
   This violates the nested control-flow prohibition.
2. **Oversized Function Violation:** `RunPeTriage` spans lines 58 to 157 (100 lines), severely violating Rule 2 (maximum 15 lines per function). It conflates flag parsing, repo listing, probe loops, sorting, and terminal output rendering into a single monolithic block.

---

## 3. Concrete Implementation Details

### 3.1 Flatten `worstFailureDetail` in `cli/cmdpe/pe_triage.go`
Decompose into `worstFailureDetail`, `findFirstFailedJobAndStep`, and `findFirstFailedStep`:

```go
// worstFailureDetail returns the first failed job and its first failed step.
func worstFailureDetail(slug string, runID int64) (string, string) {
	var view struct {
		Jobs []ghJob `json:"jobs"`
	}
	id := fmt.Sprintf("%d", runID)
	if viewErr := ghJSON(&view, "run", "view", id, "-R", slug, "--json", "jobs"); viewErr != nil {
		return "unknown", "unknown"
	}

	return findFirstFailedJobAndStep(view.Jobs)
}

func findFirstFailedJobAndStep(jobs []ghJob) (string, string) {
	for _, j := range jobs {
		if strings.EqualFold(j.Conclusion, "failure") {
			step := findFirstFailedStep(j.Steps)

			return j.Name, step
		}
	}

	return "unknown", "unknown"
}

func findFirstFailedStep(steps []ghStep) string {
	for _, s := range steps {
		if strings.EqualFold(s.Conclusion, "failure") {
			return s.Name
		}
	}

	return "unknown"
}
```

### 3.2 Decompose `RunPeTriage` into Focused <= 15 Line Functions
Split `RunPeTriage` into discrete orchestrator steps:
1. `checkPeTriageHelp`: checks help flags.
2. `parsePeTriageFlags`: parses flagset.
3. `collectPeFailures`: iterates over repos within limit and collects failing runs.
4. `probeRepoRunFailure`: probes single repo with `ghJSON`.
5. `findMostRecentFailure`: selects the failure with the newest `CreatedAt` timestamp.
6. `renderWorstFailureSummary`: prints the triage bundle to stdout.

```go
// RunPeTriage implements gitmap pe triage (alias: pet).
func RunPeTriage(args []string) *appfault.AppError {
	if isHelp := checkPeTriageHelp(args); isHelp {
		return nil
	}

	cfg, parseErr := parsePeTriageFlags(args)
	if parseErr != nil {
		return parseErr
	}

	if !hasGh() {
		return appfault.NewValidationError("pe triage requires the GitHub CLI (gh) on PATH")
	}

	return executePeTriageFlow(cfg)
}

func executePeTriageFlow(cfg peTriageConfig) *appfault.AppError {
	repos, errRepos := listGitRepos(cfg.workdir)
	if errRepos != nil {
		return apperror.WrapSimple(errRepos, "pe.triage.listGitRepos")
	}

	if len(repos) == 0 {
		fmt.Printf("pe triage: no git repos found in %s\n", cfg.workdir)

		return nil
	}

	failures, probed := collectPeFailures(repos, cfg)
	if len(failures) == 0 {
		fmt.Printf("pe triage: probed %d repos, no failing latest runs. all green.\n", probed)

		return nil
	}

	worst := findMostRecentFailure(failures)
	renderWorstFailureSummary(worst, len(failures))

	return nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** All errors return `*appfault.AppError` directly.
- **Rule 2 (Function Limit):** Every helper is <= 15 lines (e.g. `findFirstFailedStep` is 10 lines; `findFirstFailedJobAndStep` is 12 lines).
- **Rule 3 (Zero Nested Loops & Ifs):** Nested loops are completely eradicated.
- **Rule 6 (Vertical Whitespace):** Proper blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [x] No nested loops exist in `pe_triage.go`.
- [x] `worstFailureDetail` delegates to `findFirstFailedJobAndStep` and `findFirstFailedStep`.
- [x] `RunPeTriage` decomposed into functions <= 15 lines.
- [x] Output formatting matches exact existing CLI behavior.
- [x] All file paths are strictly relative lowercase Git paths.

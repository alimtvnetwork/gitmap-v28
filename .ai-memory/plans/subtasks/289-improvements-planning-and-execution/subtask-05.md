# Subtask Plan 05: Replace Raw Error Discard with Structured `appfault.Wrap` in `health_checks.go`

- **Subtask ID:** `subtask-05`
- **Subtask Name:** `05-cmdhealth-structured-appfault-wrap`
- **Parent Spec:** `02-spec/21-app/289-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdhealth/health_checks.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdhealth/health_checks.go`, health check probes gather status on repositories, pipelines, fleet nodes, and database files:

1. In `probeRepoDirty`:
   `out, runErr := gitutil.RunGitInDir(repoPath, "status", "--porcelain")`
   If `runErr != nil`, the error is discarded without diagnostics, masking permission errors, locked index files, or missing git binaries.
2. In `collectRepoFacts`:
   `readErr` from `os.ReadDir(workdir)` and `lookErr` from `exec.LookPath(constants.GitBin)` are silently discarded with empty `return` statements.
3. In `probePipelineSlug`:
   `runErr` from `cmd.Run()` (`gh run list`) is silently dropped.
4. In `collectFleetFacts`:
   `openErr` from `store.OpenDefault()` and database failures from `db.GetSSHConnections` are flattened into an opaque string note without logging underlying database errors.

Under **Rule 1 (Error Return Sovereignty)**:
- Health fact collection must track underlying diagnostic errors using structured `appfault.Wrap(err, code, msg)` or record them into `healthFacts.diagErrors []error`.
- Silent blank returns must be eliminated; if a probe fails due to system error, the root cause must be preserved for verbose reporting and audit capture.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Enhance `healthFacts` Struct with Diagnostic Error Capture
- **File:** `cli/cmdhealth/health_checks.go`
- **Enhancement:**
  ```go
  type healthFacts struct {
  	repoTotal    int
  	repoDirty    int
  	dirtyNames   []string
  	pipelineRed  int
  	redNames     []string
  	pipelineNote string
  	fleetNodes   int
  	fleetNote    string
  	staleDBs     int
  	staleNames   []string
  	diagErrors   []error
  }
  ```

### 2.2 Refactor `probeRepoDirty` with Structured Wrapping
- **File:** `cli/cmdhealth/health_checks.go`
- **Implementation:**
  ```go
  func probeRepoDirty(repoPath, name string, facts *healthFacts) {
  	facts.repoTotal++
  	out, runErr := gitutil.RunGitInDir(repoPath, "status", "--porcelain")
  	if runErr != nil {
  		facts.diagErrors = append(facts.diagErrors, apperror.WrapSimple(runErr, "health.probeRepoDirty: "+name))
  		return
  	}

  	if strings.TrimSpace(string(out)) == "" {
  		return
  	}

  	facts.repoDirty++
  	appendDirtyName(facts, name)
  }
  ```

### 2.3 Refactor `collectRepoFacts` to Preserve Workspace Read Errors
- **File:** `cli/cmdhealth/health_checks.go`
- **Implementation:**
  ```go
  func collectRepoFacts(workdir string, facts *healthFacts) {
  	entries, readErr := os.ReadDir(workdir)
  	if readErr != nil {
  		facts.diagErrors = append(facts.diagErrors, apperror.WrapSimple(readErr, "health.collectRepoFacts: "+workdir))
  		return
  	}

  	if _, lookErr := exec.LookPath(constants.GitBin); lookErr != nil {
  		facts.diagErrors = append(facts.diagErrors, apperror.WrapSimple(lookErr, "health.gitBinNotFound"))
  		return
  	}

  	for _, e := range entries {
  		inspectRepoEntry(workdir, e, facts)
  	}
  }
  ```

### 2.4 Refactor `probePipelineSlug` to Record Subprocess Failures
- **File:** `cli/cmdhealth/health_checks.go`
- **Implementation:**
  - When `cmd.Run()` returns non-nil error, wrap with `apperror.WrapSimple(runErr, "health.probePipeline: "+slug)` and append to `facts.diagErrors`.

---

## 3. Strict Coding Guidelines Compliance Checklist

- [ ] **Error Sovereignty:** All discards replaced with structured wrapping via `apperror.WrapSimple`.
- [ ] **15-Line Limit:** Function bodies remain strictly <= 15 lines.
- [ ] **Zero Nested Ifs:** All branches flattened with early returns.
- [ ] **Positive Booleans:** Positive check logic used throughout.
- [ ] **Return Styling:** Standard vertical whitespace before returns.
- [ ] **Zero Git / Zero Build / Zero Test:** Strictly spec and plan authoring.

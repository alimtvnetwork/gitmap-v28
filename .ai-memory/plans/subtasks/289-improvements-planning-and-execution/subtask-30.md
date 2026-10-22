# Subtask 30: Graceful Handling for GitHub CLI Absence

> **Subtask ID:** Subtask-30  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdpe/pe_gh.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement proactive dependency checking in `ghJSON` in `cli/cmdpe/pe_gh.go` to gracefully handle environments where the GitHub CLI (`gh`) is absent.
2. Replace unhelpful OS binary missing execution errors with structured, user-friendly `*appfault.AppError` messages.
3. Validate that callers in `cli/cmdpe/` (`pe_triage.go`, `pe_rerun.go`) degrade safely when `!hasGh()`.
4. Ensure leaf error return sovereignty, function length limits (<= 15 lines), and zero nested conditionals.

---

## 2. Detailed Technical Plan

### 2.1 Pre-Flight Guard in `ghJSON`
- Update `ghJSON` in `cli/cmdpe/pe_gh.go`:
  ```go
  // ghJSON runs gh with --json output and decodes stdout as JSON into out.
  func ghJSON(out any, ghArgs ...string) *appfault.AppError {
  	if !hasGh() {
  		return apperror.NewSimple("github CLI (gh) is not installed or not found on PATH", "E_GH_NOT_FOUND")
  	}

  	cmd := exec.Command("gh", ghArgs...)
  	var stdout bytes.Buffer
  	var stderr bytes.Buffer
  	cmd.Stdout = &stdout
  	cmd.Stderr = &stderr
  	if runErr := cmd.Run(); runErr != nil {
  		return apperror.WrapSimple(runErr, "pe.gh: "+strings.TrimSpace(stderr.String()))
  	}

  	if decodeErr := json.Unmarshal(stdout.Bytes(), out); decodeErr != nil {
  		return apperror.WrapSimple(decodeErr, "pe.gh.decode")
  	}

  	return nil
  }
  ```

### 2.2 Safe Degradation in Triage Callers
- Ensure triage commands check `hasGh()` and return early with informative remediation:
  ```go
  if !hasGh() {
  	return appfault.NewValidationError("pe triage requires the GitHub CLI (gh) on PATH; install from https://cli.github.com/")
  }
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Returns `*appfault.AppError` directly. Does not call exit handlers.
- **Rule 2 (Function Limit):** `ghJSON` remains <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guard clauses for `!hasGh()`, `runErr != nil`, and `decodeErr != nil`.
- **Rule 6 (Vertical Whitespace):** Single blank line after each `}` closure and before return.

---

## 4. Verification Criteria
- [x] `ghJSON` returns structured `E_GH_NOT_FOUND` error when `gh` binary is not found.
- [x] No unhandled panic or OS binary execution error thrown.
- [x] Callers in `cmdpe` degrade gracefully without crash.
- [x] Zero compiler builds and zero tests executed during planning.

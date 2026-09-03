# 28 - Error Management & Architecture Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/03-error-manage/97-acceptance-criteria.md)

### AC-01: Structured Error Response
- **GIVEN** any CLI backend encounters an error during request processing  
- **WHEN** the error response is generated  
- **THEN** it contains: `Code` (numeric), `Message` (human-readable), `Details` (technical), and `Stack` (up to 40 frames)  
- **AND** the error code falls within the tool's assigned range per the Error Code Registry

**Edge Cases:**
- **GIVEN** the error originates from a third-party library **WHEN** the stack trace is captured **THEN** both the library frames and the application frames are included with clear delineation
- **GIVEN** the error code is not registered in the Error Code Registry **WHEN** it is returned **THEN** a fallback generic code within the tool's range is used and a warning is logged
- **GIVEN** the error `Details` field contains sensitive data (file paths, credentials) **WHEN** the response is generated **THEN** sensitive values are redacted before sending to the client

### AC-02: Frontend-Backend Verification Protocol
- **GIVEN** the frontend receives an error response from the backend  
- **WHEN** the error is displayed in the error modal  
- **THEN** the user can see the backend error code, the frontend component that triggered the request, and the timestamp  
- **AND** "Copy All" copies both frontend and backend context

**Edge Cases:**
- **GIVEN** the clipboard API is unavailable **WHEN** "Copy All" is clicked **THEN** a fallback textarea is shown with the content pre-selected
- **GIVEN** the backend returns an error with no `Code` field **WHEN** the frontend processes it **THEN** a synthetic code `GEN-1000` is assigned and a parsing warning is logged

### AC-03: Retrospective Document Structure
- **GIVEN** a production bug has been resolved  
- **WHEN** a retrospective document is created  
- **THEN** it contains: Root Cause, Timeline, Resolution Steps, Prevention Measures, and Related Error Codes

### AC-04: Verification Pattern Application
- **GIVEN** a developer implements a fix for a known error pattern  
- **WHEN** they consult the verification patterns documentation  
- **THEN** they find step-by-step verification instructions specific to the error category

### AC-05: Debugging Guide Coverage
- **GIVEN** a developer encounters a backend error  
- **WHEN** they follow the language-specific debugging guide  
- **THEN** they can identify common issues and each issue links to the relevant specification

### AC-06: Quick Resolution
- **GIVEN** a common error scenario  
- **WHEN** the developer consults the cheat sheet  
- **THEN** they find a 3-step resolution procedure: Identify → Diagnose → Fix  
- **AND** each step includes the exact command or code to run

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule EM-1 (Error Return Sovereignty & Zero Dual-Handling):** Any leaf function or helper with an `error` or `*apperror.AppError` return signature MUST return the error instance directly (`return err` or `return apperror.Wrap(...)`). Calling an exit handler (`cliexit.HandleError`) internally and then returning `nil` is strictly forbidden.
2. **Rule EM-2 (Strongly-Typed Exit Code Enums):** Raw integer literals (`1`, `2`, `0`) must never be passed to exit handlers. Exit codes must use strongly typed enums (`cliexit.ExitCodeType`, `constants.ExitGeneralError`, `constants.ExitUsageError`).
3. **Rule EM-3 (Specialized Exit Helpers & Parameter Reduction):** Frequently repeated exit invocations must be abstracted into specialized helpers (`cliexit.HandleValidationError`, `cliexit.HandleUsageError`, `cliexit.HandleGeneralError`).
4. **Rule EM-4 (Outer Caller Handling):** Only the top-level command entrypoint, router, or root dispatcher handles the exit transition. Intermediate functions propagate errors up the call stack.
5. **Rule EM-5 (Universal Envelope & Diagnostics):** Errors must preserve contextual metadata (`Op`, `Code`, `Details`, `Stack`) without swallowing underlying causes.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/vscode_cmd.go | 27 | `cliexit.HandleError(err, 1); return nil` | Return `err` directly to caller; eliminate dual handling | PENDING |
| V-02 | gitmap/cmd/vscode_cmd.go | 65 | `cliexit.HandleError(err, 1)` | Return error from dispatch function to root caller | PENDING |
| V-03 | gitmap/cmd/reinstall.go | 42 | `cliexit.HandleError(err, 1)` followed by dispatch | Return `err` directly on abort; do not proceed with reinstall | PENDING |
| V-04 | gitmap/cmd/macro_cmd.go | 101 | `cliexit.HandleError(nil, 1); return nil` | Return `apperror.NewValidationError` directly | PENDING |
| V-05 | gitmap/cmd/macro_cmd.go | 141 | `cliexit.HandleError(nil, 1); return nil` | Return `apperror.NewValidationError` directly | PENDING |
| V-06 | gitmap/cmd/macro_cmd.go | 145 | `cliexit.HandleError(nil, 1); return nil` | Return wrapped error directly | PENDING |
| V-07 | gitmap/cmd/macro_cmd.go | 186 | `cliexit.HandleError(nil, 1); return nil` | Return `apperror.NewValidationError` directly | PENDING |
| V-08 | gitmap/cmd/macro_cmd.go | 190 | `cliexit.HandleError(nil, 1); return nil` | Return wrapped error directly | PENDING |
| V-09 | gitmap/cmd/workflow_open_pr.go | 58 | `cliexit.HandleError(appErr, 2); return nil` | Return `appErr` directly to caller | PENDING |
| V-10 | gitmap/cmd/revertscript.go | 87 | `cliexit.HandleError(err, 1); return nil` | Return wrapped error directly | PENDING |
| V-11 | gitmap/cmd/sshcat.go | 45 | `cliexit.HandleError(appErr, 1); return nil` | Return `appErr` directly to caller | PENDING |
| V-12 | gitmap/cmd/sshgen.go | 32 | `cliexit.HandleError(appErr, 1); return nil` | Return `appErr` directly to caller | PENDING |
| V-13 | gitmap/cmd/sshgen.go | 49 | `cliexit.HandleError(appErr, 1); return nil` | Return `appErr` directly to caller | PENDING |
| V-14 | gitmap/cmd/sshgen.go | 73 | `cliexit.HandleError(appErr, 1); return nil` | Return `appErr` directly to caller | PENDING |
| V-15 | gitmap/cliexit/handle.go | 24 | Missing specialized exit helpers and `ExitCodeType` | Add `ExitCodeType` enum and `HandleValidationError`, `HandleUsageError`, `HandleGeneralError` | PENDING |
| V-16 | linter-scripts/check-error-management.py | 86 | Linter only checks `_ = err` | Add AST check for dual-handling (`cliexit.HandleError` followed by `return nil`) | PENDING |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/28-error-management/01-cliexit-specialized-helpers.md`)**:
   Implement `ExitCodeType` enum, specialized exit helpers (`HandleValidationError`, `HandleUsageError`, `HandleGeneralError`, `HandleSuccess`), and update `cliexit/handle.go`.
2. **Subtask 02 (`.lovable/plans/subtasks/28-error-management/02-leaf-error-returns-refactoring.md`)**:
   Refactor leaf command functions in `vscode_cmd.go`, `reinstall.go`, `macro_cmd.go`, `workflow_open_pr.go`, `revertscript.go`, `sshcat.go`, and `sshgen.go` to enforce error return sovereignty.
3. **Subtask 03 (`.lovable/plans/subtasks/28-error-management/03-linter-enhancement-and-ci-verification.md`)**:
   Enhance `linter-scripts/check-error-management.py` to continuously verify dual-handling prevention and run `03-ai-scripts/06-cicd-local-runner.py` for 100% green verification.

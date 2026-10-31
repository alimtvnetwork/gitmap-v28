# 370 Acceptance Criteria: AI Execution & Error Handling Instructions

> **Task:** 370-ai-execution-and-error-handling-instructions  
> **Subsystem:** PowerShell Command Interceptor & AI Diagnostic Feedback  
> **Date:** 31-Oct-26  
> **Status:** APPROVED  

---

## 1. Overview

This document specifies the measurable, testable acceptance criteria for the PowerShell path argument interception and AI diagnostic guidance engine.

---

## 2. Testable Acceptance Criteria

| ID | Category | Criterion Description | Verification Method | Pass Condition |
| :--- | :--- | :--- | :--- | :--- |
| `AC-370-01` | Pre-Flight Intercept | Directory path passed to `gitmap ps <dir>` must be intercepted before spawning PowerShell. | Execute `gitmap ps 02-spec/21-app` | Intercept occurs in `< 5ms`; no `pwsh.exe` spawned; exits code 2. |
| `AC-370-02` | Pre-Flight Intercept | Non-script document path (`.md`, `.json`, `.go`) must be intercepted before spawning PowerShell. | Execute `gitmap ps cli/cmd/root.go` | Diagnostic guidance displayed; suggests `gitmap cat cli/cmd/root.go`. |
| `AC-370-03` | Diagnostic Inquiry | Output must prominently display the reflective questions: "Why did you run this command?" and "What was your expectation?". | Inspect terminal stdout for `TC-PS-01` | Both diagnostic inquiry questions are present in output. |
| `AC-370-04` | Semantic Explanation | Output must explain PowerShell execution semantics (command mode vs expression mode, syntax parser). | Inspect terminal stdout for `TC-PS-01` | Explains why `02-spec` triggers ParserError due to `-` operator. |
| `AC-370-05` | Command Suggestions | Output must present replacement table with `gitmap ls`, `dir`, `gitmap agent subtask ls`, `gitmap cat`, `gitmap heatmap`. | Inspect suggestion table in terminal | All 5 canonical alternatives present in table. |
| `AC-370-06` | Learning Roadmap | Output or documentation must provide 2-3 step learning solutions for autonomous self-remediation. | Verify section 6 of architecture spec | 3 structured learning workflows documented with Mermaid diagrams. |
| `AC-370-07` | Pass-Through Scripts | Valid `.ps1` scripts and properly formatted PowerShell expressions must execute without interception. | Execute `gitmap ps "Get-Date"` | Successfully returns date output with exit code 0. |
| `AC-370-08` | Fallback Interceptor | If an argument bypasses pre-flight and fails in PowerShell with `ParserError` or `CommandNotFoundException`, the fallback interceptor catches it. | Execute invalid shell cmdlet | Raw Go stack trace suppressed; educational advice attached. |
| `AC-370-09` | Database Audit | Every intercepted invocation must record an audit entry in `CommandAudit` with error code `E_PS_PATH_INTERCEPTED`. | Query SQLite audit table after execution | Audit record present with matching timestamp and metadata. |
| `AC-370-10` | Relative Paths | All specification documentation and code references must strictly use relative paths without absolute drive letters. | Grep spec files for absolute drive prefixes | Zero absolute path occurrences found. |

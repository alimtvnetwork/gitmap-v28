---
plan: .lovable/plans/pending/02-gitmap-installer.md
domain: Cli
phase: Implement
target_files: ["gitmap/cmd/install_comp_096.go", "gitmap/cmd/install_comp_096_test.go"]
depends_on: ["Task 095"]
citations:
  app_spec: "spec/commands/01-gitmap-installer.md §Command Specifications"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  language_guideline: "spec/02-coding-guidelines/03-golang/00-overview.md"
  boolean_styling: "n/a - no boolean parsing"
  folder_naming: "spec/02-coding-guidelines/08-file-folder-naming/03-golang.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  error_codes: "spec/21-app/07-error-and-logging/01-error-code-allocation.md"
  logging_traces: "spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md"
  response_envelope: "spec/21-app/07-error-and-logging/03-response-envelope.md"
  golden_fixture: "n/a - no fixtures needed"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "n/a - no DB yet"
  ui_surface: "n/a"
  tests: "unit TestInstallComp096"
  ci_cd_guard: "linter-scripts/check-go-build"
  ambiguity: "n/a"
  issue_rca: "n/a"
---
# Task 096 — Scaffold Installer Component 096

## 1. Learn
- Read `spec/commands/01-gitmap-installer.md` to understand the goal.
- Read `spec/02-coding-guidelines/00-canonical-size-tier.md` for sizing rules.
- Read `spec/03-error-manage/02-error-architecture/00-overview.md` for `apperror` usage.

## 2. Goal
Implement `InstallComp096` which processes a subset of the installer specifications logic. This task will scaffold the unit, add the interface and strict typing for installer commands, and ensure it builds correctly.

## 3. Inputs and Contracts
Input: `Input096` struct
Output: `Output096` struct
Error: standard `apperror.Result` envelope.
Data Uniqueness: `28be84291480` (for testing).

## 4. Execute
1. Create `gitmap/cmd/install_comp_096.go`.
2. Define `type Input096 struct { Data string }` and `type Output096 struct { Result string }`.
3. Define `func HandleInstallComp096(in Input096) (Output096, error)`.
4. Create `gitmap/cmd/install_comp_096_test.go`.
5. Write unit test `TestInstallComp096` covering success and failure.

## 5. Constraints
- Functions under 50 lines (canonical size).
- Errors wrapped in `apperror`.
- No negatives in booleans.

## 6. Verify
```bash
go test ./cmd/... -run TestInstallComp096
```
Expected output: `PASS` and `ok`.

## 7. Done When
- [ ] `HandleInstallComp096` exists and passes unit tests.
- [ ] No magic strings used outside test constants.
- [ ] Function is correctly typed.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

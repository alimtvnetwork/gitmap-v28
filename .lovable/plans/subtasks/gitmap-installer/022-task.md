---
plan: .lovable/plans/pending/02-gitmap-installer.md
domain: Plugin
phase: Implement
target_files: ["gitmap/installer/create.go", "gitmap/installer/create_test.go"]
depends_on: ["Task 021"]
citations:
  app_spec: ".lovable/spec/commands/01-gitmap-installer.md §Core Requirements"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  language_guideline: "spec/02-coding-guidelines/03-golang/00-overview.md"
  boolean_styling: "spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md"
  folder_naming: "spec/02-coding-guidelines/08-file-folder-naming/03-golang.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  error_codes: "spec/21-app/07-error-and-logging/01-error-code-allocation.md"
  logging_traces: "spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md"
  response_envelope: "spec/21-app/07-error-and-logging/03-response-envelope.md"
  golden_fixture: "n/a — no wire format in this step"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/04-database-conventions/01-sqlite-schema.md"
  ui_surface: "n/a — no direct UI in this step"
  tests: "unit TestManagerCreate"
  ci_cd_guard: "linter-scripts/check-go-build"
  ambiguity: "n/a — no ambiguity filed"
  issue_rca: "n/a — not a bugfix"
---
# Task 022 — Manager Create Logic

## 1. Learn
- Read `.lovable/spec/commands/01-gitmap-installer.md` to understand the overarching requirement for (m *Manager) Create(name, desc).
- Read `spec/02-coding-guidelines/00-canonical-size-tier.md` to ensure `gitmap/installer/create.go` remains concisely sized.
- Review `spec/03-error-manage/02-error-architecture/00-overview.md` for proper `apperror` context wrapping.
- Inspect `gitmap/installer/create.go` dependencies to see how (m *Manager) Create(name, desc) interacts with its callers.

## 2. Goal
The objective is to implement `(m *Manager) Create(name, desc)` natively in `gitmap/installer/create.go`. This explicitly unblocks downstream operations dependent on `Manager Create Logic` in the Plugin domain. No other files should be manipulated.

## 3. Inputs and Contracts
- Exported Symbols: `(m *Manager) Create(name, desc)`
- Package: `installer`
- Error wrapping MUST use the `E_INSTALLER_*` code family.

## 4. Execute
1. Open `gitmap/installer/create.go`.
2. Implement the required structure, type, or function for `(m *Manager) Create(name, desc)`.
3. Write unit tests for success and failure boundaries in `gitmap/installer/create_test.go`.
4. Ensure no cross-domain pollution.

## 5. Constraints
- Must adhere strictly to `spec/02-coding-guidelines/00-canonical-size-tier.md` (keep logic segmented under 60 lines).
- Error wrapping must include stack traces.

## 6. Verify
```bash
go test ./installer -run TestManagerCreate
```
Expected output: The test suite passes cleanly with no panics.

## 7. Done When
- [ ] `(m *Manager) Create(name, desc)` is successfully mapped and tested in `gitmap/installer/create.go`.
- [ ] All CI and `go test` commands exit zero.
- [ ] No hardcoded or dummy assumptions are left in the code.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

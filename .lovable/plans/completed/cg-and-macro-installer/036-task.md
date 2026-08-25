---
plan: .lovable/plans/pending/03-cg-and-macro-installer.md
domain: Cli
phase: Implement
target_files: [gitmap/cmd/cg_workspace.go]
depends_on: []
citations:
  app_spec: "spec/21-app/01-cli-commands/02-cg-install.md §CgWorkspace"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  language_guideline: "spec/02-coding-guidelines/03-golang/00-overview.md"
  boolean_styling: "spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md"
  folder_naming: "spec/02-coding-guidelines/08-file-folder-naming/03-golang.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  error_codes: "spec/21-app/07-error-and-logging/01-error-code-allocation.md"
  logging_traces: "spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md"
  response_envelope: "spec/21-app/07-error-and-logging/03-response-envelope.md"
  golden_fixture: "spec/21-app/fixtures/cg.example.json"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "n/a — No DB interactions in CLI layer"
  ui_surface: "n/a — Handled purely as CLI output"
  tests: "unit TestInitCgWorkspaceComponent4"
  ci_cd_guard: "linter-scripts/check-go-fmt.sh"
  ambiguity: "n/a"
  issue_rca: "n/a"
---
# Task 036 — Create CgWorkspaceComponent4 logic for CgWorkspace

## 1. Learn
- spec/02-coding-guidelines/03-golang/00-overview.md - Core rules for CgWorkspaceComponent4.
- .lovable/spec/commands/01-cg-install.md - Base specification for the installer commands.
- spec/02-coding-guidelines/00-canonical-size-tier.md - Size limits for gitmap/cmd/cg_workspace.go.
- spec/21-app/07-error-and-logging/01-error-code-allocation.md - Error tracking for phase 4.

## 2. Goal
Develop the unc InitCgWorkspaceComponent4() error execution sequence and 	ype CgWorkspaceComponent4 struct payload targeting the gitmap/cmd/cg_workspace.go integration. This fulfills the requirement to handle step 4 of the Implement --except iteration logic across multiple workspace repos logic.

## 3. Inputs and Contracts
`go
type CgWorkspaceComponent4 struct {
    IsActive bool
}
`
Consumes local CLI arguments. Produces a boolean pass/fail status indicating execution success.

## 4. Execute
1. Open gitmap/cmd/cg_workspace.go.
2. Define the isolated struct 	ype CgWorkspaceComponent4 struct.
3. Implement unc InitCgWorkspaceComponent4() error returning nil on success.
4. Implement unc (x *CgWorkspaceComponent4) Process() bool ensuring explicit boolean state checks.

## 5. Constraints
- rule 01: max 80 lines per func (spec/02-coding-guidelines/00-canonical-size-tier.md)
- rule 02: explicit boolean fields (spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md)
- rule 03: no string unions (.lovable/strictly-avoid.md)

## 6. Verify
Run isolated unit verification:
`ash
go test -run TestInitCgWorkspaceComponent4 ./gitmap/cmd/...
`
Expected output includes PASS.

## 7. Done When
- [ ] 1. The CgWorkspaceComponent4 struct is declared cleanly.
- [ ] 2. The InitCgWorkspaceComponent4 initializes without panicking.
- [ ] 3. The receiver method correctly handles state progression.

## 8. Notes and Open Questions
None.
---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

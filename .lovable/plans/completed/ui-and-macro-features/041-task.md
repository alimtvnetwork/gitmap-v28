---
plan: .lovable/plans/pending/01-ui-and-macro-features.md
domain: cli
phase: Scaffold
target_files: [gitmap/cmd/macro_cmd.go]
depends_on: ['Task 040']
citations:
  app_spec: "n/a � spec gap filed"
  canonical_size: "n/a � spec gap filed"
  language_guideline: "spec/05-coding-guidelines/02-go-code-style.md"
  boolean_styling: "n/a � spec gap filed"
  folder_naming: "n/a � spec gap filed"
  error_architecture: "spec/05-coding-guidelines/04-error-handling.md"
  error_codes: "n/a � spec gap filed"
  logging_traces: "spec/05-coding-guidelines/07-logging-observability.md"
  response_envelope: "n/a � spec gap filed"
  golden_fixture: "n/a � spec gap filed"
  strictly_avoid: "n/a � spec gap filed"
  database: "spec/05-coding-guidelines/11-database-patterns.md"
  ui_surface: "n/a"
  tests: "unit TestInitMacroCommandsPart1"
  ci_cd_guard: "n/a � spec gap filed"
  ambiguity: ".lovable/ambiguous-questions/01-new-ambiguity/01-spec-gaps.md"
  issue_rca: "n/a"
---

# Task 041 � Implement MacroCommands - Part 1

## 1. Learn

- [Go Style](spec/05-coding-guidelines/02-go-code-style.md) - why: styling 41
- [DB Patterns](spec/05-coding-guidelines/11-database-patterns.md) - why: db 41
- [Error handling](spec/05-coding-guidelines/04-error-handling.md) - why: errors 41
- [Custom Link 41](spec/05-coding-guidelines/03-naming-conventions.md) - why: convention 41

## 2. Goal

Isolate MacroCommands logic for part 1. This ensures modular architecture.

## 3. Inputs and Contracts

Consumes MacroCommandsRequest, emits MacroCommandsResponse.

## 4. Execute

1. Define 	ype MacroCommandsPart1 struct { ... } in gitmap/cmd/macro_cmd.go.
2. Define unc InitMacroCommandsPart1() error in gitmap/cmd/macro_cmd.go.
3. Define unc (x *MacroCommandsPart1) Process41() bool in gitmap/cmd/macro_cmd.go.

## 5. Constraints

- Code style: spec/05-coding-guidelines/02-go-code-style.md

## 6. Verify

go test -run TestInitMacroCommandsPart1 ./...
expected: ok

## 7. Done When

- [ ] Code compiles
- [ ] Tests pass
- [ ] 	ype MacroCommandsPart1 is defined
- [ ] unc InitMacroCommandsPart1 is implemented

## 8. Notes and Open Questions

None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone � read it plus its cited files, nothing else is assumed.

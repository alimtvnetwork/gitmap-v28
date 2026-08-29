---
plan: .lovable/plans/pending/01-ui-and-macro-features.md
domain: cli
phase: Scaffold
target_files: [gitmap/cmd/macro_tui.go]
depends_on: ['Task 036']
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
  tests: "unit TestInitMacroInteractivePart1"
  ci_cd_guard: "n/a � spec gap filed"
  ambiguity: ".lovable/ambiguous-questions/01-new-ambiguity/01-spec-gaps.md"
  issue_rca: "n/a"
---

# Task 037 � Implement MacroInteractive - Part 1

## 1. Learn

- [Go Style](spec/05-coding-guidelines/02-go-code-style.md) - why: styling 37
- [DB Patterns](spec/05-coding-guidelines/11-database-patterns.md) - why: db 37
- [Error handling](spec/05-coding-guidelines/04-error-handling.md) - why: errors 37
- [Custom Link 37](spec/05-coding-guidelines/03-naming-conventions.md) - why: convention 37

## 2. Goal

Isolate MacroInteractive logic for part 1. This ensures modular architecture.

## 3. Inputs and Contracts

Consumes MacroInteractiveRequest, emits MacroInteractiveResponse.

## 4. Execute

1. Define 	ype MacroInteractivePart1 struct { ... } in gitmap/cmd/macro_tui.go.
2. Define unc InitMacroInteractivePart1() error in gitmap/cmd/macro_tui.go.
3. Define unc (x *MacroInteractivePart1) Process37() bool in gitmap/cmd/macro_tui.go.

## 5. Constraints

- Code style: spec/05-coding-guidelines/02-go-code-style.md

## 6. Verify

go test -run TestInitMacroInteractivePart1 ./...
expected: ok

## 7. Done When

- [ ] Code compiles
- [ ] Tests pass
- [ ] 	ype MacroInteractivePart1 is defined
- [ ] unc InitMacroInteractivePart1 is implemented

## 8. Notes and Open Questions

None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone � read it plus its cited files, nothing else is assumed.

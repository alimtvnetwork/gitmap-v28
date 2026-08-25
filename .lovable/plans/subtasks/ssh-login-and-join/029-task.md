---
plan: .lovable/plans/pending/01-ssh-login-and-join.md
domain: Cli
phase: Wire+Test
target_files: ["gitmap/cmd/ssh_e2e_test.go"]
depends_on: [Task 028]
citations:
  app_spec: "spec/19-ssh-executor/01-spec.md §Section"
  canonical_size: "spec/05-coding-guidelines/01-code-quality-improvement.md"
  language_guideline: "spec/05-coding-guidelines/02-go-code-style.md"
  boolean_styling: "spec/05-coding-guidelines/03-naming-conventions.md"
  folder_naming: "spec/05-coding-guidelines/05-file-project-structure.md"
  error_architecture: "spec/05-coding-guidelines/04-error-handling.md"
  error_codes: "spec/05-coding-guidelines/04-error-handling.md"
  logging_traces: "spec/05-coding-guidelines/07-logging-observability.md"
  response_envelope: "spec/05-coding-guidelines/10-api-design.md"
  golden_fixture: "spec/08-json-schemas/ssh-list.schema.json"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/05-coding-guidelines/11-database-patterns.md"
  ui_surface: "n/a — cli tool"
  tests: "unit TestTestDispatchSSHLogin"
  ci_cd_guard: ".github/workflows/ci.yml"
  ambiguity: "n/a"
  issue_rca: "n/a"
---
# Task 029 — E2E Test SSH Login Dispatcher

## 1. Learn
- [SSH Commands](file:///d:/work/gitmap/.lovable/spec/commands/01-ssh-commands.md) — Why: Defines required behaviour.
- [App Error Docs](file:///d:/work/gitmap/spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [gitmap/cmd/ssh_e2e_test.go](file:///d:/work/gitmap/gitmap/cmd/ssh_e2e_test.go) — Why: Target file.

## 2. Goal
Deliver the Wire+Test step for `TestDispatchSSHLogin` to support the E2E Test SSH Login Dispatcher feature. This is isolated logic for the SSH/IP subdomains.

## 3. Inputs and Contracts
- Types: `string`, `context.Context`
- Outputs: `error`
- Codes: `E_INTERNAL_ERROR`
- Signature:
  ```go
  func TestDispatchSSHLogin(t *testing.T, args []string, mockCtx context.Context)
  ```

## 4. Execute
1. Call `dispatchSSH` with mock args `[]string{"login", "a@b"}`.
2. Assert correct handlers are invoked without actually spawning SSH.

## 5. Constraints
- **Canonical Size**: spec/05-coding-guidelines/01-code-quality-improvement.md.
- **Error Types**: Must use `apperror`.
- **No Globals**: .lovable/strictly-avoid.md.

## 6. Verify
```bash
go test ./... -v -run TestDispatchSSHLogin
```
Expected output:
```text
PASS
```

## 7. Done When
- [ ] 1. `TestDispatchSSHLogin` is fully functional.
- [ ] 2. Tests pass successfully.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

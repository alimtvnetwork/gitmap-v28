---
plan: .lovable/plans/pending/01-ssh-login-and-join.md
domain: Cli
phase: Scaffold
Status: completed
target_files: ["gitmap/store/models_ssh_hist.go"]
depends_on: [None]
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
  tests: "unit TestSSHHistory"
  ci_cd_guard: ".github/workflows/ci.yml"
  ambiguity: "n/a"
  issue_rca: "n/a"
---

# Task 007 — Define Go struct for SSHHistory

## 1. Learn

- [SSH Commands](.lovable/spec/commands/01-ssh-commands.md) — Why: Defines required behavior.
- [App Error Docs](spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [gitmap/store/models_ssh_hist.go](gitmap/store/models_ssh_hist.go) — Why: Target file.

## 2. Goal

Deliver the Scaffold step for `SSHHistory` to support the Define Go struct for SSHHistory feature. This is isolated logic for the SSH/IP subdomains.

## 3. Inputs and Contracts

- Types: `string`, `context.Context`
- Outputs: `error`
- Codes: `E_INTERNAL_ERROR`
- Signature:
  ```go
  type SSHHistory struct { ID string; HostIP string; JoinedAt time.Time }
  ```

## 4. Execute

1. Declare `SSHHistory` with `ID`, `HostIP`, `JoinedAt`, `User`.
2. Ensure proper mapping for time formats.

## 5. Constraints

- **Canonical Size**: spec/05-coding-guidelines/01-code-quality-improvement.md.
- **Error Types**: Must use `apperror`.
- **No Globals**: .lovable/strictly-avoid.md.

## 6. Verify

```bash
go test ./... -v -run SSHHistory
```
Expected output:
```text
PASS
```

## 7. Done When

- [ ] 1. `SSHHistory` is fully functional.
- [ ] 2. Tests pass successfully.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions

None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

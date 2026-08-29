---
plan: .lovable/plans/pending/01-ssh-login-and-join.md
domain: Cli
phase: Implement
Status: completed
target_files: ["gitmap/store/ssh_repo.go"]
depends_on: [Task 001, Task 002]
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
  tests: "unit TestDeleteHostByIP"
  ci_cd_guard: ".github/workflows/ci.yml"
  ambiguity: "n/a"
  issue_rca: "n/a"
---

# Task 005 — Implement SSHHost Delete

## 1. Learn

- [SSH Commands](.lovable/spec/commands/01-ssh-commands.md) — Why: Defines required behavior.
- [App Error Docs](spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [gitmap/store/ssh_repo.go](gitmap/store/ssh_repo.go) — Why: Target file.

## 2. Goal

Deliver the Implement step for `DeleteHostByIP` to support the Implement SSHHost Delete feature. This is isolated logic for the SSH/IP subdomains.

## 3. Inputs and Contracts

- Types: `string`, `context.Context`
- Outputs: `error`
- Codes: `E_INTERNAL_ERROR`
- Signature:
  ```go
  func DeleteHostByIP(ctx context.Context, ip string, db *sql.DB) error
  ```

## 4. Execute

1. Write `DeleteHostByIP(ctx context.Context, ip string) error`.
2. Execute `DELETE FROM ssh_hosts WHERE ip = ?`.
3. Return nil if 0 rows affected (idempotency).

## 5. Constraints

- **Canonical Size**: spec/05-coding-guidelines/01-code-quality-improvement.md.
- **Error Types**: Must use `apperror`.
- **No Globals**: .lovable/strictly-avoid.md.

## 6. Verify

```bash
go test ./... -v -run DeleteHostByIP
```
Expected output:
```text
PASS
```

## 7. Done When

- [ ] 1. `DeleteHostByIP` is fully functional.
- [ ] 2. Tests pass successfully.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions

None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.

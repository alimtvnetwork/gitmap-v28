---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/installer/execute_order.go"]
depends_on: ["Task 003" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 004 — Execution Order Engine

## 1. Goal
Execute unix-first, os-first, and fallback ordering in `installer/execute_order.go`.

## 2. Inputs and Contracts
- Package: derived from `installer/execute_order.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestExecutionOrderEngine
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/commitin/help.go"]
depends_on: ["Task 026" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---

# Task 027 — Commit-In CLI Help Integration

## 1. Goal

Deep interactive help display for commit-in and commit-write in `cmd/commitin/help.go`.

## 2. Inputs and Contracts

- Package: derived from `cmd/commitin/help.go`
- Strict error wrapping with `apperror`.

## 3. Verify

```bash
go test ./... -run TestCommit-InCLIHelpIntegration
```

## 4. Done When

- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

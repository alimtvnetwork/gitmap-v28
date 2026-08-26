---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/ssh_multi_parser.go"]
depends_on: ["Task 013" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 014 — SSH Multi-IP Parser

## 1. Goal
Parse comma-separated or space-separated IP lists in `cmd/ssh_multi_parser.go`.

## 2. Inputs and Contracts
- Package: derived from `cmd/ssh_multi_parser.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestSSHMulti-IPParser
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/sshjoin_broadcast.go"]
depends_on: ["Task 016" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 017 — SSH Broadcast Command Core

## 1. Goal
Execute remote commands across multiple machines in `cmd/sshjoin_broadcast.go`.

## 2. Inputs and Contracts
- Package: derived from `cmd/sshjoin_broadcast.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestSSHBroadcastCommandCore
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

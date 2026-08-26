---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/sshjoin_profile_sync.go"]
depends_on: ["Task 017" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 018 — SSH ZSH Profile Sync

## 1. Goal
Deploy ZSH and environment configs to remote nodes in `cmd/sshjoin_profile_sync.go`.

## 2. Inputs and Contracts
- Package: derived from `cmd/sshjoin_profile_sync.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestSSHZSHProfileSync
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

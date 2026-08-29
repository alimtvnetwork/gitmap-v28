---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/os_full_upgrade.go"]
depends_on: ["Task 020" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---

# Task 021 — Universal OS Full Upgrade

## 1. Goal

Engine for OS version full upgrades in `cmd/os_full_upgrade.go`.

## 2. Inputs and Contracts

- Package: derived from `cmd/os_full_upgrade.go`
- Strict error wrapping with `apperror`.

## 3. Verify

```bash
go test ./... -run TestUniversalOSFullUpgrade
```

## 4. Done When

- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

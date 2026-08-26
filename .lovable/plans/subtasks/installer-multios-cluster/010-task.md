---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/installer_smart_install.go"]
depends_on: ["Task 009" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 010 — Smart Auto-Detect Install CLI

## 1. Goal
Auto-detect active host OS and run matching installer block in `cmd/installer_smart_install.go`.

## 2. Inputs and Contracts
- Package: derived from `cmd/installer_smart_install.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestSmartAuto-DetectInstallCLI
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

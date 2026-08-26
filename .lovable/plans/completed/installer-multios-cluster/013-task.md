---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/installer_export_git.go"]
depends_on: ["Task 012" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 013 — Git Export CLI

## 1. Goal
CLI for export-git and export-all-git in `cmd/installer_export_git.go`.

## 2. Inputs and Contracts
- Package: derived from `cmd/installer_export_git.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestGitExportCLI
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

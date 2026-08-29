---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/cmd/installer_os_cmds.go"]
depends_on: ["Task 007" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---

# Task 008 — Dedicated OS Install CLI

## 1. Goal

CLI commands install-ubuntu, install-arch, install-centos, install-debian, install-fedora, install-mac, install-unix in `cmd/installer_os_cmds.go`.

## 2. Inputs and Contracts

- Package: derived from `cmd/installer_os_cmds.go`
- Strict error wrapping with `apperror`.

## 3. Verify

```bash
go test ./... -run TestDedicatedOSInstallCLI
```

## 4. Done When

- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

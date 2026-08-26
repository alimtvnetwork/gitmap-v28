---
plan: .lovable/plans/pending/03-installer-multios-cluster.md
domain: Plugin
phase: Implement
target_files: ["gitmap/osutil/detector.go"]
depends_on: ["Task 018" if int(num) > 1 else "None"]
citations:
  app_spec: ".lovable/spec/commands/02-installer-multios-and-cluster.md"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  strictly_avoid: ".lovable/strictly-avoid.md"
---
# Task 019 — Universal OS Detector

## 1. Goal
Detailed detection for Windows, Ubuntu, Debian, CentOS, Fedora, Arch, macOS in `osutil/detector.go`.

## 2. Inputs and Contracts
- Package: derived from `osutil/detector.go`
- Strict error wrapping with `apperror`.

## 3. Verify
```bash
go test ./... -run TestUniversalOSDetector
```

## 4. Done When
- [ ] Task logic implemented cleanly.
- [ ] Unit tests pass with zero compilation errors.

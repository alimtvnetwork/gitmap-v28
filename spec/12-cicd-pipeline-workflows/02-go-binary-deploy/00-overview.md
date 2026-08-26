# Go Binary Deploy — Overview

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Pipeline specifications for building, testing, and releasing cross-compiled Go binaries. These pipelines handle static compilation for 6 platform/architecture targets, SHA-based build deduplication, platform-specific install scripts, and code signing.

---

## Key Characteristics

| Property | Value |
|----------|-------|
| Language | Go |
| Build Mode | Static linking (`CGO_ENABLED=0`) |
| Targets | 6 platforms (windows/linux/darwin × amd64/arm64) |
| Compression | `.zip` (Windows), `.tar.gz` (Linux/macOS) |
| Distribution | GitHub Releases + install scripts |
| Version Embedding | `-ldflags -X` at compile time |

---

## Pipeline Architecture

```
CI Pipeline:
  sha-check → [lint, vulncheck] → test (matrix: N suites) → test-summary → build (matrix: 6 targets) → build-summary

Release Pipeline:
  setup (version) → build all binaries → compress → checksums → install scripts → changelog → GitHub Release
```

The CI pipeline uses a **SHA-based passthrough gate** to skip redundant validation of already-tested commits.

---

## Feature Inventory

| # | File | Description | Status |
|---|------|-------------|--------|
| 01 | [01-ci-pipeline.md](./01-ci-pipeline.md) | CI pipeline: SHA dedup, lint, vulncheck, test matrix, cross-compile | ✅ Active |
| 02 | [02-release-pipeline.md](./02-release-pipeline.md) | Release pipeline: binary build, compression, install scripts, GitHub Release | ✅ Active |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Shared Conventions | `../01-shared-conventions.md` |
| GitHub Release Standard | `../02-github-release-standard.md` |
| Vulnerability Scanning | `../03-vulnerability-scanning.md` |
| Install Script Generation | `../04-install-script-generation.md` |
| Code Signing | `../05-code-signing.md` |
| Self-Update Mechanism | `../06-self-update-mechanism.md` |
| Release Body & Changelog | `../07-release-body-and-changelog.md` |
| Self-Update Full Specs | `../../14-update/00-overview.md` |

---

*Overview — updated: 2026-04-10*

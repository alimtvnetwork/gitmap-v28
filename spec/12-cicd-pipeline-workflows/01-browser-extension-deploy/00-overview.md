# Browser Extension Deploy — Overview

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Pipeline specifications for building, testing, and releasing browser extensions (Chrome/Chromium) built with Node.js and a package manager (pnpm/npm). These pipelines handle multi-component dependency graphs where an SDK must be built before dependent modules, and all modules must be assembled into a final extension package.

---

## Key Characteristics

| Property | Value |
|----------|-------|
| Language | TypeScript / JavaScript |
| Package Manager | pnpm (or npm) |
| Build Tool | Vite, Webpack, or similar bundler |
| Output | `.zip` archive of extension `dist/` contents |
| Distribution | GitHub Releases, Chrome Web Store (manual upload) |

---

## Pipeline Architecture

```
CI Pipeline:
  setup (lint + test) → build-sdk → [build-module-A, build-module-B, build-module-C] → build-extension

Release Pipeline:
  setup (lint + test + version) → build-sdk → [build-modules...] → build-extension → package → release
```

The build graph has a **diamond dependency**: the SDK is built first, then multiple standalone modules build in parallel (each downloading the SDK artifact), then the final extension build assembles everything.

---

## Feature Inventory

| # | File | Description | Status |
|---|------|-------------|--------|
| 01 | [01-ci-pipeline.md](./01-ci-pipeline.md) | CI pipeline: lint, test, dependency-graph builds | ✅ Active |
| 02 | [02-release-pipeline.md](./02-release-pipeline.md) | Release pipeline: version, build, package, GitHub Release | ✅ Active |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Shared Conventions | `../01-shared-conventions.md` |
| GitHub Release Standard | `../02-github-release-standard.md` |
| Vulnerability Scanning | `../03-vulnerability-scanning.md` |

---

*Overview — updated: 2026-04-09*

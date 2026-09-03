# Centralized Version.json Architecture

- **Date**: 2026-08-09
- **Topic**: Release procedure, Gitmap CI/CD decoupling.

## The Problem

Previously, AI agents were instructed to run massive PowerShell array string replacements across multiple files (like `constants.go`, `latest.json`, `index.ts`, `plan.md`) to bump the version. This was incredibly brittle and caused the `v6.87.1` regression where the source code `constants.go` failed to sync with the Git tag, leading to broken CI smoke tests.
The AI agents should never touch core Golang source files during a standard release bump.

## The Solution

1. **Central `version.json`**: The absolute source of truth for the project version is now strictly `version.json` at the root of the repository.
2. **Go `var Version`**: Inside `gitmap/constants/constants.go`, the `Version` field was converted from `const` to `var Version = "0.0.0-dev"`. This enables the GitHub Action to natively inject the version using `-ldflags` during compilation.
3. **CI/CD Integration**: The `release.yml` GitHub Action now parses the version directly from `version.json` using `jq` rather than deriving it from the Git tag.
4. **AI Rules**: All specifications in `spec/` and `plan.md` have been rewritten. AI agents must **ONLY** bump `version.json` and `changelog.md` during a release.

## Enforcement

Do NOT attempt to use PowerShell to mutate `constants.go` or `.gitmap/release/latest.json`. Just bump `version.json`, update `changelog.md` and `package.json`, commit, and push. The pipeline will automatically compile the binary with the injected version.

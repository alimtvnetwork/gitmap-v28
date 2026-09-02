---
name: release-and-versioning
description: >-
  Autonomously manage repository versioning, changelog synchronization, and release assets following
  the Single Source of Truth architecture and CI release pipeline protection.
---

# Release Architecture & Versioning Management Skill

Autonomously execute and audit version bumps, changelog synchronization, and release assets adhering to `.lovable/strictly-avoid.md`, `spec/17-consolidated-guidelines/17-self-update-app-update.md`, and `.lovable/memory/release-architecture-map.md`.

## Core Checkpoints & Mandatory Invariants

1. **Single Source of Truth (SSoT) Versioning:**
   - `version.json` at root is the sole authoritative source of truth: `{"version": "X.Y.Z"}`.
   - Go binaries inject the version dynamically via build flags: `-ldflags "-X ...Version=..."`.
   - Never hand-edit version variables across Go code without synchronizing `version.json`.

2. **Synchronized Version Bump Checklist:**
   - Update root `version.json`.
   - Add new `## [vX.Y.Z]` release section at the top of `changelog.md`.
   - Repin active version mentions in root `readme.md`.
   - Update `package.json` version field.
   - Verify sync with `bunx vitest run src/test/version-sync.test.ts` and `.github/scripts/check-changelog-version-sync.sh`.

3. **Release Pipeline Protection (CODE RED):**
   - NEVER edit `.github/workflows/release.yml` or `.github/scripts/smoke-installer.*`.
   - Missing or failing release uploads are ALWAYS due to upstream build/test gating; fix the application code, never the release workflow.
   - NEVER manually create, modify, or delete files within `.gitmap/release/` or `.gitmap/release-assets/`.

4. **Release Install-Snippet Gating:**
   - Non-gitmap repositories must never have gitmap installer snippets appended to their release bodies.
   - Always guard `AppendPinnedInstallSnippet` with `ShouldPrintInstallHint(getRemoteURL())`.

5. **Legacy Reference Annotation:**
   - Any historical mentions of legacy version names in documentation or changelogs MUST be annotated with `<!-- gitmap-legacy-ref-allow -->` to pass automated linter scans.

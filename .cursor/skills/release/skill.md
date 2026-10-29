---
name: release
description: Unified release management orchestrator governing major, minor, patch version bumps, release notes generation, changelog synchronization, and Git tag release ceremonies.
---

# Unified Release Management & Version Bump Orchestrator (`release`)

> **Invocation:** `@[release](file;.cursor/skills/release)`

## Overview

`release` is the canonical skill governing all version increments and release deployments across the meta-repository. It consolidates and replaces all individual bump wrappers (`major-bump`, `minor-bump`, `patch-bump`, `version-bump-docs`, `release-now`).

Always execute with strictly relative Git paths and zero CI artifact uploads.

---

## 1. Version Increment Types & Semantic Rules

| Release Type | Source Specification / Prompt | Version Arithmetic | When to Use |
| :--- | :--- | :--- | :--- |
| **Minor Bump (Default)** | `01-prompts/17-release-management/02-minor-bump.md` | `X.Y.Z -> X.(Y+1).0` | Feature additions, backward-compatible improvements, new skills, coding guideline upgrades. Patch resets to 0. |
| **Patch Bump** | `01-prompts/17-release-management/03-patch-bump.md` | `X.Y.Z -> X.Y.(Z+1)` | Urgent bug fixes, hotfixes, typo corrections, minor lint remediations. |
| **Major Bump** | `01-prompts/17-release-management/01-major-bump.md` | `X.Y.Z -> (X+1).0.0` | Breaking API changes, non-backward-compatible schema overhauls. Minor and patch reset to 0. |
| **Docs Sync** | `01-prompts/17-release-management/05-version-bump-docs.md` | `vX.Y.Z` documentation parity | Syncing version strings across README, documentation, and spec trees without incrementing SemVer. |

---

## 2. Non-Negotiable Release Invariants

1. **Mandatory Pre-Flight Pull:** Always run `gitmap pull` (fallback `git pull`) to verify that the local tracking branch is fully synchronized with remote upstream before making version modifications.
2. **Dedicated Release Commit Verb (`gitmap cpr`):**
   * Releases MUST be committed using the sovereign release verb:
     ```bash
     gitmap cpr "vX.Y.Z - <concise-release-summary>"
     ```
   * **TOTAL BAN ON CPF:** NEVER use `gitmap cpf` for releases (it attaches a conflicting `Feature:` prefix).
   * **TOTAL BAN ON NESTED PREFIXES:** NEVER write `chore(release):` or `Release:` inside the message argument (GitMap prepends `Release:` automatically).
   * **TOTAL BAN ON BRACKETS & COLONS:** Format message with a hyphen: `"vX.Y.Z - <summary>"`.
3. **Strict Relative Git Paths:** All release notes and changelog markdown links MUST use strictly relative paths (`.ai-memory/release/release-notes-vX.Y.Z.md`). TOTAL BAN on `file:///` URIs and absolute filesystem paths.
4. **Git Tagging & Upstream Completion:**
   * Create an annotated or lightweight tag matching the exact version: `git tag vX.Y.Z`.
   * Push branch and tags upstream:
     ```bash
     git push origin main
     git push origin vX.Y.Z
     ```
   * Under AGENTS.md §13 ('No Push = Not Done'), a release is strictly incomplete until pushed upstream.

---

## 3. Release Execution Protocol

1. **Verify Clean Tree:** Ensure all active implementation work is committed and verified.
2. **Increment SemVer:** Update canonical manifests:
   - `version.json`
   - `package.json`
   - `prompt-version.template.json`
   - `readme.md` (header badge and quick-install one-liners)
   - `changelog.md` (add version header and categorized changes)
   - `public/health-score.json` (if present)
3. **Author Release Notes:** Generate `.ai-memory/release/release-notes-vX.Y.Z.md` with:
   - Categorized changes (`Added`, `Changed`, `Fixed`, `Removed`).
   - Copy-pasteable quick install one-liners (PowerShell `irm ... | iex` and Bash `curl ... | bash`).
   - List of release assets with SHA256 checksum references.
4. **Commit, Tag, and Push:**
   ```bash
   gitmap cpr "vX.Y.Z - <summary>"
   git tag vX.Y.Z
   git push origin main --tags
   ```

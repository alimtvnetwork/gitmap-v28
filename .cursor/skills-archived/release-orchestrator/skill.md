---
name: release-orchestrator
description: >-
  Execute full automated release orchestration, semantic version bumping, branch management, and tag creation using Python scripts.
---

# Automated Release Orchestrator

Execute full automated release orchestration, semantic version bumping, branch management, and tag creation using Python scripts.

## Core Directives & 5-Step Release Branching Mandate

1. Determine bump tier (MINOR default per Rule 0, reset PATCH to 0).
2. Verify git clean status before release execution.
3. **Mandatory Smart Targeted Pre-Release Verification:** Build and test ONLY packages failed in the stack trace and packages changed between the last git hash and current working tree (`git diff --name-only HEAD~1`), persisting modified files to `.ai-memory/temp/recent-file-changes.json`. Execute targeted verification via `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or alias `--smart`, `-s`), `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`, or `python 03-ai-scripts/06-cicd-local-runner.py --pkg <affected_pkg>` (with optional `--fast` heatmap filtering) to verify 100% green passing (`exit 0`) without running extraneous test suites, spellcheckers, or unrelated checks. In-flight heartbeats emit every 25s (`--heartbeat-interval 25.0`); wait dynamically via `.ai-memory/temp/runner-eta.json` rather than busy-polling.
4. **Test Inventory Validation:** Cross-reference `.ai-memory/temp/recent-file-changes.json` with `.ai-memory/test-inventory.json` to verify that all test suites covering recently modified files pass completely.
5. **Mandatory 5-Step Release Branching Lifecycle:**
   - **Step 1:** Create and switch to a dedicated release branch: `git checkout -b release/vX.Y.Z`.
   - **Step 2:** Bump the version using the dedicated Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) on the release branch.
     - **CRITICAL REPOSITORY ADAPTATION & SCRIPT REPAIR:** Inspect the target repository architecture to identify where versions are defined (`version.json`, `package.json`, `Cargo.toml`, `pyproject.toml`, `go.mod`, etc.), where and how they will change (`readme.md`, `changelog.md`, install scripts), and post-bump synchronization commands (`npm run sync`, `go generate ./...`). If the bump script is missing, outdated, or lacks support for this repository's version pin sites, **the agent MUST fix or recreate the Python bump script immediately** before executing the release.
   - **Step 3:** Commit version bump changes in the release branch: `release: vX.Y.Z <scope>`.
   - **Step 4:** Create the annotated git tag on the release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
   - **Step 5:** Put the release commit back to the `main` branch (`git checkout main && git merge release/vX.Y.Z`), push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore the starting branch.
6. **Automated Release Execution (2-Tier Toolchain):**
   - **Tier 1 (GitMap Native Release - PRIMARY):** `gitmap release --bump patch -y` or `gitmap release --bump minor -y` (alias: `gitmap r -y`). Automatically manages release branch, tag generation, and remote push.
   - **Tier 2 (Python Release Orchestrator - FALLBACK):** `python 03-ai-scripts/29-release-orchestrator.py --tier patch` (or `--tier minor`) to automate the 5-step lifecycle.

---

## Fast File Discovery & Release Context Toolchain (Mandatory Acceleration)

To rapidly discover version manifests, changelog entries, release notes, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the 2-tier toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "*version*" [-ext <ext>]` (alias `gitmap f`)
- **Inspect Changelog & Release Notes:** `gitmap changelog` (alias `gitmap cl [ver]`)
- **List Prior Release Tags:** `gitmap list-versions --limit 5` (alias `gitmap lv`)
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Step Error Logs:** `gitmap pipeline error-logs` (or alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Pipeline Runner Targets & Cache Table:** `gitmap pipeline details` (or alias `gitmap pd`)
- **Stream Manifest or Config:** `gitmap cat version.json` (zero disk writes)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 20`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`

> [!NOTE]
> **Release Verification Allowance:** Release workflows are explicitly authorized to execute targeted pre-release quality gates (e.g. `python 03-ai-scripts/06-cicd-local-runner.py run-smart`, `--changed-only`, or `--pkg <target>`) and create release branches, tags, and commits.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the completion of the release workflow, after version bumping, release notes generation, and tagging, verify that the release branch, release tag, and updated main branch are all pushed to `origin`.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. All modified files across the turn MUST be accumulated and committed together in a SINGLE grouped atomic commit at the final step before pushing!

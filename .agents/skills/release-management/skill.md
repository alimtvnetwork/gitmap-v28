---
name: release-management
description: Execute full release ceremony, SemVer version bumps, package synchronization, and changelog updates.
---

# Release Management

Conducts the standardized release ceremony with the mandatory 5-step release branching lifecycle.

## Steps & 5-Step Release Branching Mandate

1. Determine bump tier (MINOR default per Rule 0, reset PATCH to 0).
2. **Mandatory Pre-Release Unit Tests & CI/CD Verification:** Execute `python 03-ai-scripts/06-cicd-local-runner.py --run-tests` and verify all unit test suites, AST checks, and quality gates pass 100% green (`exit 0`).
3. **Test Inventory Validation:** Cross-reference `.ai-memory/temp/recent-file-changes.json` with `.ai-memory/test-inventory.json` to verify that all test suites covering recently modified files pass completely.
4. **Mandatory 5-Step Release Branching Lifecycle:**
   - **Step 1:** Create and switch to a dedicated release branch: `git checkout -b release/vX.Y.Z`.
   - **Step 2:** Bump the version using the dedicated Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) on the release branch, updating `package.json`, `version.json`, `readme.md`, `changelog.md`, and running sync scripts.
   - **Step 3:** Commit version bump changes in the release branch: `release: vX.Y.Z <scope>` (or `chore(release): bump version to X.Y.Z`).
   - **Step 4:** Create the annotated git tag on the release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
   - **Step 5:** Put the release commit back to the `main` branch (`git checkout main && git merge release/vX.Y.Z`), push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore the starting branch.
5. Use `03-ai-scripts/29-release-orchestrator.py` to automate this complete 5-step lifecycle.

---

## Fast File Discovery via Python Toolchain (Mandatory Acceleration)

To rapidly discover version manifests, changelog entries, release notes, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the Python discovery scripts first:
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 20`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`

> [!NOTE]
> **Release Verification Allowance:** Release workflows are explicitly authorized to execute pre-release quality gates (`python 03-ai-scripts/06-cicd-local-runner.py --run-tests` or `--skip-tests` for emergency runs) and create release branches, tags, and commits.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the completion of the release workflow, after version bumping, release notes generation, and tagging, verify that the release branch, release tag, and updated main branch are all pushed to `origin`.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. All modified files across the turn MUST be accumulated and committed together in a SINGLE grouped atomic commit at the final step before pushing!

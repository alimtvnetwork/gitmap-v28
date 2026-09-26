---
name: release-now-without-test
description: >-
  Execute instant release orchestration without tests, version bumping, release branching, tagging, and pushing while strictly avoiding builds and CI/CD monitoring.
---

# Release Now Without Test (Instant Release Orchestrator)

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

## Core Directives & 5-Step Release Branching Mandate

1. Determine bump tier (MINOR default per Rule 0, reset PATCH to 0; override only on explicit user request).
2. Commit any outstanding uncommitted files on the current working branch before switching branches.
3. **Mandatory 5-Step Release Branching Lifecycle:**
   - **Step 1:** Create and switch to a dedicated release branch: `git checkout -b release/vX.Y.Z`.
   - **Step 2:** Bump the version using the dedicated Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) on the release branch. Update `version.json`, `package.json`, `readme.md`, `changelog.md`, and run sync.
   - **Step 3:** Commit version bump changes on the release branch: `release: vX.Y.Z <scope>`.
   - **Step 4:** Create annotated git tag on the release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
   - **Step 5:** Put the release commit back to `main` (`git checkout main && git merge release/vX.Y.Z`), push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore the starting branch.
4. **Execute Automated Release (2-Tier Toolchain):**
   - **Tier 1 (GitMap Native Release - PRIMARY):**
     ```bash
     gitmap release --bump <minor|patch> -y
     ```
   - **Tier 2 (Python Release Orchestrator - FALLBACK):**
     ```bash
     python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major> --scope "<scope>" --skip-tests
     ```

---

## Three Mandatory Strict Avoidances (The Core Exceptions)

> [!CAUTION]
> **TOTAL BAN ON BUILDING, LOCAL TESTING, AND POST-PUBLISH CI/CD CHECKING:**
>
> 1. **STRICTLY AVOID BUILDING THE CODE (TOTAL BAN):**
>    - NEVER execute build verification commands (`go build`, `npm run build`, `cargo build`, `make build`, compiler invocations).
>    - Fast release requires zero local build overhead.
>
> 2. **STRICTLY AVOID RUNNING ANY CI/CD LOCALLY (TOTAL BAN):**
>    - NEVER run `python 03-ai-scripts/06-cicd-local-runner.py`, unit test suites (`go test ./...`, `pytest`, `npm test`), or local gate checkers.
>    - Pre-release verification is completely bypassed via `--skip-tests`.
>
> 3. **STRICTLY AVOID CHECKING CI/CD AFTER PUBLISH (TOTAL BAN):**
>    - NEVER monitor, query, or poll remote CI/CD pipelines (GitHub Actions, GitLab CI) after publishing.
>    - DO NOT run `gitmap pipeline-ai`, `gitmap pl-ai status`, or loop `gh run view`.
>    - Once `git push` is completed and the starting branch is restored, end the turn immediately and emit the final release summary.

---

## Fast File Discovery via 2-Tier Toolchain (Mandatory Acceleration)

To rapidly discover version manifests, changelog entries, release notes, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "*version*" [-ext <ext>]` (alias `gitmap f`)
- **Inspect Changelog & Release Notes:** `gitmap changelog` (alias `gitmap cl [ver]`)
- **List Prior Release Tags:** `gitmap list-versions --limit 5` (alias `gitmap lv`)
- **Stream Manifest or Config:** `gitmap cat version.json` (zero disk writes)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 20`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **TOP-INSTRUCTION PRIORITY MANDATE:** Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the completion of the release workflow, verify that the release branch, release tag, and updated main branch are all pushed to `origin`.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** All modified version files across the turn MUST be committed together in a single atomic commit on the release branch before merging and pushing!
- [ ] **RESTORE STARTING BRANCH:** The active branch MUST be restored to the starting branch before concluding the turn.

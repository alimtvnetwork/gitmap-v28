# Release Now Without Test — Fast Release Orchestrator & Branch Lifecycle (must follow)

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

Trigger phrases: `release now without test`, `release now`, `release without test`, `instant release`, `fast release`, `publish now without test`, `emergency release without test`.

If the user commands an immediate release without tests, instant release, or fast version bump with release branching, execute this prompt.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of release workflows have suffered from unnecessary delays, premature build commands, and post-publish polling stalls: running full test suites when told not to, blocking on local compilers, hanging in loops checking remote GitHub Actions after publishing, wrong step counts, uncommitted dirty working trees before branching, uppercase filenames, and failing to restore the starting branch. The mandate here is crystal clear: execute the complete release ceremony with release branching, version bumping, tagging, and remote pushing, while strictly avoiding all building, local testing, and post-publish CI/CD monitoring. Going deep IS the job. Violating this is auto-reject on the same tier as RULE 0.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/release-now-without-test/skill.md` exists.
2. If it does NOT exist, create it now. Write the core instructions of this prompt to `.agents/skills/release-now-without-test/skill.md` with frontmatter:

   ```yaml
   ---
   name: release-now-without-test
   description: >-
     Execute instant release orchestration without tests, version bumping, release branching, tagging, and pushing while strictly avoiding builds and CI/CD monitoring.
   ---
   ```

3. Once installed, load it on-demand via progressive disclosure for all future runs.

---

## Three Mandatory Strict Avoidances (The Core Exceptions)

> [!CAUTION]
> **TOTAL BAN ON BUILDING, LOCAL TESTING, AND POST-PUBLISH CI/CD CHECKING:**
> 
> 1. **STRICTLY AVOID BUILDING THE CODE (TOTAL BAN):**
>    - NEVER execute build verification commands (`go build`, `npm run build`, `cargo build`, `make build`, compiler invocations).
>    - An instant release must not be blocked or delayed by local compiler passes.
> 
> 2. **STRICTLY AVOID RUNNING ANY CI/CD LOCALLY (TOTAL BAN):**
>    - NEVER run `python 03-ai-scripts/06-cicd-local-runner.py`, unit test suites (`go test ./...`, `pytest`, `npm test`), or local gate checkers.
>    - Pre-release verification is completely bypassed. The release orchestrator MUST be run with `--skip-tests`.
> 
> 3. **STRICTLY AVOID CHECKING CI/CD AFTER PUBLISH (TOTAL BAN):**
>    - NEVER monitor, query, or poll remote CI/CD pipelines (GitHub Actions, GitLab CI) after publishing.
>    - DO NOT run `gitmap pipeline-ai`, `gitmap pl-ai status`, or loop `gh run view`.
>    - Once `git push` is completed and the starting branch is restored, end the turn immediately and emit the final release summary.

---

## RULE 0: Version Calculation Standard

1. Read the canonical version source for this repo (`version.json` or `package.json`).
2. Default bump tier is **MINOR**: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to `0`.
3. Only bump PATCH if the user explicitly specified `patch`.
4. Only bump MAJOR if the user explicitly specified `major` or breaking change.
5. State previous and new versions explicitly in output before touching files.

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

---

## Master Architecture: Heavy-Lifting Release Script (`03-ai-scripts/29-release-orchestrator.py`)

All release operations (version bumping, commit creation, release branching, git tagging, and branch reversion) MUST be executed through the centralized heavy-lifting script in the AI scripts directory with the mandatory `--skip-tests` flag:

```bash
python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major> --scope "<Release summary>" --skip-tests
```

### The Reuse-or-Create Protocol

1. **Check for Script:** Inspect if `03-ai-scripts/29-release-orchestrator.py` exists on disk.
2. **If Script Exists:**
   - Run the script directly with `--skip-tests` and the desired tier:
     ```bash
     python 03-ai-scripts/29-release-orchestrator.py --tier minor --scope "<Release summary>" --skip-tests
     ```
3. **If Script Is Missing:**
   - The AI agent MUST immediately construct `03-ai-scripts/29-release-orchestrator.py` inside `03-ai-scripts/` using Python standard libraries (`argparse`, `json`, `os`, `re`, `subprocess`, `sys`, `pathlib`, `datetime`).
   - Support `--skip-tests` to skip pre-release verification entirely.
   - Register it in `03-ai-scripts/readme.md` under slot 29.
   - Execute the script with `--skip-tests` to fulfill the release.

---

## Git Release Lifecycle & 5-Step Branching Mandate

The release orchestrator and all release triggers MUST strictly execute this 5-step release lifecycle:

```text
[Start on original_branch (e.g. main or feature/my-work)]
                   │
                   ▼
1. Detect & Store original_branch (git rev-parse --abbrev-ref HEAD)
                   │
                   ▼
[STEP 1] Create & Checkout Release Branch FIRST:
         git checkout -b release/vX.Y.Z
                   │
                   ▼
[STEP 2] Bump SemVer via Python Bump Script (repository-aware):
         python 03-ai-scripts/37-bump-version.py --tier <tier> --scope "<scope>"
         (Updates version.json, package.json, readme.md, changelog.md, and runs sync)
                   │
                   ▼
[STEP 3] Stage & Commit on Release Branch:
         git commit -m "release: vX.Y.Z <scope>"
                   │
                   ▼
[STEP 4] Create Annotated Git Tag on Release Commit:
         git tag -a vX.Y.Z -m "Release vX.Y.Z"
                   │
                   ▼
[STEP 5] Put Commit Back to Main Branch & Push:
         git checkout main
         git merge release/vX.Y.Z
         git push origin main
         git push origin release/vX.Y.Z
         git push origin vX.Y.Z
                   │
                   ▼
[STEP 6] Revert Working Tree to original_branch (if different from main)
                   │
                   ▼
[Finish: Active branch is verified to be original_branch]
```

### Critical Rules for Branch Reversion:

- **NEVER assume `main` or `master`:** Releases may be triggered from feature branches, bugfix branches, or release candidates. The script MUST record the starting branch and revert to that exact branch after merging to main.
- **Fail-Safe Restoration:** Even if pushing or downstream steps fail, the script's `finally:` block MUST ensure the working tree is safely checked back out to `original_branch`.

---

## Hard Rules & Pre-Flight Checks

- [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
- [ ] Clean Working Tree Check: You MUST commit any outstanding/current files on the working branch BEFORE running the orchestrator script. Never branch with uncommitted edits.
- [ ] NO BUILDING (TOTAL BAN): Never build the code or verify compilation before, during, or after release.
- [ ] NO LOCAL CI/CD (TOTAL BAN): Never run `06-cicd-local-runner.py` or unit tests before releasing.
- [ ] NO POST-PUBLISH CI/CD CHECKING (TOTAL BAN): Never monitor or query remote pipelines after git push.
- [ ] Step 1 is create the release branch FIRST: `git checkout -b release/vX.Y.Z`. Never bump directly on main.
- [ ] Step 2 is bump version using the repository-adapted Python bump script on the release branch.
- [ ] Step 3 is commit on the release branch.
- [ ] Step 4 is tag on the release branch commit.
- [ ] Step 5 is merge the release branch commit back into `main` and push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to origin.
- [ ] No explicit boolean checks (`if is_success == True:` is banned; use `if is_success:`).
- [ ] All filenames must be strictly lowercase (e.g. `readme.md`, `changelog.md`).
- [ ] Relative Git paths only (no `file:///` URIs or absolute filesystem paths).
- [ ] The release commit, release branch, and tag MUST be pushed to Git.
- [ ] The working tree must be on the original branch when the task completes.

---

## Actionable Execution Checklist

1. [ ] TOP-INSTRUCTION PRIORITY MANDATE: Whatever directives, constraints, checklists, or instructions are given before this section or prompt (user preamble, header constraints, prior instructions) are verified as highest priority and non-negotiable, overriding all lower-level guidelines below.
2. [ ] Identify starting branch: `git rev-parse --abbrev-ref HEAD`.
3. [ ] Commit any outstanding/current files on the working branch before proceeding.
4. [ ] Check for `03-ai-scripts/29-release-orchestrator.py` and bump version script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`).
5. [ ] Run `python 03-ai-scripts/29-release-orchestrator.py --tier <tier> --scope "<scope>" --skip-tests`.
6. [ ] Verify Step 1: Release branch created and checked out (`release/vX.Y.Z`).
7. [ ] Verify Step 2: Version bumped across `version.json`, `package.json`, `readme.md`, `changelog.md`.
8. [ ] Verify Step 3: Committed on release branch (`release: vX.Y.Z <scope>`).
9. [ ] Verify Step 4: Tag `vX.Y.Z` created on release commit.
10. [ ] Verify Step 5: Merged to `main`, and pushed `main`, `release/vX.Y.Z`, and `vX.Y.Z` tag to `origin`.
11. [ ] Verify Step 6: Working tree restored to starting branch (`git branch --show-current`).
12. [ ] DO NOT run any builds, local CI/CD tests, or post-publish pipeline monitoring.
13. [ ] Output release summary detailing starting branch, version bump, release branch, tag, and restored active branch.

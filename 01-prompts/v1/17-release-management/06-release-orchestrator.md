# Automated Release Orchestrator & Branch Lifecycle — Release Management (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Trigger phrases: `orchestrate release`, `release orchestrator`, `bump with release branch`, `release and tag`, `automated release`.

If the user requests a release, version bump with release branching, or automated tag orchestration, execute this prompt.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with placeholders, folders skimmed, open ambiguities ignored, CI/CD issues forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. Read the whole codebase, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/release-orchestrator/skill.md` exists.
2. If it does NOT exist, create it now. Write the core instructions of this prompt to `.agents/skills/release-orchestrator/skill.md` with frontmatter:

   ```yaml
   ---
   name: release-orchestrator
   description: >-
     Execute full automated release orchestration, semantic version bumping, branch management, and tag creation using Python scripts.
   ---
   ```

3. Once installed, load it on-demand via progressive disclosure for all future runs.

---

## RULE 0: Version Calculation Standard

1. Read the canonical version source for this repo (`version.json` or `package.json`).
2. Default bump tier is **MINOR**: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to `0`.
3. Only bump PATCH if user explicitly specified `patch`.
4. Only bump MAJOR if user explicitly specified `major` or breaking change.
5. State previous and new versions explicitly in output before touching files.

---

## Fast File Discovery via Python Toolchain (Mandatory Acceleration)

To rapidly discover version manifests, changelog entries, release notes, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the Python discovery scripts first:
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Step Error Logs:** `gitmap pipeline error-logs` (or alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Pipeline Runner Targets & Cache Table:** `gitmap pipeline details` (or alias `gitmap pd`)
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 20`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`

> [!IMPORTANT]
> **SMART TARGETED PRE-RELEASE TESTING (FASTEST PATH TO RELEASE):**
> When verifying code prior to release, do NOT run heavy full repository test suites, spellcheckers, or unrelated packages.
> 1. **Priority Incremental Runner:** Execute smart incremental Go tests and gates via `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or alias `smart`, `--smart`, `-s`), which inspects Git changed files, builds ONLY changed packages into OS temp, and runs the Quad Runner.
> 2. **Specific Package Targeting:** Isolate and test ONLY packages and functions identified in failing stack traces or modified packages: `python 03-ai-scripts/06-cicd-local-runner.py --pkg <target_package_or_file>`.
> 3. **Heatmap & Fast-Path Testing:** Use `--fast` to run only hot and warm tests based on `.ai-memory/test-heatmap.json`, skipping cold tests (`python 03-ai-scripts/06-cicd-local-runner.py --fast`).
> 4. **Changed Packages from Last Git Hash:** Isolate packages changed between the previous git hash and current working tree (`git diff --name-only HEAD~1` or `git status --porcelain`) using `python 03-ai-scripts/06-cicd-local-runner.py --changed-only`.
> 5. **Change State Persistence:** Record modified files into `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).
> 6. **In-Flight Heartbeats & ETA Wait:** The local runner emits heartbeats every 25s (`--heartbeat-interval 25.0`) and writes status to `.ai-memory/temp/runner-eta.json`. Agents must sleep for 60s or remaining ETA rather than busy-polling.
> 7. Once targeted checks pass green, trigger `python 03-ai-scripts/29-release-orchestrator.py` immediately without delaying the release.

---

## Master Architecture: Heavy-Lifting Release Script (`03-ai-scripts/29-release-orchestrator.py`)

All release operations (version bumping, commit creation, release branching, git tagging, and branch reversion) MUST be executed through the centralized heavy-lifting script in the AI scripts directory:

```bash
python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major>
```

### The Reuse-or-Create Protocol

1. **Check for Script:** Inspect if `03-ai-scripts/29-release-orchestrator.py` exists on disk.
2. **If Script Exists:**
   - Run the script directly with the required arguments:
     ```bash
     python 03-ai-scripts/29-release-orchestrator.py --tier minor --scope "<Release summary>"
     ```
3. **If Script Is Missing:**
   - The AI agent MUST immediately construct `03-ai-scripts/29-release-orchestrator.py` inside `03-ai-scripts/` using Python standard libraries (`argparse`, `json`, `os`, `re`, `subprocess`, `sys`, `pathlib`, `datetime`).
   - Register it in `03-ai-scripts/readme.md` under slot 29.
   - Execute the freshly created script to fulfill the release.

---

## Standalone Bump-Version Recovery & Repository Adaptation Mandate

> [!IMPORTANT]
> **REPOSITORY ARCHITECTURE INSPECTION & SCRIPT REPAIR:**
> If the repository does not have a project-level bump version script (such as `03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`), or if the existing script is outdated or incomplete:
>
> 1. **Inspect Repository Architecture:**
>    - Determine where versions are defined (`version.json`, `package.json`, `Cargo.toml`, `pyproject.toml`, `go.mod`, etc.).
>    - Determine where and how versions will change (`readme.md`, `changelog.md`, badges, install scripts).
>    - Determine post-bump synchronization commands (`npm run sync`, `go generate ./...`).
> 2. **Fix or Create the Bump Script:**
>    - The orchestrator script must not fail or halt. It MUST autonomously create, fix, or adapt `03-ai-scripts/37-bump-version.py` (or `.ai-memory/release/bump_versions.py`) based specifically on how the target repository tracks versions and where files need to change.
> 3. **In-Place File Updates:** The bump script must update:
>    - `version.json`: Update `"version"` to `next_version` and `"releaseDate"` to today's UTC date (`YYYY-MM-DD`).
>    - `package.json`: Update `"version"` to `next_version`.
>    - `readme.md`: Pin the new release version in badges, install snippets, and header versions.
>    - `changelog.md`: Prepend the release entry under `# Changelog` with the release version, UTC date, and changelog bullets.
>    - Run synchronization commands (`npm run sync`) to regenerate spec trees and manifests.

---

## Git Release Lifecycle & 5-Step Branching Mandate

The release orchestrator and all release triggers MUST strictly execute this 5-step release lifecycle:

```
[Start on original_branch (e.g. main or feature/my-work)]
                   │
                   ▼
1. Detect & Store original_branch (git rev-parse --abbrev-ref HEAD)
                   │
                   ▼
[STEP 1] Create & Checkout Release Branch:
         git checkout -b release/vX.Y.Z
                   │
                   ▼
[STEP 2] Bump SemVer via Python Bump Script (repository-aware):
         python 03-ai-scripts/37-bump-version.py --tier <tier> --scope "<scope>"
         (or python .ai-memory/release/bump_versions.py)
         Updates version.json, package.json, readme.md, changelog.md, and runs sync.
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

- [ ] You MUST resolve the Git state by committing any outstanding/current files on the working branch BEFORE running the orchestrator script. The orchestrator script will NOT stage all changes; it will only stage version-related files (`version.json`, `package.json`, `changelog.md`, `readme.md`).
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

- [ ] 1. Identify starting branch: `git rev-parse --abbrev-ref HEAD`.
- [ ] 2. Resolve the Git state by committing any outstanding/current files on the working branch before proceeding.
- [ ] 3. Check for `03-ai-scripts/29-release-orchestrator.py` and bump version script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`). If missing, create/adapt them based on repository structure.
- [ ] 4. Step 1: Create and switch to release branch: `git checkout -b release/vX.Y.Z`.
- [ ] 5. Step 2: Bump version on the release branch using the repository-aware Python bump script.
- [ ] 6. Step 3: Commit version bump changes on the release branch: `release: vX.Y.Z <scope>`.
- [ ] 7. Step 4: Create annotated tag on the release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
- [ ] 8. Step 5: Switch to `main`, merge `release/vX.Y.Z`, and push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`.
- [ ] 9. Verify that the active git branch is restored to the starting branch (`git branch --show-current`).
- [ ] 10. Output release summary detailing starting branch, version bump, release branch, tag, and restored active branch.

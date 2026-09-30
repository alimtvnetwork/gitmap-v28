# Minor Version Bump — Release Management (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Trigger Keywords & Aliases: `minor bump`, `minor release`, `bump minor version`, `bump version`, `default release`, `feature bump`

Execute this prompt whenever a minor version bump or feature addition release is requested. Per Rule 0, MINOR is the default release tier when no specific tier is designated.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.ai-memory/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.ai-memory/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before any execution, check if this prompt is installed as a native Antigravity Skill.

1. Check if `.agents/skills/release-management/skill.md` exists.
2. If it does NOT exist, create it now. Write the core instructions of this prompt to `.agents/skills/release-management/skill.md` with frontmatter:

   ```yaml
   ---
   name: release-management
   description: >-
     Execute full release ceremony, SemVer version bumps, package synchronization, and changelog updates.
   ---
   ```

3. Once installed, load it on-demand via progressive disclosure for all future runs.

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
> **Release Verification Allowance:** Release workflows are explicitly authorized to execute pre-release quality gates (`python 03-ai-scripts/06-cicd-local-runner.py --run-tests` or `--skip-tests` for emergency runs) and create release branches, tags, and commits.

---

## RULE 0: Minor Version Calculation Standard (Default Tier)

1. Read the canonical version source for THIS repo (discover it: `version.json`, `package.json` `"version"`, or whatever single file the repo treats as the version of record). Do not guess.
2. **Minor Bump Arithmetic:** `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`.
   - **MANDATORY RESET:** `PATCH` MUST reset strictly to `0`!
   - Example: `v6.42.3` -> `v6.43.0`.
3. State the previous version and new version explicitly in the reply before touching any file.
4. Do NOT ask "minor or patch?". Do NOT open plan mode. When in doubt: MINOR.

---

## Repository-Aware Python Bump Script Creation & Repair Mandate

> [!IMPORTANT]
> **CRITICAL SCRIPT CREATION & REPAIR REQUIREMENT:**
> The AI agent MUST NOT assume that a static bump script works across different repositories without inspection.
> You MUST inspect the target repository's architecture and **create, adapt, or fix the Python version bump script** (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) based specifically on how the repository tracks versions and where files need to change:
>
> 1. **Identify Where Versions Are Defined:**
>    - Root manifests: `version.json`, `package.json`, `Cargo.toml`, `pyproject.toml`, `go.mod`, `pom.xml`, `build.gradle`, or `.version`.
>    - If multi-package or monorepo: check for workspace packages, sub-packages, and template files (e.g. `prompt-version.template.json`).
> 2. **Identify Where and How Versions Must Change:**
>    - Root `readme.md`: badges, header versions, curl/irm install snippets, and download links.
>    - Changelogs: `changelog.md` and spec changelogs (e.g. `02-spec/19-main-worker-service/98-changelog.md`).
>    - Install scripts: `install.sh`, `install.ps1`, or platform release assets.
> 3. **Identify Post-Bump Synchronization Steps:**
>    - Does the repository have `npm run sync` to regenerate spec trees, manifests, or health badges?
>    - Does it need `go generate ./...` or code generation?
> 4. **Fix or Create the Python Script:**
>    - If `03-ai-scripts/37-bump-version.py` (or `.ai-memory/release/bump_versions.py`) is missing, outdated, or lacks support for this repository's version pin sites, **the AI agent MUST fix or rewrite the script immediately** to handle all identified files before running the release.
>    - Test the bump script with `--dry-run` or direct invocation to verify it executes cleanly without error.

---

## Mandatory 5-Step Release Branching Lifecycle

All minor version bumps and releases MUST strictly follow this 5-step Git release lifecycle:

```text
[Start on original_branch (e.g. main or feature/work)]
                   │
                   ▼
1. Detect & Store original_branch (git rev-parse --abbrev-ref HEAD)
                   │
                   ▼
[STEP 1] Create & Switch to Dedicated Release Branch:
         git checkout -b release/vX.Y.0
         (NEVER bump directly on main or active work branch)
                   │
                   ▼
[STEP 2] Bump Version on Release Branch via Python Script:
         python 03-ai-scripts/37-bump-version.py --tier minor --scope "<Feature release scope>"
         (or python .ai-memory/release/bump_versions.py --type minor)
         Updates version.json, package.json, readme.md, changelog.md, and runs sync.
                   │
                   ▼
[STEP 3] Stage & Commit on Release Branch:
         git commit -m "release: vX.Y.0 <feature release scope>"
                   │
                   ▼
[STEP 4] Create Annotated Git Tag on Release Commit:
         git tag -a vX.Y.0 -m "Release vX.Y.0"
                   │
                   ▼
[STEP 5] Put Commit Back to Main Branch & Push:
         git checkout main
         git merge release/vX.Y.0
         git push origin main
         git push origin release/vX.Y.0
         git push origin vX.Y.0
                   │
                   ▼
[RESTORE] Revert Working Tree to original_branch (if different from main):
          git checkout <original_branch>
```

### Detailed Breakdown of the 5 Steps:

1. **Step 1: Create & Switch to Release Branch:**
   Create and checkout the dedicated release branch first: `git checkout -b release/vX.Y.0`. Releases MUST NOT be committed directly to `main` without a release branch.
2. **Step 2: Bump Version on Release Branch via Python Script:**
   Execute the dedicated Python bump script (`03-ai-scripts/37-bump-version.py --tier minor` or `.ai-memory/release/bump_versions.py --type minor`) on the release branch.
   - The bump version script must be adapted based on the target repository architecture to update `version.json`, `package.json`, `readme.md`, `changelog.md`, and install snippets.
   - If the script needs changes for this repo, fix it before running.
3. **Step 3: Commit in Release Branch:**
   Stage and commit all version bump and generated release files on the release branch: `git commit -m "release: vX.Y.0 <feature release scope>"`.
4. **Step 4: Create Annotated Git Tag:**
   Create the annotated tag on that release commit: `git tag -a vX.Y.0 -m "Release vX.Y.0"`.
5. **Step 5: Put Commit Back to Main Branch & Push:**
   Switch to `main` (`git checkout main`), merge the release branch commit (`git merge release/vX.Y.0`), push `main`, `release/vX.Y.0`, and tag `vX.Y.0` to `origin`, then restore the starting branch if different from `main`.

---

## Pre-Flight Checks (Fail-Fast)

- **Git Sync & Clean State:** Run `git status`. If there are pending uncommitted changes, resolve and commit them first. Then run `git pull` to fetch and merge upstream changes.
- **Force Override:** If user said `force`, ignore cached scripts and regenerate the Python bump script from scratch.
- **Idempotency Guard:** If the canonical version file already equals the computed new version, STOP. Someone half-ran a release. Detect what is already done, resume from the first incomplete step, do NOT double-bump.
- **Placeholder Guard:** If the previous version's changelog entry is empty or a placeholder (`TBD`, `WIP`, no bullets), refuse to release until it is filled or the user overrides.
- **Date Source:** The release date is UTC today. Get it from `date -u +%Y-%m-%d`. Do not invent it.

---

## Mandatory Pinning & Changelog (Fatal if missed)

1. **Changelog:** You MUST read the `"changelog"` configuration from `version.json` (e.g. `file_path` and `format`) and append the proper changelog correctly according to that format.
   - Header format: `## [vX.Y.0] YYYY-MM-DD <short headline>`
   - Document new features and changes under `### Added` or `### Changed`.
2. **Root README Pinning (Fatal if missed):** You MUST pin the latest release version into the root `readme.md` file. Rewrite every occurrence of the previous version (`vX.Y.Z` and bare `X.Y.Z`) in badges, install snippets, and inline references. It is FATAL if you do not update the version pins in the root README file!
3. **Test File Ban:** You MUST NOT read, scan, or modify test files (e.g., `*_test.*`, `*.spec.*`, `test/*`) when discovering or updating versions. Test files contain mock data, and updating mock data corrupts the tests.
4. **Version Inheritance Protocol:** The root `version.json` file is the strict Single Source of Truth. If sub-components are set to `"inherit"`, DO NOT modify them independently; they automatically scale with the global version. Always bump the global root `"version"` property.

---

## Install Section (Non-Negotiable)

When generating release notes, changelogs, or README updates, you MUST include the install snippet:

### Install <Project Name> vX.Y.0

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.0/install.sh | bash -s -- ".ai-memory/prompts" "vX.Y.0"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "vX.Y.0"`

*(Note: You MUST dynamically discover the `<owner>/<repo>` by running `git config --get remote.origin.url`.)*

---

## Publish Platform Release with Quick Install One-Liners (FATAL IF MISSED ON GITHUB/GITLAB)

When creating the GitHub / GitLab release (via `bump_versions.py --create-release` or CLI):
1. **Quick Install One-Liners:**
   - Binary repos: `irm https://github.com/<owner>/<repo>/releases/download/vX.Y.0/install.ps1 | iex` (PowerShell) and `curl -fsSL https://github.com/<owner>/<repo>/releases/download/vX.Y.0/install.sh | bash` (Bash).
   - Script/Prompt repos: `Invoke-WebRequest ...` (PowerShell) and `curl -sL ... | bash -s -- ".ai-memory/prompts" "vX.Y.0"` (Bash).
2. **Extracted Changelog:** The exact `[vX.Y.0]` section from `changelog.md`.
3. **Pass Notes File:**
   ```bash
   gh release create "vX.Y.0" --title "vX.Y.0" --notes-file ".ai-memory/release/release-notes-vX.Y.0.md" --generate-notes
   ```
4. **STRICT PROHIBITION:** NEVER execute `gh release create <tag> --generate-notes` without `--notes-file`. Bare `--generate-notes` dumps commit hashes only and strips the installation one-liners!

---

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (e.g. `.ai-memory/release/bump_versions.py`, runner, autofixer, linter) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, release auto-bumpers, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python Release Scripts: `.ai-memory/release/` (e.g. `bump_versions.py`) and `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `06-cicd-local-runner.py`, `37-bump-version.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/`, `.ai-memory/cicd-issues/`, and `.ai-memory/release/issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** Pre-release quality gates and test orchestrators must run concurrently using a worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python release scripts (`bump_versions.py`, `06-cicd-local-runner.py`, `37-bump-version.py`) from scratch**, regardless of whether the file already exists on disk.
> 5. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.ai-memory/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `02-spec/03-error-manage/readme.md`, `.ai-memory/plans/readme.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
> 6. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

---

## Issue Logging (MUST, when anything goes wrong)

Path: `.ai-memory/release/issues/xx-<new-version>-<slug>.md` (lowercase). Body:
- Previous version and new version
- Step that failed (number and name)
- Command run and full error output
- Files involved
- Resolution or workaround, or `unresolved`

Then link it from the `### Issues` bullet under the changelog entry.

---

## Ambiguity Handling (Open Questions & Answers)

Ambiguity is not a license to guess. It is a file to write.
- Open: `.ai-memory/ambiguous-questions/01-new-ambiguity/xx-<slug>.md`
- Answered: `.ai-memory/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md`

When answered: `mv` from `01-new-ambiguity/` to `02-ambiguity-resolved/`, flip `Status: resolved`, and append a `## Resolution` block.

---

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ingest canonical version source from disk (`version.json` or `package.json`).
- [ ] Determine new version: `MAJOR.MINOR.PATCH` -> `MAJOR.(MINOR+1).0` (PATCH resets to 0).
- [ ] State previous and new version explicitly in chat.
- [ ] Ensure git repository starts clean (`git status`, commit pending changes, `git pull`).
- [ ] Inspect repository architecture and adapt/fix the Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) based on where versions are tracked and how they change.
- [ ] Step 1: Create and switch to dedicated release branch: `git checkout -b release/vX.Y.0`.
- [ ] Step 2: Bump version on release branch using the repository-aware Python bump script (`python 03-ai-scripts/37-bump-version.py --tier minor` or `python .ai-memory/release/bump_versions.py --type minor`).
- [ ] Read `version.json` for Changelog formatting rules and prepend entry in `changelog.md`.
- [ ] Pin the latest version into root `readme.md` (badges, text, install snippets; FATAL IF MISSED).
- [ ] Run post-bump sync scripts (`npm run sync` or repository generator) to update spec trees and manifests.
- [ ] Step 3: Stage and commit all release changes on release branch: `git commit -m "release: vX.Y.0 <feature release scope>"`.
- [ ] Step 4: Create annotated tag on release commit: `git tag -a vX.Y.0 -m "Release vX.Y.0"`.
- [ ] Step 5: Switch to `main`, merge `release/vX.Y.0`, push `main`, `release/vX.Y.0`, and tag `vX.Y.0` to `origin`, then restore starting branch.
- [ ] Generate release notes file with Quick Install one-liners and changelog, and publish via `gh release create --notes-file` (NEVER bare `--generate-notes`).
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] **File Change Summary:** Provide a highly detailed summary in chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the completion of the release workflow, verify that the release branch, release tag, and updated main branch are all pushed to `origin`.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. All modified files across the turn MUST be accumulated and committed together in a SINGLE grouped atomic commit at the final step before pushing!

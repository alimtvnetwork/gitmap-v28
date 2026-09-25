# Version Bump Documentation Update — Release Management (must follow)

> **Prompt Version:** 2.2.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Trigger Keywords & Aliases: `version bump docs`, `bump docs version`, `update doc versions`, `sync docs version`, `bump documentation version`

Execute this prompt whenever instructed to bump or synchronize version numbers across repository documentation, specs, installation guides, and changelogs.

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

## Fast File Discovery via Python Toolchain (Mandatory Acceleration)

To rapidly discover documentation files, version manifests, changelog entries, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the Python discovery scripts first:
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins in Docs:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 50`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`

---

## Documentation Version Synchronization Protocol

1. **Discover Canonical Version:**
   Read current SemVer from `version.json` or `package.json`. If a new target version is specified by the user, adopt it; otherwise calculate according to SemVer rules (default MINOR per Rule 0).
2. **Search and Replace Across Documentation:**
   - Find all instances of the old version number across documentation files (`readme.md`, `installation.md`, `docs/**/*.md`, `02-spec/**/*.md`).
   - Update version badges, shields.io URLs, release download URLs, and install one-liners.
   - Root `readme.md` MUST be pinned to the new version. After updating, searching for the old version in `readme.md` MUST return zero occurrences outside historic changelog references.
3. **Changelog Formatting:**
   - Prepend new entry under `# Changelog` in `changelog.md`: `## [vX.Y.Z] - YYYY-MM-DD`.
   - Update spec changelogs (e.g. `02-spec/19-main-worker-service/98-changelog.md`) if maintained by the repository.

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
>    - Changelogs: `changelog.md` and spec changelogs.
>    - Documentation files: `02-spec/`, installation guides, tutorials.
> 3. **Identify Post-Bump Synchronization Steps:**
>    - Does the repository have `npm run sync` to regenerate spec trees, manifests, or health badges?
>    - Does it need `go generate ./...` or documentation site generators?
> 4. **Fix or Create the Python Script:**
>    - If `03-ai-scripts/37-bump-version.py` (or `.ai-memory/release/bump_versions.py`) is missing, outdated, or lacks support for this repository's version pin sites, **the AI agent MUST fix or rewrite the script immediately** to handle all identified files before running the release.
>    - Test the bump script with `--dry-run` or direct invocation to verify it executes cleanly without error.

---

## Mandatory 5-Step Release Branching Lifecycle

All documentation version bumps and releases MUST strictly follow this 5-step Git release lifecycle:

```text
[Start on original_branch (e.g. main or feature/work)]
                   │
                   ▼
1. Detect & Store original_branch (git rev-parse --abbrev-ref HEAD)
                   │
                   ▼
[STEP 1] Create & Switch to Dedicated Release Branch:
         git checkout -b release/vX.Y.Z
         (NEVER bump directly on main or active work branch)
                   │
                   ▼
[STEP 2] Bump Version on Release Branch via Python Script:
         python 03-ai-scripts/37-bump-version.py --tier <tier> --scope "<Documentation version sync>"
         (or python .ai-memory/release/bump_versions.py --type <tier>)
         Updates version.json, package.json, readme.md, changelog.md, and runs sync.
                   │
                   ▼
[STEP 3] Stage & Commit on Release Branch:
         git commit -m "docs(release): bump version in documentation to vX.Y.Z"
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
[RESTORE] Revert Working Tree to original_branch (if different from main):
          git checkout <original_branch>
```

### Detailed Breakdown of the 5 Steps:

1. **Step 1: Create & Switch to Release Branch:**
   Create and checkout the dedicated release branch first: `git checkout -b release/vX.Y.Z`. Releases MUST NOT be committed directly to `main` without a release branch.
2. **Step 2: Bump Version on Release Branch via Python Script:**
   Execute the dedicated Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) on the release branch.
   - The bump version script must be adapted based on the target repository architecture to update `version.json`, `package.json`, `readme.md`, `changelog.md`, and install snippets.
   - If the script needs changes for this repo, fix it before running.
3. **Step 3: Commit in Release Branch:**
   Stage and commit all version bump and generated release files on the release branch: `git commit -m "docs(release): bump version in documentation to vX.Y.Z"`.
4. **Step 4: Create Annotated Git Tag:**
   Create the annotated tag on that release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
5. **Step 5: Put Commit Back to Main Branch & Push:**
   Switch to `main` (`git checkout main`), merge the release branch commit (`git merge release/vX.Y.Z`), push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore the starting branch if different from `main`.

---

## Pre-Flight Checks (Fail-Fast)

- **Git Sync & Clean State:** Run `git status`. If there are pending uncommitted changes, resolve and commit them first. Then run `git pull` to fetch and merge upstream changes.
- **Force Override:** If user said `force`, ignore cached scripts and regenerate the Python bump script from scratch.
- **Idempotency Guard:** If the canonical version file already equals the computed new version, STOP. Someone half-ran a release. Detect what is already done, resume from the first incomplete step, do NOT double-bump.
- **Placeholder Guard:** If the previous version's changelog entry is empty or a placeholder (`TBD`, `WIP`, no bullets), refuse to release until it is filled or the user overrides.
- **Date Source:** The release date is UTC today. Get it from `date -u +%Y-%m-%d`. Do not invent it.

---

## Install Section (Non-Negotiable)

When generating release notes, changelogs, or README updates, you MUST include the install snippet:

### Install <Project Name> vX.Y.Z

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.sh | bash -s -- ".ai-memory/prompts" "vX.Y.Z"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "vX.Y.Z"`

*(Note: You MUST dynamically discover the `<owner>/<repo>` by running `git config --get remote.origin.url`.)*

---

## Publish Platform Release with Quick Install One-Liners (FATAL IF MISSED ON GITHUB/GITLAB)

When creating the GitHub / GitLab release (via `bump_versions.py --create-release` or CLI):
1. **Quick Install One-Liners:**
   - Binary repos: `irm https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.ps1 | iex` (PowerShell) and `curl -fsSL https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.sh | bash` (Bash).
   - Script/Prompt repos: `Invoke-WebRequest ...` (PowerShell) and `curl -sL ... | bash -s -- ".ai-memory/prompts" "vX.Y.Z"` (Bash).
2. **Extracted Changelog:** The exact `[vX.Y.Z]` section from `changelog.md`.
3. **Pass Notes File:**
   ```bash
   gh release create "vX.Y.Z" --title "vX.Y.Z" --notes-file ".ai-memory/release/release-notes-vX.Y.Z.md" --generate-notes
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
- [ ] Determine new version according to SemVer rules.
- [ ] State previous and new version explicitly in chat.
- [ ] Ensure git repository starts clean (`git status`, commit pending changes, `git pull`).
- [ ] Inspect repository architecture and adapt/fix the Python bump script (`03-ai-scripts/37-bump-version.py` or `.ai-memory/release/bump_versions.py`) based on where versions are tracked and how documentation files reference them.
- [ ] Step 1: Create and switch to dedicated release branch: `git checkout -b release/vX.Y.Z`.
- [ ] Step 2: Bump version across documentation and manifests on release branch using the repository-aware Python bump script (`python 03-ai-scripts/37-bump-version.py` or `python .ai-memory/release/bump_versions.py`).
- [ ] Search and replace old version strings across `readme.md`, `installation.md`, and docs.
- [ ] Pin the latest version into root `readme.md` (badges, text, install snippets; FATAL IF MISSED).
- [ ] Read `version.json` for Changelog formatting rules and prepend entry in `changelog.md`.
- [ ] Run post-bump sync scripts (`npm run sync` or repository generator) to update spec trees and manifests.
- [ ] Step 3: Stage and commit all release changes on release branch: `git commit -m "docs(release): bump version in documentation to vX.Y.Z"`.
- [ ] Step 4: Create annotated tag on release commit: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`.
- [ ] Step 5: Switch to `main`, merge `release/vX.Y.Z`, push `main`, `release/vX.Y.Z`, and tag `vX.Y.Z` to `origin`, then restore starting branch.
- [ ] Generate release notes file with Quick Install one-liners and changelog, and publish via `gh release create --notes-file` (NEVER bare `--generate-notes`).
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] **File Change Summary:** Provide a highly detailed summary in chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the completion of the release workflow, verify that the release branch, release tag, and updated main branch are all pushed to `origin`.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them. All modified files across the turn MUST be accumulated and committed together in a SINGLE grouped atomic commit at the final step before pushing!

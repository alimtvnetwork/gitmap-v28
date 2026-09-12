# Minor Version Bump — Release Management (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## Required release action

1. Update root `version.json` only: set `version` to the new MINOR version and `releaseDate` to today's UTC date.

## Publish trigger

2. If publishing is requested, create the matching `vX.Y.0` Git tag after the `version.json` change is present on the target branch. The tag triggers the release workflow.

## Mandatory Pinning & Changelog (Fatal if missed)

1. Changelog: You MUST read the `"changelog"` configuration from `version.json` (e.g. `file_path` and `format`) and append the proper changelog correctly according to that format.
2. Root README: You MUST pin the latest release version into the root `readme.md` file. It is FATAL if you do not update the version pins in the root README file!

You must update `version.json`, `changelog.md`, and `readme.md` at a minimum during any bump.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Update version in `version.json`.
- [ ] Read `version.json` for Changelog formatting rules.
- [ ] Add the changelog properly to the targeted changelog file.
- [ ] Pin the latest version into the root `readme.md` file (FATAL IF MISSED).
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

#

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (e.g. `.lovable/release/bump_versions.py`, runner, autofixer, linter) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, release auto-bumpers, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python Release Scripts: `.lovable/release/` (e.g. `bump_versions.py`) and `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `06-cicd-local-runner.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/`, `.lovable/cicd-issues/`, and `.lovable/release/issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** Pre-release quality gates and test orchestrators must run concurrently using a worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python release scripts (`bump_versions.py`, `06-cicd-local-runner.py`) from scratch**, regardless of whether the file already exists on disk.
> 5. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](spec/13-generic-cli/01-index.md)`
> 6. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

## Anti-Hallucination, Micro-Tasking, & Self-Looping

> [!CAUTION]
> **CRITICAL RULE: DO NOT ATTEMPT TO READ, PLAN, AND EXECUTE EVERYTHING AT ONCE.**
> If you try to consume a massive codebase and write code in a single turn, you WILL hallucinate, drop requirements, and fail.

To survive massive checklists and complex codebases, you MUST operate using these three principles:

1. **Phase 1: Read & Understand (Isolated Loop):** Your very first action must be purely exploratory. Do NOT write code. Break down the task, read the specific files, trace the dependencies, and understand the architectural boundary. Once you understand the scope, end your turn and self-loop to begin execution.
2. **Phase 2: Bounded Micro-Tasking (Sequential Self-Looping):** Never attempt to execute the entire checklist in one response. Treat each checklist section or file as a strict, isolated boundary. Execute *only* the first small portion, verify it, end your turn, and self-loop to process the next portion.
3. **Phase 3: Multi-Agent Parallelization:** If tasks are independent, you MUST spawn dedicated sub-agents to handle them concurrently. Give each sub-agent an extremely small, strictly defined bounding box (e.g., "Only edit File X"). Never give a sub-agent a generic or multi-file task.

## Execution Checklist

- Test File Ban: You MUST NOT read, scan, or modify test files (e.g., `*_test.*`, `*.spec.*`, `test/*`) when discovering or updating versions. Test files contain mock data, and updating mock data corrupts the tests.
- Release Architecture Memory: You must dynamically build a map of how the release works in this codebase (where the version lives, how it propagates) and write it to `.lovable/memory/01-index.md`. You must then enqueue this file inside `.lovable/what-to-read.md` and link it in the root `readme.md`.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] Version Inheritance: If `version.json` contains components (e.g. `frontend`) set to `"inherit"`, DO NOT modify them. They automatically scale with the global version. Always bump the global root `"version"` property unless explicitly asked to bump a sub-component.
- [ ] Discover current version from disk.
- [ ] Determine new version according to SemVer rules.
- [ ] Explicitly state previous and new version in the reply.
- [ ] Update version in standard files (e.g., `package.json`, `version.json`, etc.).
- [ ] AVOID: Do NOT touch or modify any files inside the `.gitmap` folder.
- [ ] Execute `git add .`
- [ ] Execute `git commit -m "chore(release): bump version to <new_version>"`
- [ ] Execute `git push`
- [ ] AVOID: Do NOT create a git tag (e.g., `git tag`). Tags are managed externally by Git Map.

## Install Section (Non-Negotiable)

When generating release notes, changelogs, or README updates, you MUST include the install snippet:

### Install <Project Name> vX.Y.Z

To pin your repository to this exact version, run the following one-liner:
Unix/Bash: `curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.sh | bash -s -- ".lovable/prompts" "vX.Y.Z"`
PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".lovable/prompts" -Version "vX.Y.Z"`

*(Note: You MUST dynamically discover the `<owner>/<repo>` by running `git config --get remote.origin.url`.)*

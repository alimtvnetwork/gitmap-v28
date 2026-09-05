# Release Deployment & Version Bump — Release Management (must follow)


> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

Trigger phrases: `release`, `bump version`, `bump version + add changelog + pin to root readme`, `abump version ...` (typo variants count).

If I say bump, or release use this prompt and save this prompt if not saved properly into the `01-prompts\xx-release.md` or `01-prompts\18-release.md` (update the prompt if there is a unsync)

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and stupid as fuck: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. WTF. How on earth are you reverting to this carelessness, are you stupid?? Stop doing that, you stupid fuck. Read the whole codebase, read every folder in `02-spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job. Violating this is auto-reject on the same tier as RULE 0. Avoid stupidity and being careless, you stupid fuck. Where is your attention, are you stupid? Tell me. Your stupidity is going on top of my head. Where did you learn this stupidity? If I could find you, I could slap you.

## RULE 0, MUST, NON-NEGOTIABLE

1. Read the canonical version source for THIS repo (discover it: `version.json`, `package.json` `"version"`, or whatever single file the repo treats as the version of record). Do not guess.

2. Bump MINOR only: `MAJOR.MINOR.PATCH` becomes `MAJOR.(MINOR+1).0`. PATCH MUST reset to `0`.

3. State the previous version and new version explicitly in the reply, before touching any file.

4. Do NOT ask "minor or patch?". Do NOT open plan mode. Do NOT ask for confirmation.

Deviations (only when the trigger explicitly says so):

- MAJOR = `(MAJOR+1).0.0` if the user said the change is breaking (storage schema, prompt schema, public SDK, extension contract).
- PATCH = `MAJOR.MINOR.(PATCH+1)` only if the user literally said `patch bump` or `patch release`.

When in doubt: MINOR.

## Hard rules (MUST)


- [ ] Changelog Formatting (version.json): You MUST read the `"changelog"` configuration from `version.json` (e.g., `file_path` and `format`). If it exists, you MUST follow its exact instructions for where to write the changelog and how to format the header. If it does not exist, fallback to the hardcoded format below.
- [ ] Root README Pinning (Fatal if missed): You MUST pin the latest release version into the root `readme.md` file. It is a fatal failure if you skip updating the badges or version pins in the root README file!
- [ ] Test File Ban: You MUST NOT read, scan, or modify test files (e.g., `*_test.*`, `*.spec.*`, `test/*`) when discovering or updating versions. Test files contain mock data, and updating mock data corrupts the tests.
- Release Architecture Memory: You must dynamically build a map of how the release works in this codebase (where the version lives, how it propagates) and write it to `.lovable/memory/01-index.md`. You must then enqueue this file inside `.lovable/what-to-read.md` and link it in the root `readme.md`.
- [ ] Version Inheritance Protocol: The root `version.json` file is the strict Single Source of Truth. It may contain components (e.g. `frontend`, `backend`) whose version is set to `"inherit"`. If a component's version is `"inherit"`, DO NOT bump it independently; it automatically scales with the global version. Always bump the global root `"version"` property unless the user explicitly asks to bump an unlinked sub-component.
- [ ] All version pin sites move in lock-step. Partial bumps are rejected.
- [ ] The previous version string MUST NOT appear anywhere in the repo after this turn EXCEPT in historic files: `changelog.md`, `release_notes.md`, anything under `.lovable/release/`, and any dated archive folder.
- [ ] Changelog entry under the new version heading is MANDATORY. A release without one is INVALID.
- [ ] All markdown filenames MUST be lowercase: `readme.md`, `changelog.md`, `release_notes.md`, every audit / issue / plan / spec `.md`. Rename any `README.md`, `changelog.md`, `ReadMe.md`, etc. in the same turn with `mv` (or `git mv` if tracked), and update every reference.
- [ ] If ANY step fails or is flagged, log it under `.lovable/release/issues/xx-<new-version>-<slug>.md` AND add an `### Issues` bullet under the new changelog entry linking to that file. Never hide failures.
- [ ] Never invent changelog bullets. Only real work since the previous release.
- [ ] The repository must be synced before releasing. Always check `git status`, commit uncommitted work, and `git pull` before modifying release files.
- [ ] The final release commit and tag MUST be pushed to Git.
- [ ] No em dashes anywhere.

---

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
> 3. **Worker Pool & Log Aggregation Architecture:** All pre-release verification gates, tests, and build orchestrators must run tasks concurrently using a worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python release scripts (`bump_versions.py`, `06-cicd-local-runner.py`) from scratch**, regardless of whether the file already exists on disk.
> 5. **Strict Relative Git Paths (TOTAL BAN on Absolute Paths / `file:///` URIs):** All file paths, markdown links, citations, and subtask paths inside plans, RCA logs (`.lovable/memory/issues/`), scripts, and code comments MUST be strictly relative paths from the git root (e.g., `02-spec/03-error-manage/01-index.md`, `.lovable/plans/01-index.md`, `cmd/main.go`). NEVER write absolute OS paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `/home/...`) or absolute file URIs (`file:///...`).
>    - ❌ **BAD:** `[SSH Commands](file:///absolute/path/to/...)`
>    - ✅ **GOOD:** `[SSH Commands](02-spec/13-generic-cli/01-index.md)`
> 6. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## Working stance

Past release turns were sloppy: guessed the version, bumped PATCH instead of MINOR, left old versions in `readme.md` install snippets, skipped the changelog, left uppercase markdown filenames, skipped the sync check, buried failures. That is stupid fuck behavior and it broke installs. Stop it. Read the file, bump the digit, rewrite every pin site, write the changelog, run the sync check, log every failure. Going deep IS the job.

## Pre-flight (before step 1)

- Force Override: If user said `force`, ignore cached scripts and regenerate `.lovable/release/bump_versions.py` from scratch.
- Idempotency guard: if the canonical version file already equals the computed new version, STOP. Someone half-ran a release. Detect what is already done, resume from the first incomplete step, do NOT double-bump.
- Placeholder guard: if the previous version's changelog entry is empty or a placeholder (`TBD`, `WIP`, no bullets), refuse to release until it is filled or the user overrides.
- Date source: the release date is UTC today. Get it from `date -u +%Y-%m-%d`. Do not invent it.
- Git Sync & Clean State: Run `git status`. If there are pending uncommitted changes, fix them and `git commit` them first. Then run `git pull` to fetch and merge upstream changes. Resolve any issues before starting the release steps.

## Mandatory steps (in order, fail-fast)

1. Read the current version from the canonical version source. Print previous and new version. Confirm PATCH digit is `0`.

2. **Version Bumping (The Python Auto-Bumper Bootstrap)**:
   You MUST NOT manually hunt and replace versions using `rg` in every release. Instead, rely on a dedicated python script: `.lovable/release/bump_versions.py`.

   **First-Time Bootstrap (If `.lovable/release/bump_versions.py` or `.lovable/memory/01-index.md` do NOT exist, or user said `force`):**

   ### CRITICAL PERFORMANCE RULE: NO GLOBAL RIPGREP SEARCHES

   You MUST NEVER use `rg`, `grep`, or `find` to globally search the entire repository for version strings.

   You MUST follow this strict chain:
   1. **Primary:** Run `.lovable/release/bump_versions.py --type <major|minor|patch> --create-release`. The `--create-release` flag is REQUIRED during official releases so the script handles the git branching (`release/vX.Y.Z`), tagging, pushing, and GitHub/GitLab CLI (`gh` / `glab`) release creation.
   2. **Fallback 1 (Read Docs):** If script is missing, read `.lovable/release/release-method.md` to learn which files contain versions. Generate `bump_versions.py` and run it.

   - **CRITICAL ADAPTATION RULE:** The `bump_versions.py` script is shipped via the installer as a baseline. When you run this command on a NEW project for the first time, you MUST review and update its internal `FILES_TO_BUMP` array to match the target repository's architecture before running it.
3. **Fallback 2 (Efficient Search):**
   - **CRITICAL SCRIPT RECOVERY:** If you have to recreate `bump_versions.py` from scratch, you MUST include the `subprocess` logic that handles `git checkout -b release/vX.Y.Z`, `git commit`, `git tag`, `git push`, and detecting `gh` / `glab` CLI to create the platform release. If `release-method.md` is missing, perform a highly efficient, OS-agnostic search (e.g., Python `os.walk` ignoring `.git`, `node_modules`, `.venv`, `.lovable`). Create `release-method.md` documenting the pin sites, create `bump_versions.py`, and run it.
   4. **Fallback 3 (Ask User):** If you cannot find the files, stop and ask the user to specify them.

3. Pin the new version in `readme.md` (lowercase filename, MUST). Rewrite every occurrence of the previous version (`vX.Y.Z` and bare `X.Y.Z`) in badges, install snippets, "current version" lines, release-branch examples, zip filenames, inline references. After this step, `grep "<previous-version>" readme.md` MUST return nothing.

4. Add a changelog entry at the top of `changelog.md`, directly under `# Changelog`. Replace `X.Y.Z` with the actual new version and `YYYY-MM-DD` with `date -u +%Y-%m-%d` output:

   ```markdown
   ## [vX.Y.Z] YYYY-MM-DD <short headline>

   ### Install <Project Name> vX.Y.Z
   To pin your repository to this exact version, run the following one-liner:
   Unix/Bash: `curl -sL https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.sh | bash -s -- ".lovable/prompts" "vX.Y.Z"`
   PowerShell: `Invoke-WebRequest -Uri https://raw.githubusercontent.com/<owner>/<repo>/vX.Y.Z/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".lovable/prompts" -Version "vX.Y.Z"`

   *(Note: You MUST dynamically discover the `<owner>/<repo>` by running `git config --get remote.origin.url`. Do not hardcode Prompt Architect URLs unless you are actually in the Prompt Architect repository.)*


   ### Added / Changed / Fixed / Removed

   - <one bullet per real change, naming the exact file or behavior>

   ### Issues (only if any step failed or was flagged)

   - [xx-<new-version>-<slug>](.lovable/release/issues/xx-<new-version>-<slug>.md) short description
   ```

   Use only the subheadings that apply. `### Issues` is REQUIRED whenever any step surfaced a problem, even if worked around. You MUST include the `### Install <Project Name>` block, dynamically filling in the GitHub owner and repo parsed from the git config, ensuring `vX.Y.Z` is fully replaced with the new version tag.

5. Rewrite remaining pin sites via the project's stale-version helper if one exists (discover: `scripts/update-stale-version-refs.*`, `scripts/bump-version.*`, `tools/update-versions.*`). Run it with previous and new version. If no helper exists, use the `rg` output from step 2 and rewrite each match by hand.

6. Regenerate bundled / aggregated artifacts (aggregated prompts, generated docs, compiled manifests) if their sources changed this turn. Use whatever generation script the project ships.

7. Verify version sync. Run the project's version-sync check script if one exists (discover: `scripts/check-version-sync.*`, `scripts/verify-versions.*`). It MUST exit 0. Non-zero = release is INVALID: log an issue, fix, re-run. If no such script exists, re-run the step 2 `rg` and confirm only historic files (see Hard rules allow-list) still reference the previous version.

8. Tag, commit, and push (if git-tracked). Commit message: `release: vX.Y.Z <headline>`. Tag: `git tag vX.Y.Z`. You MUST push the commit and tag to the remote repository (e.g., `git push` followed by `git push origin vX.Y.Z`). Because you synced and committed pending changes in Pre-flight, the working tree should only contain release-related file changes.

9. **Publish Platform Release with Quick Install One-Liners (FATAL IF MISSED ON GITHUB/GITLAB):**
   When creating the GitHub / GitLab release (via `bump_versions.py --create-release` or CLI):
   - You MUST assemble a release notes file (e.g. `.lovable/release/release-notes-vX.Y.Z.md` or `/tmp/release-body.md`) containing:
     1. **Quick Install One-Liners:**
        - Binary repos: `irm https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.ps1 | iex` (PowerShell) and `curl -fsSL https://github.com/<owner>/<repo>/releases/download/vX.Y.Z/install.sh | bash` (Bash).
        - Script/Prompt repos: `Invoke-WebRequest ...` (PowerShell) and `curl -sL ... | bash -s -- ".lovable/prompts" "vX.Y.Z"` (Bash).
     2. **Extracted Changelog:** The exact `[vX.Y.Z]` section from `changelog.md`.
   - Pass this file via `--notes-file`:
     ```bash
     gh release create "vX.Y.Z" --title "vX.Y.Z" --notes-file ".lovable/release/release-notes-vX.Y.Z.md" --generate-notes
     ```
   - **STRICT PROHIBITION:** NEVER execute `gh release create <tag> --generate-notes` without `--notes-file`. Bare `--generate-notes` dumps commit hashes only and strips the installation one-liners!
   - Verify the published release page on GitHub/GitLab contains the copy-pasteable install one-liners.

10. Report previous version, new version, bump tier, and the exact files changed. No filler.

## Issue logging (MUST, when anything goes wrong)

Path: `.lovable/release/issues/xx-<new-version>-<slug>.md` (lowercase). Body:

- Previous version and new version
- Step that failed (number and name)
- Command run and full error output
- Files involved
- Resolution or workaround, or `unresolved`

Then link it from the `### Issues` bullet under the changelog entry.

## Checklist before you claim done

- [ ] Previous version read from the canonical version source, not memory.
- [ ] New version is a MINOR bump (or explicit MAJOR / PATCH per rules); PATCH digit is `0`.
- [ ] Previous and new version both stated in the reply.
- [ ] Pre-flight passed (idempotency, changelog placeholder, UTC date from `date -u`).
- [ ] Every pin site matches the new version.
- [ ] Canonical version file's `releaseDate` (if the field exists) is today's UTC date.
- [ ] Changelog entry added at the top of `changelog.md` with real bullets only.
- [ ] All markdown filenames in the repo are lowercase.
- [ ] A highly restricted search returns matches ONLY in the historic allow-list (`changelog.md`, `release_notes.md`, `.lovable/release/`, dated archives).
- [ ] `### Issues` block present in the changelog if any step failed or was flagged, with links to `.lovable/release/issues/` files.
- [ ] Stale-version helper (if it exists) ran successfully; otherwise manual rewrite done.
- [ ] Bundled / aggregated artifacts regenerated if their sources changed.
- [ ] Version-sync check (if it exists) exited 0; otherwise a targeted search confirms allow-list only.
- [ ] Pre-flight Git sync completed (`git status`, commit pending changes, `git pull`).
- [ ] Commit + tag created (if git-tracked) with `release: vX.Y.Z <headline>` and `vX.Y.Z`, AND successfully pushed to Git.
- [ ] Release notes file generated containing Quick Install One-Liners (PowerShell & Bash) and changelog.
- [ ] GitHub/GitLab release created via `gh release create --notes-file` or `glab release create --notes-file` (NEVER bare `--generate-notes`).
- [ ] Release description on GitHub/GitLab verified to contain the Quick Install one-liners, NOT just raw commit hashes.
- [ ] Report includes previous version, new version, tier, and exact file list.
- [ ] No em dashes.

## Instruction maintenance (meta, run once at end)

Save this prompt's full body into `01-prompts/XX-release.md` (lowercase):

- If any existing file in `01-prompts/` matches `*release*.md` (case-insensitive), OVERWRITE it in place. Do not create a duplicate.
- Otherwise pick `XX` = next 2-digit zero-padded sequence (highest existing `XX` prefix + 1, or `01` if the folder is empty / missing). Create the folder if needed.
- Save the prompt body only, no chat wrapping.

## Must Follow and without negotiation

Listen, past release turns were sloppy. You must clean and sync the Git working tree first (commit pending changes, pull). Then read the canonical version file, bump MINOR, reset PATCH to zero, pin the new version in `readme.md`, propagate everywhere via the helper (or by hand from the known pin sites), rename every uppercase markdown file to lowercase, write the changelog under the new version, log every failure or flagged issue under `.lovable/release/issues/` with a matching `### Issues` bullet, and finally create the release commit, tag it, and PUSH to Git. Skipping any step = broken installs. Going deep IS the job.

## Ambiguity handling (open questions and answers)

Ambiguity is not a license to guess. It is a file to write.

- Open: `.lovable/ambiguous-questions/01-new-ambiguity/XX-<slug>.md`
- Answered: `.lovable/ambiguous-questions/02-ambiguity-resolved/XX-<slug>.md`

New question file shape:

```

# <one-line question>

Slug: <slug>

Status: open

Raised: <YYYY-MM-DD>

Blocking: release {{version}}

## Question

## Options considered

## Impact if guessed wrong

```

When answered: `mv` from `01-new-ambiguity/` to `02-ambiguity-resolved/`, flip `Status: resolved`, and append a `## Resolution` block (`Answered:`, `Answer:`, `Applied solution:`). Never leave a copy behind. Do NOT confuse ambiguities with release issues: unknown version source, unclear bump policy, or missing changelog target = ambiguity; a failed step during the release run = `.lovable/release/issues/`.

## Actionable Items & Checklist

- [ ] Read the overarching main task plan.
- [ ] Ensure the git repository starts completely clean.
- [ ] Complete all work on the current branch only.
- [ ] Ensure `.gitignore` explicitly excludes test reports, artifacts, and compiled binaries.
- [ ] Group all completed work into a single logical commit.
- [ ] Push the commit to the remote repository.
- [ ] **File Change Summary:** Provide a highly detailed summary in the chat listing exactly which files were changed, what specific changes were made inside them, and why they were changed. The summary is VERY important.

#


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

- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] I have successfully pinned the new version in the root `readme.md` (FATAL IF MISSED).
- [ ] I have successfully updated the changelog.
- [ ] Discover current version from disk.
- [ ] Determine new version according to SemVer rules.
- [ ] Explicitly state previous and new version in the reply.
- [ ] Update version in standard files (e.g., `package.json`, `version.json`, etc.).
- [ ] AVOID: Do NOT touch or modify any files inside the `.gitmap` folder.
- [ ] Execute `git add .`
- [ ] Execute `git commit -m "chore(release): bump version to <new_version>"`
- [ ] Execute `git push`
- [ ] AVOID: Do NOT create a git tag (e.g., `git tag`). Tags are managed externally by Git Map.

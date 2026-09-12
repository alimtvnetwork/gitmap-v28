# Strictly Avoid

Items in this file MUST NEVER be suggested, recommended, asked about, or built again.

## 02-spec/19-main-worker-service implementation — TOTAL BAN

🔴 **NEVER write, scaffold, propose, or suggest implementation code for `02-spec/19-main-worker-service/` (the Main-Worker Service).**

This repo is **spec-only** for Spec/19. Allowed work:
- ✅ Authoring / editing markdown under `02-spec/19-main-worker-service/**`
- ✅ Audits, consistency reports, changelogs, diagrams, glossary
- ✅ Cross-spec references that *describe* the worker

Forbidden:
- ❌ Any Go / Rust / TypeScript / PowerShell / shell source files implementing the worker
- ❌ Scaffolding service binaries, DB migrations, REST handlers, JWT/JID flows, push update logic for Spec/19
- ❌ "Phase 1 implementation", "begin coding", "ship the service", "starter skeleton" suggestions
- ❌ Asking the user whether they want to start implementing Spec/19
- ❌ Listing Spec/19 implementation as an "optional next-phase candidate" or follow-up

**If a `next` command would otherwise land on Spec/19 implementation, skip it and propose spec-level work instead.**

**Why:** User explicitly stated this repo only writes the spec — implementation belongs elsewhere. Re-suggesting it is a hard failure.

---

## readme.txt timestamp generator — TOTAL BAN

🔴 **NEVER build, suggest, propose, design, spec, or even mention any feature that writes a timestamp / date / time / "Malaysia-formatted" content into `readme.txt` (or any other file).**

This includes — but is not limited to:
- ❌ A `refresh-readme.ps1` / `refresh-readme.sh` / any script that writes time into readme.txt
- ❌ A `readme` sub-command on `run.ps1` / `run.sh` that touches readme.txt timestamps
- ❌ An npm script (`npm run refresh-readme`, etc.) that writes time into readme.txt
- ❌ Hooking timestamp-writing into `npm run sync` or any other workflow
- ❌ Hard-coded prefix variants (`let's start now`, etc.), configurable prefix, curated lists, random phrases
- ❌ Any timezone discussion (Asia/Kuala_Lumpur, UTC, local) tied to readme.txt
- ❌ Any 12-hr / 24-hr / `dd-MMM-yyyy` format discussion tied to readme.txt
- ❌ Any idempotency variant (always rewrite, skip if same day, write if missing)
- ❌ "Instructions" / "how it works" / "how to run" / "how to test" sections about any such generator
- ❌ Asking clarifying questions about any of the above
- ❌ Offering it as a follow-up, alternative, or "while we're at it" suggestion

**If the user asks for this feature again, do nothing except acknowledge that this entry forbids it. Do not negotiate. Do not propose a "smaller" version. Do not ask "did you mean X". Just stop.**

The only acceptable interaction with `readme.txt` is a one-shot manual edit when the user explicitly types the exact content they want in that turn.

**Why:** User has rejected this feature, the suggestion of this feature, the discussion of this feature, and the documentation of this feature multiple times across sessions, with escalating frustration. Re-raising it is a hard failure.

---

## Absolute File System Paths

🔴 **NEVER use `file:///` absolute paths in markdown files, artifacts, or code.**

Everything must be standalone relative to the repo root (`/`). See: `.lovable/memory/avoid/03-absolute-file-system-paths.md`

---

## Committing Generated Artifacts and Test Reports — TOTAL BAN

🔴 **NEVER commit test results, test reports, temporary test data, or compiled binaries.**

This includes — but is not limited to:
- ❌ `.test-report.*`, HTML coverage reports, or JSON test outputs
- ❌ `.exe`, `.dll`, `.so`, `.class`, `.out`, or any compiled binary
- ❌ Committing the `build/`, `bin/`, `obj/`, or `dist/` folder unless explicitly permitted by a deploy spec.

**If you generate these files during a run or compilation, verify they are ignored by `.gitignore`. If not, update `.gitignore` or delete them before running `git add`.**

## Release on Every Commit — TOTAL BAN

🔴 **NEVER trigger a release (version bump, release tagging, `scripts/release.mjs`) on every commit or every chat turn.**

Forbidden:
- ❌ Running `scripts/release.mjs` or `npm run release` for standard tasks, documentation updates, bug fixes, or minor features.
- ❌ Creating `release: vX.Y.Z` commits and `git tag vX.Y.Z` on every single conversation turn.
- ❌ Treating the end of an AI turn as the "End of Tunnel Release" unless the user EXPLICITLY commands a version release.

Allowed work:
- ✅ Standard semantic commits (`feat: ...`, `fix: ...`, `docs: ...`, `chore: ...`) for all work.
- ✅ Pushing standard commits to the branch (`git push`).
- ✅ Executing a release **ONLY** when the user explicitly says "release", "bump version", or something explicitly confirming a version bump is needed for distribution.

**Why:** Releasing on every commit completely pollutes the git history and version tags.

---

## Explicit `== true` Checks — TOTAL BAN

🔴 **NEVER evaluate boolean variables explicitly against `true` (e.g., `== true` or `=== true`).**

Forbidden:
- ❌ `if isValid == true {`
- ❌ `if (hasMatch === true) {`
- ❌ `return isSuccess == true`

Allowed work:
- ✅ Implicit positive checks: `if isValid {`
- ✅ Implicit positive checks: `if (hasMatch) {`
- ✅ Returning directly: `return isSuccess`
- ✅ Using `== false` or `=== false` as a replacement for the banned `!` operator (if permitted by the language's specific guideline).

**Why:** Implicit boolean evaluation is a universal standard. The AI incorrectly generalized the `=== false` rule into `=== true`. `true` is redundant and prohibited.

## British English Spelling — TOTAL BAN

🔴 **NEVER use British English spelling (e.g., `behavior`, `recognise`) in the codebase.**

Forbidden:
- ❌ `behavior`
- ❌ `recognise`
- ❌ `colour`, `initialise`

Allowed work:
- ✅ US English spelling: `behavior`, `recognize`, `color`, `initialize`

**Why:** The codebase strictly enforces US English to pass the misspell linter and ensure global consistency.

## Disabling Linters or CI/CD Checks — TOTAL BAN

🔴 **NEVER modify linter configurations or CI/CD pipelines to bypass or disable failing checks.**

Forbidden:
- ❌ Changing `enable:` to `disable:` in `.golangci.yml` or removing strict rules (`gosec`, `revive`, `misspell`, etc.).
- ❌ Modifying `.eslintrc`, `eslint.config.js`, `.prettierrc`, or `ruff.toml` to turn off failing rules.
- ❌ Adding `//nolint`, `@ts-ignore`, or `eslint-disable` globally or excessively to bypass errors.
- ❌ Modifying GitHub Actions/GitLab CI `.yml` files to skip steps or ignore failures.

Allowed work:
- ✅ Fixing the actual source code that is violating the linter rule.
- ✅ Modifying linter configurations ONLY if the user explicitly commands you to "configure the linter" or "add this rule to golangci".

**Why:** When instructed to "fix CI errors," the AI sometimes takes the lazy route of disabling the linter rather than fixing the code. This defeats the entire purpose of quality gates and is strictly prohibited.

## Golang Underscore Variable Naming — TOTAL BAN

🔴 **NEVER use underscores (`snake_case`) for variable, struct, or function names in Golang.**

Forbidden:
- ❌ `user_id`, `has_error`, `api_key`
- ❌ `type user_model struct`

Allowed work:
- ✅ `userId` or `userID` (camelCase for variables)
- ✅ `hasError`, `apiKey`
- ✅ `type UserModel struct` (PascalCase for structs/exported types)

**Why:** Go conventions explicitly dictate `camelCase` or `PascalCase`. Underscores violate the language's core style guide and will fail standard linters.

## Modifying Version Information — TOTAL BAN

🔴 **NEVER manually edit `version.json`, `package.json` version strings, or changelog dates.**

Forbidden:
- ❌ Modifying `version.json` manually to bump a version or change an `updated` date.
- ❌ Updating the `version` field in `package.json` by hand.

Allowed work:
- ✅ Allowing the `scripts/release.mjs` (or similar node scripts) to manage versions and dates.
- ✅ Inspecting the `.git/config` or running `git remote -v` if you need to verify repository information.

**Why:** Version information is strictly managed by its own synchronization scripts and source-of-truth repositories. Manual AI edits cause synchronization drift and pipeline failures.

### Strict Relative Git Paths (Zero Tolerance)

Absolute filesystem paths (e.g., `/absolute/path/to/...`, `/Users/.../`, `/home/...`) and absolute file URI schemes (`file:///absolute/path/to/...`, `file:///absolute/path/to/`) are **strictly forbidden** inside committed repository files, specifications, markdown plans, subtask files, code comments, and citations. All paths must be relative to the git repository root.

### No `spec/` Inside `.lovable/` (Total Ban on `.lovable/spec/`)

🔴 **NEVER create or store specifications inside `.lovable/spec/`.**
- All canonical specifications must live under the root `spec/` directory.
- All repo-specific / application-specific specifications must reside under `02-spec/21-app/`.
- The `.lovable/` directory is reserved exclusively for AI metadata (`memory/`, `plans/`, `prompts/`, `ai-fix-scripts/`, `assets/`, `procedures/`, `suggestions/`, `question-and-ambiguity/`).

---

## Go Interface Suffix — TOTAL BAN

🔴 **NEVER suffix Go interfaces with `Interface` (e.g., `WriterInterface`, `StreamerInterface`).**

Forbidden:
- ❌ `type WriterInterface[T any] interface`
- ❌ `type StreamerInterface[T any] interface`
- ❌ `type HandlerInterface interface`

Allowed work:
- ✅ Idiomatic Go `-er` interfaces: `type Writer[T any] interface`, `type Streamer[T any] interface`, `type Reader interface`, `type Formatter[T any] interface`.

**Why:** Go conventions mandate concise, idiomatic `-er` naming for single- or few-method interfaces representing behavior. Suffixing with `Interface` is an anti-pattern imported from other languages and strictly prohibited in this repository.

---

## Uppercase ID Acronym in Identifiers — TOTAL BAN

🔴 **NEVER use all-caps `ID` in variable names, struct fields, method names, or function parameters.**

Forbidden:
- ❌ `UserID`, `OrderID`, `AccountID`, `TraceID`, `ID`
- ❌ `GetID()`, `SetID()`, `traceID`

Allowed work:
- ✅ PascalCase `Id`: `UserId`, `OrderId`, `AccountId`, `TraceId`, `Id`
- ✅ camelCase `id`: `userId`, `orderId`, `accountId`, `traceId`, `id`

**Why:** Acronym casing must be normalized to `Id` in PascalCase and `id` in camelCase across all languages to eliminate capitalization inconsistencies and pass repository naming linters.

---

## Boolean Fields Without Positive Prefixes — TOTAL BAN

🔴 **NEVER define boolean fields, variables, or properties without an explicit positive prefix (`is`, `has`, `should`, `can`).**

Forbidden:
- ❌ `Active bool`, `Success bool`, `Match bool`, `Ready bool`
- ❌ `active: boolean`, `success: boolean`

Allowed work:
- ✅ `IsActive bool`, `IsSuccess bool`, `HasMatch bool`, `IsReady bool`
- ✅ `isActive: boolean`, `isSuccess: boolean`

**Why:** Bare boolean identifiers violate the repository's positive-polarity naming convention and impair readability in conditional guard clauses.

---

## Noisy Passing Quality Gate Output in CI Runners — TOTAL BAN

🔴 **NEVER flood developer terminals with stdout/stderr logs from passing quality gates in local CI test runners.**

Forbidden:
- ❌ Dumping passing command outputs to the console when all gates succeed.
- ❌ Interleaving asynchronous stdout streams from parallel workers across terminal lines.

Allowed work:
- ✅ Real-time single-line status ticker for completion progress (`[ 1/21] ✅ [PASS] <Gate> (<duration>s)`).
- ✅ Selective log suppression: print stdout/stderr ONLY for gates that exit with a non-zero status code or timeout.
- ✅ Full verbose logs emitted ONLY when the user explicitly passes the `--all` (`-a`) flag.

---

## Modifying Chrome Local State Without Running Process Guards & Incomplete Attribute Schemas — TOTAL BAN

🔴 **NEVER modify Chrome's `Local State` without detecting active Chrome processes or write incomplete Chromium profile schemas.**

Forbidden:
- ❌ Mutating `%LOCALAPPDATA%\Google\Chrome\User Data\Local State` without checking `isChromeRunning()`. Active Chrome processes overwrite in-memory cache to disk, destroying external changes.
- ❌ Registering profiles in `profile.info_cache` with partial schemas (only `name` and `user_name`), which causes Chrome's profile picker UI to discard the tile during startup validation.
- ❌ Wiping `account_info` from profile `Preferences` during snapshot imports, destroying user email bindings.

Allowed work:
- ✅ Check `isChromeRunning()` before modifying browser state, issuing clear user notices.
- ✅ Populate all 13 Chromium attributes (`name`, `shortcut_name`, `user_name`, `avatar_icon`, `default_avatar_fill_color`, `default_avatar_stroke_color`, `profile_highlight_color`, `profile_color_seed`, `active_time`, `is_using_default_avatar`, `is_using_default_name`, `is_ephemeral`, `is_consented_primary_account`, `signin.with_credential_provider`).
- ✅ Maintain `profile.profiles_order` synchronicity and provide `gitmap chrome profile reconcile` for automated orphan profile restoration.

---

## Bare Platform-Specific `open` in Macro Steps Without Cross-Platform Shims — TOTAL BAN

🔴 **NEVER execute raw platform-specific commands like `open` in macro runners without cross-platform shimming or OS-aware dispatching.**

Forbidden:
- ❌ Assuming `open` is a universal shell executable across Windows and Linux.
- ❌ Spawning unshimmed `open <target>` on Windows where `open` does not exist, causing exit status 1 failures (`E9000:EXECUTION`).

Allowed work:
- ✅ Provide internal shimming in the macro execution engine: translate `open` on Windows to `explorer.exe`, `cmd.exe /c start ""`, or `rundll32 url.dll,FileProtocolHandler`, and on Linux to `xdg-open`.
- ✅ Delegate URL/browser operations to GitMap's built-in commands (`gitmap open`, `gitmap chrome open`).

---

## Executing Git Command Chains via Shell String Concatenation (`cmd /c`) on Windows — TOTAL BAN

🔴 **NEVER execute composite Git command chains using shell string concatenation (`cmd /c`, `sh -c`) without structured argument isolation.**

Forbidden:
- ❌ Running `exec.Command("cmd", "/c", "git -C <path> commit -m \"...\"")`. Windows `cmd.exe` strips or mishandles double quotes around commit messages, causing Git to mistake message words for file pathspecs (`error: pathspec '<word>' did not match any file(s) known to git`).
- ❌ Relying on shell quote-escaping for compound operations (`&&`, `;`) containing dynamic commit messages or file paths with spaces.
- ❌ Compiling only `bin/gitmap.exe` while neglecting `C:\Users\Alim\AppData\Local\gitmap-cli\gitmap.exe` and root `gitmap.exe`, causing user terminals to invoke stale binaries through PowerShell wrappers.

Allowed work:
- ✅ Execute discrete command steps natively via `RemediationStep{Name: "git", Args: []string{"-C", repoPath, "commit", "-m", msg}}` with `exec.Command(step.Name, step.Args...)`.
- ✅ Synthesize structured steps dynamically when recipes are loaded from legacy serializations missing `steps`.
- ✅ Synchronize all 4 Gitmap execution paths upon compilation (`bin/gitmap.exe`, root `gitmap.exe`, `AppData\Local\gitmap-cli\gitmap.exe`, and `AppData\Local\gitmap\gitmap.exe`).

---

## Unverified VMware Mount Invocations Without Pre-Mount Cleanup or Diagnostics — TOTAL BAN

🔴 **NEVER execute `vmhgfs-fuse` or VMware mount commands without pre-mount stale cleanup, service readiness checks, fallback mounting, and Error -107 diagnostics.**

Forbidden:
- ❌ Running `vmhgfs-fuse` on a mount point without unmounting stale or broken FUSE endpoints first (`fusermount -u` / `umount -l`).
- ❌ Failing silently or returning raw `Error -107 cannot open connection!` without explaining that VMware host Shared Folders is disabled in VM settings.
- ❌ Installing only `open-vm-tools` on Linux guests while omitting `open-vm-tools-desktop`.

Allowed work:
- ✅ Always unmount stale FUSE mount points before mounting.
- ✅ Ensure `open-vm-tools` service is running via `systemctl start open-vm-tools`.
- ✅ Attempt fallback mount `mount -t fuse.vmhgfs-fuse .host:/ <mountPoint> -o allow_other` if primary `vmhgfs-fuse` fails.
- ✅ Check `vmware-hgfsclient` and provide clear, actionable VM Settings -> Options -> Shared Folders remediation instructions when Error -107 occurs.
- ✅ Always install `open-vm-tools` and `open-vm-tools-desktop` together.

---

## Unshimmed Interactive Prompts Lacking PWD Awareness & Local File Inspection — TOTAL BAN

🔴 **NEVER force users into blind interactive CLI step prompts without displaying the current working directory (`PWD`) or allowing directory file inspection (`ls`/`dir`).**

Forbidden:
- ❌ Displaying a bare `Step N> ` prompt without showing where in the filesystem the command will run.
- ❌ Registering file inspection commands (`ls`, `dir`) as permanent macro steps when the user is trying to inspect their workspace.

Allowed work:
- ✅ Render a clean `[PWD: ...]` banner above prompt lines, with user toggles (`pwd on`/`pwd off`).
- ✅ Provide in-builder `ls`, `dir`, `find`, `search`, and `replace` commands to assist recipe construction.

---

## Strict Network Release Dependencies in Local CI/CD Smoke Tests — TOTAL BAN

🔴 **NEVER force local CI/CD runners to query remote GitHub release downloads for unreleased local tags without local mock server fallback.**

Forbidden:
- ❌ Looping 5 times with 10-second sleeps against GitHub for unpublished tags, stalling local test execution for 176+ seconds before hard failing.
- ❌ Bypassing or disabling installer tests in CI/CD instead of implementing local mock release packaging.

Allowed work:
- ✅ Automatically probe if the release asset exists on GitHub.
- ✅ When absent, spin up an offline Python HTTP server serving local snapshot binaries packaged into the expected `.zip` / `.tar.gz` and `checksums.txt` format via `GITMAP_DOWNLOAD_URL`.

---

## Depositing Binaries in Parent Workspaces or Outside Canonical Paths — TOTAL BAN

🔴 **NEVER deposit, build, or write executable binaries into workspace parent directories or anywhere outside canonical `bin/` and user AppData directories.**

Forbidden:
- ❌ Running `go build -o ../gitmap` or `go build -o ../gitmap.exe` resulting in stray binaries in parent workspaces (e.g. `../gitmap.exe`).
- ❌ Leaving unversioned compiled binaries scattered across repository directories or outside `bin/`.

Allowed work:
- ✅ Always direct build output strictly to `bin/gitmap` (or `bin/gitmap.exe`).
- ✅ User deployment targets must strictly be `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` or `/usr/local/bin/gitmap`.
- ✅ Ensure `make clean` purges `bin/gitmap` and `bin/gitmap.exe`.

---

## Running Tests Without Owner Explicit Command — TOTAL BAN

🔴 **NEVER run unit tests, test suites, or CI test jobs (`go test`, `pytest`, `npm test`, or test jobs in CI runners) during standard development tasks or prompt execution unless explicitly commanded by the repository owner.**

Forbidden:
- ❌ Executing `go test`, `pytest`, `npm test`, or `cargo test` during standard development tasks, prompt executions, coding guideline fixes, or refactoring loops.
- ❌ Running `python 03-ai-scripts/06-cicd-local-runner.py` without `--no-tests` during standard development; always pass `--no-tests`.
- ❌ Adding automatic test execution steps to non-release workflows or prompts.

Allowed work:
- ✅ Run tests when the repository owner explicitly requests it in their prompt (e.g., "run tests", "execute unit tests", "fix failing tests").
- ✅ **ALL CI/CD Fix Workflows (`ci-cd-fix`, `16-ci-cd/*`):** MUST run all unit test suites, integration tests, linters, and quality gates properly (`python 03-ai-scripts/06-cicd-local-runner.py`) to diagnose, surface, and repair pipeline failures. Skipping tests with `--no-tests` in CI/CD fix tasks is strictly prohibited.
- ✅ **ALL Release Workflows (`release-management`, `release-orchestrator`, `01`, `03`, `07`, `16-ci-cd/04`):** MUST run all unit test suites (`--run-tests`) and verify 100% green passing before cutting any release.
- ✅ Always use `--no-tests` (or `--skip-tests`) when running standard routine development quality gate checks (`06-cicd-local-runner.py`) unless running CI/CD fixes, release ceremonies, or explicitly instructed by the owner.

**Why:** Unit test suites can be slow, resource-heavy, and disruptive during rapid iterative development loops. Running tests without explicit owner authorization wastes resources. Quality gates in standard turns focus on static analysis, linting, and structural integrity.

---

## Test Inventory & Atomic Recent File Changes Locking Mandate

🔴 **NEVER record or modify recent file change logs without cross-platform atomic file locking, and NEVER bypass `.lovable/test-inventory.json`.**

Forbidden:
- ❌ Writing directly to `.lovable/temp/recent-file-changes.json` without acquiring `.lovable/temp/recent-file-changes.lock`.
- ❌ Failing to release the lock or failing to handle stale locks properly.
- ❌ Guessing or manually hard-coding test file relationships without checking `.lovable/test-inventory.json`.

Allowed work:
- ✅ Use `python 03-ai-scripts/33-test-inventory-generator.py --record <relative-path>...` to atomically record modified files and resolve associated tests.
- ✅ Maintain and synchronize `.lovable/test-inventory.json` when adding, moving, or deleting test files by running `python 03-ai-scripts/33-test-inventory-generator.py`.
- ✅ Ensure all recorded paths are distinct, lowercase, and strictly relative to the repository root.

**Why:** Concurrent multi-agent orchestration and asynchronous script runs will corrupt `recent-file-changes.json` if writes are uncoordinated. Centralized test inventory mapping guarantees reproducible test discovery when an authorized release or targeted test fix is executed.

---

## Running Full CI/CD Runner During Routine Development or Coding Guideline Turns — TOTAL BAN

🔴 **NEVER run `python 03-ai-scripts/06-cicd-local-runner.py` during routine development tasks, coding guideline fixes (`15-cg-execute/*`), micro-loops, or sub-agent turns.**

Forbidden:
- ❌ Running `python 03-ai-scripts/06-cicd-local-runner.py` (with or without `--no-tests`) during routine task execution loops, coding standard audits, single-file refactoring, or micro-batches.
- ❌ Re-running the heavy 28-38 gate pipeline repeatedly for routine edits, wasting minutes across unrelated files and packages.
- ❌ Using the full pipeline runner to verify a single guideline edit (e.g. nested-if or boolean condition) when a targeted linter is available.

Allowed work:
- ✅ Run targeted file-level linters/autofixers directly on the modified file(s) (e.g., `python linter-scripts/check-nested-ifs.py <file>`, `python 03-ai-scripts/08-naming-autofixer.py <file>`, `python linter-scripts/check-boolean-guidelines.py <file>`).
- ✅ **Owner Explicit Command:** Run the runner if and only if the repository owner explicitly requests running the pipeline.
- ✅ **CI/CD Fix Tasks (`ci-cd-fix`, `16-ci-cd/*`):** May run `python 03-ai-scripts/06-cicd-local-runner.py` because the primary goal of those tasks is specifically repairing CI/CD infrastructure.
- ✅ **Release Ceremonies (`release-orchestrator`, `01`, `03`, `07`, `16-ci-cd/04`):** Run `python 03-ai-scripts/06-cicd-local-runner.py --run-tests` as the mandatory final pre-release gate before cutting a release.

**Why:** The local CI/CD runner runs up to 38 segments (linters, cross-OS compilation, snapshot builds, web builds) across the entire codebase. Executing this massive suite on every micro-turn or coding guideline edit causes immense latency, hits unrelated files, and wastes substantial developer and compute time.

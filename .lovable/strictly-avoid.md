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


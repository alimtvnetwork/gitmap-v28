# 🔴 CODE RED: No Generated Code, Artifacts, or Test Results (AI Execution Prompt)

> **/goal** Prevent repository bloat, secret leakage, and diff pollution by strictly prohibiting the git staging and commitment of generated code, build binaries, test reports, and temporary test artifacts.
> **/learn** Master clean repository hygiene, namespace temporary host execution within `<temp_dir>/gitmap/<category>/`, isolate local artifacts exclusively inside `.ai-memory/temp/`, enforce mandatory pre-build cleanup, and maintain comprehensive `.gitignore` coverage.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify that no compiled binaries (`.exe`, `.so`, `.dylib`), build folders (`dist/`, `bin/`, `obj/`), or test reports (`.test-report.*`, coverage data) are staged or committed.
- [ ] `/learn` Never run `git add .` blindly; actively audit staged paths to ensure generated files and test outputs are safely ignored or pruned.
- [ ] `/goal` Restrict all host OS temp operations to `<temp_dir>/gitmap/<category>/` and in-repository temporary storage to `.ai-memory/temp/`.
- [ ] `/learn` Perform mandatory pre-build cleanup before compiling code to prevent disk exhaustion and prevent artifact collisions.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 🚫 The Rule: NEVER Commit Generated Code or Artifacts

Under no circumstances should any generated code (e.g., ORM models, gRPC stubs, OpenAPI clients), compiled binary, test result, test report, or temporary test data ever be committed to the Git repository.

### What is Forbidden?

- **Generated Code:** `*.generated.*` or `*_generated.*` files produced by code generators (unless explicitly permitted as test fixtures).
- **Test Results & Reports:** Any outputs from test runs, such as HTML reports, coverage data, JSON summaries, `.test-report.*`, or CSV results.
- **Compiled Binaries & Executables:** `.exe`, `.dll`, `.so`, `.dylib`, `.out`, `.class`, or any pre-compiled binaries resulting from a build step.
- **Build Directories:** `build/`, `bin/`, `obj/`, `dist/` (unless explicitly allowed for a specific deployment pipeline in a separate spec).
- **Temporary Data:** Any mock data, database dumps, or local logs generated during a local run.

### Why?

1. **Repository Bloat:** Binary files and continuous test reports exponentially increase the Git repository size, slowing down clones and pulls for all team members and CI/CD pipelines.
2. **Noise in History:** Committing transient data pollutes the Git history and diffs, making it impossible to perform meaningful code reviews.
3. **Security Risks:** Test results or compiled binaries can accidentally leak sensitive environment data, secrets, or internal architectural layouts.

---

## AI Agent Directives

> 🛑 **MANDATORY INSTRUCTION FOR AI:**
>
> If you are instructed to run a test suite, compile a binary, or execute a script that produces an output file, you MUST NOT commit those outputs to the Git repository.
>
> If you notice that an action you performed has generated a new test report or binary, you must ensure that it is either ignored via `.gitignore` or deleted before you stage files using `git add`.
>
> **Never blindly run `git add .` without verifying that generated artifacts are safely ignored.**

## Enforcement

This rule is enforced globally via `.gitignore` patterns. If a new type of artifact is introduced, you must update `.gitignore` before committing anything else.

---

## 🧹 Storage Hygiene, Temp Directory Isolation & Pre-Build Cleanup

1. **Repository-Scoped OS Temp**: Any execution requiring the host OS/user temporary directory (`os.TempDir()`, `tempfile.gettempdir()`, `$env:TEMP`, `$TMPDIR`) must strictly namespace operations within `<temp_dir>/gitmap/<category>/` (e.g. `build/`, `test/`, `purge/`, `downloads/`). Un-namespaced files or loose directories in the root of OS temp are strictly forbidden.
2. **Mandatory Pre-Build Cleanup**: Before executing compilation (`go build`, `npm run build`), previous build artifacts in the target build directory must be deleted. Storage reuse must be strictly maintained to prevent disk exhaustion.
3. **Workspace Isolation**: In-repository temporary artifacts must reside exclusively inside `.ai-memory/temp/`. Creating `.tmp/` at the repository root is banned.
4. **GitHub Actions Zero Storage (Total Ban on `upload-artifact` in CI)**: Never use `actions/upload-artifact` in CI pipelines. Free accounts are capped at 0.5 GB (500 MB) across the account. CI builds are ephemeral compilation verifications; binary artifacts belong exclusively in GitHub Releases via `release.yml`, never in Actions artifact storage. See `02-spec/02-coding-guidelines/01-cross-language/30-actions-zero-storage.md`.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ARCH-029: Zero Generated Code, Test Artifacts, and Binary Staging

**Given** Development workspaces, build pipelines, and automated test runners producing transient outputs.
**When** Repository git staging, file hygiene audits, and CI/CD preflight linters scan the repository.
**Then** Zero compiled binaries, test summaries, coverage reports, un-namespaced temp files, or generated code stubs are tracked or committed, achieving exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.

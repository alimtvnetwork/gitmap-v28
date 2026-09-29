# Release-Triggered CI/CD Fix Tweak with Targeted Smart Testing — Workflow (must follow)

Trigger Keywords & Aliases: `cicd fix release tweak`, `ci release tweak`, `fix and release tweak`, `smart release fix`

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

- [A] Base Workflow: `01-prompts/16-ci-cd/06-ci-cd-fix-with-release.md` (or V2: `01-prompts/16-ci-cd/09-ci-cd-fix-with-release-n-steps.md`)
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Release Script: `03-ai-scripts/29-release-orchestrator.py`
- [E] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

---

## Smart Targeted Fix & Release Execution (Fastest Path)

Follow the release-triggered RCA structure from [A] and review past failures in [B]. The primary goal is the fastest possible resolution of the failing stack trace followed by automated release ceremony without delays.

### Mandatory Smart Testing Rules

1. **Top-Instruction Precedence:** Any custom instruction provided before this section or in chat overrides default procedures.
2. **Stack Trace Targeting:** Isolate and build/test ONLY the packages and test functions directly cited in the provided failure or error stack trace (extract bounded lines via `gitmap pipeline error-logs` / `gitmap pe`).
3. **Priority Incremental Runner:** Run `python [E] run-smart` (or alias `--smart`, `-s`) to build only changed Go packages into OS temp and run Quad Runner.
4. **Changed Packages from Last Git Hash:** Compare changes against the last known git hash (`git diff --name-only HEAD~1` or `git status --porcelain`). Re-run tests ONLY for packages (Go, TS, Rust, Python) that contain actual modifications using `python [E] --changed-only` or `python [E] --pkg <target>`.
5. **Change Tracking & State Persistence:** Every time a fix is applied, write the modified file list and state to [C] (or run `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`) so the build system knows exactly which targets require re-verification.
6. **Strict Ban on Extraneous Runs:** NEVER run the complete test suite, global linters, spellcheckers, or unrelated packages during debugging. Verify strictly using `python [E] run-smart`, `python [E] --pkg <affected_package>`, or `python [E] --changed-only` with optional `--fast` heatmap filtering.
7. **Mandatory Same-Turn Tool Chaining:** Do not end the turn after outputting the diagnostic plan; invoke the diagnostic or fix tool in the exact same turn.
8. **Total Ban on Interim Per-File Commits:** Never commit individual files during debugging. All changes are accumulated and committed atomically by the release orchestrator.
9. **Issue & RCA Destination Routing:** Store CI/CD issue post-mortems in `.ai-memory/cicd-issues/` (indexed in `.ai-memory/cicd-index.md`). If an issue is an application bug rather than a CI/CD failure, document it in `02-spec/22-app-issues/` (indexed in `02-spec/22-app-issues/readme.md`).
10. **Targeted Release Verification:** Once the isolated fix passes green, execute [D] (`python [D]`) to finalize the automated release ceremony.

---

## Execution Prompt (Copy & Paste Trigger)

```markdown
- [A] Base Workflow: `01-prompts/16-ci-cd/06-ci-cd-fix-with-release.md`
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Release Script: `03-ai-scripts/29-release-orchestrator.py`
- [E] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

Follow workflow [A] and check past post-mortems in [B]. Perform a grounded 4-part RCA on the failure below (storing CI/CD RCAs in [B], or application bug RCAs in 02-spec/22-app-issues/), isolate the broken area, and apply a surgical fix.
SMART TEST & RELEASE: Do NOT run the full test suite during debugging—build and test ONLY the specific packages/files affected by the stack trace and files changed from the last git hash. Persist changes to [C] and verify targeted packages via `python [E] run-smart`, `python [E] --changed-only`, or `python [E] --pkg <target>` (with optional `--fast` heatmap filtering). Once green, execute release via [D]. Do not commit file-by-file; release orchestrator handles atomic commit and push.

<paste pipeline error / stack trace here>
```

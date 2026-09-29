# CI/CD Fix Tweak with Targeted Smart Testing & RCA — Workflow (must follow)

Trigger Keywords & Aliases: `cicd fix tweak`, `ci fix tweak`, `fix tweak`, `smart ci fix`

> [!IMPORTANT]
> Prompt Version: 2.5.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Whatever directives, constraints, checklists, or instructions are given before this section or prompt (including in the prompt preamble, header blocks, or incoming user request) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below.

- [A] Base Workflow: `01-prompts/16-ci-cd/03-ci-cd-fix.md` (or V2: `01-prompts/16-ci-cd/08-ci-cd-fix-with-n-steps.md`)
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

---

## Smart Targeted Fix Execution (Fastest Path)

Follow the core RCA structure from [A] and review past failures in [B]. The primary goal is the fastest possible resolution of the failing stack trace with zero wasted test cycles.

### Mandatory Smart Testing Rules

1. **Top-Instruction Precedence:** Any custom instruction provided before this section or in chat overrides default procedures.
2. **Stack Trace Targeting:** Isolate and build/test ONLY the packages and test functions directly cited in the provided failure or error stack trace.
3. **Changed Packages from Last Git Hash:** Compare changes against the last known git hash (`git diff --name-only HEAD~1` or `git status --porcelain`). Re-run tests ONLY for packages (Go, TS, Rust, Python) that contain actual modifications.
4. **Change Tracking & State Persistence:** Every time a fix is applied, write the modified file list and state to [C] (or run `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`) so the build system knows exactly which targets require re-verification.
5. **Strict Ban on Extraneous Runs:** NEVER run the complete test suite, global linters, spellcheckers, or unrelated packages. Verify strictly using `python [D] --pkg <affected_package>` or `python [D] --changed-only` until green.
6. **Mandatory Same-Turn Tool Chaining:** Do not end the turn after outputting the diagnostic plan; invoke the diagnostic or fix tool in the exact same turn.
7. **Final Step Atomic Commit Mandate:** Accumulate all fixes and stage/commit them in a single atomic commit at the end (`git commit -m "fix(ci): <summary>"`). TOTAL BAN on per-file commits.
8. **Issue & RCA Destination Routing:** Store CI/CD issue post-mortems in `.ai-memory/cicd-issues/` (indexed in `.ai-memory/cicd-index.md`). If an issue is an application bug rather than a CI/CD failure, document it in `02-spec/22-app-issues/` (indexed in `02-spec/22-app-issues/readme.md`).

---

## Execution Prompt (Copy & Paste Trigger)

```markdown
- [A] Base Workflow: `01-prompts/16-ci-cd/03-ci-cd-fix.md`
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

Follow workflow [A] and check past post-mortems in [B]. Perform a grounded 4-part RCA on the failure below (storing CI/CD RCAs in [B], or application bug RCAs in 02-spec/22-app-issues/) and apply a surgical fix without altering CI definitions or business logic.
SMART TEST ONLY: Do NOT run the full test suite. Identify packages from the stack trace and files changed from the last git hash, persist modified files to [C], and run targeted builds/tests ONLY for affected packages via `python [D] --changed-only` or `--pkg <target>` to verify the fix with maximum speed. Do not commit file-by-file; commit atomically at the final step.

<paste pipeline error / stack trace here>
```

# CI/CD Fix Tweak with Targeted Smart Testing & RCA — Workflow (must follow)

Trigger Keywords & Aliases: `cicd fix tweak`, `ci fix tweak`, `fix tweak`, `smart ci fix`

> **Prompt Version:** 1.0.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

- [A] Base Workflow: `01-prompts/16-ci-cd/03-ci-cd-fix.md`
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

---

## Smart Targeted Fix Execution (Fastest Path)

Follow the core RCA structure from [A] and review past failures in [B]. The primary goal is the fastest possible resolution of the failing stack trace with zero wasted test cycles.

### Mandatory Smart Testing Rules

1. **Stack Trace Targeting:** Isolate and build/test ONLY the packages and test functions directly cited in the provided failure or error stack trace.
2. **Changed Packages from Last Git Hash:** Compare changes against the last known git hash (`git diff --name-only HEAD21^1 or `sit status --porcelain`). Re-run tests ONLY for packages (Go, TS, Rust, Python) that contain actual modifications.
3. **Change Tracking & State Persistence:** Every time a fix is applied, write the modified file list and state to [C] (or run `python 03-ai-scripts/33-test-inventory-generator.py --record <files>...`) so the build system knows exactly which targets require re-verification.
4. **Strict Ban on Extraneous Runs:** NEVER run the complete test suite, global linters, spellcheckers, or unrelated packages. Verify strictly using `python [D] --pkg <affected_package>` or `python [D] --changed-only` until green.

---

## Execution Prompt (Copy & Paste Trigger)

```markdown
- [A] Base Workflow: `01-prompts/16-ci-cd/03-ci-cd-fix.md`
- [B] Prior Issues: `.ai-memory/cicd-issues/`
- [C] Change State Tracker: `.ai-memory/temp/recent-file-changes.json`
- [D] Local Runner: `03-ai-scripts/06-cicd-local-runner.py`

Follow workflow [A] and check past post-mortems in [B]. Perform a grounded 4-part RCA on the failure below and apply a surgical fix without altering CI definitions or business logic.
SMART TEST ONLY: Do NOT run the full test suite. Identify packages from the stack trace and files changed from the last git hash, persist modified files to [C], and run targeted builds/tests ONLY for affected packages via `python [D] --changed-only` or `--pkg <target>` to verify the fix with maximum speed.

<paste pipeline error / stack trace here>
```

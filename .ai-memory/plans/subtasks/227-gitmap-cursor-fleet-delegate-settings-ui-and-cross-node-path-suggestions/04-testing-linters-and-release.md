# Subtask 227.04: Quality Verification, Linters, Minor Version Bump & Release Ceremony

**Subtask ID:** `227.04`  
**Parent Plan:** [227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md](../../227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md)  
**Spec References:**
- [01-architecture-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/01-architecture-spec.md)
- [02-component-and-ui-spec.md](../../../../02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/02-component-and-ui-spec.md)
**Status:** Queued (Ready for execution following Subtasks 01-03)  
**Target Area:** `linter-scripts/`, `cli/`, `version.json`, `changelog.md`  

---

## 1. Objective

Conduct rigorous quality assurance across all newly authored and refactored components in Task 227:
1. Enforce strict adherence to coding guidelines via Python and Go linters with zero warnings or errors.
2. Execute targeted Go unit tests covering `cli/cmdcursor`, `cli/cmdui`, `cli/completion`, and `cli/cmdpipeline`.
3. Perform the canonical SemVer minor release ceremony bumping `6.490.0` to `6.491.0` in `version.json` and `changelog.md`.
4. Stage, commit, and push all modifications atomically using `gitmap cpf`.

---

## 2. Step-by-Step Execution Plan

### Step 1: Coding Guideline Linters & Static Analysis
1. Run Nested-Ifs Linter:
   - Command: `python linter-scripts/check-nested-ifs.py`
   - Scope: `cli/cmdcursor/`, `cli/cmdui/`, `cli/completion/`
   - Gate: Exit code 0 (zero nested `if` statements allowed; early guard returns required).
2. Run Enums and Booleans Linter:
   - Command: `python linter-scripts/check-enum-and-boolean.py`
   - Scope: Modified Go files and TypeScript/JS UI scripts
   - Gate: Exit code 0 (all boolean identifiers must carry positive prefixes `is`, `has`, `can`, `auto`; no explicit `== true` or mixed polarity).
3. Run Relative Paths Linter:
   - Command: `python linter-scripts/check-relative-paths.py`
   - Scope: Entire repository
   - Gate: Exit code 0 (zero absolute filesystem paths permitted in code or specs).
4. Run CI/CD Local Quality Gates:
   - Command: `python 03-ai-scripts/06-cicd-local-runner.py --failed`
   - Gate: All configured preflight quality gates report passing status.

### Step 2: Go Unit Test Suite Execution
1. Test Cursor Subsystem:
   - Command: `go test -v ./cli/cmdcursor`
   - Verifies: Flag parsing, mutually exclusive validation (`--skip-settings` vs `--settings-only`), default parameter synthesis, and command routing.
2. Test Web UI & Settings Subsystem:
   - Command: `go test -v ./cli/cmdui`
   - Verifies: `loadSettings` / `saveSettingsData` attribute preservation, JSON envelope encapsulation, and `/api/terminal/exec` context timeout cancellation.
3. Test Completion Subsystem:
   - Command: `go test -v ./cli/completion`
   - Verifies: Dynamic completion with cross-node candidates (`u1:<repo>`), prefix filtering, and candidate deduplication.
4. Gate: 100% test success rate across all packages with zero flaky or unmocked OS failures.

### Step 3: SemVer Minor Bump Ceremony
1. Update `version.json`:
   - Inspect current version (`6.490.0`).
   - Increment minor version: `"Version": "6.491.0"`.
2. Update `changelog.md`:
   - Prepend new release entry at the top:
     ```markdown
     ## [v6.491.0] - 2026-10-05

     ### Added
     - native cursor fleet delegation cli command (`gitmap cursor migrate` / `gitmap cur delegate`)
     - modernized settings ui with 6 sub-tabs, full 19-setting two-way binding, and embedded terminal widget
     - installer candidate path discovery and cross-node path suggestions engine
     ```
3. Update Root `readme.md`:
   - Synchronize any version badges or references to `6.491.0`.

### Step 4: Atomic Commit and Push Ceremony
1. Verify clean repository state:
   - Confirm no untracked scratch files or unneeded build artifacts exist.
2. Execute GitMap Fast Atomic Commit:
   - Command: `gitmap cpf "feat(cursor-fleet): native cursor migrate command, modern settings ui with embedded terminal, cross-node path suggestions"`
3. Verify remote tracking:
   - Ensure working tree is clean and local branch is synchronized with origin.

### Step 5: Master Ledger & Memory Finalization
1. Update `02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/00-master-audit-ledger.md`:
   - Mark Subtask 01, Subtask 02, Subtask 03, and Subtask 04 as `Completed`.
   - Update master status to `Completed`.
2. Update `.ai-memory/plans/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md` to `COMPLETED`.

---

## 3. Evidence Checklist
- [ ] `check-nested-ifs.py` returns exit code 0.
- [ ] `check-enum-and-boolean.py` returns exit code 0.
- [ ] `check-relative-paths.py` returns exit code 0.
- [ ] `go test ./cli/cmdcursor ./cli/cmdui ./cli/completion` passes 100%.
- [ ] `version.json` updated to `6.491.0`.
- [ ] `changelog.md` updated with release notes.
- [ ] `gitmap cpf` executes atomic commit and push.
- [ ] `00-master-audit-ledger.md` marked completed.

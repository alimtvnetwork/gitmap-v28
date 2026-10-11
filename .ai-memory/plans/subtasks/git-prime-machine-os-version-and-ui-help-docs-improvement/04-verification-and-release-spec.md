# Subtask 04: Quality Gate, Minor Version Bump & Release Specification

> **Subtask ID:** Task-04  
> **Parent Plan:** `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md`  
> **Spec Reference:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/02-component-and-cli-spec.md`  
> **Assigned Role:** Release Orchestrator & Quality Gate Lead  
> **Status:** PENDING  
> **Owned Files:**  
> - `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/04-verification-and-release-spec.md`  

---

## 1. Deliverables & Release Objectives

This subtask defines the comprehensive quality verification, version synchronization, changelog generation, plan consolidation, and release ceremony for the **Git Prime Machine OS Version, HD Docs UI & Merge-AI Improvement** milestone.

### Key Deliverables:
1. **Pre-Flight Quality Gate Execution:**
   - Execute FastGate AST linter (`python 03-ai-scripts/50-fastgate.py`) to verify all Go, Python, and TypeScript code passes with 0 syntax or architectural errors in $<1.5\text{s}$.
   - Verify strict relative path hygiene (`python 03-ai-scripts/49-verify-privacy-and-relative-paths.py`) to confirm zero absolute paths or drive letters exist.
2. **SemVer Minor Version Bump (`v6.533.0` &rarr; `v6.534.0`):**
   - Synchronize `Version: "6.534.0"` in `version.json`.
   - Update `constants.Version = "6.534.0"` in `cli/constants/constants.go`.
   - Update `version: "6.534.0"` in `package.json`.
   - Run `python 03-ai-scripts/14-version-sync-checker.py` to confirm 100% version synchronization across the repository.
3. **Changelog Authoring (`changelog.md`):**
   - Document new features, improvements, and fixes for `v6.534.0` under standard release headings.
4. **Plan Consolidation & Memory Archival:**
   - Verify all 4 subtasks are marked COMPLETED.
   - Archive parent plan from `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md` to `.ai-memory/plans/completed/`.
   - Update `.ai-memory/plans/readme.md` index.
5. **Atomic Commit, Tagging & Upstream Push:**
   - Stage all changed files.
   - Execute atomic Git commit:
     `gitmap cpf "feat(machine,hd-docs,merge-ai): structured os version, interactive docs, and parameterized examples"`
   - Create annotated Git tag `v6.534.0`.
   - Push commits and tags to upstream remote repository ("No Push = Not Done" invariant).

---

## 2. Concrete Verification Commands & Test Matrix

| Step | Target Subsystem | Verification Command | Expected Outcome | Failure Remediation |
|---|---|---|---|---|
| **V-01** | **FastGate AST Linter** | `python 03-ai-scripts/50-fastgate.py` | Exit code 0, 0 AST errors, execution time $<1.5\text{s}$ | Fix syntax, import, or type errors reported |
| **V-02** | **Relative Path Hygiene** | `python 03-ai-scripts/49-verify-privacy-and-relative-paths.py` | Exit code 0, zero drive letters (`C:\`, `D:\`), zero `file://` URIs | Replace with forward-slash relative Git paths |
| **V-03** | **Version Synchronization** | `python 03-ai-scripts/14-version-sync-checker.py` | Exit code 0, versions in `version.json`, `package.json`, and `changelog.md` all match `6.534.0` | Re-run `python 03-ai-scripts/37-bump-version.py` |
| **V-04** | **Machine CLI Terminal Output** | `go run ./cli/cmd/gitmap machine` | Displays formatted terminal block with structured OS release, build, and edition | Verify `renderSingleMachineIdentity` formatting |
| **V-05** | **Machine CLI JSON Output** | `go run ./cli/cmd/gitmap machine --json` | Emits valid JSON array containing `osVersionDetails` with `isWindowsServer`, `breakdown.major`, etc. | Verify `MachineIdentity` JSON tags and struct values |
| **V-06** | **Merge-AI Help Interception** | `go run ./cli/cmd/gitmap merge-ai --help` | Displays colorized PTerm manual with all 5 input modalities, exit code 0 | Fix help routing in `cli/cmdmergeai/merge_ai_core.go` |
| **V-07** | **Merge-AI Unit Tests** | `go test -v ./cli/cmdmergeai/...` | 100% tests passing, zero failures | Fix test assertions in `merge_ai_test.go` |
| **V-08** | **OS Info Unit Tests** | `go test -v ./cli/cmdos/...` | 100% tests passing, zero failures | Fix test assertions in `os_info_test.go` |
| **V-09** | **Web UI Build / Typecheck** | `npm run build` (or `tsc --noEmit`) | TypeScript compiles without errors, new commands registered in `commands.ts` | Fix TypeScript typings in `src/data/commands.ts` |
| **V-10** | **Git Status Verification** | `git status --porcelain` | Clean working tree after commit, only untracked scratch files if any | Ensure all owned files staged and committed |

---

## 3. Full Acceptance Criteria Traceability Matrix (AC-01 through AC-12)

| Spec ID | Category | Description | Verification Step | Status |
|---|---|---|---|---|
| **AC-01** | OS Identity | Structured OS Version Object in `MachineIdentity` (`osVersionDetails`) | V-05 | Pending Execution |
| **AC-02** | OS Identity | Positive Affirmative Booleans (`isWindowsServer`, `isWindowsPro`, `isWindowsHome`) | V-05 | Pending Execution |
| **AC-03** | OS Identity | SemVer Version Breakdown (`major`, `minor`, `build`, `revision`) | V-05 | Pending Execution |
| **AC-04** | OS Identity | Cross-Platform Parity (Linux `/etc/os-release`, macOS `sw_vers`) | V-08 | Pending Execution |
| **AC-05** | Terminal UI | Human ANSI Block Layout with OS Breakdown | V-04 | Pending Execution |
| **AC-06** | Merge-AI CLI | Help Interception (`-h`, `--help`, `help`) with exit code 0 | V-06 | Pending Execution |
| **AC-07** | Merge-AI CLI | Named Parameter Flags (`-d`, `--destination`, `-s`, `--source`, `--url`, `--config`) | V-07 | Pending Execution |
| **AC-08** | Merge-AI CLI | Support for All 5 Input Modalities | V-07 | Pending Execution |
| **AC-09** | Merge-AI CLI | Collision Staging Architecture (`01_`/`02_`, manifest, instruction.md) | V-07 | Pending Execution |
| **AC-10** | HD Docs UI | Command Catalog Additions in `src/data/commands.ts` (9 commands) | V-09 | Pending Execution |
| **AC-11** | HD Docs UI | Interactive `src/pages/MergeAIPage.tsx` with Live Preview | V-09 | Pending Execution |
| **AC-12** | Release Gate | FastGate $<1.5\text{s}$, Version Bump (`v6.534.0`), Zero Absolute Paths | V-01, V-02, V-03 | Pending Execution |

---

## 4. Release Ceremony & Execution Protocol

### Step 1: Pre-Release Quality Verifications
```bash
python 03-ai-scripts/50-fastgate.py
python 03-ai-scripts/49-verify-privacy-and-relative-paths.py
go test -v ./cli/cmdos/... ./cli/cmdmergeai/...
```

### Step 2: Version Bump Ceremony
```bash
# Update version.json to 6.534.0 and propagate
python 03-ai-scripts/37-bump-version.py --minor
python 03-ai-scripts/14-version-sync-checker.py
```

### Step 3: Changelog Entry (`changelog.md`)
Append the following release notes under `## [6.534.0]`:
```markdown
## [6.534.0] - 2026-10-11

### Added
- **Deep OS Version Object (`osVersionDetails`):** Structured operating system inspection in `gitmap machine --json` providing `osType`, `marketingName`, `edition`, positive affirmative booleans (`isWindowsServer`, `isWindowsPro`, `isWindowsHome`), and numerical breakdown (`major`, `minor`, `build`, `revision`).
- **Enhanced Machine Terminal Block:** Human ANSI formatting displaying OS Platform, Release, Build, and Edition breakdown.
- **Merge-AI Help & Parameterized Examples:** Comprehensive help flag interception (`-h`, `--help`, `help`) and support for all 5 input modalities (positional, named flags, comma-separated, text manifest, JSON config).
- **HD Docs UI Catalog Expansion:** Added interactive documentation for 9 recent commands (`merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, `cpar`) in `src/data/commands.ts`.
- **Dedicated Interactive Merge-AI Page (`src/pages/MergeAIPage.tsx`):** Web UI parameter builder, live CLI generator, and collision staging simulator.

### Changed
- Improved `parseMergeAIArgs` in `cli/cmdmergeai/merge_ai_core.go` to handle named flags `-d`, `-s`, `-u`, and `-c`.
```

### Step 4: Atomic Commit & Remote Push
```bash
gitmap cpf "feat(machine,hd-docs,merge-ai): structured os version, interactive docs, and parameterized examples"
git tag -a v6.534.0 -m "Release v6.534.0: Machine OS version object, HD docs UI, and merge-ai parameterized examples"
git push origin main --tags
```

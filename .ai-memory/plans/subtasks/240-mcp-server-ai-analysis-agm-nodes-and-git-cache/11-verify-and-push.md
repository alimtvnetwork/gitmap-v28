# Subtask 11: End-to-End Verification, Coding Guideline Auditing & Release Readiness

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `11-verify-and-push`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 8)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Execute end-to-end quality assurance, automated coding standard compliance checks, relative path linter validations, and prepare atomic commit readiness across all subsystems implemented in task `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`.

### The Scope
Validate that all newly implemented and refactored components meet GitMap's stringent production standards:
1. **Full Test Suite Execution:** Verify unit and integration test suites across `cli/cmdchildpath`, `cli/cmdinstall`, `cli/cmd`, `cli/cmdnodes`, `cli/cmdgitcache`, and `cli/cmdai`.
2. **Coding Guideline Linter Audit:** Ensure zero violations of Go coding standards (positive booleans, PascalCase structs/tables, early returns, under 100-line file sizing where applicable) using `python 03-ai-scripts/05-guideline-autofixer.py --check-only`.
3. **Relative Path Hygiene Audit:** Run `python linter-scripts/check-relative-paths.py` to ensure zero absolute path leakage in code, database schemas, specs, or plans.
4. **Visual Parity Audit:** Validate output and behavior against verified screenshots:
   - `assets/screenshots/240-muse-meta-installer.png`
   - `assets/screenshots/240-get-childitem-filter.png`
   - `assets/screenshots/240-git-grep-func-runfix.png`
5. **Git Hygiene & Clean Status:** Inspect repository status, stage modified files, verify commit messages adhere to atomic conventional commits, and confirm upstream readiness.

---

## 2. Concrete Files to Audit & Verify

| Subsystem | Core Packages / Paths | Verification Target |
| :--- | :--- | :--- |
| **Child-Path Scanner** | `cli/cmdchildpath/` | Traversal speed (<15ms), glob filtering, JSON output integrity. |
| **Meta Muse Installer** | `cli/cmdinstall/` | Cross-platform commands, dry-run parsing, version probe. |
| **Smart Repo Create** | `cli/cmd/` | Elimination of force-push, remote probe, interactive/headless gates. |
| **AGM Fleet Purge** | `cli/cmdnodes/` | AES-256-GCM encryption, 3-pass zero shredder, all-nodes success gate. |
| **Git Commit Cache** | `cli/cmdgitcache/` | SQLite schema, 0ms fast path, ahead/behind calculation, terminal UI. |
| **AI Analysis & MCP** | `cli/cmdai/`, `cli/store/` | Zero-loss removal vault, Split SQLite `ai-analysis.db`, LLM training export. |

---

## 3. Verification Protocol & Quality Gates

```
+----------------------------------------------------------------------------------------------------+
|                                    QUALITY GATE PIPELINE                                           |
+----------------------------------------------------------------------------------------------------+
|                                                                                                    |
|    [Step 1: Go Unit Tests]                                                                         |
|    go test -v ./cli/cmdchildpath/... ./cli/cmdinstall/... ./cli/cmdgitcache/...                   |
|    ---------------------------------------------------------------------------> PASS               |
|                                                                                                    |
|    [Step 2: Coding Guideline Linter]                                                               |
|    python 03-ai-scripts/05-guideline-autofixer.py --check-only                                      |
|    ---------------------------------------------------------------------------> 0 Violations       |
|                                                                                                    |
|    [Step 3: Relative Path Hygiene]                                                                 |
|    python linter-scripts/check-relative-paths.py                                                   |
|    ---------------------------------------------------------------------------> 0 Absolute Paths   |
|                                                                                                    |
|    [Step 4: CLI Execution Sanity]                                                                  |
|    gitmap child-path . "*.go" --depth 1                                                            |
|    gitmap cmp main HEAD                                                                            |
|    ---------------------------------------------------------------------------> Clean Execution    |
+----------------------------------------------------------------------------------------------------+
```

---

## 4. Acceptance Criteria

- [ ] All unit tests across modified packages pass with 0 failures and 0 panics.
- [ ] `05-guideline-autofixer.py` reports clean compliance across modified Go packages.
- [ ] `check-relative-paths.py` reports zero absolute path violations.
- [ ] No force-push code exists anywhere in repository creation paths.
- [ ] All subtask documentation files (`01` through `11`) are complete, non-empty, and reference correct specs.
- [ ] All git working tree modifications are accounted for.

---

## 5. Verification Commands

```powershell
# 1. Run all relevant package test suites
go test -v ./cli/cmdchildpath/...
go test -v ./cli/cmdinstall/... -run TestInstallMuse
go test -v ./cli/cmd/... -run TestRepoCreateRemote
go test -v ./cli/cmdnodes/... -run TestShred
go test -v ./cli/cmdgitcache/...

# 2. Run coding guideline verification
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdchildpath cli/cmdinstall cli/cmd cli/cmdnodes cli/cmdgitcache --check-only

# 3. Run relative path hygiene linter
python linter-scripts/check-relative-paths.py

# 4. Check git status
git status --short
```

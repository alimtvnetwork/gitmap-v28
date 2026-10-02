# Subtask 06: Verification, Minor Bump & Release Ceremony

> **Parent Plan:** [67-nodes-cfr-remote-fleet-clone-enhancement](../../pending/67-nodes-cfr-remote-fleet-clone-enhancement.md)
> **Component Spec:** [02-component-spec.md](../../../../02-spec/21-app/67-nodes-cfr-remote-fleet-clone-enhancement/02-component-spec.md)
> **Primary Targets:**
> - `version.json`
> - `readme.md`
> - `cli/cmdnodes/`
> - `cli/cmd/`

---

## 1. Objective

Execute comprehensive verification, coding guideline enforcement, automated lint remediation, SemVer minor version bump (`v6.454.0`), and atomic commit ceremony for task 67 (`67-nodes-cfr-remote-fleet-clone-enhancement`).

---

## 2. Detailed Technical Scope

### 2.1 File-Scoped Guideline Remediation
Execute the canonical guideline autofixer across all files created or modified in task 67:
```powershell
python 03-ai-scripts/05-guideline-autofixer.py --files cli/cmdnodes/nodes_clone_types.go cli/cmdnodes/nodes_clone_file.go cli/cmdnodes/nodes_clone.go cli/cmdnodes/nodes_clone_remote.go cli/cmdnodes/nodes_clone_table.go cli/cmdnodes/nodes_clone_help.go cli/cmd/nodes_cmd.go cli/helptext/catalog.go
```

Verify compliance against core standards:
1. **Positive Boolean Naming:** All boolean fields and variables must use positive prefixes (`is*`, `has*`). Negative names (`disable`, `not*`, `no*`) are forbidden.
2. **Function Length Limit:** Every function must be $\le 15$ lines of executable code.
3. **Structured Error Management:** All returned errors must be wrapped using `apperror.WrapSimple()` or `apperror.WrapWithDetails()`.
4. **Clean Builds:** Zero compiler warnings or errors during package compilation.

### 2.2 Compilation & Test Suite Verification
Verify package integrity:
```powershell
go build ./cli/...
go test -v ./cli/cmdnodes/
```

### 2.3 SemVer Minor Release Bump (`v6.454.0`)
1. **`version.json`:**
   - Update `Version` field from `6.453.0` to `6.454.0`.
2. **`readme.md`:**
   - Synchronize release badges and version strings to `6.454.0`.
   - Document new capabilities:
     - Modernized CFR Terminal UI with active machine status and reachability metrics.
     - Remote GitMap version discovery over JSON protocol.
     - High-performance asynchronous fleet clone dispatch.
     - Intelligent CWD relative work directory preservation and `-d/--dest` flags.
     - Comprehensive `gitmap help nodes` manual and post-execution actionable suggestions.

### 2.4 Atomic Commit Protocol
Execute the final commit ceremony using GitMap:
```powershell
gitmap cpf "nodes - modernize cfr ui remote version probe and async dispatch"
```

> [!CAUTION]
> **STRICT BAN ON RAW GIT COMMANDS:** Never execute raw `git` commands (`git add`, `git commit`, `git status`). Always use `gitmap cpf`.

---

## 3. Implementation Steps

1. **Step 1:** Run file-scoped guideline autofixer across modified packages.
2. **Step 2:** Inspect any remaining lint warnings and surgically refactor functions exceeding 15 lines.
3. **Step 3:** Compile all packages via `go build ./cli/...` and run unit tests.
4. **Step 4:** Bump `version.json` and sync `readme.md` to `6.454.0`.
5. **Step 5:** Execute atomic commit via `gitmap cpf "nodes - modernize cfr ui remote version probe and async dispatch"`.

---

## 4. Acceptance Criteria

- [ ] All modified files adhere to positive boolean naming (`is*`, `has*`).
- [ ] All functions $\le 15$ lines.
- [ ] All errors wrapped via `apperror`.
- [ ] `cli/helptext/nodes.md` displays properly under `gitmap help nodes`.
- [ ] `go build ./cli/...` passes with zero errors.
- [ ] `version.json` is bumped to `6.454.0`.
- [ ] Atomic commit is recorded cleanly with no uncommitted changes remaining.

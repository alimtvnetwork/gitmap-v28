# Plan 201: Cluster & SC Compare Matrix Dispatch & Help Text Parity

## Task Lifecycle & Completion Summary
- **Initiated**: User requested full architectural parity across `cluster` and `sc` (servers-clients) execution, including comparison matrix commands (`gitmap cluster compare`, `gitmap sc compare`, `matrix`), explanation of why each command is different, how to join, and how to monitor machines across the Triad architecture.
- **Loops/Steps**: Completed across 2 autonomous loops within budget.
- **Status**: COMPLETED (100% Verified)

---

## 1. Overview & Problem Statement
The user requested:
> *"So similar should happen with the cluster execution, also server times (`sc`) commands. That should have an explanation why each command is different. So there should be a comparing table that would explain which one is using one when and how to use them, how to use the join, how to monitor the machines, things like that. These are kind of missing from your documentation."*

To satisfy this:
1. **CLI Routing**:
   - `gitmap cluster compare` and `gitmap cluster matrix` routed to `cmdssh.RunSSHCompareCLI`.
   - `gitmap sc compare` and `gitmap sc matrix` routed to `cmdssh.RunSSHCompareCLI`.
2. **Help Text Documentation Parity**:
   - `cli/helptext/cluster.md`: added the full Subsystems Architecture Comparison Table and command guidance.
   - `cli/helptext/sc.md`: added the full Subsystems Architecture Comparison Table and command guidance.
3. **Coding Guidelines Compliance**:
   - Decomposed `dispatchSCNodeOps` and `dispatchSCMetaOps` in `cli/cmd/rootcore.go` to strictly adhere to $\le 15$ line function limits.
   - Verified `go vet` and file-level linters.

---

## 2. Changes Made & Consolidated Subtasks

### Subtask 01: Export & Wire Comparison Matrix across Cluster & SC
- **`cli/cmdssh/ssh_compare_table.go`** (88 lines):
  - Exported `RunSSHCompareCLI(args []string) error` for use by `cmd/cluster.go` and `cmd/rootcore.go`.
- **`cli/cmd/cluster.go`**:
  - Added `case "compare", "matrix":` to `routeClusterBootstrapOrExec`.
- **`cli/cmd/rootcore.go`**:
  - Decomposed `dispatchSCNodeOps` into `dispatchSCMetaOps` and `dispatchSCNodeOps` ($\le 15$ lines each).
  - Added `case "compare", "matrix":` to route directly to `cmdssh.RunSSHCompareCLI`.

### Subtask 02: Documentation & Help Text Synchronization
- **`cli/helptext/cluster.md`**:
  - Added the complete Subsystems Architecture Comparison table under `## Subsystems Architecture Comparison (compare / matrix)`.
  - Added CLI example `gitmap cluster compare`.
  - Included detailed breakdown of when to use `ssh` vs `cluster` vs `sc`.
- **`cli/helptext/sc.md`**:
  - Added the complete Subsystems Architecture Comparison table under `## Subsystems Architecture Comparison (compare / matrix)`.
  - Added CLI example `gitmap sc compare`.
  - Included detailed breakdown of when to use `ssh` vs `cluster` vs `sc`.

---

## 3. Verification & Compliance Matrix

| Quality Gate | Tool / Command | Result |
|---|---|---|
| **Cluster Compare CLI** | `gitmap cluster compare` & `matrix` | **100% PASS** (aligned terminal table output) |
| **SC Compare CLI** | `gitmap sc compare` & `matrix` | **100% PASS** (aligned terminal table output) |
| **Cluster Help Text** | `gitmap cluster help` | **100% PASS** (renders comparison table) |
| **SC Help Text** | `gitmap sc help` | **100% PASS** (renders comparison table) |
| **Control Flow (Nested Ifs)** | `python linter-scripts/check-nested-ifs.py` | **100% PASS** (0 violations across 3,098 files) |
| **Booleans & Enums** | `python linter-scripts/check-enum-and-boolean.py` | **100% PASS** (0 violations across 2,334 files) |
| **Go Vet** | `go vet ./cmd ./cmdssh` | **100% PASS** (0 warnings) |
| **Function Length** | Canonical limit: $\le 15$ lines per function | **100% PASS** |
| **Inventory Recording** | `python 03-ai-scripts/33-test-inventory-generator.py --record` | **Recorded** (46 tracked files) |

# Plan 310: Strict Relative Git Paths & Absolute Path Elimination

## User Request (Verbatim)
```text
@[01-prompts/15-cg-execute/12-relative-paths.md]
```

## 1. Context & Goal
Eliminate all hardcoded absolute filesystem paths (`D:\...`, `C:\...`) and `file:///` URIs from documentation, specifications, plans, and helpdocs across the GitMap codebase. Replace them with standard repository-relative forward-slash paths.

## 2. Invariants & Guardrails
- **Zero Builds / Zero Tests:** Ban on `go build`, `npm run build`, and test suite execution.
- **Worker Git Ban:** Subagents never run git commands.
- **Tree Swarm Concurrency:** $A = 3$, $H = 2$, Capacity = 6.
- **Strict Relative Paths:** All committed files must contain zero absolute drives or file URIs.
- **FastGate & Linter Cleanliness:** All checks must pass with exit 0.

## 3. Disjoint Subtask Breakdown
| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | Completed Plan Path Sanitization | Worker 01 | `.ai-memory/plans/completed/297-result-wrapper-stutter-eradication.md` | DONE | Line 93 converted to relative path `01-prompts/...` |
| Task-02 | Spec 298/299 Path Sanitization | Worker 02 | `02-spec/21-app/298-result-wrapper-slice-collection-and-apperror/01-architecture-spec.md`, `02-spec/21-app/299-result-wrapper-store-ssh-and-templates/01-architecture-spec.md` | DONE | Line 5 in both specs converted to relative path `01-prompts/...` |
| Task-03 | CLI Helpdoc URI Link Sanitization | Worker 03 | `cli/helpdoc/devtool.md`, `cli/helpdoc/pulle.md` | DONE | `file:///` URIs replaced with relative markdown links |

## 4. Verification Gates
- `python linter-scripts/check-relative-paths.py` (PASS — exit 0 across 9,834 files)
- `python 03-ai-scripts/50-fastgate.py` (PASS — exit 0)
- `python linter-scripts/check-forbidden-strings.py` (PASS — exit 0)
- Final atomic push via `gitmap cpf "paths - standardize strict relative git paths and eliminate absolute paths"`

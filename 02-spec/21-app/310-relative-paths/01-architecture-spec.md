# Plan 310 Architecture Specification: Strict Relative Git Paths & Absolute Path Elimination

## User Request (Verbatim)
```text
@[01-prompts/15-cg-execute/12-relative-paths.md]
```

## 1. Executive Summary & Architectural Intent
This specification governs the systematic elimination of machine-specific, absolute filesystem paths (`D:\...`, `C:\...`, `/home/...`) and `file:///` URIs across the repository, replacing them with clean, repository-root-relative paths.

Hardcoded absolute paths compromise cross-platform interoperability, pollute version control with developer-specific host environments, and violate repository coding guidelines (Rule R11: Strict Relative Git Paths).

## 2. In-Scope Files & Target Violations
The initial automated audit identified 5 files containing absolute paths and `file:///` URIs:
1. `.ai-memory/plans/completed/297-result-wrapper-stutter-eradication.md`:
   - Line 93: Hardcoded absolute path `@[d:\work\gitmap\01-prompts\15-cg-execute\19-result-wrapper-and-apperror-returns.md]` &rarr; replace with `@[01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md]`.
2. `02-spec/21-app/298-result-wrapper-slice-collection-and-apperror/01-architecture-spec.md`:
   - Line 5: Hardcoded absolute path `@[d:\work\gitmap\01-prompts\15-cg-execute\19-result-wrapper-and-apperror-returns.md]` &rarr; replace with `@[01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md]`.
3. `02-spec/21-app/299-result-wrapper-store-ssh-and-templates/01-architecture-spec.md`:
   - Line 5: Hardcoded absolute path `@[d:\work\gitmap\01-prompts\15-cg-execute\19-result-wrapper-and-apperror-returns.md]` &rarr; replace with `@[01-prompts/15-cg-execute/19-result-wrapper-and-apperror-returns.md]`.
4. `cli/helpdoc/devtool.md`:
   - Line 3: URI link `(file:///cli/helptext/devtools.md)` &rarr; replace with relative link `(cli/helptext/devtools.md)`.
5. `cli/helpdoc/pulle.md`:
   - Line 3: URI link `(file:///cli/helptext/pull-error.md)` &rarr; replace with relative link `(cli/helptext/pull-error.md)`.

## 3. Disjoint Subtask Allocation
Execution is partitioned across three autonomous workers in disjoint file boxes:
- **Task-01 (Worker 01):** Completed plan path sanitization in `.ai-memory/plans/completed/297-result-wrapper-stutter-eradication.md`.
- **Task-02 (Worker 02):** Active architecture spec path sanitization in `02-spec/21-app/298-*/01-architecture-spec.md` and `02-spec/21-app/299-*/01-architecture-spec.md`.
- **Task-03 (Worker 03):** CLI helpdoc link sanitization in `cli/helpdoc/devtool.md` and `cli/helpdoc/pulle.md`.

## 4. Verification Evidence Contract
- `python linter-scripts/check-relative-paths.py` MUST exit 0 with 0 violations across all tracked files.
- `python 03-ai-scripts/50-fastgate.py` MUST exit 0.
- `git status --porcelain` MUST only contain targeted files.

# 35 - Relative Git Paths & Absolute Path Elimination Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md)

### AC-01: Guideline Coverage
- [ ] Boolean principles define naming, evaluation, and composition patterns
- [ ] Casting elimination patterns cover type-safe alternatives to type assertions
- [ ] Code style defines formatting, naming, and structural conventions

### AC-02: Enforcement
- [ ] All guidelines include ❌ (forbidden) and ✅ (compliant) code examples
- [ ] ESLint/linter rules are documented for automated enforcement
- [ ] Master guidelines document consolidates all standards for AI reference

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule RP-1 (Strict Relative Git Paths):** All markdown links, file paths, citations, subtask paths, and code references MUST be strictly relative to the Git repository root (e.g. `cmd/main.go`, `02-spec/03-error-manage/01-index.md`).
2. **Rule RP-2 (Zero Absolute OS Paths):** Absolute filesystem paths (`/absolute/path/to/...`, `/absolute/path/to/...`, `C:\Users\...`, `/home/...`) are strictly forbidden in committed documentation and source code.
3. **Rule RP-3 (Zero `file:///` URIs):** Absolute file URI schemes (`file:///...`) are strictly prohibited in all markdown files, citations, and source code.
4. **Rule RP-4 (CI/CD Local Runner Verification):** All changes must pass `python linter-scripts/check-relative-paths.py` and `python 03-ai-scripts/06-cicd-local-runner.py` with exit code 0.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | Repository-wide | Various | Absolute paths and `file:///` URIs | Scanned 6,765 tracked files; verified zero absolute paths | VERIFIED |
| V-02 | linter-scripts/check-relative-paths.py | N/A | Relative path check | Verified exit code 0 | VERIFIED |
| V-03 | 03-ai-scripts/06-cicd-local-runner.py | 34 | Relative path quality gate | Verified in Batch 1 of local runner | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/35-relative-paths/01-absolute-path-audit.md`)**:
   Deeply scan repository files for absolute paths, drive letters, and `file:///` schemes.
2. **Subtask 02 (`.lovable/plans/subtasks/35-relative-paths/02-linter-and-ci-verification.md`)**:
   Execute `python linter-scripts/check-relative-paths.py` and full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.

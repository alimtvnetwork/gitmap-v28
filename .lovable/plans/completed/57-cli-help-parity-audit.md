# 36 - CLI Commands, Help Text Parity & Help UI Audit Specification

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

1. **Rule CLI-1 (100% Command Discoverability):** Every executable command and subcommand MUST be registered in the root CLI tree and displayed in `--help`.
2. **Rule CLI-2 (Flag & Option Coverage):** All flags MUST have clear human-readable descriptions, types, and defaults.
3. **Rule CLI-3 (Usage Examples):** Every command and subcommand `--help` MUST contain real-world terminal invocation examples.
4. **Rule CLI-4 (Help Text Parity & AST Consistency):** All registered command constants must match AST definitions and pass `03-ai-scripts/09-cli-help-auditor.py` and Go AST tests.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/ | N/A | Subcommand registration and help descriptions | Verified all CLI commands across 459 files contain help strings | VERIFIED |
| V-02 | gitmap/helptext/ | N/A | Helptext markdown fixtures & AST sync | Verified `go test -C gitmap ./helptext/...` passes | VERIFIED |
| V-03 | gitmap/constants/ | N/A | Top-level command constants uniqueness & AST parity | Verified `go test -C gitmap ./constants/...` passes | VERIFIED |
| V-04 | 03-ai-scripts/09-cli-help-auditor.py | N/A | CLI help auditor verification | Verified exit code 0 | VERIFIED |
| V-05 | 03-ai-scripts/06-cicd-local-runner.py | 37 | CLI Help Parity quality gate | Verified in Batch 1 of local runner | VERIFIED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/36-cli-help/01-command-registration-and-examples.md`)**:
   Verify that all CLI commands, subcommands, flags, and options have complete descriptions and practical terminal examples.
2. **Subtask 02 (`.lovable/plans/subtasks/36-cli-help/02-helptext-ast-parity.md`)**:
   Verify Go AST parity and helptext consistency tests.
3. **Subtask 03 (`.lovable/plans/subtasks/36-cli-help/03-ci-runner-verification.md`)**:
   Execute `python 03-ai-scripts/09-cli-help-auditor.py` and full CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py`.

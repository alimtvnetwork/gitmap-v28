---
name: cg-error-management
description: >-
  Autonomously audits, refactors, and validates repository-wide error management against spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
---

# Error Management & Architecture Coding Guidelines (`cg-error-management`)

This skill provides autonomous audit, refactoring, and validation of repository-wide error handling based on `spec/03-error-manage/` and `.lovable/coding-guidelines/coding-guidelines.md`.

## Core Invariants

1. **No Bare Panics or Bare Exits**: Zero calls to `panic("...")`, `panic(err)`, or `os.Exit(...)` outside the central dispatcher (`cliexit.HandleError`).
2. **Context-Rich `AppError` Wrappers**: All errors MUST be wrapped with `Op`, `Code`, `Type`, `Severity`, `Creator`, `Message`, `Ctx`, and `Cause`.
3. **Universal Response Envelope**: All API endpoints return `{ "data": ..., "errors": [...], "meta": ... }`.
4. **Never Swallow Errors**: Every catch block and error return must be recorded and handled explicitly.
5. **Targeted Verification**: Continuous verification via `python linter-scripts/check-error-management.py <files>`. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine fixes.


## Change Tracking & Test Avoidance
- Test execution disabled; pass `--no-tests` to local runner.
- Append modified files to `.lovable/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`).


## Routine Execution Policy
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.

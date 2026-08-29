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
5. **CI/CD Linter Hook**: Continuous verification via `linter-scripts/check-error-management.py` and `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.

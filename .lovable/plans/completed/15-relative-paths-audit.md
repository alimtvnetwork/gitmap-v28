# Plan 15: Strict Relative Git Paths & Absolute Path Elimination Audit

## Objective

Autonomously scan, plan, refactor, and fix all absolute filesystem paths (`D:\...`, `C:\...`, `/home/...`) and `file:///` URIs across the codebase, directly updating markdown documents, specifications, plans, subtasks, citations, and code comments to use strictly relative Git repository paths until 100% green without stopping.

## Background & Rationale

1. **Portability Destruction:** Absolute paths tie repository files to one specific machine, drive letter, or user directory, instantly breaking links on macOS, Linux, Windows, and CI/CD pipelines.
2. **Reviewer Friction:** Broken absolute file URIs prevent developers and other AI agents from navigating cross-references and citations.
3. **CI/CD Flakiness:** Tests or scripts relying on local machine paths fail deterministically on GitHub Actions and containerized runners.

## Target Areas

1. **Memory Logs & RCA Files:** `.lovable/memory/index.md`, `.lovable/memory/issues/*.md`, `.lovable/cicd-issues/*.md`.
2. **Plans & Subtasks:** `.lovable/plans/completed/**/*.md`, `.lovable/plans/subtasks/**/*.md`.
3. **Specifications:** `spec/**/*.md`.
4. **CI/CD & Linter Scripts:** `linter-scripts/check-relative-paths.py`, `.lovable/ai-fix-scripts/04-relative-path-fixer.py`.
5. **CI Runner Integration:** Register `Relative Path Check` in `.lovable/ai-fix-scripts/03-cicd-local-runner.py`.

## Subtasks

- [ ] Subtask 01: Build and test `.lovable/ai-fix-scripts/04-relative-path-fixer.py` and `linter-scripts/check-relative-paths.py`.
- [ ] Subtask 02: Execute `04-relative-path-fixer.py` across `.lovable/`, `spec/`, and repository files.
- [ ] Subtask 03: Verify zero remaining violations with `linter-scripts/check-relative-paths.py`.
- [ ] Subtask 04: Integrate linter check into `.lovable/ai-fix-scripts/03-cicd-local-runner.py` and verify all quality gates pass.

## Acceptance Criteria

- [ ] `python linter-scripts/check-relative-paths.py` reports ✅ PASS with 0 violations.
- [ ] `python .lovable/ai-fix-scripts/03-cicd-local-runner.py` passes with exit code 0.
- [ ] Working tree clean and properly tracked.

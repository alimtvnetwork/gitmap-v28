# Plan 18: Repository-Wide Coding Guidelines Remediation

## Overview

Audited 3019 repository files against master coding guidelines. Computed baseline compliance score of **99.0 / 100** with 1109 identified guideline gaps across sizing, booleans, enums, and error handling.

## Key Objectives

1. **Phase 1:** Enforce affirmative booleans and eliminate explicit `== true` comparisons.
2. **Phase 2:** Decompose monolithic React components (`src/pages/*.tsx`, `src/components/*.tsx`) exceeding the 100-line limit into atomic presentation and logic sub-modules.
3. **Phase 3:** Align casing conventions (`Id`, `Url`, `Api`) and enum suffixes (`*Type`).

## Subtasks Enqueued

- Target directory: `.lovable/plans/subtasks/18-coding-guidelines/`
- Total atomic subtasks: 5 modules

## Verification Plan

- Run `python linter-scripts/check-file-sizes.py`
- Run `python linter-scripts/check-boolean-guidelines.py`
- Run `python linter-scripts/check-nested-ifs.py`
- Run `python .lovable/ai-fix-scripts/03-cicd-local-runner.py`

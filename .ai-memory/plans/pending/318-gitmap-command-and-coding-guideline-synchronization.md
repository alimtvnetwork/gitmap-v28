# Master Plan: GitMap Command and Coding Guideline Synchronization

## Metadata
- **Task ID:** 318-gitmap-command-and-coding-guideline-synchronization
- **Slug:** gitmap-command-and-coding-guideline-synchronization
- **Spec Path:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/`
- **Status:** In-Progress (Phase 1: Planning & Discovery)
- **Concurrency Budget:** A = 3, H = 2, Concurrency Capacity = 6

---

## High-Level Objective
Implement native `gitmap cargo`, `gitmap rust`, and `gitmap kill` toolchain commands with automated LLVM/libclang environment resolution; modernize commit formatting to uppercase `TYPE(MODULE): <summary>` without square brackets (third brackets) and add `cfr` (refactor) support; eliminate magic strings across codebase by introducing strongly-typed enums ending in `Type`; embed `02-spec/03-error-manage/` guidelines into the V6 execute prompt; and synchronize updated skills and guidelines across all 41 fleet repositories.

---

## Discrete Task Breakdown
- [ ] **Task 318-01:** Native Toolchain Commands (`gitmap cargo`, `gitmap rust`, `gitmap kill`) with Automatic LLVM / `LIBCLANG_PATH` Detection.
- [ ] **Task 318-02:** Modernize Commit Formatter to Uppercase Semantic Types, Parentheses, No Third Brackets (`FEAT(MODULE)`, `FIX(MODULE)`, `REFACTOR(MODULE)`, `CHORE(MODULE)`), and Register `cfr` Alias.
- [ ] **Task 318-03:** Codebase Magic String Elimination & Strongly-Typed Enums (`*Type`) in `cli/constants/`.
- [ ] **Task 318-04:** Embed Error Management & Leaf Sovereignty Guidelines into V6 Execute Prompt & Coding Guidelines.
- [ ] **Task 318-05:** Unit & E2E Test Suite for Toolchains, Uppercase Formatter, and Enums.
- [ ] **Task 318-06:** HD Help Documentation, LLM Training Curriculum, and Fleet-Wide Skills Synchronization.

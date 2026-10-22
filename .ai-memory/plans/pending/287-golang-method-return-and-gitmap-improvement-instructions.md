# Plan 287: Golang Method Return & GitMap Improvement Instructions

## Status: Completed
- **Plan ID:** 287
- **Slug:** `287-golang-method-return-and-gitmap-improvement-instructions`
- **Spec Reference:** [02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/01-architecture-spec.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/01-architecture-spec.md)
- **Scope:** Golang Method Return Types, Zero Error Swallowing, 200 Core GitMap Improvements, 200 Implementation Steps, 100 UI Improvements
- **Created At:** 2026-10-22
- **Completed At:** 2026-10-22
- **Parent Goal:** Review the codebase for Golang methods that do not return an error or result wrapper, eliminate error swallowing (`_ = ...`, silent void returns), document 200 potential improvements with 200 actionable implementation steps, and design 100 UI/UX improvement ideas in dedicated specifications.

---

## 1. Executive Summary

This plan addressed all non-negotiable requirements:
1. **Golang Error Sovereignty & Return Types:** Converted void methods (`updateRepoTransportInDB`, `clearResolvedPullError`, `applySSHTransportRemediation`, `remediateCredentialAndConvertToSSH`, `loadStoreRepos`, and `ensureWorkDirsTableExists`) to return explicit `error` types. Eliminated blank discards (`_ = ...`) and silent returns. Callers now inspect, log warnings, and propagate failures cleanly.
2. **200 GitMap Architectural Improvements & 200 Concrete Steps:** Surveyed all core subsystems and generated a comprehensive catalog of 200 improvements and 200 sequential execution steps in [03-gitmap-200-improvements-and-200-steps.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/03-gitmap-200-improvements-and-200-steps.md).
3. **100 GitMap UI & Terminal UX Improvements:** Formulated 100 concrete UI/UX initiatives across 10 dimensions in [04-gitmap-100-ui-improvements.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/04-gitmap-100-ui-improvements.md).

---

## 2. Subtasks Breakdown

- [x] **Subtask 287.1:** Audit & Fix Void Error Swallowing in `cli/cmdpull/pull_remediation_ssh.go`
- [x] **Subtask 287.2:** Codebase-Wide Survey & Remediation of Swallowed Errors
- [x] **Subtask 287.3:** 200 GitMap Architectural Improvements & 200 Implementation Steps
- [x] **Subtask 287.4:** 100 GitMap UI & Terminal UX Improvement Ideas
- [x] **Subtask 287.5:** Quality Gate Verification & Atomic Git Commit

---

## 3. Specifications Index

1. [01-architecture-spec.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/01-architecture-spec.md)
2. [02-component-and-error-return-spec.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/02-component-and-error-return-spec.md)
3. [03-gitmap-200-improvements-and-200-steps.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/03-gitmap-200-improvements-and-200-steps.md)
4. [04-gitmap-100-ui-improvements.md](../../../02-spec/21-app/287-golang-method-return-and-gitmap-improvement-instructions/04-gitmap-100-ui-improvements.md)

---

## 4. Verification & Quality Gates

- **FastGate AST Linter:** 100% PASS (< 0.02s) on all modified files.
- **Unit Tests:** `cmdpull` and `cmdworkdir` passed with 100% green.
- **Path Hygiene:** 100% strictly relative Git paths (zero absolute path violations).

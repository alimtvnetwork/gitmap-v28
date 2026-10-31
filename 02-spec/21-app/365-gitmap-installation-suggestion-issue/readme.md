# Specification Index: 365 — GitMap Installation Suggestion Issue

- **Task Code:** `365-gitmap-installation-suggestion-issue`
- **Slug:** `gitmap-installation-suggestion-issue`
- **Root Directory:** `02-spec/21-app/365-gitmap-installation-suggestion-issue/`
- **Plan File:** `.ai-memory/plans/pending/365-gitmap-installation-suggestion-issue.md`

---

## 1. Documents

- [01-root-cause-analysis-and-architecture.md](01-root-cause-analysis-and-architecture.md): 6-factor grounded Root Cause Analysis and suggestion state engine.
- [02-implementation-and-verification-spec.md](02-implementation-and-verification-spec.md): Code specifications, installer scripts refactoring, Cobra alias completion, and test vectors.

---

## 2. Invariants

- **Always-Report Guarantee:** Terminal suggestions state must be visibly reported on every install, update, and check.
- **Offline / Local Parity:** `local-install.ps1` must maintain 100% parity with remote `install.ps1`.
- **Zero Unrequested Builds/Tests:** No intermediate compiler or test suite executions in downstream repositories.

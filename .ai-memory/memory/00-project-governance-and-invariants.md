# 00 — Project Governance & Architectural Invariants

- **Domain:** Project Governance, Strict Invariants & Quality Standards
- **Authoritative Status:** Active & Enforced
- **Canonical Specifications:**
  - [01-cli-architecture](../../02-spec/21-app/01-cli-architecture/01-architecture-spec.md)
  - [02-coding-guidelines](../../02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md)

---

## 1. Non-Negotiable Architectural Invariants

Every agent, tool, and subtask operating on the GitMap codebase must strictly adhere to these core governance rules:

### Invariant 1: Positive Boolean Naming Only
- **Rule:** All boolean variables, struct fields, configuration keys, and function names must use affirmative polarity.
- **Allowed Prefixes:** `is*`, `has*`, `can*`, `should*`. (e.g., `isValid`, `isReady`, `hasData`, `isPendingIsolated`).
- **Forbidden:** Negative boolean identifiers (`isNotValid`, `hasNoData`, `isUnready`) and implicit truth comparisons (`if isValid == true` in Go/TS). Double negatives like `!isNotReady` are strictly prohibited.

### Invariant 2: 100% Relative Git Paths
- **Rule:** Hardcoding absolute filesystem paths (such as `/home/user/...` or `C:\...`) is strictly prohibited in source code, configuration files, and documentation.
- **Implementation:** Always resolve paths dynamically using `filepath.Join`, `os.UserHomeDir`, or relative paths rooted at repository boundaries. Cross-references between markdown files must use relative Git paths (e.g., `../../02-spec/...`).

### Invariant 3: Zero-Build & Zero-Test Execution During Maintenance
- **Rule:** Maintenance tasks (guideline remediation, specification audits, documentation consolidation, memory updates) must NEVER trigger compiler builds, integration tests, or unit test suites.
- **Rationale:** Running heavy builds or test suites consumes massive resources, risks lock contention on SQLite databases, and pollutes git working trees. Verification must be performed via targeted static analysis linters (`check-relative-paths.py`, `check-forbidden-strings.py`).

### Invariant 4: Total Git Command Ban for Subagents
- **Rule:** Subagents MUST NEVER execute any git commands (`git add`, `git commit`, `git push`, `git rm`, `git status`, `git diff`).
- **Rationale:** Prevents index lock contention (`.git/index.lock`), micro-commit spam, and accidental dirty-state commits. All Git state manipulation is strictly owned by the lead orchestrator.

### Invariant 5: Strict Isolation of Pending & Open Work
- **Rule:** The following directories are untouchable during consolidation tasks:
  - `.ai-memory/plans/pending/` (all enqueued plans remain isolated).
  - Open active plans in `.ai-memory/plans/` (e.g., `219`, `226`, `236`).
  - Active subtask directories under `.ai-memory/plans/subtasks/`.

---

## 2. Append-Only Strictly Prohibited Registry

These rules are permanent and unalterable:

1. **No Timestamps in Readmes:** Never insert dynamic date/time/clock stamps into `readme.txt` or general documentation files.
2. **Untouchable CI Release Pipeline:** `.github/workflows/release.yml`, `.github/scripts/smoke-installer.*`, and anything under `.gitmap/release/` are strictly off-limits. Never edit, refactor, or "improve" release workflow files.
3. **Commit-In Payloads:** Never store raw file contents, file hashes, or diff payloads in commit-in database tables. Only store `RelativePath` strings (`git cat-file` resolves content on demand).
4. **Append-Only History:** Never rewrite Git history of active working repositories during standard operations. History alterations belong strictly to dedicated purge tools (`gitmap history-purge`).
5. **No Negative Preconditions:** Conditional filters must never rely on inverse boolean guards.
6. **No Raw Secrets in Git:** Plaintext passwords, tokens, API keys, or private SSH keys must never be committed. All credentials must route through the encrypted RSA vault.

---

## 3. Codebase Structure & Quality Guidelines

- **File Size Cap:** Source files must not exceed 100 lines of executable code. Split oversized files into discrete helper modules.
- **Function Body Cap:** Function implementations should be capped at 8–15 lines. Extract discrete conditions and transformations into dedicated helper functions.
- **Flat Conditional Branching:** Maximum nesting depth of `1`. Use guard clauses and early returns to flatten nested `if` statements.
- **Error Propagation:** Never swallow errors using `_ = err`. All Go errors must be wrapped with structured caller context via `*appfault.AppError`.
- **Domain Centralization:** Type definitions and data structs must be centralized in dedicated `types.go` files per package.

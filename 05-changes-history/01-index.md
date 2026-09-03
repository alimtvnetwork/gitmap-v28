# Transaction History & Task Log Index

> **Location:** `05-changes-history/`  
> **Purpose:** Canonical persistent audit trail of all repository modifications, architectural tasks, cross-repo synchronizations, and system changes.  
> **Rule:** Every AI work session must record its operations in a dedicated, numbered subfolder so that subsequent AI agents can immediately understand the project trajectory, decisions made, and pending steps.

---

## 1. Directory Structure & Naming Conventions

All entries under `05-changes-history/` must follow strict sequential naming:

```
05-changes-history/
├── 01-index.md                                    # This master index and convention specification
├── 01-<task-slug>/                                # Dedicated directory for Task 01
│   ├── 01-transaction-log.md                      # Comprehensive transaction log for Task 01
│   └── ...                                        # Supporting diffs, manifests, or intermediate reports
├── 02-<task-slug>/                                # Dedicated directory for Task 02
│   └── 01-transaction-log.md
└── ...
```

### Folder & File Rules
1. **Strict Lowercase:** Every folder and file inside `05-changes-history/` must be strictly lowercase (e.g., `01-transaction-log.md`, `01-gitmap-sync-and-ai-scripts/`).
2. **Sequential Numbering:** Two-digit zero-padded sequence prefix (`01-`, `02-`, etc.) based on execution order.
3. **Relative Paths Only:** Never include absolute filesystem paths (`C:\...` or `file:///...`) inside log files. All paths must be relative to git root.
4. **Append-Only Integrity:** Past transaction logs are historical records; never truncate or delete past logs.

---

## 2. Standard Structure of `01-transaction-log.md`

Each task transaction log must contain:
1. **Header & Metadata:** Task title, date, author/agent, affected modules.
2. **Context & Goals:** Why this task was performed and what problem it solves.
3. **Files Changed / Created:** Complete itemized list of all files modified, added, or deleted.
4. **Architectural Decisions & Rationale:** Specific reasons for non-obvious choices.
5. **Cross-Repo Actions:** Any synchronization with sibling repositories (e.g., `gitmap`, `prompt-architect-v2`).
6. **Verification & Quality Gate Results:** Tests, linters, or scripts executed to prove correctness.
7. **Next Steps / Hand-off Context:** Clear pointers for subsequent tasks.

---

## 3. Transaction Log Register

| # | Task Directory | Summary of Changes | Status | Date |
| 01 | [`01-gitmap-ai-scripts-and-spec-sync`](01-gitmap-ai-scripts-and-spec-sync/01-transaction-log.md) | Sync AI scripts with Gitmap, add Go formatter and linters, align spec architecture | Completed | 2026-09-03 |
| 02 | [`02-pluggable-writer-architecture-research`](02-pluggable-writer-architecture-research/01-transaction-log.md) | Research and specification of composable writer contracts, BaseWriter, RestAPIWriter, and package decomposition | Completed | 2026-09-03 |

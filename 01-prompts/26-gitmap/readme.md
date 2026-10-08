# GitMap AI Training & Autonomous Automation Prompts (`26-gitmap`)

This directory contains the canonical AI training prompts and specifications for **GitMap**, the ultra-fast autonomous developer companion and CLI automation engine.

- **Category:** `26-gitmap`
- **Scope:** Autonomous Developer Companion, Fast Code Search, Zero-Storage CI/CD, Split-DB SQLite, and Multi-Agent Orchestration
- **Lead Architect & Author:** MD ALIM UL KARIM (alimtvnetwork)
- **Sponsored By:** RISEUP ASIA LLC (https://riseup-asia.com)

---

## 1. Directory Index

| # | Prompt File | Target Mode | Description | Companion Skill |
| :---: | :--- | :--- | :--- | :--- |
| **01** | [`01-gitmap-core-engine.md`](01-gitmap-core-engine.md) | Canonical Specification | Full 9-feature sequence, command reference matrix, and AI directives | `gitmap` |

---

## 2. Directory Tree Representation

```text
01-prompts/26-gitmap/
├── 01-gitmap-core-engine.md
└── readme.md
```

---

## 3. Summary of GitMap Core Capabilities

1. **Multi-Core Streaming Regex Search (`gitmap aum search`):** Fast streaming regex across files; replaces slow shell grep with mandatory directory and extension scoping.
2. **Indexed Global Symbol Search (`gitmap search`):** Sub-millisecond split-db SQLite symbol lookups across indexed codebases.
3. **Rapid File Finding (`gitmap find` / `gitmap ff` / `gitmap ffa`):** Sub-10ms filename discovery matching exact names, substrings, prefixes, and suffixes.
4. **File Inventory & Listing (`gitmap lf`):** Emits relative paths for targeted directory structures.
5. **Deterministic Terminal File Streaming (`gitmap cat`):** Zero-disk file content streaming directly to stdout.
6. **Script Execution & Cache Offloading (`gitmap py` / `gitmap ps` / `gitmap rc`):** Cross-platform script execution with automatic temporary script offloading into `repo-cache`.
7. **Python Toolchain Locator & Cache Backup (`gitmap aum locate python`):** Ultra-fast (<15ms) runtime discovery persisted in `installation.db`.
8. **LLM Training Curriculum & Pipeline Diagnostics (`gitmap llm train` / `gitmap ld` / `gitmap pe`):** Chained AI curriculum and CI/CD error log extraction for 4-part RCA remediation.
9. **Semantic Hyphen-Separated Atomic Commits (`gitmap cpf` / `gitmap cpb` / `gitmap cpr`):** Clean hyphen-separated commit messages without colons.

---

## 4. Non-Negotiable AI Rules

- **Total Ban on Shell Search:** TOTAL BAN on `Select-String`, `rg`, `ripgrep`, `grep`, `git grep`, `Get-ChildItem -Recurse`, and `findstr`. Use GitMap commands exclusively.
- **Strict Relative Paths:** All file paths, markdown links, and subtask paths must be relative git paths from repository root; zero `file:///` URIs. Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
- **Positive Booleans:** Implicit positive booleans only (`isReady`, `hasCache`).

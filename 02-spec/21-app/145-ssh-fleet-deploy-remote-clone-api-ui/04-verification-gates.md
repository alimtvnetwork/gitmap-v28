# Spec 145: Verification Gates & Quality Acceptance

Spec Reference: [01-overview.md](01-overview.md)

---

## 1. Acceptance Criteria Index

| ID | Title | Verification Condition |
|---|---|---|
| **AC-SPEC145-01** | Macro Fleet Deployment | `gitmap macro deploy ssh --except <nodes>` broadcasts macros across online nodes with zero errors. |
| **AC-SPEC145-02** | Fleet Update Engine | `gitmap update --all` / `gitmap ua` aggregate remote JSON summaries and render an aligned terminal table. |
| **AC-SPEC145-03** | Software Inventory LS | `gitmap update ls` queries all nodes and displays multi-column software versions. |
| **AC-SPEC145-04** | Self-Healing Remote Clone | `gitmap ssh clone` detects missing auth access, injects keys, logs each step, and completes clone. |
| **AC-SPEC145-05** | VS Code Remote Integration | `gitmap vscode remote <node> [path]` constructs proper executable flags; `gitmap vscode remote fix` repairs stale server state. |
| **AC-SPEC145-06** | Bidirectional File Transfer | `gitmap ssh cp` copies files between local and remote nodes preserving content integrity. |
| **AC-SPEC145-07** | Embedded Web UI Server | `gitmap ui` and sub-module UI commands launch browser at available port and serve SPA cleanly. |
| **AC-SPEC145-08** | In-Browser Remote Text Editor | `gitmap editor ui <node> <file>` reads and saves remote files with syntax highlighting and zero corruption. |
| **AC-SPEC145-09** | Inter-Node REST Triad Daemon | Mutual authenticated REST endpoints serve ping/exec/update with transparent fallback to SSH execution. |
| **AC-SPEC145-10** | Coding Guidelines Compliance | All Go and script files pass boolean, nested-if, and newline styling linters with exit code 0. |

---

## 2. Quality Gate Verification Commands

```bash
# Boolean Guidelines Linter
python linter-scripts/check-boolean-guidelines.py

# Nested-If Linter (Changed Files)
python linter-scripts/check-nested-ifs.py --changed-only

# Newline Styling Linter (Unix LF)
python linter-scripts/check-newline-styling.py
```

# Architecture Spec: Cursor Memories, Conversations & Projects Migration to Ubuntu (u1)

## 1. Architectural Overview & Context
This specification details the end-to-end migration of Cursor IDE internal state—including AI memories, custom persona agents, skills, per-project conversation transcripts, global AI chat state (`state.vscdb`), conversation index (`conversation-search.db`), and workspace project configurations—from the local Windows development workstation to Ubuntu fleet node `u1` (`192.168.1.22`).

---

## 2. Asset Classification & Mapping Matrix

```
Windows Source                                          Ubuntu Target (POSIX Node u1)
────────────────────────────────────────────────────    ────────────────────────────────────────────────────
C:\Users\Administrator\.cursor\                         ~/.cursor/
├── agents/                                             ├── agents/
├── ai-tracking/                                        ├── ai-tracking/
├── projects/                                           ├── projects/
│   └── d-work-<repo>/                                  │   └── home-a-git-work-<repo>/
├── skills-cursor/                                      ├── skills-cursor/
├── argv.json                                           ├── argv.json
└── cli-config.json                                     └── cli-config.json

C:\Users\Administrator\AppData\Roaming\Cursor\User\     ~/.config/Cursor/User/
├── settings.json                                       ├── settings.json (Dracula + styling)
├── globalStorage/                                      ├── globalStorage/
│   ├── alefragnani.project-manager/projects.json       │   ├── alefragnani.project-manager/projects.json (POSIX paths)
│   ├── conversation-search.db                          │   ├── conversation-search.db (SQLite)
│   ├── state.vscdb                                     │   ├── state.vscdb (SQLite global AI chat state)
│   └── storage.json                                    │   └── storage.json
└── workspaceStorage/                                   └── workspaceStorage/
    └── <hash>/ (49 workspaces)                             └── <hash>/ (workspace.json + state.vscdb)
```

---

## 3. Path Sanitization & Invariant Preservation
During staging, text and configuration files are processed through deterministic regex normalization:

| Pattern (Source) | Replacement (POSIX Target) | Context |
|---|---|---|
| `(?i)[dD]:[/\\]work[/\\]([a-zA-Z0-9_\-]+)` | `/home/a/git-work/\1` | Workspace and repository roots |
| `(?i)d-work-([a-zA-Z0-9_\-]+)` | `home-a-git-work-\1` | Storage identifiers inside `.cursor/projects/` |
| `(?i)[cC]:[/\\]Users[/\\]Administrator[/\\]\.cursor` | `/home/a/.cursor` | Cursor home directory references |
| `(?i)[cC]:[/\\]Users[/\\]Administrator` | `/home/a` | User home directory references |
| `statsig-cache.json` | `{}` | Reset feature flags to prevent payload bloat |

---

## 4. SQLite Migration Strategy
For binary SQLite files (`state.vscdb`, `conversation-search.db`, and `workspaceStorage/<hash>/state.vscdb`):
- SQLite databases are transferred as binary assets.
- In `workspace.json`, internal workspace paths are sanitized to POSIX format under `/home/a/git-work/<repo>`.
- Database integrity checks (`PRAGMA integrity_check;`) are performed before and after extraction.

---

## 5. Non-Destructive Snapshots & Rollback Guarantee
Before extraction on remote node `u1`, the migration engine creates pre-flight backups:
- `/home/a/.config/Cursor/User.bak.<timestamp>/`
- `/home/a/.cursor.bak.<timestamp>/`

Active development repositories under `/home/a/git-work/` are protected from modification.
All extracted files and directories are normalized to:
- Directories: `0755` (`drwxr-xr-x`)
- Files: `0644` (`-rw-r--r--`)
- Ownership: User `a:a`

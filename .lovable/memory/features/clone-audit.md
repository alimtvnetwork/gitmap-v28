---
name: clone-audit
description: gitmap clone --audit reads a manifest, computes planned git clone/pull commands, and prints a diff-style summary without executing or touching the network.
type: feature
---

# Clone Audit Mode (v3.99.0)

`gitmap clone --audit <source>` is a read-only planner. It parses a
manifest (json/csv/text/path), runs the same branch-selection strategy
as the live cloner (`pickCloneStrategy`), and prints one diff-style row
per record:

| Marker | Action   | Trigger                                          |
|--------|----------|--------------------------------------------------|
| `+`    | clone    | target path missing                              |
| `~`    | pull     | target is a git repo (cache miss)                |
| `=`    | cached   | clone-cache fingerprint matches local HEAD       |
| `?`    | conflict | target exists but is not a git repo              |
| `!`    | invalid  | record has no HTTPSUrl/SSHUrl                    |

## Why
- Lets users review a manifest before a batch clone.
- CI dry-runs against generated `.gitmap/output/` files.
- Works offline; `requireOnline` and SSH-key resolution are skipped.

## Constraints
- Refuses direct git URLs — manifest only.
- Never invokes git, never writes outside stdout.
- Output format is stable (constants in `constants_clone_audit.go`)
  so downstream grep/awk pipelines can rely on it.

## Source map
- `gitmap/cloner/audit.go` — planner + report printer.
- `gitmap/cloner/audit_path.go` — stat-only existence helper.
- `gitmap/cloner/audit_test.go` — coverage for classification, command
  shape, marker mapping, and printer formatting.
- `gitmap/cmd/cloneaudit.go` — CLI dispatcher invoked from `runClone`
  before `requireOnline`.
- `gitmap/constants/constants_clone_audit.go` — every user-facing string.

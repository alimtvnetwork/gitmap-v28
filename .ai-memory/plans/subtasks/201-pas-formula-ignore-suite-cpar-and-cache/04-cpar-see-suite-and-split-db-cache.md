# Subtask 04: CPAR, Observability Suite, and Split-DB Repo Cache Engine

## Scope
- Files:
  - `cli/cmdcpar/` (`gitmap commit-push-all-repos`, `cpar`)
  - `cli/cmdsee/` (`gitmap see` suite, `ses`)
  - `cli/cmdcache/` (`gitmap cache` suite)
  - `cli/cmd/rootcore.go`, `cli/cmd/rootgit.go`, `cli/cmd/rootnodes.go`
- Objective:
  - Implement `cpar` with `-y`, `--review` (`-r`), `--review --commit-only` (`-co`).
  - Implement `see` subcommands: `commit pending`, `git-ignore`/`ig issues`, `errors`, `history`, `errors ssh` (`ses`), `repo-manage ui`.
  - Implement `cache` subcommands: `create`, `ls`, `add`, `remove` (`rm`), `help`, `search`, `search-multi`, `search-multi-grep`, `recache`/`reconcile`/`sync`.
  - Split-DB SQLite structure: root file index DB + folder slug DBs, indexing files <= 200KB.
  - Minor release bump and CI/CD verification.

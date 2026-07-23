---
Slug: reclone-folder-and-url-flow
Parent: 03-reclone-transport-and-vscode-open
Status: pending
Created: 2026-06-07
---

# Subtask 03 — Reclone: read → store → reuse → log

## Goal
Every reclone-class command (`cfr`, `cfrp`, `clone-now`/`reclone`, direct-URL `clone` when target already exists) must:
1. **Resolve** the repo path (cwd, explicit folder arg, or DB lookup by URL).
2. **Read** `.git/config` `remote.origin.url`.
3. **Classify** transport via the shared classifier (same one `scan` uses — extract to `gitmap/transport/classify.go` if not already shared).
4. **Persist** the classified transport via `UpsertRepoByPath` (subtask 02).
5. **Reuse** transport in the URL picker (mirror `pickURLForTransport` from v6.21) so the actual `git clone` argv matches the original transport.
6. **Log** the reclone event into the history table so `gitmap history` shows it.

## Resolve rules
- No arg + cwd inside a git repo → use repo root (walk up until `.git` found).
- Folder arg → `filepath.Abs(arg)`; error if not a git repo.
- URL arg → look up `Repo` by `HTTPSUrl`/`SSHUrl`, else clone fresh into default location (existing behavior).

## Read rules
- Prefer `git -C <path> config --get remote.origin.url` (handles worktrees + alternate config locations).
- Fall back to parsing `.git/config` directly only if the `git` invocation fails.

## Log rules
- Reuse the existing history table (`mem://features/history-rewrite` lists `gitmap history`); if no row schema exists for reclone events, add one in this subtask (`ReencloneHistory` table, columns: `Id`, `RepoId`, `Source`, `Target`, `Transport`, `Command`, `ExitCode`, `OccurredAt`).
- One row per reclone invocation; failed clones also logged with non-zero `ExitCode`.

## Error policy
- Per `mem://style/code-constraints` + `mem://tech/code-red-error-management`: every failure logged to `os.Stderr` with the standardized format. No swallowed errors. `errors.Is` for sentinels.

## Tests
- Fixture repo with SSH origin → `cfr` → assert emitted `git clone` command contains `git@…`, `Repo.IdentifiedTransport == "ssh"` after run, history row written.
- Fixture repo with HTTPS origin → same flow, transport `"https"`.
- Folder-arg form (`gitmap cfr ./some/repo`) → resolves and persists.
- URL-arg form against a DB-known repo → reuses stored transport.

## Definition of done
- All four invocation forms (cwd, folder arg, URL arg known, URL arg unknown) tested.
- `gitmap history` lists the reclone events.
- No reclone path bypasses the picker.

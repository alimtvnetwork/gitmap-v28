# Subtask 05 — Implement `merge write` / `mr` (program 273, WS3)

## File box

- `cli/cmdmerge/` — NEW package, ONE file (≤300 lines): dispatch entry,
  flags, GitHub host call, help text, example-registry entry.
- Dispatch table — ONE surgical edit: add `merge write` + `mr` alias.
- Touch nothing else (no pull/commit internals).

## Context

Spec: `02-spec/21-app/gitmap-help-text-and-ui-enhancement/05-merge-write-command.md`.
P0 for the Lovable migration: create a GitHub pull request from the CLI.
Flags: `--title` (required), `--body`, `--source` (default: current branch),
`--target` (default: repo default branch), `--draft`, `--yes`. Flow: resolve
repo+host from git remote (GitHub only v1 — other hosts get a clear
tracked error, never silent failure); reuse the existing credential path
(`gh` CLI if present, else the GitHub App token flow); print confirmation
summary, prompt unless `--yes`; on success print PR number + URL; on failure
surface the host error verbatim. Rules: gitmap search tools only; relative
paths; ≤300 lines; `go build ./...` from `cli/` — NEVER `go test`.

## Checklist

- [ ] `gitmap merge write --title ... --body ... --source ... --target ... --draft`
  opens a PR; output prints PR number + URL.
- [ ] `gitmap mr` alias identical behavior.
- [ ] Missing `--title` → tracked error naming the missing flag.
- [ ] Non-GitHub remote → clear "unsupported host" error (no panic).
- [ ] Help example entry registered (spec 03 format).
- [ ] Commit: `gitmap cpf "merge - implement merge write and mr alias"` + push
  (verify `git ls-remote`).

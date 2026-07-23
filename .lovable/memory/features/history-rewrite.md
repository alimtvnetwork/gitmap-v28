---
name: History Rewrite
description: gitmap history-purge (hp) and history-pin (hpin) wrap git filter-repo in a mirror-clone sandbox; working repo never rewritten in place (v4.17.0)
type: feature
---

`gitmap history-purge` / `hp` and `gitmap history-pin` / `hpin` wrap
`git filter-repo` in a temp `os.MkdirTemp` mirror-clone sandbox. The
user's working repo is never rewritten in place. Spec:
`spec/04-generic-cli/16-history-rewrite.md`.

## Pipeline (5 phases, both commands)

1. Identify origin via `git remote get-url origin` (cwd).
2. Mirror-clone into temp sandbox.
3. Run `git filter-repo` inside sandbox: purge uses `--invert-paths --path P`; pin uses `--blob-callback` with a JSON manifest of (path, current bytes, historical SHA set) per requested path.
4. Verify: purge -> `git log --all -- P` empty; pin -> every `git show <sha>:P` hashes to same SHA-256.
5. Push prompt: `--force-with-lease --mirror`. `--yes` skips prompt; `--no-push` short-circuits with manual command.

## Path parsing

Multi-form: separate args, comma, or comma-space all accepted. Folders and files both work. Quoting irrelevant. See `historyrewrite_paths.go`.

## Flags

`-y/--yes`, `--no-push` (mutually exclusive with `--yes`), `--dry-run`, `--message <s>` (rewrites touched-commit messages only, via `--commit-callback`), `--keep-sandbox`, `-q/--quiet`.

## Exit codes

`0` ok / `2` not-in-repo / `3` filter-repo-not-installed (with OS install hint) / `4` bad-args / `5` filter-repo-failed / `6` verify-failed / `7` push-failed.

## Files

- `gitmap/cmd/historyrewrite.go` + `_paths.go` + `_flags.go` + `_sandbox.go` + `_pin.go` + `_verify.go` + `_push.go`
- `gitmap/constants/constants_historyrewrite.go`
- `gitmap/helptext/history-purge.md`, `history-pin.md`
- `src/pages/HistoryRewrite.tsx` (`/history-rewrite`)
- `.github/workflows/history-rewrite-smoke.yml` + `.github/scripts/smoke-history-{purge,pin}.sh`

## Dependency

`git filter-repo` not bundled. Missing binary -> exit 3 with OS install hint (pip / brew / scoop). No auto-install.

## Safety guarantees (non-negotiable)

- Working repo never rewritten — all mutation in temp sandbox.
- Verification gates push prompt; failed verify cannot push.
- `--dry-run` short-circuits before push regardless of `--yes`.
- Sandbox path printed on every error path; `--keep-sandbox` keeps it.
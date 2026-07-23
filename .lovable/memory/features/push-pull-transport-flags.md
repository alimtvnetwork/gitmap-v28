---
name: push-pull-transport-flags
description: `gitmap push` (new) + `gitmap pull` accept `--ssh`/`--https` (aliases `-ssh`/`-https`/`--sh`/`--ht`) that PERSIST the rewrite via `git remote set-url origin` before invoking git.
type: feature
---

# `gitmap push` / `gitmap pull` `--ssh` / `--https` (v5.29.0+)

Adds `gitmap push` (cwd short-circuit, mirrors v5.28.0 `pull` cwd
behavior) and `--ssh` / `--https` transport flags on both `push` and
the existing `pull` cwd short-circuit.

## Contract

- Flags rewrite `remote.origin.url` of the cwd repo PERSISTENTLY via
  `git remote set-url origin <converted>` — not a transient
  per-invocation env override.
- Mutually exclusive; `--ssh` wins on conflict + one-line stderr warn.
- Aliases: `-ssh` / `--sh` and `-https` / `--ht`. Single-dash forms
  exist because Go's `flag` rejects them and we route them via
  `reorderFlagsBeforeArgs` + manual recognition in
  `parseTransportFlags`.
- Unrecognised URLs fail OPEN: print a warning, skip the rewrite,
  still run git push/pull.
- Extra positionals after the flags forward verbatim to git
  (`gitmap push --ssh origin main` ⇒ `git push origin main`).

## Wiring

- `gitmap/cmd/remotetransport.go` — `ApplyTransportFlag(dir, useSSH, useHTTPS)`.
- `gitmap/cmd/push.go` — `runPush`, `parseTransportFlags` shared with pull.
- `gitmap/cmd/pull.go::runPullCWD` — accepts the same flags.
- Dispatch in `rootcore.go`: `{CmdPush, CmdPushAlias}` → `runPush`.
- `CmdPush="push"`, `CmdPushAlias="ph"` (avoid clash with `p`=pull).
- E2E in `gitmap/cmd/pushpull_transport_e2e_test.go` (skipped when
  `git` missing from PATH).

## Spec

`spec/01-app/111-push-pull-transport-flags.md`.

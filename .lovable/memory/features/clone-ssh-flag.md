---
name: clone-ssh-flag
description: `gitmap clone --ssh` / `--https` force every URL into the requested transport before git runs; conversion via ConvertURLToSSH / ConvertURLToHTTPS; skips non-URL positionals.
type: feature
---

# `gitmap clone --ssh` / `--https` (v5.20.0+)

`gitmap clone --ssh <url>` rewrites the URL into its
`git@host:owner/repo.git` SSH-shorthand form before git is invoked.
`--https` is the symmetric counterpart. Both flags also flow through
the multi-URL form (`clone url1,url2,url3 --ssh`).

## Wiring

- `UseSSH` / `UseHTTPS` on `CloneFlags` in `gitmap/cmd/rootflags.go`.
- Conversion helpers in `gitmap/cmd/cloneurlconvert.go`:
  `ConvertURLToSSH(url)` and `ConvertURLToHTTPS(url)`. Both return
  `(string, bool)` — `ok=false` means the input wasn't a recognised
  Git URL and is returned unchanged so callers can fall through.
- Dispatch in `gitmap/cmd/clone.go::applyURLSchemeFlags` runs after
  `applySSHKey` and BEFORE the multi-URL / direct-URL routers, so the
  multi-URL detector sees the converted URLs.

## Contract

- `--ssh` and `--https` are mutually exclusive; when both are set
  `--ssh` wins + one-line stderr warning is printed.
- Non-URL positionals (folder names, `json` / `csv` / `text`
  shorthands) are passed through unchanged — gated by the same
  `isDirectURL` predicate the multi-URL detector uses.
- Port hints in `ssh://` URLs are dropped (SSH-shorthand has no port).
- Already-shorthand URLs are normalized (`.git` suffix appended).
- Manifest-mode clones (`json` / `csv` / `text`) currently IGNORE the
  flags — those formats already carry both `httpsUrl` and `sshUrl`
  columns and select via `mode` upstream. Honouring `--ssh` / `--https`
  as a per-record override in `cloner` is tracked as follow-up in
  `spec/01-app/110-clone-ssh-flag.md`.

Spec: `spec/01-app/110-clone-ssh-flag.md`.

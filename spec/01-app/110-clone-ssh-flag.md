# `gitmap clone --ssh` / `--https` (URL scheme coercion)

## Purpose

Let users force every clone URL into a specific transport — SSH-shorthand
(`git@host:owner/repo.git`) or HTTPS (`https://host/owner/repo.git`) —
without manually rewriting the URL on the command line. The flag also
flows through to the multi-URL form so `gitmap clone url1,url2,url3 --ssh`
converts the whole batch in one shot.

## Synopsis

```
gitmap clone <url|json|csv|text|path> [folder] [--ssh | --https] [flags]
```

- `--ssh` — rewrite every recognised Git URL into `git@host:owner/repo.git`
  SSH-shorthand. HTTPS URLs (`https://github.com/owner/repo`) and
  `ssh://git@host[:port]/owner/repo` URLs are both converted.
  Already-shorthand URLs are normalized (`.git` suffix appended).
- `--https` — symmetric counterpart, forces every URL into
  `https://host/owner/repo.git` form.
- The two flags are mutually exclusive — when both are set, `--ssh`
  wins and a one-line stderr warning is printed.

## Conversion contract

| Input shape                              | `--ssh` output                       | `--https` output                     |
|------------------------------------------|--------------------------------------|--------------------------------------|
| `https://github.com/owner/repo`          | `git@github.com:owner/repo.git`      | `https://github.com/owner/repo.git`  |
| `https://github.com/owner/repo.git`      | `git@github.com:owner/repo.git`      | `https://github.com/owner/repo.git`  |
| `http://gitlab.example/owner/repo`       | `git@gitlab.example:owner/repo.git`  | `https://gitlab.example/owner/repo.git` |
| `git@github.com:owner/repo.git`          | (unchanged)                          | `https://github.com/owner/repo.git`  |
| `ssh://git@github.com:22/owner/repo.git` | `git@github.com:owner/repo.git`      | `https://github.com/owner/repo.git`  |
| Non-URL token (folder name, `json`)      | (unchanged)                          | (unchanged)                          |

Port hints in `ssh://` URLs are intentionally dropped — SSH-shorthand
has no port slot. Users that need a non-default SSH port should pin it
in `~/.ssh/config` and continue to use the explicit `ssh://` URL.

## Implementation

- Flag wiring lives in `gitmap/cmd/rootflags.go` (`UseSSH` + `UseHTTPS`
  on `CloneFlags`, both registered against the existing clone
  `flag.FlagSet`).
- Conversion helpers live in `gitmap/cmd/cloneurlconvert.go`:
  `ConvertURLToSSH(url)` and `ConvertURLToHTTPS(url)`, both returning
  `(string, bool)` so callers can detect unrecognised inputs.
- Dispatch lives in `gitmap/cmd/clone.go` (`applyURLSchemeFlags` runs
  after `applySSHKey` and before the multi-URL / direct-URL routers).
  Non-URL positionals (`json`, folder names) are skipped via the same
  `isDirectURL` predicate used for the multi-URL detector — so a stray
  `--ssh` cannot corrupt a manifest-style invocation.

## Behaviour notes

- The flag converts the URL **before** the multi-URL detector runs, so
  `clone url1,url2 --ssh` works exactly the same as
  `clone url1 url2 --ssh` and `clone url1 --ssh` (single-URL form).
- When a URL is rewritten, a `↪ --ssh rewrite: <before> → <after>`
  breadcrumb is printed to stdout so the user can verify the
  substitution before git runs.
- Manifest-mode clones (`json` / `csv` / `text`) ignore `--ssh` /
  `--https` for now — those formats already carry both `httpsUrl` and
  `sshUrl` columns and are selected by `mode=https|ssh` upstream in
  `gitmap scan`. Future work: honour `--ssh` / `--https` as a
  per-record override in the cloner. Tracked in
  `spec/01-app/110-clone-ssh-flag.md` (this file) as a follow-up.

## Related

- `spec/01-app/96-clone-replace-existing-folder.md` — replace-on-clone flow
- `spec/01-app/104-clone-multi.md` — multi-URL batch form
- `mem://features/navigation-helper` — `cn` / `cd` shell handoff

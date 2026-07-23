# Issue 02 — Reclone (`cfr` / `cfrp` / `cn`) loses SSH transport

**Status:** open
**Created:** 2026-06-07
**Related:** `.lovable/issues/01-ssh-repo-cloned-as-https.md`, `.lovable/spec/commands/04-reclone-honors-stored-transport.md`

## Symptom
When a repo was originally cloned over SSH (`origin = git@github.com:…`), running `gitmap clone-fix-repo` / `cfrp` / `clone-now` re-clones it over HTTPS, triggering the browser-auth prompt on private remotes.

## Repro
1. `cd` into a repo whose `origin` is `git@github.com:owner/repo.git`.
2. Run `gitmap cfr` (or `cfrp`, or `cn` against a stored entry).
3. Observe the printed `git clone` command uses `https://github.com/owner/repo.git`.
4. On private remotes: browser-auth prompt fires.

## Expected
- Reclone reads `.git/config` `remote.origin.url`, classifies transport, persists it on the `Repo` row, and reuses SSH for the actual `git clone`.
- History (`gitmap history`) shows the reclone event with `transport: ssh`.

## Actual
- URL picker in `clonefixrepo.go` / `clone.go` / `clonenow.go` falls back to HTTPS-first.
- No `IdentifiedTransport` column on `Repo` → nothing to remember even after a successful manual fix.
- Reclone events are not logged into the history table.

## Related files (suspected)
- `gitmap/cmd/clonefixrepo.go`
- `gitmap/cmd/clone.go`
- `gitmap/cmd/clonenow.go`
- `gitmap/db/*` (schema migration 007 for `IdentifiedTransport`)
- `gitmap/cloner/summary.go` (`pickURL`)

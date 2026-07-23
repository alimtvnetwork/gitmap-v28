# Command: reclone / clone-fix-repo must honor identified transport

**Status:** open
**Created:** 2026-06-07

## Verbatim
> If we are in the repo or point a folder, then use a reclone command, it will also check how the repo is already cloned. And based on that, it would clone it. If it was done with SSH, it would remember that and also save it in its SQLite DB, and also reclone using SSH. The reclone can be done with a URL. Reclone can be done when we are inside a folder. Reclone can be done when we point to a folder — then it would go into the .git folder inside it and try to find what the repo path is, and also log everything. We could see the histories in our gitmap history.

## Scope
- Applies to: `clone-fix-repo` (`cfr`), `clone-fix-repo-pub` (`cfrp`), `clone-now` / `reclone` (`cn`), and any future reclone variant.
- Trigger forms:
  1. inside a git repo (cwd = repo root or any subdir)
  2. explicit folder argument pointing at a repo
  3. direct URL argument (HTTPS or SSH)
- Behavior:
  - Read `.git/config` `remote.origin.url`, classify as `ssh` / `https` / `other` (same classifier as scan).
  - Persist the classified transport on the `Repo` row (column `IdentifiedTransport`).
  - When recloning, the URL picker must prefer the stored/identified transport over the scan-wide `--mode` (mirrors v6.19–v6.22 rule).
  - Log every reclone event (source, target, transport, command, exit) so it shows up in `gitmap history`.

## When it applies
Every reclone-class command, on every platform. No flag opt-in — this is the default.

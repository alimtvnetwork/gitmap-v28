# Reclone-class URL picker audit — 2026-06-21

Plan 03 step 1: classify every reclone-class URL picker as ✅ honors-identified-transport / ⚠️ HTTPS-first / ❌ broken. File:line evidence below.

## Verdicts

| Command | Entry file:line | Picker | Verdict | Notes |
|---|---|---|---|---|
| Batch `clone-now` / `reclone` / `clone-pick` (scan-driven) | `gitmap/cloner/cloner.go:178` | `pickURL(rec)` @ `gitmap/cloner/summary.go:41` | ✅ Honors `rec.Transport == "ssh"` and returns `SSHUrl`; falls back to HTTPS only when SSH is empty. |
| Cache freshness check (batch) | `gitmap/cloner/cache.go:131,178` | `pickURL(rec)` | ✅ Same helper — cache key and freshness compare use the transport-aware URL, so a transport flip invalidates the cache row. |
| Per-row audit (batch) | `gitmap/cloner/audit.go:87` | `pickURL(rec)` | ✅ Same helper. |
| Direct-URL `clone <url>` | `gitmap/cmd/clone.go` (executeDirectClone) | URL is user-supplied verbatim; `--ssh` / `--https` flags coerce via `ConvertURLToSSH` / `ConvertURLToHTTPS`. | ✅ User intent wins; no implicit downgrade. |
| `cfr` / `cfrp` (direct-URL reclone-class) | `gitmap/cmd/clonefixrepo.go:71` | `preferExistingFolderTransport(url, absPath)` @ `gitmap/cmd/clonefixrepofoldertransport.go:26` | ✅ When destination already contains `.git/`, reads `remote.origin.url` via `gitutil.RemoteURL`, classifies via `isSSHURL`, and rewrites the positional URL to match before clone. No-ops on fresh clones (intentional — user URL is source of truth). |
| `cfr` / `cfrp` (fresh clone, no `.git/`) | `gitmap/cmd/clonefixrepo.go:73` → `executeDirectClone` | n/a | ⚠️ No DB lookup. If the repo was previously cloned and the folder was deleted, the stored transport on `Repo` is ignored — user must pass `--ssh` / `--https` themselves. **This is the gap Plan 03 step 2 (`IdentifiedTransport` column) + step 3 (DB-driven reuse) closes.** |
| `clone-replace` temp-swap fallback | `gitmap/cmd/clonereplace.go` | Reuses positional URL + `preferExistingFolderTransport` via the shared `runCloneCommandPretty` runner | ✅ Inherits the cfr fix transitively (v6.49.0+ unified runner). |
| `clone-next` | `gitmap/cmd/clonenext.go` | Resolves next `-vN+1` URL by string-rewriting the **current** origin URL, preserving its transport. | ✅ Transport-preserving by construction. |

## Gap summary (drives Plan 03 step 2 + step 3)

Only **one** picker class is still HTTPS-first by omission: `cfr` / `cfrp` on a **fresh** clone where no on-disk `.git/` exists but the `Repo` row in `gitmap.db` remembers a prior SSH transport. The fix requires:

1. **Step 2** — `Repo.IdentifiedTransport TEXT NOT NULL DEFAULT ''` (next free migration is `008`; `007` is taken by commit-in tag-replay per `gitmap/constants/constants_commitin_tagreplay_sql.go:3`).
2. **Step 3** — before `executeDirectClone`, resolve the `Repo` row by canonical URL slug; if `IdentifiedTransport='ssh'` and positional URL is HTTPS, coerce via `ConvertURLToSSH` and log a history event analogous to `MsgCFRFolderTransport`.

## Cross-references

- Spec: `.lovable/spec/commands/04-reclone-honors-stored-transport.md`
- Plan: `.lovable/plans/pending/03-reclone-transport-and-vscode-open.md` (step 1 ← this audit)
- Prior audit: `.lovable/audits/2026-06-07-reclone-pickers.md` (referenced by `clonefixrepofoldertransport.go:14`)
- Memory: `mem://features/reclone-transport-handling` (none yet — create when step 3 lands)

---
name: v15 Legacy Compat Shim Audit (v3.12.1)
description: Decision record for the pre-v15 draft/preRelease JSON+SQLite backward-compat shim — keep through v3.x, drop in v4.0
type: feature
---

# v15 Legacy Compat Shim Audit — v3.12.1

**Date:** 2026-04-20 (Malaysia time)
**Scope:** `gitmap/release/metadata.go::ReadReleaseMeta` JSON overlay + `gitmap/store/migrate_v15phase5.go` SQLite column rename
**Decision:** **KEEP both shims through the v3.x line. Schedule removal for v4.0.0.**

---

## What the two shims do

### 1. `release/metadata.go::ReadReleaseMeta` — JSON overlay (lines 156-167)

After the primary `json.Unmarshal` into `ReleaseMeta` (which writes the new `isDraft` / `isPreRelease` keys), a second unmarshal into a tiny anonymous struct reads the **legacy** keys `draft` / `preRelease`. If the new field is still zero **and** the legacy pointer is non-nil **and** points to `true`, the new field is set. Re-saved (post-v15) files are never downgraded because the legacy-pointer-is-nil guard skips the overlay.

### 2. `store/migrate_v15phase5.go` — SQLite column rename (lines 22-63)

Detects a pre-v15 `Release` table with `Draft` / `PreRelease` columns and rebuilds it into the canonical v15 shape with `IsDraft` / `IsPreRelease`. Idempotent: bails out early if `Release` doesn't exist (fresh install) or if `Draft` column is already gone (already migrated). Triple-guarded against clobbering when both old and new columns coexist.

---

## On-disk audit (this repo, v3.12.1)

| Metric | Value |
|---|---|
| Total release JSON files in `.gitmap/release/` | 74 |
| Files using new `isDraft` / `isPreRelease` keys | **0** |
| Files using legacy `draft` / `preRelease` keys | ~10+ |
| Files where `"draft": true` | **0** |
| Files where `"preRelease": true` | **0** |

**Implication:** for every release ever cut from this repo, the JSON overlay is currently a **no-op** — the legacy field is always present-but-false, so the `*legacy.Draft` dereference resolves to `false` and the override branch never fires. The shim's *behavior* in this repo is identical to deleting the entire overlay block.

## Why we're keeping it anyway

1. **External / forked installs are unobservable.** The v3.4.x → v3.5.0 transition shipped to whoever was running gitmap at the time. Any of those users could have a `.gitmap/release/v3.4.x.json` with `"draft": true` (e.g. a user who cut a draft release before the `IsDraft` rename landed). Removing the JSON overlay would silently flip those drafts to "published" on next read.
2. **The SQLite migration is non-negotiable.** A user upgrading from v3.4.x to v3.12.x runs the binary against an existing `gitmap.db` that physically has `Draft` / `PreRelease` columns. Without `migrateV15Phase5`, every read against `Release` errors with `no such column: IsDraft`. Removing this is a hard data-loss bug.
3. **Both shims are cheap.** JSON overlay is ~12 lines, runs once per `ReadReleaseMeta` call, no I/O. SQLite migration is one-shot per install (guarded by `columnExists` early-return), runs zero times after first successful upgrade.
4. **`v3.12.x` is a minor line.** Per project policy (`.lovable/user-preferences` line 8: *"Code changes must bump at least minor version"*) and standard semver, we cannot drop backward compat in a patch or minor release. The earliest defensible removal point is **v4.0.0**.

## Removal plan (deferred to v4.0.0)

When v4.0.0 is cut:

1. **Delete `gitmap/store/migrate_v15phase5.go`** entirely. The v3.x line will have run the migration on every active install for >12 months by then; any database still on the v3.4.x shape is by definition abandoned. If a holdout user surfaces, the v3.12.x binary remains downloadable from the GitHub release page and acts as a one-shot upgrade tool (`gitmap` v3.12.x → run once → upgrade to v4.0.x).
2. **Delete the JSON legacy-overlay block** in `gitmap/release/metadata.go::ReadReleaseMeta` (lines 153-167 plus the accompanying comment). The function reduces to a single `json.Unmarshal` + return.
3. **Add a v4.0 migration note** to `CHANGELOG.md` under "Removed": *"Pre-v15 SQLite `Release.Draft` / `Release.PreRelease` column rename and JSON `draft` / `preRelease` key compat overlay. Users on v3.4.x or earlier must run any v3.5.0+ binary once before upgrading to v4.0."*
4. **Remove `migrateV15Phase5()` from the migration chain** in `gitmap/store/migrations.go` (or wherever it's invoked).
5. **Drop the v15 phase tests** that exercise the legacy code path, if any exist under `gitmap/store/*_test.go`.

## What does NOT need to change in v3.12.x

- `release.Version.PreRelease` (semver suffix on `v1.2.3-rc.1`) is a **completely unrelated field** on a different struct. It is not legacy and stays.
- The `--draft` CLI flag (`gitmap release --draft`) is the user-facing flag. Stays.
- Help-text / message strings referencing "DRAFT" column headers in `gitmap ls releases` output. Stay.
- All v15 phase 1.2 / 1.3 / 1.4 migrations. They are independent of this audit and remain required for v3.x → v4.x as well (audit them separately when planning v4.0).

## Cross-reference

- Spec: this audit should eventually be summarized into `spec/01-app/` as part of the v4.0 breaking-change matrix when that doc is created.
- Related migration: `gitmap/store/migrate_v15phase4.go` (the `Id` → `{Table}Id` PK rename) is the sibling phase that ships in the same v15 generation. Audit it on the same v4.0 schedule.

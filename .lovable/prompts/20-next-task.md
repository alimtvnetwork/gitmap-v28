# Next 20 Task — Plan 03 Step 2 still pending (migration 007 + `Repo.IdentifiedTransport`)

v6.27.0 stamped the prompt but did not execute the migration. This prompt re-queues the same Step 2 work and bumps to v6.28.0 for the planning artifact.

## Next step (N=1)
Implement Plan 03 Step 2: schema migration 007 `ALTER TABLE Repo ADD COLUMN IdentifiedTransport TEXT NOT NULL DEFAULT ''` (idempotent via `PRAGMA table_info`), extend `model.Repo`, update `UpsertRepoByPath` + `SelectAll` / `SelectBySlug` / `SelectByPath`, lazy backfill from URL prefix on first read.

## Why now
Plan 03 Step 3's history-log half and `clone-now`/`reclone` durable memory both block on this column. Without it the v6.26.0 `cfr`/`cfrp` fix only works when the destination `.git/` already exists locally.

# Next 19 Task — Plan 03 Step 2: persist `IdentifiedTransport` on `Repo` (migration 007)

Saved prompt for the v5 "Next N Steps" run that produced v6.27.0.

## Next step (N=1)
Implement Step 2 of `.lovable/plans/pending/03-reclone-transport-and-vscode-open.md` — schema migration 007 + `model.Repo.IdentifiedTransport` + `Select*`/`UpsertRepoByPath` updates + lazy backfill.

## Why now
Step 3's "read → store → reuse → log" half cannot land without the column to write into. The `cfr`/`cfrp` half of step 3 already shipped in v6.26.0 but the durable memory is still missing, so the same SSH→HTTPS downgrade re-appears across sessions.

## Status snapshot
- v6.27.0 cut for the planning/prompt artifact (this file + plan accounting).
- Plan 03 remaining: Step 2 (this prompt), Step 3 history-log half, Step 4 `gitmap code`, Step 5 closeout.

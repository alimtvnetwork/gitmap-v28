---
Slug: db-identifiedtransport-column
Parent: 03-reclone-transport-and-vscode-open
Status: pending
Created: 2026-06-07
---

# Subtask 02 — Persist `IdentifiedTransport` on the `Repo` row

## Goal
Add a durable column so reclone (and any future consumer) can remember how a repo was originally cloned, even across sessions.

## Schema migration (007)
- File: `gitmap/db/migrations/007_repo_identified_transport.go` (or the package's existing migration registration pattern — match what 006 used).
- SQL:
  ```sql
  ALTER TABLE Repo ADD COLUMN IdentifiedTransport TEXT NOT NULL DEFAULT '';
  ```
- Idempotent guard (skip if column already exists via `PRAGMA table_info`), per `mem://tech/database-architecture`.

## Model + accessor updates
- `model.Repo` gains `IdentifiedTransport string`.
- `UpsertRepoByPath` accepts the new field; existing callers pass `""` (no behavior change) — only the reclone path (subtask 03) and `scan` will set a non-empty value.
- `SelectAll`, `SelectBySlug`, `SelectByPath` — extend the `SELECT` column list + `Scan` targets. Per `mem://tech/database-constraints`, re-query IDs after upsert.

## Backfill
- On first read after the migration, if `IdentifiedTransport == ""` AND `HTTPSUrl/SSHUrl` are populated, derive from URL prefix (`git@` → `ssh`, `https://` → `https`) and persist back. Lazy, not a one-shot migration script.

## Tests
- Migration runs twice without error.
- Round-trip: upsert with `"ssh"`, select, value preserved.
- Backfill: insert row with `IdentifiedTransport=""` and `SSHUrl=git@…`, then `SelectByPath` returns `"ssh"` and the row on disk is updated.

## Definition of done
- All three `Select*` helpers return the column.
- `gitmap scan` populates the column on every row it touches (piggyback on the existing classifier).
- No call site references `Repo` struct fields by positional index (guard against silent breakage).

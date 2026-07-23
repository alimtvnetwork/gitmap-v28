# Subtask 01 — Schema + model for identified transport

**Parent:** 02-ssh-aware-clone
**Status:** pending
**Created:** 2026-06-07

## Goal
Persist identified transport and ensure both URL variants are always populated when derivable.

## Work
- Add `IdentifiedTransport string` (`ssh` | `https` | `other`) to `model.ScanRecord`.
- Mapper: when only one URL form is present on `origin`, synthesize the counterpart (github/gitlab/bitbucket host known patterns) so both `HTTPSUrl` and `SSHUrl` are stored.
- Set `IdentifiedTransport` from the actual `origin` URL classification (reuse `classifyTransport`).
- Store migration: add `IdentifiedTransport` column to `Repos` table; bump `SchemaVersionCurrent`; update `SQLUpsertRepoByPath` + `SQLSelectAll/BySlug/ByPath` + `scanOneRow`.
- Unit tests: SSH-only origin yields populated HTTPS counterpart and vice versa; identified transport preserved across upsert.

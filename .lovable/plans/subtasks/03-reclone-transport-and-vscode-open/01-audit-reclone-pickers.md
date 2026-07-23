---
Slug: audit-reclone-pickers
Parent: 03-reclone-transport-and-vscode-open
Status: completed
Created: 2026-06-07
---

# Subtask 01 — Audit reclone URL pickers

## Goal
Produce a verdict table for every reclone-class command's URL selection: does it honor identified transport, or does it hard-code HTTPS-first? Output is a checked-in markdown file at `.lovable/audits/2026-06-07-reclone-pickers.md` with file:line evidence.

## Files to inspect
- `gitmap/cmd/clonefixrepo.go` — both `cfr` and `cfrp` entrypoints; locate the picker that builds the actual `git clone` argv.
- `gitmap/cmd/clone.go` — direct-URL clone path; check if transport is inferred from the URL prefix itself.
- `gitmap/cmd/clonenow.go` — `clone-now` / `reclone` dispatch.
- `gitmap/cloner/summary.go` — shared `pickURL` already fixed in v6.20; confirm reclone callers route through it.
- `gitmap/cloner/*.go` — any other shared helpers.

## Per-command verdict shape
For each command write one row:
```
- <cmd> (<file>:<line>) — <honors|ignores> transport — <one-sentence justification>
```

## Definition of done
- Audit file committed under `.lovable/audits/`.
- Every "ignores" row links to the exact source line and proposes the one-line picker swap (mirror v6.21's `pickURLForTransport`).
- Plan step 3 implementation can be driven directly off this table with no re-reading.

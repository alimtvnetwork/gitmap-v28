# SSH-aware scan, persistence, and clone behavior

**Slug:** ssh-aware-clone
**Steps:** 5
**Status:** pending
**Created:** 2026-06-07

## Context
Scanning two SSH repos in D:\Work surfaced two defects: (1) only HTTPS URLs are persisted/used, so SSH repos get pulled/probed over HTTPS and trigger a browser auth prompt; (2) the new `make-all-public` / `make-all-private` commands are not reaching the docs UI cleanly. Fix the data model + every consumer to respect the identified transport, then deliver the docs/UI for the new commands, then cut a minor release.

Captured inputs:
- Issue: `.lovable/issues/01-ssh-repo-cloned-as-https.md`
- Command: `.lovable/spec/commands/03-respect-identified-transport.md`
- Existing pending: `.lovable/plans/pending/01-bulk-visibility-mapub-mapri.md` (carry forward — its UI delivery overlaps with step 4).

## Steps
1. Persist both URL variants + identified transport in scan/model/store. See ./subtasks/02-ssh-aware-clone/01-schema-and-model.md
2. Update every URL consumer (formatter, clone scripts, probe, pull, terminal report) to honor identified transport. See ./subtasks/02-ssh-aware-clone/02-consumers-honor-transport.md
3. Update terminal scan report to show `transport:`, `https:`, `ssh:`, and `command:` (identified transport) per repo; refresh affected golden fixtures.
4. Finish docs/UI delivery for `make-all-public` / `make-all-private` (routes, sidebar, search, dedicated pages) and add a "Transport awareness" note to the scan docs page. See ./subtasks/02-ssh-aware-clone/03-docs-and-ui.md
5. Bump minor: `gitmap/constants/constants.go` Version → 6.19.0, `src/constants/index.ts` VERSION → "v6.19.0", add `## v6.19.0` entry to `CHANGELOG.md`, pin `v6.19.0` in `README.md`. Do NOT touch `.gitmap/release/`.

## Verification
- `gitmap scan` on a repo with SSH origin shows SSH `command:` line and SSH-form `clone.ps1`/probe — no browser auth prompt.
- SQLite `Repos` row for that repo has BOTH `HTTPSUrl` and `SSHUrl` populated and `IdentifiedTransport='ssh'`.
- `go test ./...` green (including refreshed goldens + version-sync test).
- Docs preview: `/docs/make-all-public` and `/docs/make-all-private` render; both appear in sidebar + search.
- `README.md` shows v6.19.0; `CHANGELOG.md` has the new entry.

## Appended from prior pending tasks
- `01-bulk-visibility-mapub-mapri.md` — UI/data wiring for the bulk visibility commands is absorbed into step 4 of this plan.

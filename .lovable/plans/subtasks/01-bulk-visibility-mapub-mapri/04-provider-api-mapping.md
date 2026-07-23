---
Slug: provider-api-mapping
Status: pending
Created: 2026-06-06
Parent: 01-bulk-visibility-mapub-mapri
---

# Subtask 04 — Provider repo-list + visibility API

## Provider detection
Reuse `visibilityresolve.go::ResolveProvider`. Add `ResolveOwnerOnly(owner-or-url)`
that returns `(provider, owner)` without requiring a repo segment.

## Repo list

### GitHub (`gh`)
```
gh repo list <owner> --limit 4000 --json name,visibility --jq '.[].name'
```
Auth pre-check: `gh auth status` (existing helper `mustEnsureProviderCLI`).
Pagination: `--limit 4000` (gh paginates internally up to that cap).

### GitLab (`glab`)
```
glab repo list --group <owner> --output json
```
or `--user <owner>` when the owner is a user not a group. Detect by first
trying `--group`; on `404 group not found`, retry with `--user`.

## Apply visibility (per repo, sequential)
Reuse existing `visibilityapply.go::applyVisibility(ctx, repoFullName, target)`
called inside a loop. On success → `ResultStatusSucceeded`; on error →
`ResultStatusFailed` + log to `os.Stderr` per error-mgmt rule (Core memory:
`fmt.Fprintf(os.Stderr, "✗ %s: %v\n", repo, err)`), continue to next repo
(do NOT abort the batch — one failure should not kill the rest).

## Concurrency
Default: sequential (safe for rate limits + readable output).
Future: add `--parallel N` flag — out of scope for v1.

## Verification
- Dry-run mode (`--dry-run`): list what WOULD be flipped, do not call API,
  persist results as `Skipped` with note "dry-run". (Optional v1.1.)
- Smoke test: run against a throwaway personal repo.

# `gitmap hd` Hosted-Docs Fallback (No More Hard-Exit on Missing docs-site)

**Status:** Fixed in v6.0.x
**Affects:** `gitmap help-dashboard` (`hd`) on any install where `docs-site/`
is absent AND `docs-site.zip` cannot be auto-downloaded.

**Supersedes the user-visible failure mode of:**
- spec/02-app-issues/26-docs-site-not-bundled-and-swallowed-errors.md (Issue 1)
- spec/02-app-issues/27-docs-site-not-deployed-by-runscripts.md

> **Number history:** this spec is numbered **34**, not 33. Number 33 was
> previously occupied by `spec/02-app-issues/33-stale-binary-clone-folder-url-guard.md`
> (see `CHANGELOG.md` v3.95.0 entry, ~line 2230). That file was deleted
> upstream but the CHANGELOG link is permanent, so number 33 is
> reserved-historical and MUST NOT be reused. Future specs continue at 35+.

---

## Symptom (before this fix)

```
> gitmap hd
  ✗ Could not auto-download docs-site.zip (tried 2 source(s)). Last error: http 404
  ✗ Docs site directory not found at <binaryDir>\docs-site
```

Exit code 1. User had no usable docs at all when a release shipped without
`docs-site.zip` (older versions, hot-fix releases, or self-builds that
skipped the bundling step).

## Root Cause

`cmd/helpdashboard.go` treated the three failure modes — missing dir,
failed download, failed extract — as fatal (`os.Exit(1)`). The fixes in
issues #26 and #27 addressed the **bundling/install** side but left the
**runtime** side brittle: any release without `docs-site.zip` killed `hd`.

## Fix

`runHelpDashboard` now calls `openHostedDocsFallback()` instead of
`os.Exit(1)` whenever the local docs site cannot be materialized:

1. `downloadDocsSiteArchive` fails → fallback.
2. `extractDocsSiteZip` fails → fallback.
3. `docs-site/` still missing after the extract path → fallback.

`openHostedDocsFallback` prints `MsgHDHostedFallback` (the URL) to
stderr and launches the OS default browser at `constants.DocsURL`
(`https://gitmap.dev/docs`). The URL is always printed first, so the
fallback still works if `start` / `open` / `xdg-open` is unavailable.

The existing happy paths (local `docs-site/dist/` → `serveStatic`;
local `docs-site/` source → `serveDev`) are unchanged.

## Files

- `gitmap/cmd/helpdashboard.go` — `runHelpDashboard`, new helpers
  `openHostedDocsFallback` + `openURL` (extracted from `openBrowser`).
- `gitmap/constants/constants_helpdashboard.go` — `MsgHDHostedFallback`.
- `gitmap/constants/constants_messages.go` — `DocsURL` (reused, unchanged).

## Verification Checklist

- [x] `go vet ./cmd/... ./constants/...` clean.
- [x] `go build ./...` clean.
- [x] Unit pin: `constants/constants_helpdashboard_test.go`
      (`TestHostedDocsFallbackContract`) guards `DocsURL` +
      `MsgHDHostedFallback` format string.
- [x] Runtime pin: `cmd/helpdashboard_fallback_test.go`
      (`TestOpenHostedDocsFallbackPrintsURL`,
      `TestOpenURLNonFatalOnMissingLauncher`) — verifies the URL is
      written to stderr BEFORE the browser launch attempt, and that a
      missing OS launcher (`start`/`open`/`xdg-open`) is non-fatal.
- [ ] On a machine without `docs-site/` and a release whose
      `docs-site.zip` 404s: `gitmap hd` prints the hosted URL and
      opens the browser, exit code 0. *(manual; covered behaviorally
      by the two runtime pins above)*
- [ ] Happy path (`docs-site/dist/` present): still serves locally
      on port 5173 (unchanged). *(manual smoke test)*

---
name: cross-platform-install-update
description: Single source of truth for Windows/macOS/Linux gitmap install, update, uninstall, verify. Lives at spec/01-app/108-cross-platform-install-update.md and is rendered at /install-gitmap.
type: feature
---

# Cross-Platform Install / Update Reference (v3.100.0)

The README, the in-app `/install-gitmap` page, and `--help` text all
reference the same canonical matrix at
`spec/01-app/108-cross-platform-install-update.md`.

## Surfaces wired

- Spec: `spec/01-app/108-cross-platform-install-update.md`
- React page: `src/pages/InstallGitmap.tsx` (route `/install-gitmap`)
- Sidebar entry: "Install / Update gitmap" right under "Getting Started"
- README: top-of-file callout linking to the spec + page
- Existing helptext kept (`self-install.md`, `update.md`,
  `self-uninstall.md`) — they're now leaves of this canonical matrix.

## What lives in the matrix

For each of `install (default | prompt | pinned)`, `update`,
`uninstall`, `verify`, both PowerShell and bash/zsh one-liners exist.
PATH-activation modes (`auto`, `both`, `zsh+pwsh`, …) are listed once
with the profiles each touches.

## Why a separate page from `/install`

`/install` documents `gitmap install <tool>` (the third-party tool
installer for node, postgres, etc.). `/install-gitmap` documents how to
install gitmap itself. Distinct concerns, distinct pages.

## Update contract

`gitmap update` falls back in this order:
1. Linked source repo → `git pull` + build.
2. `gitmap-updater` binary → release asset download.
3. Manual one-liner panel from `gitmap/helptext/update.md`.

Phase 3 cleanup always writes a durable handoff log
(`<TMP>/gitmap-update-handoff-YYYYMMDD.log`) so failures stay
recoverable on Windows where stdout can be swallowed.

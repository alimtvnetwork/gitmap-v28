# Acceptance Criteria — gitmap-help-text-and-ui-enhancement

**Version:** 1.0.0
**Updated:** 2026-10-16

---

## Overview

33 testable criteria across 6 areas. Every criterion is verifiable by a
follow-up AI without asking the user.

---

## Area 1: Spec structure (health)

| # | Criterion | Source |
|---|-----------|--------|
| AC-001 | `readme.md` present with metadata, keywords, scoring, file table, cross-refs | `04-required-files.md` |
| AC-002 | `99-consistency-report.md` present with 100/100 health | `04-required-files.md` |
| AC-003 | All files lowercase kebab-case, unique numeric prefixes | `03-naming-conventions.md` |
| AC-004 | All cross-references relative with `.md` extension, zero broken | `10-cross-references.md` |
| AC-005 | Every spec file ≤300 lines | program 273 hard rules |

## Area 2: Theme design (WS2)

| # | Criterion | Source |
|---|-----------|--------|
| AC-006 | All theme palettes trace to coding-guidelines-v24 catalogue/tokens; zero invented colors | `02-theme-design.md` |
| AC-007 | Registry exposes `navy-purple`, `tokyo-night`, `vscode-dark-plus`, `monokai-pro` via `--theme` | `02-theme-design.md` |
| AC-008 | Legacy `bright`/`standard`/`monochrome` still accepted (aliases kept) | `02-theme-design.md` |
| AC-009 | Single registry table; docs UI help screens render from the same data | `02-theme-design.md` |
| AC-010 | Piped output byte-identical across themes (no ANSI in pipes) | `02-theme-design.md` |
| AC-011 | Help/UI code references semantic roles only, never literal colors | `02-theme-design.md` |

## Area 3: Help example format (WS2)

| # | Criterion | Source |
|---|-----------|--------|
| AC-012 | `gitmap help example <cmd>` works for every dispatched command and alias | `03-help-example-format.md` |
| AC-013 | Each example shows what-it-does + invocation + expected output in order | `03-help-example-format.md` |
| AC-014 | Examples live in one data registry; helpdoc markdown generated, never hand-edited | `03-help-example-format.md` |
| AC-015 | Every example invocation verified runnable; output shape matches normalized form | `03-help-example-format.md` |
| AC-016 | Examples render clean under `--theme mono` and when piped | `03-help-example-format.md` |

## Area 4: Missing commands (WS3)

| # | Criterion | Source |
|---|-----------|--------|
| AC-017 | Audit enumerates dispatcher verbs vs help entries (A/B/C diffs) using gitmap search only | `04-missing-command-audit.md` |
| AC-018 | `gitmap merge write` creates a GitHub PR with title/body/source/target/draft | `05-merge-write-command.md` |
| AC-019 | `gitmap mr` alias works identically | `05-merge-write-command.md` |
| AC-020 | `--draft`, `--yes`, and failure surfacing behave per §4–§5 | `05-merge-write-command.md` |
| AC-021 | New package ≤300 lines; no pull/commit internals touched | `05-merge-write-command.md` |
| AC-022 | Help example present for `merge write`/`mr` (spec 03 format) | `05-merge-write-command.md` |

## Area 5: Proposals + RCA (WS4, WS5)

| # | Criterion | Source |
|---|-----------|--------|
| AC-023 | Exactly 20 proposals, all template fields filled, in `.ai-memory/plans/` | `06-twenty-new-command-proposals.md` |
| AC-024 | Zero implementation code in WS4 | `06-twenty-new-command-proposals.md` |
| AC-025 | RCA names the 385aba38e triggers (registry Run, VSS, autologin, lock probing) | `07-windows-virus-rca.md` |
| AC-026 | RCA lists prevention: signing, vendor submission, reproducible builds, CI AV gate, avoidance list | `07-windows-virus-rca.md` |
| AC-027 | RCA written for a continuing AI (open questions present) | `07-windows-virus-rca.md` |

## Area 6: Verification + program rules (WS6)

| # | Criterion | Source |
|---|-----------|--------|
| AC-028 | `go build ./...` from `cli/` green after all waves | hard rules |
| AC-029 | Every commit atomic with `gitmap cpf "<module> - <summary>"` and pushed (`git ls-remote` verified) | hard rules |
| AC-030 | No binaries committed; no files deleted; all new files ≤300 lines | hard rules |
| AC-031 | Relative paths only in all artifacts (incl. proposals) | hard rules |
| AC-032 | No `go test` run locally at any point | hard rules |
| AC-033 | No version bump and no release in program 273 | WS6 scope |

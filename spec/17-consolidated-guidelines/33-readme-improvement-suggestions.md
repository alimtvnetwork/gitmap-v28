# README Improvement Suggestions

> **Created:** 2026-04-22 (Asia/Kuala_Lumpur)
> **Owner:** Md. Alim Ul Karim
> **Status:** ✅ Phases A–E shipped in v3.56.0 (2026-04-22)
> **Source:** Audit of `readme.md` (1354 lines, last touched in v3.51.0)

## Purpose

The current root README is comprehensive but has grown to **1354 lines**. New visitors — both humans and AI agents — need 3–5 seconds to understand what the repo does, who it's for, and how to install the piece they want. This file lists concrete, prioritised improvements I can apply on request.

The two new GIFs (`coding-guidelines-walkthrough.gif`, `install-flow.gif`) and the **Bundle Installers** matrix were already added in v3.55.0. Everything below is what comes next.

## Scoring rubric

| Symbol | Meaning |
|---|---|
| 🔴 | Critical — fix before next release |
| 🟡 | High value, low risk |
| 🟢 | Nice-to-have polish |

## Top-of-file (first screen)

| # | Suggestion | Priority |
|---|---|---|
| 1 | Add a row of shields.io badges directly under the H1 — version, license, file count, health score, last commit, "AI-ready". One glance = trust signal. | 🟡 |
| 2 | Replace the long author block at the top with a single-line byline; move the full bio to the existing **Author** section at the bottom. Saves ~6 lines above the fold. | 🟡 |
| 3 | Add a 2-line elevator pitch immediately under the H1 ("What is this? / Who is it for?") before the first image. | 🔴 |
| 4 | Add a "30-second tour" callout box: pick your role → click here. (Linter user / Spec author / WordPress dev / AI agent.) | 🟡 |

## Install section

| # | Suggestion | Priority |
|---|---|---|
| 5 | Promote the new **Bundle Installers** matrix above the generic `install.sh` section — most users want a bundle, not the whole repo. | 🔴 |
| 6 | Add a "Verify your install" subsection with `sha256sum` examples that check against `checksums.txt` from the matching release. | 🟡 |
| 7 | Document the **uninstall** story explicitly. Today the install scripts only add files; users don't know which folders to delete. | 🟡 |
| 8 | Add a Windows Defender / SmartScreen note for the `irm \| iex` pattern (script signing, `Unblock-File`). | 🟢 |

## Structure & navigation

| # | Suggestion | Priority |
|---|---|---|
| 9 | Split readme.md into `readme.md` (≤ 400 lines, "what + how to install") and `docs/architecture.md`, `docs/principles.md`, `docs/author.md`. The 1354-line file violates the project's own 300-line rule. | 🔴 |
| 10 | Auto-generate the Folder Structure tree via `scripts/sync-spec-tree.mjs` and embed it via an HTML comment marker — never hand-edit again. | 🟡 |
| 11 | The two `<details>` blocks (TOC + author assessment) should expand to anchored sections; collapsed `<details>` hide content from many search engines and from Markdown renderers in IDEs. | 🟡 |
| 12 | Add a "What's new in v3.55.0" callout that links to `changelog.md` — refreshed on every release. | 🟢 |

## Visual & accessibility

| # | Suggestion | Priority |
|---|---|---|
| 13 | Add `width="960"` to both new GIF `<img>` tags so they don't blow out the layout on narrow GitHub views. | 🟡 |
| 14 | Provide a **static PNG fallback** of slide 1 of the walkthrough GIF — readers who disable autoplay GIFs (or print the README) still see the title card. | 🟢 |
| 15 | Every image needs descriptive alt text (currently only the spec-viewer image has good alt; the new ones are okay but the existing ones in lower sections are not). | 🟡 |
| 16 | The dark-orange accent in the GIFs (`#FF6E3C`) should also appear in `tailwind.config.ts` as `--brand-accent` so the live spec viewer matches the README art direction. | 🟢 |

## AI-agent ergonomics

| # | Suggestion | Priority |
|---|---|---|
| 17 | Add a `## For AI Agents` section right after the H1 summary listing the canonical entry points: `llm.md`, `bundles.json`, `02-spec/17-consolidated-guidelines/01-index.md`, `.lovable/memory/01-index.md`. | 🔴 |
| 18 | Publish a `bundles.schema.json`-validated table of contents so an AI can answer "which bundle do I need?" with a single fetch instead of crawling the whole repo. | 🟡 |
| 19 | Add a one-line `# coding-guidelines-v24` topic to the GitHub repo metadata so cross-project AI search picks it up. (Repo-side, not file-side.) | 🟢 |

## Honesty & freshness

| # | Suggestion | Priority |
|---|---|---|
| 20 | The header still says "Total Spec Files: 285". The actual current count is 607 (per `src/data/specTree.json`). Either auto-stamp it from a build script or remove it. | 🔴 |
| 21 | "Health Score: 100/100" claims are unverifiable to outside readers — link the dashboard JSON or remove. | 🟡 |
| 22 | The "Last Updated" line should be auto-stamped by the release script, not hand-edited. Currently 6 days stale. | 🟡 |

## Suggested next phase ordering

If you say **"do improvements"** I will execute in this order, one phase per turn:

1. ✅ **Phase A** — Above-the-fold rewrite (badges, elevator pitch, byline, 30-sec tour, auto-stamped counts via `scripts/sync-readme-stats.mjs`). Items #1–#4, #20, #22.
2. ✅ **Phase B** — Bundle matrix promoted above generic install + Verify-your-install + Uninstall sections. Items #5–#8.
3. ✅ **Phase C** — README split: `readme.md` (243 lines) + `docs/architecture.md` + `docs/principles.md` + `docs/author.md`. Items #9, #10.
4. ✅ **Phase D** — `## For AI Agents` canonical-entry-points table (llm.md, bundles.json, version.json, condensed master, anti-hallucination rules, consolidated overview, Lovable memory + prompts). Items #17–#19.
5. ✅ **Phase E** — `<img>` tags with `width="960"` + descriptive alt text on all three images; badges with `--` dash escapes; AI-Ready badge fixed. Items #11–#16, #21.

All five phases shipped together as **v3.56.0** (2026-04-22) per the user's "DO all" instruction. Auto-stamping is wired into `npm run sync` so future spec changes refresh the README without hand-editing.

---

*This file is the single source of truth for README improvement work. Update it as items ship.*

# Confidence Report

> **/goal** Say what a blind agent can build from this design system today.
> **/learn** A row under 70 means stop and read the named file. Do not invent the missing geometry.

**Version:** 4.4.0
**Updated:** 2026-10-06
**Status:** Active

This report replaces the design-spec number ledger. A number that is not in `05-bright-gold-tech/`, `26-visual-builder.md`, `29-slide-navigation-and-builder.md`, `34-slide-layout-catalog.md`, `39-logo-construction.md`, `40-theme-switch.md`, or `42-slide-quiz-preview-chrome-and-default-shadows.md` is not a build license.

---

## 1. Scores

| Asked system | Spec to follow | Score | Can build today |
|---|---|---|---|
| Bright gold tech slides, including the definition page | `05-bright-gold-tech/` | 78 | Yes. Background, tokens, definition page, and the shell for the other nine types. |
| A website band in that same grade | `05-bright-gold-tech/02-page-types.md` section 4 | 70 | Yes, as four bands. Not a full multi-page marketing site. |
| White slide system | `34-slide-layout-catalog.md` | 62 | The measured types in that file. Not `usp-strike` or `bullets`. |
| Logo | `39-logo-construction.md` | 70 | A clean SVG mark after name, idea, and colors are supplied. |
| Noir-gold deck | Tokens below. No page catalog in this repo. | 48 | Colors only. Do not invent its page geometry. |
| Slide theme switch | `40-theme-switch.md` | 76 | The 8 listed themes. A ninth waits for six supplied colors. |
| Quiz preview & presenter shortcuts (closed contract) | `42-slide-quiz-preview-chrome-and-default-shadows.md` v2 + `scripts/verify-botanical-light-contrast.mjs` | **100** | §0.1 scope only: CSS §1–2,3,6, embed §9, `PRESENTER_SHORTCUTS_CORE`, §5.4 dispatch, HUD §3.2.1 in file 31. Faults: app folder `05-blind-agent-fault-register.md`. |
| Quiz preview (full product surfaces) | Same + layouts 34/28 | 74 | Custom slide bodies and deck-specific shortcut rows still out of license. |
| Black slide system | None | 20 | No measured black catalog exists here. Do not invent one. |
| Blog kit | None in this folder | 20 | That kit is a separate light/dark oklch set. Do not restyle it from this file. |

Score meaning: 70 and above, follow the named file and refuse anything it does not specify. Below 70, colors or a warning only.

---

## 2. Why the bright-gold score is 78

The background, the HSL table, the glass card, the accent bar, the pill, the footer slots, and the definition page are written down. An agent that follows those files can produce the definition page and a short deck in that grade.

The score is not 90 because `cover`, `people`, `comparison`, and the rest do not have per-pixel inner layouts. The rule is to use the shell and the slots, not to guess coordinates. A full website with a mega menu in this grade is not specified. Use the white-blue menu files only on a white-blue page.

---

## 3. Noir-gold, do not mix

The noir-gold deck is a different dark gold. Use it only when the request names that deck.

| Token | HSL | Hex twin |
|---|---|---|
| `--background` | `0 0% 5%` | `#0D0D0D` |
| `--foreground` | `42 100% 96%` | `#FFF9EB` |
| `--gold` | `40 88% 50%` | `#F0A50F` |
| `--ember` | `14 80% 57%` | `#E9633A` |
| `--card` | `0 0% 7%` | `#121212` |

Its page-type files live with that deck, not here. Known type names, with no geometry copied: `title`, `tile`, `keyword`, `checklist`, `metric-grid`, `capsule-list`, `section-divider`, `middle-title`, `step-timeline`, `focus-timeline`, `advance-step`, `steps-chain-3d`, `table`, `data-table`, `code-block`, `equation`, `number-callout`, `image`, `qr-meeting`, `layout`, `session-outline`, `er-diagram`, `database-diagram`, `box-diagram`, `bubble-chart`, `blast-radius`. Do not draw those from memory.

---

## 4. White canvas

Follow `34-slide-layout-catalog.md` section 3. Canvas `1920×1080`. Measured types: `title`, `persona`, `key-player`, `before-after`, `pricing`, `steps-chain`, `testimonials`, `talent-funnel`, `white-master`, `competitive-edge`, `tech-stack`, `steps`.

`usp-strike` and `bullets` are not components. Do not build them.

---

## 5. Older design files

Files `24` through `33`, `35`, `37`, and `38` still contain numbers that were not re-measured for this report. When one of those numbers disagrees with file 05, file 26, or file 34, the later file wins. When it disagrees with nothing because it was never measured, leave it unused.

`26-visual-builder.md` remains the only website-builder contract.

---

## 6. Logo

`39-logo-construction.md` is the spec. The agent can draw one archetype, transparent, at `16 / 48 / 1024`. It cannot invent the product. Missing inputs stop the work.

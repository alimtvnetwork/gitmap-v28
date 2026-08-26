# Root README Conventions

**Version:** 1.0.0
**Updated:** 2026-04-22
**Status:** Active — **MANDATORY** for every repository that ships this spec system
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## §1 — Purpose

The root `readme.md` is the **first surface** any human or AI sees. It is
the cover of the book. This file defines the **non-negotiable structure**
every repository in this family MUST follow so that:

- A new contributor recognises the project family within 3 seconds.
- An AI agent can locate canonical entry points without scanning the tree.
- The visual identity (icon, badges, author block) is consistent across
  all sibling repositories (e.g. `coding-guidelines-v24`, `gitmap-v6`).

Any deviation from §3 is a release blocker. Use `linter-scripts/check-root-readme.py`
(roadmap) or manual review against §9's checklist before tagging a release.

---

## §2 — Mandatory Brand Icon

Every root README MUST open with a **centered brand icon**.

### §2.1 — File location

- Path: `public/images/<repo-slug>-icon.png`
- Format: **PNG with transparent background** (so it works on light + dark GitHub themes)
- Size: **512×512** source, displayed at **160×160** in the README.
- Style: gradient geometric mark, no text inside the icon (the H1 below carries the wordmark).

### §2.2 — Markup (verbatim template)

```html
<p align="center">
  <a href="https://github.com/<org>/<repo>">
    <img
      src="public/images/<repo-slug>-icon.png"
      alt="<Project Name> brand icon — <one-line visual description>"
      width="160"
      height="160"
    />
  </a>
</p>
```

Rules:

- `width` and `height` MUST both be set (prevents layout shift).
- `alt` MUST describe the visual, not just repeat the project name.
- The `<a>` wrapper MUST point to the canonical GitHub repo URL.

---

## §3 — Mandatory Hero Block (in order)

Immediately after the icon, the README MUST contain — **in this exact order**:

1. **Centered H1 title** (`<h1 align="center">`) — the wordmark.
2. **Centered tagline** (`<p align="center"><strong>…</strong></p>`) — one or two lines, ≤ 200 chars total. Italicise the language/scope keywords.
3. **Centered primary badges row** — version, file count, folders, lines, license, AI ready, last-updated. Use the `<!-- STAMP:BADGES -->` marker so `scripts/sync-readme-stats.mjs` keeps them current.
4. **Centered platform/quality badges row** — languages, platforms, bundle count, health score, audit score, PR-welcome, etc. Use the `<!-- STAMP:PLATFORM_BADGES -->` marker.
5. **Centered author block** — see §4.
6. **Centered stats line** — auto-stamped counts (`STAMP:FILES`, `STAMP:FOLDERS`, `STAMP:LINES`, `STAMP:VERSION`, `STAMP:UPDATED`).
7. **Horizontal rule** (`---`).

No section between the icon and the first `---` may be left-aligned. Centering
is enforced by wrapping each block in `<p align="center">` (or `<h1 align="center">`).

---

## §4 — Author & Company Block (verbatim template)

The author block MUST be centered and use **this exact template**, with
only the bracketed values substituted:

```html
<p align="center">
  <strong>By <a href="https://alimkarim.com/">Md. Alim Ul Karim</a></strong> — Chief Software Engineer, <a href="https://riseup-asia.com/">Riseup Asia LLC</a><br/>
  <a href="https://www.linkedin.com/in/alimkarim">LinkedIn</a> ·
  <a href="https://stackoverflow.com/users/513511/md-alim-ul-karim">Stack Overflow</a> ·
  <a href="https://github.com/alimtvnetwork">GitHub</a> ·
  <a href="docs/author.md">Full bio</a>
</p>
```

### §4.1 — Author rules

- **Author name** MUST be `Md. Alim Ul Karim` (with the period after `Md`).
- **Title** MUST be `Chief Software Engineer` (no abbreviations).
- **Author homepage** MUST link to `https://alimkarim.com/`.
- The author block MUST appear **once and only once** in the hero. A
  longer bio belongs in `docs/author.md`.

### §4.2 — Company rules

- **Company name** MUST be written as `Riseup Asia LLC` — full legal
  name, no abbreviations like "RAL", no missing `LLC`, no extra
  qualifiers.
- **Company URL** MUST link to `https://riseup-asia.com/`.
- The company name MUST appear **once** in the hero block (in the author
  line). Repeating it in the tagline or badges is forbidden — it reads
  as marketing spam.

### §4.3 — Required link order (left to right)

`LinkedIn` → `Stack Overflow` → `GitHub` → `Full bio`. This order matches
the audience's likely click priority. Do not re-order.

---

## §5 — Badge Strategy

### §5.1 — Two rows, by purpose

| Row | Purpose | Marker | Examples |
|-----|---------|--------|----------|
| Primary | Identity & release state | `STAMP:BADGES` | version, files, folders, lines, license, AI ready, updated |
| Platform | Reach & quality signals | `STAMP:PLATFORM_BADGES` | languages, OS, bundle count, health score, audit score, PRs welcome, stars |

Total badge count: **between 12 and 16**. Fewer feels empty; more is noise.

### §5.2 — Color discipline

Use **one color per semantic meaning** across all repos in the family:

| Color | Hex | Used for |
|-------|-----|----------|
| Blue | `3B82F6` | Version, identity |
| Green | `10B981` / `22C55E` | Counts, health, license |
| Purple | `8B5CF6` / `6366F1` | Folders, platform |
| Amber | `F59E0B` | Lines, stars |
| Cyan | `0EA5E9` | Updated date |
| Orange | `FF6E3C` | AI / brand accent |
| Pink | `EC4899` | Languages |
| Teal | `14B8A6` | Bundles |

### §5.3 — Auto-stamping

Badges MUST be auto-stamped by `scripts/sync-readme-stats.mjs`. Hand-editing
badge values is allowed but will be overwritten on the next `npm run sync`.

---

## §6 — Required Sections (after the hero)

After the hero `---`, the README MUST contain these sections, **in this
exact order**. This sequence is designed for the typical reader journey:
*"What does it look like? → How do I install it? → What is it really? →
How do I use it? → Who made it?"*

| # | Section | Centered? | Purpose |
|---|---------|-----------|---------|
| 1 | Title (H1 wordmark) | ✅ | Already rendered in the hero (§3) — repeated here only conceptually as the anchor |
| 2 | Badges (primary + platform rows) | ✅ | Already rendered in the hero (§3) |
| 3 | 📸 Screenshot or animated walkthrough (GIF/PNG) | ✅ | Visual proof of concept — single centered `<img>`. GIF preferred; static PNG fallback allowed |
| 4 | 🛠️ Install Scripts | ❌ | One-liner installers FIRST (curl/iwr commands), THEN the per-bundle matrix, THEN the flag table. Readers want to *try it* before reading prose |
| 5 | 📖 About this Repo | ✅ (heading + intro) | Centered `<h2 align="center">` heading and a centered intro paragraph, then left-aligned long-form prose telling the **story**: origin, motivation, problem it solves, references to related repos/docs. This is the narrative heart of the README |
| 6 | 🤖 For AI Agents | ❌ | Canonical entry-point table for AI tooling |
| 7 | 📦 Bundle Installers | ❌ | Per-bundle one-liner matrix (deeper than §4's headline installer) |
| 8 | 🚀 Usage / Commands | ❌ | Common commands, scripts, and workflows beyond install |
| 9 | 📚 Documentation | ❌ | Links to `docs/*` deep-dives |
| 10 | 🔄 What's New | ❌ | Last 3 versions; full history in `changelog.md` |
| 11 | 🤝 Contributing | ❌ | How to add specs, modify rules, run health checks |
| 12 | 👤 Author & Company | ✅ | Re-state the §4 author block (or a fuller variant) at the bottom for readers who scrolled past the hero. Centered |
| 13 | Auto-stamp footer | ✅ | Italic note explaining the sync script |

### §6.1 — Section ordering rationale

1. **Visual first** — a screenshot/GIF immediately after the hero proves the
   project is real and working before the reader commits to reading prose.
2. **Install before explanation** — most readers arrive wanting to *try*
   the tool. Putting installers above the long About section respects
   their time. The deeper Bundle Installers matrix (§6 row 7) can repeat
   and expand on the headline installer from row 4.
3. **About in the middle** — once the reader has seen it work and knows
   how to install it, they're ready to invest in the *why* and the
   *story*. The About heading and intro paragraph MUST be centered;
   long-form prose underneath is left-aligned for readability.
4. **Author at the end (and the top)** — the hero has the compact author
   line (§4); the bottom Author section is a fuller, centered re-statement
   so the README closes on attribution.

### §6.2 — Section heading format

Section emoji prefixes (📸 🛠️ 📖 🤖 📦 🚀 📚 🔄 🤝 👤) are **mandatory** —
they aid visual scanning and AI section-detection.

Section headings MAY use either Markdown (`## …`) or centered HTML
(`<h2 align="center">…</h2>`). Sections marked ✅ in the table above
MUST use the centered HTML form for the heading (and centered intro
content where noted). The README linter
(`linter-scripts/check-root-readme.py`) accepts both forms.

---

## §7 — Length & Modularity

- The root README has **no hard line limit** — it can be as long as the
  story and content require.
- That said, prefer extracting deep-dive content into `docs/` to keep the
  root README scannable:
  - `docs/principles.md` — coding principles, CODE RED rules
  - `docs/architecture.md` — folder structure, conventions
  - `docs/author.md` — extended author bio, company background, FAQ
- Cross-reference extracted docs from the §6 "Documentation" section.

---

## §8 — Accessibility & Visual Polish

- Every `<img>` MUST have descriptive `alt` text (not "image" or
  filename).
- Every `<img>` MUST set `width` (and `height` for icons) to prevent
  layout shift.
- GIFs SHOULD have a static PNG fallback for printed/PDF rendering
  (roadmap).
- All hero blocks MUST be centered using `<p align="center">` or
  `<h1 align="center">`. Do not use `<center>` (deprecated).

---

## §9 — Acceptance Checklist (release blocker)

Before tagging a release, verify **all** of the following on the root
`readme.md`:

- [ ] Centered brand icon at `public/images/<repo-slug>-icon.png`, 160×160 display.
- [ ] Centered H1 wordmark immediately below the icon.
- [ ] Centered tagline (≤ 200 chars).
- [ ] Two badge rows: `STAMP:BADGES` (primary) + `STAMP:PLATFORM_BADGES` (platform).
- [ ] Total badge count is 12–16.
- [ ] Centered author block uses the §4 verbatim template.
- [ ] Author name reads exactly `Md. Alim Ul Karim`.
- [ ] Company name reads exactly `Riseup Asia LLC`.
- [ ] Author/company each link to their canonical URL.
- [ ] Centered auto-stamped stats line.
- [ ] Hero block ends with `---` before any left-aligned content.
- [ ] All §6 sections present, in the §6 order, with emoji prefixes.
- [ ] §6 ordering: screenshot/GIF → install scripts → centered About (with story) → AI agents → bundle installers → usage → docs → what's new → contributing → centered author block → auto-stamp footer.
- [ ] Centered sections (screenshot, About heading + intro, bottom Author block, footer) use `<p align="center">` / `<h2 align="center">`.
- [ ] Every `<img>` has `alt`, `width` (and `height` for icons).
- [ ] `npm run sync` was run and committed (badges + stats are current).

---

## §10 — Cross-References

- Sync script: [`scripts/sync-readme-stats.mjs`](../../scripts/sync-readme-stats.mjs)
- Long-form docs: [`docs/principles.md`](../../docs/principles.md), [`docs/architecture.md`](../../docs/architecture.md), [`docs/author.md`](../../docs/author.md)
- Sibling-repo reference implementation: `gitmap-v6` root README
- Folder structure rules: [`01-folder-structure.md`](./01-folder-structure.md)
- Naming conventions: [`02-naming-conventions.md`](./02-naming-conventions.md)
- Required files (per spec module): [`03-required-files.md`](./03-required-files.md)

---

*Root README Conventions — v1.0.0 — 2026-04-22*
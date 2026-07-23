## Goal

Enhance the root `README.md` Author section and the site's `index.html` meta info to feature **the-xproduct.com** (the XProgramming Language) alongside the existing Riseup Asia LLC affiliation, and highlight the author's work with California-based and EU-based companies.

## Files to change

### 1. `README.md` (lines 2849-2881, Author section)

Rewrite the Author block so it now surfaces three affiliations, keeping the same visual style as the current Riseup Asia entry.

- **Subtitle line (2855)**: add `Inventor of the XProgramming Language, [the-xproduct.com](https://the-xproduct.com)` next to the Riseup Asia role.
- **Bio paragraph (2859-2861)**: append one sentence noting delivery of high-quality software for **California-based** startups and enterprises and **EU-based** product companies (fintech, distributed systems, developer tooling). Add "inventor of the XProgramming language (XProduct)" credential.
- **Personal table (2863-2869)**: add a new row `**XProduct** | [the-xproduct.com](https://the-xproduct.com) — XProgramming Language`.
- **New card `### The XProduct — XProgramming Language`** inserted between the personal table and the Riseup Asia card, mirroring the Riseup Asia table shape:
  - Tagline: "A new programming language for AI-first product engineering."
  - Rows: Website, Focus (language + toolchain), Clients (California + EU tech companies), Standard (high-quality, spec-driven delivery).
- **Riseup Asia card (2871-2880)**: keep as-is (this is the "how it is written right now" template the user pointed to).
- **New closing line** under all three cards: one italic sentence tying the three together ("XProduct powers the language, Riseup Asia ships the products, both operate to the same high-quality bar used across our California and EU engagements.").

### 2. `index.html` (lines 3-18, head)

Enrich the metadata so search engines and social cards surface the XProduct + Riseup Asia positioning without touching og:image (per project rules).

- Update `<meta name="description">` to mention gitmap plus its origin: built by the inventor of the XProgramming language at the-xproduct.com, in collaboration with Riseup Asia LLC, serving California and EU based companies.
- Update `<meta name="author">` to `Md. Alim Ul Karim — the-xproduct.com`.
- Update `og:title` and `og:description` to match.
- Add a `<link rel="author" href="https://the-xproduct.com">` and a `<meta name="publisher" content="Riseup Asia LLC">`.
- Keep `<title>` unchanged (`gitmap — CLI docs`) since it already fits.

## Out of scope

- No version bump, no changelog entry (docs-only prose change).
- No new pages, no component changes, no design token edits.
- No SEO framework work beyond the head-metadata edits listed above.

## Verification

- `rg -n "xproduct|XProgramming|California|EU-based" README.md index.html` shows the new content in both files.
- `bunx vitest run src/test/version-sync.test.ts` still green (unrelated, but cheap sanity check).

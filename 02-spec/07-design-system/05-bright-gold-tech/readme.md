# Bright Gold Tech

> **/goal** Build a dark gold technology deck or a website band in that same grade.
> **/learn** Read `01-tokens-background-and-chrome.md`, then `02-page-types.md`. Do not import a second palette onto these pages.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active
**Theme id:** `bright-gold-tech`
**AI Confidence:** See `../98-confidence-report.md`.
**Ambiguity:** Page types other than `definition` share the shell. Their inner coordinates are not specified. Do not invent them.

---

## What this theme is

`bright-gold-tech` is a near-black field, white type, and one vivid gold accent. It is for a technology company keynote, a definition page, or a website section that must look like that keynote.

It is not the noir-gold deck. Noir-gold uses a duller gold and a cream foreground. Do not mix the two token sets on one page.

It is not the white canvas deck. Do not place this background behind a white-canvas layout.

---

## Reading order

1. [01-tokens-background-and-chrome.md](./01-tokens-background-and-chrome.md)
2. [02-page-types.md](./02-page-types.md)
3. [../39-logo-construction.md](../39-logo-construction.md) when the page needs a mark
4. [../98-confidence-report.md](../98-confidence-report.md) before claiming a second deck

---

## Checklist

- [ ] Read file 01, then file 02. Do not open a second palette.
- [ ] Colors are the HSL tokens. Components use `hsl(var(--token))` inside this theme's own classes, or `var(--canvas)` after `40-theme-switch.md` maps them.
- [ ] The page type is one of the ten names in file 02. Any other name stops the work.
- [ ] `definition` uses the type sizes in file 02. Other types use the shell and their slots.
- [ ] Copy comes from the deck. This folder does not supply a company name or a sample sentence.
- [ ] A theme change uses `40-theme-switch.md` and does not edit these layout sizes.

---

## Closed rules

1. Store colors as HSL triplets. Consume them as `hsl(var(--token))`.
2. One accent on a page. Gold is the accent. Ember is not a second brand color on this theme.
3. Text is live HTML. Do not bake a headline into a picture.
4. Copy is supplied by the deck. This folder does not contain a company name, a product name, or a sample sentence to reuse.
5. If a size is not in these two files, do not invent it. Use the shared shell and leave the unknown region empty.

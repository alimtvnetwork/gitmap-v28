# Bright Gold Tech: Page Types

> **/goal** Give each page one job so the deck stays in this grade.
> **/learn** `definition` is fully specified. Every other type uses the same background, card, and footer, and must not invent a second coordinate grid.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active

---

## 1. Closed set

Build only these types. A request for another type stops. Do not alias a missing type onto `definition`.

| Type | Job | Inner geometry |
|---|---|---|
| `definition` | One term, one statement card | Specified in section 2 |
| `cover` | Name and one subtitle | Shell only |
| `divider` | Section name, nothing else | Shell only |
| `comparison` | Two labeled columns | Shell plus two `.glass-card` columns |
| `steps` | Ordered steps, one row | Shell plus one card per step |
| `people` | One person or a short row of people | Shell plus one `.glass-card` per person |
| `quotes` | Quotations | Shell plus `.rec-card` |
| `metrics` | Up to four numbers | Shell plus `.gold-emphasis` figures |
| `stack` | Named tools, no logos baked as text | Shell plus pills |
| `close` | One next action and a contact line | Shell only |

"Shell only" means the background from file 01, the footer from file 01, and one text block. Do not add a chart, a photo collage, or a second accent.

A website may reuse `definition` as a hero band and `metrics` as the next band. It still uses these tokens. It does not switch to the white-blue theme halfway down the page.

---

## 2. `definition`

This is the page in the supplied reference: eyebrow, two-line title, one statement card, footer. Slots are filled by the deck. Do not copy a sample company, a sample URL, or a sample sentence into the component.

Regions, top to bottom, all left aligned:

1. Eyebrow. `{SERIES}` in `--pres-accent`, then a middle dot, then `{INDEX}` and `{SECTION}` in `--pres-text-muted`. Under `{SERIES}` only, one `.accent-bar` at the `72px` width.
2. Title. Line one `{HEADLINE_LEAD}` in `--pres-text`. Line two `{HEADLINE_ACCENT}` in `--pres-accent`. Display font. Do not put both lines in gold.
3. One `.glass-card`. Inside it: a `72px` `.accent-bar`, then the statement, then one `.gold-pill`.
4. Footer from file 01.

The statement is one paragraph. Phrases that carry the claim use `.gold-emphasis`. The rest stays `--pres-text`. One pill. The pill label is a short noun such as the section name. It is not a second headline.

```html
<section class="slide-dark-bg definition-page">
  <p class="eyebrow">
    <span class="gold-emphasis">{SERIES}</span>
    <span class="dot">·</span>
    <span>{INDEX} · {SECTION}</span>
  </p>
  <div class="accent-bar"></div>
  <h1>
    <span class="lead">{HEADLINE_LEAD}</span>
    <span class="gold-emphasis">{HEADLINE_ACCENT}</span>
  </h1>
  <article class="glass-card">
    <div class="accent-bar"></div>
    <p>{STATEMENT_WITH_EMPHASIS}</p>
    <span class="gold-pill">{PILL}</span>
  </article>
  <footer>
    <div class="accent-bar footer-rule"></div>
    <img alt="" height="48" />
    <p class="address">{ADDRESS}</p>
    <p class="index">{CURRENT} / {TOTAL}</p>
  </footer>
</section>
```

```css
.definition-page {
  width: 1920px;
  height: 1080px;
  padding: 72px 80px 48px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.definition-page h1 {
  margin: 28px 0 36px;
  font-size: 72px;
  line-height: 1.05;
  font-weight: 700;
}
.definition-page h1 .lead,
.definition-page h1 .gold-emphasis {
  display: block;
}
.definition-page .glass-card {
  width: min(1080px, 100%);
  padding: 36px 40px 32px;
}
.definition-page .glass-card p {
  margin: 28px 0;
  font-size: 32px;
  line-height: 1.35;
  font-weight: 500;
}
.definition-page .eyebrow {
  font-size: 14px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}
.definition-page .footer-rule {
  width: 100%;
}
```

The `72px` title, `32px` statement, `1080px` card width, and `72px 80px 48px` padding are the authoring contract for `definition` only. They were set in this file so every agent uses one scale. They were not copied from a second component. Do not resize `cover` or `people` to this type scale, and do not replace these numbers with a different guess.

Do not add a right-hand illustration, a second card, or a filled gold panel behind the title.

---

## 3. The other types

Each type still mounts `.slide-dark-bg` and the footer.

| Type | Required slots | Refuse |
|---|---|---|
| `cover` | `{SERIES}`, `{TITLE}`, `{SUBTITLE}` | A statement card |
| `divider` | `{SECTION}` only, centered | Body copy |
| `comparison` | `{TITLE}`, two cards, each with `{LABEL}` and up to four rows | A third column |
| `steps` | `{TITLE}`, three to five cards, each `{INDEX}` plus `{LABEL}` | A vertical novel |
| `people` | `{NAME}`, `{ROLE}`, one portrait per person | A portrait wider than the card |
| `quotes` | `{QUOTE}`, `{NAME}`, `{ROLE}` inside `.rec-card` | A star rating row |
| `metrics` | Up to four pairs of `{VALUE}` and `{LABEL}` | A chart library |
| `stack` | `{TITLE}` and `.gold-pill` labels | Invented product logos |
| `close` | `{TITLE}`, one `.gold-pill`, one contact line | A second call to action |

Portrait files are photographs. Names and quotes are deck copy. Do not generate a face.

---

## 4. Website band

A marketing page that should look like this deck stacks bands in this order only:

1. `definition` or `cover`
2. `metrics` or `steps`
3. `comparison` or `quotes`
4. `close`

The page background between bands stays `--pres-bg`. Do not insert a white band. Navigation, when present, is a `72px` bar on `--pres-bg` with the mark at `48px` and links in `--pres-text-muted`. The current link is `--pres-accent`.

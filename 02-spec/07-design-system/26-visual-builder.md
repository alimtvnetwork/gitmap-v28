# 26 — Visual Builder Overlay

> **/goal** Specify a review-only overlay that edits wording, images, icons, menu labels, and in-group order on the live page.
> **/learn** Gate, modes, sanitizer, storage bucket, and export. This file is not the slide builder. Slides use `29-slide-navigation-and-builder.md`.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active

---

## 0. Anti-hallucination

If a limit is not written here, do not invent it. The public site does not change until a developer applies the export. The overlay never publishes by itself.

Placeholders in examples are `{SITE_DOMAIN}` and `{OWNER_EMAIL}`. Do not replace them with a real organization name.

---

## 1. Gate

The editor opens only when both query parameters are present and the email matches the configured owner:

```text
{SITE_DOMAIN}/?builder=1&email={OWNER_EMAIL}
```

Wrong or missing parameters: ship no builder script and no builder UI. The query stays on the URL through in-app navigation so the next page stays editable. Strip the query from storage keys and from exported URLs.

---

## 2. Modes

One mode is active at a time.

| Mode | Click target | Edit |
|---|---|---|
| Text | Headings, paragraphs, list items, button and link labels | In place. Default. |
| Images | Content `img`, and small inline SVG icons | Replace media or icon |
| Menu | Header and footer links | Label, description, href |
| Layout | Repeated card or list groups | Move up or down inside the group |
| Off | Nothing | Chrome hides. Edits stay visible |

Structural chrome (chevrons, menu icon, close icon) is not editable. Builder UI nodes are not editable.

A plain click on a link edits the label. Navigation uses an Open control on the hovered link, or Cmd, Ctrl, or Alt click. The two query parameters travel with that navigation.

---

## 3. Text

Double-click or click in Text mode sets `contentEditable` to `plaintext-only` where supported, otherwise `true`. Place the caret with `caretPositionFromPoint` or `caretRangeFromPoint`.

- `Escape` restores the previous value.
- `Ctrl+Enter` or `Cmd+Enter` commits. Plain `Enter` inserts a line break and does not commit.
- The toolbar chip states the commit keys.
- Empty text is rejected.
- If the sanitized value equals the original, delete the saved record.

Allowed tags, exactly: `STRONG`, `B`, `EM`, `I`, `BR`, `A`, `SPAN`.

- Keep `SPAN` only when `class="gradient-text"`. Unwrap every other span.
- Strip every attribute except `href` on `A`, and `class` plus `data-text` on the accent span.
- Drop `href` values that start with `javascript:` or `data:`.
- No `style`. No `on*` handlers.
- Any other element is replaced by its children.

`slugify` is lowercase, NFKD, non-alphanumerics collapsed to `-`, trimmed, truncated to 80 characters. Do not change 80 once ids are stored.

While `[data-bm-editing="true"]`, flatten masked or gradient headings to solid text so the caret is visible.

A touch counts as a tap only when the finger moves less than `10px` and lifts within `600ms`.

---

## 4. Images and icons

Accepted types, exactly: `image/png`, `image/jpeg`, `image/webp`, `image/svg+xml`. Reject AVIF. Hard-reject any file over `2 MB`. Do not warn and continue.

Alt text is required. Suggested file name is `{slug}-{short-hash}.{ext}`. Applying a change sets `src` and `alt` and clears `srcset` and `sizes`.

Icons stay inside Images mode. The stored type is still `icon`.

- Pick by name from the project icon catalogue, or bundle that catalogue. Do not fetch an unpinned remote icon index.
- Custom SVG max `200 KB`. Parse with `DOMParser` in the SVG namespace. Keep `viewBox`, or fall back to `0 0 24 24`. Recreate children with `createElementNS`.
- Reset deletes the record. It does not store a "back to default" row.

---

## 5. Menu

Opt-in attributes:

| Attribute | Node | Meaning |
|---|---|---|
| `data-bm-menu-label` | Text inside the link | Label |
| `data-bm-menu-description` | Secondary line | Description |
| `data-bm-description` | The link | Fallback description |

Href must match:

```text
/^(\/|https?:\/\/|mailto:|tel:|#)/
```

Menu changes live in one site-wide store, not under a page. v1 does not add, remove, or reorder menu items.

---

## 6. Layout

Layout mode shows up and down controls on each item of a tagged group.

- Reorder the live DOM nodes.
- Store the new label sequence against the original order captured at tag time.
- If the sequence returns to the original order, delete the record.
- Reposition controls on scroll and resize inside `requestAnimationFrame`.
- After each move, recompute index, inset, and `z-index` on stacked slots, and renumber visible `01` / `02` / `03` labels in the same frame.
- Replay the saved order on the next load before those effects run.

Items cannot move between groups, and cannot be added or deleted.

---

## 7. Panel, history, export

Desktop: right rail. Below `768px`: bottom dock, `max-height: 62vh`, collapsed at start.

Panel shows the mode switch, save status (*Saving…* / *Saved* / *Couldn't save*), change counts, editing on/off, show original, undo, reset this page, history, and export.

A normal edit saves after `250ms`. Deleting a record saves after `0ms`. A conflict retries after `40ms`. Backfill of breadcrumb and context waits `1500ms`. Do not use a different delay unless this file is revised.

Each element keeps `before` immutable and a `history` list. Undo pops the latest timestamp across buckets and reloads. Reset page deletes that page bucket after confirm.

Identity attributes: `data-builder-id`, `data-builder-type`, `data-bm-section-name`, `data-bm-order`, `data-bm-source`. Wait until the DOM is quiet for `700ms`, and never wait longer than `4s`, before the first scan. Rescan debounce is `120ms`.

Export is a zip:

```text
content-changes--all-pages--YYYY-MM-DD-HHmm.zip
  SUMMARY.md
  content-changes--{page-slug}--YYYY-MM-DD-HHmm.md
  menu-changes.md
  images/
  icons/
```

Include the time in the file name so two exports on the same day do not overwrite. Sort changes by `data-bm-order`. Each change shows before and after. Do not put colors, fonts, or new sections in the export. The overlay cannot edit those.

---

## 8. Identity

Do not use a positional counter (`text-7`). Do not use `data-edit-id` or `data-source-file`.

1. If the rendered text matches a string in the content modules, the id is that source key (`home.hero.titleLines[1]`). The stored source is the content file and that key.
2. Otherwise hash the seed `type|section|tagName|text` with DJB2 starting at `5381`, unsigned, base-36. Freeze the algorithm once a draft exists.
3. Reserve ids already on the page, then de-duplicate with a `-1`, `-2` suffix.
4. Store the origin text in `data-bm-origin` at tag time. A later scan re-binds a lost id from that text.
5. Nav ids are seeded from the label. Renaming a nav label starts a fresh history for that link.

`{OWNER_EMAIL}` is compared on the server. It is not shipped in the client bundle.

---

## 9. Session and shared draft

Activation calls `openBuilderSession({ email })` and stores the returned token in `sessionStorage`. A cached token wins over the email in the URL. Tokens shorter than `24` characters are rejected. Lookup hashes the token in application code and compares with a constant-time compare. It is not a SQL filter on the hash.

The store version is `2`. Page buckets use the path with `builder`, `email`, and `draft` removed. Menu changes are one flat bucket.

The default store is one server-only draft table. No `anon` or `authenticated` grant. Deny-all row security. The only doors are `loadBuilderDraft` and `saveBuilderDraft`. A save sends `revision` and writes only when it still matches.

On conflict, merge shallow and let the local copy win, per page for `pages` and once for `menu`. On network failure, keep a recovery copy under `{BRAND}.builder.shared-recovery.v2.{first 12 chars of the token}`.

A browser-only `localStorage` key `{BRAND}.builder.v2` is the single swap for the two server calls. It is not a second product, and it does not change ids, modes, or the export.

Screenshots are session-only. They are not written into the draft.

---

## 10. Screenshots

Photograph the section (`section`, else `header`, else `footer`, else `main`), not the element alone. After: `3px` red outline, `3px` offset. Before: the same treatment in grey.

Before capture: scroll the section to center, wait one frame, `180ms`, then another frame, decode images, and wait for `document.fonts.ready`. Add `.bm-shooting`, which sets `animation`, `transition`, `transform`, and `filter` to none, and `opacity` and `visibility` to visible, including `[data-reveal]`.

Render with `html-to-image` `toPng`, `pixelRatio: 1`, white background, cap `1440 × 2400`. Exclude `[data-bm-root]`. Blank check: scale to `48px` wide; if no channel differs from the first pixel by more than `6`, discard the frame. Restore styles in `finally`.

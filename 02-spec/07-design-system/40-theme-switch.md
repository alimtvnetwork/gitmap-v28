# 40 Theme Switch

> **/goal** Change a deck's colors without changing its layout, and know exactly how many themes exist.
> **/learn** There are 8 slide themes. Each one sets the same six variables. A ninth theme needs six colors from the user before it is added.

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** An agent can switch these 8. It cannot invent a ninth.
**Ambiguity:** File 27 also mentions a three-token bright-gold row. That row is not a switch id. The switch id for the definition-page gold is `bright-gold-tech`.

---

## 1. Count

Slide switcher count: **8**. Website themes are not in this count. They stay on their own pages.

| # | Id | Kind | Where the colors are defined |
|---|---|---|---|
| 1 | `bright-gold-tech` | Dark gold deck | `05-bright-gold-tech/01-tokens-background-and-chrome.md` |
| 2 | `noir-gold` | Dark muted gold | `27-slide-canvas-and-themes.md` section 4 |
| 3 | `snow` | Black and white | `27-slide-canvas-and-themes.md` section 3 |
| 4 | `midnight` | Near-black, gold mark | `27-slide-canvas-and-themes.md` section 3 |
| 5 | `paper` | Warm light | `27-slide-canvas-and-themes.md` section 3 |
| 6 | `sunset` | Dark rose | `27-slide-canvas-and-themes.md` section 3 |
| 7 | `print` | White print | `27-slide-canvas-and-themes.md` section 3 |
| 8 | `playbook` | Warm light, orange mark | `27-slide-canvas-and-themes.md` section 3 |

The nine families in the design-system readme are reference palettes. They are not this switcher. Do not add them here.

`white-blue` is a website theme. It is not one of the 8. Do not set it on a slide stage.

---

## 2. Shared variables

Every theme block sets these six, and only these six, on the stage:

| Variable | Role |
|---|---|
| `--canvas` | Page ground |
| `--ink` | Primary type |
| `--ink-muted` | Captions and indexes |
| `--accent` | The one accent |
| `--accent-ink` | Type that sits on the accent |
| `--card` | Card fill |

Components read `var(--canvas)`, `var(--ink)`, `var(--ink-muted)`, `var(--accent)`, `var(--accent-ink)`, and `var(--card)`. They do not wrap those variables in `hsl()` again. They do not read a hex from a component file.

When a source row has no card color, set `--card` to the same value as `--canvas`. That is the rule. It is not a new color.

---

## 3. Blocks

```css
.slide-stage {
  background: var(--canvas);
  color: var(--ink);
  transition: background-color 0.45s cubic-bezier(0.4, 0, 0.2, 1), color 0.45s cubic-bezier(0.4, 0, 0.2, 1);
}
[data-theme="bright-gold-tech"] {
  --canvas: hsl(240 20% 4%);
  --ink: hsl(0 0% 100%);
  --ink-muted: hsl(218 11% 75%);
  --accent: hsl(41 100% 50%);
  --accent-ink: hsl(240 20% 4%);
  --card: hsl(240 15% 11%);
}
[data-theme="noir-gold"] {
  --canvas: hsl(0 0% 5%);
  --ink: hsl(0 0% 100%);
  --ink-muted: hsl(0 0% 65%);
  --accent: hsl(45 56% 54%);
  --accent-ink: hsl(0 0% 8%);
  --card: hsl(0 0% 5%);
}
[data-theme="snow"] {
  --canvas: #000000;
  --ink: #ffffff;
  --ink-muted: #b8b8b8;
  --accent: #ffffff;
  --accent-ink: #000000;
  --card: #000000;
}
[data-theme="midnight"] {
  --canvas: #101010;
  --ink: #ffffff;
  --ink-muted: #b8b8b8;
  --accent: #ffd83a;
  --accent-ink: #1a1100;
  --card: #101010;
}
[data-theme="paper"] {
  --canvas: #f5f0e6;
  --ink: #1a1a1a;
  --ink-muted: #615a4f;
  --accent: #1d4ed8;
  --accent-ink: #f5f0e6;
  --card: #f5f0e6;
}
[data-theme="sunset"] {
  --canvas: #1b0d1f;
  --ink: #ffeaf0;
  --ink-muted: #c89aa6;
  --accent: #ff7a59;
  --accent-ink: #1b0d1f;
  --card: #1b0d1f;
}
[data-theme="print"] {
  --canvas: #ffffff;
  --ink: #000000;
  --ink-muted: #444444;
  --accent: #000000;
  --accent-ink: #ffffff;
  --card: #ffffff;
}
[data-theme="playbook"] {
  --canvas: #faf7f3;
  --ink: #141414;
  --ink-muted: #6b6b6b;
  --accent: #e8701a;
  --accent-ink: #1a1a1a;
  --card: #faf7f3;
}
@media (prefers-reduced-motion: reduce) {
  .slide-stage { transition-duration: 0.01ms; }
}
```

`bright-gold-tech` also keeps the dot-field background from its own file. The other seven themes use a flat `--canvas`. Do not copy the gold dot field onto `paper` or `print`.

Switching themes does not change the page type, the type size, or the logo geometry.

---

## 4. How a control switches

1. Render one control that lists the 8 ids, in the table order.
2. The active id has a `1px` border of `var(--accent)` and a label in `var(--ink)`.
3. Choosing an id sets `data-theme` on `.slide-stage` only.
4. Store the id on the deck. Do not store a second copy of the six colors.
5. Default, when the deck has no id, is `bright-gold-tech`.

---

## 5. Adding a theme

Stop unless the user supplied all six values. Then:

1. Add the next number and the new id to the table in section 1.
2. Add one `[data-theme="..."]` block with those six values.
3. Change the count in section 1 and in the checklist.
4. Leave every existing block unchanged.

Do not derive a palette from a mood word. Do not copy `white-blue` onto a slide.

---

## 6. Checklist

- [ ] Count the `[data-theme]` blocks. The count is 8, unless section 5 was completed in this same edit.
- [ ] Each block sets `--canvas`, `--ink`, `--ink-muted`, `--accent`, `--accent-ink`, and `--card`.
- [ ] No component contains a hex or an `hsl()` of its own.
- [ ] The stage transition is `0.45s` with `cubic-bezier(0.4, 0, 0.2, 1)`.
- [ ] Reduced motion uses `0.01ms`.
- [ ] The gold dot field is only on `bright-gold-tech`.
- [ ] Layout files were not edited during the switch.
- [ ] `white-blue` is not in the slide control.
- [ ] A requested id that is not in the table is refused.

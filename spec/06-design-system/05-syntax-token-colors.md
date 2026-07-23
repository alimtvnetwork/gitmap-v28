# 05 — Syntax Token Colors

> These are the Highlight.js token colors used inside code blocks.
> They are defined in `src/index.css` as global CSS rules.
> Code block background uses `--terminal` from `01-colors-and-themes.md`.

---

## Base

```css
.hljs { color: hsl(220, 20%, 92%); }
```

---

## Token Map

| Token Class | Color HSL | Visual |
|-------------|-----------|--------|
| `.hljs-keyword`, `.hljs-type`, `.hljs-built_in`, `.hljs-selector-tag` | `207, 82%, 66%` | Blue |
| `.hljs-string`, `.hljs-attr`, `.hljs-property` | `95, 38%, 62%` | Green |
| `.hljs-number`, `.hljs-variable`, `.hljs-regexp`, `.hljs-template-variable` | `29, 54%, 61%` | Orange |
| `.hljs-comment`, `.hljs-quote` | `220, 10%, 45%` | Gray (italic) |
| `.hljs-title`, `.hljs-section`, `.hljs-tag` | `220, 20%, 85%` | Light gray |
| `.hljs-title.function_`, `.hljs-title.class_` | `39, 67%, 69%` | Gold |
| `.hljs-literal`, `.hljs-symbol` | `207, 82%, 66%` | Blue |
| `.hljs-meta`, `.hljs-meta .hljs-keyword` | `355, 65%, 65%` | Red |
| `.hljs-params` | `220, 20%, 92%` | Default |
| `.hljs-attribute` | `39, 67%, 69%` | Gold |
| `.hljs-selector-class`, `.hljs-selector-id` | `95, 38%, 62%` | Green |
| `.hljs-addition` | `135, 52%, 60%` | Green + 10% green bg |
| `.hljs-deletion` | `355, 65%, 65%` | Red + 10% red bg |

---

## Design Rationale

This palette is inspired by **GitHub Dark** with adjustments for the
terminal-green brand aesthetic:

- **Keywords** are blue to stand out from the green brand color.
- **Strings** are muted green to complement without clashing with `--primary`.
- **Comments** are low-contrast gray italic — de-emphasized.
- **Functions/classes** are gold for quick identification.
- **Additions/deletions** use tinted backgrounds for diff visibility.

---

## Customization

To change the syntax theme:

1. Edit the `.hljs-*` rules in `src/index.css` (lines 230–290).
2. Keep contrast ratio ≥ 4.5:1 against `--terminal` background.
3. Test with Go, TypeScript, Bash, and JSON samples.

---

## Constraints

- All token colors are HSL values — no hex or named colors.
- Selection inside code blocks is always white text on green bg (see `03-code-blocks.md`).
- Token colors do not change between light/dark mode — code blocks always use dark terminal bg.

# Code Line Alignment Mismatch

**Created:** 2026-04-03  
**Status:** Active  
**Severity:** High  
**Area:** Docs Viewer / MarkdownRenderer / Syntax Highlighting

---

## Summary

In rendered code blocks, the line numbers in the left gutter do not align vertically with the actual code rows on the right.

This is visible in structure/tree examples and standard fenced code blocks inside the docs viewer.

---

## Symptoms

- First code row starts lower than line number `1`
- Every subsequent row appears shifted relative to the matching line number
- The mismatch is most obvious in short tree/structure examples
- Font-size controls do not fully solve the offset because the root cause is extra inner padding, not only text sizing

---

## Root Cause

The imported `highlight.js` theme (`github-dark.css`) applies its own default spacing:

```css
pre code.hljs {
  display: block;
  overflow-x: auto;
  padding: 1em;
}
code.hljs {
  padding: 3px 5px;
}
```

Our custom renderer already provides spacing at the container level:

- `.code-line-numbers { padding: 1rem 0; }`
- `.code-content { padding: 1rem 1.25rem; }`

Because `code.hljs` still keeps its internal padding, the code text gets pushed downward and inward while the line number gutter does not. That creates a persistent vertical mismatch between line numbers and code rows.

---

## Resolution

Override highlight.js block-level spacing inside the docs viewer so the rendered `<code class="hljs">` uses:

- `padding: 0`
- `margin: 0`
- `line-height: inherit`

This ensures that only the outer layout containers control spacing and both gutters stay synchronized.

---

## Validation

After the fix:

- Line `1` aligns with the first visible code row
- Tree examples align correctly
- Alignment remains correct after font-size changes
- Hover/pinned/range selection still lines up with the correct rows

---

## Affected Files

- `src/components/MarkdownRenderer.tsx`
- `src/index.css`
- `node_modules/highlight.js/styles/github-dark.css` (source of conflicting default theme padding)

---

*Issue documented: 2026-04-03*

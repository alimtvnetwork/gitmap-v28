# Checklist Blocks — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Consecutive lines matching the checklist pattern (`- [ ]` or `- [x]`) are grouped into a styled checklist block with a header, copy button, and export button.

---

## Detection Pattern

```regex
/^(\s*)([-*+]|\d+\.)\s+\[([ xX])\]\s+(.+)$/
```

Matches:
- Optional leading whitespace (for nesting)
- List marker (`-`, `*`, `+`, or `1.`)
- Checkbox `[ ]` or `[x]`/`[X]`
- Content text

---

## HTML Structure

```html
<div class="checklist-block" data-checklist-id="{ID}">
  <div class="checklist-header">
    <span class="checklist-title">Checklist</span>
    <div class="checklist-actions">
      <button class="checklist-copy-btn" data-checklist="{ENCODED_MD}">
        <svg><!-- copy icon --></svg>
        <span>Copy</span>
      </button>
      <button class="checklist-export-btn" data-checklist="{ENCODED_MD}">
        <svg><!-- share/export icon --></svg>
        <span>Export</span>
      </button>
    </div>
  </div>
  <ul class="checklist-items">
    <li class="spec-checkbox checked" style="margin-left: {LEVEL}rem">
      <span class="checkbox-box checked-box">✓</span>
      <span class="checklist-item-content">{CONTENT}</span>
    </li>
    <li class="spec-checkbox" style="margin-left: 0rem">
      <span class="checkbox-box empty-box"></span>
      <span class="checklist-item-content">{CONTENT}</span>
    </li>
  </ul>
</div>
```

---

## Nesting

Indent level is calculated as `Math.floor(indent / 2)` where indent is the number of spaces (tabs = 2 spaces). Each level adds `1rem` left margin.

---

## Copy Button

- Copies the **raw markdown** (not HTML) of the checklist
- Markdown is HTML-encoded in `data-checklist` attribute
- Uses same encoding as code blocks: `"` → `&quot;`, `'` → `&#39;`, `\n` → `&#10;`
- Shows "Copied!" feedback for 2000ms

---

## Export Button

- Downloads the checklist as `checklist.md`
- Creates `Blob` with `type: "text/markdown;charset=utf-8"`
- Same temporary `<a>` download pattern as code blocks
- Shows "Exported!" feedback for 2000ms

---

## Inline Formatting in Items

Checklist item content supports:
- Inline code: `` `code` `` → `<code class="inline-code">`
- Bold: `**text**` → `<strong>`
- Italic: `*text*` → `<em>`
- Links: `[text](url)` → `<a class="spec-link">`

---

## Styling

### `.checklist-block`

```css
border: 1px solid hsl(var(--border));
border-radius: 0.65rem;
margin: 0.75rem 0;
overflow: hidden;
background: hsl(var(--card));
transition: box-shadow 0.3s ease, border-color 0.2s ease;
```

**Hover:**
```css
border-color: hsl(var(--primary) / 0.3);
box-shadow: 0 4px 16px hsl(var(--primary) / 0.06);
```

### `.checklist-header`

```css
display: flex;
align-items: center;
justify-content: space-between;
padding: 0.4rem 0.85rem;
background: hsl(var(--muted) / 0.4);
border-bottom: 1px solid hsl(var(--border) / 0.6);
```

### `.checklist-title`

```css
font-size: 0.7rem;
font-weight: 600;
text-transform: uppercase;
letter-spacing: 0.05em;
color: hsl(var(--muted-foreground));
/* Pseudo-element ::before adds "☐" icon */
```

### `.checklist-copy-btn` / `.checklist-export-btn`

```css
display: flex;
align-items: center;
gap: 0.3rem;
font-size: 0.65rem;
padding: 0.2rem 0.5rem;
border-radius: 0.35rem;
background: hsl(var(--muted) / 0.5);
color: hsl(var(--muted-foreground));
border: 1px solid hsl(var(--border));
cursor: pointer;
transition: all 0.2s ease;
```

**Hover:**
```css
background: hsl(var(--primary) / 0.1);
color: hsl(var(--primary));
border-color: hsl(var(--primary) / 0.3);
```

### `.spec-checkbox`

```css
display: flex;
align-items: center;
gap: 0.5rem;
padding: 0.3rem 0.4rem;
list-style: none;
border-radius: 0.35rem;
transition: transform 0.15s ease, background 0.2s ease;
```

**Hover:**
```css
transform: translateX(3px);
background: hsl(var(--primary) / 0.04);
```

### `.checkbox-box`

```css
display: inline-flex;
align-items: center;
justify-content: center;
width: 1.15rem;
height: 1.15rem;
border-radius: 4px;
font-size: 0.7rem;
font-weight: 700;
flex-shrink: 0;
transition: transform 0.2s ease, box-shadow 0.2s ease;
```

- `.checked-box`: green gradient background
- `.empty-box`: bordered empty square

---

*Checklist Blocks — updated: 2026-04-08*

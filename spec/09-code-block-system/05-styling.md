# Styling — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Design Philosophy

Code blocks use a **permanently dark theme** regardless of the app's light/dark mode. The background is always `hsl(220, 14%, 11%)`. This creates a consistent reading experience and avoids jarring contrast shifts.

Each language gets a unique **accent color** (HSL) that tints the badge, hover shadow, and glow effects — giving each block a subtle personality.

---

## CSS Custom Properties

### On `.code-block-wrapper`

| Property | Default | Modified By |
|----------|---------|-------------|
| `--code-font-size` | `18px` | Font control buttons |
| `--code-line-height` | `1.6` | Fixed |
| `--lang-accent` | Per-language HSL | Set at build time from `LANG_COLORS` |
| `--badge-color` | Same as `--lang-accent` | Set on `.code-lang-badge` |

---

## Complete Class Reference

### `.code-block-wrapper`

```css
border-radius: 0.75rem;
overflow: hidden;
border: 1px solid hsl(220, 13%, 22%);
background: hsl(220, 14%, 11%);
font-family: 'Ubuntu Mono', 'JetBrains Mono', ui-monospace, monospace;
position: relative;
transition: box-shadow 0.3s ease, transform 0.2s ease;
```

**Hover:**
```css
box-shadow: 0 8px 32px hsl(var(--lang-accent) / 0.1),
            0 0 0 1px hsl(var(--lang-accent) / 0.15);
transform: translateY(-2px);
```

### `.code-block-header`

```css
display: flex;
align-items: center;
justify-content: space-between;
padding: 0.5rem 1rem;
background: hsl(220, 14%, 14%);
border-bottom: 1px solid hsl(220, 13%, 20%);
flex-shrink: 0;
```

### `.code-lang-badge`

```css
display: flex;
align-items: center;
gap: 0.4rem;
font-size: 0.7rem;
font-weight: 600;
color: hsl(var(--badge-color));
text-transform: uppercase;
letter-spacing: 0.05em;
```

### `.code-lang-dot`

```css
width: 7px;
height: 7px;
border-radius: 50%;
background: hsl(var(--badge-color));
box-shadow: 0 0 6px hsl(var(--badge-color) / 0.5);
```

### `.code-line-count`

```css
font-size: 0.65rem;
color: hsl(220, 10%, 45%);
font-weight: 400;
margin-right: 0.5rem;
```

### `.code-selection-label`

```css
font-size: 0.6rem;
font-weight: 600;
color: hsl(var(--primary));
background: hsl(var(--primary) / 0.1);
padding: 0.15rem 0.45rem;
border-radius: 0.3rem;
letter-spacing: 0.02em;
animation: slideUpBar 0.2s ease;
```

### `.code-tool-btn`

```css
display: flex;
align-items: center;
gap: 0.3rem;
font-size: 0.65rem;
padding: 0.25rem 0.5rem;
border-radius: 0.375rem;
background: hsl(220, 13%, 20%);
color: hsl(220, 10%, 65%);
border: 1px solid hsl(220, 13%, 25%);
cursor: pointer;
font-family: 'Ubuntu Mono', ui-monospace, monospace;
white-space: nowrap;
transition: all 0.2s ease;
```

**Hover:**
```css
background: hsl(220, 13%, 28%);
color: hsl(0, 0%, 95%);
border-color: hsl(var(--lang-accent) / 0.4);
box-shadow: 0 0 8px hsl(var(--lang-accent) / 0.15);
```

**Copied state (`.copied`):**
```css
background: hsl(152, 60%, 18%);
color: hsl(152, 70%, 60%);
border-color: hsl(152, 50%, 30%);
```

### `.code-font-controls`

```css
display: flex;
gap: 1px;
background: hsl(220, 13%, 18%);
border-radius: 0.375rem;
overflow: hidden;
border: 1px solid hsl(220, 13%, 25%);
```

Inner buttons:
```css
border: none;
border-radius: 0;
font-weight: 700;
font-size: 0.6rem;
padding: 0.25rem 0.4rem;
min-width: 1.75rem;
justify-content: center;
```

### `.code-block-body`

```css
display: flex;
overflow-x: auto;
```

### `.code-line-numbers`

```css
padding: 1rem 0;
text-align: right;
user-select: none;
min-width: 3rem;
background: hsl(220, 14%, 9%);
border-right: 1px solid hsl(220, 13%, 18%);
margin: 0;
overflow: hidden;
flex-shrink: 0;
font-size: var(--code-font-size);
line-height: var(--code-line-height);
```

### `.code-line-number`

```css
display: flex;
align-items: center;
justify-content: flex-end;
height: calc(var(--code-font-size) * var(--code-line-height));
padding: 0 0.75rem;
font-size: calc(var(--code-font-size) * 0.7);
line-height: 1;
color: hsl(220, 10%, 35%);
cursor: pointer;
user-select: none;
transition: color 0.15s ease;
```

### `.code-content`

```css
padding: 1rem 1.25rem 1rem 0.5rem;
margin: 0;
flex: 1;
overflow-x: auto;
background: transparent;
font-size: var(--code-font-size);
```

### `.code-line`

```css
display: block;
white-space: pre;
line-height: var(--code-line-height);
padding: 0 0.25rem;
border-radius: 2px;
transition: background 0.15s ease;
```

**Hover / `.line-highlight`:**
```css
background: hsl(220 15% 16%);
```

**`.line-pinned`:**
```css
background: hsl(var(--primary) / 0.12);
```

**`.code-line-number.line-highlight`:**
```css
color: hsl(var(--primary));
background: hsl(220 15% 12%);
```

**`.code-line-number.line-pinned`:**
```css
color: hsl(var(--primary));
border-right: 2px solid hsl(var(--primary) / 0.6);
background: hsl(var(--primary) / 0.12);
```

---

## Selection Bar

### `.copy-selected-bar`

```css
display: flex;
align-items: center;
gap: 0.5rem;
padding: 0.35rem 1rem;
background: hsl(var(--primary) / 0.08);
border-top: 1px solid hsl(var(--primary) / 0.2);
animation: slideUpBar 0.2s ease;
```

### `.copy-selected-label`

```css
font-size: 0.7rem;
color: hsl(var(--primary));
font-weight: 600;
letter-spacing: 0.02em;
```

### `.copy-selected-btn`

```css
background: hsl(var(--primary) / 0.15);
border-color: hsl(var(--primary) / 0.3);
color: hsl(var(--primary));
```

### `.clear-selected-btn`

```css
margin-left: auto;
background: transparent;
border-color: transparent;
color: hsl(220, 10%, 50%);
font-size: 0.75rem;
padding: 0.2rem 0.4rem;
```

**Hover:** `color: hsl(var(--destructive)); background: hsl(var(--destructive) / 0.1);`

---

## Fullscreen

### `.code-fullscreen-overlay`

```css
position: fixed;
inset: 0;
background: hsl(0 0% 0% / 0.7);
backdrop-filter: blur(4px);
z-index: 998;
```

### `.code-block-wrapper.code-fullscreen`

```css
position: fixed !important;
inset: 2rem;
z-index: 999;
border-radius: 1rem;
transform: none !important;
display: flex;
flex-direction: column;
max-height: calc(100vh - 4rem);
box-shadow: 0 25px 80px hsl(var(--lang-accent) / 0.25),
            0 0 0 1px hsl(var(--lang-accent) / 0.3);
```

---

## Inline Code

### `.inline-code`

```css
background: hsl(var(--code-bg));
color: hsl(var(--code-text));
padding: 0.2em 0.45em;
border-radius: 5px;
font-size: 0.85em;
font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, monospace;
font-weight: 500;
border: 1px solid hsl(var(--border) / 0.5);
white-space: nowrap;
transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.15s ease;
```

**Hover:**
```css
box-shadow: 0 0 0 2px hsl(var(--highlight-glow) / 0.15);
transform: translateY(-1px);
```

---

## Animation Keyframes

### `slideUpBar`

```css
@keyframes slideUpBar {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}
```

---

## Font Stack

| Context | Fonts |
|---------|-------|
| Code blocks | `'Ubuntu Mono', 'JetBrains Mono', ui-monospace, monospace` |
| Tool buttons | `'Ubuntu Mono', ui-monospace, monospace` |
| Inline code | `'JetBrains Mono', 'Fira Code', ui-monospace, monospace` |

---

## Color Constants (Hardcoded in CSS)

These are the few hardcoded HSL values used in code blocks (not derived from theme variables):

| Color | HSL | Usage |
|-------|-----|-------|
| Block background | `220, 14%, 11%` | `.code-block-wrapper` |
| Header background | `220, 14%, 14%` | `.code-block-header` |
| Line numbers background | `220, 14%, 9%` | `.code-line-numbers` |
| Header border | `220, 13%, 20%` | Header bottom border |
| Wrapper border | `220, 13%, 22%` | Outer border |
| Button background | `220, 13%, 20%` | `.code-tool-btn` |
| Button border | `220, 13%, 25%` | `.code-tool-btn` border |
| Button hover bg | `220, 13%, 28%` | Hover state |
| Line number color | `220, 10%, 35%` | Muted line numbers |
| Line count color | `220, 10%, 45%` | "N lines" text |
| Button text color | `220, 10%, 65%` | Default button text |
| Hover line bg | `220, 15%, 16%` | Line highlight |
| Hover line-num bg | `220, 15%, 12%` | Line number highlight |
| Copied feedback | `152, 70%, 50%` | Check icon color |
| Copied bg | `152, 60%, 18%` | Copied button bg |
| Copied text | `152, 70%, 60%` | Copied button text |
| Copied border | `152, 50%, 30%` | Copied button border |

---

*Styling — updated: 2026-04-08*

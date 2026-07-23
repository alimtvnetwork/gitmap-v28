# Accessibility

## Purpose

Ensure every interface — web UI and CLI — is usable by people with
disabilities. Accessibility is a requirement, not an enhancement.
These rules apply to all new features and retroactive fixes.

## Semantic HTML

### Rules

1. Use the **correct element** for the job — `<button>` for actions,
   `<a>` for navigation, `<nav>`, `<main>`, `<article>`, `<aside>`.
2. Never use `<div>` or `<span>` as interactive elements.
3. Headings follow a **strict hierarchy** — one `<h1>` per page,
   no skipped levels (`h1 → h3`).
4. Lists use `<ul>`/`<ol>`, not styled `<div>` stacks.
5. Tables use `<th scope="col|row">` and `<caption>`.

### Correct vs Incorrect

```tsx
// ✗ Wrong — div as button
<div onClick={handleClick} className="btn">Save</div>

// ✓ Correct — semantic button
<button onClick={handleClick} className="btn">Save</button>
```

```tsx
// ✗ Wrong — skipped heading level
<h1>Dashboard</h1>
<h3>Recent Activity</h3>

// ✓ Correct — sequential headings
<h1>Dashboard</h1>
<h2>Recent Activity</h2>
```

## ARIA Patterns

### Rules

1. **First rule of ARIA** — don't use ARIA if a native element works.
2. Every interactive custom widget needs `role`, `aria-label`, and
   keyboard handling.
3. Dynamic content updates use `aria-live="polite"` (non-urgent) or
   `aria-live="assertive"` (urgent errors).
4. Toggle states use `aria-expanded`, `aria-pressed`, or `aria-checked`.
5. Loading states use `aria-busy="true"` on the container.
6. Never use `aria-label` on non-interactive elements that already
   have visible text.

### Common Patterns

| Widget         | Role              | Required ARIA                        |
|----------------|-------------------|--------------------------------------|
| Accordion      | `region`          | `aria-expanded`, `aria-controls`     |
| Modal          | `dialog`          | `aria-modal`, `aria-labelledby`      |
| Tab panel      | `tablist`/`tab`   | `aria-selected`, `aria-controls`     |
| Toast          | `status`          | `aria-live="polite"`                 |
| Error alert    | `alert`           | `aria-live="assertive"`              |
| Search input   | `search`          | `aria-label` (when no visible label) |

### Pattern Example

```tsx
<div
  role="alert"
  aria-live="assertive"
  className="text-destructive"
>
  {errorMessage}
</div>
```

## Keyboard Navigation

### Rules

1. Every interactive element must be **focusable** and **operable**
   with keyboard alone.
2. Focus order follows **visual reading order** (no positive `tabIndex`).
3. Custom widgets implement standard key bindings:
   - `Enter`/`Space` — activate.
   - `Escape` — close/dismiss.
   - `Arrow keys` — navigate within composite widgets.
   - `Tab` — move between widgets.
4. **Focus trap** inside modals — `Tab` cycles within, `Escape` closes.
5. **Visible focus indicator** — never `outline: none` without a
   replacement. Use `focus-visible` ring styles.
6. Skip links: provide "Skip to main content" as the first focusable
   element on pages with navigation.

### Focus Management

```tsx
// Trap focus in modal
useEffect(() => {
  if (!isOpen) return;
  const prev = document.activeElement as HTMLElement;
  dialogRef.current?.focus();
  return () => prev?.focus();
}, [isOpen]);
```

### CLI Keyboard Rules

1. Interactive TUI respects `Escape` to cancel, `Enter` to confirm.
2. Arrow keys navigate lists; typed characters filter.
3. `Ctrl+C` always exits gracefully.

## Color Contrast

### Rules

1. **WCAG AA minimum** — 4.5:1 for normal text, 3:1 for large text
   (18px+ or 14px+ bold).
2. **Never convey information by colour alone** — pair with icons,
   text labels, or patterns.
3. Use the project's **semantic design tokens** (`--foreground`,
   `--muted-foreground`, `--destructive`) which are pre-validated
   for contrast.
4. Interactive states (hover, focus, active) must also meet contrast
   ratios.
5. Test both **light and dark modes** — contrast can differ.

### Validation Checklist

| Check                    | Tool / Method                    |
|--------------------------|----------------------------------|
| Text vs background       | Browser DevTools contrast picker |
| Icon-only buttons        | Must have `aria-label`           |
| Disabled states          | 3:1 minimum against background   |
| Error states             | Red + icon + text, not red alone |
| Charts / data viz        | Patterns + labels + colour       |

## Screen Reader Support

### Rules

1. All images have **meaningful `alt` text** or `alt=""` if decorative.
2. Icon-only buttons have `aria-label` or visually hidden text.
3. Form inputs have associated `<label>` elements (explicit `htmlFor`
   or wrapping).
4. Page title (`<title>`) updates on route change for SPAs.
5. Route transitions announce new content via `aria-live` region.
6. SVG icons use `aria-hidden="true"` when paired with text.

### Visually Hidden Utility

```css
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
```

### Image Alt Text Rules

| Image Type    | Alt Text                                      |
|---------------|-----------------------------------------------|
| Informative   | Describe content: `"Bar chart of monthly revenue"` |
| Decorative    | Empty: `alt=""`                               |
| Functional    | Describe action: `"Close dialog"`             |
| Complex       | Short alt + long description via `aria-describedby` |

## Testing Requirements

| Method              | Frequency      | Tool                        |
|---------------------|----------------|-----------------------------|
| Automated audit     | Every PR       | axe-core / Lighthouse       |
| Keyboard walkthrough| Every feature  | Manual                      |
| Screen reader test  | Major releases | NVDA / VoiceOver            |
| Contrast check      | Every PR       | Browser DevTools             |
| Heading outline     | Every page     | HeadingsMap extension        |

## Constraints

- Files ≤ 200 lines, functions 8–15 lines.
- No `div` or `span` as interactive elements.
- No `outline: none` without a visible focus replacement.
- No colour-only information signalling.
- All forms must have associated labels.
- WCAG AA is the minimum standard.

## References

- [03 Naming Conventions](../05-coding-guidelines/03-naming-conventions.md)
- [12 Documentation Standards](../05-coding-guidelines/12-documentation-standards.md)
- [19 Internationalization](../05-coding-guidelines/19-internationalization.md)

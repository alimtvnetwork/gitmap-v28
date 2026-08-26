# Sidebar System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The sidebar provides file tree navigation, search, and context. It uses theme tokens for all colors and supports collapse, mobile sheet, and keyboard toggling.

---

## Layout

| Property | Value |
|----------|-------|
| Width (expanded) | `16rem` (256px) |
| Width (collapsed) | Icon-only via `data-state="collapsed"` |
| Position | Fixed left, full height |
| Background | `hsl(var(--sidebar-background))` |
| Border right | `1px solid hsl(var(--sidebar-border))` |
| Toggle shortcut | `Ctrl+B` / `Cmd+B` |

---

## File Tree Items

### Default State

```css
.tree-item {
  color: hsl(var(--sidebar-foreground));
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
  border-radius: calc(var(--radius) - 2px);
  transition: background 0.15s ease, color 0.15s ease;
}
```

### Hover State

```css
.tree-item:hover {
  background: hsl(var(--sidebar-accent));
  color: hsl(var(--sidebar-accent-foreground));
}
```

### Active/Selected State

```css
.tree-item.active {
  background: hsl(var(--sidebar-primary));
  color: hsl(var(--sidebar-primary-foreground));
  font-weight: 600;
}
```

---

## Folder Expand/Collapse

- Folders use chevron icons that rotate on expand
- Transition: `transform 0.2s ease` on the chevron
- Collapsed: `rotate(0deg)`, Expanded: `rotate(90deg)`
- Children animate via accordion: `0.2s ease-out`

---

## Search

```css
.sidebar-search {
  background: hsl(var(--sidebar-background));
  border: 1px solid hsl(var(--sidebar-border));
  border-radius: var(--radius);
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
  color: hsl(var(--sidebar-foreground));
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.sidebar-search:focus {
  border-color: hsl(var(--sidebar-ring));
  box-shadow: 0 0 0 2px hsl(var(--sidebar-ring) / 0.2);
}
```

- Case-insensitive filtering on filenames and content
- Maximum 20 results
- Highlights matched text in results

---

## Mobile Behavior

On screens < 768px:
- Sidebar becomes a `Sheet` overlay (slides from left)
- Overlay: `hsl(0 0% 0% / 0.5)` backdrop
- Triggered by hamburger button
- Closes on file selection or outside click

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Sidebar Component | `src/components/ui/sidebar.tsx` |
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |

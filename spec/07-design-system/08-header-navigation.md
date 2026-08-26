# Header & Navigation

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The header and navigation system defines layout, spacing, interaction patterns, and hover behavior for the site's primary navigation. All colors reference theme tokens. All transitions use CSS3.

---

## Header Layout

```
┌──────────────────────────────────────────────────────────┐
│  Logo/Title     Nav Items...        Actions (icons/btns) │
└──────────────────────────────────────────────────────────┘
```

| Element | Alignment | Spacing |
|---------|-----------|---------|
| Logo/Title | Left | `gap: 0.5rem` from nav |
| Nav Items | Center or left-aligned after logo | `gap: 0.5rem` between items |
| Action Icons | Right | `gap: 0.35rem` between icons |

### Structure

```css
header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  background: hsl(var(--background));
  border-bottom: 1px solid hsl(var(--border));
  position: sticky;
  top: 0;
  z-index: 50;
}
```

---

## Navigation Menu Items

### Default State

```css
.nav-item {
  color: hsl(var(--foreground) / 0.8);
  font-family: 'Ubuntu', sans-serif;
  font-size: 0.85rem;
  font-weight: 500;
  padding: 0.4rem 0.75rem;
  border-radius: var(--radius);
  position: relative;
  transition: color 0.2s ease, background 0.2s ease;
}
```

### Hover State

```css
.nav-item:hover {
  color: hsl(var(--foreground));
  background: hsl(var(--muted) / 0.5);
}
```

### Active State

```css
.nav-item.active {
  color: hsl(var(--primary));
  background: hsl(var(--primary) / 0.08);
}
```

---

## Menu Hover Underline

Navigation items use the **right-to-left underline sweep** pattern (same as prose links):

```css
.nav-item::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  bottom: 0;
  left: 0;
  background: linear-gradient(
    90deg,
    hsl(var(--heading-gradient-from)),
    hsl(var(--heading-gradient-to))
  );
  transform: scaleX(0);
  transform-origin: bottom right;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 1px;
}

.nav-item:hover::after {
  transform: scaleX(1);
  transform-origin: bottom left;
}
```

---

## Icon Hover Transitions

Header action icons (theme toggle, fullscreen, shortcuts) follow this pattern:

```css
.header-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border-radius: var(--radius);
  color: hsl(var(--muted-foreground));
  transition: color 0.2s ease, background 0.2s ease, transform 0.15s ease;
}

.header-icon:hover {
  color: hsl(var(--foreground));
  background: hsl(var(--muted) / 0.5);
  transform: scale(1.05);
}

.header-icon:active {
  transform: scale(0.95);
}
```

**No bounce, no spring** — just a subtle scale pulse on hover.

---

## Dropdown / Search Menu

### Dropdown Container

```css
.dropdown-menu {
  background: hsl(var(--popover));
  border: 1px solid hsl(var(--border));
  border-radius: var(--radius);
  box-shadow: 0 4px 16px hsl(var(--foreground) / 0.08);
  padding: 0.25rem;
  animation: fade-in 0.15s ease-out;
}
```

### Dropdown Item Hover

```css
.dropdown-item {
  padding: 0.4rem 0.75rem;
  border-radius: calc(var(--radius) - 2px);
  color: hsl(var(--foreground));
  transition: background 0.15s ease, color 0.15s ease;
}

.dropdown-item:hover {
  background: hsl(var(--primary) / 0.08);
  color: hsl(var(--primary));
}
```

This is the preferred hover behavior — subtle primary-tinted background shift with text color change.

---

## Breadcrumb Navigation

```
📖  Spec authoring guide  ›  Cli module template
```

| Element | Style |
|---------|-------|
| Icon | `BookOpen` from Lucide, `hsl(var(--primary))` |
| Path segments | `hsl(var(--muted-foreground))`, `0.8rem` |
| Separator | `›` character |
| Current page | `hsl(var(--foreground))`, `font-weight: 500` |

---

## View Mode Tabs

Preview / Edit / Split tabs in the header:

```css
.view-tab {
  font-size: 0.8rem;
  padding: 0.3rem 0.75rem;
  border-radius: var(--radius);
  color: hsl(var(--muted-foreground));
  transition: all 0.2s ease;
}

.view-tab:hover {
  color: hsl(var(--foreground));
  background: hsl(var(--muted) / 0.5);
}

.view-tab.active {
  color: hsl(var(--primary-foreground));
  background: hsl(var(--primary));
}
```

---

## Responsive Behavior

| Breakpoint | Header Change |
|-----------|--------------|
| Desktop (> 1024px) | Full nav items visible, icon buttons with labels |
| Tablet (768-1024px) | Nav items visible, icon-only buttons |
| Mobile (< 768px) | Hamburger menu, sidebar becomes sheet overlay |

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Motion Transitions | [06-motion-transitions.md](./06-motion-transitions.md) |
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
| Button System | [09-button-system.md](./09-button-system.md) |

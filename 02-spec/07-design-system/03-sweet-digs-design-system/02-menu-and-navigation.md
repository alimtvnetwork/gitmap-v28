# Sweet Digs Menu & Navigation System

> **/goal** Specify the layout, sticky behavior, and micro-interactions of the Sweet Digs navigation system, focusing on the 75% expanding underline animation on hover.
> **/learn** Master the simultaneous 300ms transitions for text color, background tint, and centered pseudo-element underline expansion.

---

## 1. Header Layout Architecture

The Sweet Digs header floats at the top of the viewport with a sticky position and subtle glassmorphic backdrop blur:

```css
.site-header {
  position: sticky;
  top: 0;
  z-index: 50;
  width: 100%;
  height: 70px;
  display: flex;
  align-items: center;
  background-color: hsl(var(--background) / 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid hsl(var(--border) / 0.5);
  transition: background-color 300ms ease, border-color 300ms ease;
}

.header-container {
  max-width: 1200px;
  width: 100%;
  margin-inline: auto;
  padding-inline: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
```

---

## 2. Nav Item Default State

Navigation links operate as relative containers so the animated `::after` underline can position precisely along the bottom edge:

```css
.nav-link {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;              /* 6px spacing for icons/chevrons */
  padding: 0.5rem 0.875rem;    /* 8px vertical, 14px horizontal */
  font-size: 0.875rem;        /* 14px text-sm */
  font-weight: 500;           /* medium weight */
  color: hsl(var(--foreground) / 0.8);
  border-radius: var(--radius);
  text-decoration: none;
  cursor: pointer;
  background-color: transparent;
  transition: color 300ms ease, background-color 300ms ease;
}
```

### The Symmetrical Underline Pseudo-Element

The underline is anchored at the exact horizontal center (`left: 50%; transform: translateX(-50%);`) with an initial width of 0:

```css
.nav-link::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  background-color: hsl(var(--primary));
  border-radius: 9999px;
  transition: width 300ms ease;
  pointer-events: none;
}
```

---

## 3. Hover State: Simultaneous 3-Part Transition

When hovered, three distinct visual shifts execute simultaneously over **300ms ease**:

1. **Text Color:** Transitions from subdued text (`hsl(var(--foreground) / 0.8)`) to pure brand accent (`hsl(var(--primary))`).
2. **Background Tint:** Transitions from transparent to soft secondary mint tint (`hsl(var(--secondary))`).
3. **Underline Expansion:** The `::after` pseudo-element expands horizontally from center from `0` to `75%` of the item width.

```css
/* Hover State */
.nav-link:hover {
  color: hsl(var(--primary));
  background-color: hsl(var(--secondary));
}

.nav-link:hover::after {
  width: 75%;
}

/* Active Page State */
.nav-link.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
}

.nav-link.is-active::after {
  width: 75%;
}
```

### Transition Specifications Matrix

| Element Property | Default State | Hover State | Transition Timing |
|:---|:---|:---|:---|
| `color` | `hsl(var(--foreground) / 0.8)` | `hsl(var(--primary))` | `300ms ease` |
| `background-color` | `transparent` | `hsl(var(--secondary))` | `300ms ease` |
| `::after width` | `0` | `75%` | `300ms ease` |
| `::after left / transform` | `50% / translateX(-50%)` | `50% / translateX(-50%)` | Maintained (centered) |

---

## 4. Dropdown Chevron Icon Behavior

Items featuring dropdown menus display a micro-chevron icon that rotates smoothly on interaction:

```html
<button class="nav-link" aria-expanded="false">
  <span>Solutions</span>
  <svg class="chevron-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
    <polyline points="6 9 12 15 18 9"></polyline>
  </svg>
</button>
```

```css
.chevron-icon {
  opacity: 0.6;
  transition: transform 300ms ease, opacity 300ms ease;
}

.nav-link:hover .chevron-icon {
  opacity: 1;
}

.nav-link.has-dropdown-open .chevron-icon {
  transform: rotate(180deg);
}
```

---

## 5. Mobile Responsive Header & Drawer

On viewports below 768px (`< md`):
1. Navigation items collapse into a slide-over mobile drawer or sheet.
2. The trigger button displays an animated hamburger-to-close icon.
3. Mobile drawer items preserve the 75% expanding underline interaction on tap/hover.
4. Positive boolean tracking (`isMenuOpen`) dictates drawer visibility without mixed polarities.

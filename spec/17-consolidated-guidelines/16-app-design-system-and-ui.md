# Consolidated: App Design System & UI

**Version:** 3.3.0  
**Updated:** 2026-04-16

---

## Purpose

Application-specific design system and UI specifications. This module covers component patterns, theming decisions, layout conventions, and visual standards **unique to the application** — extending the foundational design system's core tokens. An AI reading only this file and `07-design-system.md` must be able to reproduce every app-specific UI element.

---

## Relationship to Core Design System

The core design system (`07-design-system`) defines foundational CSS variables, motion standards, and theming primitives. This module builds **on top of** those foundations. It does **not** duplicate or override core tokens — it extends them with app-specific layout, component variants, error UI colors, and interaction patterns.

---

## App-Specific Layout Tokens

```css
:root {
  --app-sidebar-width: 280px;
  --app-sidebar-collapsed-width: 0px;
  --app-header-height: 64px;
  --app-content-max-width: 1200px;
  --app-content-padding: 2rem;
}
```

### App Shell Layout

```
┌────────────────────────────────────────────────────┐
│ Header (--app-header-height: 64px)                 │
├──────────────┬─────────────────────────────────────┤
│ Sidebar      │          Main Content               │
│ (280px)      │   (max-width: 1200px, centered)     │
│              │   padding: 2rem                      │
└──────────────┴─────────────────────────────────────┘
```

---

## App-Specific Semantic Tokens

These tokens extend the core design system for app-specific UI needs:

```css
:root {
  /* Sidebar */
  --sidebar: 0 0% 98%;
  --sidebar-foreground: 240 5.3% 26.1%;
  --sidebar-primary: 240 5.9% 10%;
  --sidebar-primary-foreground: 0 0% 98%;
  --sidebar-accent: 240 4.8% 95.9%;
  --sidebar-accent-foreground: 240 5.9% 10%;
  --sidebar-border: 220 13% 91%;
  --sidebar-ring: 217.2 91.2% 59.8%;

  /* Status indicators */
  --success: 142 76% 36%;
  --success-foreground: 0 0% 100%;
  --warning: 38 92% 50%;
  --warning-foreground: 0 0% 100%;

  /* Code viewer */
  --code-bg: 220 13% 95%;
  --code-border: 220 13% 87%;
  --code-line-highlight: 210 40% 96%;
  --code-selection: 214 95% 93%;
}

.dark {
  --sidebar: 240 5.9% 10%;
  --sidebar-foreground: 240 4.8% 95.9%;
  --sidebar-accent: 240 3.7% 15.9%;
  --sidebar-accent-foreground: 240 4.8% 95.9%;
  --sidebar-border: 240 3.7% 15.9%;

  --success: 142 70% 45%;
  --success-foreground: 0 0% 100%;
  --warning: 38 92% 50%;
  --warning-foreground: 0 0% 100%;

  --code-bg: 240 6% 10%;
  --code-border: 240 4% 16%;
  --code-line-highlight: 240 6% 14%;
  --code-selection: 214 60% 20%;
}
```

### Tailwind Config Mapping for App Tokens

```typescript
// tailwind.config.ts — extend the colors object
sidebar: {
  DEFAULT: "hsl(var(--sidebar))",
  foreground: "hsl(var(--sidebar-foreground))",
  primary: "hsl(var(--sidebar-primary))",
  "primary-foreground": "hsl(var(--sidebar-primary-foreground))",
  accent: "hsl(var(--sidebar-accent))",
  "accent-foreground": "hsl(var(--sidebar-accent-foreground))",
  border: "hsl(var(--sidebar-border))",
  ring: "hsl(var(--sidebar-ring))",
},
success: {
  DEFAULT: "hsl(var(--success))",
  foreground: "hsl(var(--success-foreground))",
},
warning: {
  DEFAULT: "hsl(var(--warning))",
  foreground: "hsl(var(--warning-foreground))",
},
```

---

## Typography Scale

### App-Specific Type Tokens

| Token | Size | Weight | Line Height | Use Case |
|-------|------|--------|-------------|----------|
| `text-page-title` | `1.875rem` (30px) | `700` (bold) | `1.2` | Page headings |
| `text-section-title` | `1.25rem` (20px) | `600` (semibold) | `1.3` | Section headings |
| `text-card-title` | `1rem` (16px) | `600` (semibold) | `1.4` | Card headings |
| `text-body` | `0.875rem` (14px) | `400` (normal) | `1.6` | Body text |
| `text-caption` | `0.75rem` (12px) | `400` (normal) | `1.5` | Captions, timestamps |
| `text-code` | `0.8125rem` (13px) | `400` (normal) | `1.7` | Monospaced code |

### Font Stack

```css
:root {
  --font-sans: "Inter", ui-sans-serif, system-ui, sans-serif;
  --font-mono: "JetBrains Mono", "Fira Code", ui-monospace, monospace;
}
```

---

## Spacing System

The app uses Tailwind's default 4px-based spacing scale. App-specific spacing conventions:

| Context | Spacing | Tailwind Class |
|---------|---------|----------------|
| Between sections | `32px` | `space-y-8` |
| Between cards | `16px` | `gap-4` |
| Card internal padding | `24px` | `p-6` |
| Form field gap | `16px` | `space-y-4` |
| Inline icon-to-text | `8px` | `gap-2` |
| Button icon-to-label | `8px` | `gap-2` |
| Sidebar item padding | `8px 12px` | `px-3 py-2` |
| Header horizontal padding | `24px` | `px-6` |
| Modal body padding | `24px` | `p-6` |
| Toast padding | `16px` | `p-4` |

---

## Z-Index Scale

| Layer | z-index | Use Case |
|-------|---------|----------|
| Content (base) | `0` | Default stacking |
| Sticky elements | `10` | Sticky headers, floating toolbars |
| Sidebar | `20` | Sidebar overlay on mobile |
| Dropdown/Popover | `30` | Dropdown menus, popovers, tooltips |
| Modal overlay | `40` | Dialog backdrop |
| Modal content | `50` | Dialog content |
| Toast | `60` | Toast notifications |
| Global error modal | `70` | Error queue overlay (always on top) |

**Rule:** Never use arbitrary z-index values. All layers must use this scale.

---

## Shadow Scale

| Token | Value | Use Case |
|-------|-------|----------|
| `--shadow-sm` | `0 1px 2px 0 rgb(0 0 0 / 0.05)` | Subtle card depth |
| `--shadow-md` | `0 4px 6px -1px rgb(0 0 0 / 0.1)` | Card hover, dropdowns |
| `--shadow-lg` | `0 10px 15px -3px rgb(0 0 0 / 0.1)` | Modals, floating panels |
| `--shadow-glow` | `0 0 15px hsl(var(--primary) / 0.15)` | Primary button hover glow |

Dark mode shadows use lower opacity (`0.3` → `0.15`) since dark backgrounds need subtler depth cues.

---

## Component Variants

### Button Variants

| Variant | Token Mapping | Use Case | Hover Behavior |
|---------|--------------|----------|----------------|
| `default` | `bg-primary text-primary-foreground` | Primary actions | `translateY(-1px)` + glow shadow |
| `secondary` | `bg-secondary text-secondary-foreground` | Secondary actions | `bg-secondary/80` |
| `outline` | `border-border text-foreground` | Tertiary actions | `bg-accent text-accent-foreground` |
| `ghost` | transparent | Toolbar/nav actions | `bg-accent text-accent-foreground` |
| `destructive` | `bg-destructive text-destructive-foreground` | Delete/danger actions | `bg-destructive/90` |
| `link` | `text-primary underline` | Inline text links | Underline sweep animation |

All buttons share: `focus-visible: 2px solid hsl(var(--ring))`, `disabled: opacity-50 pointer-events-none`.

### Button Sizes

| Size | Padding | Height | Font Size | Icon Size |
|------|---------|--------|-----------|-----------|
| `sm` | `px-3` | `h-8` (32px) | `text-xs` | 14px |
| `default` | `px-4 py-2` | `h-10` (40px) | `text-sm` | 16px |
| `lg` | `px-8` | `h-11` (44px) | `text-base` | 18px |
| `icon` | `p-0` | `h-10 w-10` | — | 20px |

### Badge Variants

| Variant | Token Mapping | Use Case |
|---------|--------------|----------|
| `default` | `bg-primary text-primary-foreground` | Status indicators |
| `secondary` | `bg-secondary text-secondary-foreground` | Category labels |
| `outline` | `border-border text-foreground` | Subtle tags |
| `destructive` | `bg-destructive text-destructive-foreground` | Error badges |
| `success` | `bg-success/15 text-success border-success/30` | Success status |
| `warning` | `bg-warning/15 text-warning border-warning/30` | Warning status |

### Card Component

| Property | Token |
|----------|-------|
| Background | `bg-card` |
| Text | `text-card-foreground` |
| Border | `border` (1px solid) |
| Radius | `rounded-lg` (var(--radius)) |
| Hover | `shadow-md`, optional `translateY(-1px)` |

### Card Sub-Components

| Part | Token Mapping |
|------|--------------|
| `CardHeader` | `p-6 pb-0` |
| `CardTitle` | `text-card-foreground font-semibold leading-none tracking-tight` |
| `CardDescription` | `text-muted-foreground text-sm` |
| `CardContent` | `p-6 pt-0` |
| `CardFooter` | `p-6 pt-0 flex items-center` |

### Input / Form Elements

| Property | Token |
|----------|-------|
| Background | `bg-background` |
| Border | `border-input` |
| Focus | `ring-2 ring-ring ring-offset-2` |
| Placeholder | `text-muted-foreground` |
| Disabled | `opacity-50 cursor-not-allowed` |
| Radius | `rounded-md` |

### Form Validation States

| State | Border | Message Color | Icon |
|-------|--------|---------------|------|
| Default | `border-input` | — | — |
| Error | `border-destructive` | `text-destructive text-sm` | `AlertCircle` (14px) |
| Success | `border-success` | `text-success text-sm` | `CheckCircle` (14px) |
| Warning | `border-warning` | `text-warning text-sm` | `AlertTriangle` (14px) |

Error messages appear below the input with `mt-1.5` spacing. Icon sits inline before the message text.

### Select / Combobox

| Property | Token |
|----------|-------|
| Trigger | Same as Input tokens + `cursor-pointer` |
| Content | `bg-popover text-popover-foreground border rounded-md shadow-md` |
| Item hover | `bg-accent text-accent-foreground` |
| Selected item | `bg-accent font-medium` + `Check` icon (16px) |
| Placeholder | `text-muted-foreground` |
| Scroll buttons | `text-muted-foreground` centered |

### Tabs Component

| Part | Token Mapping |
|------|--------------|
| TabsList | `bg-muted rounded-md p-1` |
| TabsTrigger (inactive) | `text-muted-foreground` |
| TabsTrigger (active) | `bg-background text-foreground shadow-sm` |
| TabsContent | `mt-2` |

### Table Component

| Part | Token Mapping |
|------|--------------|
| Header row | `bg-muted/50 text-muted-foreground font-medium text-sm` |
| Body row | `border-b border-border` |
| Row hover | `bg-muted/50` |
| Cell padding | `p-4` |
| Empty state | `text-muted-foreground text-center py-10` |

### Dialog / Modal

| Property | Token |
|----------|-------|
| Overlay | `bg-black/80` with `backdrop-blur-sm` |
| Content bg | `bg-background` |
| Content border | `border` |
| Content radius | `rounded-lg` |
| Title | `text-foreground font-semibold text-lg` |
| Description | `text-muted-foreground text-sm` |
| Footer | `flex justify-end gap-2 pt-4` |

### Toast / Notification

| Variant | Token Mapping |
|---------|--------------|
| Default | `bg-background border text-foreground` |
| Destructive | `bg-destructive text-destructive-foreground border-destructive` |
| Success | `bg-success text-success-foreground` |
| Warning | `bg-warning text-warning-foreground` |

Enter: `300ms ease-out translateY + opacity`. Exit: `200ms ease-in translateY + opacity`.

### Sidebar Item States

| State | Visual Treatment |
|-------|-----------------|
| Default | `text-sidebar-foreground` |
| Hover | `bg-sidebar-accent/50 text-sidebar-accent-foreground` |
| Active | `bg-sidebar-accent text-sidebar-accent-foreground` + left accent border |
| Expanded folder | Chevron rotated 90° (`200ms ease-out`) |
| Collapsed folder | Chevron at 0° |

### Dropdown Menu

| Property | Token |
|----------|-------|
| Background | `bg-popover` |
| Text | `text-popover-foreground` |
| Border | `border` |
| Item hover | `bg-accent text-accent-foreground` |
| Item destructive | `text-destructive` on hover: `bg-destructive/10` |
| Separator | `bg-muted` |
| Shortcut text | `text-muted-foreground text-xs ml-auto` |
| Disabled item | `opacity-50 pointer-events-none` |
| Sub-trigger | Chevron right icon (14px) |

### Tooltip

| Property | Token |
|----------|-------|
| Background | `bg-primary` |
| Text | `text-primary-foreground` |
| Radius | `rounded-md` |
| Font size | `text-xs` |
| Animation | `fade-in 150ms ease-out` |
| Max width | `max-w-xs` (320px) |
| Delay | `200ms` open, `0ms` close |

### Skeleton / Loading Placeholder

| Property | Token |
|----------|-------|
| Background | `bg-muted` |
| Animation | `pulse` (2s infinite) — alternates `opacity: 1` ↔ `0.5` |
| Radius | Matches target element (`rounded-md` for text, `rounded-full` for avatars) |

### Avatar

| Size | Dimensions | Font Size | Use Case |
|------|-----------|-----------|----------|
| `sm` | `h-6 w-6` | `text-xs` | Inline mentions, compact lists |
| `default` | `h-10 w-10` | `text-sm` | User cards, comments |
| `lg` | `h-16 w-16` | `text-lg` | Profile headers |

Fallback: `bg-muted text-muted-foreground` with initials centered.

### Progress / Loading Bar

| Property | Token |
|----------|-------|
| Track | `bg-secondary h-2 rounded-full` |
| Fill | `bg-primary rounded-full` |
| Indeterminate | `animate-indeterminate` — fill slides left-to-right continuously |
| Width | `100%` of container |

---

## Error UI Color System

### Error Level Badge Colors (GlobalErrorModal header)

```typescript
const levelColors = {
  error: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
  warn:  "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  info:  "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
};
```

### Error Level Icons

```typescript
const levelIcons = { error: AlertCircle, warn: AlertTriangle, info: Info };
const levelColors = { error: "text-red-500", warn: "text-yellow-500", info: "text-blue-500" };
```

### Header Icon Color (GlobalErrorModal)

```tsx
<AlertCircle className={cn(
  "h-5 w-5 sm:h-6 sm:w-6 shrink-0",
  selectedError.level === LogLevel.Error ? "text-destructive"
    : selectedError.level === LogLevel.Warn ? "text-warning" : "text-muted-foreground"
)} />
```

### Two-Tier Stack Trace Colors

| Tier | Icon Color | Background | Text Color | Used For |
|------|-----------|------------|-----------|----------|
| **Go Backend** | — (Server icon) | `bg-muted` | `text-blue-500 dark:text-blue-400` | Go stack traces, methods stack |
| **PHP / Delegated** | `text-orange-500` (AlertTriangle) | `bg-orange-500/5` | `text-orange-500/600/700` | PHP frames, delegated service errors |

### Frontend Section Colors

| Element | Classes |
|---------|---------|
| Trigger badge | `bg-primary/5 border-primary/20` |
| Source badge | `variant="secondary" font-mono` |
| First call chain entry | `text-primary font-semibold` |
| React execution chain area | `bg-blue-500/5` border |
| First parsed frame row | `bg-primary/5`, `text-primary font-semibold` |
| Internal frames | `opacity-50` |
| Fixes numbering | `bg-primary/10 text-primary` |
| Request chain dots | `bg-blue-500` (React→Go), `bg-orange-500` (Go→PHP) |

### Error History Drawer

| Element | Classes |
|---------|---------|
| Selected item | `bg-accent border-primary` |
| Default item | `bg-card hover:bg-accent/50` |
| Header icon | `text-destructive` |
| Delete buttons | `text-destructive hover:text-destructive` |

### Error Queue Badge

```tsx
<Button className="relative text-destructive hover:text-destructive hover:bg-destructive/10">
  <Badge variant="destructive" className="absolute -top-1 -right-1">{count}</Badge>
</Button>
```

### BackendLogEntry (Semantic Tokens Only)

| Level | Background | Border |
|-------|-----------|--------|
| Error | `bg-destructive/10` | `border-destructive/30` |
| Warning | `bg-warning/10` | — |
| Info | `bg-primary/10` | — |

### Error UI Color Rules

1. **Never use raw color classes** — use design tokens or documented tier colors
2. **Tier colors are fixed**: Blue = Go session frames, Orange = PHP/Delegated, Neutral = Go raw stacks
3. **Error levels**: `text-destructive` (error), `text-warning` (warn), `text-muted-foreground` (info)
4. **Backgrounds use low opacity**: `bg-destructive/5`, `bg-orange-500/5`, `bg-blue-500/5`, `bg-primary/5`
5. **Borders use medium opacity**: `border-destructive/30`, `border-orange-500/30`
6. **Dark mode overrides**: `dark:text-amber-400`, `dark:text-blue-400`, `dark:text-orange-400`

---

## Animation Specifications

### Page Transitions

| Animation | Duration | Easing | Property |
|-----------|----------|--------|----------|
| Content fade-in | `200ms` | `ease-out` | `opacity` |
| Sidebar slide | `300ms` | `ease-out` | `transform: translateX` |
| Modal overlay | `200ms` | `ease-out` | `opacity` |

### Micro-Interactions

| Animation | Duration | Easing | Property |
|-----------|----------|--------|----------|
| Button hover | `150ms` | `ease-out` | `background-color` |
| Link underline | `150ms` | `ease-out` | `scaleX` via `::after` |
| Chevron rotation | `200ms` | `ease-out` | `transform: rotate` |
| Focus ring | `150ms` | `ease-out` | `box-shadow` |
| Toast enter | `300ms` | `ease-out` | `translateY + opacity` |
| Toast exit | `200ms` | `ease-in` | `translateY + opacity` |

### Code Block Animations

| Animation | Duration | Easing | Property |
|-----------|----------|--------|----------|
| Line highlight | `150ms` | `ease-out` | `background-color` |
| Selection bar appear | `200ms` | `ease-out` | `opacity + translateY` |
| Fullscreen enter | `300ms` | `ease-out` | `transform: scale` |
| Copy feedback | `1500ms` | — | Text swap + fade |

### Reduced Motion

When `prefers-reduced-motion: reduce` is active:

- All transitions set to `0ms`
- Skeleton pulse replaced with static `opacity: 0.5`
- Toast enter/exit: instant (no transform)
- Sidebar slide: instant

---

## Responsive Breakpoints

| Breakpoint | Width | Behavior |
|------------|-------|----------|
| Mobile | `< 768px` | Sidebar hidden, hamburger menu, full-width content |
| Tablet | `768px – 1024px` | Sidebar collapsible overlay, content adjusts |
| Desktop | `> 1024px` | Sidebar persistent, content with max-width |

### Mobile-Specific Adaptations

- Touch targets: minimum 44px
- Sidebar becomes `Sheet` overlay (slides from left)
- Modals: full-width on mobile (`max-w-full mx-4`)
- Tables: horizontal scroll wrapper
- Code blocks: horizontal scroll, font controls hidden
- Buttons: full-width on mobile for primary actions (`w-full sm:w-auto`)

---

## Icon System

- Lucide React icons throughout
- Size: `16px` inline, `20px` buttons, `24px` navigation
- Color: inherits from parent `text-*` class
- No custom SVGs — Lucide only for consistency

### Common Icon Mapping

| Context | Icon | Size |
|---------|------|------|
| Error | `AlertCircle` | 20-24px |
| Warning | `AlertTriangle` | 20-24px |
| Info | `Info` | 20-24px |
| Success | `CheckCircle` | 20-24px |
| Close | `X` | 16-20px |
| Settings | `Settings` | 20px |
| Search | `Search` | 16-20px |
| Folder | `Folder`, `FolderOpen` | 16px |
| File | `File`, `FileText` | 16px |
| Copy | `Copy`, `Check` (after copy) | 16px |
| External link | `ExternalLink` | 14px |
| Expand/Collapse | `ChevronDown`, `ChevronRight` | 16px |
| Menu | `Menu` | 24px |
| Back | `ArrowLeft` | 20px |

---

## State Language (All Interactive Elements)

| State | Visual Change | Transition |
|-------|--------------|-----------|
| Default | Base token values | — |
| Hover | Lighten bg, subtle lift, glow | `150-200ms ease-out` |
| Active/Pressed | Slightly darker than hover | Instant |
| Focus | Ring outline (`--ring` token, 2px offset) | Instant |
| Disabled | 50% opacity, `pointer-events: none` | — |
| Loading | Spinner icon, disabled state | — |
| Selected | Primary accent color, stronger border | `150ms ease-out` |
| Empty | Muted text centered, optional illustration | — |
| Error | Destructive border + message below | `150ms ease-out` |

---

## Dark Mode Rules

1. **Token-only switching**: Components never check dark mode directly — tokens handle everything
2. **No `dark:` prefixes in components** except for documented error UI tier colors
3. **Shadows reduce opacity** in dark mode (core `--shadow-*` tokens handle this)
4. **Borders become subtler** — dark mode borders use lower-contrast values
5. **Success/Warning** tokens have slightly boosted lightness in dark mode for readability

---

*Consolidated app design system & UI — v3.3.0 — 2026-04-16*

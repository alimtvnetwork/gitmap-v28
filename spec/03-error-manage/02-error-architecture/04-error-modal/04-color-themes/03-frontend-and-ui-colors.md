# Frontend Section & UI Element Colors

> **Parent:** [Color Themes Index](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## Frontend Section Colors

| Element | Classes |
|---------|---------|
| Trigger badge | `bg-primary/5 border-primary/20` |
| Source badge | `variant="secondary" font-mono` |
| First call chain entry | `text-primary font-semibold` |
| Last click in path | `text-primary` |
| React execution chain area | `bg-blue-500/5` border |
| Activity icon | `text-blue-500` |
| First parsed frame row | `bg-primary/5`, `text-primary font-semibold` |
| Internal frames | `opacity-50` |
| Fixes numbering | `bg-primary/10 text-primary` |

---

## Section Toggle Buttons

Active: `variant="default"` (primary bg). Inactive: `variant="outline"`. Bar: `bg-muted/30`.

## Error History Drawer

- Selected item: `bg-accent border-primary`
- Default item: `bg-card hover:bg-accent/50`
- Header icon: `text-destructive`
- Delete buttons: `text-destructive hover:text-destructive`

## Error Queue Badge

```tsx
<Button className="relative text-destructive hover:text-destructive hover:bg-destructive/10">
  <Badge variant="destructive" className="absolute -top-1 -right-1">{count}</Badge>
</Button>
```

## App Error Boundary

Standard card with `bg-background`, `text-muted-foreground`, `variant="outline"` for secondary action, default variant for primary.

---

## Color Usage Rules

1. **Never use raw color classes** — use design tokens or documented tier colors
2. **Tier colors are fixed**: Blue = Go session frames, Orange = PHP/Delegated, Neutral = Go raw stacks
3. **Error levels**: `text-destructive` (error), `text-warning`/`text-amber-*` (warn), `text-muted-foreground`/`text-blue-*` (info)
4. **Backgrounds use low opacity**: `bg-destructive/5`, `bg-orange-500/5`, `bg-blue-500/5`, `bg-primary/5`
5. **Borders use medium opacity**: `border-destructive/30`, `border-orange-500/30`
6. **Dark mode overrides**: `dark:text-amber-400`, `dark:text-blue-400`, `dark:text-orange-400`
7. **BackendLogEntry**: semantic tokens only (`bg-destructive/10`, `bg-warning/10`, `bg-primary/10`)
8. **Request chain dots**: `bg-blue-500` (React→Go), `bg-orange-500` (Go→PHP)

---

*Frontend & UI element colors — updated: 2026-03-31*

# Design System Tokens & Error Level Colors

> **Parent:** [Color Themes Index](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## 1. Design System Tokens (index.css)

### Light Mode (`:root`)

```css
--destructive: 0 84% 60%;           /* Red — errors, delete actions */
--destructive-foreground: 0 0% 100%;
--warning: 38 92% 50%;              /* Amber — warnings */
--warning-foreground: 0 0% 100%;
--success: 142 76% 36%;             /* Green — success states */
--success-foreground: 0 0% 100%;
--info: 217 91% 60%;                /* Blue — informational */
--info-foreground: 0 0% 100%;
--muted: 210 40% 96.1%;
--muted-foreground: 215.4 16.3% 46.9%;
--primary: 222.2 47.4% 11.2%;
--primary-foreground: 210 40% 98%;
```

### Dark Mode (`.dark`)

```css
--destructive: 0 72% 51%;
--warning: 38 92% 50%;
--success: 120 45% 39%;
--info: 217 91% 60%;
--muted: 0 0% 18%;
--muted-foreground: 0 0% 60%;
--primary: 120 45% 39%;
--primary-foreground: 0 0% 100%;
```

---

## 2. Error Level Color Mapping

### Level Badge Colors (GlobalErrorModal header)

```typescript
const levelColors = {
  error: "bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400",
  warn:  "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400",
  info:  "bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400",
};
```

### Level Icon Colors (ErrorHistoryDrawer)

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

---

*Design tokens — updated: 2026-03-31*

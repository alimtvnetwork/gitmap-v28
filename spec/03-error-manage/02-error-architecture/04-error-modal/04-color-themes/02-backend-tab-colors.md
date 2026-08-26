# Backend Section — Tab-Specific Color Themes

> **Parent:** [Color Themes Index](./00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31

---

## Overview Tab

| Element | Tailwind Classes |
|---------|-----------------|
| Error banner | `border-destructive/30 bg-destructive/5`, `text-destructive` |
| Missing delegation warning | `border-amber-500/30 bg-amber-500/5`, `text-amber-600 dark:text-amber-400` |
| HTTP status badge (≥400) | `variant="destructive"` |
| HTTP status badge (<400) | `variant="outline"` |
| Site URL box | `bg-muted` |

## Stack Tab

| Element | Tailwind Classes |
|---------|-----------------|
| PHP Delegated Error Stack | `bg-orange-500/5`, `text-orange-700 dark:text-orange-300` |
| PHP Stack Trace (DelegatedRequestServer) | `bg-orange-500/5`, icon `text-orange-500` |
| PHP Stack Frames table header | `bg-orange-500/10` |
| First PHP frame | `text-orange-600 dark:text-orange-400 font-semibold` |
| Go Backend Stack | `bg-muted` (no color coding) |
| Session Go frames | `text-blue-500 dark:text-blue-400` |
| Session PHP frames | `text-orange-500 dark:text-orange-400` |
| PHP stacktrace.txt log | `bg-orange-500/5`, `text-orange-700 dark:text-orange-300` |

## SessionLogsTab — StackFrameRow

```tsx
function StackFrameRow({ frame, variant }: { frame: SessionStackFrame; variant: "golang" | "php" }) {
  const fnColor = variant === "golang" ? "text-blue-500 dark:text-blue-400" : "text-orange-500 dark:text-orange-400";
  return <span className={cn("font-semibold", fnColor)}>{frame.class ? `${frame.class}::${frame.function}` : frame.function}</span>;
}
```

## Request Tab — 2-Node Chain

```
Node 1: React → Go     → Blue   (bg-blue-500 dot, bg-blue-500/10 badge)
Node 2: Go → PHP       → Orange (bg-orange-500 dot, bg-orange-500/10 badge, text-orange-600)
```

## Traversal Tab

- PHP badge: `bg-orange-500/10 border-orange-500/30 text-orange-600 dark:text-orange-400`
- Delegated error stack: icon `text-orange-500`, bg `bg-orange-500/5`
- Methods stack & backend trace: `bg-muted`, first row `bg-primary/5`, first method `text-primary font-semibold`

## Log Tab (LogLine)

```typescript
// Stage headers → text-primary font-semibold
// [ERROR]/[FATAL] → text-destructive
// [WARN] → text-amber-600 dark:text-amber-400
// ✓/success → text-green-600 dark:text-green-400
```

## Execution Tab

```tsx
// BackendLogEntry levels:
log.level === LogLevel.Error && "bg-destructive/10 text-destructive"
log.level === LogLevel.Warn  && "bg-warning/10 text-warning"
log.level === LogLevel.Info  && "bg-primary/10 text-primary"
log.level === LogLevel.Debug && "bg-muted text-muted-foreground"

// Go Call Chain: first row bg-primary/5, first method text-primary font-semibold
```

---

*Backend tab colors — updated: 2026-03-31*

# Error Modal — Color Theme & Design Token Reference

> **Version:** 2.1.0  
> **Updated:** 2026-04-01
> **Status:** Active  
> **Purpose:** Definitive color mapping for every error-related UI element, verified against the actual React components. Any AI or developer must use this document to replicate the exact visual appearance.

---

## 1. Design System Tokens (index.css)

All error-related components consume these CSS custom properties from `src/index.css`. **Never hardcode raw color values in components** — always use semantic tokens.

### Light Mode (`:root`)

```css
--destructive: 0 84% 60%;           /* Red — errors, delete actions */
--destructive-foreground: 0 0% 100%; /* White text on destructive bg */
--warning: 38 92% 50%;              /* Amber — warnings */
--warning-foreground: 0 0% 100%;
--success: 142 76% 36%;             /* Green — success states */
--success-foreground: 0 0% 100%;
--info: 217 91% 60%;                /* Blue — informational */
--info-foreground: 0 0% 100%;
--muted: 210 40% 96.1%;             /* Light gray — backgrounds */
--muted-foreground: 215.4 16.3% 46.9%;
--primary: 222.2 47.4% 11.2%;       /* Dark navy — primary actions */
--primary-foreground: 210 40% 98%;
```

### Dark Mode (`.dark`)

```css
--destructive: 0 72% 51%;
--destructive-foreground: 0 0% 100%;
--warning: 38 92% 50%;
--warning-foreground: 0 0% 100%;
--success: 120 45% 39%;
--success-foreground: 0 0% 100%;
--info: 217 91% 60%;
--info-foreground: 0 0% 100%;
--muted: 0 0% 18%;
--muted-foreground: 0 0% 60%;
--primary: 120 45% 39%;
--primary-foreground: 0 0% 100%;
```

---

## 2. Error Level Color Mapping

### Level Badge Colors (GlobalErrorModal header)

Used in `GlobalErrorModal` for the error code badge:

```typescript
// src/components/errors/GlobalErrorModal.tsx
// v2.1.0: Use semantic tokens instead of hardcoded colors
const levelColors = {
  error: "bg-destructive/10 text-destructive",
  warn:  "bg-warning/10 text-warning",
  info:  "bg-primary/10 text-primary",
};
```

### Level Icon Colors (ErrorHistoryDrawer)

Used in `ErrorHistoryDrawer` for the per-item icon:

```typescript
// src/components/errors/ErrorHistoryDrawer.tsx
const levelIcons = {
  error: AlertCircle,
  warn: AlertTriangle,
  info: Info,
};

const levelColors = {
  error: "text-destructive",
  warn:  "text-warning",
  info:  "text-primary",
};
```

### Header Icon Color (GlobalErrorModal)

```tsx
// src/components/errors/GlobalErrorModal.tsx — header icon
// NOTE: Use LogLevel enum (e.g., LogLevel.Error, LogLevel.Warn) — not raw strings.
<AlertCircle className={cn(
  "h-5 w-5 sm:h-6 sm:w-6 shrink-0",
  selectedError.level === LogLevel.Error ? "text-destructive"
    : selectedError.level === LogLevel.Warn ? "text-warning" : "text-muted-foreground"
)} />
```

---

## 3. Backend Section — Tab-Specific Color Themes

### 3.1 Two-Tier Color System

**IMPORTANT:** The codebase uses a **two-tier** color system, not three:

| Tier | Icon Color | Background | Text Color | Used For |
|------|-----------|------------|-----------|----------|
| **Go Backend** | — (uses `Server` icon, no color) | `bg-muted` | `text-blue-500 dark:text-blue-400` (session frames only) | Go stack traces, raw stack, methods stack |
| **PHP / Delegated** | `text-orange-500` (`AlertTriangle`) | `bg-orange-500/5` | `text-orange-500 dark:text-orange-400`, `text-orange-600 dark:text-orange-400`, `text-orange-700 dark:text-orange-300` | PHP frames, delegated service errors, delegated server traces |

> ⚠ There is **no purple theme** in the current codebase. All delegated/PHP-related UI uses orange.

### 3.2 Overview Tab

| Element | Tailwind Classes | Source Component |
|---------|-----------------|------------------|
| Error banner | `border-destructive/30 bg-destructive/5`, `text-destructive` | `OverviewContent` — BackendMessage |
| Missing delegation warning | `border-warning/30 bg-warning/5`, `text-warning` | `OverviewContent` — missing DelegatedRequestServer fields |
| HTTP status badge (≥400) | `variant="destructive"` | `OverviewContent` |
| HTTP status badge (<400) | `variant="outline"` | `OverviewContent` |
| Session badge | `variant="outline"` | `OverviewContent` |
| Stack traces badge | `variant="outline"` | `OverviewContent` |
| Site URL box | `bg-muted` | `OverviewContent` |

#### Overview Tab Code (BackendMessage banner)

```tsx
// src/components/errors/BackendSection.tsx — OverviewContent
{error.envelopeErrors?.BackendMessage && (
  <div className="rounded-md border border-destructive/30 bg-destructive/5 p-3 space-y-2">
    <h4 className="text-xs font-medium text-destructive uppercase tracking-wider flex items-center gap-1.5">
      <Server className="h-3 w-3" />
      Backend Error
    </h4>
    <p className="text-sm font-mono break-all">{error.envelopeErrors.BackendMessage}</p>
  </div>
)}
```

#### Overview Tab Code (Missing delegation warning)

```tsx
// src/components/errors/BackendSection.tsx — OverviewContent
// v2.1.0: Uses semantic warning token
<div className="rounded-md border border-warning/30 bg-warning/5 p-3 space-y-1">
  <h4 className="text-xs font-medium text-warning uppercase tracking-wider flex items-center gap-1.5">
    <AlertTriangle className="h-3 w-3" />
    Missing Delegation Data
  </h4>
  <p className="text-xs text-muted-foreground">...</p>
</div>
```

### 3.3 Stack Tab

#### PHP Delegated Error Stack (orange)

```tsx
// src/components/errors/BackendSection.tsx — StackContent
<h4 className="... flex items-center gap-2">
  <AlertTriangle className="h-4 w-4 text-orange-500" />
  PHP Delegated Error Stack
</h4>
<ScrollArea className="h-[200px] rounded-md border bg-orange-500/5">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap break-all text-orange-700 dark:text-orange-300">
    ...
  </pre>
</ScrollArea>
```

#### PHP Stack Trace — Delegated Server (orange)

```tsx
// src/components/errors/BackendSection.tsx — StackContent (DelegatedRequestServer.StackTrace)
<AlertTriangle className="h-4 w-4 text-orange-500" />
<Badge variant={statusCode >= 400 ? "destructive" : "secondary"} className="text-xs">
  {Method} {StatusCode}
</Badge>
<ScrollArea className="h-[200px] rounded-md border bg-orange-500/5">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap break-all text-orange-700 dark:text-orange-300">
    ...
  </pre>
</ScrollArea>
```

#### PHP Stack Frames Table (orange)

```tsx
// src/components/errors/BackendSection.tsx — StackContent (phpStackFrames)
<thead className="bg-orange-500/10">...</thead>
<tr className={cn("border-t border-border/50", index === 0 && "bg-orange-500/5")}>
  <td className="p-2 font-mono">
    <span className={cn(index === 0 && "text-orange-600 dark:text-orange-400 font-semibold")}>
      {frame.class}::{frame.function}()
    </span>
  </td>
</tr>
```

#### Go Backend Stack (muted, no color coding)

```tsx
// src/components/errors/BackendSection.tsx — StackContent (envelopeBackendStack)
<ScrollArea className="h-[200px] rounded-md border bg-muted">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap break-all">
    {envelopeBackendStack.join('\n')}
  </pre>
</ScrollArea>
```

#### Session Go Stack Frames (blue text)

```tsx
// src/components/errors/BackendSection.tsx — StackContent (sessionGoFrames)
<ScrollArea className="h-[200px] rounded-md border bg-muted">
  <span className="font-semibold text-blue-500 dark:text-blue-400">
    {frame.class}::{frame.function}
  </span>
</ScrollArea>
```

#### Session PHP Stack Frames (orange text)

```tsx
// src/components/errors/BackendSection.tsx — StackContent (sessionPhpFrames)
<ScrollArea className="h-[200px] rounded-md border bg-orange-500/5">
  <span className="font-semibold text-orange-500 dark:text-orange-400">
    {frame.class}::{frame.function}()
  </span>
</ScrollArea>
```

#### PHP Log (stacktrace.txt) (orange)

```tsx
// src/components/errors/BackendSection.tsx — StackContent
<ScrollArea className="h-[200px] rounded-md border bg-orange-500/5">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap break-all text-orange-700 dark:text-orange-300">
    {sessionDiag.phpStackTraceLog}
  </pre>
</ScrollArea>
```

### 3.4 SessionLogsTab — StackFrameRow Colors

```tsx
// src/components/errors/SessionLogsTab.tsx — StackFrameRow
function StackFrameRow({ index, frame, variant }: { index: number; frame: SessionStackFrame; variant: "golang" | "php" }) {
  const fnColor = variant === "golang" ? "text-blue-500 dark:text-blue-400" : "text-orange-500 dark:text-orange-400";
  return (
    <span className={cn("font-semibold", fnColor)}>
      {frame.class ? `${frame.class}::${frame.function}` : frame.function}
    </span>
  );
}
```

#### SessionLogsTab — PHP Log View

```tsx
<pre className="p-3 text-xs font-mono whitespace-pre-wrap break-words text-orange-600 dark:text-orange-400">
  {phpStackTraceLog}
</pre>
```

### 3.5 Request Tab — 2-Node Request Chain

The request chain has **two nodes** (not three):

```
Node 1: React → Go     → Blue   (bg-blue-500 dot, bg-blue-500/10 badge)
Node 2: Go → PHP       → Orange (bg-orange-500 dot, bg-orange-500/10 badge, text-orange-600)
```

#### Node 1 Code

```tsx
// src/components/errors/RequestDetails.tsx
<div className="w-2 h-2 rounded-full bg-blue-500 shrink-0" />
<Badge variant="outline" className="text-xs font-mono bg-blue-500/10 border-blue-500/30">React → Go</Badge>
```

#### Node 2 Code

```tsx
// src/components/errors/RequestDetails.tsx
<div className="w-2 h-2 rounded-full bg-orange-500 shrink-0" />
<Badge variant="outline" className="text-xs font-mono bg-orange-500/10 border-orange-500/30 text-orange-600 dark:text-orange-400">Go → PHP</Badge>
<p className="text-xs font-mono text-orange-600 dark:text-orange-400 break-all">{phpEndpointUrl}</p>
```

#### PHP Error Stack in Request Tab

```tsx
<pre className="... text-orange-600 dark:text-orange-400 whitespace-pre-wrap">
  {error.envelopeErrors.DelegatedServiceErrorStack.join('\n')}
</pre>
```

### 3.6 Traversal Tab

#### Endpoint Flow — PHP Badge

```tsx
// src/components/errors/TraversalDetails.tsx
<Badge variant="outline" className="shrink-0 text-xs bg-orange-500/10 border-orange-500/30 text-orange-600 dark:text-orange-400">PHP</Badge>
```

#### Delegated Service Error Stack

```tsx
// src/components/errors/TraversalDetails.tsx
<AlertTriangle className="h-4 w-4 text-orange-500" />
<span className="text-orange-600 dark:text-orange-400">Delegated Service Error Stack</span>
<ScrollArea className="h-[200px] rounded-md border border-orange-500/30 bg-orange-500/5">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap break-all">...</pre>
</ScrollArea>
```

#### Methods Stack & Backend Trace

Both use `bg-muted` backgrounds with `bg-primary/5` for the first row and `text-primary font-semibold` for the first method name, matching the Execution tab pattern.

### 3.7 Log Tab (SessionLogsTab — LogLine)

```typescript
// src/components/errors/SessionLogsTab.tsx
// v2.1.0: Uses semantic design tokens — no hardcoded amber/green

const LOG_LINE_STYLES = {
  stage: 'text-primary font-semibold',
  error: 'text-destructive',
  warning: 'text-warning',
  success: 'text-success',
  default: 'text-foreground',
} as const;

function getLogLineStyle(line: string): string {
  if (line.includes("STAGE:") || line.match(/^[─═]+$/)) return LOG_LINE_STYLES.stage;
  if (line.includes("[ERROR]") || line.includes("[FATAL]")) return LOG_LINE_STYLES.error;
  if (line.includes("[WARN]")) return LOG_LINE_STYLES.warning;
  if (line.includes("✓") || line.includes("completed") || line.includes("success")) return LOG_LINE_STYLES.success;

  const stageEndMatch = line.match(/STAGE END: (\w+) - (\w+) \((\d+)ms\)/);
  if (stageEndMatch) {
    const [, , status] = stageEndMatch;
    return status === "success" ? LOG_LINE_STYLES.success : LOG_LINE_STYLES.error;
  }

  return LOG_LINE_STYLES.default;
}

function LogLine({ line }: { line: string }): JSX.Element {
  return <div className={getLogLineStyle(line)}>{line}</div>;
}
```

### 3.8 Execution Tab — BackendLogEntry Level Colors

```tsx
// src/components/errors/BackendSection.tsx — BackendLogEntry
// NOTE: Use LogLevel enum constants — not raw strings.
<div className={cn(
  "text-xs font-mono py-1 px-2 rounded",
  log.level === LogLevel.Error && "bg-destructive/10 text-destructive",
  log.level === LogLevel.Warn  && "bg-warning/10 text-warning",
  log.level === LogLevel.Info  && "bg-primary/10 text-primary",
  log.level === LogLevel.Debug && "bg-muted text-muted-foreground"
)}>
  <span className="text-muted-foreground">[{formatTs(log.timestamp)}]</span>
  {log.step && <span className="text-primary ml-1">[{log.step}]</span>}
  <span className="ml-1 whitespace-pre-wrap break-words">{message}</span>
</div>
```

### 3.9 Execution Tab — Go Call Chain Table

```tsx
// src/components/errors/BackendSection.tsx — ExecutionContent
<thead className="bg-muted">...</thead>
<tr className={cn("border-t border-border/50", index === 0 && "bg-primary/5")}>
  <td className="p-2 font-mono font-semibold">{frame.Method}</td>
  <td className="p-2 font-mono text-muted-foreground truncate max-w-[200px]">{frame.File}</td>
</tr>
```

---

## 4. Frontend Section — Color Themes

### 4.1 Overview Tab

| Element | Classes | Source |
|---------|---------|--------|
| Trigger badge | `bg-primary/5 border-primary/20` | `FrontendSection` |
| Source badge | `variant="secondary"` + `font-mono` | `FrontendSection` |
| First call chain entry | `text-primary font-semibold` | `FrontendSection` — invocationChain |
| Last click in path | `text-primary` | `FrontendSection` — uiClickPath |

### 4.2 Stack Tab

| Element | Classes | Source |
|---------|---------|--------|
| React execution chain area | `bg-blue-500/5` border | `FrontendSection` — executionLogs |
| Activity icon | `text-blue-500` | `FrontendSection` — Activity icon |
| First parsed frame row | `bg-primary/5` highlight | `FrontendSection` — displayFrames |
| First frame function | `text-primary font-semibold` | `FrontendSection` — displayFrames |
| Internal frames | `opacity-50` dimmed | `FrontendSection` — `frame.isInternal` |
| Debug mode tip | `bg-muted text-muted-foreground` | `FrontendSection` — executionLogsEnabled=false |

#### React Execution Chain Code

```tsx
// src/components/errors/FrontendSection.tsx
<Activity className="h-4 w-4 text-blue-500" />
<ScrollArea className="h-32 rounded-md border bg-blue-500/5">
  <pre className="text-xs p-3 font-mono whitespace-pre-wrap">
    {error.executionLogsFormatted}
  </pre>
</ScrollArea>
```

#### Parsed Stack Table Code

```tsx
// src/components/errors/FrontendSection.tsx
<tr className={cn("border-t border-border/50", index === 0 && "bg-primary/5", frame.isInternal && "opacity-50")}>
  <td className="p-2 font-mono">
    <span className={cn(index === 0 && "text-primary font-semibold")}>{frame.function}</span>
  </td>
</tr>
```

### 4.3 Fixes Tab

```tsx
// src/components/errors/FrontendSection.tsx
<span className="flex-shrink-0 w-5 h-5 rounded-full bg-primary/10 text-primary text-xs flex items-center justify-center">
  {index + 1}
</span>
```

---

## 5. Section Toggle Buttons

```tsx
// src/components/errors/GlobalErrorModal.tsx
<Button variant={activeSection === "backend" ? "default" : "outline"} size="sm"
  className="gap-1.5 sm:gap-2 text-xs sm:text-sm flex-1 sm:flex-none">
  <Server className="h-3.5 w-3.5 sm:h-4 sm:w-4" /> Backend
</Button>
<Button variant={activeSection === "frontend" ? "default" : "outline"} size="sm"
  className="gap-1.5 sm:gap-2 text-xs sm:text-sm flex-1 sm:flex-none">
  <Monitor className="h-3.5 w-3.5 sm:h-4 sm:w-4" /> Frontend
</Button>
```

Active section uses `variant="default"` (primary bg), inactive uses `variant="outline"`. Section toggle bar has `bg-muted/30` background.

---

## 6. Error History Drawer Colors

```typescript
// src/components/errors/ErrorHistoryDrawer.tsx — item card states
// Selected:
className="p-3 rounded-lg border cursor-pointer transition-colors bg-accent border-primary"
// Default:
className="p-3 rounded-lg border cursor-pointer transition-colors bg-card hover:bg-accent/50"

// Header icon:
<AlertCircle className="h-5 w-5 text-destructive" />

// Total count badge:
<Badge variant="secondary">{total}</Badge>

// Delete/clear buttons:
className="text-destructive hover:text-destructive"

// Empty state icon:
<AlertCircle className="h-12 w-12 mb-3 opacity-50" />
```

---

## 7. Error Queue Badge

```tsx
// src/components/errors/ErrorQueueBadge.tsx
<Button
  variant="ghost"
  size="sm"
  className="relative h-8 px-2 text-destructive hover:text-destructive hover:bg-destructive/10"
>
  <AlertCircle className="h-4 w-4" />
  <Badge
    variant="destructive"
    className="absolute -top-1 -right-1 h-5 min-w-5 px-1 text-xs"
  >
    {displayCount > 99 ? "99+" : displayCount}
  </Badge>
</Button>
```

---

## 8. App Error Boundary

```tsx
// src/components/errors/AppErrorBoundary.tsx — fallback UI
<div className="flex min-h-[60vh] items-center justify-center px-6 py-12">
  <div className="w-full max-w-lg space-y-4 rounded-lg border bg-background p-6">
    <h1 className="text-lg font-semibold">Something went wrong</h1>
    <p className="text-sm text-muted-foreground">...</p>
    <Button variant="outline">View error details</Button>
    <Button>Reload</Button>  {/* Primary action — default variant */}
  </div>
</div>
```

---

## 9. Color Usage Rules

1. **Never use raw color classes** in error components — always reference design tokens or the documented tier colors
2. **Tier colors are fixed**: Blue = Go session frames (`text-blue-500`), Orange = PHP/Delegated (`text-orange-500`), Neutral (`bg-muted`) = Go raw stacks
3. **Error levels use semantic tokens**: `text-destructive` (error), `text-warning` (warn), `text-primary` (info)
4. **Log lines use semantic tokens**: `text-success` (success), `text-warning` (warn), `text-destructive` (error) — never `text-amber-*` or `text-green-*`
5. **Backgrounds use low opacity**: `bg-destructive/5`, `bg-orange-500/5`, `bg-orange-500/10` (table headers), `bg-blue-500/5` (React execution), `bg-primary/5` (first row highlight)
6. **Borders use medium opacity**: `border-destructive/30`, `border-orange-500/30`
7. **BackendLogEntry uses semantic tokens**: `bg-destructive/10`, `bg-warning/10`, `bg-primary/10`, `bg-muted`
8. **Request chain dots**: `bg-blue-500` (Node 1: React→Go), `bg-orange-500` (Node 2: Go→PHP) — these are intentionally raw colors as they represent infrastructure tiers, not semantic states

---

## 10. Version History

| Version | Date | Changes |
|---------|------|---------|
| v1.0.0 | 2026-03-03 | Initial creation |
| v2.0.0 | 2026-03-03 | Full audit against actual components: removed incorrect purple tier, fixed 3-hop→2-hop request chain, added BackendLogEntry/LogLine/StackFrameRow actual classes, added all SessionLogsTab colors, corrected Go Backend stack bg to `bg-muted` |

---

## Cross-References

- [Error Modal Overview](./00-overview.md) — Modal architecture overview
- [Error Modal Reference](./03-error-modal-reference.md) — Full modal structure, data model, visual layout diagrams
- [Notification Colors](../03-notification-colors.md) — Toast/notification color tokens
- [LogLevel Enum](../../../02-coding-guidelines/02-typescript/10-log-level-enum.md) — Canonical `LogLevel` enum (Debug, Info, Warn, Error, Fatal) used in color mappings

---

*Color Theme Reference v2.0.0 — updated: 2026-03-03*

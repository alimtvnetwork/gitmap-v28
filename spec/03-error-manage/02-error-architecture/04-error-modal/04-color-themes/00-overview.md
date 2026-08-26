# Color Theme & Design Token Reference (Index)

> **Parent:** [Error Modal Spec](../00-overview.md)  
> **Version:** 2.1.0  
> **Updated:** 2026-03-31  
> **Purpose:** Definitive color mapping for every error-related UI element.

---

## File Index

| # | File | Description |
|---|------|-------------|
| 01 | [01-design-tokens.md](./01-design-tokens.md) | CSS custom properties (light/dark) + error level color mapping |
| 02 | [02-backend-tab-colors.md](./02-backend-tab-colors.md) | Backend section tab-specific colors (Overview, Stack, Session, Request, Traversal, Execution) |
| 03 | [03-frontend-and-ui-colors.md](./03-frontend-and-ui-colors.md) | Frontend section color themes + UI element colors (section toggle, error history drawer, queue badge, error boundary) |

---

## Two-Tier Color System

| Tier | Icon Color | Background | Text Color | Used For |
|------|-----------|------------|-----------|----------|
| **Go Backend** | — (Server icon) | `bg-muted` | `text-blue-500 dark:text-blue-400` (session frames) | Go stack traces, methods stack |
| **PHP / Delegated** | `text-orange-500` (AlertTriangle) | `bg-orange-500/5` | `text-orange-500/600/700` | PHP frames, delegated service errors |

> ⚠ There is **no purple theme** in the current codebase. All delegated/PHP-related UI uses orange.

---

---

- [Error Modal Reference](../03-error-modal-reference/00-overview.md)
- [LogLevel Enum](../../../../02-coding-guidelines/02-typescript/10-log-level-enum.md)

---

*Color theme index — updated: 2026-03-31*

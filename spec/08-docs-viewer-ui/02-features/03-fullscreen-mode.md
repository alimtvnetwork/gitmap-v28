# Feature: Fullscreen Mode

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Specification

### Trigger

- Button in the doc header bar (Maximize2 icon from lucide-react)
- Keyboard shortcut: `F` key (when no input focused)

### Behavior

- Fullscreen hides the sidebar and expands the content area to fill the viewport
- Header remains visible with a minimize button to exit
- `Escape` key also exits fullscreen
- Implemented via React state (not browser Fullscreen API) for better control

---

*Fullscreen mode — updated: 2026-04-03*

# Feature: Copy Markdown Button

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Specification

### Trigger

- Button in the doc header bar (Copy icon from lucide-react)
- Copies the raw markdown content of the active file to clipboard

### Feedback

- Button text/icon changes to "Copied!" with a checkmark for 2 seconds
- Uses `navigator.clipboard.writeText()`

---

*Copy markdown — updated: 2026-04-03*

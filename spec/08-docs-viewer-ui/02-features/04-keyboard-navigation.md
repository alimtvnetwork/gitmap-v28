# Feature: Keyboard Navigation

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Specification

### Key Bindings

| Key | Action |
|-----|--------|
| `→` | Next file within current folder |
| `←` | Previous file within current folder |
| `↓` | Next folder (jump to first file) |
| `↑` | Previous folder (jump to first file) |

### Guards

- Only active when no `<input>`, `<textarea>`, or `[contenteditable]` is focused
- Only active when a file is currently selected (no-op on welcome screen)

### Edge Cases

- At last file in folder → wraps to first file in same folder
- At last folder → wraps to first folder

---

*Keyboard navigation — updated: 2026-04-03*

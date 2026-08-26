# Clipboard Utility — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## File: `src/lib/clipboard.ts`

### Strategy: Fallback-First

The clipboard utility prioritizes the **synchronous `document.execCommand('copy')`** method because it works reliably in sandboxed iframes and non-secure contexts. The modern `navigator.clipboard.writeText` is used only as a secondary fallback.

### Flow

```
copyTextToClipboard(text)
  │
  ├── Empty text? → toast.error("Nothing to copy") → throw
  │
  ├── fallbackCopyText(text)
  │     ├── Save current active element + selection
  │     ├── Create hidden textarea (fixed, off-screen, opacity 0)
  │     ├── Append to body, focus, select all
  │     ├── document.execCommand("copy")
  │     ├── Remove textarea
  │     ├── Restore previous selection + focus
  │     └── Return boolean success
  │
  ├── Fallback succeeded? → toast.success("Copied to clipboard") → return
  │
  ├── navigator.clipboard unavailable? → toast.error("Copy failed") → throw
  │
  └── navigator.clipboard.writeText(text)
        ├── Success → toast.success("Copied to clipboard")
        └── Failure → toast.error("Copy failed") → throw
```

### Toast Feedback

All copy operations show a **sonner toast** notification:

| State | Toast |
|-------|-------|
| Success | `toast.success("Copied to clipboard")` |
| Empty content | `toast.error("Nothing to copy", { description: "No content available." })` |
| Clipboard blocked | `toast.error("Copy failed", { description: "Clipboard access is blocked..." })` |
| Write failed | `toast.error("Copy failed", { description: "Unable to access the clipboard..." })` |

### Hidden Textarea Implementation

```typescript
function createHiddenTextarea(text: string): HTMLTextAreaElement {
  const textarea = document.createElement("textarea");
  textarea.value = text;
  textarea.setAttribute("readonly", "true");
  textarea.style.position = "fixed";
  textarea.style.top = "0";
  textarea.style.left = "-9999px";
  textarea.style.opacity = "0";
  textarea.style.pointerEvents = "none";
  return textarea;
}
```

Key properties:
- `readonly` prevents mobile keyboards
- `position: fixed` + `left: -9999px` keeps it invisible
- `opacity: 0` + `pointerEvents: none` ensures no visual/interaction impact
- Previous selection is saved and restored after copy

---

*Clipboard Utility — updated: 2026-04-08*

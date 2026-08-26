# Interactions — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All interactions use **event delegation**. A single set of listeners is attached to the markdown container element. Each handler uses `.closest(selector)` to check if the click target is relevant.

Source: `useCodeBlockEvents.ts` → creates handlers from 3 files:
- `codeBlockActionHandlers.ts` — copy, download, fullscreen, checklist
- `codeBlockLineHandlers.ts` — line click, pin, range, keyboard
- `codeBlockDragHandlers.ts` — drag-select, hover, font size

---

## 1. Copy Code

**Trigger:** Click `.copy-code-btn`

**Behavior:**
1. Read `data-code` attribute (HTML-escaped raw code)
2. Decode escaped characters (`&#10;` → `\n`, `&amp;` → `&`, etc.)
3. Copy to clipboard via `copyTextToClipboard()` (see Clipboard spec)
4. Show feedback:
   - Hide copy icon SVG, show check icon SVG
   - Change label text to "Copied!"
   - Add `.copied` class (green background)
   - After **2000ms**, revert all changes

**Decoding map:**
| Encoded | Decoded |
|---------|---------|
| `&#10;` | `\n` |
| `&#39;` | `'` |
| `&amp;` | `&` |
| `&lt;` | `<` |
| `&gt;` | `>` |
| `&quot;` | `"` |

---

## 2. Download Code

**Trigger:** Click `.download-code-btn`

**Behavior:**
1. Read `data-code` (escaped raw code) and `data-ext` (file extension)
2. Decode escaped characters
3. Create `Blob` with `type: "text/plain"`
4. Create temporary `<a>` element with `download="code.{ext}"`
5. Trigger click, then revoke object URL

**Download filename:** Always `code.{ext}` where `ext` comes from `LANG_EXTENSIONS` map.

---

## 3. Fullscreen Toggle

**Trigger:** Click `.fullscreen-code-btn`

**Behavior:**
1. Read `data-block-id` from button
2. Toggle state: if current fullscreen block matches ID → exit, otherwise enter
3. `useSyncFullscreenClass` adds/removes `.code-fullscreen` class on the wrapper
4. An overlay `<div class="code-fullscreen-overlay">` is rendered behind the block

**Exit methods:**
- Click fullscreen button again
- Click the overlay
- Press **Escape** key (`useEscapeFullscreen` hook)

**Fullscreen layout:**
- Block becomes `position: fixed; inset: 2rem; z-index: 999`
- Overlay is `position: fixed; inset: 0; z-index: 998; background: hsl(0 0% 0% / 0.7); backdrop-filter: blur(4px)`
- Body has `flex-direction: column; max-height: calc(100vh - 4rem)`
- Code body gets `flex: 1; overflow: auto`

---

## 4. Font Size Controls

**Trigger:** Click `.font-increase-btn`, `.font-decrease-btn`, or `.font-reset-btn`

**Behavior:**
1. Read `data-block-id` to find the target wrapper
2. Read current `--code-font-size` from computed style
3. Compute new size:
   - Increase: `current + 2` (max **32px**)
   - Decrease: `current - 2` (min **12px**)
   - Reset: **18px**
4. Set `--code-font-size` on the wrapper element

**Constants:**
| Name | Value |
|------|-------|
| `DEFAULT_FONT_SIZE` | 18px |
| `MIN_FONT_SIZE` | 12px |
| `MAX_FONT_SIZE` | 32px |
| `FONT_SIZE_STEP` | 2px |

---

## 5. Line Selection — Click

**Trigger:** Click on a `.code-line` or `.code-line-number` (not on tool buttons)

**Single click:**
1. Clear all `.line-pinned` classes in the wrapper
2. Add `.line-pinned` to the clicked line AND its corresponding line number
3. Update selection bar with "Line {N}"
4. Store as anchor for future Shift+click

**Shift+click:**
1. Use previously stored anchor line
2. Pin all lines from anchor to clicked line (inclusive)
3. Update selection bar with "Lines {from}–{to}"

---

## 6. Line Selection — Keyboard

**Trigger:** Arrow Up / Arrow Down when a wrapper is active

**Without Shift:**
- Move pin to adjacent line (single line selected)
- Scroll line into view

**With Shift:**
- Extend/contract the range from anchor
- Update selection bar

**Guards:**
- Ignores input if `document.activeElement` is `INPUT`, `TEXTAREA`, or `contentEditable`
- Prevents default scroll behavior

---

## 7. Line Selection — Drag

**Trigger:** Mousedown on `.code-line-number`

**Behavior:**
1. `mousedown` on a line number starts drag: record anchor, clear existing pins, pin anchor
2. `mousemove` (on document): resolve line under cursor, pin range from anchor to current
3. `mouseup` (on document): end drag state

**Resolution:** Uses `document.elementFromPoint(clientX, clientY)` to find the element under the cursor during drag.

---

## 8. Line Hover

**Trigger:** `mouseover` on any element inside `.code-block-wrapper`

**Behavior:**
1. Resolve line index from target
2. Add `.line-highlight` class to that line and its line number
3. On `mouseout`, remove all `.line-highlight` classes in the wrapper

---

## 9. Copy Selected Lines

**Trigger:** Click `.copy-selected-btn` in the selection bar

**Behavior:**
1. Find all `.code-line.line-pinned` in the wrapper
2. Extract `textContent` from each
3. Join with `\n`
4. Copy via clipboard utility
5. Show "Copied!" feedback (same pattern as main copy button)

---

## 10. Clear Selection

**Trigger:** Click `.clear-selected-btn` (✕) in the selection bar

**Behavior:**
1. Remove `.line-pinned` from all elements in the wrapper
2. Hide the selection bar (`display: none`)
3. Hide the header selection label
4. Clear the wrapper from `lastPinned` map

---

## Selection Bar

The selection bar appears at the bottom of the code block when lines are pinned:

- Shows "Line {N}" or "Lines {from}–{to}"
- Contains "Copy selected" button and "✕" clear button
- Animated in with `slideUpBar` keyframe (opacity 0→1, translateY 4px→0)
- The header also shows a selection label badge next to the line count

---

*Interactions — updated: 2026-04-08*

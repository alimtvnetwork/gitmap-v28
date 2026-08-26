# HTML Structure — Code Block System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Code Block HTML

Every fenced code block produces this exact HTML structure:

```html
<div class="code-block-wrapper my-5"
     style="--lang-accent: {HSL_VALUE}"
     data-block-id="{ID}">

  <!-- HEADER -->
  <div class="code-block-header">
    <!-- Left: Language badge -->
    <div class="code-lang-badge" style="--badge-color: {HSL_VALUE}">
      <span class="code-lang-dot"></span>
      <span>{LANGUAGE_LABEL}</span>
    </div>

    <!-- Right: Tools -->
    <div class="code-header-right">
      <span class="code-line-count">{N} line(s)</span>
      <span class="code-selection-label" style="display:none"></span>

      <!-- Font controls group -->
      <div class="code-font-controls">
        <button class="code-tool-btn font-decrease-btn" data-block-id="{ID}">A-</button>
        <button class="code-tool-btn font-reset-btn" data-block-id="{ID}">A</button>
        <button class="code-tool-btn font-increase-btn" data-block-id="{ID}">A+</button>
      </div>

      <!-- Copy button -->
      <button class="code-tool-btn copy-code-btn" data-code="{ESCAPED_CODE}">
        <svg class="copy-icon"><!-- clipboard icon --></svg>
        <svg class="check-icon" style="display:none"><!-- checkmark icon --></svg>
        <span class="copy-label">Copy</span>
      </button>

      <!-- Download button -->
      <button class="code-tool-btn download-code-btn"
              data-code="{ESCAPED_CODE}"
              data-ext="{FILE_EXTENSION}"
              data-lang="{LABEL}">
        <svg><!-- download icon --></svg>
        <span>Download</span>
      </button>

      <!-- Fullscreen button -->
      <button class="code-tool-btn fullscreen-code-btn" data-block-id="{ID}">
        <svg><!-- expand icon --></svg>
      </button>
    </div>
  </div>

  <!-- BODY -->
  <div class="code-block-body">
    <!-- Line numbers column -->
    <pre class="code-line-numbers" aria-hidden="true">
      <span class="code-line-number">1</span>
      <span class="code-line-number">2</span>
      ...
    </pre>

    <!-- Code content column -->
    <pre class="code-content">
      <code class="hljs language-{LANG}">
        <span class="code-line" data-line="1">{HIGHLIGHTED_HTML}</span>
        <span class="code-line" data-line="2">{HIGHLIGHTED_HTML}</span>
        ...
      </code>
    </pre>
  </div>

  <!-- SELECTION BAR (hidden by default) -->
  <div class="copy-selected-bar" style="display:none">
    <span class="copy-selected-label"></span>
    <button class="code-tool-btn copy-selected-btn">
      <svg class="copy-icon"><!-- clipboard --></svg>
      <svg class="check-icon" style="display:none"><!-- checkmark --></svg>
      <span class="copy-label">Copy selected</span>
    </button>
    <button class="code-tool-btn clear-selected-btn">✕</button>
  </div>
</div>
```

---

## Data Attributes Reference

| Attribute | Element | Purpose |
|-----------|---------|---------|
| `data-block-id` | `.code-block-wrapper`, font buttons, fullscreen button | Unique block identifier (integer) |
| `data-code` | Copy button, Download button | HTML-escaped raw code content |
| `data-ext` | Download button | File extension for download (e.g., `ts`, `go`) |
| `data-lang` | Download button | Display label (e.g., `TypeScript`) |
| `data-line` | `.code-line` | 1-indexed line number |

---

## CSS Custom Properties on Wrapper

| Property | Source | Purpose |
|----------|--------|---------|
| `--lang-accent` | Per-language HSL from `LANG_COLORS` | Hover shadow, badge color, glow |
| `--badge-color` | Same as `--lang-accent` | Badge dot + text color |
| `--code-font-size` | Default `18px`, modified by font controls | Font size for code + line numbers |
| `--code-line-height` | Default `1.6` | Line height for code + line numbers |

---

## SVG Icons Used

All icons are inline SVGs at 14×14, `viewBox="0 0 24 24"`, stroke-based:

| Icon | Used In | SVG Description |
|------|---------|-----------------|
| Copy | Copy button, Copy selected | Two overlapping rounded rectangles |
| Check | Copy feedback | Polyline checkmark (points: 20 6 9 17 4 12) |
| Download | Download button | Arrow pointing down into a tray |
| Fullscreen | Fullscreen button | Two diagonal arrows (expand) |

---

## Inline Code HTML

Single backtick produces:

```html
<code class="inline-code">{ESCAPED_CONTENT}</code>
```

---

## Fullscreen Overlay

When a block is fullscreen, an overlay `div` is injected after the prose container:

```html
<div class="code-fullscreen-overlay" onClick="exitFullscreen"></div>
```

The target block gets the `.code-fullscreen` class added.

---

*HTML Structure — updated: 2026-04-08*

# WordPress Migration Compatibility

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Ambiguity:** Medium — WordPress implementation approach not yet decided

---

## Overview

The design system is structured to support future WordPress migration. This document explains how the variable-driven architecture maps to WordPress theming, and what to preserve or adapt during migration.

---

## Migration Paths

| Path | Description | Compatibility |
|------|-------------|--------------|
| **Block Theme (FSE)** | Uses `theme.json` for design tokens, block templates for layout | High — CSS variables map directly to `theme.json` settings |
| **Classic Theme** | Uses `style.css` + PHP templates | High — CSS variables work in any CSS file |
| **Plugin-Driven** | React components rendered via shortcodes or blocks | Medium — requires build tooling |

**Current recommendation:** Block Theme (FSE) for maximum token compatibility.

---

## Token Mapping to theme.json

CSS custom properties map to WordPress `theme.json` settings:

```json
{
  "settings": {
    "color": {
      "palette": [
        { "slug": "primary", "color": "hsl(252 85% 60%)", "name": "Primary" },
        { "slug": "accent", "color": "hsl(330 85% 60%)", "name": "Accent" },
        { "slug": "background", "color": "hsl(0 0% 100%)", "name": "Background" },
        { "slug": "foreground", "color": "hsl(230 25% 15%)", "name": "Foreground" },
        { "slug": "muted", "color": "hsl(220 20% 96%)", "name": "Muted" }
      ]
    },
    "spacing": {
      "units": ["rem"],
      "spacingSizes": [
        { "slug": "sm", "size": "0.5rem", "name": "Small" },
        { "slug": "md", "size": "1rem", "name": "Medium" },
        { "slug": "lg", "size": "2rem", "name": "Large" }
      ]
    },
    "typography": {
      "fontFamilies": [
        { "slug": "heading", "fontFamily": "'Ubuntu', sans-serif", "name": "Heading" },
        { "slug": "body", "fontFamily": "'Poppins', sans-serif", "name": "Body" }
      ]
    },
    "custom": {
      "radius": "0.5rem",
      "headingGradientFrom": "252 85% 60%",
      "headingGradientTo": "330 85% 60%"
    }
  }
}
```

---

## What Transfers Directly

| Design System Element | WordPress Equivalent |
|----------------------|---------------------|
| CSS custom properties | `theme.json` palette + custom settings |
| Font stacks | `theme.json` fontFamilies |
| Spacing scale | `theme.json` spacingSizes |
| Border radius | `theme.json` custom.radius |
| Component CSS classes | `style.css` with same class names |
| CSS3 transitions | `style.css` — works unchanged |
| Dark mode via `.dark` class | Body class toggle via JavaScript or plugin |

---

## What Needs Adaptation

| Design System Element | WordPress Adaptation |
|----------------------|---------------------|
| Tailwind utility classes | Replace with WordPress block styles or custom CSS |
| React components | Convert to block patterns or shortcodes |
| `hsl(var(--token))` syntax | Works in CSS but not in `theme.json` — use direct values |
| Sidebar component | WordPress native sidebar widget area |
| Monaco editor | Not applicable — WordPress has its own editor |

---

## Admin Configurable Theming

Theme variables should be exposable through WordPress admin:

1. **Customizer panels** — Color pickers that write to CSS variables
2. **theme.json presets** — Multiple preset themes selectable from admin
3. **Custom CSS field** — Advanced users can override variables directly

The variable architecture makes this possible because:
- All visual changes flow from ~40 CSS custom properties
- No component-level overrides needed
- WordPress Customizer can write `<style>` blocks with variable overrides

---

## Ambiguities (Not Yet Decided)

| Question | Notes |
|----------|-------|
| Block theme vs. classic theme? | Block recommended but not confirmed |
| Plugin or theme? | Design system components could be either |
| Multi-theme presets from day one? | Currently single base theme; presets can be added later |
| Gutenberg block mapping? | Which React components become blocks vs. patterns? |
| Dark mode toggle mechanism? | Body class toggle, cookie, or WP user preference? |

These should be resolved before implementation begins.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Theme Variable Architecture | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |
| Design Principles (Portability) | [01-design-principles.md](./01-design-principles.md) |

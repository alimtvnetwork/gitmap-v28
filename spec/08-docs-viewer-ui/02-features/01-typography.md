# Feature: Typography System

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Specification

### Google Fonts Integration

Add to `index.html` `<head>`:

```html
<link href="https://fonts.googleapis.com/css2?family=Ubuntu:wght@400;500;600;700&family=Poppins:wght@300;400;500;600&display=swap" rel="stylesheet">
```

### CSS Variables (index.css)

```css
--font-heading: 'Ubuntu', sans-serif;
--font-body: 'Poppins', sans-serif;
```

### Tailwind Config

```ts
fontFamily: {
  heading: ['Ubuntu', 'sans-serif'],
  body: ['Poppins', 'sans-serif'],
}
```

### Application Rules

| Selector | Font | Class |
|----------|------|-------|
| `h1, h2, h3, h4` | Ubuntu | `font-heading` |
| `body, p, li, td, span` | Poppins | `font-body` |
| `code, pre` | System monospace | (unchanged) |

---

*Typography — updated: 2026-04-03*

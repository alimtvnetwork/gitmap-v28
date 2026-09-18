# CSS Modularity and Styling

## 1. Utility-First Preference

- Prefer utility-first CSS frameworks (like Tailwind CSS) for rapid development and constrained design token usage.
- Avoid writing custom, bespoke CSS classes unless strictly necessary for complex animations or components that utility classes cannot cleanly handle.

## 2. Scoped CSS

- If custom CSS must be written, it MUST be scoped to the component (e.g., CSS Modules `Button.module.css` or styled-components).
- Global CSS styles are strictly forbidden, except for CSS resets (e.g., normalize.css) and root design token variable definitions.

## 3. Design Tokens

- Hardcoding colors, fonts, or spacing values in components is a CODE-RED violation.
- Always reference the centralized Design System tokens (e.g., `var(--color-primary-500)` or Tailwind's `text-primary-500`).

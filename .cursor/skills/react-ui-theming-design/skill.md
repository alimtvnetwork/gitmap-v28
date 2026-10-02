---
name: react-ui-theming-design
description: Develop, theme, and maintain modern React UI components using Tailwind CSS, Radix UI primitives, and design tokens.
---

# React UI, Components & Theming Design System

This skill guides frontend component engineering, theme token management, and accessibility standards across `src/`.

## Architecture & Technology Stack

- **Framework:** React 18 + Vite + TypeScript.
- **Styling:** Tailwind CSS with CSS variables for light/dark theme transitions.
- **Component Primitives:** Radix UI (`@radix-ui/react-*`), Sonner, Lucide React icons.
- **Data & State:** TanStack React Query (`@tanstack/react-query`) for server state management.

## Guidelines & Best Practices

1. **Component Boundedness:**
   - Keep individual component files under 100 lines.
   - Decompose complex views into focused subcomponents under `src/components/`.

2. **Theming & Color Tokens:**
   - Always reference theme variables (e.g. `bg-background`, `text-foreground`, `border-border`, `text-muted-foreground`) rather than hardcoded hex codes.
   - Support dark mode seamlessly via `ThemeProvider`.

3. **Accessibility & Usability:**
   - Ensure interactive elements are keyboard navigable with visible focus rings.
   - Utilize proper ARIA labels and semantic HTML tags.

4. **Screenshot & Print Screen Image Ingestion:**
   - If UI implementation tasks provide a screenshot URL, print screen link, or base64 data URI:
   - Immediately decode and persist the image into `assets/screenshots/<feature>-<NN>.png` or `assets/ui/<feature>-<NN>.png`.
   - Never keep raw base64 or remote URLs in code comments, specs, or component documentation.
   - Reference the image file using strict relative paths (e.g. `assets/screenshots/<feature>-01.png`) and inspect it as the visual ground truth for padding, layout, responsive breakpoints, and color palette tokens.

5. **Linting & Type Verification:**
   - Run type check and ESLint:
     ```bash
     npm run build
     ```

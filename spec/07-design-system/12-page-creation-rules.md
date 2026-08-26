# Page Creation Rules

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

When building new pages, follow these rules to maintain visual consistency with the design system. Every page inherits the same tokens, motion, and component patterns — no page-specific overrides.

---

## Rule 1: No Page-Specific Colors

Every color on every page MUST come from the theme token registry. No page may introduce:
- Hardcoded hex/rgb values
- Tailwind literal color classes (`bg-purple-600`)
- Inline style color values

If a new semantic color is needed, add it to the token registry first.

---

## Rule 2: Page Structure Template

Every page should follow this vertical flow:

```
1. Header (sticky, consistent across all pages)
2. Hero / Introduction Section
3. Content Sections (from section patterns)
4. CTA / Action Section (optional)
5. Footer (if applicable)
```

---

## Rule 3: Consistent Heading Hierarchy

| Level | Usage |
|-------|-------|
| H1 | Page title — one per page, gradient |
| H2 | Major section headings — gradient, border-bottom |
| H3 | Subsection — left border accent |
| H4 | Minor labels — muted, hover brightens |

Never skip heading levels (no H1 → H3).

---

## Rule 4: Motion Budget

Each page should use no more than:
- 3 unique transition durations
- 2 unique easing functions
- 1 keyframe animation (beyond shared ones)

This prevents visual noise. Stick to the shared transitions from [06-motion-transitions.md](./06-motion-transitions.md).

---

## Rule 5: Responsive First

Design for mobile (< 768px) first, then enhance:
- Single column on mobile
- Grid layouts on tablet/desktop
- Sidebar collapses to sheet on mobile
- Touch targets minimum 44px
- No hover-only features (hover enhances, never gates)

---

## Rule 6: Section Composition

Build pages by composing section patterns from [11-section-patterns.md](./11-section-patterns.md):

```
Page = Header + [Hero] + N × [Section Pattern] + [CTA] + [Footer]
```

Each section is visually self-contained and can be reordered.

---

## Rule 7: Spacing Rhythm

Maintain vertical rhythm between sections:
- Section-to-section gap: `4rem` to `6rem`
- Heading-to-content gap: `1rem` to `1.5rem`
- Content-to-content gap: `0.5rem` to `1rem`
- Use the spacing scale — never arbitrary pixel values

---

## Rule 8: Dark Mode Compatibility

Every page must work in both light and dark modes without modification. This is automatic if all colors use tokens. Verify by toggling themes and checking:
- Text contrast (WCAG AA: 4.5:1 minimum)
- Border visibility
- Hover state visibility
- Image/icon visibility

---

## Rule 9: Font Consistency

- Headings: Ubuntu
- Body: Poppins
- Code: Ubuntu Mono / JetBrains Mono
- No additional fonts without updating the design system

---

## Rule 10: Interaction Consistency

All interactive elements on new pages must follow the state language from [01-design-principles.md](./01-design-principles.md):

| State | Behavior |
|-------|---------|
| Hover | Lighten + subtle lift + optional glow |
| Active | Slightly darker than hover |
| Focus | Ring outline with `--ring` |
| Disabled | 50% opacity |

---

## Checklist for New Pages

- [ ] All colors from theme tokens
- [ ] Heading hierarchy H1 → H2 → H3 (no skips)
- [ ] Sections use patterns from section-patterns spec
- [ ] Transitions use shared durations and easings
- [ ] Works in both light and dark mode
- [ ] Responsive at mobile/tablet/desktop breakpoints
- [ ] Hover effects follow state language
- [ ] Fonts match design system stacks
- [ ] Spacing uses the defined scale

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Section Patterns | [11-section-patterns.md](./11-section-patterns.md) |
| Motion System | [06-motion-transitions.md](./06-motion-transitions.md) |
| Design Principles | [01-design-principles.md](./01-design-principles.md) |
| Theme Variables | [02-theme-variable-architecture.md](./02-theme-variable-architecture.md) |

export interface ReferenceRepoItem {
  name: string;
  repo: string;
  url: string;
  category: string;
  description: string;
  badge: string;
}

export const referenceRepositories: ReferenceRepoItem[] = [
  {
    name: "shadcn/ui",
    repo: "shadcn-ui/ui",
    url: "https://github.com/shadcn-ui/ui",
    category: "Design System & Accessible Primitives",
    description: "Copy-paste accessible React components built on Radix UI and Tailwind CSS with first-class HSL token architecture.",
    badge: "Radix + Tailwind",
  },
  {
    name: "Radix Primitives",
    repo: "radix-ui/primitives",
    url: "https://github.com/radix-ui/primitives",
    category: "Headless Component Architecture",
    description: "Unstyled, fully accessible UI components (dialogs, tooltips, dropdowns, popovers) with WAI-ARIA compliance.",
    badge: "WAI-ARIA Standard",
  },
  {
    name: "Aceternity UI",
    repo: "mannupaaji/aceternity-ui",
    url: "https://github.com/mannupaaji/aceternity-ui",
    category: "Visual Craft & Canvas Animations",
    description: "High-craft modern interactive components, border beams, glowing cards, and canvas visual effects.",
    badge: "Framer Motion",
  },
  {
    name: "Magic UI",
    repo: "magicuidesign/magicui",
    url: "https://github.com/magicuidesign/magicui",
    category: "Micro-Interactions & Landing UI",
    description: "20+ copy-paste components for high-conversion landing pages: animated beams, retro grids, and particle UI.",
    badge: "Creative Engineering",
  },
  {
    name: "Tremor",
    repo: "tremorlabs/tremor",
    url: "https://github.com/tremorlabs/tremor",
    category: "Data Telemetry & Dashboards",
    description: "Modular dashboard and telemetry components with clean dark mode neutral depth metrics cards and charts.",
    badge: "Analytics UI",
  },
  {
    name: "Lucide Icons",
    repo: "lucide-icons/lucide",
    url: "https://github.com/lucide-icons/lucide",
    category: "Vector Iconography",
    description: "Consistent, clean vector iconography system with uniform 24px viewports and zero runtime overhead.",
    badge: "Vector Icons",
  },
  {
    name: "Modern CSS Solutions",
    repo: "moderncss.dev",
    url: "https://moderncss.dev",
    category: "CSS3 Architecture & Layout Patterns",
    description: "Deep-dive CSS3 architectural recipes for modern grid, flexbox, fluid typography, and accessible form styling.",
    badge: "CSS3 Standard",
  },
  {
    name: "CodyHouse Framework",
    repo: "CodyHouse/codyhouse-framework",
    url: "https://github.com/CodyHouse/codyhouse-framework",
    category: "Token Architecture & Typography Scales",
    description: "Production-ready design system tokens, typography scales, spacing units, and modular CSS utility classes.",
    badge: "Tokens & Scales",
  },
  {
    name: "CSS-Tricks Guides",
    repo: "css-tricks.com",
    url: "https://css-tricks.com",
    category: "Subgrid & Custom Properties",
    description: "Authoritative guides on fluid typography, CSS custom properties, responsive subgrid layouts, and micro-animations.",
    badge: "Layout Guide",
  },
];

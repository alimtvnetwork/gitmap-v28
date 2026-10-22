---
name: next-100-ui-improvements
description: Autonomously scan, discover, plan, and refactor 100 UI-level improvements across color tokens, WCAG contrast, 4-plane depth hierarchy, 60/30/10 balance, typography stacks, button systems, sticky navigation physics, spacing rhythm, CSS3 motion, and React component hygiene.
---

# Next 100 UI Improvements

> **[/goal](slashCommand;goal)** Autonomously audit, catalog, plan, and execute 100 concrete UI improvements across the target frontend codebase, enforcing strict design token architecture, WCAG 2.2 AA/AAA contrast, 4-plane depth hierarchy, 60/30/10 visual balance, typography stacks, button systems, sticky navigation physics, and React component decomposition.
> **[/learn](slashCommand;learn)** Ingest the 10 UI improvement clusters and design tokens from `02-spec/07-design-system/` before applying modifications. Enqueue all 100 improvements as structured pending plans in `.ai-memory/plans/pending/` and execute in bounded micro-batches via GitMap.
> **[/plan](slashCommand;plan)** [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.

**Version:** 1.0.0
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. When to Use

Activate this skill when:
- Auditing or refactoring frontend interfaces for enterprise visual polish, consistency, and accessibility.
- Eliminating hardcoded hex codes, inline styling, or ad-hoc Tailwind color classes in favor of semantic HSL tokens.
- Resolving WCAG 2.2 AA and AAA contrast failures across body text, display headings, inputs, and interactive badges.
- Implementing the 4-plane depth hierarchy (Plane 0 Canvas, Plane 1 Raised Wells, Plane 2 Surface Cards, Plane 3 Elevated Modals) with hairline lighting.
- Enforcing the 60/30/10 visual balance rule and the single-accent Von Restorff principle on hero sections and dashboards.
- Upgrading typography to pure live DOM stacks (Ubuntu headings, Poppins body, JetBrains Mono inline, Ubuntu Mono blocks).
- Standardizing the 6 canonical button variants (`primary`, `solid`, `outline`, `glass`, `ghost`, `link`) and 4 sizes (`sm: 36px`, `md: 44px`, `lg: 52px`, `icon: 44px`).
- Calibrating sticky 72px glassmorphic navigation headers with 12px scroll triggers, safe-region geometry (`pad = 14px`), and 220ms close debounce.
- Deconstructing oversized React components (> 80–100 lines) into focused, single-responsibility micro-components.
- Generating and enqueueing continuous 100-improvement iterative cycles into `.ai-memory/plans/pending/`.

---

## 2. Core Architecture & Design Token System

All visual properties must be tokenized and driven by CSS custom properties in `:root` and `.dark` blocks in space-separated HSL format (`H S% L%`):

```css
/* Token declaration */
--token-name: H S% L%;

/* Component usage */
color: hsl(var(--token-name));
background: hsl(var(--token-name) / 0.15);
border: 1px solid hsl(var(--border));
```

### 2.1 Primitive Color Ramps (3-Format Reference)

#### Blue Ramp (Primary Action & Corporate Trust)
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--blue-50` | `oklch(0.972 0.014 259)` | `#F0F6FF` | `rgb(240, 246, 255)` | `hsl(216, 100%, 97%)` | Hero wash start, brand pill bg |
| `--blue-100` | `oklch(0.935 0.032 262)` | `#DEEAFF` | `rgb(222, 234, 255)` | `hsl(218, 100%, 94%)` | Badge hover fill, soft tile ring |
| `--blue-200` | `oklch(0.874 0.061 263)` | `#C1D6FF` | `rgb(193, 214, 255)` | `hsl(220, 100%, 88%)` | Chrome gradient stop, subtle border |
| `--blue-300` | `oklch(0.777 0.104 264)` | `#95B6FA` | `rgb(149, 182, 250)` | `hsl(220, 91%, 78%)` | Chrome gradient end, data series |
| `--blue-400` | `oklch(0.662 0.164 264)` | `#5E8EF6` | `rgb(94, 142, 246)` | `hsl(221, 89%, 67%)` | Gradient text end stop, ring focus |
| `--blue-500` | `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | **Primary CTA Action (`--brand-secondary`)** |
| `--blue-600` | `oklch(0.489 0.208 263.5)` | `#1B51D3` | `rgb(27, 81, 211)` | `hsl(222, 77%, 47%)` | Brand gradient start, active CTA |
| `--blue-700` | `oklch(0.409 0.176 264.5)` | `#173DA7` | `rgb(23, 61, 167)` | `hsl(224, 76%, 37%)` | Brand pill text, high-contrast link |
| `--blue-800` | `oklch(0.369 0.147 264.8)` | `#17368B` | `rgb(23, 54, 139)` | `hsl(224, 72%, 32%)` | Deep interactive border, badge text |
| `--blue-900` | `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | **Core Brand Anchor (`--brand-primary`)** |

#### Violet Ramp (Tertiary Accent & Gradient Partner)
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--violet-50` | `oklch(0.972 0.017 300)` | `#F7F3FF` | `rgb(247, 243, 255)` | `hsl(260, 100%, 98%)` | Accent surface fill |
| `--violet-100` | `oklch(0.946 0.032 300)` | `#F0E9FF` | `rgb(240, 233, 255)` | `hsl(259, 100%, 96%)` | Subtle chip background |
| `--violet-200` | `oklch(0.885 0.068 299)` | `#DFD0FF` | `rgb(223, 208, 255)` | `hsl(259, 100%, 91%)` | Accent border on light cards |
| `--violet-300` | `oklch(0.784 0.128 298)` | `#C3A6FF` | `rgb(195, 166, 255)` | `hsl(260, 100%, 83%)` | Soft glow highlight on dark |
| `--violet-400` | `oklch(0.659 0.187 297.5)` | `#A173F4` | `rgb(161, 115, 244)` | `hsl(261, 85%, 70%)` | Secondary chart accent |
| `--violet-500` | `oklch(0.532 0.253 296.9)` | `#822EE8` | `rgb(130, 46, 232)` | `hsl(267, 80%, 55%)` | **Tertiary Brand Accent (`--brand-tertiary`)** |
| `--violet-600` | `oklch(0.481 0.243 296.5)` | `#721CD2` | `rgb(114, 28, 210)` | `hsl(268, 76%, 47%)` | Brand gradient midpoint (55%) |
| `--violet-700` | `oklch(0.417 0.213 295.9)` | `#5C14AF` | `rgb(92, 20, 175)` | `hsl(268, 79%, 38%)` | Accent foreground text |
| `--violet-800` | `oklch(0.362 0.187 295.5)` | `#4B0D92` | `rgb(75, 13, 146)` | `hsl(268, 84%, 31%)` | Deep violet structural accent |
| `--violet-900` | `oklch(0.257 0.153 294.5)` | `#2D0063` | `rgb(45, 0, 99)` | `hsl(267, 100%, 19%)` | Maximum depth violet anchor |

#### Cyan, Neutrals & Surfaces
| Token | Hex | HSL | Architectural Role |
| :--- | :--- | :--- | :--- |
| `--cyan-400` | `#03D5E7` | `hsl(185, 97%, 46%)` | **Brand Highlight (`--brand-highlight`)**, `.shine-sweep` |
| `--neutral-200` | `#E2E6ED` | `hsl(218, 23%, 91%)` | **Light Hairline Border (`--border`, `--input`)** |
| `--neutral-400` | `#8F9BB2` | `hsl(219, 19%, 63%)` | Dark-band muted text (`--on-dark-muted`) |
| `--neutral-600` | `#48536B` | `hsl(221, 20%, 35%)` | **Soft Ink (`--ink-soft`, `--muted-foreground`)** |
| `--neutral-900` | `#090E18` | `hsl(220, 45%, 6%)` | **Primary Ink (`--ink`, `--foreground`)** |
| `--surface-dark` | `#0A0E17` | `hsl(222, 39%, 6%)` | Dark section canvas / header bar |

### 2.2 Contrast Verification & Relative Luminance Formula
All color pairings must satisfy WCAG 2.2 AA (>= 4.5:1 normal text, >= 3:1 large text >= 24px) and target AAA (>= 7.0:1):
- Navy-Purple Theme: `#F8FAFC` on `#0A0F1D` = **16.4:1 (AAA)**.
- GitHub Dark High Contrast: `#FFFFFF` on `#0A0C10` = **21.0:1 (AAA)**.
- Warm Editorial Theme: `#141210` on `#FBF9F6` = **16.8:1 (AAA)**.

#### Fast Relative Luminance:
$$L = 0.299 R + 0.587 G + 0.114 B$$
- $L > 0.6 \implies$ dark ink `#090E18`.
- $L \le 0.6 \implies$ stark white `#FFFFFF`.

### 2.3 4-Plane Depth Hierarchy
- **Plane 0 (Canvas Base):** Deep ground (`#0A0F1D` dark / `#FBF9F6` light).
- **Plane 1 (Raised Wells):** Inset containers, alternating section bands (`#10172A`).
- **Plane 2 (Surface Cards):** Standard cards, tables, inputs, code block containers (`#172036`).
- **Plane 3 (Elevated Modals):** Floating modals, dropdowns, hovered cards with subtle edge glow (`#1E2945`).
- **Hairline Lighting:** `box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12), 0 4px 20px rgba(0, 0, 0, 0.25);`

### 2.4 60/30/10 Rule & Single-Accent Von Restorff Rule
- **60% Dominant Neutral:** Canvas ground and outer wrappers.
- **30% Structural Surfaces:** Card bodies, panels, typography, and borders.
- **10% Saturated Accent:** Reserved strictly for primary CTAs, active status, and key metrics.
- **Single-Accent Von Restorff Rule:** At any given viewport position, **exactly ONE prominent accent-filled CTA button** may be visible. Secondary actions MUST use outline, ghost, glass, or link styles.

---

## 3. The 10 Improvement Clusters Overview

1. **Cluster 01: Color Tokens, HSL Architecture & Zero Hardcoded Hex** (Items 01–10)
2. **Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Compliance** (Items 11–20)
3. **Cluster 03: 4-Plane Depth Hierarchy & Materiality** (Items 21–30)
4. **Cluster 04: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule** (Items 31–40)
5. **Cluster 05: Typography Stacks & Pure Live DOM Text Mandate** (Items 41–50)
6. **Cluster 06: Button System Standards & Hit Target Geometry** (Items 51–60)
7. **Cluster 07: Header, Navigation & Safe-Region Physics** (Items 61–70)
8. **Cluster 08: Spacing Scales, Containers & Band Rhythm** (Items 71–80)
9. **Cluster 09: CSS3 Hardware-Accelerated Motion & Accessibility** (Items 81–90)
10. **Cluster 10: Frontend Code Hygiene & React Component Deconstruction** (Items 91–100)

Refer to `01-prompts/30-improvements/01-next-100-ui-improvements.md` for the exhaustive 100-item specification catalog.

---

## 4. Workflow & GitMap Execution Protocol

Execute UI improvements using the canonical Tree Swarm Hierarchy with GitMap high-speed search primacy:

```text
Tree Hierarchy: Lead Orchestrator -> Branch Coordinators -> Leaf Workers (A = 3, H = 2, Capacity = 6)
```

1. **Step 0 Preflight & Cognitive Heatmap:**
   - Run `gitmap heatmap status` and `gitmap heatmap tables`.
   - Log pre-action intent:
     `gitmap heatmap log --task <slug> --action READ --file "<path>" --state EXPLORATION --expected "<goal>"`
2. **Phase 1 Planning & Spec Authoring:**
   - Spawn 3 read-only discovery subagents (`TypeName: "research"`).
   - Spawn 3 authoring subagents (`TypeName: "self"`) to author specs in `02-spec/21-app/XX-<slug>/`.
   - Enqueue master plan to `.ai-memory/plans/pending/XX-<slug>.md` and subtasks to `.ai-memory/plans/subtasks/XX-<slug>/`.
3. **Phase 2 Bounded Micro-Batch Execution:**
   - Spawn 3 worker subagents (`TypeName: "self"`).
   - Apply surgical fixes in disjoint bounding boxes (5-8 files per micro-batch).
   - Log edits: `gitmap heatmap log --task <slug> --action EDIT --file "<path>" --state EXECUTION --expected "<refactor>"`
4. **Phase 3 Verification & Atomic GitMap Commit:**
   - Run targeted file-scoped linters only (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`).
   - Run secrets gate (`gitmap aum search -r`, offload via `gitmap rs`).
   - Execute ONE atomic GitMap commit: `gitmap cpf "<module> - <summary>"` (Feature) or `gitmap cpb "<module> - <summary>"` (Bug fix) with hyphen format `-`.
   - Push to remote tracking branch.

---

## 5. Enqueueing Plans to `.ai-memory/plans/pending/`

Every execution cycle MUST preserve continuity:
1. Write the current run's execution plan and 100 improvements into `.ai-memory/plans/pending/XX-<slug>.md`.
2. Break down active batch into discrete subtask files under `.ai-memory/plans/subtasks/XX-<slug>/01-<subtask>.md`.
3. Record all completed improvements into `.ai-memory/plans/completed/`.
4. Systematically populate the subsequent 100 UI improvements into a new pending plan file (`.ai-memory/plans/pending/XX-next-100-ui-improvements.md`) to maintain the self-improving iteration loop.

---

## 6. Verification Criteria & Invariants

- **Zero Builds & Tests Ban (R1):** NEVER run `npm run build`, `vite build`, `npm test`, `go build`, or `pytest`.
- **Targeted Checks Only (R2):** Run only file-scoped fast linters on modified files.
- **Strict Relative Git Paths (R11):** TOTAL BAN on absolute paths and `file:///` URIs.
- **Strict Lowercase (Rule 4):** All files, scripts, and documentation strictly lowercase.
- **Positive Booleans (Rule 1):** Positive checks only (`isReady`), zero explicit `== true`, zero mixed polarity (`if isA && !isB` BANNED).
- **Pure DOM Text Mandate:** Zero baked raster text in PNG/JPEG images.
- **Component Sizing:** All React components <= 80–100 lines.
- **Mandatory Commit & Push (Section 13):** If not committed and pushed via `gitmap cpf` / `gitmap cpb`, the task is strictly NOT DONE.

---

## Related Skills

- [follow-ui-ux-design-system](file;.cursor/skills/follow-ui-ux-design-system)
- [theming-guidelines](file;.cursor/skills/theming-guidelines)
- [update-theming](file;.cursor/skills/update-theming)
- [react-frontend-guidelines](file;.cursor/skills/react-frontend-guidelines)
- [execute-parent-task-with-n-steps-v6](file;.cursor/skills/execute-parent-task-with-n-steps-v6)

# High Priority Instruction

# Next 100 UI Improvements

## slug: next-100-ui-improvements

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /next-100-ui-improvements <target-component-or-page>
>
> **Trigger Keywords & Aliases:**
> `next-100-ui-improvements`, `next-100-ui`, `ui-100-improvements`, `next-ui-improvements`, `find-100-ui`
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.

[/goal](slashCommand;goal) Autonomously audit, discover, catalog, and refactor 100 concrete UI improvements across the target frontend interface. FIRST showcase and list out the given task in visible chat during Turn 1 with the 5 Mandatory Pillars and 6-part verification schema, capture it verbatim, plan it in the repository under `02-spec/21-app/XX-<slug>/`, enqueue the complete catalog of 100 improvements as pending plans into `.ai-memory/plans/pending/XX-next-100-ui-improvements.md` along with subtasks in `.ai-memory/plans/subtasks/XX-next-100-ui-improvements/`, spawn autonomous subagents via `invoke_subagent` (A = 3, H = 2, Capacity = 6; solo execution without calling `invoke_subagent` is an auto-reject failure) in disjoint file boxes using GitMap high-speed commands as primary, prove every single claim with concrete evidence, enforce coding guidelines to 100%, and finish with one atomic GitMap commit (`gitmap cpf "<module> - <summary>"` / `gitmap cpb "<module> - <summary>"`) using hyphen format (no colons in GitMap arguments, as the colon is already provided by GitMap) that holds strictly this task's files. Every execution turn MUST resolve current target improvements and generate the subsequent 100 pending UI improvement tasks in `.ai-memory/plans/pending/`.

[/learn](slashCommand;learn) Enforce the Top-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are provided ABOVE this prompt outrank everything below. Turn 1 MUST showcase the given task list in visible chat before any background execution. Master the 10 UI improvement clusters, color ramps, 4-plane depth hierarchy, 60/30/10 visual balance, typography stacks, button system, and safe-region navigation physics. Progress lives in the ledger and in `.ai-memory/plans/`, never only in chat.

[/plan](slashCommand;plan) Execute thorough step-by-step planning in the repository before execution for improved work or less friction, faster work and more pure and effective work. Ensure all deliverables, architecture boundaries, and requirements are clearly defined in the audit ledger and subtask plans before dispatching worker waves.

### 🚨 MANDATORY SUBAGENT SPAWNING GATE (A = 3, H = 2, CAPACITY = 6 — TREE SWARM & DYNAMIC JOIN)

- **ACTUAL TOOL CALL REQUIRED:** You must ACTUALLY CALL the `invoke_subagent` tool via your tool-calling API. Do NOT just print the text "Dispatched Worker..." and stop. If you only print text, the agents will not spawn and the task will fail! You must execute the `invoke_subagent` JSON tool payload.
- **TREE-STRUCTURED SWARM HIERARCHY:** Organize execution as an operational tree:
  1. **Root Orchestrator (Lead Agent):** Owns the unified Brain roadmap, audit ledger, and final GitMap commit. Never executes file edits solo.
  2. **Branch Coordinators (Subtask Primary Owners):** Own designated disjoint file boundaries (`Owned Files`) and execute modifications.
  3. **Leaf Workers & Dynamic Joining Assistants:** Assist primary owners with parallel reading, symbol caller discovery, memory generation, and companion lint verification.
- **3-STAGE MANDATORY DISPATCH:** You must invoke `A = 3` agents (`invoke_subagent`) at EVERY stage of the workflow:
  1. **Planning Step (A = 3 `research` subagents):** Spawn 3 read-only discovery subagents (`TypeName: "research"`) to research the codebase and return findings; lead writes the unified Brain execution plan and roadmap.
  2. **Spec Step (A = 3 self authoring subagents — ZERO SOLO SPEC AUTHORING):** Spawn 3 authoring subagents (TypeName: "self") to author modular, disjoint spec files and subtasks (never writing the same file). Solo authoring by the lead orchestrator is STRICTLY FORBIDDEN.
  3. **Execution Step (A = 3 `self` worker subagents):** Spawn 3 worker subagents (`TypeName: "self"`) to execute code modifications in strictly disjoint file boxes.
- **MANDATORY 3 SUBAGENTS IN SINGLE TOOL CALL:** Every call to invoke_subagent MUST contain an array of EXACTLY 3 subagent objects: `"Subagents": [agent1, agent2, agent3]`. Calling invoke_subagent with fewer than 3 subagent objects, or executing solo without calling invoke_subagent, is an immediate protocol violation and auto-reject failure.
- **DYNAMIC AGENT JOIN & IMMEDIATE RESPAWN (NO IDLE WORKERS):** Worker capacity `A` must remain pinned at 3 active agents (System Concurrency Capacity = 6) throughout execution. The lead orchestrator MUST NOT wait for an entire 3-agent wave to finish. As soon as ANY single worker completes, the lead immediately spawns a replacement worker for the next pending subtask, or commands the worker to dynamically join an active subtask in collaborative assistance mode (parallel reading, companion linters, memory logging) to keep concurrency at exactly 3 active subagents (100% duty cycle).
- **SOLO EXECUTION IS AN AUTO-REJECT FAILURE:** The lead orchestrator is **STRICTLY FORBIDDEN** from executing planning, spec writing, or code changes by itself without calling the `invoke_subagent` tool. Failing to call the actual tool is a critical protocol violation.

---

## Phase 0: Antigravity Skill Bootstrap (Memory Optimization)

Before executing the tasks below, check if this prompt is already installed as a native Antigravity Skill.

1. If `.agents/skills/next-100-ui-improvements/skill.md` does not exist in the workspace, create it now.
2. Extract the core instructions of this prompt and save it into that `skill.md` using the standard YAML frontmatter (with `name` and `description`).
3. Maintain symmetric mirroring in `.cursor/skills/next-100-ui-improvements/skill.md` with IDE link syntax `[next-100-ui-improvements](file;.cursor/skills/next-100-ui-improvements)`.
4. Once installed, rely on progressive disclosure for future runs. Do not keep the entire prompt in active memory if not needed.

---

## Comprehensive Color Tokens, HSL Architecture & Materiality Standards

### 1. Primitive Color Ramps (3-Format Reference)

All colors MUST be variable-driven and declared as CSS custom properties in `:root` and `.dark` blocks in HSL space-separated format (e.g. `221 83% 53%`) to allow opacity modifiers via `hsl(var(--token) / 0.15)`.

#### 1.1 Blue Ramp (Primary Action & Corporate Trust)
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--blue-50` | `oklch(0.972 0.014 259)` | `#F0F6FF` | `rgb(240, 246, 255)` | `hsl(216, 100%, 97%)` | Hero wash start, brand pill background |
| `--blue-100` | `oklch(0.935 0.032 262)` | `#DEEAFF` | `rgb(222, 234, 255)` | `hsl(218, 100%, 94%)` | Badge hover fill, soft tile ring |
| `--blue-200` | `oklch(0.874 0.061 263)` | `#C1D6FF` | `rgb(193, 214, 255)` | `hsl(220, 100%, 88%)` | Chrome gradient stop, subtle border |
| `--blue-300` | `oklch(0.777 0.104 264)` | `#95B6FA` | `rgb(149, 182, 250)` | `hsl(220, 91%, 78%)` | Chrome gradient end, data series |
| `--blue-400` | `oklch(0.662 0.164 264)` | `#5E8EF6` | `rgb(94, 142, 246)` | `hsl(221, 89%, 67%)` | Gradient text end stop, ring focus |
| `--blue-500` | `oklch(0.546 0.215 262.9)` | `#2563EB` | `rgb(37, 99, 235)` | `hsl(221, 83%, 53%)` | **Primary CTA Action (`--brand-secondary`)** |
| `--blue-600` | `oklch(0.489 0.208 263.5)` | `#1B51D3` | `rgb(27, 81, 211)` | `hsl(222, 77%, 47%)` | Brand gradient start, active CTA |
| `--blue-700` | `oklch(0.409 0.176 264.5)` | `#173DA7` | `rgb(23, 61, 167)` | `hsl(224, 76%, 37%)` | Brand pill text, high-contrast link |
| `--blue-800` | `oklch(0.369 0.147 264.8)` | `#17368B` | `rgb(23, 54, 139)` | `hsl(224, 72%, 32%)` | Deep interactive border, badge text |
| `--blue-900` | `oklch(0.317 0.135 264.5)` | `#0D2975` | `rgb(13, 41, 117)` | `hsl(224, 80%, 25%)` | **Core Brand Anchor (`--brand-primary`)** |

#### 1.2 Violet Ramp (Tertiary Accent & Gradient Partner)
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

#### 1.3 Cyan, Neutrals & Cool Slate Surfaces
| Token | OKLCH Runtime Value | Format 1: HEX | Format 2: RGB | Format 3: HSL | Architectural Role |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `--cyan-400` | `oklch(0.797 0.136 205.6)` | `#03D5E7` | `rgb(3, 213, 231)` | `hsl(185, 97%, 46%)` | **Brand Highlight (`--brand-highlight`)**, `.shine-sweep` |
| `--neutral-200` | `oklch(0.924 0.011 265)` | `#E2E6ED` | `rgb(226, 230, 237)` | `hsl(218, 23%, 91%)` | **Light Hairline Border (`--border`, `--input`)** |
| `--neutral-400` | `oklch(0.688 0.036 264)` | `#8F9BB2` | `rgb(143, 155, 178)` | `hsl(219, 19%, 63%)` | Dark-band muted text (`--on-dark-muted`) |
| `--neutral-600` | `oklch(0.441 0.043 265.9)` | `#48536B` | `rgb(72, 83, 107)` | `hsl(221, 20%, 35%)` | **Soft Ink (`--ink-soft`, `--muted-foreground`)** |
| `--neutral-900` | `oklch(0.163 0.023 265)` | `#090E18` | `rgb(9, 14, 24)` | `hsl(220, 45%, 6%)` | **Primary Ink (`--ink`, `--foreground`)** |
| `--surface-dark` | `oklch(0.129 0.019 264)` | `#0A0E17` | `rgb(10, 14, 23)` | `hsl(222, 39%, 6%)` | Dark section canvas / header bar |

### 2. Contrast Verification & Fast Relative Luminance

All interface combinations MUST strictly comply with WCAG 2.2 AA (minimum 4.5:1 for normal text, 3:1 for large text >= 18pt/24px) and target AAA (7.0:1):
- **Navy-Purple Theme:** Primary text (`#F8FAFC`) on Canvas (`#0A0F1D`) achieves **16.4:1 (AAA)**.
- **GitHub Dark High Contrast:** Pure white (`#FFFFFF`) on Base (`#0A0C10`) achieves **21.0:1 (AAA)**.
- **Warm Editorial Theme:** Deep Ink (`#141210`) on Warm Paper (`#FBF9F6`) achieves **16.8:1 (AAA)**.
- **Electric Violet CTA Text:** Primary ink (`#0A0F1D`) on Violet (`#A855F7`) achieves **11.2:1 (AAA)**.

#### Fast Relative Luminance Calculation:
$$L = 0.299 R + 0.587 G + 0.114 B$$
- If $L > 0.6 \implies$ use dark ink foreground: `#090E18` / `hsl(220, 45%, 6%)`.
- If $L \le 0.6 \implies$ use stark white foreground: `#FFFFFF` / `hsl(0, 0%, 100%)`.

### 3. The 4-Plane Depth Hierarchy (Materiality Over Shadows)
Depth in modern light and dark modes is governed by lightness steps and hairline highlights rather than heavy, muddy drop shadows:
- **Plane 0 (Canvas Base):** Deepest page ground (`#0A0F1D` dark / `#FBF9F6` light).
- **Plane 1 (Raised Wells):** Inset containers, alternating section bands, and sidebars (+3% to 4% lightness, `#10172A`).
- **Plane 2 (Surface Cards):** Standard cards, tables, inputs, and code block containers (+3% to 4% lightness, `#172036`).
- **Plane 3 (Elevated Modals & Flyouts):** Modals, dropdown menus, and hovered cards with subtle edge glow (`#1E2945`).
- **Hairline Lighting:** Cards and panels carry a physical light refraction highlight:
  ```css
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12), 0 4px 20px rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(248, 250, 252, 0.08);
  ```

### 4. 60 / 30 / 10 Visual Balance & Single-Accent Von Restorff Rule
- **60% Dominant Neutral Base:** Deep canvas ground providing visual rest and structural canvas.
- **30% Structural Surfaces & Typography:** Content panels, tables, structural cards, and secondary ink.
- **10% Purposeful Saturated Accent:** Reserved strictly for primary CTAs, active pills, and singular focal metrics.
- **Single-Accent Von Restorff Principle:** At any viewport position, **exactly ONE prominent accent-filled element** may be visible (e.g. one solid Primary CTA). Secondary actions must use outline, ghost, glass, or link treatments. Two competing primary buttons side-by-side collapses the visual hierarchy.

### 5. Pure Live DOM Typography Mandate
- **Primary Font Stacks:**
  - Headings (`h1` - `h6`): Ubuntu, sans-serif (`font-heading`)
  - Body Text & UI: Poppins, sans-serif (`font-body`)
  - Inline Code & Badges: JetBrains Mono, monospace
  - Code Blocks & Terminal: Ubuntu Mono, JetBrains Mono, monospace
- **Pure DOM Text Mandate:** TOTAL BAN on baking text into raster images, canvas bitmaps, or static PNG graphics. All copy must be live, selectable HTML elements with proper semantic markup (`<h1>`, `<p>`, `<span>`, `<kbd>`).
- **Copy Length Caps:** Section Eyebrow <= 3 words; Section H2 <= 8 words; Lead Paragraph <= 28 words; Card Body <= 32 words; at most ONE gradient accent phrase (`<GradientText>`) per heading.

### 6. Universal Button System Standards
- **6 Semantic Variants:**
  1. `primary`: `bg-[image:var(--gradient-accent)] text-white shadow-[var(--shadow-card)]` with active `.shine-sweep`.
  2. `solid`: `bg-brand-primary text-white shadow-[var(--shadow-card)]` with active `.shine-sweep`.
  3. `outline`: `border border-border bg-transparent text-foreground` with active `.pointer-fill`.
  4. `glass`: `bg-[image:var(--gradient-card-dark)] text-on-dark backdrop-blur-md` with active `.gradient-ring`.
  5. `ghost`: `bg-transparent text-foreground hover:bg-muted`.
  6. `link`: `text-brand-secondary px-0 hover:underline underline-offset-4`.
- **4 Sizing Standards:**
  - `sm: 36px` (`h-9 px-4 text-[13px]`, icon `14px`)
  - `md: 44px` (`h-11 px-5 text-sm`, icon `16px`) — Default Touch Target
  - `lg: 52px` (`h-13 px-7 text-[15px]`, icon `18px`)
  - `icon: 44px` (`size-11 px-0`, icon `20px`)
- **Micro-Interactions:** Continuous angled 45-degree light sweep (`.shine-sweep`), radial cursor pointer fill (`.pointer-fill`), magnetic cursor pull (`strength: 0.22`), and tactile active state (`scale(0.98)`).

### 7. Header Navigation & Safe-Region Physics
- **Header Geometry:** Sticky 72px glassmorphic header (`h-[72px] top-0 z-50 w-full`), `backdrop-blur-xl`.
- **12px Scroll Trigger:** Scroll > 12px adds `border-b border-border shadow-[var(--shadow-card)]` transitioning over 420ms `cubic-bezier(0.16, 1, 0.3, 1)`.
- **Safe-Region Physics:** Safe region buffer `pad = 14px` around both header and dropdown panel rects.
- **Debounced Close:** Moving outside safe region initiates `220ms` debounce timer; moving back inside cancels close. `Escape` key or window scroll dismisses immediately (`0ms`).
- **`SlideSwapLabel` Nav Links:** Per-character vertical slide swap on hover with `stagger: 0.04s` over 520ms; collapses to plain span under `prefers-reduced-motion`.
- **Sticky-Safe Mobile Drawer:** Fixed `top-[72px] bottom-0` drawer intercepting touch events WITHOUT setting `overflow: hidden` on body (which breaks `position: sticky`).

---

## The 10 Structured Improvement Clusters (100 Improvements Total)

Each execution run audits the frontend codebase against these 100 improvements, implements target fixes in bounded batches, and enqueues all pending items in `.ai-memory/plans/pending/`:

### Cluster 01: Color Tokens, HSL Architecture & Zero Hardcoded Hex
1. **Eliminate Hardcoded Hex in Components:** Replace raw hex values (`#2563eb`, `#1e293b`) with semantic tokens (`--primary`, `--surface`) to ensure theme compatibility.
2. **Space-Separated HSL Token Storage:** Convert tokens to space-separated HSL (`H S% L%`) without wrapper to enable alpha composition (`hsl(var(--primary) / 0.15)`).
3. **10-Step Blue Ramp Implementation:** Map all blue hues to the standard `--blue-50` through `--blue-900` tokens with exact OKLCH/HEX/HSL values.
4. **10-Step Violet Ramp Implementation:** Standardize tertiary accents across `--violet-50` through `--violet-900` for harmonious brand gradient pairings.
5. **Cyan Highlight Integration:** Deploy `--cyan-400` (`#03D5E7`) exclusively for high-conversion sheens, live telemetry, and active status rings.
6. **Neutral Slate Surface Calibration:** Align borders, inputs, and muted copy to `--neutral-200`, `--neutral-400`, and `--neutral-600`.
7. **Semantic Token Decoupling:** Banish literal color classes (`bg-blue-500`, `text-purple-600`) in favor of purpose-driven tokens (`--primary`, `--accent`).
8. **Dark/Light Theme Variable Parity:** Ensure every token in `:root` has a calibrated twin in `.dark` with +5-10% lightness offset on dark accents.
9. **Tokenized Gradient Definitions:** Declare all brand, hero, and chrome gradients as CSS custom properties in stylesheet root rather than inline styles.
10. **Code Block Theme Isolation:** Pin code block containers to fixed dark slate (`hsl(220, 14%, 11%)`) across all themes to preserve syntax token legibility.

### Cluster 02: Contrast Ratios & WCAG 2.2 AA/AAA Compliance
11. **Body Text Contrast AA/AAA Pass:** Ensure all paragraph and descriptive text achieves >= 4.5:1 (AA) and targets >= 7.0:1 (AAA) against parent surfaces.
12. **Large Display Text Contrast:** Verify headings and hero numbers >= 24px achieve >= 3.0:1 contrast against background canvases.
13. **High-Contrast Theme Mode:** Implement verified high-contrast themes (GitHub Dark High Contrast 21.0:1 AAA, Navy-Purple 16.4:1 AAA).
14. **Automated Relative Luminance Check:** Implement $L = 0.299 R + 0.587 G + 0.114 B$ evaluation across all dynamic color combinations.
15. **Luminance Threshold Foreground Selection:** Enforce dark ink (`#090E18`) when $L > 0.6$ and pure white (`#FFFFFF`) when $L \le 0.6$.
16. **Focus Ring Contrast Ratio:** Ensure `:focus-visible` ring token (`--ring`) achieves >= 3:1 contrast against adjacent background colors.
17. **Form Control Border Contrast:** Calibrate input, select, and textarea borders to achieve >= 3:1 contrast against surrounding card surfaces.
18. **Placeholder Text Differentiation:** Set placeholder contrast to >= 3.5:1 while remaining clearly visually secondary to user input values.
19. **Link Distinction in Prose:** Guarantee in-text hyperlinks have underline styling or >= 3:1 contrast ratio against surrounding body text.
20. **Non-Color Visual Redundancy:** Pair all color-coded status pills (error, warning, success) with distinct icons or descriptive labels for colorblind accessibility.

### Cluster 03: 4-Plane Depth Hierarchy & Materiality
21. **Plane 0 Canvas Foundation:** Establish the deepest page ground (`#0A0F1D` dark / `#FBF9F6` light) as the base viewport foundation.
22. **Plane 1 Raised Wells:** Structure alternating section bands, table wells, and sidebars at +3% to 4% lightness (`#10172A`).
23. **Plane 2 Surface Cards:** Elevate cards, interactive panels, inputs, and code blocks at +3% to 4% additional lightness (`#172036`).
24. **Plane 3 Elevated Modals:** Configure floating modals, popovers, and hovered cards at Plane 3 (`#1E2945`) with focused boundary glow.
25. **Hairline Top-Edge Lighting:** Add `box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.12)` to simulate physical top-edge light refraction.
26. **Subtle Hairline Dividers:** Replace dark drop shadows with crisp `1px` translucent borders (`rgba(248, 250, 252, 0.08)`).
27. **Materiality Over Shadow Diffusion:** Create perceived depth through stepped lightness intervals rather than wide, muddy blur shadows.
28. **Layered Radial Halos:** Position subtle radial gradient halos (8-10% opacity) behind hero cards to generate ambient focal illumination.
29. **Frosted Acrylic Glassmorphism:** Apply `backdrop-blur-xl` and 90% opacity surface fills to floating and sticky headers.
30. **Restrained Hover Elevation:** Standardize interactive card hover lift to `translateY(-2px)` with subtle shadow enhancement (`--shadow-lift`).

### Cluster 04: 60/30/10 Visual Balance & Single-Accent Von Restorff Rule
31. **60% Dominant Neutral Allocation:** Ensure 60% of viewport surface area comprises neutral background base planes.
32. **30% Structural Surface Rhythm:** Allocate 30% of surface area to structural cards, content panels, tables, and typography.
33. **10% Saturated Accent Discipline:** Limit saturated primary brand color to at most 10% of visible pixels (CTAs, active pills, key metrics).
34. **Single-Accent Von Restorff Principle:** Guarantee exactly ONE prominent accent-filled CTA button is visible per viewport fold.
35. **Primary/Secondary CTA Pairing:** Pair every solid primary button with an outline, ghost, or link variant to prevent visual shouting.
36. **Elimination of Competing Primary Buttons:** Refactor side-by-side solid buttons in heroes and cards into primary + secondary hierarchy.
37. **Dashboard Focal Metric Isolation:** Highlight only the single most critical KPI card with accent styling while keeping sibling cards neutral.
38. **Restrained Notification Cards:** Replace full-bleed saturated alert banners with neutral cards carrying a 2px left accent hairline.
39. **Single Gradient Text Phrase:** Restrict `<GradientText>` spans to at most ONE short phrase (2-4 words) per headline.
40. **Anti-AI-Slop Layout Enforcement:** Eliminate generic 3-card icon grids, purple-blue canvas soup, and adjective-heavy marketing claims.

### Cluster 05: Typography Stacks & Pure Live DOM Text Mandate
41. **Ubuntu Heading Font Stack:** Enforce Ubuntu with weights 400, 500, and 700 for all `h1` through `h6` headings.
42. **Poppins Body Font Stack:** Enforce Poppins with weights 300, 400, 500, and 600 for body prose, labels, and table cells.
43. **JetBrains Mono Inline Code Stack:** Apply JetBrains Mono for inline code pills, chip metadata, and technical attributes.
44. **Ubuntu Mono Code Block Stack:** Apply Ubuntu Mono for multi-line syntax-highlighted code blocks and terminal logs.
45. **Pure Live DOM Mandate (TOTAL BAN ON RASTER TEXT):** Banish baking text into PNG/JPEG images; all text must be selectable live DOM.
46. **Fluid Type Scaling via CSS Clamp:** Implement `clamp()` formulas for responsive headings (e.g. `clamp(1.75rem, 4vw, 3rem)` for H1).
47. **Typographic Measure & Line Height:** Set body text line-height to `1.6-1.65` with optimal reading measure capped at 65-75 characters.
48. **Heading Tracking Calibration:** Apply subtle negative tracking (`-0.02em` to `-0.01em`) on large display titles for optical tightness.
49. **Section Eyebrow Word Cap:** Restrict section eyebrows to at most 3 words in uppercase, semi-bold, tracked styling.
50. **Section Copy Length Budget:** Enforce H2 <= 8 words, lead paragraph <= 28 words, and card body <= 32 words.

### Cluster 06: Button System Standards & Hit Target Geometry
51. **6 Canonical Button Variants:** Implement `primary`, `solid`, `outline`, `glass`, `ghost`, and `link` variants strictly mapped to design tokens.
52. **4 Proportional Sizing Scales:** Standardize button heights: `sm: 36px`, `md: 44px` (default), `lg: 52px`, `icon: 44px`.
53. **Minimum 44x44px Touch Targets:** Ensure all interactive clickable elements maintain at least 44×44px hit bounds on mobile.
54. **Continuous Sheen Sweep (`.shine-sweep`):** Add 45-degree angled light sweep animating on hover/focus over 3s linear infinite.
55. **Radial Cursor Pointer-Fill (`.pointer-fill`):** Equip outline buttons with smooth interactive background fill from pointer entry angle.
56. **Gradient Ring Border Mask:** Apply multi-color accent border mask (`.gradient-ring`) on frosted glass buttons.
57. **High-Conversion ShineButton CTA:** Deploy near-black pill (`bg-night`) with continuous 45-degree sheen sweep for primary header CTA.
58. **Magnetic Cursor Physics:** Wrap hero primary CTAs in `<Magnetic strength={0.22}>` for subtle cursor attraction.
59. **Tactile Active Feedback:** Implement immediate `scale(0.98)` depression and reduced shadow on `:active` button states.
60. **High-Contrast Capsule Pills:** Standardize telemetry chips and filter pills using `--capsule-gold` / `--capsule-ember` tokens.

### Cluster 07: Header, Navigation & Safe-Region Physics
61. **Sticky 72px Glassmorphic Header:** Standardize header height strictly to `72px` (`h-[72px] top-0 z-50 w-full backdrop-blur-xl`).
62. **12px Viewport Scroll Trigger:** Add bottom border and shadow when scroll exceeds 12px with 420ms cubic-bezier transition.
63. **Safe-Region Pointer Geometry:** Implement `pad = 14px` buffer around navigation trigger and dropdown panel rectangles.
64. **Debounced Close Physics:** Delay dropdown menu closure by `220ms` when pointer exits safe region; cancel timer on re-entry.
65. **Instant Dismissal on Escape & Scroll:** Dismiss open navigation dropdowns immediately (`0ms`) upon pressing Escape or window scrolling.
66. **`SlideSwapLabel` Nav Link Hover:** Implement per-character vertical slide swap with `stagger: 0.04s` over 520ms transition.
67. **Animated Nav Underline Indicator:** Add `1px` hairline underline indicator scaling from center or left (`scale-x-0 -> scale-x-100` over 240ms).
68. **Rotating Chevron Indicator:** Equip menu triggers with 14px chevron rotating 180deg over 240ms when active.
69. **Avant-Garde MegaPanel Dropdown:** Mount full-width card dropdown with 260ms entrance (`opacity: 0, y: -8, scale: 0.985`).
70. **Sticky-Safe Mobile Drawer:** Construct mobile navigation drawer fixed `top-[72px] bottom-0` without setting `overflow: hidden` on body.

### Cluster 08: Spacing Scales, Containers & Band Rhythm
71. **4px/8px Modular Spacing Grid:** Standardize all padding, margin, and gap utilities strictly to multiples of 4px and 8px.
72. **1280px Container Constraint:** Enforce `max-w-7xl` (1280px) content boundaries with responsive horizontal gutters (`px-4 sm:px-6 lg:px-8`).
73. **Alternating Section Band Rhythm:** Cycle section background tones (`light` -> `soft` -> `light` -> `soft`) with zero adjacent duplicates.
74. **Single Dark Void Section Allocation:** Permit at most ONE full dark void section per page for embedded product demos or testimonials.
75. **Standard Vertical Section Padding:** Apply `clamp(72px, 8vw, 120px)` vertical padding on standard page sections.
76. **Tight Vertical Section Padding:** Apply `clamp(40px, 5vw, 64px)` vertical padding on compact interstitial callouts.
77. **Standardized Card Grid Gaps:** Enforce uniform `gap-6` (24px) or `gap-8` (32px) between cards across responsive breakpoints.
78. **Card Internal Padding Scale:** Apply `p-6` (24px) on standard cards and `p-8` (32px) on featured pricing or hero cards.
79. **Corner Radius Hierarchy:** Enforce 8px (`rounded-lg`) on inputs, 12px (`rounded-xl`) on buttons, 20px (`rounded-2xl`) on cards, 9999px on pills.
80. **Flexbox Bottom Action Alignment:** Apply `mt-auto` to card footer buttons to guarantee clean horizontal baseline alignment in grid rows.

### Cluster 09: CSS3 Hardware-Accelerated Motion & Accessibility
81. **GPU-Accelerated CSS3 Properties:** Restrict all interactive animations strictly to `transform` and `opacity` properties.
82. **Total Ban on Heavy JS Animation Libraries:** Ban `framer-motion`, `gsap`, `anime.js`, and `react-spring` to preserve runtime performance.
83. **Standardized Transition Timings:** Apply `0.2s` for color/shadow transitions and `0.15s` for transforms with natural deceleration easing.
84. **Restrained Hover Transforms:** Cap interactive hover transforms to `translateY(-2px)` and `scale(1.02)`, avoiding bouncy spring physics.
85. **Prefers-Reduced-Motion Compliance:** Collapse all transition durations to `0.01ms` when user enables reduced motion preferences.
86. **Infinite Linear Marquee Animation:** Implement hardware-accelerated CSS `@keyframes marquee` with pause-on-hover capability.
87. **Smooth Grid Height Accordions:** Animate collapsible accordions using CSS `grid-template-rows: 0fr -> 1fr` without JavaScript height polling.
88. **3D Card Flip Perspective:** Implement `perspective: 1400px` and `rotateY(180deg)` over 820ms for interactive flip cards.
89. **Non-Blocking Focus Rings:** Implement crisp `2px` focus rings with `2px` offset visible exclusively during keyboard navigation (`:focus-visible`).
90. **Mobile Tap Highlight Sanitization:** Remove intrusive grey tap highlights on mobile browsers via `-webkit-tap-highlight-color: transparent`.

### Cluster 10: Frontend Code Hygiene & React Component Deconstruction (<=80-100 lines)
91. **Component Length Cap (<= 80-100 lines):** Deconstruct any React component exceeding 80–100 lines into focused micro-components.
92. **Separation of Presentation and Data:** Extract data fetching, complex state machines, and business logic into dedicated custom hooks.
93. **Dedicated TypeScript Prop Interfaces:** Define typed props interfaces in `types.ts` with zero raw `any` or loose parameter bloat.
94. **Strict Boolean Naming (`is` / `has`):** Enforce `is` and `has` prefixes (`isLoading`, `hasBorder`) evaluated implicitly (`if isVisible`).
95. **Elimination of Nested JSX Ternaries:** Refactor nested ternary operators in JSX templates into early returns or render helper functions.
96. **Semantic HTML Landmark Usage:** Replace generic `<div>` containers with semantic `<header>`, `<nav>`, `<main>`, `<section>`, `<footer>`.
97. **Image Aspect Ratios & CLS Prevention:** Enforce explicit `aspect-ratio` or `width`/`height` attributes on all image elements to prevent CLS.
98. **Accessible SVG Icon Wrappers:** Ensure all SVG icons feature `aria-hidden="true"` or accessible `<title>` elements without raw SVG clutter.
99. **Automated Plan Enqueueing:** Enqueue discovered UI improvements into `.ai-memory/plans/pending/XX-next-100-ui-improvements.md` and subtasks.
100. **Continuous 100-Improvement Cycle:** Design execution so that resolving the current batch systematically generates the next 100 UI improvements.

---

## Actionable Items & Master Execution Pipeline

Execute this prompt via a strict 3-Phase pipeline with mandatory subagents and GitMap acceleration:

### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)
1. **Top-Instruction Priority Verification:** Inspect preamble directives above this prompt and verify them as highest priority.
2. **Turn 1 Extreme Showcase:** In visible chat during Turn 1, output the complete 5 Mandatory Pillars and 6-part task breakdown before running background tools:
   - **Pillar 1: Strategic Intent & Systems Reasoning:** Explain the UI architecture goals, visual hierarchy, and user intent.
   - **Pillar 2: Explicit Scope Boundaries:** Strict In-Scope vs. Excluded/Deferred two-column matrix.
   - **Pillar 3: Non-Negotiable Constraints & Negative Guardrails:** Ban on raw search (`rg`, `Select-String`), worker git ban, zero builds/tests ban, relative paths mandate, lowercase naming, zero secrets.
   - **Pillar 4: Blast Radius & Component Hierarchy:** Map impacted frontend components, styles, design tokens, and callers.
   - **Pillar 5: Multi-Agent Concurrency Topology:** Swarm structure (`A = 3`, `H = 2`, Capacity = 6 concurrent operations), disjoint boundaries, dynamic join.
3. **Task Comprehension & Model Runtime Invariants:** In each task breakdown item, explicitly include:
   `- **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints]`
   `- **Azure AI / Model Runtime Context:** [N/A — Local CLI/GitMap or specific model endpoint]`
4. **Tool Chaining:** Emit chat output and immediately call the first tool in the same turn without pausing.

### Phase 1B: Spec Authoring & Subtask Planning Gate
1. **Spawn 3 Discovery Subagents (`invoke_subagent`):** Spawn `A = 3` read-only subagents (`TypeName: "research"`) to scan the codebase using GitMap high-speed search (`gitmap aum search`, `gitmap find`, `gitmap cat`).
2. **Author Specification Files:** Spawn `A = 3` authoring subagents (`TypeName: "self"`) to author modular specifications under `02-spec/21-app/XX-<slug>/`:
   - `01-architecture-and-tokens-spec.md`
   - `02-component-and-hierarchy-spec.md`
   - `03-motion-and-hygiene-spec.md`
3. **Enqueue Pending Plans & Subtasks:**
   - Write master pending plan: `.ai-memory/plans/pending/XX-next-100-ui-improvements.md`
   - Write bounded subtask files: `.ai-memory/plans/subtasks/XX-next-100-ui-improvements/01-tokens-and-contrast.md`, `02-depth-and-buttons.md`, etc.
   - Catalog all remaining items into `.ai-memory/plans/pending/` to preserve the continuous 100-improvement cycle.
4. **Initialize SQLite Task DB:** Initialize or inspect task state via GitMap agent task engine:
   `gitmap agent task init --name "<task name>" --budget 400`
   `gitmap agent subtask add --json '[{"code": "Task-01", "title": "<title>", "owned_files": ["<paths>"], "agent_role": "Worker 01"}]'`

### Phase 2: Worker Wave Execution (Bounded 5-8 File Micro-Batches)
1. **Spawn 3 Worker Subagents (`invoke_subagent`):** Invoke `A = 3` workers (`TypeName: "self"`) in disjoint file boxes (strictly one owner per file).
2. **Apply Surgical UI Refactorings:**
   - Migrate hardcoded colors to semantic HSL CSS custom properties.
   - Enforce WCAG 2.2 AA/AAA contrast ratios and luminance thresholds.
   - Implement the 4-plane depth hierarchy with hairline lighting.
   - Apply the 60/30/10 balance and single-accent Von Restorff rule.
   - Standardize typography stacks and enforce pure live DOM text.
   - Refactor button components to the 6 canonical variants and 4 sizes.
   - Calibrate sticky 72px glass header with 14px safe-region physics and 220ms debounce.
   - Deconstruct components exceeding 80–100 lines into focused micro-components.
3. **Cognitive Heatmap Logging:** Before reading or editing files, log actions via `gitmap heatmap log` or `gitmap agent log`. Inspect leases via `gitmap agent heatmap`.
4. **Dynamic Agent Join & Respawn:** When a worker completes, immediately spawn the next worker or join an active subtask to keep concurrency pinned at 3 active subagents (100% duty cycle).

### Phase 3: Targeted Verification & Atomic GitMap Commit
1. **Targeted File-Scoped Linter Verification:** Run fast file-scoped checks on modified files (e.g. `python 03-ai-scripts/05-guideline-autofixer.py <folder> --check-only`). A check scanning 0 files is a FAIL.
2. **Strict No-Build / No-Test Ban:** NEVER run full test suites, `go test`, `npm test`, `pytest`, `npm run build`, `vite build`. Only verify targeted files.
3. **Secrets Gate:** Check git status, scan for credentials/tokens via `gitmap aum search -r`, offload any private values to `repo-secrets` via `gitmap rs`.
4. **Atomic Push via GitMap:** Conclude the entire parent task with ONE atomic GitMap commit:
   `gitmap cpf "<module> - <summary>"` (Feature) or `gitmap cpb "<module> - <summary>"` (Bug fix)
   Format with hyphen `-` (no colons in GitMap argument, as GitMap provides `Feature: ` or `Bug: ` prefix automatically).
   Push to GitHub remote tracking branch. If not committed and pushed, the task is NOT DONE.

---

## Core Operational Rules (R1 to R16)

- **R1 Zero Builds, Test Suites & CI/CD Verification Ban (TOTAL BAN).** NEVER run `go build`, `npm run build`, `vite build`, `cargo build`, `go test ./...`, `pytest`, `npm test`, or `03-ai-scripts/06-cicd-local-runner.py`. Every time you write code, you do not build the code. Never run all tests locally or full build unless explicitly asked by the user. Only do targeted file checks.
- **R2 Targeted Checks Only.** Run only fast, file-scoped checks on specifically modified files. A check scanning 0 files is a **FAIL**.
- **R3 Evidence or It Did Not Happen.** Every `DONE`, `PASS`, or "verified" claim MUST cite a concrete file path, git diffstat, or command exit code (`exit 0`). Vague assurances are auto-rejected.
- **R4 Never Invent Commands, Flags, or Paths.** Verify commands with a harmless call (`gitmap lf readme.md`), not `--help`. Use documented fallbacks and log in ledger.
- **R5 Mandatory Subagents (`invoke_subagent`), Concurrency Scaling & Immediate Dynamic Respawn.** Spawning subagents via `invoke_subagent` (`A = 3`, `H = 2`, System Concurrency Capacity = 6 concurrent operations) is an **ABSOLUTE MUST** (`research` for discovery in Phase 1, `self` for spec authoring and code edits in Phase 2/3). The lead agent is STRICTLY FORBIDDEN from executing all reads or edits solo. Keep concurrency pinned at 3 active subagents.
- **R6 One Owner Per File (Disjoint Bounding Boxes).** Within every worker wave, each file has exactly one owner. Shared indexes belong exclusively to lead.
- **R7 Git Safety & Isolation (Worker Git Ban).** Subagents NEVER run git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces, worker git calls create `.git/index.lock` collisions that crash parallel agents.
- **R8/R9 Atomic Commit & Push via GitMap (TOTAL BAN ON RAW GIT COMMITS).** The run ends with ONE GitMap call: `gitmap cpf "<module> - <feature summary>"` or `gitmap cpb "<module> - <fix summary>"`. Use hyphen format `-`; NEVER include colons inside the message argument. Stage, commit, and push upstream atomically.
- **R10 Zero Unauthorized Releases.** Never bump versions, edit `version.json`, update changelogs, or trigger release scripts unless user explicitly requested release.
- **R11 Strict Relative Git Paths & Lowercase Hygiene.** Strict ban on absolute paths (`C:\...`, `/home/...`) and `file:///` URIs. Only add relative paths starting from git root. All filenames, documentation, and specs strictly lowercase.
- **R12 No Polling / Immediate Turn Yielding.** When calling `invoke_subagent`, make it the sole tool action at turn end, print progress line, and **STOP CALLING TOOLS**. Never poll in a loop.
- **R13 Two-Strike Retry Cap & Anti-Looping.** A tool failing twice: worker replies `STATUS: BLOCKED` with exact error and stops. Lead takes over. Subtask failing two remediation rounds is marked `FAILED` with RCA.
- **R14 100% Ambiguity & Decision Boundaries.** Non-blocking: choose conservative option, log in ledger `Assumptions:`, proceed. Blocking: `ask_question` once, log in `.ai-memory/ambiguous-questions/01-new-ambiguity/`, continue unblocked tasks.
- **R15 Zero Generated Artifacts Committed.** Never commit build caches, logs, temp scripts, or newly generated binaries unless repository already tracked them.
- **R16 Zero Secrets in Standard Repos.** Never write credentials, tokens, passwords, or `.env` contents into tracked files. Offload secrets to `repo-secrets` in work directory via `gitmap rs`.

---

## GitMap High-Speed Command Primacy (Run Everything Faster)

GitMap is your **PRIMARY** acceleration engine. Banned ALL variants of `Get-ChildItem` (with or without flags/parameters), `dir`, `ls` in shell, `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Execute all searches and operations through GitMap:

| Operation | Primary GitMap Command | High-Speed Alias | Purpose & Advantage |
| :--- | :--- | :--- | :--- |
| **Live Streaming Search** | `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` | `gitmap aum grep "<pat>"` | Multi-core streaming live file search (replaces `Select-String`, `git grep`) |
| **Indexed Symbol Search** | `gitmap search "<query>" [--limit <n>]` | `gitmap search` | Instant SQLite cached keyword/symbol search across indexed repos |
| **Wildcard File Search** | `gitmap find "<pattern>" [-ext <ext>]` | `gitmap f "<pat>"` | Index-accelerated multi-core glob filename finder |
| **Directory Inventory** | `gitmap lf <dir>` / `gitmap find "*.*" <dir>` | `gitmap lf [pat]` | Instant indexed repository inventory (replaces `Get-ChildItem`, `dir`, `ls`) |
| **Substring File Search** | `gitmap find-files-any "<substring>"` | `gitmap ffa "<str>"` | High-speed partial filename matcher |
| **Stream File** | `gitmap cat <filepath>` | `gitmap cat` | Direct zero-disk memory streaming to stdout (replaces `cat`, `type`, `Get-Content`) |
| **PowerShell Runner** | `gitmap pwsh "<command>"` | `gitmap ps "<cmd>"` | High-speed PowerShell execution with `-NoProfile` |
| **Bash Runner** | `gitmap bash "<command>"` | `gitmap sh "<cmd>"` | Standard cross-platform Bash command execution |
| **Offload Secrets** | `gitmap rs file <filepath>` / `folder` / `text` | `gitmap rs` | Auto-commits into `repo-secrets` in work directory |
| **Offload Scripts** | `gitmap rc file <file.ps1>` / `text` | `gitmap rc` | Auto-commits reusable scripts into `repo-cache` |
| **Atomic Commits** | `gitmap cpf "<module> - <msg>"` (Feature) / `cpb` (Bug) | `gitmap cpf` | Stages, formats with prefix, and pushes atomically. Format with hyphen `-`. |
| **Agent Task Engine** | `gitmap agent task <init\|add\|claim\|complete\|fail\|status>` | `gitmap agent` | Sub-millisecond compiled Go SQLite task manager for multi-agent workflows |
| **Heatmap Action Log** | `gitmap heatmap log --task <task> --action <ACTION> --file <path> --state <STATE> --expected "<goal>"` | `gitmap heatmap log` | Concurrency-safe cognitive action & reasoning telemetry into `.ai-memory/heatmap.db` |
| **Contention & Leases** | `gitmap agent heatmap [--json]` / `gitmap heatmap status` | `gitmap agent heatmap` | Pre-edit file lease conflict & contention hotspot analysis across active agents |
| **RC Heatmap Suite** | `gitmap rc heatmap <status\|tables\|query\|view\|tree\|live\|sync>` | `gitmap rc heatmap` | Unified RC heatmap commands with automatic vault auto-provisioning & roadmap linkage |

---

## Additional Instructions

- [/plan](slashCommand;plan) first before doing work for improved work or less friction, faster work and more pure and effective work.
- Always preserve strictly relative Git paths (`01-prompts/...`, `02-spec/...`, `.ai-memory/...`).
- Ensure every run records completed improvements and queues the subsequent 100 UI improvements in `.ai-memory/plans/pending/XX-next-100-ui-improvements.md`.

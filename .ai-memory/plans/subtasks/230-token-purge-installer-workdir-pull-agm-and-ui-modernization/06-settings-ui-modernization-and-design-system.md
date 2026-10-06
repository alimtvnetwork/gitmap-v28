# Subtask 06: Settings UI Modernization and Design System Alignment

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md)  
> **Status:** `QUEUED`  
> **Target Subsystems:**  
> - `src/pages/Settings.tsx`  
> - `src/components/settings/`  
> - `cli/cmdui/ui_assets.go`  
> - `cli/cmdui/ui_server.go`  
> - `02-spec/07-design-system/`  

---

## 1. Technical Objective

Overhaul the Settings user interfaces across both the React application (`src/pages/Settings.tsx`) and the embedded Go Web UI (`cli/cmdui/ui_assets.go` `#tab-settings`) to strictly adhere to the **4-Plane Neutral Depth Hierarchy** and **HSL token architecture** defined in `02-spec/07-design-system/18-dark-mode-and-materiality.md`. Eliminate amateur dark-mode anti-patterns (pure black `#000000`, purple-blue gradients, low-contrast labels, and flat cards), implement bidirectional synchronization with `/api/settings`, and compile an authoritative reference list of modern UI/UX and CSS3 repositories.

---

## 2. File Modification Inventory

| File | Action | Responsibilities |
|:---|:---|:---|
| `src/pages/Settings.tsx` | **Modify** | Restructure layout into semantic elevation cards (Plane 2) with hairline top-light borders, responsive categories, and positive boolean controls. |
| `src/components/settings/StorageTempCard.tsx` | **Modify** | Modernize form inputs with HSL tokens, subtle focus rings, and validation feedback. |
| `src/components/settings/PipelineMonitorCard.tsx` | **Modify** | Align polling controls with tokenized slider and stepper styles. |
| `src/components/settings/TerminalThemeCard.tsx` | **Modify** | Add theme preset preview chips using HSL color ramps. |
| `cli/cmdui/ui_assets.go` | **Modify** | Refactor embedded `#tab-settings` markup and CSS to use CSS custom properties (`var(--card)`, `var(--border)`, `var(--primary)`). |
| `cli/cmdui/ui_server.go` | **Modify** | Validate `/api/settings` serialization to guarantee zero attribute drops during updates. |

---

## 3. Step-by-Step Implementation Plan

### Step 3.1: Enforce 4-Plane Elevation Tokens
1. Verify token availability in `src/index.css`:
   - Canvas (Plane 0): `--background: 230 25% 8%;`
   - Raised (Plane 1): `--muted: 230 18% 18%;`
   - Surface (Plane 2): `--card: 230 20% 12%;`
   - Elevated (Plane 3): `--popover: 230 20% 16%;`
2. Configure hairline borders: `border: 1px solid hsl(var(--border) / 0.4);` with subtle top highlight.

### Step 3.2: React Settings Page Modernization (`src/pages/Settings.tsx`)
1. Organize settings into categorized panels:
   - **Storage & Temp:** Custom `.ai-memory/temp` path with browse/reset actions.
   - **Pipeline Telemetry:** Polling interval, failure tree depth, and error limit toggles.
   - **Antigravity Multi-Instance:** Default instance selection, language server discovery toggle.
   - **Terminal & Styling:** Theme selection and interactive terminal font sizes.
2. Enforce positive booleans on all state hooks:
   - `isAutoSyncActive`, `isTelemetryEnabled`, `hasLocalFallback`.
3. Wire bidirectional save:
   - On save, issue `POST /api/settings` with JSON payload.
   - Provide optimistic state update and toast notifications.

### Step 3.3: Embedded Web UI Modernization (`cli/cmdui/ui_assets.go`)
1. In `#tab-settings` container:
   - Replace legacy table layouts with CSS grid (`display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 1.5rem;`).
   - Wrap control groups inside cards with `.settings-card` class referencing `var(--card)`.
   - Update input fields to use `var(--input)` and `var(--foreground)`.
2. Ensure buttons have +4% hover lightness and tactile active states.

### Step 3.4: Reference UI/UX & CSS3 Repositories Index
Document and verify the curated list of reference repositories within the specification and UI guide:
- `shadcn-ui/ui`: Copy-paste accessible components and token architecture.
- `radix-ui/primitives`: Unstyled accessible headless component primitives.
- `mannupaaji/aceternity-ui`: Modern interactive canvas animations and glowing borders.
- `magicuidesign/magicui`: High-craft animations, retro grids, and particle UI.
- `tremorlabs/tremor`: Clean data telemetry cards and analytical dashboards.
- `lucide-icons/lucide`: Consistent, clean vector iconography.
- `moderncss.dev`: Modern CSS3 layout and selector architectures.

---

## 4. Verification Commands & Expected Output

```bash
# 1. Verify TypeScript types and build hygiene
npx tsc --noEmit

# 2. Verify embedded Web UI server starts cleanly
go test -v -run TestSettingsAPI ./cli/cmdui/

# 3. Inspect settings endpoint JSON serialization
curl -s http://127.0.0.1:42120/api/settings | jq .
```

---

## 5. Acceptance Criteria

- [ ] Settings page renders without pure black `#000000` or violet gradient soup.
- [ ] 4-plane depth hierarchy visually distinguishable between background, cards, and popovers.
- [ ] Both React and Go embedded Web UI settings forms persist changes reliably.
- [ ] Positive boolean naming conventions enforced repo-wide for settings state.
- [ ] Curated UI/UX repository index published in component spec.

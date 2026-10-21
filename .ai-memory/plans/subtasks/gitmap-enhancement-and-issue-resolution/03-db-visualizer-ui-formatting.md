# Subtask 03: Database Visualizer UI Binary Formatting & Inspection Drawer

## Status: READY
- **Subtask Slug**: `03-db-visualizer-ui-formatting`
- **Parent Plan**: `.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md`
- **Architecture Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md`
- **Owner**: Worker 01 (Database Operations & Cache Engineer)
- **Target Subsystems**: `cli/cmddb` and `src/pages`
- **Exclusive Owned Files**:
  - `cli/cmddb/rest_visualizer_html.go`
  - `cli/cmddb/rest_visualizer_js1.go`
  - `cli/cmddb/rest_visualizer_js2.go`
  - `src/pages/Database.tsx`

---

## 1. Context & Objectives

The GitMap Database Visualizer provides web-based table browsing and query execution. However, inspection of tables containing binary hashes, Protobuf metadata, serialized envelopes, or long text currently suffers from critical usability flaws:
1. **Unformatted / Corrupted Binary Output**: Binary BLOBs currently render as raw strings, causing Unicode replacement characters (`\uFFFD` / black diamonds) or unreadable text.
2. **Missing Inspection Drawer**: Developers have no side-by-side tool to inspect binary blobs as structured Hex Dumps, formatted JSON trees, or decoded strings.
3. **Unbounded Pagination**: Clicking "Next" increments page count indefinitely beyond available rows, displaying blank tables.
4. **Disjoint Visualizer UIs**: The embedded web visualizer (`rest_visualizer_*.go`) and the React desktop UI (`src/pages/Database.tsx`) must share identical formatting semantics, badges, and inspection drawer capabilities.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Update Visualizer HTML Structure & CSS (`cli/cmddb/rest_visualizer_html.go`)

1. **Add CSS Styles for Cell Badges & Drawer**:
   - `.blob-badge`: Purple/blue monospace badge (`background: rgba(139, 92, 246, 0.2); color: #a78bfa; border: 1px solid rgba(139, 92, 246, 0.4); padding: 2px 6px; border-radius: 4px; cursor: pointer; font-size: 11px;`).
   - `.json-badge`: Cyan monospace badge (`background: rgba(6, 182, 212, 0.2); color: #22d3ee; border: 1px solid rgba(6, 182, 212, 0.4); padding: 2px 6px; border-radius: 4px; cursor: pointer; font-size: 11px;`).
   - `.bool-badge-true`: Green badge (`color: #10b981; font-weight: 600;`).
   - `.bool-badge-false`: Muted gray badge (`color: #64748b; font-weight: 600;`).
   - `.cell-text-truncated`: Ellipsis overflow with pointer cursor and title tooltip.
   - `.inspection-drawer`: Sliding side drawer:
     ```css
     .inspection-drawer {
       position: fixed; top: 0; right: 0; width: 520px; max-width: 90vw; height: 100vh;
       background: #0c111d; border-left: 1px solid #1e293b; z-index: 250;
       display: flex; flex-direction: column; transform: translateX(100%);
       transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
       box-shadow: -8px 0 24px rgba(0, 0, 0, 0.6);
     }
     .inspection-drawer.open { transform: translateX(0); }
     .drawer-backdrop {
       position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 240; display: none;
     }
     .drawer-backdrop.open { display: block; }
     ```

2. **Add Inspection Drawer DOM Elements**:
   - Drawer Header: Cell title, column name, data type, and close icon (`&times;`).
   - Mode Tabs: `Hex Dump`, `Formatted JSON`, `String Decoder`.
   - Content Panes:
     - `#drawer-hex-content`: `<pre class="hex-dump-pre">` for 16-byte aligned hex view.
     - `#drawer-json-content`: `<pre class="json-dump-pre">` for formatted syntax-colored JSON.
     - `#drawer-text-content`: Decoded text area with character counts.
   - Drawer Footer: Action buttons (`Copy to Clipboard`, `Download as .bin`, `Close`).

---

### Step 2: Implement Smart Cell Rendering & Pagination Bounds (`cli/cmddb/rest_visualizer_js1.go`)

1. **Implement `formatCellContent(val, colName)`**:
   - Check if `val === null || val === undefined`: return `<span class="null-val">NULL</span>`.
   - Check for BLOB object (`val.__blob === true`):
     - Return `<button class="blob-badge" onclick="window.openBlobDrawer(...)">[BLOB ` + val.size + `B]</button>`.
   - Check for Booleans (`val === true || val === false || (colName.startsWith('is') && (val === 1 || val === 0))`):
     - Return `<span class="bool-badge-` + (val ? 'true' : 'false') + `">` + (val ? '✓ true' : '✗ false') + `</span>`.
   - Check for JSON strings (starts with `{` and ends with `}`, or starts with `[` and ends with `]`):
     - Attempt `JSON.parse(val)`. On success, return `<button class="json-badge" onclick="window.openJsonDrawer(...)">[JSON: ` + Object.keys(parsed).length + ` keys]</button>`.
   - Check for long strings (`String(val).length > 80`):
     - Truncate to 80 characters with ellipsis. Add title attribute with full content. Add click-to-copy handler.

2. **Fix Pagination Boundary Controls**:
   - In `renderGrid(env)`:
     - Extract `totalRecords = (env.Attributes && env.Attributes.TotalRecords) ? env.Attributes.TotalRecords : (rows.length ? state.page * state.limit : 0)`.
     - Calculate `totalPages = Math.max(1, Math.ceil(totalRecords / state.limit))`.
     - Update UI buttons:
       ```javascript
       var btnPrev = document.getElementById('btn-prev');
       var btnNext = document.getElementById('btn-next');
       btnPrev.disabled = (state.page <= 1);
       btnNext.disabled = (state.page >= totalPages);
       btnPrev.style.opacity = state.page <= 1 ? '0.4' : '1.0';
       btnNext.style.opacity = state.page >= totalPages ? '0.4' : '1.0';
       document.getElementById('page-info').textContent = 'Page ' + state.page + ' of ' + totalPages + ' (' + totalRecords + ' records)';
       ```
     - Guard click handlers:
       - `btn-next`: `if (state.page < totalPages) { state.page++; loadGridData(); }`.
       - `btn-prev`: `if (state.page > 1) { state.page--; loadGridData(); }`.

---

### Step 3: Implement Drawer Controllers & Encoders (`cli/cmddb/rest_visualizer_js2.go`)

1. **Implement `renderHexDump(hexStr)`**:
   - Parse hex string into 16-byte chunks.
   - Format each row: `Offset (8 hex) | 16 Hex Bytes (separated by space) | ASCII Printable (dots for non-printable)`.
   - Example line: `00000000  47 49 54 4d 41 50 00 01  02 03 04 05 06 07 08 09  |GITMAP..........|`.
2. **Implement `renderJsonTree(jsonObj)`**:
   - Format with 2-space indentation.
   - Apply syntax coloring spans for keys, strings, numbers, and booleans.
3. **Implement Drawer Interactivity**:
   - `openInspectionDrawer(data, type, title)`: Slide drawer open, activate selected tab, populate pane.
   - `closeInspectionDrawer()`: Slide drawer closed and hide backdrop.
   - `copyDrawerContent()`: Copy active pane text to `navigator.clipboard`.
   - `downloadDrawerBlob(hexStr, filename)`: Convert hex to `Uint8Array`, construct `Blob(['application/octet-stream'])`, and trigger download.

---

### Step 4: Author React Table Browser & Inspection Drawer (`src/pages/Database.tsx`)

1. **Component Architecture**:
   - Built using React 18, Tailwind CSS, Lucide Icons, and `@tanstack/react-query`.
   - Header with Database selector, table switcher, query execution latency, and refresh button.
   - Split layout:
     - Left Sidebar: Tables and Views with record count badges and filter input.
     - Center Grid: Dynamic columns, sticky headers, cell classification (BLOB badge, JSON badge, boolean indicator, truncated string).
     - Right Slide-out Drawer: Built with Radix Sheet or Tailwind transition drawer for Hex / JSON / String inspection.
2. **Features Parity**:
   - Identical 16-byte hex dump formatting.
   - Interactive collapsible JSON viewer.
   - Copy to clipboard with toast notification.
   - Bounded pagination bar (`Page X of Y`).
   - Download raw binary action.

---

## 3. Verification & Quality Gates

- **Visual Quality**: Zero `\uFFFD` black diamonds rendered when opening tables with binary columns.
- **Drawer Responsiveness**: Inspection drawer opens in < 50ms with smooth 60fps CSS transform transitions.
- **Pagination Safety**: Next button disables when on the last page; Prev button disables on page 1.
- **Zero Nested Ifs & Function Limits**: Helper functions <= 15 lines in Go backend files.
- **Relative Paths Only**: Strictly relative paths in all imports, scripts, and documentation.

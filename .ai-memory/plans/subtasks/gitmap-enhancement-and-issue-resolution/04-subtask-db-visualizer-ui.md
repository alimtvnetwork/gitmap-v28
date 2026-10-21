# Subtask 04: Database Web Visualizer UI Modernization & Safe BLOB Rendering

## Status: READY
- **Parent Plan**: `.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md`
- **Specification Reference**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md`
- **Owner**: Worker 04 (Visualizer UI & Backend Formatting Engineer)
- **Target Subsystem**: `cli/cmddb`
- **Owned Files**:
  - `cli/cmddb/rest_handlers_crud.go`
  - `cli/cmddb/rest_visualizer_html.go`
  - `cli/cmddb/rest_visualizer_js1.go`
  - `cli/cmddb/rest_visualizer_js2.go`
  - `cli/cmddb/rest_visualizer_test.go` (new file)

---

## 1. Context & Objectives

1. **Eliminate Corrupted Black Diamond Replacement Glyphs (`` / `\uFFFD`)**:
   - In `cli/cmddb/rest_handlers_crud.go` (`scanAllRows`), raw byte slices (`[]byte`) are unconditionally cast to `string(b)`.
   - When binary data (Protobuf payloads, cryptographic hashes, BLOB columns) contains invalid UTF-8 byte sequences, web browsers and JSON encoders substitute the Unicode replacement character `\uFFFD`, corrupting data presentation.
   - Implement backend binary detection (`isBinaryPayload`) using `utf8.Valid(b)` and byte ratio checks, encoding binary data into structured tagged envelopes (`{ "__blob": true, "size": N, "hex": "...", "preview": "..." }`).
2. **Frontend Smart Cell Formatting**:
   - Modernize cell rendering in `cli/cmddb/rest_visualizer_js1.go` and `rest_visualizer_js2.go`.
   - **BLOB Preview & Modal**: Render binary fields as interactive `[BLOB 128B]` badges. Clicking opens a Hex Viewer modal with a dual-pane 16-byte hex dump and ASCII printables.
   - **Collapsible JSON Tree**: Detect valid JSON objects and arrays. Render collapsible badges `[JSON: 4 keys]` that expand to an inline interactive tree.
   - **Boolean Badges**: Style booleans as positive green `✓ true` or muted gray `✗ false`.
   - **Text Truncation**: Truncate strings exceeding 100 characters with ellipsis and click-to-expand.
3. **Pagination Boundary Controls**:
   - Calculate `totalPages = Math.max(1, Math.ceil(totalRecords / state.limit))`.
   - Prevent navigation past `totalPages` on `btn-next` and before page 1 on `btn-prev`.
   - Visually disable buttons when boundary is reached.
   - Display `Page X of Y (Total Z records)`.
4. **Code Quality & Coding Guidelines**:
   - All Go functions must be <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Positive boolean naming only (`is*` and `has*`).
   - Strictly relative paths and lowercase filenames.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Refactor `scanAllRows` in `cli/cmddb/rest_handlers_crud.go`
1. Implement `isBinaryPayload(b []byte) bool`:
   - Guard 1: Verify `utf8.Valid(b)`. If false, return `true`.
   - Guard 2: Count null bytes (`\x00`). If `hasNullBytes(b)`, return `true`.
   - Guard 3: Calculate control character ratio. If non-printable characters exceed 5%, return `true`.
   - Return `false` (safe UTF-8 text).
2. Implement helper `formatBinaryEnvelope(b []byte) map[string]any`:
   - Return tagged map:
     ```go
     map[string]any{
         "__blob":  true,
         "size":    len(b),
         "hex":     hex.EncodeToString(b),
         "preview": fmt.Sprintf("[BLOB %s]", formatByteSize(len(b))),
     }
     ```
3. Update `scanAllRows`:
   - Replace direct `row[c] = string(b)` with:
     ```go
     if isBinaryPayload(b) {
         row[c] = formatBinaryEnvelope(b)
     } else {
         row[c] = string(b)
     }
     ```

### Step 2.2: Add Hex Viewer Modal in `cli/cmddb/rest_visualizer_html.go`
1. Add modal container `#hex-viewer-modal`:
   ```html
   <div id="hex-viewer-modal" class="modal-backdrop" style="display:none;">
     <div class="modal-card hex-modal-card">
       <div class="modal-header">
         <span id="hex-modal-title">BLOB Hex Inspector</span>
         <span id="hex-modal-size" class="item-badge">0 B</span>
         <button class="btn-close" id="btn-hex-modal-close">&times;</button>
       </div>
       <div class="hex-viewer-body">
         <pre id="hex-dump-output" class="hex-dump-pre"></pre>
       </div>
       <div class="modal-footer">
         <button class="btn-sm" id="btn-hex-copy">Copy Hex</button>
         <button class="btn-sm" id="btn-hex-download">Download Binary</button>
       </div>
     </div>
   </div>
   ```
2. Add CSS styles for `.hex-modal-card`, `.hex-dump-pre`, `.badge-blob`, `.badge-bool-true`, `.badge-bool-false`, and `.json-tree-node`.

### Step 2.3: Modernize Cell Rendering in `cli/cmddb/rest_visualizer_js1.go`
1. Implement `formatTableCell(val, colName)`:
   - Check 1: If `val === null || val === undefined`, return `<span class="val-null">NULL</span>`.
   - Check 2: If `val && val.__blob`, return `<span class="badge-blob" onclick="openHexViewer('${val.hex}', ${val.size}, '${colName}')">${val.preview}</span>`.
   - Check 3: If `typeof val === 'boolean'`, return boolean badge `<span class="badge-bool-${val}">${val ? '✓ true' : '✗ false'}</span>`.
   - Check 4: If `typeof val === 'object'`, return collapsible JSON tree toggle.
   - Check 5: If string length > 100, return truncated span with title attribute and click-to-expand.
2. Update `renderGrid` in `rest_visualizer_js1.go`:
   - Apply `td.innerHTML = formatTableCell(r[c], c)`.
3. Update pagination boundary handling in `renderGrid`:
   - Read `totalRecords = (env.Attributes && env.Attributes.TotalRecords) ? env.Attributes.TotalRecords : rows.length`.
   - Compute `totalPages = Math.max(1, Math.ceil(totalRecords / state.limit))`.
   - Update prev/next button disabled attributes:
     - `document.getElementById('btn-prev').disabled = (state.page <= 1)`.
     - `document.getElementById('btn-next').disabled = (state.page >= totalPages)`.
   - Set `#page-info` text: `Showing ${rows.length} of ${totalRecords} (Page ${state.page} of ${totalPages})`.

### Step 2.4: Implement Hex Dump & Modal Logic in `cli/cmddb/rest_visualizer_js2.go`
1. Implement `renderHexDump(hexStr)`:
   - Break hex string into 16-byte chunks (32 hex characters).
   - Format each row: `Offset (8 hex digits): Hex Bytes (spaced) | ASCII representation`.
   - Replace unprintable ASCII characters with `.`.
2. Implement `openHexViewer(hexStr, size, colName)`:
   - Set modal title and size badge.
   - Populate `#hex-dump-output` with formatted hex dump.
   - Display modal with smooth transition.
3. Implement `copyHexToClipboard` and `downloadBinaryBlob`.
4. Implement boundary-guarded pagination listeners:
   - `btn-next`: `if (state.page < totalPages) { state.page++; loadGridData(); }`.
   - `btn-prev`: `if (state.page > 1) { state.page--; loadGridData(); }`.

### Step 2.5: Implement Tests in `cli/cmddb/rest_visualizer_test.go`
1. Test `isBinaryPayload`:
   - Valid ASCII string -> `false`.
   - Valid UTF-8 multi-byte string (e.g. Japanese, emojis) -> `false`.
   - Raw binary with null bytes (`\x00`) -> `true`.
   - Invalid UTF-8 sequence (`\xff\xfe\xfd`) -> `true`.
2. Test `scanAllRows` with table containing BLOB and TEXT columns:
   - Verify BLOB column produces `{ "__blob": true, ... }`.
   - Verify TEXT column produces clean string without `\uFFFD`.
3. Test REST visualizer HTML and JS bundling:
   - Verify assets compile cleanly and endpoints serve updated scripts.

---

## 3. Verification & Quality Gates

- **Function Length**: <= 15 lines per function (<= 8 lines preferred).
- **Branching**: Zero nested `if` statements. Guard clauses only.
- **Naming**: Strict positive booleans (`is*` and `has*`).
- **Data Integrity**: Zero corrupted black diamond glyphs (`\uFFFD`) across all test tables.
- **Relative Path Hygiene**: Strictly relative paths in code, comments, and tests.

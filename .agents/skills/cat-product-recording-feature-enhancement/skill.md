---
name: cat-product-recording-feature-enhancement
description: >-
  Enhanced recording feature capturing UI telemetry, intelligent hover deduplication, ROI drag/resize coordinates, automatic V6/V7 prompt generation, export dropdown (ZIP, Copy V6/V7), direct email, Supabase sync with 450MB/14-day retention, error management modal, and local run speed/terminal lifecycle fixes.
---

# CAT Product Recording Feature Enhancement

## Strategic Overview
Enhance the CAT product recording and bug reporting engine to capture deep UI telemetry (XPath, bounding box, padding, colors, distance, sequence number, visual spring animation) and serialize recordings into structured Markdown using standard V6 and V7 prompt templates. Images are stored locally near the binary in `bug-reports/<seq>-<date>/` and referenced via `[sequence_instance_what](file:;<absolute_path>)` without polluting codebase assets. Provide multi-export options (ZIP, Copy V6, Copy V7), email dispatch, Supabase rolling synchronization (<=450MB limit, 14-day data eviction), top-level error management inspection, direct local execution without redundant zipping, and clean terminal lifecycle teardown.

## Key Directives
1. **Local Binary-Adjacent Storage:** Store all screenshots in `bug-reports/<sequence>-<date>/` adjacent to the executable binary. Never save screenshots inside repository assets.
2. **Prompt Generation (V6 & V7):** Automatically generate complete, copyable prompt templates adhering to GitMap/Letterly V6 and V7 standards with full context and image references.
3. **UI Telemetry & Spring Animation:** Capture element geometry, computed styles, distances, and trigger an animated sequence counter on click.
4. **Intelligent Hover Recording:** Debounce and deduplicate hover captures, restricting captures to clickable elements only (buttons, links, clickable images).
5. **Region of Interest (ROI) Drag Telemetry:** Record start/end coordinates, pixel movement, XPath, and paired before/after screenshots for drag and resize actions.
6. **Multi-Export Dropdown:** Provide a clean dropdown for ZIP export, Copy as V6 Prompt, and Copy as V7 Prompt.
7. **Email Integration:** Support direct email dispatch with markdown report and screenshot attachments.
8. **Supabase Quota Governance:** Enforce 450MB quota cap and auto-prune data older than 14 days.
9. **Error Management Visibility:** Add error manager modal and top navigation error button with one-click markdown copy.
10. **Local Runner Optimization:** Eliminate redundant zipping during local development runs and fix terminal breakage upon termination.

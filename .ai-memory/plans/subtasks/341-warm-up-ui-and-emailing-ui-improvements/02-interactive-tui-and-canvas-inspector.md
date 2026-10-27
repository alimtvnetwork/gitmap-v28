# Subtask 02: Interactive TUI Operator Studio, Live Visual Canvas & Render Pipeline

## Status: READY
- **Subtask Slug:** `02-interactive-tui-and-canvas-inspector`
- **Parent Plan:** `.ai-memory/plans/pending/341-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/341-warm-up-ui-and-emailing-ui-improvements/02-improvements-431-to-460-interactive-tui-and-canvas-inspector.md`
- **Target Range:** Improvements 431–460 (Clusters 44–46)
- **Target Subsystem:** `pkg/cli/tui/`, `pkg/web/canvas/`, `pkg/mail/telemetry/`
- **Exclusive Owned Files:**
  - `pkg/cli/tui/canvas_studio.go`
  - `pkg/cli/tui/queue_monitor.go`
  - `pkg/web/canvas/inspector.go`
  - `pkg/web/canvas/dom_renderer.go`
  - `pkg/mail/telemetry/render_profiler.go`
  - `cli/cmdwarmup/canvas_cmd.go`

---

## 1. Context & Objectives

This subtask delivers the interactive terminal operator studio, visual template inspector, and streaming render pipeline telemetry:
1. **Interactive TUI Operator Studio (Cluster 44, Items 431–440):** Full-screen Lipgloss terminal dashboard offering real-time queue observation, pause/resume hotkeys, and ISP delivery latency monitors.
2. **Visual Live Canvas Inspector (Cluster 45, Items 441–450):** Embedded web inspector rendering responsive HTML email layouts, previewing CSS box models, and testing dark-mode inversion.
3. **MIME Render Pipeline Telemetry (Cluster 46, Items 451–460):** Sub-millisecond profiling of MIME body generation, attachment streaming buffers, and memory allocations per send operation.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isStudioActive`, `hasRenderCompleted`, `isCanvasInspected`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Lipgloss Operator Studio TUI
- Implement multi-pane terminal studio with live event bus listeners in `pkg/cli/tui/canvas_studio.go`.
- Build real-time queue inspector and active connection monitors in `pkg/cli/tui/queue_monitor.go`.

### Step 2: Embedded Live Visual Canvas Inspector
- Construct zero-CDN DOM inspector and box model visualizer in `pkg/web/canvas/inspector.go`.
- Implement responsive viewport preview renderer in `pkg/web/canvas/dom_renderer.go`.

### Step 3: High-Speed MIME Pipeline Profiler
- Build microsecond benchmark profiler in `pkg/mail/telemetry/render_profiler.go`.
- Integrate visual studio commands into `cli/cmdwarmup/canvas_cmd.go`.

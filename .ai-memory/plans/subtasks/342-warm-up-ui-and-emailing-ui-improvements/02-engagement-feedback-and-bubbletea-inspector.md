# Subtask 02: Engagement Feedback & Bubbletea Inspector Studio

## Status: READY
- **Subtask Slug:** `02-engagement-feedback-and-bubbletea-inspector`
- **Parent Plan:** `.ai-memory/plans/pending/342-warm-up-ui-and-emailing-ui-improvements.md`
- **Architecture Spec:** `02-spec/21-app/342-warm-up-ui-and-emailing-ui-improvements/02-improvements-531-to-560-engagement-and-tui-studio.md`
- **Target Range:** Improvements 531–560 (Clusters 54–56)
- **Target Subsystem:** `pkg/cli/tui/`, `pkg/mail/engagement/`, `pkg/web/canvas/`
- **Exclusive Owned Files:**
  - `pkg/cli/tui/bubbletea_inspector.go`
  - `pkg/cli/tui/inspector_model.go`
  - `pkg/cli/tui/viewport_renderer.go`
  - `pkg/mail/engagement/seed_tracker.go`
  - `pkg/mail/engagement/reply_analyzer.go`
  - `pkg/mail/engagement/inbox_placement.go`
  - `pkg/web/canvas/canvas_debugger.go`
  - `cli/cmdwarmup/inspector_cmd.go`

---

## 1. Context & Objectives

This subtask delivers the interactive Bubbletea terminal studio, real-time engagement telemetry, and visual template canvas debugger:
1. **Interactive Bubbletea Operator Studio (Cluster 54, Items 531–540):** Full-screen terminal dashboard featuring split-pane inspectors, live dispatch queues, dynamic log viewports, and keyboard event handlers.
2. **Seed Account Engagement Telemetry (Cluster 55, Items 541–550):** Automatic IMAP polling of warmup seed accounts, inbox placement verification, open rate tracking, and automated positive reply generation.
3. **Canvas Template Debugger & Rendering Profiler (Cluster 56, Items 551–560):** Microsecond rendering benchmarks, CSS inlining validation, dark mode rendering simulation, and responsive layout testing.
4. **Architectural Invariants:**
   - Functions <= 15 lines (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Strict positive boolean naming (`isInspecting`, `hasDelivered`, `isInboxPlacement`).
   - Typed error returns via `*apperror.AppError` (zero exits).
   - Zero external CDN dependencies (100% embedded assets).
   - 100% relative repository paths.

---

## 2. Step-by-Step Implementation Breakdown

### Step 1: Interactive Bubbletea Terminal Inspector
- Construct Bubbletea Elm-architecture model in `pkg/cli/tui/inspector_model.go`.
- Implement responsive Lipgloss terminal layout in `pkg/cli/tui/bubbletea_inspector.go`.
- Build viewport scroller for real-time dispatch streams in `pkg/cli/tui/viewport_renderer.go`.

### Step 2: Seed Account Engagement Feedback Engine
- Implement IMAP inbox poll worker in `pkg/mail/engagement/seed_tracker.go`.
- Build sentiment extraction and reply parser in `pkg/mail/engagement/reply_analyzer.go`.
- Calculate inbox vs spam vs promo placement ratios in `pkg/mail/engagement/inbox_placement.go`.

### Step 3: Visual Canvas Debugger & CLI Integration
- Construct local web canvas preview debugger in `pkg/web/canvas/canvas_debugger.go`.
- Wire CLI interactive inspector command into `cli/cmdwarmup/inspector_cmd.go`.

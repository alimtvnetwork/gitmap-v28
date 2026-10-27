# Specification 337: Warm-Up UI and Emailing UI Improvements (Overview & Architecture)

- **Sequence:** 337
- **Slug:** `337-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## 1. Executive Summary & System Vision

Modern outbound email infrastructure requires rigorous warmup protocols, reputation tracking, and deliverability safeguarding before mass or automated transactional sending can succeed. The **Warm-Up and Emailing UI Engine** within GitMap provides a unified, zero-CDN, enterprise-grade control plane inspired by modern autonomous manager dashboards (such as Google Antigravity Manager).

It delivers end-to-end visualization, automated ramp scheduling, real-time seed list interaction simulation, and ISP inbox/spam placement gauges.

```
                          ┌──────────────────────────────────────────────┐
                          │    GitMap Warm-Up & Outbound Control Plane   │
                          └──────────────────────┬───────────────────────┘
                                                 │
                   ┌─────────────────────────────┼─────────────────────────────┐
                   ▼                             ▼                             ▼
       ┌───────────────────────┐     ┌───────────────────────┐     ┌───────────────────────┐
       │   Cluster 1 (01-10)   │     │   Cluster 2 (11-20)   │     │   Cluster 3 (21-30)   │
       │   Ramp Curves &       │     │   Deliverability      │     │   Seed List & Reply   │
       │   Visual Timelines    │     │   Health & Gauges     │     │   Simulation Physics  │
       └───────────┬───────────┘     └───────────┬───────────┘     └───────────┬───────────┘
                   │                             │                             │
                   └─────────────────────────────┼─────────────────────────────┘
                                                 │
                                                 ▼
                         ┌──────────────────────────────────────────────┐
                         │      Core Architecture & Shared Engine       │
                         │   - Split-DB SQLite Storage (Positive State) │
                         │   - Error Return Sovereignty (No Leaf Exits) │
                         │   - Zero External CDN / Pure SVG & Vanilla   │
                         │   - CLI Counterparts (`gitmap warmup/mail`)  │
                         └──────────────────────────────────────────────┘
```

The system is partitioned into three core operational clusters:
1. **Cluster 1: Warm-Up Ramp Curves & Visual Schedule Timeline (Improvements 01–10):** Dynamic mathematical progression profiles (Linear, Exponential, S-Curve, Fibonacci, Wavefront), adaptive throttle physics, timeline heatmaps, timezone dispersion bands, and interactive milestone gates.
2. **Cluster 2: Deliverability Health Gauges & ISP Inbox/Spam Placement (Improvements 11–20):** Tri-color radial SVG dial gauges, multi-ISP placement breakdowns (Google Workspace, Microsoft 365, Yahoo/AOL, Apple iCloud, Private Corporate MX), DNS authentication telemetry (SPF/DKIM/DMARC/BIMI), reputation trend lines, and automated pause circuit-breakers.
3. **Cluster 3: Seed List Scheduler, Reply Simulation & Automated Interaction Physics (Improvements 21–30):** Seed mailbox cluster orchestration, automated open/read velocity simulation, multi-turn reply thread weaving, mark-as-important & unspam physics, inbox category nudging, and sentiment-scored conversational feedback loops.

---

## 2. Inspiration from Antigravity Manager

The user interface and state machine take architectural inspiration from the **Google Antigravity Manager** workspace control plane:
- **Card-Deck Layout & Glassmorphism Tokens:** High-density, dark-mode first terminal and web dashboards utilizing CSS custom properties for color ramps, ambient border glows, and crisp 1px borders.
- **Micro-State Visualization:** Every domain entity displays immediate status indicators (`isActive`, `hasWarning`, `isPaused`, `isHealthy`) with deterministic positive boolean states.
- **Zero-Latency Telemetry & Polling Invariants:** Telemetry events stream via reactive websockets or local SQLite change journals, completely avoiding blocking loops or external telemetry servers.
- **Coupled CLI & GUI Sovereignty:** Every visual widget in the UI has an exact 1:1 CLI counterpart (`gitmap warmup ...` or `gitmap mail ...`), allowing headless server administration, CI/CD pipeline automation, and local desktop dashboarding without architectural divergence.

---

## 3. High-Level Data Flow & State Machine

```
   [User / CLI Dispatcher]
             │
             ▼
   [Command Orchestrator (gitmap warmup/mail)]
             │
             ├──► [Profile & Curve Evaluator] ──► Computes hourly & daily quotas
             │
             ├──► [Seed Node Network] ─────────► Dispatches synthetic peer exchanges
             │
             ├──► [Deliverability Inspector] ──► Queries DNS, IMAP seed boxes, headers
             │
             ▼
   [SQLite Warmup Store (`.gitmap/warmup.db`)]
             │
             ▼
   [Zero-CDN Web / Terminal Dashboard] (Pure SVG, CSS Variables, Native DOM updates)
```

### 3.1 State Transitions
- `Unconfigured` -> `Scheduled` (Profile selected, seed pool mapped)
- `Scheduled` -> `Ramping` (Execution active, volume accelerating)
- `Ramping` -> `Paused` (Automated circuit-breaker triggered by bounce > 3% or spam placement > 10%)
- `Paused` -> `Ramping` (Manual operator resolution or cool-down recovery)
- `Ramping` -> `Completed` (Target daily volume reached and stability threshold maintained for 7 days)
- `Completed` -> `Maintenance` (Background baseline pulse of 10-15% volume to preserve IP/domain warmth)

---

## 4. Design System & CSS Invariants (Zero-CDN)

The UI adheres strictly to the GitMap Zero-CDN Design System. No Google Fonts, unpkg, cdnjs, or third-party web fonts are loaded. All icons are inline SVG vectors.

### 4.1 Color Tokens
```css
:root {
  --bg-primary: #0f172a;
  --bg-secondary: #1e293b;
  --bg-tertiary: #334155;
  --border-subtle: #1e293b;
  --border-active: #3b82f6;
  --text-primary: #f8fafc;
  --text-muted: #94a3b8;
  --accent-blue: #3b82f6;
  --accent-green: #10b981;
  --accent-amber: #f59e0b;
  --accent-red: #ef4444;
  --accent-purple: #8b5cf6;
  --curve-gradient: linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(16, 185, 129, 0.05) 100%);
}
```

### 4.2 Component Architecture
- **Cards:** Border-radius `8px`, background `var(--bg-secondary)`, border `1px solid var(--border-subtle)`.
- **Gauges:** Mathematical inline `<svg viewBox="0 0 100 100">` with `stroke-dasharray` and `stroke-dashoffset` computations.
- **Curves:** Native SVG `<path d="..." fill="none" stroke="var(--accent-blue)" stroke-width="2.5" />` rendered from cubic Bezier points calculated client-side.
- **Tooltips:** Pure CSS `:hover` or lightweight JavaScript hover coordinates without popper.js or external libraries.

---

## 5. Coding & Architectural Invariants

In accordance with `02-spec/02-coding-guidelines/` and GitMap repository rules:
1. **Error Return Sovereignty:** Leaf functions returning `error` or `Result[T]` must return errors directly. No calls to `os.Exit` or CLI exit handlers inside leaf functions.
2. **Boolean Naming Invariant:** All booleans MUST start with `is` or `has` (e.g. `isActive`, `hasSmtpAuth`, `isInboxPlaced`). Zero negative booleans (`isNotValid` is strictly forbidden).
3. **Function Length:** Functions must remain between 8 and 15 lines max.
4. **Conditional Flattening:** Zero nested `if` statements. Guard clauses with early returns only.
5. **Database Architecture:** Split-DB architecture using SQLite with PascalCase tables, camelCase columns, and strongly typed enums (`*Type`).
6. **Path Discipline:** 100% relative paths across documentation, specifications, and code imports.

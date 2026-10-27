# Specification 338: Warm-Up UI and Emailing UI Improvements (Overview & Architecture)

- **Sequence:** 338
- **Slug:** `338-warm-up-ui-and-emailing-ui-improvements`
- **File:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
- **Status:** APPROVED & SPECIFIED
- **Related Specs:**
  - `02-spec/21-app/337-warm-up-ui-and-emailing-ui-improvements/01-architecture-and-overview.md`
  - `02-spec/21-app/01-cli-architecture/readme.md`
  - `02-spec/02-coding-guidelines/`
  - `02-spec/03-error-manage/`
- **Assigned Author:** Spec Author Subagent 1

---

## 1. Executive Summary & Vision

Sequence 338 elevates the GitMap Warm-Up & Outbound Emailing ecosystem into an enterprise-grade autonomous control plane. Expanding beyond foundational scheduling and telemetry (Sequence 337), Sequence 338 introduces:
1. **Advanced Interactive Canvas Physics & Delivery Path Routing** (Cluster 11, Items 101–110): Direct-manipulation force-directed network topology nodes, real-time Bezier spline packet animations, multi-hop SMTP relay tracing, and dynamic throughput load balancing across distributed outbound gateways.
2. **AI Deliverability Co-Pilot & Subject Line Entropy** (Cluster 12, Items 111–120): Localized LLM inference engines analyzing semantic spam triggers, calculating Shannon token entropy, detecting spam trap syntaxes, scoring inbox sentiment, and recommending instantaneous template adaptations.
3. **Multi-Persona Seed Simulation & Conversational Dynamics** (Cluster 13, Items 121–130): Realistic multi-turn dialog simulations among heterogeneous synthetic personas (B2B Enterprise Buyer, Technical Architect, Executive Decision Maker, Consumer Retail), simulating authentic click velocities, thread branching, and inbox folder engagement physics.

```
                          ┌──────────────────────────────────────────────┐
                          │    GitMap Sequence 338 Architecture Matrix   │
                          └──────────────────────┬───────────────────────┘
                                                 │
                   ┌─────────────────────────────┼─────────────────────────────┐
                   ▼                             ▼                             ▼
       ┌───────────────────────┐     ┌───────────────────────┐     ┌───────────────────────┐
       │  Cluster 11 (101-110) │     │  Cluster 12 (111-120) │     │  Cluster 13 (121-130) │
       │  Interactive Node     │     │  AI Deliverability    │     │  Multi-Persona Seed   │
       │  Canvas & Delivery    │     │  Co-Pilot & Subject   │     │  Simulation & Dialog  │
       │  Path Topology        │     │  Entropy Engine       │     │  Dynamics Engine      │
       └───────────┬───────────┘     └───────────┬───────────┘     └───────────┬───────────┘
                   │                             │                             │
                   └─────────────────────────────┼─────────────────────────────┘
                                                 │
                                                 ▼
                         ┌──────────────────────────────────────────────┐
                         │      Core Architecture & Shared Invariants   │
                         │   - Split-DB SQLite Store (`warmup.db`)      │
                         │   - Error Return Sovereignty (Zero Exits)    │
                         │   - Pure HTML5/SVG/CSS3 (Zero External CDN)  │
                         │   - Coupled CLI Tools (`gitmap warmup ...`)  │
                         │   - Positive Booleans (`is*`, `has*` Only)   │
                         └──────────────────────────────────────────────┘
```

---

## 2. Antigravity Design Parity & UI Foundations

Sequence 338 enforces 100% design parity with the Google Antigravity Manager UX framework:
- **Zero External CDNs:** Pure local CSS3 custom properties, embedded SVG vectors, and native DOM event listeners. No external scripts, unpkg fonts, or web assets.
- **Glassmorphic Cyber-Card Deck:** 1px hairline borders (`rgba(255, 255, 255, 0.08)`), high-contrast backdrop blurs (`backdrop-filter: blur(16px)`), deep obsidian surfaces (`#0b0f19`), and focused ambient accent glows.
- **Micro-State Determinism:** UI states expose strict positive boolean attributes (`isActive`, `hasWarning`, `isProcessing`, `hasSpamRisk`).
- **Kinetic Physics Engine:** Canvas elements employ dampening springs, particle flow representations, and cubic-bezier transitions (`300ms cubic-bezier(0.16, 1, 0.3, 1)`).

---

## 3. High-Level System Architecture & Component Interaction

```
   [Operator Web GUI / Terminal HUD] ◄────────┐
                 │                            │ WebSocket / Local SQLite Journal
                 ▼                            │
   [CLI Dispatcher: `gitmap warmup`]          │
                 │                            │
                 ▼                            │
   ┌──────────────────────────────────────────┴─────────────────────────┐
   │ GitMap Autonomous Warm-Up Orchestration Core                       │
   │                                                                    │
   │ ┌────────────────────────┐  ┌───────────────────────────────────┐  │
   │ │ Canvas Physics Engine  │  │ AI Deliverability Engine (Co-Pilot)│ │
   │ │ - Force-Directed Graph │  │ - Shannon Entropy Calculator      │  │
   │ │ - Bezier Flow Routing  │  │ - Spam Trap Syntax Heuristics     │  │
   │ │ - Gateway Weight Nodes │  │ - Multi-Turn Prompt Synthesizer   │  │
   │ └───────────┬────────────┘  └─────────────────┬─────────────────┘  │
   │             │                                 │                    │
   │             ▼                                 ▼                    │
   │ ┌─────────────────────────────────────────────────────────────┐    │
   │ │ Multi-Persona Seed Exchange Hub                             │    │
   │ │ - Synthetic Thread Generator & Branching Tree Graph         │    │
   │ │ - Autonomous IMAP / SMTP Simulation Workers                 │    │
   │ └─────────────────────────────┬───────────────────────────────┘    │
   └───────────────────────────────┼────────────────────────────────────┘
                                   │
                                   ▼
             ┌───────────────────────────────────────────┐
             │ SQLite Split-DB Storage (`.gitmap/warmup`)│
             │ - `CanvasNodes`, `AiEvaluations`          │
             │ - `PersonaPools`, `DeliveryThreads`       │
             └───────────────────────────────────────────┘
```

---

## 4. Design System Tokens & Invariants (Zero-CDN)

```css
:root {
  --bg-primary: #07090e;
  --bg-surface: #0e131f;
  --bg-surface-elevated: #161e31;
  --bg-overlay: rgba(14, 19, 31, 0.75);

  --border-subtle: rgba(255, 255, 255, 0.08);
  --border-accent: rgba(59, 130, 246, 0.35);
  --border-danger: rgba(239, 68, 68, 0.35);

  --text-primary: #f8fafc;
  --text-muted: #94a3b8;
  --text-dim: #64748b;

  --accent-cyan: #06b6d4;
  --accent-blue: #3b82f6;
  --accent-indigo: #6366f1;
  --accent-emerald: #10b981;
  --accent-amber: #f59e0b;
  --accent-rose: #f43f5e;

  --font-mono: "JetBrains Mono", "Cascadia Code", "Fira Code", monospace;
  --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;

  --shadow-node: 0 4px 20px -2px rgba(0, 0, 0, 0.5);
  --glow-accent: 0 0 15px -3px rgba(6, 182, 212, 0.4);
}
```

---

## 5. Specification Structure & Cluster Mapping

- `01-architecture-and-overview.md` (this file): System architectural blueprint, component mapping, and token dictionary.
- `02-improvements-101-to-130-canvas-ai-copilot-personas.md`:
  - **Cluster 11: Interactive Node Canvas & Delivery Path Routing** (Items 101–110)
  - **Cluster 12: AI Deliverability Co-Pilot & Subject Line Entropy** (Items 111–120)
  - **Cluster 13: Multi-Persona Seed Simulation & Conversational Dynamics** (Items 121–130)

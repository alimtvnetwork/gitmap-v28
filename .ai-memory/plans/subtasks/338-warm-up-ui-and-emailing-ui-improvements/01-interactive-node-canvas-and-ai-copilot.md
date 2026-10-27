# Subtask 01: Interactive Node Canvas & AI Co-Pilot

- **Parent Plan:** `.ai-memory/plans/pending/338-warm-up-ui-and-emailing-ui-improvements.md`
- **Spec Path:** `02-spec/21-app/338-warm-up-ui-and-emailing-ui-improvements/`
- **Sequence:** 338
- **Subtask Index:** 01
- **Status:** PLANNED
- **Created At:** 2026-10-27
- **Completed At:** PENDING

## Objectives

Implement and verify Improvements 101 through 130 spanning Clusters 11, 12, and 13:
1. **Cluster 11: Interactive Node Canvas & Delivery Path Topology (Items 101–110):**
   - Force-directed visual network topology canvas representing outbound SMTP relay nodes, edge hops, and delivery latencies.
   - Dynamic Bezier flow curves illustrating real-time email dispatch packets with color-coded throughput velocities.
   - Interactive drag-and-drop gateway node weight assignment and routing failover controls.

2. **Cluster 12: AI Deliverability Co-Pilot & Subject Line Entropy (Items 111–120):**
   - Local LLM prompt synthesizer and Shannon token entropy calculator to eliminate repetitive spam syntax patterns.
   - Real-time spam trigger keyword highlighter with contextual synonym replacement suggestions.
   - Inbox sentiment scoring engine analyzing simulated recipient engagement probability.

3. **Cluster 13: Multi-Persona Seed Simulation & Conversational Dynamics (Items 121–130):**
   - Heterogeneous seed persona engine (Executive, Technical, Buyer, Consumer) with randomized interaction schedules.
   - Realistic multi-turn conversation thread synthesizer generating authentic reply dialogues.
   - Seed engagement simulation physics: click velocities, scroll depths, and inbox folder moves (Mark as Important, Star, Unspam).

## Concrete Deliverables

- `pkg/web/canvas/topology_mesh.go`: Force-directed graph physics engine and Bezier routing calculator.
- `pkg/mail/ai/entropy_scorer.go`: Token entropy and spam trigger scoring logic.
- `pkg/mail/seed/persona_simulator.go`: Multi-persona seed agent generator and thread dialog weaver.
- `web/static/js/canvas_topology.js`: Local canvas renderer with particle animations (zero external CDNs).
- Unit tests covering physics dampening, Shannon entropy edge cases, and persona reply generation.

## Quality Standards

- 100% positive booleans (`isNodeActive`, `hasValidRoute`, `hasHighEntropy`).
- Functions strictly bounded between 8 and 15 lines.
- Pure Go standard library + embedded static assets; zero external network requests.

# Collision Resolution Summary Report

> Generated: 2026-02-28 | Source: `spec/03-error-manage/03-error-code-registry/01-registry.md`
**Version:** 3.2.0  

## Overview

13 error code collisions have been identified and resolved across the project ecosystem. This report provides a consolidated before/after view of all resolutions.

---

## Resolution Table

| # | Collision | Before (Stale) | After (Current) | Type |
|---|-----------|---------------|-----------------|------|
| **1** | AB Lovable Reasoning vs WebSocket Resilience | AB LR: 9830-9835 | AB LR: 10500-10519 → **19000-19019** (superseded by #13) | Reassign |
| **2** | Nexus Flow Reset API vs State Errors | NF Reset: 8301-8308 | NF Reset: **8350-8369** | Reassign |
| **3** | Link Manager vs AI Transcribe | LM: 14000-14999 | LM: **15000-15999** | Reassign |
| **4** | AIT Voice Commands vs Model Download | Model Download: 14200+ | Model Download: **14470-14489** | Reassign |
| **5** | WP Plugin Publish Local Codes | WPP: `E{x}xxx` (local) | WPP: **13000-13999** | Formalize |
| **6** | WP Plugin Builder Range Compressed | WPB: 10000-10999 | WPB: **10000-10499** (compressed) | Compress |
| **7** | SM Code Generation vs WSP | SM-CG: 12000-12799 | SM-CG: **16000-16799** | Reassign |
| **8** | SM Project Editor vs WPP | SM-PE: 13000-13999 | SM-PE: **17000-17999** | Reassign |
| **9** | PS/AB SEO 9500 Range Overlap | PS: 9500-9599, AB SEO: 9500-9540 | **No change** (intentional, format-separated) | Document |
| **10** | AIT Voice Codes vs WPP | AIT: 13200-13308 | Voice Cmd: **14200-14206**, Cloning: **14300-14307**, TTS: **14150-14158** | Reassign |
| **11** | SM Realtime vs WSP | SM-RT: 12001-12031 | SM-RT: **2800-2849** | Reassign |
| **12** | SM GSearch CLI vs Multiple | SM-GS: 1xxx-12xxx (92 codes) | SM-GS: **18000-18249** | Reassign |
| **13** | AB Lovable Reasoning vs WPB | AB LR: 10500-10519, WPB: 10500-10899 | AB LR: **19000-19019**, WPB: **10000-10499** | Reassign + Compress |

---

## Resolution Types

| Type | Count | Description |
|------|-------|-------------|
| Reassign | 9 | Module moved to a new, non-colliding range |
| Compress | 1 | Range narrowed with internal code remapping |
| Reassign + Compress | 1 | Both reassignment and compression applied |
| Formalize | 1 | Local codes formalized into the registry |
| Document | 1 | Intentional overlap documented (no change needed) |

---

## Current Active Allocations (Post-Resolution)

```
Range           Owner   Description
──────────────  ──────  ─────────────────────────────
0000-0999       GEN     General/Cross-cutting
2000-2999       SM      Spec Management (base)
2800-2849       SM-RT   SM Realtime (Resolution 11)
6010-6027       SM      Spec Editor + Error Recovery
7000-7099       GS      GSearch CLI Core
7100-7599       BR      BRun CLI
7600-7919       GS      GSearch Extensions
8000-8399       NF      Nexus Flow
8350-8369       NF      NF Reset API (Resolution 2)
9000-9599       AB/PS   AI Bridge + PowerShell
9600-9999       AB      AI Bridge Extended
10000-10499     WPB     WP Plugin Builder (Resolution 6/13)
10500-10999     ---     [UNALLOCATED]
11000-11999     SRC     Spec Reverse
12000-12599     WSP     WP SEO Publish CLI
13000-13999     WPP     WP Plugin Publish (Resolution 5)
14000-14499     AIT     AI Transcribe CLI
14500-14999     EQM     Exam Manager
15000-15999     LM      License Manager (Resolution 3)
16000-16799     SM-CG   SM Code Generation (Resolution 7)
17000-17999     SM-PE   SM Project Editor (Resolution 8)
18000-18249     SM-GS   SM GSearch CLI (Resolution 12)
19000-19019     AB      Lovable Reasoning (Resolution 13)
```

---

## Notes

- Resolution 1 was **superseded** by Resolution 13: AB Lovable Reasoning moved from 10500-10519 to 19000-19019.
- Resolution 9 is the only **no-action** resolution (intentional format-separated overlap).
- All resolutions are logged in the canonical registry at `spec/03-error-manage/03-error-code-registry/01-registry.md`.
- Remediation tracking: `spec/23-how-app-issues-track/16-error-code-collision-remediation.md`.

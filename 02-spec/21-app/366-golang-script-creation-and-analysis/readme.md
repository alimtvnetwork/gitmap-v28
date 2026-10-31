# Specification Index: 366 — Golang Script Creation and Analysis

- **Task Code:** `366-golang-script-creation-and-analysis`
- **Slug:** `golang-script-creation-and-analysis`
- **Root Directory:** `02-spec/21-app/366-golang-script-creation-and-analysis/`
- **Plan File:** `.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md`
- **Status:** `in-progress`
- **Date:** `31-Oct-26`

---

## 1. Executive Overview

This specification establishes a native, zero-dependency Chrome DevTools Protocol (CDP) client and inspection command suite within GitMap (`gitmap cdp`). It permanently eliminates the need for autonomous AI agents to construct ad-hoc 70-line Python scripts using low-level socket programming and manual RFC 6455 byte framing. 

By centralizing CDP inspection into a native Go package (`cli/cmdcdp/`), GitMap provides robust port and target discovery, high-capacity frame defragmentation for payloads exceeding 64 KB, typed DOM probing, and Tokyo Night Lipgloss terminal visualization with JSON fallback.

---

## 2. Specification Document Inventory

| # | Document | Scope & Purpose |
|---|----------|-----------------|
| 01 | [01-root-cause-and-architecture.md](01-root-cause-and-architecture.md) | Grounded 4-part Root Cause Analysis (RCA), failure modes of ad-hoc Python scripts, architectural flowcharts, and pure Go RFC 6455 WebSocket framing design. |
| 02 | [02-cdp-engine-and-cli-spec.md](02-cdp-engine-and-cli-spec.md) | Package decomposition (`cli/cmdcdp/`), concrete Go struct/interface definitions, CLI command syntax (`targets`, `probe`, `eval`), Tokyo Night Lipgloss specifications, and test verification matrix. |
| 99 | [99-consistency-report.md](99-consistency-report.md) | Module consistency report, structural health scorecard (100/100 A+), and coding guideline verification matrix. |

---

## 3. Core System Invariants

- **Zero Raw Sockets by AI:** Autonomous agents are strictly banned from authoring ad-hoc Python/Node raw socket scripts to probe browser environments; all browser inspections must utilize `gitmap cdp`.
- **Pure Standard Library RFC 6455 Client:** The WebSocket implementation in `cli/cmdcdp/client.go` uses exclusively standard library packages (`net`, `crypto/rand`, `encoding/binary`, `encoding/json`), maintaining zero third-party dependencies.
- **Dynamic Port & Target Resolution:** No hardcoded ports (e.g. `49994`) or target GUIDs. The engine automatically queries `/json/list` across standard and active debug endpoints.
- **Large Payload Safety:** Frame assembly handles chunked frames and responses up to 16 MB without truncation or socket timeouts.
- **Strict Coding Guidelines:** All functions adhere to the <= 15-line limit, enforce positive booleans (`is*`, `has*`), use zero nested `if` statements, and uphold leaf error sovereignty.

---

## 4. Cross-References

- Coding Guidelines: [02-spec/02-coding-guidelines/](../../02-coding-guidelines/)
- Error Management: [02-spec/03-error-manage/](../../03-error-manage/)
- Plan File: [.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md](../../../.ai-memory/plans/pending/366-golang-script-creation-and-analysis.md)
- Target Package: `cli/cmdcdp/`

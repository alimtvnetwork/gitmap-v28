# 38 — Fleet Nodes Clone Table Spacing and W3 Reachability Resilience

## Status: Learned
- **Domain**: Terminal Output Formatting, Fleet Orchestration, Diagnostics
- **Date**: 2026-10-01

---

## Key Learnings

1. **Table Column Gutter Margin**:
   - Column format widths must always exceed the maximum expected token length by at least 2 characters.
   - For duration columns, `%-10s` fails for 10-character tokens like `"in-process"`. Allocating `%-12s` ensures a minimum 2-character padding plus the inter-column delimiter space.
   - Replacing literal placeholders with real measurements (e.g. `14ms`) ensures consistent semantics across heterogeneous rows.

2. **Accurate Failure Classification vs Assumptions**:
   - Never assert a physical/virtual machine state (e.g. "machine is off") based solely on network connection timeout or failed ARP resolution.
   - Use non-presumptive phrasing: `(unreachable or port 22 closed)` or `node unreachable (host offline, adapter isolated, or IP changed)`.

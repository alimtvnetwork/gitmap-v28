# Subtask 03: Skills & Prompts Synchronization across GitMap and Coding-Guidelines

- **Parent Plan:** `359-gitmap-feature-implementation-and-update`
- **Slug:** `gitmap-feature-implementation-and-update - Task 03`
- **Status:** `pending`
- **Assigned Subagent:** Worker 03

---

## Deliverables
1. Update `.agents/skills/gitmap/skill.md`:
   - Add Section on Autonomous AI Identity (`gitmap login set ai-auto [agent-name]`, `gitmap login set global ai-auto`, `gitmap cpf/cpb --ai-agent <name>`).
   - Add Heatmap Attribution explanation (`gitmap heatmap`, `gitmap rc heatmap`, agent attribution).
   - Add operational guardrail / band on proper skill usage.
2. Update `.agents/skills/gitmap-automation-bot/skill.md`:
   - Document `gitmap login set global ai-auto` and `gitmap login set ai-auto [agent-name]`.
3. Update `execute-parent-task-with-n-steps-v7` (skill & prompt) in both repos:
   - Use `gitmap login set global ai-auto` & `gitmap login set ai-auto [agent-name]`.
   - Add completion tables checklist: Task Execution & Verification Matrix, Coding Guidelines & Invariants Compliance Table, File Ownership & Boundary Audit, Git Delivery & Upstream Invariant Table.
   - Add masked verification script requirement at the end.
   - Add conditional E2E testing checklist and explicit rule: if NOT asked, DO NOT run build or E2E tests!
4. Update `execute-parent-task-with-n-steps-v6` (skill & prompt) in both repos:
   - Add completion tables checklist, masked verification script, and conditional E2E testing rule.
5. Mirror changes to `d:/work/coding-guidelines/`.

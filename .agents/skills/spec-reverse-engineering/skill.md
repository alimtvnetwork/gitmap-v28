---
name: spec-reverse-engineering
description: >-
  Autonomously scan, reverse-engineer, and synthesize comprehensive architectural specs and security audits from any codebase.
---

# Instruction (must follow): Reverse Engineering Spec Writer & Architecture Discovery

/goal Autonomously scan, reverse-engineer, and synthesize an exhaustive, multi-file architectural specification of any target codebase into `spec/21-app/` (or `<spec-folder>/21-app/`), detecting programming languages, isolating architectural boundaries, analyzing security vulnerabilities and risk tiers, and structuring all findings with strictly lowercase naming and relative git paths without stopping until 100% green.

```text
N = 80
```

Source Prompt Reference: `01-prompts/03-read-write/07-spec-write-from-reverse-engineering.md`

## Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Step 0 (Dynamic Environment Bootstrap): Detect if `.lovable/` and `.lovable/temp/` exist; if missing, dynamically create `.lovable/temp/`.
2. [ ] /goal Step 0 (Spec Folder Initialization): Detect or create the target spec folder (`spec/21-app/` or `<spec-root>/21-app/`) adhering to two-digit monotonic numbering and strictly lowercase directory conventions.
3. [ ] /goal Step 1 (Tooling Script Creation): Create or verify the Python inventory tool (`03-ai-scripts/<NN>-codebase-file-lister.py` or `.lovable/<temp>/list_files.py`) to crawl all repository files while excluding `.git`, `node_modules`, `vendor`, and build artifacts.
4. [ ] /goal Step 2 (Codebase Crawl & Language Topology): Execute the inventory script to generate `.lovable/<temp>/files-inventory.json` with relative paths, file sizes, and language breakdown percentages.
5. [ ] /goal Phase 1 (Planning & File Management Ledger, Steps 1..N/2): Partition discovered files into balanced batches across 2–3 concurrent sub-agents, recording active locks in `.lovable/<temp>/file-assignments.json`.
6. [ ] /goal Phase 1 (Micro-Tasking Subtasks): Decompose file analysis into granular subtask files under `.lovable/plans/subtasks/<reverse-engineering>/` tracking which sub-agent processes which file paths.
7. [ ] /goal Phase 1 (Zero-Stop Transition): Transition directly into Phase 2 execution mode without pausing or requesting user input.
8. [ ] /goal Phase 2 (Parallel File Reverse Engineering, Steps N/2+1..N): Dispatch 2–3 execution sub-agents in parallel on disjoint file sets to analyze code semantics, exported types, data flow, functions, and external dependencies.
9. [ ] /goal Phase 2 (Component Specification Generation): Write dedicated, modular specification files under `spec/21-app/` (e.g., `spec/21-app/xx-domain-models.md`, `spec/21-app/xx-service-layer.md`, `spec/21-app/xx-api-endpoints.md`).
10. [ ] /goal Phase 3 (Master Index Synthesis): Author `spec/21-app/01-index.md` summarizing the overall application architecture, behavior, technology stack, architectural health score, and component topology.
11. [ ] /goal Phase 3 (Security Audit & Risk Assessment): Identify hardcoded credentials, unauthenticated endpoints, input sanitization flaws, and dependency vulnerabilities, publishing an exhaustive risk evaluation in `spec/21-app/xx-security-and-risks.md`.
12. [ ] /goal Phase 3 (Final Structure Communication): Output a clean, viewable markdown/ASCII folder tree in the final chat response illustrating the complete generated specification layout.
13. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
14. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
15. [ ] /learn Ingest `spec/02-coding-guidelines/08-file-folder-naming/` for lowercase naming and continuous file sequencing.
16. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
17. [ ] /learn Ingest `spec/21-app/01-index.md` for baseline application documentation standards.

## Strict Rules & Invariants

- **Sequential Self-Looping:** Never attempt to process all files in a single turn. Work across the $N=80$ budget.
- **Strict Relative Git Paths:** TOTAL BAN on absolute filesystem paths or `file:///` URIs.
- **Strict Lowercase File Naming:** All created files and directories must use lowercase characters and hyphens.
- **File Management System:** Maintain `.lovable/<temp>/file-assignments.json` to guarantee zero concurrent file collisions between agents.

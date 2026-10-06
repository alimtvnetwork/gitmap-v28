# Master Execution Plan: 236-deep-spec-consolidation-and-canonical-reduction

## User Request (Verbatim)
```text
# High Priority Instruction

I want you to go through the specs, plans, and subtasks. Wherever there is aging information, information which is not necessary to the spec, try to consolidate the files if possible. Try to reduce definition if not necessary. Try to reduce as much spec file as possible, and only keep whatever is the final version. For example, let's say a spec started saying X, Y, and then it said A, B, so we will just keep the A, B, the last one. You have to understand and go very deep. I want you to read all this spec, probably spend 300 steps properly, try to learn the code base, and then try to reduce the spec. We need to go inside the spec folder, spec folder 21 only, and also the AI memory files memory, and also the plans file, which are completed, not the pending files. Pending files, do not touch. Only those are completed. You will try to consolidate and reduce the writing. Before doing that, you will release and keep a backup branch, and then after you do the changes, you will make a final release as well.

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Go through the specs, plans, and subtasks to consolidate files and reduce unnecessary definitions.
5. Keep only the final version of specs, removing outdated information.
6. Spend 300 steps to understand the codebase and reduce the spec.
7. Focus on spec folder 21, AI memory files, and completed plans files only.
8. Do not touch pending files.
9. Release and keep a backup branch before making changes.
10. Make a final release after changes are completed.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## 1. Specification & Subtask Index

- **Master Ledger**: [02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/00-master-audit-ledger.md](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/00-master-audit-ledger.md)
- **Architecture Spec**: [02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)
- **Component Spec**: [02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/02-component-and-cli-spec.md](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/02-component-and-cli-spec.md)
- **Subtasks**:
  1. `01-pre-consolidation-safety-release-and-backup-branch.md`
  2. `02-deep-survey-and-aging-specs-inventory.md`
  3. `03-deep-consolidation-of-spec-folder-21.md`
  4. `04-deep-consolidation-of-completed-plans-and-memory.md`
  5. `05-registries-sync-and-post-release.md`

---

## 2. Multi-Agent Delegation Plan

- **Phase 0: Safety & Baseline Release**:
  - Lead: Bump minor version to `v6.500.0`, push release branch and tag `v6.500.0`. Create and push `backup/pre-deep-spec-consolidation-20261006` branch to origin.
- **Phase 1: Discovery (`A = 2` `research` subagents)**:
  - Research 01: Audit `02-spec/21-app/` surviving standalone folders to map into the 8 Canonical Domain Clusters.
  - Research 02: Audit `.ai-memory/plans/completed/` and `.ai-memory/memory/` for deeper compaction opportunities.
- **Phase 1: Spec Authoring (`A = 2` `self` subagents)**:
  - Author 01: Author `01-architecture-spec.md` and Subtasks 01-02.
  - Author 02: Author `02-component-and-cli-spec.md` and Subtasks 03-05.
- **Phase 2: Execution Waves (`A = 2` `self` worker subagents)**:
  - Wave 1:
    - Worker 01: Fold surviving legacy spec directories in `02-spec/21-app/` into the 8 Canonical Clusters.
    - Worker 02: Deep-compact `.ai-memory/plans/completed/` and `.ai-memory/memory/` while strictly isolating pending files.
  - Wave 2:
    - Worker 01: Synchronize registries (`02-spec/21-app/readme.md`, `.ai-memory/plans/readme.md`, `.ai-memory/what-to-read.md`).
    - Worker 02: Execute verification gates (relative paths, linters, `go vet`).
- **Phase 3: Final Release**:
  - Lead: Bump version to `v6.501.0`, tag, release branch, atomic GitMap commit, and showcase report.

# Subtask 04: Minor Version Bump & Release Ceremony

- **Parent Plan:** `359-gitmap-feature-implementation-and-update`
- **Slug:** `gitmap-feature-implementation-and-update - Task 04`
- **Status:** `completed`
- **Assigned Subagent:** Lead Orchestrator

---

## Deliverables
1. Run `python 03-ai-scripts/37-bump-version.py --tier minor --scope "Implement gitmap login set ai-auto, commit --ai-agent flag, heatmap attribution, and V6/V7 verification checklist synchronization"`.
2. Confirm version bumped from `6.632.1` to `6.633.0` across:
   - `version.json`
   - `package.json`
   - `changelog.md`
   - `readme.md`
   - `what-to-read.md`
   - `cli/constants/constants.go`
   - `.gitmap/release/latest.json`
   - `.ai-memory/release/release-notes-v6.633.0.md`
3. Execute `gitmap cpr "v6.633.0 implement gitmap login set ai-auto commit ai-agent flag and sync v6 v7 verification checklists"`.
4. Push upstream to remote.

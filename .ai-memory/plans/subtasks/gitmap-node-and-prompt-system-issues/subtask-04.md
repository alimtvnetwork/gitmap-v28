# Subtask 04: Isolated Testing, FastGate & Minor Bump Release Ceremony

**Parent Plan:** `.ai-memory/plans/gitmap-node-and-prompt-system-issues.md`  
**Subtask ID:** `subtask-04`  
**Owner:** Lead Agent  
**Status:** Pending  

## Scope of Ownership
- Test harnesses (`cli/tests/`)
- `03-ai-scripts/50-fastgate.py`
- `03-ai-scripts/37-bump-version.py`
- Release notes and changelogs

## Implementation Checklist
- [ ] Safeguard running instance: NEVER touch, kill, or close `default-copy-814`. Use isolated test instances only.
- [ ] Run targeted Go tests across modified packages (`cmdnodes`, `cmdagy`, `cmdsupabase`, `cmdui`).
- [ ] Run TypeScript checks: `npx tsc --noEmit`.
- [ ] Run FastGate: `python 03-ai-scripts/50-fastgate.py`.
- [ ] Execute minor version bump per release ceremony: `python 03-ai-scripts/37-bump-version.py --minor`.
- [ ] Commit with hyphenated message: `gitmap cpf "node-prompt-fixes - resolve fleet telemetry prompt tree and supabase throttling"`.
- [ ] Monitor CI/CD pipeline via `gitmap pe -t` until green.

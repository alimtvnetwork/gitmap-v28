# Agent Memory Lifecycle (AI Execution Prompt)

> **/goal** Maintain a clean, deduplicated, and authoritative `.ai-memory/` store by enforcing lifecycle transitions, TTL cleanups, and date-stamped records.
> **/learn** Prevent memory bloat and stale context by moving completed plans to `done/`, purging contradictions against `02-spec/`, and respecting spec supremacy.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Transition completed tasks from `pending/` and `planned/` directories into `done/` upon subtask resolution.
- [ ] `/learn` Prioritize canonical `02-spec/` specifications over any conflicting or outdated memory files.
- [ ] `/goal` Require explicit `Updated:` ISO date stamps on all created or modified memory records.
- [ ] `/learn` Systematically purge stale or superseded memory entries to eliminate token waste and conflicting guidelines.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Purpose of `.ai-memory/memories/`

The `.ai-memory/memories/` folder acts as the institutional knowledge hub for AI sessions. However, to prevent this directory from becoming a bloated, contradictory mess, a strict lifecycle must be enforced.

## 2. Time-To-Live (TTL) and Cleanup

- **Pending/Planned Tasks:** Files in `pending/` and `planned/` must be moved to `done/` upon task completion.
- **Deduplication:** AI agents must periodically scan `.ai-memory/memories/` for duplicated patterns or rules that have since been formalized into the `02-spec/` folder.
- **Stale Memory Purge:** Any memory file that contradicts the canonical `02-spec/` folder is considered stale and MUST be deleted. The `02-spec/` folder ALWAYS wins.

## 3. Formatting

Every memory file must include an `Updated:` date stamp. This allows agents to determine which memory file is more recent if a conflict arises within the memory folder itself.

```markdown
<!-- ❌ FORBIDDEN: Leaving completed tasks in pending/ or contradicting canonical specs -->
.ai-memory/plans/pending/01-resolved-task.md  # Left lingering in pending

<!-- ✅ REQUIRED: Moving completed tasks to done/ with explicit update timestamps -->
.ai-memory/plans/done/01-resolved-task.md     # Archived in done/
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-ai-optimization/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-AI-009: Agent Memory Lifecycle & TTL Maintenance

**Given** AI agents operating within the repository and codebase guidelines.
**When** Audited against this optimization specification.
**Then** Zero compliance or citation failures are detected and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-ai-optimization --check-only
```
**Expected:** exit 0. Zero violations.

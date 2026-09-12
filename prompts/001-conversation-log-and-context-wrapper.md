# Conversation Log and Context Wrapper Specification

## 1. Context & Purpose
This prompt defines the execution workflow for persisting session conversation histories, auditing project memory against coding standards, staging follow-up instructions into `prompts/`, and generating comprehensive audit reports without prematurely executing unconfirmed work.

- **Roadmap Source of Truth:** `.lovable/plan.md`
- **Memory Index:** `mem://01-index.md` (`.lovable/memory/01-index.md`)
- **Coding Guidelines:** `.lovable/coding-guidelines.md` and `spec/02-coding-guidelines/`
- **Error Management Standard:** `spec/03-error-manage/`

---

## 2. Inputs & Prerequisites
1. Full prior conversation turns across the active session.
2. Verified clean working tree (`git status` exits clean).
3. Active git branch only (no unauthorized branch switching).

---

## 3. Scope & Execution Checklist

### Phase 1: Conversation Persistence
- [ ] Ensure repository root directory `conversation/` exists.
- [ ] Enumerate existing files in `conversation/` and assign the next sequential zero-padded 3-digit prefix (`001`, `002`, ...).
- [ ] Partition chat history into logical topic chunks (or one file if single-topic).
- [ ] Write each chunk file using lowercase kebab-case naming (`conversation/NNN-<topic-slug>.md`).
- [ ] Include header metadata: `Sequence`, `CapturedUtc` (ISO-8601 UTC), `Span` (user message count), `Topic`.
- [ ] Capture user instructions **verbatim** within blockquotes. Do not summarize or alter user words.
- [ ] Provide terse factual bullet points for assistant actions (no chain-of-thought or internal reasoning).
- [ ] Document concrete outcomes/decisions and open threads.

### Phase 2: Memory & Standards Verification
- [ ] Inspect `mem://01-index.md` (`.lovable/memory/01-index.md`).
- [ ] Verify reference to `.lovable/coding-guidelines.md`; if missing, stage a proposed write for `mem://standards/coding-guidelines.md`.
- [ ] Verify reference to `.lovable/plan.md`; if missing or referencing legacy paths, stage a proposed update.
- [ ] Verify reference to `mem://workflow/conversation-log`; if missing, stage a proposed write.
- [ ] List all memory files that the follow-up task will touch for user audit prior to execution.
- [ ] **Do NOT** silently auto-create or overwrite memory files without user confirmation.

### Phase 3: Instruction Staging
- [ ] Stage rewritten follow-up instructions into `prompts/NNN-<slug>.md` with 3-digit prefix.
- [ ] Append staged instructions to `01-prompts/01-prompt-library-setup/01-prompt-library-setup.md` table without overwriting historical rows.
- [ ] Enforce definition of done, input/output contracts, and explicit verification criteria.

---

## 4. Definition of Done
1. All prior user instructions persisted verbatim in `conversation/NNN-<slug>.md`.
2. Memory references audited and proposed writes clearly reported.
3. Follow-up prompt staged in `prompts/` and indexed in `01-prompts/01-prompt-library-setup/01-prompt-library-setup.md`.
4. Staging report presented to the user with zero execution of the follow-up task.
5. All 36 CI/CD quality gates remain green (`python 03-ai-scripts/06-cicd-local-runner.py` exits 0).
6. Clean git commit and push using `--no-verify`.

---

## 5. Mandatory Self-Instructions

> **Self-Instruction:** Before acting, re-read `mem://01-index.md` and `.lovable/coding-guidelines.md`; restate which rules apply.
>
> **Continuous Improvement:** After executing this task, analyze the workflow and suggest further structural improvements to this instruction file.

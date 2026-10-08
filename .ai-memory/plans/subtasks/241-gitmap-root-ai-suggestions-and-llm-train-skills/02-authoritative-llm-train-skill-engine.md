# Subtask Plan 02: Authoritative LLM Train & Self-Skill Creation Engine

## Overview
Elevate `gitmap llm train` from a basic textual summary into an authoritative skill delivery system. It emits the complete Antigravity skill with full YAML frontmatter and cheatsheets directly to stdout and writes/updates `.agents/skills/gitmap/SKILL.md`. It explicitly directs the calling LLM to autonomously create and persist its own skill from this training output.

---

## Owned Files
- `cli/cmd/llm/llm_train.go`
- `cli/cmd/llm/llm_skill.go`
- `cli/cmd/llm/llm_types.go`
- `cli/cmd/llm/llm_train_test.go`

---

## Step-by-Step Implementation

1. **Update `cli/cmd/llm/llm_types.go`:**
   - Add `IsURLs bool` to `TrainFlags`.
   - Ensure clean struct definitions with positive booleans.

2. **Refactor `cli/cmd/llm/llm_train.go`:**
   - In `RunTrain(args []string)`:
     * When executed without flags (or standard training), emit the complete Antigravity skill to stdout ("Skill Method Primacy").
     * Write/update `.agents/skills/gitmap/SKILL.md` (unless `--text-only` is provided).
     * Add banner instructing the AI: "🤖 Direct Skill Ingestion: Review the generated skill below. As an autonomous AI model, you must persist or adapt this skill into your active skill directory (.agents/skills/gitmap/SKILL.md) and adhere to all command substitutions."
     * Parse `--urls` flag to emit the public documentation catalog.
     * Incorporate recursive learning directives.

3. **Update `cli/cmd/llm/llm_skill.go`:**
   - Ensure `SkillTemplate` contains complete, up-to-date GitMap cheatsheet and command replacement matrix (e.g. banning `Get-ChildItem -Filter`, `git grep`, raw `Remove-Item`, `rg`, `findstr`).
   - Support custom skill paths and verify atomic writes.

4. **Update Unit Tests in `cli/cmd/llm/llm_train_test.go`:**
   - Test default `gitmap llm train` emits skill content containing YAML frontmatter (`name: gitmap`).
   - Test `--urls` flag emits Markdown documentation links.
   - Test `--text-only` skips writing to disk.

---

## Acceptance Criteria
- [ ] `gitmap llm train` emits full Antigravity skill directly to stdout.
- [ ] Automatically updates `.agents/skills/gitmap/SKILL.md`.
- [ ] Directs AI models to persist/adapt the skill into their workspace skills.
- [ ] Full command replacement matrix included in output.
- [ ] All unit tests in `cli/cmd/llm/` pass.

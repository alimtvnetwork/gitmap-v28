# Subtask 336.03 — Coding Guidelines & Antigravity Skill Synchronization

> **Sequence ID:** 336  
> **Parent Plan:** [336-automatic-binary-swap-and-skill-synchronization.md](.ai-memory/plans/pending/336-automatic-binary-swap-and-skill-synchronization.md)  
> **Target Packages:** `.agents/skills/`, `.agents/rules/`, `03-ai-scripts/`  
> **Status:** Pending  

---

## 1. Objectives

1. Create a new specialized Antigravity skill `.agents/skills/gitmap-binary-swap/skill.md`:
   - Document how agents can invoke `gitmap swap-binary` when recompiling `gitmap` from source.
   - Detail the safety contract: automatic PID discovery, detached supervisor script, retry backoff, and `--version` verification.
   - Provide command snippets for Windows (`powershell`) and Unix (`bash`).
2. Update `.agents/rules/gitmap_guidelines.md`:
   - Document `gitmap swap-binary` as the mandatory command for replacing active GitMap binaries after local source compilation.
   - Enforce positive boolean naming (`is`/`has`) and 100% relative paths across all plan and spec documents.
3. Synchronize canonical skills, prompts, and rules across connected workspace repositories:
   - Run `03-ai-scripts/38-sync-prompts-skills-scripts.py` to mirror updated skills and rules to connected targets (`antigravity-manager`, `spec-builder`, `macro-ahk`, `movie-cli`, etc.).
   - Verify zero regression or uncommitted drift across synced repositories.

---

## 2. File Manifest

- `.agents/skills/gitmap-binary-swap/skill.md` — New Antigravity skill definition.
- `.agents/rules/gitmap_guidelines.md` — Updated operational guidelines for binary swapping.
- `03-ai-scripts/38-sync-prompts-skills-scripts.py` — Execution of multi-repo sync tool.

---

## 3. Verification Criteria

- [ ] `.agents/skills/gitmap-binary-swap/skill.md` contains valid YAML frontmatter and markdown documentation.
- [ ] `.agents/rules/gitmap_guidelines.md` documents `gitmap swap-binary` usage rules.
- [ ] Multi-repo sync script completes with status code 0 and reports all connected repositories in sync.

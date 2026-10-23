# Subtask 05: Coding Guidelines Authoring, Index Update & GitMap Skills Synchronization

## Status: READY
- **Task ID:** `315-subtask-05`
- **Parent Spec:** [02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/00-overview.md](../../../../02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/00-overview.md)
- **Module:** `02-spec/02-coding-guidelines/06-ai-optimization`, `.agents/skills/gitmap`
- **Owned Files:**
  - `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md`
  - `02-spec/02-coding-guidelines/06-ai-optimization/readme.md`
  - `.agents/skills/gitmap/skill.md`

---

## 1. Objectives

1. **Author Canonical Coding Guideline in `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md`:**
   - Establish repository-wide architectural standards for universal AI commit provenance, ambient fallback resolution, automatic flag stripping, and cross-database audit ingestion.
   - Enforce Rule 1: Zero Unattributed AI Commits in automated swarms (AI agents must pass `--agent-id`, `--heatmap-id`, and `--task` or define ambient environment variables).
   - Enforce Rule 2: Automatic Flag Stripping (provenance flags must never contaminate or leak into git commit messages).
   - Enforce Rule 3: Ambient Fallback Precedence Hierarchy (`CLI Flag` > `GITMAP_* Env` > `Platform Agent Env`).
   - Enforce Rule 4: Hyphen-separated commit format (`<module> - <summary>`) strictly preserved without colons.
   - Enforce Rule 5: External Split-DB isolation (`repodb/commit_history.db` outside repository tree).
   - Enforce Rule 6: Origin attribution disambiguation (local machine AI vs external worker/node/human).

2. **Update Coding Guidelines Index in `02-spec/02-coding-guidelines/06-ai-optimization/readme.md`:**
   - Register Guideline 12 in the master guidelines catalog table with link and summary description.
   - Increment guideline inventory counts and update cross-references.

3. **Update GitMap Skill Specification in `.agents/skills/gitmap/skill.md`:**
   - Document `gitmap commit info <sha>` (and aliases `gitmap ci <sha>`, `gitmap commit-info <sha>`) with argument options, flags (`--json`, `--repo`), and output samples.
   - Document commit provenance flags (`--agent-id`, `--heatmap-id`, `--ai`, `--spec`, `--task`) across commit commands (`commit`, `cm`, `cpf`, `cpb`, `cpr`, `cpc`, `pcp`).
   - Add ASCII tree visualizer examples illustrating author details, AI attribution, origin disambiguation, and file touch statistics.
   - Add operational cheatsheet workflows for developers and AI agents inspecting commit provenance.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Author `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md`
1. Draft document metadata, problem statement, and executive intent.
2. Formulate the 6 core non-negotiable rules for commit provenance:
   - Specification of `--agent-id`, `--heatmap-id`, `--ai`, `--spec`, and `--task`.
   - Resolution hierarchy: CLI arguments $\to$ `GITMAP_*` env vars $\to$ `ANTIGRAVITY_*` / `AI_MODEL` env vars.
   - Prohibition of flag tokens in git commit messages.
   - Split-DB `<EffectiveDataDir>/repodb/commit_history.db` storage rules.
   - Origin attribution notice logic.
3. Provide concrete code and command examples (Go snippets, CLI examples, bad vs good patterns).
4. Include anti-patterns to avoid (e.g., hardcoding commit messages with unstripped flags, storing databases in repository root).

### Step 2: Register Guideline 12 in `02-spec/02-coding-guidelines/06-ai-optimization/readme.md`
1. Locate the catalog table in `readme.md`.
2. Add row for Guideline 12:
   ```markdown
   | 12 | [12-commit-provenance-and-tracking.md](./12-commit-provenance-and-tracking.md) | Universal AI commit provenance, task linkage, ambient fallback hierarchy, and audit logging |
   ```
3. Verify all relative links and formatting in the table.

### Step 3: Update `.agents/skills/gitmap/skill.md`
1. Locate the commit commands section in `.agents/skills/gitmap/skill.md`.
2. Add comprehensive documentation for `gitmap commit info <sha>`:
   - Syntax, arguments, shorthands (`ci`, `commit-info`).
   - Supported SHA formats: full 40-char SHA, 7-char short SHA, prefix, suffix, substring.
   - Output fields and ASCII tree representation.
3. Add provenance flag reference table for `gitmap commit / cm / cpf / cpr / cpc`:
   - `--agent-id <id>`: Identifier of the executing subagent.
   - `--heatmap-id <id>`: Cognitive heatmap action or session ID.
   - `--ai <name>`: Model identifier (e.g., `gemini-2.5-pro`, `claude-3-7-sonnet`).
   - `--spec <path>`: Relative path to related application specification.
   - `--task <slug>`: Associated parent task or subtask slug.
4. Add sample terminal session demonstrating `gitmap commit info dc8e0d8` with local AI match and external commit origin fallback.

---

## 3. Verification & Quality Gates

| Verification Target | Tool / Script | Expected Outcome |
| :--- | :--- | :--- |
| **Relative Path Hygiene** | `03-ai-scripts/01-check-relative-paths.py` | 100% PASS: Zero absolute paths or drive letters in guideline and skill files. |
| **Boolean & Enum Conventions** | `03-ai-scripts/18-check-enum-and-boolean.py` | 100% PASS: All boolean examples adhere to `is*`/`has*` positive naming. |
| **Catalog Link Integrity** | Manual Link Check | 100% PASS: Relative link `./12-commit-provenance-and-tracking.md` resolves cleanly from `readme.md`. |
| **GitMap Skill Cheatsheet** | Terminal rendering check | 100% PASS: Cheatsheet blocks format cleanly without broken markdown fences. |

---

## 4. Definition of Done
- [ ] `12-commit-provenance-and-tracking.md` authored with exhaustive rules, rationale, and examples.
- [ ] `readme.md` catalog updated with Guideline 12.
- [ ] `.agents/skills/gitmap/skill.md` updated with `commit info` command and provenance flags documentation.
- [ ] Python relative paths linter confirms 100% compliance.

# Subtask 05: Muse Master Prompt & GitMap Skills Modernization

> **Subtask ID:** Subtask-05  
> **Parent Plan:** `.ai-memory/plans/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md`  
> **Associated Specs:**  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-architecture-spec.md`  
> - `02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/02-component-spec.md`  
> **Owned Files:**  
> - `01-prompts/27-muse-prompts/01-muse-master-prompt.md`  
> - `.agents/skills/muse-master-prompt/skill.md`  
> - `.cursor/skills/muse-master-prompt/skill.md`  
> - `.agents/skills/gitmap/SKILL.md`  
> - `.cursor/skills/gitmap/skill.md`  
> - `02-spec/02-coding-guidelines/06-ai-optimization/readme.md`  
> - `02-spec/02-coding-guidelines/06-ai-optimization/02-anti-hallucination-rules.md`  
> - `02-spec/02-coding-guidelines/06-ai-optimization/03-ai-quick-reference-checklist.md`  

---

## 1. Objectives

- [ ] 1. Modernize `01-prompts/27-muse-prompts/01-muse-master-prompt.md`:
  - Integrate native GitMap tool primacy into agent directives and instructions:
    - Mandate `gitmap aum search` over slow shell grep/Select-String.
    - Mandate `gitmap task` for SQLite multi-agent state tracking.
    - Mandate `gitmap run` for universal script execution with audit logging.
    - Mandate `gitmap pe --ai` for all AI CI/CD error inspection (clipboard-free).
    - Mandate `gitmap sync` for multi-repo synchronization.
    - Mandate `gitmap cpc "<module> - <summary>"` for chore commits alongside `cpf`, `cpb`, and `cpr`.
  - Include full forms in help and command descriptions.
- [ ] 2. Synchronize Muse skill definitions:
  - Update `.agents/skills/muse-master-prompt/skill.md` and `.cursor/skills/muse-master-prompt/skill.md` ensuring 1:1 parity with the updated master prompt.
- [ ] 3. Modernize GitMap engineering skills:
  - In `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md`:
    - Update Non-Negotiable Command Replacement Matrix:
      - Add `gitmap pe --ai` as mandatory replacement for `gitmap pe` / `gh run view` in AI workflows.
      - Add `gitmap run <file>` as universal execution engine.
    - Add documentation for new commands:
      - `gitmap run <file>` (auto-interpreter resolution, task audit DB, errors DB).
      - `gitmap run errors` / `gitmap run history` (error diagnostics and retry).
      - `gitmap supabase add <alias> <url> <anon> <service> [db_url]` (encrypted vault).
      - `gitmap supabase list`, `gitmap supabase test`, `gitmap supabase remove`.
      - `gitmap cpc "<module> - <summary>"` (commit-push-chore).
      - `gitmap pe -ud` (until-done continuous polling).
      - `gitmap pe -t` (2-minute poll with early error abort).
- [ ] 4. Update Coding Guidelines in `02-spec/02-coding-guidelines/06-ai-optimization/`:
  - In `02-anti-hallucination-rules.md`:
    - Add rule `AH-GM1: Mandate gitmap pe --ai for AI Error Telemetry`: AI agents must never run plain `gitmap pe` without `--ai` to prevent clobbering the host system clipboard.
    - Add rule `AH-GM2: Ban Raw Shell File Traversal & Ambient Interpreters`: Mandate `gitmap aum search`, `gitmap find`, and `gitmap run`.
  - In `03-ai-quick-reference-checklist.md`:
    - Add checklist items for GitMap command substitution and clipboard-free AI inspection.
  - In `readme.md`:
    - Update catalog and file descriptions to reflect new rules.
- [ ] 5. Verify documentation integrity and cross-references:
  - Ensure zero absolute paths and zero `file:///` URIs across all markdown files.
  - Verify skill parity between `.agents/skills/` and `.cursor/skills/`.

---

## 2. Updated Coding Guidelines Rules

```markdown
### AH-GM1: Mandate gitmap pe --ai for AI Error Telemetry

❌ **Never generate/execute:** `gitmap pe`, `gitmap pipeline error-logs`, or `gh run view` from automated AI agents without the `--ai` flag.
✅ **Always generate/execute:** `gitmap pe --ai` or `gitmap pe -t --ai`.
📖 Running `gitmap pe` without `--ai` writes the full log dump to the system clipboard, clobbering the user's active clipboard data and generating unnecessary terminal notices. The `--ai` flag cleanly pipes output directly to stdout/stderr.

### AH-GM2: Ban Raw Shell File Traversal & Ambient Interpreters

❌ **Never generate/execute:** `python script.py`, `powershell -File script.ps1`, `bash script.sh`, `Select-String`, `rg`, or `findstr`.
✅ **Always generate/execute:** `gitmap run <script>`, `gitmap aum search "<pattern>"`, or `gitmap find "<pattern>"`.
📖 `gitmap run` automatically validates file extensions, discovers verified interpreters, maintains an execution audit trail in SQLite, and captures errors in a queryable errors database.
```

---

## 3. Acceptance Criteria

- **AC-MUSE-001 (Master Prompt Modernization):** `01-prompts/27-muse-prompts/01-muse-master-prompt.md` incorporates `gitmap aum search`, `gitmap task`, `gitmap run`, `gitmap pe --ai`, `gitmap sync`, and `gitmap cpc`.
- **AC-SKILL-002 (GitMap Skill Completeness):** `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md` document all new commands (`run`, `supabase`, `cpc`, `pe --ai`, `pe -ud`, `pe -t`) in the command replacement matrix and cheat sheet.
- **AC-CG-003 (Coding Guidelines Alignment):** `02-spec/02-coding-guidelines/06-ai-optimization/02-anti-hallucination-rules.md` contains `AH-GM1` and `AH-GM2`, strictly mandating `gitmap pe --ai` and `gitmap run`.
- **AC-PARITY-004 (Skill Parity):** All updates to `.agents/skills/` are mirrored identically in `.cursor/skills/`.
- **AC-HYGIENE-005 (Relative Path Hygiene):** All documentation files use strictly relative Git paths with zero absolute paths and zero `file:///` URIs.

---

## 4. Verification Instructions

1. Verify skill and prompt files exist and are populated:
   ```bash
   gitmap cat 01-prompts/27-muse-prompts/01-muse-master-prompt.md
   gitmap cat .agents/skills/gitmap/SKILL.md
   ```
2. Verify relative path hygiene across the repository:
   ```bash
   python linter-scripts/check-relative-paths.py
   ```

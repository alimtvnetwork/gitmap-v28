# Subtask 08: Universal Readme Text & Timestamp Generator Eradication

- **Subtask Path**: `.ai-memory/plans/subtasks/gitmap-enhancement-and-issue-resolution/08-eradicate-readme-text.md`
- **Parent Plan**: [.ai-memory/plans/gitmap-enhancement-and-issue-resolution.md](../../gitmap-enhancement-and-issue-resolution.md)
- **Specification Reference**: [02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md](../../../02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md)
- **Primary Owner**: Worker 02 (Governance, Specifications & Code Quality Specialist)
- **Target Files (Exclusive)**:
  - `.ai-memory/strictly-avoid.md`
  - `02-spec/01-spec-authoring-guide/11-exceptions.md`
  - `02-spec/17-consolidated-guidelines/03-strictly-avoid-quickref.md`
  - `.ai-memory/memory/00-project-governance-and-invariants.md`
- **Status**: READY
- **Target Release**: GitMap v6.590.0

---

## 1. Context & Objectives

The user issued an explicit, non-negotiable instruction:
> "Never use or add the README.text for adding the timestamp. If that memory exists, correct it everywhere, not only in GitMap but also in the coding guideline. Make sure README text never exists. Find out how many places the code is not following the coding guideline."

Past runs of autonomous AI agents periodically hallucinated generating `readme.text` (or `readme.txt`) containing dynamic date/time strings (e.g. Asia/Kuala_Lumpur or UTC timestamps) or proposed scripts like `refresh-readme.ps1` to touch readme files. While `.ai-memory/strictly-avoid.md` banned `readme.txt` timestamps, it did not explicitly name `readme.text` and contained ambiguous language permitting manual edits.

### Key Objectives:
1. **Universal Ban on `readme.text`**: The filename `readme.text` (and any casing: `README.TEXT`, `Readme.text`, `README.TXT`) is strictly prohibited repository-wide. All documentation must strictly use `readme.md`.
2. **Absolute Ban on Timestamp Injection**: No script, automation, command, hook, or skill may ever write, append, suggest, or format timestamps, dates, or clocks into any `readme.*` file.
3. **Memory & Specification Synchronization**: Review and update `.ai-memory/strictly-avoid.md`, `02-spec/01-spec-authoring-guide/11-exceptions.md`, `02-spec/17-consolidated-guidelines/03-strictly-avoid-quickref.md`, and `.ai-memory/memory/00-project-governance-and-invariants.md`.
4. **Coding Guidelines Verification**: Audit repository adherence to coding standards (positive booleans, 15-line limit, zero nested ifs, strictly relative paths, error return sovereignty).

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Repository-Wide Scan & Inventory

1. Search the repository for all occurrences of `readme.text` and `readme.txt` across specifications and memory files.
2. Verify that zero active source files or scripts generate or reference `readme.text`.
3. Catalog every governance and memory document where prohibition rules must be updated.

---

### Step 2.2: Update `.ai-memory/strictly-avoid.md`

1. Retitle the prohibition section:
   ```markdown
   ## readme.text and readme.txt timestamp generator — TOTAL BAN
   ```
2. Broaden the prohibition clause:
   - Explicitly name `readme.text` alongside `readme.txt`.
   - Prohibit creating, suggesting, proposing, designing, or specifying any feature, script, or hook that writes a timestamp, date, or time into `readme.text`, `readme.txt`, or any documentation file.
   - Clarify that `readme.text` must never exist under any circumstance; all documentation files must use canonical lowercase `readme.md`.
   - Remove any legacy ambiguous phrasing that permitted manual edits of `readme.text`.

---

### Step 2.3: Update `02-spec/01-spec-authoring-guide/11-exceptions.md`

1. Update Sequence 01:
   ```markdown
   ### Sequence 01 — readme.text & readme.txt time suggestions
   
   - Never suggest, recommend, hint at, ask about, or auto-add anything related to timestamps inside readme.text or readme.txt.
   - The file readme.text is completely forbidden; all documentation files must strictly use canonical readme.md.
   - Never mention "git update time", commit time, last-modified time, scheduled regeneration, cron jobs, or any time-flavoured workflow in or about readme files.
   ```

---

### Step 2.4: Update `02-spec/17-consolidated-guidelines/03-strictly-avoid-quickref.md`

1. Update Item 01 in the Quick Reference table:
   - Retitle prohibition to: `readme.text & readme.txt time suggestions`.
   - Enforce total prohibition: "Never suggest/recommend/hint/auto-add timestamps, 'git update time', last-modified metadata, or time-flavoured workflows in or about readme.text or readme.txt. The filename readme.text is completely banned. All documentation must use readme.md."

---

### Step 2.5: Update `.ai-memory/memory/00-project-governance-and-invariants.md`

1. Add explicit invariant section:
   ```markdown
   ### Universal Readme Casing & Timestamp Invariant
   - Documentation files must strictly use canonical lowercase `readme.md`.
   - Files named `readme.text` or `readme.txt` are completely prohibited.
   - Automatic timestamp generation, date injection, or time-stamping scripts are permanently banned.
   ```

---

### Step 2.6: Coding Guidelines Audit Across Modified Packages

1. Audit modified packages (`cli/cmddb/`, `cli/cmdnodes/`, `cli/cmdos/`, `cli/cmdfix/`, `cli/gitutil/`) against the 5 core quality pillars:
   - **Pillar 1 (Positive Booleans)**: Verify booleans use `is*` or `has*` only. Zero `can`, `should`, `was`, `will`, or negative booleans.
   - **Pillar 2 (Function Length)**: Verify all functions <= 15 lines (<= 8 lines preferred).
   - **Pillar 3 (Zero Nested Ifs)**: Verify all conditionals use guard clauses and early returns.
   - **Pillar 4 (Strictly Relative Paths)**: Verify 100% relative Git paths (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero `file:///` or drive letters.
   - **Pillar 5 (Error Return Sovereignty)**: Verify leaf functions return direct errors; no calls to `cliexit.HandleError` or `os.Exit`.

---

## 3. Coding Guidelines & Invariant Verification

- [ ] **Prohibition Clarity**: Explicit total ban on `readme.text` and timestamp generators in all governance files.
- [ ] **Path Hygiene**: Strictly relative paths in all updated markdown documents.
- [ ] **Documentation Integrity**: Preserve existing numbered sequence rules and cross-references.
- [ ] **Audit Completeness**: Verify zero uncorrected positive mentions of `readme.text` across the codebase.

---

## 4. Verification & Validation Commands

1. **Verify No Positive Usage of `readme.text`**:
   Search across `.ai-memory/` and `02-spec/` to confirm that any match for `readme.text` is strictly a prohibition statement:
   ```bash
   gitmap aum search "readme.text"
   ```
2. **Verify All Documentation Uses `readme.md`**:
   Confirm that all new and existing documentation files adhere to lowercase `readme.md`.

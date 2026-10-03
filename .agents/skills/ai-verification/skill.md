---
name: ai-verification
description: >-
  Autonomously conduct retrospective code and spec quality audits across recent tasks, checking acceptance criteria, coding guidelines, relative path hygiene, and CI/CD status.
---

# Retrospective AI Verification & Code Quality Audit

Use this skill when tasked with performing a retrospective verification of recent tasks (the last 2–3 completed tasks or 30–40 minutes of work). It audits generated specifications for acceptance criteria, inspects completed subtasks, verifies coding guideline compliance, and monitors CI/CD health via GitMap telemetry.

**Canonical Prompt:** `01-prompts/25-ai-verification/01-retrospective-ai-verification.md`  
**Automation Engine:** `python 03-ai-scripts/47-retrospective-ai-verification.py`  

## Quick Start & Execution

```bash
# Run automated retrospective verification across recent 3 tasks
python 03-ai-scripts/47-retrospective-ai-verification.py --tasks 3 --since-minutes 34

# Output full JSON diagnostic report
python 03-ai-scripts/47-retrospective-ai-verification.py --format json
```

## Audit Dimensions

1. **Specification Acceptance Criteria:**
   - Verify every new spec file under `02-spec/21-app/` contains a non-empty `## Acceptance Criteria` section with binary checkboxes (`- [ ]`).
2. **Coding Guideline Hygiene:**
   - **Booleans:** Implicit checks only (`if isReady`), zero `== true` / `== false`, zero mixed polarity (`if isA && !isB`).
   - **Paths:** Strict relative paths starting from git root; zero absolute paths (no URI schemes, no drive letters).
   - **Vertical Line Spacing:** Blank line before `if`, blank line after `}`, blank line before `return`.
3. **CI/CD Pipeline Telemetry:**
   - Inspect live pipeline health using `gitmap pe -t`.
   - If failures exist, extract exact failure lines via `gitmap pe` and perform surgical 4-part RCA remediation.

## Verification Checklist

- [ ] Recent 2–3 commits retrieved and file inventory cataloged.
- [ ] All touched specifications audited for `## Acceptance Criteria`.
- [ ] No absolute paths or `file:///` URIs present in touched files.
- [ ] No explicit boolean `== true` or mixed polarity conditions.
- [ ] CI/CD pipeline telemetry is verified green.

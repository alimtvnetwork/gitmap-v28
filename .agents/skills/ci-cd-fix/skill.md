---
name: ci-cd-fix
description: >-
  Use this skill to autonomously diagnose, fix, and verify CI/CD pipelines using local runner scripts, 4-part RCA, and self-looping.
---

# Instruction (must follow): Autonomous CI/CD Fix Loop (with Local Runner & RCA)

Trigger Keywords & Aliases: `fix with RCA`, `fix`, `fix, fix`, `CI/CD fix`, `cicd fix`

```text
N = 200
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to reading all CI/CD configuration files and generating `.lovable/ai-fix-scripts/03-cicd-local-runner.py` with native host commands (Docker stripped out).
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to running that script in an autonomous loop, applying surgical fixes to the codebase after each failure, until exit code = 0.
- [ ] /learn Ingest `.lovable/cicd-issues/`, `.lovable/strictly-avoid.md`, and `spec/03-error-manage/` before touching any code so past mistakes are never repeated.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

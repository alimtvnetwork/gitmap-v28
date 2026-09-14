---
name: scripts-fixer-alignment-and-pull-progress-bar
description: Autonomously align GitMap with scripts-fixer profiles (git-compact, dev profile, antigravity) and enhance git pull / git pull all UI with animated progress bar.
---

# Scripts Fixer Alignment & Git Pull Progress Bar Suite

## Objective
Autonomously align GitMap profile installations and Antigravity tooling with recent scripts-fixer enhancements, integrate git-compact into the dev profile and profile installation workflows, and enhance the GitMap pull UI (single-repo and multi-repo / pull all) with rich progress bar tracking.

## Key Capabilities
1. **Scripts-Fixer Intelligence & Profile Alignment:**
   - Ingest scripts-fixer repository changes across the last 50 commits.
   - Align profile installations, including git-compact, dev profile, simple-dev, and Antigravity installation fixes.
   - Ensure profile installation idempotency and tree rendering.

2. **Git Pull UI & Progress Bar:**
   - Enhance gitmap pull and gitmap pull --all / git pull all UI.
   - Provide visual progress bars, percentage indicators, repository step streaming, and clean summary tables.

3. **Strict Compliance & Architecture:**
   - Functions <= 8-15 lines.
   - Affirmative booleans (is*, has*).
   - Universal *apperror.AppError return wrapping.
   - Strict Unix LF line endings.
   - Ban on running go test, go build, or local runner during loops.

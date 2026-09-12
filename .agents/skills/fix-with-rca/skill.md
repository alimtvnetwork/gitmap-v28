---
name: fix-with-rca
description: Diagnose, analyze, and fix bugs using grounded 4-part Root Cause Analysis (RCA) without guessing.
---

# Fix with Root Cause Analysis (RCA)

Solves bugs and pipeline failures through structured 4-part Root Cause Analysis.

## 4-Part RCA Structure
1. **Symptom:** Raw error output, failing command, and reproduction steps.
2. **Root Cause:** One-sentence root cause identifying the exact mechanism of failure.
3. **Resolution:** Direct minimal fix applied to code or configuration.
4. **Prevention & Learnings:** Specific rule or avoidance logged to prevent recurring failures.

## Logging
- General bugs: `.lovable/issues/01-<slug>.md`
- CI/CD failures: `.lovable/cicd-issues/01-<slug>.md`

# Spec 278: GitMap Pull, Fix Remediation & Architecture Optimization

## Overview
Comprehensive remediation, parallelization, and modernization across:
1. Git Credential Store fix (`wincred` -> `wincredman`) eliminating Windows GCM fatal errors.
2. `gitmap reclone` URL and single repository support.
3. Generic typed worker pool for `gitmap fix --all` with CPU-adaptive concurrency (default 4 up to 20).
4. `gitmap pull-all` concurrency auto-scaling up to CPU cores, stream buffer preservation for color/progress, and redesigned post-pull resolution menu.
5. `gitmap update` daily probing throttle cache, version transition UI, title casing, and left padding.
6. `repo-cache` repository integration, multi-AI memory storage, and LLM Muse training migration from `repo-secrets` to `repo-cache`.
7. Subcommand normalization (hyphenated aliases eliminating literal space requirements).
8. GitMap Heatmap command documentation in `.agents/skills/gitmap/SKILL.md` and coding guidelines.

## Specifications
- [01-scope-and-plan.md](01-scope-and-plan.md)
- [02-acceptance-criteria-and-evidence.md](02-acceptance-criteria-and-evidence.md)

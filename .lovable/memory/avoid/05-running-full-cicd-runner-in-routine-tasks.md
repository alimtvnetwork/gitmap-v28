# Strictly Avoid: Running Full CI/CD Runner During Routine Tasks

- **Decision Date:** 2026-09-12
- **Context:** Agents executing coding guideline tasks (`15-cg-execute/*`) and routine loops were invoking `python 03-ai-scripts/06-cicd-local-runner.py --force --no-tests` on every micro-turn, executing 28+ quality gates across thousands of repository files, causing 30-70+ second delays and flagging unrelated test files.

## Invariant Rules
1. **Total Ban in Routine Turns:** Never run `06-cicd-local-runner.py` during routine development, single-file edits, or coding guideline batches.
2. **Targeted Verification Only:** Use targeted single-file linters (e.g. `check-nested-ifs.py <file>`, `check-boolean-guidelines.py <file>`, `08-naming-autofixer.py <file>`).
3. **Runner Whitelist:** The full CI/CD runner is strictly reserved for:
   - Explicit commands from the repository owner.
   - Dedicated CI/CD repair tasks (`ci-cd-fix`).
   - Final pre-release ceremonies (`--run-tests`).
4. **Runner Debounce Caching:** The runner caches execution results in `.lovable/cicd/last_run_cache.json`. Invocations within 15 seconds return `'Here is the result from the previous run.'` with cached status.

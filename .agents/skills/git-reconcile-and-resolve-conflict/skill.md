---
name: git-reconcile-and-resolve-conflict
description: Autonomously reconcile diverged Git branches, mechanically resolve merge conflicts across known file patterns, and safely push to remote tracking branches without rewriting Git history.
---

# Git Reconciliation & Mechanical Conflict Resolver

Autonomously resolves git branch divergence, applies deterministic domain-specific conflict resolution rules, and pushes to remote tracking branches without rewriting published Git history.

## Core Capabilities

1. **Non-Destructive Synchronization:**
   - Standard merge semantics (`git merge --no-ff`) strictly avoiding force-pushing, rebasing, or squashing published commits (preserving Lovable and CI/CD pipelines).
2. **Domain-Specific Mechanical Conflict Solvers:**
   - `.gitignore`, `.gitattributes`, `.npmignore`, `.dockerignore`: Deduplicated line union.
   - `version.json`, `package.json`, `prompt-version.template.json`: Monotonic SemVer comparison (highest version wins).
   - `changelog.md`, `release-notes.md`: Chronological release section preservation.
   - Markdown task lists & checklists: Deduplicated union with `- [x]` completion dominance.
   - Python `requirements*.txt`: Package union and higher version constraint arbitration.
   - Declarative code imports (`.py`, `.ts`, `.js`, `.go`, `.rs`): Lexicographical union and deduplication.
3. **Safety Verification & Syntax Invariants:**
   - Verifies zero remaining conflict markers (`<<<<<<<`, `=======`, `>>>>>>>`).
   - Runs pre-commit AST and schema syntax validation (Python `ast.parse()`, JSON `json.loads()`).
   - Pushes cleanly to remote tracking branch and verifies 0 ahead/behind delta post-push.

## CLI Execution

```bash
# 1. Run internal mechanical conflict resolver self-tests
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --self-test

# 2. Check divergence status and simulate merge (dry-run)
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --check

# 3. Reconcile diverged branches, mechanically resolve conflicts, and push to origin
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --push

# 4. Mechanically resolve active conflict markers in working directory
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py resolve-conflicts --strategy smart

# 5. Programmatic JSON output for agent automation
python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --check --json
```

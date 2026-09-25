---
name: prompts-and-skills-sync
description: Autonomously synchronize canonical prompts (01-prompts/), Antigravity skills (.agents/skills/), and AI scripts (03-ai-scripts/) from coding-guidelines across all connected repositories using 03-ai-scripts/38-sync-prompts-skills-scripts.py.
---

# Prompts, Skills & AI Scripts Multi-Repository Synchronizer

> **/goal** Propagate canonical prompt templates, Antigravity skills, and AI automation scripts across all connected target repositories, ensuring 100% parity and hygiene.
> **/learn** Execute `03-ai-scripts/38-sync-prompts-skills-scripts.py` safely with git pre-checks, branch status verification, and automated downstream commits.

**Version:** 1.0.0
**Updated:** 2026-09-24
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. When to Use

Activate this skill when:
- New prompts are added or updated in `01-prompts/`.
- New Antigravity skills are authored or updated in `.agents/skills/`.
- New shared scripts are added in `03-ai-scripts/` or `.agents/scripts/`.
- The user requests cross-repository synchronization or prompt deployment.

---

## 2. Connected Target Repositories (13 Repositories)

The synchronizer mirrors assets to the 13 connected repositories in the workspace parent directory (`d:\work\`):

1. `antigravity-manager`
2. `spec-builder`
3. `movie-cli`
4. `macro-ahk`
5. `laravel-automation`
6. `lara-publishing`
7. `lara-licensing`
8. `gitmap`
9. `wp-exam`
10. `wp-git-log`
11. `wp-html-automate`
12. `wp-link-manager`
13. `wp-onboarding`

---

## 3. Synchronized Directories

- `01-prompts/` -> `01-prompts/` (Including `v1/` classic prompts and `v2/` GitMap AUM accelerated prompts)
- `.agents/skills/` -> `.agents/skills/`
- `03-ai-scripts/` -> `03-ai-scripts/`
- `.agents/scripts/` -> `.agents/scripts/`

---

## 4. Execution Workflow & Multi-Repo Release Ceremony

For each target repository, execute this mandatory safety & release sequence:

### Step 1: Pre-Flight Pull & Safety Backup Branch
```bash
git checkout main
git pull origin main
# Create and immediately push safety backup branch
$ts = Get-Date -Format "yyyyMMdd-HHmmss"
git branch "backup/sync-$ts"
git push origin "backup/sync-$ts"
```

### Step 2: Dedicated Work Branch
```bash
git checkout -b feat/sync-prompts-v1-v2-skills
```

### Step 3: Mirror Synchronized Assets
Run the synchronization script from `coding-guidelines`:
```bash
python 03-ai-scripts/38-sync-prompts-skills-scripts.py
```

### Step 4: Atomic Commit
```bash
git add 01-prompts/ .agents/skills/ 03-ai-scripts/ .agents/scripts/
git commit -m "chore(sync): update prompts v1/v2, skills, and ai scripts"
git push origin feat/sync-prompts-v1-v2-skills
```

### Step 5: Full Release Ceremony
Bump version, generate tag, create release branch, and push all artifacts:
```bash
# Option A: GitMap automated release
gitmap release --bump patch -y

# Option B: Python release orchestrator
python 03-ai-scripts/29-release-orchestrator.py --tier patch
```

---

## 5. Non-Negotiable Quality Gates

- [ ] All mirrored skill files must use strictly lowercase `skill.md` filenames.
- [ ] No temporary files (`.pyc`, `.tmp`, `__pycache__`) are synced.
- [ ] Every target repository MUST have a backup branch pushed before any modifications.
- [ ] Every target repository commit message must follow standard semantic convention (`chore(sync): update prompts v1/v2, skills, and ai scripts`).
- [ ] All branches (work branch, release branch, backup branch) and tags MUST be pushed to remote.
- [ ] Never force-push or overwrite published git history.

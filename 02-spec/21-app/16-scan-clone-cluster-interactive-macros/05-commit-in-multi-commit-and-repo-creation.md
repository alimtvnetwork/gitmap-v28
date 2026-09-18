# Specification 16 - Part 5: Multi-Commit Sequential Replay Engine & Automatic Repository Scaffolding

## 1. Multi-Commit Sequential Processing Architecture

### 1.1 True Multi-Commit Replay Engine

`commit-in`, `commit-left`, and `commit-right` MUST execute true sequential multi-commit replay pipelines rather than squashing changes into a single monolithic commit:
- Reads commit manifests or diff logs containing multiple commit steps.
- For each discrete commit step:
  1. Stages modified / added / deleted files for that step.
  2. Synthesizes or preserves the distinct commit message, author metadata, and timestamp.
  3. Executes atomic `git commit`.
  4. Records checkpoint in the SQLite audit runlog (`commitin_runlog`).
- On failure of any intermediate commit, state is preserved cleanly with replay resume capabilities.

### 1.2 Automatic Repository Initialization & Scaffolding

When the target directory does not yet exist or is not a git repository:
- If `--create-repo` (or automatic mode) is enabled:
  1. Creates target directory tree.
  2. Runs `git init -b main`.
  3. Creates default `.gitignore`, `.gitattributes`, and `README.md`.
  4. Creates initial root commit `chore: initialize repository`.
  5. If remote origin URL is provided or GitHub CLI is configured, links remote and prepares branch tracking.
  6. Proceeds with sequential commit application.

## 2. Commit-In & SEO-Write Discovery & Examples

### 2.1 CLI Search Discovery

Searching `gitmap help --filter commit-in` or `gitmap help --filter commit` must display complete multi-commit workflows:
```bash

# Replay commits from source into target repo sequentially

gitmap commit-right --source ./proto-v1 --target ./prod-v2 --all

# Transfer discrete commits with automated changelog reconstruction

gitmap commit-in --manifest ./commits.json --target ./new-service --create-repo
```

### 2.2 SEO-Write Automated Commit Generation

- `gitmap seo-write` generates realistic time-staggered commits following specified templates, CSV keywords, and author signatures.
- Includes rich built-in templates (`seo-templates.json`) for web apps, APIs, documentation, and tooling packages.

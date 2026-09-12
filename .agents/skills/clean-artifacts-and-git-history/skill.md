---
name: clean-artifacts-and-git-history
description: Safely clean build artifacts, test outputs, and temporary files while preserving git hygiene.
---

# Clean Artifacts and Git History

Enforces repository cleanliness and guards against accidental commit of generated files.

## Actions
- Clean pycache, build artifacts, test logs, coverage dumps.
- Run `python 03-ai-scripts/19-artifact-remover.py`.
- Ensure `.gitignore` rules cover all newly introduced intermediate files.
- **Consolidated Atomic Commits:** NEVER make piecemeal 1-2 file commits. Commit all modified files and plans together as a single atomic unit.
- **Immediate Git Push:** ALWAYS push immediately to GitHub (`git push origin <branch>`) after creating any commit.
- **No Routine Builds:** NEVER execute full builds (`npm run build`, `go build ./...`) or `06-cicd-local-runner.py` during routine cleanup.

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

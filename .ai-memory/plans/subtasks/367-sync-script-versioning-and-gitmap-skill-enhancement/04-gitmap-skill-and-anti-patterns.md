# Subtask 04: GitMap Skill Enhancement, Anti-Pattern Ban & Alternative Reference Matrix

- **Task Code:** `TASK-04`
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 04`
- **Status:** completed
- **Owner:** AI Skill Architect
- **Specification:** [04-gitmap-skill-and-anti-patterns-spec.md](../../../02-spec/21-app/367-sync-script-versioning-and-gitmap-skill-enhancement/04-gitmap-skill-and-anti-patterns-spec.md)

## Actions
1. Update `.agents/skills/gitmap/skill.md` with:
   - Restrict condition ban and V7 engine enforcement rules.
   - Comprehensive Anti-Pattern Ban Matrix covering:
     * `powershell -Command "Get-Content ... -Tail N | Select-String ..."`
     * `gitmap py -c "import time; time.sleep(N)"`
     * `Get-ChildItem -Path ... | Select-Object ...`
     * `node -e "const fs = require('fs'); ..."`
     * `gitmap py -c "import os; [print(...) for ... in os.walk(...)]"`
     * `gitmap py -c "print(open('...').read())"`
   - Side-by-side GitMap command alternatives.
2. Mirror updates to `coding-guidelines/.agents/skills/gitmap/skill.md`.

## Quick Install v6.422.0

### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.422.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.422.0"
```

### Linux / macOS
```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.422.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.422.0"
```

## What's Changed in v6.422.0
- Added `gitmap gitignore agm` (`gitmap gitignore agy`, `gitmap agm`, `gitmap agy gitignore`) command to untrack (`git rm --cached`), delete, add `.antigravity_resume_task.json` and `antigravity-resume_task.json` to `.gitignore`, and commit (`Update .gitignore`).
- Added detection and interactive prompt during `gitmap scan` and `gitmap pull-all` (`pa`, `ta`, `pull-all-efficient`, `pae`) when `.antigravity_resume_task.json` is present or tracked in a repository.
- Added `.antigravity_resume_task.json` and `antigravity-resume_task.json` to `common.gitignore` and release `.gitignore` defaults.

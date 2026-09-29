## Quick Install v6.392.0

### Windows (PowerShell)

```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.392.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.392.0"
```

### Unix / Linux / macOS (Bash)

```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.392.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.392.0"
```

---

## What's Changed in v6.392.0

### Added
- `gitmap ls preview` / `--preview` / `-p` / `--gap`: Numbered repository list with repo name and path on the next line, separated by a single blank line gap between entries.
- `gitmap ls tree` / `--tree` / `-t`: Groups repositories by parent directory in an emoji tree format (`📁 <parent-folder>`).
- `gitmap folder-tree` (`ft`, `foldertree`): New command to scan and visualize folder and repository structures on any path on disk without prior database scanning. Includes sequence numbers, Git detection (`[git: <branch>]`), tree/preview rendering, and multi-format export/import (`json`, `yaml`, `tree`, `preview`) with directory and placeholder file scaffolding.

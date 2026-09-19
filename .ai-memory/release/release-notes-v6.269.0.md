## Quick Install v6.269.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.269.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.269.0/install.sh | bash
```

## Changelog v6.269.0

- Remove raw direct git command lines (`git -C <repo> ...`) from remediation plan display in `gitmap fix`
- Replace direct git plan string with human-readable strategy description in `executeFixRecipe`
- Clean step command display in `executeSingleStep` to omit `-C <repoPath>` and show concise actions (`git stash -u`, `git pull`, `git stash pop`)
- Clean remediation CLI help text and examples shown after `gitmap pull-all`
- Add unit tests for `formatStepCommand` in `cli/cmd/fix_execute_test.go`
- Update `cli/helptext/fix.md` examples to match clean plan output

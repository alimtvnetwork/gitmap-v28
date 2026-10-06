# GitMap v6.497.0

## What's Changed in v6.497.0

- **ubuntu ide and github desktop scan sync suite**:
  - Grounded 4-part RCA in `02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md` identifying root cause of repo omission on Ubuntu.
  - Multi-IDE standalone synchronization toolchain in `03-ai-scripts/40-ubuntu-ide-desktop-sync.py` with automatic directory creation (`0755`) and dedup checks.
  - Dedicated `gitmap ide` command group (`gitmap ide add`, `gitmap ide sync`, `gitmap ide remove`, `gitmap ide list`, `gitmap ide status`, `gitmap ide help`) in `cli/cmdide/`.
  - Scan integration in `cli/cmdscan/` with `--sync-ide`, `--skip-sync`, and `--exclude-sync` flags.
  - Auto-mkdir `EnsureProjectsJSONPath()` in `cli/vscodepm/path.go` preventing silent omission on clean systems.
  - Linux binary candidate resolution in `cli/desktop/resolve.go` supporting package managers and desktop paths.
  - Cluster pending commits and sends suite (`gitmap pending-commits`, `gitmap sends`, `gitmap nodes pc`, `gitmap nodes commits/cpf/cpb/cpr/commit-fix`).

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.497.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.497.0/install.sh | sh
```

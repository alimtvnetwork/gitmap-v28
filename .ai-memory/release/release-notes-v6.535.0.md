# GitMap v6.535.0

## What's Changed in v6.535.0

- **Repo-Cache JSON Injection**: Companion `repo-cache` metadata dynamically probed via Split-DB and injected into `gitmap.json` manifests (`cli/cmdscan/scanoutput.go`), ensuring workspace cache directories are reliably restored on subsequent clones.
- **Two-Tier Subtree Solutions Formatting**: Refactored `gitmap pull` and `gitmap pc` error rendering (`cli/cmdpull/pull_efficient_render.go`) to render hierarchical two-tier subtrees with native GitMap commands (`gitmap sends cp`, `gitmap stash`, `gitmap fix`, `gitmap pull`), eradicating raw `git -C ...` strings and adding fleet multi-repo Option 3 (`gitmap sends cp all "wip: save changes"`).
- **Antivirus False-Positive Heuristics Mitigations**: Eliminated plaintext `DefaultPassword` registry writes in `cli/cmdos/os_autologin_windows.go` (eradicating `Trojan:Win32/AutoLogon`) and removed telemetry service tampering in `cli/cmdos/os_tweak_windows_privacy.go` (eradicating EDR T1562.001 trigger).
- **Golang Direct Query CLI Acceleration**: Implemented native `<5ms` path query flags (`gitmap rc --path`, `gitmap rs --path`, `gitmap path <repo>`), eradicating the slow Python subprocess wrapper anti-pattern.
- **Skills & Coding Guidelines Synchronization**: Updated `.agents/skills/gitmap/SKILL.md` and `02-spec/02-coding-guidelines/` enforcing native Go query commands and prohibiting Python subprocess chaining.

### Specifications & Documentation
- Architecture Spec: `02-spec/21-app/gitmap-repository-management-and-release-process/01-architecture-spec.md`
- Component Spec: `02-spec/21-app/gitmap-repository-management-and-release-process/02-component-spec.md`
- Antivirus Investigation: `02-spec/21-app/gitmap-repository-management-and-release-process/03-antivirus-investigation.md`
- Skills & Guidelines: `02-spec/21-app/gitmap-repository-management-and-release-process/04-skills-and-guidelines-spec.md`

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.535.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.535.0/install.sh | sh
```


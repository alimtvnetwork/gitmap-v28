# 08-distribution-and-release: Release Components, Runners & Installer Specification

- **Spec ID:** `08-distribution-and-release/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Release Automation, Cross-Platform Runners, NSIS Installer, Fleet Distribution
- **Dependencies:** `cli/release`, `cli/installer`, `scripts/`, `03-ai-scripts/`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Distribution and Release cluster comprises four primary operational components:

```
08-distribution-and-release/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Release Ceremony Orchestrator** | `cli/release/`, `03-ai-scripts/37-bump-version.py` | Version calculation, changelog generation, git tag creation, manifest sync. |
| **Cross-Platform Runners** | `run.ps1`, `run.sh`, `run.config.json` | Declarative execution runners with pre-flight environment verification. |
| **Windows NSIS Installer** | `scripts/installer.nsi`, `scripts/install.ps1` | Payload detection, installation to `%LOCALAPPDATA%\GitMap`, PATH injection. |
| **Linux Packaging & Installers** | `scripts/local-install.sh`, `scripts/install.sh` | Linux package extraction, binary symlinking to `/usr/local/bin`, uninstall CLI. |

---

## 2. Release Ceremony Automation

### 2.1 Release Bump Lifecycle
The release ceremony is invoked via GitMap release scripts:
1. **Pre-flight Checks:** Verifies git working tree is clean and runs static checks (`check-relative-paths.py`, `check-forbidden-strings.py`).
2. **Version Bump:** Increments SemVer in `version.json` based on type (patch: `v6.498.1`, minor: `v6.499.0`, major: `v7.0.0`).
3. **Manifest Synchronization:** Updates:
   - `version.json`
   - `readme.md` (version badge, quick install snippets)
   - `cli/constants/version.go`
   - `changelog.md` (extracts recent atomic commits)
4. **Git Tagging:** Creates annotated tag `vX.Y.Z` with commit summary.

---

## 3. Cross-Platform Execution Runners (`run.ps1` / `run.sh`)

### 3.1 `run.config.json` Specification
The execution runners are configured declaratively via `run.config.json`:

```json
{
  "name": "gitmap",
  "defaultTask": "test",
  "tasks": {
    "build": {
      "command": "go build -o bin/gitmap ./cli",
      "description": "Compiles GitMap CLI binary"
    },
    "test": {
      "command": "go test ./cli/...",
      "description": "Runs complete Go unit test suite"
    },
    "lint": {
      "command": "python3 linter-scripts/check-relative-paths.py",
      "description": "Validates relative path integrity"
    }
  }
}
```

---

## 4. Packaging & Installer Mechanics

### 4.1 NSIS Windows Installer
- **Default Path:** `%LOCALAPPDATA%\Programs\GitMap` (user mode, requires no admin elevation).
- **Environment Integration:** Injects GitMap into user `PATH` and registers PowerShell predictive completion.
- **Uninstaller:** Removes files, cleans PATH entries, and clears temporary handoff files.

### 4.2 Linux Standalone Installer
- Fetches target architecture tarball (`gitmap-linux-amd64.tar.gz`).
- Unpacks to `~/.local/share/gitmap`.
- Creates symlink at `~/.local/bin/gitmap` or `/usr/local/bin/gitmap`.

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isReleaseCeremonyAutomated: true
  isRunConfigJsonParsedCleanly: true
  isInstallersCrossPlatformTested: true
  isPositiveBooleansUsed: true
```

- [x] Version bump synchronizes all manifests atomically.
- [x] `run.ps1` and `run.sh` execute default tasks reliably.
- [x] Package manifests generate without syntax or path errors.

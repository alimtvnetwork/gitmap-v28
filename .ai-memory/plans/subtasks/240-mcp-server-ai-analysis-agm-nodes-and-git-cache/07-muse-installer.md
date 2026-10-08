# Subtask 07: Meta Muse AI Companion Installer Subsystem

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `07-muse-installer`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-component-and-cli-spec.md` (Section 3)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **Meta Muse AI Development Companion Installer** (`gitmap install muse` and alias `gitmap muse install`).

### The Problem
AI companion tool installations require manual, platform-specific shell scripts, environment variable configurations, and verification steps. Developers working across Windows, Linux, and macOS frequently encounter broken PATH setups or missing dependencies, leading to developer friction and uneven agent environments.

### The Solution
Provide a seamless, unified installation wrapper inside GitMap:
- Automatic host OS detection (`runtime.GOOS`).
- Execution of verified official install channels:
  - **Windows:** `powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "irm https://dev.meta.ai/install.ps1 | iex"`
  - **Linux:** `bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"` (with `wget` fallback).
  - **macOS:** `bash -c "curl -fsSL https://dev.meta.ai/install.sh | bash"` (with architecture awareness).
- Post-install verification executing `muse --version` and resolving binary location.
- Visual parity and workflow alignment verified against reference screenshot: `assets/screenshots/240-muse-meta-installer.png`.
- Dry-run mode (`--dry-run`) allowing inspection of commands before execution.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmdinstall/install_muse.go` | New | Implements `RunInstallMuse()`, OS-specific script executors, PATH checking, and verification logic. |
| `cli/cmdinstall/install_muse_test.go` | New | Unit tests verifying command construction, OS branch selection, dry-run flags, and error parsing. |
| `cli/cmdinstall/install_aliases.go` | Modify | Wire `muse` sub-target into `gitmap install` argument dispatcher. |
| `cli/cmd/rootdispatch.go` | Modify | Register `muse install` and `install-muse` direct CLI aliases. |

---

## 3. Detailed Implementation Requirements

### 3.1 CLI Arguments & Flag Handling
```bash
gitmap install muse [--dry-run] [--verify-only] [--force]
gitmap muse install [--dry-run] [--verify-only] [--force]
```
- Flags:
  - `--dry-run` (`-n`): Output execution commands without running.
  - `--verify-only` (`-v`): Check if `muse` is already installed and runnable.
  - `--force` (`-f`): Reinstall even if already detected.
  - `--install-dir <path>`: Custom installation directory.

### 3.2 Dispatch Engine (`cli/cmdinstall/install_muse.go`)
```go
package cmdinstall

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunInstallMuse(args []string) error {
	opts, err := parseMuseInstallOptions(args)
	if err != nil {
		return err
	}

	if opts.IsVerifyOnly {
		return verifyExistingMuse()
	}

	fmt.Printf("%s● Installing Meta Muse AI Companion (%s)%s\n\n", constants.ColorCyan, runtime.GOOS, constants.ColorReset)
	start := time.Now()

	switch runtime.GOOS {
	case "windows":
		return installMuseWindows(opts, start)
	case "linux":
		return installMuseLinux(opts, start)
	case "darwin":
		return installMuseDarwin(opts, start)
	default:
		return apperror.NewSimple("unsupported operating system for muse install: "+runtime.GOOS, "E1088")
	}
}
```

### 3.3 Visual Parity Verification
- Installer terminal output must match layout in `assets/screenshots/240-muse-meta-installer.png`:
  - Blue/Cyan section header: `● Installing Meta Muse AI Companion`.
  - Progress spinner or status message during script execution.
  - Green completion checkmark: `✓ Meta Muse installed successfully`.
  - Version report and binary PATH verification.

---

## 4. Acceptance Criteria

- [ ] `gitmap install muse` and `gitmap muse install` dispatch correctly.
- [ ] Windows executes PowerShell `irm | iex` with `Bypass` policy.
- [ ] Linux and macOS execute curl/bash pipelines with appropriate fallbacks.
- [ ] `--dry-run` prints commands without executing them.
- [ ] Post-install verification checks binary presence and runs `muse --version`.
- [ ] Visual terminal presentation matches `assets/screenshots/240-muse-meta-installer.png`.
- [ ] Unit tests pass with >85% coverage.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests
go test -v ./cli/cmdinstall/... -run TestInstallMuse

# 2. Test dry-run mode on current workstation
gitmap install muse --dry-run

# 3. Test verify-only mode
gitmap install muse --verify-only

# 4. Guideline compliance check
python 03-ai-scripts/05-guideline-autofixer.py cli/cmdinstall --check-only
```

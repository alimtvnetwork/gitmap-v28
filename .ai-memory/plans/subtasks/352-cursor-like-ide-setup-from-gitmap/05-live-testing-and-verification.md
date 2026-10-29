# Subtask 05: Live End-to-End Verification, Cross-Platform Testing, and Test Suite Parity

> **Subtask ID:** 352-05  
> **Sequence ID:** 352  
> **Status:** PENDING  
> **Module Path:** `cli/cmdide/setup/`, `cli/cmdui/`, `.agents/skills/gitmap/`  
> **Parent Plan:** `.ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md`  
> **Specification Reference:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/readme.md`  

---

## 1. Objective & Scope

Execute comprehensive multi-tier verification and release procedures for the Cursor-like IDE setup feature:
1. Conduct complete end-to-end integration testing in `cli/cmdide/setup/e2e_test.go` covering preset resolution, dry-run simulation, pre-flight backup, atomic file mutation, and JSON output formatting.
2. Validate cross-platform path resolution across Windows, macOS, and Linux virtual environments.
3. Verify live connection testing against mock HTTP servers and validate real Moonshot Kimi endpoint compatibility.
4. Execute full coding guidelines audit across all modified and newly created files (<= 15 lines per function, zero nested ifs, positive booleans, leaf error sovereignty, typed exit codes).
5. Synchronize GitMap CLI skill documentation in `.agents/skills/gitmap/skill.md` with new alternative command recipes.
6. Conduct atomic minor version bump release ceremony via `gitmap cpr`.

---

## 2. End-to-End Integration Test Suite (`cli/cmdide/setup/e2e_test.go`)

Construct isolated integration test suite using temporary file systems:

```go
package setup_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
)

func TestE2E_FullSetupWorkflow(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)
	t.Setenv("APPDATA", filepath.Join(tempHome, "AppData", "Roaming"))

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": [{"id": "kimi-k3.1-pro"}]}`))
	}))
	defer mockServer.Close()

	args := []string{
		"all",
		"--provider", "kimi",
		"--base-url", mockServer.URL,
		"--api-key", "test-key-12345",
		"--test",
		"--backup",
	}

	err := cmdide.RunIDESetup(args)
	if err != nil {
		t.Fatalf("expected successful setup, got error: %v", err)
	}

	resolver := setup.NewPathResolver()
	continuePath, _ := resolver.GetContinueConfigPath()
	if _, err := os.Stat(continuePath); os.IsNotExist(err) {
		t.Errorf("expected continue config at %s, file missing", continuePath)
	}

	zedPath, _ := resolver.GetZedSettingsPath()
	if _, err := os.Stat(zedPath); os.IsNotExist(err) {
		t.Errorf("expected zed settings at %s, file missing", zedPath)
	}
}

func TestE2E_DryRunDoesNotMutate(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("USERPROFILE", tempHome)
	t.Setenv("HOME", tempHome)

	args := []string{
		"continue",
		"--provider", "kimi",
		"--dry-run",
	}

	err := cmdide.RunIDESetup(args)
	if err != nil {
		t.Fatalf("expected dry-run to succeed, got: %v", err)
	}

	resolver := setup.NewPathResolver()
	continuePath, _ := resolver.GetContinueConfigPath()
	if _, err := os.Stat(continuePath); !os.IsNotExist(err) {
		t.Errorf("expected dry-run NOT to write file, but file exists: %s", continuePath)
	}
}
```

---

## 3. Cross-Platform Verification Matrix

| Operating System | Target File | Expected Path Format | Verified In Test |
|:---|:---|:---|:---:|
| **Windows 10/11** | Cline MCP Settings | `%APPDATA%\Code\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json` | `path_resolver_test.go` |
| | Continue Config | `%USERPROFILE%\.continue\config.json` | `path_resolver_test.go` |
| | Zed Settings | `%APPDATA%\Zed\settings.json` | `path_resolver_test.go` |
| **macOS (Darwin)** | Cline MCP Settings | `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json` | `path_resolver_test.go` |
| | Continue Config | `~/.continue/config.json` | `path_resolver_test.go` |
| | Zed Settings | `~/.config/zed/settings.json` | `path_resolver_test.go` |
| **Linux (Ubuntu/Debian)** | Cline MCP Settings | `~/.config/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json` | `path_resolver_test.go` |
| | Continue Config | `~/.continue/config.json` | `path_resolver_test.go` |
| | Zed Settings | `~/.config/zed/settings.json` | `path_resolver_test.go` |

---

## 4. Coding Guidelines Compliance Audit Checklist

Every modified file in this feature MUST be strictly audited against repository invariants:

1. **15-Line Function Limit**:
   - `cli/cmdide/setup/presets.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/path_resolver.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/generator_cline.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/generator_continue.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/generator_zed.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/merger.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/backup.go`: All functions <= 15 lines.
   - `cli/cmdide/setup/connection_tester.go`: All functions <= 15 lines.
   - `cli/cmdide/setup_cmd.go`: All functions <= 15 lines.
   - `cli/cmdide/setup_flags.go`: All functions <= 15 lines.
   - `cli/cmdide/setup_view.go`: All functions <= 15 lines.
   - `cli/cmdui/ui_server_ide_setup.go`: All functions <= 15 lines.

2. **Positive Boolean Naming**:
   - Variables and struct fields use `is*` or `has*` prefixes exclusively:
     - `isDryRun`, `isQuiet`, `isJSON`, `isForceCreate`, `isTestEnabled`, `isBackupEnabled`
     - `isSuccess`, `isInstalled`, `hasConfig`, `isBackupCreated`, `hasStreamingSupport`
   - ZERO negative boolean identifiers (`no*`, `disable*`, `un*`, `isNot*`).

3. **Zero Nested If Statements**:
   - All nested conditionals replaced with guard clauses and early returns.

4. **Leaf Error Sovereignty**:
   - All leaf methods return `*apperror.AppError` or `error` directly.
   - Zero invocations of `os.Exit` or CLI exit handlers inside leaf functions.

5. **Strongly Typed Enums & Typed Exit Codes**:
   - Enums end in `Type` (`IDETargetType`, `IDEProviderType`).
   - CLI exits bubble up to command dispatchers using typed constants (`constants.ExitUsageError`, `constants.ExitGeneralError`).

---

## 5. Skill Documentation Synchronization (`.agents/skills/gitmap/skill.md`)

Add dedicated alternative command recipes for the Cursor-like setup in `.agents/skills/gitmap/skill.md`:

```markdown
### Cursor-Like AI IDE Setup Recipes

Configure open-source and native editors into Cursor-equivalent AI coding environments powered by Moonshot Kimi K3.1 Pro (or other LLMs):

```bash
# 1-Click universal setup for all editors using recommended Kimi K3.1 Pro preset
gitmap ide setup all --provider kimi

# Configure VS Code + Cline for Cursor Composer multi-file editing
gitmap ide setup cline --provider kimi

# Configure VS Code + Continue.dev for Cursor Tab autocomplete & inline edits
gitmap ide setup continue --provider kimi

# Configure Zed AI Editor for native ultra-fast BYOK coding
gitmap ide setup zed --provider kimi

# Dry-run preview of configuration diffs without writing to disk
gitmap ide setup all --provider kimi --dry-run

# Output structured JSON summary for automation pipelines
gitmap ide setup all --provider kimi --json

# Test provider credentials and measure endpoint latency before applying
gitmap ide setup all --provider kimi --test

# Launch the interactive Fleet Web UI to configure editors via graphical tabs
gitmap ui
```
```

---

## 6. Staged Minor Version Bump Release Ceremony

Following 100% green verification, execute the official GitMap release sequence:

1. **Local Test Suite Run**:
   ```bash
   go test -v ./cli/cmdide/... ./cli/cmdui/...
   ```
2. **Stage Commit for Implementation**:
   ```bash
   gitmap cpf "IDE - implement cursor-like ide setup with kimi k3.1 pro and open editors"
   ```
3. **Stage Commit for Skill Synchronization**:
   ```bash
   gitmap cpc "SKILL - document cursor-like ide setup commands in gitmap skill"
   ```
4. **Minor Version Bump & Tag Release Ceremony**:
   ```bash
   gitmap cpr "v6.630.0 cursor-like ide setup with kimi k3.1 pro and open editors"
   ```
5. **Verify Clean Working Tree**:
   ```bash
   git status
   ```

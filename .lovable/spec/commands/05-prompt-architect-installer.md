# Specification: Prompt Architect Installation Wrapper (gitmap ct install-prompts)

## 1. Overview
This specification defines the `gitmap ct install-prompts` command suite for downloading, updating, and verifying Prompt Architect v2 in any target repository.

## 2. Command Architecture
- **Top-level Command**: `gitmap ct install-prompts`
- **Aliases**: `gitmap ct install-prompt`, `gitmap ct update-prompts`, `gitmap install-prompts`, `gitmap ct prompts-status`, `gitmap ct prompts-version`.
- **Target Specifiers**: Supports explicit directory paths, repo slugs, numeric database IDs, or automatic child repository discovery in registered work directories.

## 3. OS-Aware Remote Installation Scripts
- **Unix / macOS / Linux**:
  `curl -sL https://raw.githubusercontent.com/alimtvnetwork/prompt-architect-v2/main/install.sh | bash`
- **Windows (PowerShell)**:
  `powershell -NoProfile -ExecutionPolicy Bypass -Command "Invoke-Expression "& { $(Invoke-RestMethod https://raw.githubusercontent.com/alimtvnetwork/prompt-architect-v2/main/install.ps1) }""`

## 4. Metadata in version.json
- Inspects and verifies `promptArchitectByRiseupAsia` or `prompt-architect` metadata section:
  ```json
  {
    "version": "6.105.0",
    "promptArchitectByRiseupAsia": {
      "version": "v2.0.0",
      "installed_at": "2026-08-26T17:30:00Z",
      "status": "active"
    }
  }
  ```

## 5. In-Place Update Protocol
- Allows seamless re-installation and updates without prompt loss.

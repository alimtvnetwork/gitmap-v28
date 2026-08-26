# PowerShell Runner Error Codes

> **Version:** 1.1.0  
> **Updated:** 2026-02-04  
> **Status:** Active

---

## Error Code Range

PowerShell integration uses error codes in the **9500-9599** range.

---

## Exit Codes

| Exit Code | Name | Description | Resolution |
|-----------|------|-------------|------------|
| 0 | SUCCESS | All steps completed successfully | N/A |
| 1 | ERR_PREREQUISITES | Prerequisites installation failed | Check winget availability |
| 2 | ERR_NPM_INSTALL | npm install failed | Clear npm cache, check network |
| 3 | ERR_NPM_BUILD | npm build failed | Check build errors in output |
| 4 | ERR_GO_RUN | Go backend failed to start | Check Go errors in output |
| 5 | ERR_CONFIG_MISSING | powershell.json not found | Create configuration file |
| 6 | ERR_CONFIG_INVALID | powershell.json parse error | Validate JSON syntax |
| 7 | ERR_PATH_NOT_FOUND | Configured path doesn't exist | Verify folder structure |
| 8 | ERR_COPY_FAILED | Build copy operation failed | Check disk space/permissions |
| 9 | ERR_GIT_FAILED | Git pull failed critically | Check git credentials |
| 10 | ERR_FIREWALL | Firewall configuration failed | Run as Administrator |

---

## Detailed Error Codes

### 9500 - Configuration Errors

| Code | Name | Message |
|------|------|---------|
| 9500 | ERR_CONFIG_NOT_FOUND | powershell.json not found in project root |
| 9501 | ERR_CONFIG_PARSE | Failed to parse powershell.json: {details} |
| 9502 | ERR_CONFIG_MISSING_FIELD | Required field '{field}' missing in powershell.json |
| 9503 | ERR_CONFIG_INVALID_PATH | Path '{path}' in config does not exist |
| 9504 | ERR_CONFIG_INVALID_PORT | Port '{port}' must be between 1 and 65535 |

### 9510 - Prerequisites Errors

| Code | Name | Message |
|------|------|---------|
| 9510 | ERR_WINGET_NOT_FOUND | winget is not available for auto-install |
| 9511 | ERR_GO_INSTALL_FAILED | Failed to install Go via winget |
| 9512 | ERR_NODE_INSTALL_FAILED | Failed to install Node.js via winget |
| 9513 | ERR_GO_NOT_IN_PATH | Go installed but not in PATH. Restart required |
| 9514 | ERR_NPM_NOT_IN_PATH | npm installed but not in PATH. Restart required |

### 9520 - Build Errors

| Code | Name | Message |
|------|------|---------|
| 9520 | ERR_NPM_INSTALL_FAILED | npm install failed with exit code {code} |
| 9521 | ERR_NPM_BUILD_FAILED | npm build failed with exit code {code} |
| 9522 | ERR_DIST_NOT_CREATED | Build completed but dist folder not found |
| 9523 | ERR_COPY_DIST_FAILED | Failed to copy dist to target directory |
| 9524 | ERR_CLEAN_FAILED | Failed to remove {path} during force clean |

### 9530 - Backend Errors

| Code | Name | Message |
|------|------|---------|
| 9530 | ERR_BACKEND_DIR_NOT_FOUND | Backend directory '{dir}' not found |
| 9531 | ERR_MAIN_GO_NOT_FOUND | main.go not found in backend directory |
| 9532 | ERR_GO_BUILD_FAILED | go build failed with exit code {code} |
| 9533 | ERR_GO_RUN_FAILED | go run failed with exit code {code} |
| 9534 | ERR_CONFIG_COPY_FAILED | Failed to copy config.example.json |

### 9540 - Firewall Errors

| Code | Name | Message |
|------|------|---------|
| 9540 | ERR_NOT_ADMIN | -OpenFirewall requires Administrator privileges |
| 9541 | ERR_FIREWALL_CMDLET | New-NetFirewallRule cmdlet not available |
| 9542 | ERR_FIREWALL_RULE_FAILED | Failed to create firewall rule for port {port} |

### 9550 - Git Errors

| Code | Name | Message |
|------|------|---------|
| 9550 | ERR_NOT_GIT_REPO | Not a git repository (skipping pull) |
| 9551 | ERR_GIT_PULL_FAILED | git pull failed with exit code {code} |
| 9552 | ERR_GIT_CONFLICT | git pull aborted due to conflicts |

---

## Error Handling Patterns

### PowerShell Error Handling

```powershell
$ErrorActionPreference = "Stop"

try {
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR [9521]: npm build failed with exit code $LASTEXITCODE" -ForegroundColor Red
        exit 3
    }
}
catch {
    Write-Host "ERROR [9521]: $($_.Exception.Message)" -ForegroundColor Red
    exit 3
}
```

### Graceful Degradation

```powershell

# Non-critical errors (warn but continue)

git pull
if ($LASTEXITCODE -ne 0) {
    Write-Host "WARNING [9551]: git pull failed, continuing anyway..." -ForegroundColor Yellow
    # Don't exit, continue with build
}
```

---

## Logging

### Console Output Levels

| Level | Color | Usage |
|-------|-------|-------|
| Info | Cyan | Step headers, status |
| Success | Green | Completion messages |
| Warning | Yellow | Non-fatal issues |
| Error | Red | Fatal errors |
| Debug | Gray | Timing, paths |

### Example Output

```
[3/5] Building React frontend...
  FORCE MODE: Cleaning build artifacts...
  Removing node_modules...
  Removing dist...
  ✓ Clean complete
  Installing npm dependencies...
  Running npm build...
  ✓ Frontend built successfully
  ⏱ 45.2s
```

---

## Cross-References

- [Script Reference](./02-script-reference.md) - CLI flags
- [Integration Guide](./03-integration-guide.md) - Setup steps
- [Error Management](../03-error-manage/00-overview.md) - Global error patterns

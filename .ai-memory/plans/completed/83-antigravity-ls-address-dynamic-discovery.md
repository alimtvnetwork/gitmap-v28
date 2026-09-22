# Master Spec & Completion: Antigravity Language Server Address Dynamic Discovery & agentapi Injection

## Problem Description & User Incident
When executing `gitmap agy pt "<text>"` from terminal or external scripts (e.g. `scripts-fixer`), the command failed with:
```text
Error: [E9031:EXECUTION] agentapi execution failed: {
  "error": "ANTIGRAVITY_LS_ADDRESS is not set"
}: (at=cmdagy/agy_agentapi.go:108)
```
The user reported:
> "Failed you didnt test it properly locally, why????"

## Root Cause Analysis
1. **Missing LS Address & CSRF Token**:
   - `language_server.exe agentapi` requires two environment variables when called: `ANTIGRAVITY_LS_ADDRESS` (host:port) and `ANTIGRAVITY_CSRF_TOKEN` (UUID token).
   - In terminal sessions outside of the Antigravity IDE, these environment variables are never exported to the OS environment.
   - `executeAgentAPICmd` in `cmdagy/agy_agentapi.go` executed `cmd.CombinedOutput()` with standard inherited OS environment, resulting in the fatal `ANTIGRAVITY_LS_ADDRESS is not set` failure.
2. **Multi-Port Probing Requirement**:
   - `language_server.exe` opens multiple TCP listening ports on `127.0.0.1` (e.g. HTTP bridge on port N and gRPC Language Server on port N+1).
   - Sending requests to the wrong port produces socket connection resets or HTTP preface errors.
   - Dynamic port selection must verify active gRPC handshake capability before binding.

## Accomplished Implementation
1. **Dynamic Language Server Discovery (`cli/cmdagy/agy_ls_discovery.go`)**:
   - **Process & Command-Line Inspection**: Inspects the running `language_server.exe` process via CIM (`Get-CimInstance Win32_Process`) and extracts the `--csrf_token <UUID>` from its command line.
   - **TCP Port Extraction & Probing**: Identifies all active listening TCP ports owned by the process on `127.0.0.1` via `Get-NetTCPConnection`.
   - **Safe Probe Handshake (`probeLSPort`)**: Probes ports using a fast `get-conversation-metadata probe-check`. Selects the port that successfully accepts the connection without socket reset or preface errors.
   - **Thread-Safe RW Caching**: Caches the discovered `ANTIGRAVITY_LS_ADDRESS` and `ANTIGRAVITY_CSRF_TOKEN` in memory with `sync.RWMutex` to avoid process query overhead on subsequent invocations.
   - **Cache Invalidation (`InvalidateAntigravityLSEnvCache`)**: Automatically invalidates the cached endpoint if an agentapi call fails, allowing automatic recovery on server restarts.
2. **Environment Variable Injection (`cli/cmdagy/agy_agentapi.go`)**:
   - Injects `ANTIGRAVITY_LS_ADDRESS` and `ANTIGRAVITY_CSRF_TOKEN` into `cmd.Env` during `executeAgentAPICmd`.
3. **Antigravity CLI Fallback (`cli/cmdagy/agy_agentapi_dispatch.go`)**:
   - Added secondary fallback to Antigravity CLI (`agy.exe -p <prompt>`) when the IDE process is offline (`pid <= 0`).
4. **Adhoc E2E Test Suite (`cli/cmdagy/agy_e2e_adhoc_test.go`)**:
   - Added `TestE2E_LSEnvDiscovery` and `TestE2E_AgentAPIExecution`.
   - Verified that `language_server.exe agentapi` discovers `127.0.0.1:54912` and executes successfully with zero errors.

## Verification & Test Results
- **Unit & Adhoc E2E Tests**:
  - `go test -v -count=1 -tags=e2e ./cmdagy -run "TestE2E_LSEnvDiscovery|TestE2E_AgentAPIExecution"` -> **PASS (0.126s)**
  - Full CLI live test: `go run . agy pt "Test prompt verification"` -> **PASS (Prompt injected into active Antigravity session)**
- **Coding Guideline & Linter Verification**:
  - `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/cmdagy` -> **PASS (All 148 files conform to boolean & newline guidelines)**
  - `python 03-ai-scripts/26-go-code-formatter.py cli/cmdagy` -> **PASS (148 Go files verified/formatted)**

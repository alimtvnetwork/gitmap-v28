# Subtask 05: Antigravity Multi-Instance Prompt Query API

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md)  
> **Status:** `QUEUED`  
> **Target Subsystems:**  
> - `cli/cmdagy/`  
> - `cli/cmdui/`  
> - `cli/store/`  

---

## 1. Technical Objective

Design and implement an instance-aware Antigravity prompt and conversation inspection engine capable of discovering and querying running prompts, in-queue prompts (`agy-prompt-queue.json`), and recent conversation dialogues across multiple concurrent Antigravity IDE instances. Expose this querying capability through:
1. CLI flags (`--instance <id>`, `--all-instances`, `--json`, `--full`, `--limit`, `--wc`) on `gitmap rp ls` and `gitmap agy prompts`.
2. REST API endpoints (`/api/instances` and `/api/prompts/instances`) mounted on the embedded Web UI server (`cli/cmdui/ui_server.go`).

---

## 2. File Modification Inventory

| File | Action | Responsibilities |
|:---|:---|:---|
| `cli/cmdagy/agy_instance_discovery.go` | **Create** | Discover active and registered Antigravity instances across system processes and user profiles; return `[]AgyInstanceInfo`. |
| `cli/cmdagy/agy_running_prompts_cmd.go` | **Modify** | Add `-i/--instance` and `-a/--all-instances` flags to `makeRunningPromptsLsCmd()`; delegate to multi-instance runner. |
| `cli/cmdagy/agy_prompts_collector.go` | **Modify** | Enhance prompt collection to accept target instance configuration; query instance-specific `brain/` transcripts and `conversation_summaries.db`. |
| `cli/cmdui/ui_types.go` | **Modify** | Declare `AgyInstanceInfo`, `AgyInstancePromptPayload`, and `AgyMultiInstancePromptResponse` REST models with positive boolean fields. |
| `cli/cmdui/ui_server.go` | **Modify** | Register `/api/instances` and `/api/prompts/instances` route handlers with query parameter parsing. |
| `cli/cmdagy/agy_instance_test.go` | **Create** | Unit tests verifying instance filtering, `--all-instances` aggregation, and JSON payload serialization. |

---

## 3. Step-by-Step Implementation Plan

### Step 3.1: Define Data Models (`cli/cmdui/ui_types.go`)
1. Implement `AgyInstanceInfo` with positive boolean fields:
   ```go
   type AgyInstanceInfo struct {
       InstanceID       string   `json:"instanceId"`
       InstanceName     string   `json:"instanceName"`
       ProcessID        int      `json:"processId"`
       LanguageServer   string   `json:"languageServer"`
       ConfigDir        string   `json:"configDir"`
       BrainDir         string   `json:"brainDir"`
       SummariesDBPath  string   `json:"summariesDbPath"`
       ActiveWorkspaces []string `json:"activeWorkspaces"`
       IsPrimary        bool     `json:"isPrimary"`
       IsRunning        bool     `json:"isRunning"`
       LastActiveAt     string   `json:"lastActiveAt"`
   }
   ```
2. Implement `AgyInstancePromptPayload` and `AgyMultiInstancePromptResponse` envelopes ensuring clean JSON marshaling.

### Step 3.2: Implement Instance Discovery Engine (`cli/cmdagy/agy_instance_discovery.go`)
1. Implement `DiscoverAllAgyInstances() ([]AgyInstanceInfo, error)`:
   - Identify primary instance from `~/.gemini/antigravity/`.
   - On Linux/macOS: scan active processes matching language server binaries.
   - On Windows: invoke CIM query via PowerShell (`Get-CimInstance Win32_Process`).
   - Read registered instances from `~/.gemini/antigravity/instances.json` if available.
2. Implement `ResolveInstance(instanceID string) (*AgyInstanceInfo, error)` with default fallback to primary instance.

### Step 3.3: Extend Running Prompts CLI (`cli/cmdagy/agy_running_prompts_cmd.go`)
1. In `setupRunningPromptsLsFlags(cmd *cobra.Command)`:
   ```go
   cmd.Flags().StringP("instance", "i", "", "Filter by Antigravity instance ID or alias (default: primary)")
   cmd.Flags().BoolP("all-instances", "a", false, "Aggregate prompts across all active Antigravity instances")
   ```
2. Update `executeRunningPromptsLsCmd(cmd *cobra.Command, args []string)`:
   - Read instance flags.
   - If `isAllInstances`: execute `RunMultiInstancePromptsLs(limit, wc, isFull, isJSON)`.
   - If `instanceID != ""`: resolve specific instance before querying.

### Step 3.4: Mount REST API Endpoints (`cli/cmdui/ui_server.go`)
1. In `mountAPIRoutes(mux *http.ServeMux)`:
   - Mount `/api/instances` -> `handleAPIInstances`
   - Mount `/api/prompts/instances` -> `handleAPIPromptsInstances`
2. Implement `handleAPIInstances`:
   - Enforce `GET` method only.
   - Return JSON list of detected instances.
3. Implement `handleAPIPromptsInstances`:
   - Parse query parameters: `instance`, `status`, `limit`, `max_words`, `include_convs`.
   - Query prompt snapshots for requested instance or all instances.
   - Return structured response envelope.

---

## 4. Verification Commands & Expected Output

```bash
# 1. Verify CLI output for primary instance in JSON format
gitmap rp ls --json

# 2. Verify all-instances query across detected environments
gitmap rp ls --all-instances --json

# 3. Test REST API endpoints via curl
curl -s http://127.0.0.1:42120/api/instances | jq .
curl -s "http://127.0.0.1:42120/api/prompts/instances?status=all&limit=5" | jq .

# 4. Run targeted unit test
go test -v -run TestInstanceDiscovery ./cli/cmdagy/
```

---

## 5. Acceptance Criteria

- [ ] `gitmap rp ls` supports `--instance` and `--all-instances` without error.
- [ ] Structured JSON output conforms to positive boolean attributes (`isRunning`, `isPrimary`, `isSuccess`).
- [ ] `/api/instances` returns HTTP 200 with all detected Antigravity instances.
- [ ] `/api/prompts/instances` partitions running and in-queue prompts correctly per instance.
- [ ] File-scoped unit tests pass with zero regressions.

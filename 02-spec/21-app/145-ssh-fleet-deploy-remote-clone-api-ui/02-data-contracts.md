# Spec 145: Data Contracts, Schemas & API Envelopes

Spec Reference: [01-overview.md](01-overview.md)

---

## 1. REST Endpoint Contracts

### 1.1 Cluster Node Daemon (`:49152`)
All inter-node REST communication utilizes mutual HMAC authentication via the `X-GitMap-Token` header.

| Method | Endpoint | Description | Request Body | Response Body |
|---|---|---|---|---|
| `GET` | `/api/v1/ping` | Liveness check | None | `{"status":"ok","time":123456789}` |
| `POST` | `/api/v1/exec` | Remote command execution | `{"command":"gitmap version","timeout":30}` | `{"exit_code":0,"stdout":"...","stderr":""}` |
| `POST` | `/api/v1/update` | Target package update | `{"target":"all","dry_run":false}` | `{"updated":["node","agy"],"errors":[]}` |
| `GET` | `/api/v1/inventory` | Installed software list | None | `{"packages":[{"name":"node","version":"22.0.0"}]}` |

### 1.2 Web UI Internal REST APIs (`127.0.0.1:8080`)
| Method | Endpoint | Description | Request Body | Response Body |
|---|---|---|---|---|
| `GET` | `/api/ssh/nodes` | List all cluster nodes | None | `[{"node_id":"worker-1","alias":"worker-1","host":"192.168.1.10","is_online":true}]` |
| `POST` | `/api/editor/read` | Fetch file content | `{"node_alias":"worker-1","file_path":"/etc/app.conf"}` | `{"is_success":true,"content":"...","language":"shell"}` |
| `POST` | `/api/editor/save` | Save file content | `{"node_alias":"worker-1","file_path":"/etc/app.conf","content":"..."}` | `{"success":true}` |
| `POST` | `/api/commitin/exec`| Execute visual commit | `{"message":"feat: add ui","amend":false}` | `{"success":true,"output":"[main abc1234]..."}` |
| `GET` | `/api/settings` | Get local UI settings | None | `{"theme":"dark","default_remote":"origin","cluster_port":49152}` |

---

## 2. Core Go Structs & Models

### 2.1 Remote File Content Contract
```go
type RemoteFileContent struct {
	NodeAlias string `json:"node_alias"`
	FilePath  string `json:"file_path"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
}
```

### 2.2 Software Inventory Matrix
```go
type FleetSoftwareInventory struct {
	NodeAlias string                 `json:"node_alias"`
	NodeHost  string                 `json:"node_host"`
	IsOnline  bool                   `json:"is_online"`
	Packages  map[string]ToolVersion `json:"packages"`
}

type ToolVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Path    string `json:"path"`
}
```

### 2.3 Self-Healing Clone Event Log
```go
type CloneAuthProbeEvent struct {
	TargetNode  string `json:"target_node"`
	RepoURL     string `json:"repo_url"`
	InitialFail string `json:"initial_fail"`
	ActionTaken string `json:"action_taken"`
	KeyFingerprint string `json:"key_fingerprint"`
	Resolved    bool   `json:"resolved"`
	DurationMs  int64  `json:"duration_ms"`
}
```

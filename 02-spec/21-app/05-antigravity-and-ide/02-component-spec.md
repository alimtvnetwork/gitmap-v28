# 05-antigravity-and-ide: Antigravity CLI Components, Prompt Manager & IDE Sync Specification

- **Spec ID:** `05-antigravity-and-ide/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Antigravity CLI Controller, Prompt Manager, Running Projects, Cursor Sync
- **Dependencies:** `cli/cmdagy`, `cli/cmdcursor`, `cli/cmdide`, `cli/workspacesync`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Antigravity and IDE cluster comprises four primary operational components:

```
05-antigravity-and-ide/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Antigravity CLI Controller** | `cli/cmdagy/agy.go`, `cli/cmdagy/deploy.go` | Antigravity management commands, settings propagation, and fleet sync. |
| **Running Projects Registry** | `cli/cmdagy/rp.go`, `cli/cmdagy/projects.go` | Discovery and indexing of active Antigravity projects across workspaces. |
| **Prompt Lifecycle Manager** | `cli/cmdagy/prompts.go`, `cli/cmdprompt/` | Prompt template management, running prompt backup/restore, decision logs. |
| **Multi-IDE Synchronizer** | `cli/cmdcursor/`, `cli/workspacesync/` | Settings and extension synchronization across Cursor, VSCode, and Desktop. |

---

## 2. Antigravity Deployment Engine (`gitmap agy deploy`)

### 2.1 Deployment Pipeline
The `gitmap agy deploy` command performs five sequential phases:
1. **Pre-flight Audit:** Reads local workstation theme seeds, plugin paths, and skill manifests.
2. **Payload Bundling:** Creates an in-memory gzipped tar archive of configuration files and plugin assets.
3. **Remote Node Handshake:** Connects via SSH using RSA-OAEP vault credentials.
4. **Remote Extraction & Path Normalization:** Unpacks payload into Linux target directories (`~/.config/Antigravity/User/settings.json`, `~/.gemini/config/`).
5. **Post-Deploy Hardening:** Sets SUID permissions on `chrome-sandbox` and restarts Antigravity daemon if running.

```go
package cmdagy

type DeployOptions struct {
    TargetNode    string `json:"targetNode"`
    SyncThemes    bool   `json:"syncThemes"`
    SyncPlugins   bool   `json:"syncPlugins"`
    SyncSkills    bool   `json:"syncSkills"`
    RestartDaemon bool   `json:"restartDaemon"`
}
```

---

## 3. Running Projects & Pinned Projects Registry

### 3.1 SQLite Project Registry Schema
Projects are cataloged in `installation.db` within the `AntigravityProject` table:

```sql
CREATE TABLE IF NOT EXISTS AntigravityProject (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    projectName TEXT NOT NULL UNIQUE,
    workspacePath TEXT NOT NULL,
    lastOpened DATETIME,
    isPinned INTEGER NOT NULL DEFAULT 0,
    turboMode INTEGER NOT NULL DEFAULT 1,
    descriptorJson TEXT
);
```

### 3.2 Operational CLI Commands
- `gitmap agy projects ls`: Lists all detected Antigravity projects and active sessions.
- `gitmap agy pin <project>`: Pins a project to the top of the selection menu.
- `gitmap agy prompts backup`: Exports all active conversation prompts to `~/.gitmap/prompts_backup/`.
- `gitmap agy prompts restore`: Restores saved prompt templates to the active Antigravity session.

---

## 4. Cursor IDE Workstation Integration

The Cursor integration subsystem manages:
- **Profile Synchronization:** Syncs Cursor `settings.json`, `keybindings.json`, and snippets.
- **Taskbar / Dock Pinning:** Generates Desktop entries with correct icon paths on Linux (GNOME dock).
- **Workspace State Migration:** Safely migrates workspace storage databases (`state.vscdb`) across host machines.

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isAgyDeployPipelineFunctional: true
  isProjectRegistryIndexed: true
  isCursorStateSyncVerified: true
  isPositiveBooleansUsed: true
```

- [x] Remote deployment completes in <10 seconds over SSH.
- [x] Project table accurately records workspace paths without duplicates.
- [x] Cursor settings match source environment after migration.

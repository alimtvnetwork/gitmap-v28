# Integration & Execution Specification: Antigravity IDE Automatic Workspace Attachment & Target Repo Cloning

> **Module:** `cli/cmdclone/`, `cli/workspacesync/`, `cli/cmdagy/`  
> **Parent Program:** `333-clone-and-fix-git-repositories`  
> **File:** `02-spec/21-app/333-clone-and-fix-git-repositories/03-antigravity-ide-integration-and-clone-execution-spec.md`  
> **Status:** Approved Specification  

---

## 1. Overview & Architectural Goals

When cloning repositories via GitMap (`gitmap clone`, `gitmap cfr`, `gitmap clonenext`), developer flow requires that cloned repositories are immediately and autonomously accessible inside the running developer environment.

This specification details:
1. **Automated Antigravity Workspace Registration:** Hooking every direct and manifest clone into Antigravity IDE workspace configuration files so projects are instantly registered.
2. **Active Session & Current Execution Integration:** Live workspace attachment into the active Antigravity session descriptor and SQLite conversation summaries so the current agent context has access without requiring a manual reload or session restart.
3. **Execution Plan for Target Repositories:** Concrete cloning, slug resolution, and directory binding for `yt-manage V1` and `go email sender`.
4. **Verification & Acceptance Criteria:** Strict evidence criteria verifying zero data loss, positive boolean conformity, and error sovereignty.

---

## 2. Auto-Registration in Antigravity IDE on Clone

### 2.1 Current State & Root Cause
In `cli/cmdclone/clone_direct.go`, the post-clone sequence is handled by `postCloneRegistration`:
```go
func postCloneRegistration(params DirectCloneParams, url, repoName, folderName, absPath string) {
	persistRecloneTransport(url)
	handleDirectCloneSuccess(params, repoName, absPath)
	upsertDirectClone(url, repoName, folderName, absPath)
	registerSingleDesktop(repoName, absPath)
	WriteShellHandoff(absPath)
}
```
While GitHub Desktop (`registerSingleDesktop`) and VS Code Project Manager (`syncSingleClonedRepoToVSCodePM`) are updated, `workspacesync.SyncAntigravity` was never called in `postCloneRegistration`. Consequently, cloned repositories remained unregistered in Antigravity IDE until manual intervention.

### 2.2 Target Implementation
Update `cli/cmdclone/clone_direct.go` to invoke `workspacesync.SyncAntigravity`:
```go
func postCloneRegistration(params DirectCloneParams, url, repoName, folderName, absPath string) {
	persistRecloneTransport(url)
	handleDirectCloneSuccess(params, repoName, absPath)
	upsertDirectClone(url, repoName, folderName, absPath)
	registerSingleDesktop(repoName, absPath)
	_ = workspacesync.SyncAntigravity(absPath, repoName)
	_ = cmdide.SyncAntigravityPinned(absPath, repoName)
	WriteShellHandoff(absPath)
}
```

### 2.3 Idempotency and Path Normalization
`workspacesync.SyncAntigravity(repoPath, repoName)` provides:
1. **Idempotent Checks:** Queries existing JSON configurations under `.gemini/config/projects/*.json` using normalized `file:///` URIs. If an entry already exists for the repository URI, it returns early with `true`.
2. **Standardized Resource Manifest:** If not present, creates a new UUID descriptor with:
   - `id`: newly generated UUID
   - `name`: `repoName`
   - `projectResources.resources[0].gitFolder.folderUri`: formatted `file:///` path
   - `projectResources.resources[0].gitFolder.defaultBranch`: `"main"`
   - `isWorkspaceOnly`: `false`
3. **Pinned Registry Synchronization:** Calls `cmdide.SyncAntigravityPinned` to register the repository in `.gemini/config/pinned_projects.json`.

---

## 3. Active Session & Current Execution Integration

To make cloned repositories immediately accessible in the **currently running execution environment**, GitMap must execute a three-tier registration protocol:

```
+-----------------------------------------------------------------------------------+
|                        Clone & Registration Pipeline                              |
+-----------------------------------------------------------------------------------+
                                        |
                   1. Direct Git Clone Succeeds
                                        v
       +-------------------------------------------------------------+
       |   Tier 1: Global Antigravity Config                         |
       |   - Write .gemini/config/projects/<uuid>.json               |
       |   - Append to .gemini/config/pinned_projects.json           |
       +-------------------------------------------------------------+
                                        v
       +-------------------------------------------------------------+
       |   Tier 2: Active Instance Multi-Root Resource Descriptor    |
       |   - Load instance project JSON (e.g. 1deb16d4-*.json)       |
       |   - Append cloned folder to projectResources.resources      |
       +-------------------------------------------------------------+
                                        v
       +-------------------------------------------------------------+
       |   Tier 3: Active Conversation Workspace Summaries           |
       |   - Inspect conversation_summaries.db                       |
       |   - Append workspace URI to workspace_uris column           |
       |   - Run `gitmap agy add <path>` CLI command                 |
       +-------------------------------------------------------------+
```

### 3.1 Tier 2: Instance Project Multi-Root Workspace Injection
When running inside an active Antigravity IDE instance (e.g., instance `v4-1464`), the active project descriptor defines the loaded workspace folders.
For instance:
- File: `<instanceDir>/home/.gemini/config/projects/1deb16d4-4e55-441d-9cda-4fe16f26e796.json`
- Action: Parse JSON, inspect `projectResources.resources` slice.
- If URI `file:///d%3A/work/<target-repo>` is not present, append:
  ```json
  {
    "gitFolder": {
      "folderUri": "file:///d%3A/work/<target-repo>",
      "defaultBranch": "main"
    }
  }
  ```
- Save updated JSON with indentation.

### 3.2 Tier 3: Conversation Summaries Database Update
The Antigravity runtime tracks active conversation workspaces in `conversation_summaries.db`:
- Table: `conversation_summaries`
- Columns: `conversation_id`, `workspace_uris` (comma-separated or JSON list of URIs), `project_id`.
- Action:
  1. Retrieve active conversation records where `killed IS NULL OR killed = 0`.
  2. If the cloned repo URI is missing from `workspace_uris`, append the URI formatted as `file:///<workspace-root>/<target-repo>`.
  3. Commit transaction via SQLite driver.

### 3.3 Tier 4: GitMap CLI Integration Command
Invoke the existing high-speed CLI tool:
```bash
gitmap agy add <path>
```
This performs directory validation, restricted path checks, project ID discovery, and queues the workspace without spawning blocking background jobs.

---

## 4. Execution Plan for Target Repositories

### 4.1 Target 1: `yt-manage V1`
1. **Repository Slug Resolution:**
   - Raw input: `yt-manage V1` / `yt-manage-v1` / `yt-manage`
   - Remote GitHub URL: `https://github.com/alimtvnetwork/yt-manage`
   - Clone Path: `d:/work/yt-manage`
   - Alias / Junction: `d:/work/yt-manage-v1` -> `d:/work/yt-manage` (if legacy scripts target the hyphenated suffix).
2. **Execution Steps:**
   - Run: `gitmap clone alimtvnetwork/yt-manage` or `gitmap clone yt-manage`
   - Verify working tree: `git -C d:/work/yt-manage status`
   - Verify branch tracking: `git -C d:/work/yt-manage branch -vv`

### 4.2 Target 2: `go email sender`
1. **Repository Slug Resolution:**
   - Raw input: `go email sender` / `go-email-sender`
   - Remote GitHub URL: `https://github.com/alimtvnetwork/go-email-reader`
   - Primary Path: `d:/work/go-email-reader`
   - Alias Path: `d:/work/go-email-sender`
2. **Execution Steps:**
   - Inspect existing `d:/work/go-email-reader` clone:
     - Remote check: `git -C d:/work/go-email-reader remote get-url origin`
     - Status check: ensure working tree is clean and synced with remote.
   - Link / Junction creation:
     - If `d:/work/go-email-sender` does not exist, create directory junction pointing to `d:/work/go-email-reader` so both directory conventions succeed unconditionally:
       ```powershell
       New-Item -ItemType Junction -Path "d:\work\go-email-sender" -Target "d:\work\go-email-reader"
       ```
   - Register both paths via `gitmap agy add d:\work\go-email-reader` and `gitmap agy add d:\work\go-email-sender`.

---

## 5. Acceptance Criteria & Verification Evidence

| ID | Criterion | Verification Method | Expected Outcome |
|---|---|---|---|
| AC-1 | Auto-registration hook in `clone_direct.go` | Code Inspection | `postCloneRegistration` includes `workspacesync.SyncAntigravity` and `cmdide.SyncAntigravityPinned`. |
| AC-2 | `yt-manage` cloned and valid | CLI Execution | `d:/work/yt-manage` exists with valid `.git` pointing to `alimtvnetwork/yt-manage`. |
| AC-3 | `go-email-reader` / `go-email-sender` available | CLI Execution | `d:/work/go-email-reader` and `d:/work/go-email-sender` exist and track origin cleanly. |
| AC-4 | Active session descriptor includes cloned folders | JSON Inspection | Project descriptor contains `gitFolder` entries for cloned repositories. |
| AC-5 | Zero test/build regressions | Pipeline / Linter | Zero regressions in `cli/cmdclone` packages. |

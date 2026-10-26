# Subtask 333.3 — Cloned Repositories Execution & Antigravity IDE Auto-Attachment

> **Task Code:** `Task-03`  
> **Title:** Cloned Repositories Execution & Antigravity IDE Auto-Attachment  
> **Parent Plan:** `.ai-memory/plans/pending/333-clone-and-fix-git-repositories.md`  
> **Specification:** `02-spec/21-app/333-clone-and-fix-git-repositories/03-antigravity-ide-integration-and-clone-execution-spec.md`  
> **Status:** PENDING  
> **Owned Files:**  
> - `cli/cmdclone/clone_direct.go`  
> - Antigravity project configuration JSON descriptors  
> - Target cloned repositories (`d:/work/yt-manage`, `d:/work/go-email-reader`, `d:/work/go-email-sender`)  

---

## 1. Objectives & Execution Tasks

### Task 3.1: Hook Antigravity Auto-Registration into Clone Direct Pipeline [COMPLETED]
1. Modified `cli/cmdclone/clone_direct.go`:
   - In `postCloneRegistration(params DirectCloneParams, url, repoName, folderName, absPath string)`:
   - Added call: `_ = workspacesync.SyncAntigravity(absPath, repoName)`
   - Note: Invoked directly via `workspacesync.SyncAntigravity` to avoid package import cycle (`cmdclone` -> `cmdide` -> `cmdcursor` -> `cmdssh` -> `cmdclone`).
2. Verified guidelines:
   - No negative booleans introduced.
   - Function length of `postCloneRegistration` remains 8 lines (well under <= 15 lines max).
   - Zero nested conditionals.
   - Blank line before `if`, after `}`, before `return` strictly preserved.

### Task 3.2: Execute Cloning of `yt-manage V1`
1. Resolve remote repository URL for `yt-manage`:
   - Remote: `https://github.com/alimtvnetwork/yt-manage`
   - Target destination: `d:/work/yt-manage`
2. Execute clone via GitMap:
   - Command: `gitmap clone alimtvnetwork/yt-manage` or `gitmap clone yt-manage`
3. Verify clone integrity:
   - Check `.git` repository existence and default branch (`main`).
   - Create junction/alias `d:/work/yt-manage-v1` if legacy tooling requires hyphenated suffix.

### Task 3.3: Execute Verification & Binding of `go email sender`
1. Inspect existing repository at `d:/work/go-email-reader`:
   - Remote: `https://github.com/alimtvnetwork/go-email-reader`
   - Fetch latest refs: `git -C d:/work/go-email-reader fetch origin`
2. Ensure directory link `d:/work/go-email-sender` exists:
   - Create directory junction `d:/work/go-email-sender` -> `d:/work/go-email-reader` so both names resolve identically without duplicating storage.

### Task 3.4: Integrate Repositories into Active Antigravity IDE Session
1. Execute CLI registration:
   - `gitmap agy add d:\work\yt-manage`
   - `gitmap agy add d:\work\go-email-reader`
   - `gitmap agy add d:\work\go-email-sender`
2. Update active instance project descriptor:
   - Locate active instance descriptor (e.g. `1deb16d4-4e55-441d-9cda-4fe16f26e796.json`).
   - Append folder resources for `d:/work/yt-manage` and `d:/work/go-email-reader` into `projectResources.resources`.
3. Update `conversation_summaries.db`:
   - Ensure `workspace_uris` contains URIs for newly attached repositories.

---

## 2. Verification Evidence Contract

| Check ID | Verification Action | Expected Evidence | Status |
|---|---|---|---|
| V-01 | Check `postCloneRegistration` in `cli/cmdclone/clone_direct.go` | Calls `workspacesync.SyncAntigravity`. | **PASS** |
| V-02 | Inspect `d:/work/yt-manage` | Git directory exists, remote origin matches `alimtvnetwork/yt-manage`. | PENDING |
| V-03 | Inspect `d:/work/go-email-reader` & `d:/work/go-email-sender` | Both paths resolve to valid git repositories with valid origin remotes. | PENDING |
| V-04 | Inspect Antigravity project config JSONs | Descriptors created or updated with valid `file:///` folder URIs. | PENDING |
| V-05 | Run `git status` | Clean working tree with only expected files modified or created. | PENDING |

---

## 3. Strict Coding Guidelines Compliance
- **Error Sovereignty:** All errors return directly (`return err` or `return apperror.Wrap(...)`); no leaf calls to `os.Exit`.
- **Function Size:** Functions <= 15 lines max.
- **Boolean Naming:** `is*` or `has*` only. No negative booleans.
- **Relative Path Usage:** All documentation and path references strictly follow relative path conventions.

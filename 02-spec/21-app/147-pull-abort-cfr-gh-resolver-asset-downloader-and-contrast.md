# Specification 147: Pull Batch Abort, CFR GitHub Resolver, Asset Downloader, Color Contrast, and Token Fleet Management

> **Spec Status:** Active  
> **Traceability IDs:** Task-01, Task-02, Task-03, Task-04, Task-05, Task-06, Task-07, Task-08, Task-09, Task-10  
> **Canonical Path:** `02-spec/21-app/147-pull-abort-cfr-gh-resolver-asset-downloader-and-contrast.md`  
> **Parent Spec:** `02-spec/21-app/01-index.md`

---

## 1. Domain Architecture & System Context

This specification defines critical resilience fixes, visual styling upgrades, asset pipelines, and repository resolution automation across GitMap CLI:

1. **Pull Batch Abort & Graceful Cancellation (`cmdpull`)**:
   - Prevent `E9000:EXECUTION` AppError stack trace dumps when users select `[q] Quit` or `[s] Skip` in dirty repository remediation.
   - Cleanly finalize batch tasks in SQLite store with non-zero exit code but zero fatal panic stack traces.

2. **Search in Untracked / Parent Directories (`cmd/search.go`)**:
   - When `gitmap search <filename>` is invoked from an untracked directory (e.g. `D:\work`), eliminate the fatal `search.getRepoDB` failure.
   - Gracefully fall back to global repository index or local workspace filesystem search to locate files (e.g. `vmpass.txt`).

3. **CFR & Clone Short-Name Resolution via GitHub CLI & SQLite `repodb` Cache**:
   - When `gitmap cfr <name>` or `gitmap clone <name>` receives a non-URL repository slug (e.g. `pwp-mobile`), automatically resolve the remote URL using `gh repo list` / `gh repo view` or cached SQLite repository metadata.
   - Persist repository names per user/org into SQLite `repodb` cache for high-speed terminal auto-completion.
   - Automatically register cloned repositories in both GitHub Desktop and VS Code Project Manager (`projects.json`).

4. **LightShot & PrintScreen Asset Downloader (`gitmap asset download-prnt` / `gitmap prnt`)**:
   - Native CLI command to download screenshot assets from `prnt.sc/<slug>` or `img.lightshot.app/<id>.png`.
   - Automatically parse HTML `og:image` or `#screenshot-image`, extract image bytes, and save to `assets/screenshots/<slug>-<timestamp>.png` and the active Antigravity brain workspace.

5. **Terminal & Installer High-Contrast Color Palette**:
   - Replace low-contrast dark blue accents (`#6272a4`, conhost 16-color dark blue downsampling) with high-contrast bright cyan (`#8be9fd` / `\033[1;96m`), bright white (`#f8f8f2` / `\033[1;97m`), or crisp yellow across `reconcile_prompt.go`, `install.ps1`, and documentation UI.

6. **Git Access Token Fleet Deployment (`gitmap token`)**:
   - Manage access tokens (`gitmap token add|remove|list|deploy`) and distribute them across remote SSH fleet nodes.

7. **Multi-Agent & E2E Prompt Synchronization**:
   - Refactor Folder 21 temporary E2E test prompt (`01-prompts/21-temp-end-to-end-tests/01-temp-end-to-end-test.md`) to establish the 2-half N-step lifecycle and strict must-follow security and isolation constraints.
   - Update Parent Task and Read Codebase prompts to mandate 2-agent orchestration with 3-4 parallel tasks during spec authoring.

---

## 2. Ingested Screenshot Visual Telemetry

The user provided four screenshot artifacts illustrating exact defects:

1. `assets/screenshots/8BnieUEKidCr.png`:
   - Command: `gitmap search vmpass.txt` executed in `D:\work>`.
   - Failure: Fatal crash `[E9000:EXECUTION] search.getRepoDB: current directory is not a tracked gitmap repository. run 'gitmap scan' first` with 8-frame Go stack trace.
   - Resolution: Catch untracked directory condition and fall back to global DB or local workspace search without stack traces.

2. `assets/screenshots/_N8xDMM-6ylG.png`:
   - View: `gitmap pull` interactive remediation menu showing `Path: <repo-root>` and remediation options in dark blue on black terminal background.
   - Failure: Text is virtually illegible; selecting `[q]` triggers batch execution failure.
   - Resolution: Upgrade `dimStyle` to high-contrast pastel white/cyan; handle `[q]` as clean user abort.

3. `assets/screenshots/0J1Th8lwMNIL.png`:
   - View: Scripts Fixer `.run agy help` and Antigravity installation status.
   - Observation: Verify clean high-contrast output and absence of dark blue blends.

4. `assets/screenshots/w66J-xV1NN3E.png`:
   - Command: `gitmap pull all` running batch lifecycle across 62 repositories.
   - Failure: Execution failure with 2 failures dumped as `E9000:EXECUTION` at `cmdpull/pull.go:498`.
   - Resolution: Clean failure reporting without stack dump.

---

## 3. Data Contracts & Interfaces

### 3.1 GitHub CLI Repository Resolution Contract
```go
type GhRepoEntry struct {
	Name          string `json:"name"`
	NameWithOwner string `json:"nameWithOwner"`
	URL           string `json:"url"`
}

type CachedRepo struct {
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	URL       string    `json:"url"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

### 3.2 Asset Downloader Request & Result
```go
type AssetDownloadOptions struct {
	URL        string
	TargetDir  string
	TargetName string
	SyncBrain  bool
}

type AssetDownloadResult struct {
	SourceURL string
	LocalPath string
	BrainPath string
	BytesSize int64
}
```

### 3.3 Git Token Fleet Deployment
```go
type TokenDeployRequest struct {
	Token    string
	TargetIP string
	Alias    string
	Protocol string
}
```

---

## 4. Verification Gates & Invariants

1. **Gate 1 (Search Fallback):** `gitmap search <filename>` in an untracked directory MUST NEVER output `[E9000:EXECUTION]` stack trace.
2. **Gate 2 (Pull Abort Grace):** Exiting interactive remediation with `q` MUST exit cleanly with code 1/0 and ZERO stack trace.
3. **Gate 3 (CFR Slug Resolution):** `gitmap cfr <name>` MUST attempt `gh` or cache resolution if `<name>` is not a valid URL.
4. **Gate 4 (Asset Downloader):** `gitmap asset download-prnt <url>` MUST resolve image and save to file system.
5. **Gate 5 (Color Contrast):** Zero usage of low-contrast `#6272a4` or 16-color dark blue in terminal interactive prompts.
6. **Gate 6 (Prompt Architecture):** Folder 21 prompt contains 2-half N-step lifecycle and parent task prompts mandate 2 agents.

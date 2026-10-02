# Architecture Specification: Chrome Profile Export/Import Subsystem, VMware Integration, Repo-Secrets Orchestration & Minor Version Bump

## 1. System Overview & Scope

This specification establishes the architectural blueprint and operational protocol for task `72-chrome-ext-test-and-vmware-spec` across GitMap. The objective of this subsystem is to provide end-to-end resilience, deterministic profile preservation, isolated guest testing, and rigorous secret-backed execution logging.

The architecture comprises four core pillars:
1. **Chrome Profile Preservation Subsystem**:
   - Comprehensive capture, serialization, sanitization, and restoration of Google Chrome browser profiles.
   - Tri-format export capabilities: JSON snapshot format (`chromeExport`), structured ZIP archive (`buildChromeProfileArchive`), and curated portable SQLite databases (`ChromeProfileSQLiteEntries`).
   - Secure preservation of authentication state: session cookies (`Network/Cookies`, `Cookies`), Web Data credentials, and Google OAuth tokens extracted from the SQLite `token_service` table with multi-layer reversible encryption (`ChromeTokenVault`).
   - Extension tracking through ID catalogs and `gitmap-pending-extensions.txt` restoration hints.
   - Bidirectional synchronization with Chrome's profile picker (`Local State`) and GitMap's SQLite tracking store (`installation.db` / `gitmap.db`).
2. **VMware Test Environment & Guest Isolation**:
   - Execution boundary between the host workstation and virtualized test nodes (e.g. VMware Workstation / ESXi guests).
   - Shared filesystem conduits (`vmware-shared`), automated testing harness execution, and isolated test profile provisioning without host session interference.
3. **Repo-Secrets Sequencing & Operational Storage**:
   - Dynamic allocation of isolated sequence folders under `<repo-secrets-root>/01-gitmap/<seq>-<task-name>/` (e.g., `../repo-secrets/01-gitmap/10-chrome-ext-test-and-vmware/`).
   - Strict two-digit sequence numbering, lowercase folder conventions, and zero underscores.
   - Segregation of sensitive runtime artifacts (`logs/`, `exports/`, `reports/`, `artifacts/`) outside committed git repositories.
4. **Structured Error Logging & Diagnostic Tracing**:
   - Deterministic execution traces capturing timestamped phase transitions, CLI invocations, stdout/stderr streams, exit codes, and failure classification.
   - Non-fatal error handling for database synchronization with explicit warning surfacing.
5. **Initial Minor Version Bump Protocol**:
   - Canonical version bump from `v6.458.0` to `v6.459.0` anchored in `version.json`.
   - Synchronized updates across `package.json`, Go constants (`cli/constants/constants.go`), and changelog ledgers.

---

## 2. Architecture & Data Flow

```
                      +---------------------------------------+
                      |         GitMap CLI Command Flow       |
                      |  [ gitmap cpe ]       [ gitmap cpi ]  |
                      +---------------------------------------+
                                          |
                +-------------------------+-------------------------+
                |                                                   |
      [ Profile Export Phase ]                            [ Profile Import Phase ]
                |                                                   |
   +---------------------------+                        +---------------------------+
   |   Source Profile Scanner  |                        |  Export Archive Inspector |
   |  - Check Chrome Running   |                        |  - Inspect JSON / ZIP     |
   |  - Read Preferences / GAIA|                        |  - Version Check Warning  |
   +---------------------------+                        +---------------------------+
                |                                                   |
   +---------------------------+                        +---------------------------+
   |   Auth & Token Extraction |                        | Base Files & Auth Restore |
   |  - SQLite token_service   |                        |  - Write Bookmarks / Prefs|
   |  - Caesar & Double Base64 |                        |  - Restore Web Data       |
   |  - Base64 Session Cookies |                        |  - Re-inject Token Service|
   +---------------------------+                        |  - Restore Network Cookies|
                |                                       +---------------------------+
   +---------------------------+                                    |
   |    Packaging & Storage    |                        +---------------------------+
   |  - JSON Snapshot Format   |                        |  System Registration      |
   |  - ZIP Archive Bundling   |                        |  - Extensions Hint List   |
   |  - Curated SQLite Files   |                        |  - Chrome Local State Reg |
   |  - Write to repo-secrets  |                        |  - Upsert GitMap SQLite DB|
   +---------------------------+                        +---------------------------+
                |                                                   |
                +-------------------------+-------------------------+
                                          |
                      +---------------------------------------+
                      |    Repo-Secrets Logging & Telemetry   |
                      |  - Execution Log (execution.log)      |
                      |  - Audit Trace (audit-trace.json)     |
                      |  - Diagnostic Summary Report          |
                      +---------------------------------------+
```

---

## 3. Detailed Component Architecture

### 3.1 Chrome Profile Export/Import Subsystem (`cli/cmdchromeprofile`)

The Chrome profile management subsystem provides robust tooling for cloning, backing up, and migrating browser profiles without manual user intervention.

#### 3.1.1 JSON Snapshot Schema (`chromeExport`)

The JSON snapshot representation serves as a portable, text-serializable profile capture format. The data schema is defined as follows:

```go
type chromeExport struct {
    SchemaVersion    int               `json:"schemaVersion" yaml:"schemaVersion"`
    GitMapVersion    string            `json:"gitmapVersion,omitempty" yaml:"gitmapVersion,omitempty"`
    Name             string            `json:"name" yaml:"name"`
    DisplayName      string            `json:"displayName,omitempty" yaml:"displayName,omitempty"`
    Email            string            `json:"email,omitempty" yaml:"email,omitempty"`
    GaiaID           string            `json:"gaiaId,omitempty" yaml:"gaiaId,omitempty"`
    GaiaName         string            `json:"gaiaName,omitempty" yaml:"gaiaName,omitempty"`
    GaiaGivenName    string            `json:"gaiaGivenName,omitempty" yaml:"gaiaGivenName,omitempty"`
    ExportedAt       string            `json:"exportedAt" yaml:"exportedAt"`
    Bookmarks        json.RawMessage   `json:"bookmarks,omitempty" yaml:"bookmarks,omitempty"`
    Preferences      json.RawMessage   `json:"preferences,omitempty" yaml:"preferences,omitempty"`
    CookiesRawBase64 string            `json:"cookiesRawBase64,omitempty" yaml:"cookiesRawBase64,omitempty"`
    WebDataRawBase64 string            `json:"webDataRawBase64,omitempty" yaml:"webDataRawBase64,omitempty"`
    ExtensionIDs     []string          `json:"extensionIds,omitempty" yaml:"extensionIds,omitempty"`
    TokenVault       *ChromeTokenVault `json:"tokenVault,omitempty" yaml:"tokenVault,omitempty"`
}
```

Key schema rules:
- **Additive Compatibility**: Schema version is currently `1`. All newly added fields must remain optional and default to zero/nil to preserve backward compatibility with previous snapshot versions.
- **Version Notice**: On import, if `exp.GitMapVersion != constants.Version`, GitMap logs a non-fatal warning alerting the operator to the version difference while continuing import execution.
- **Bookmarks & Preferences**: Serialized as raw JSON messages (`json.RawMessage`) to prevent schema distortion or key loss during round-trip unmarshaling.

#### 3.1.2 Curated Portable SQLite Databases (`constants.ChromeProfileSQLiteEntries`)

When exporting with `--format=zip` or `--format=sqlite`, GitMap extracts a curated collection of SQLite databases from the profile directory:
- `History`: Navigation history, visit counts, search terms, and download records.
- `Web Data`: Autofill profiles, credit card tokens (if allowed), and Google token service entries.
- `Shortcuts`: Omnibox search provider mappings and shortcut indices.
- `Top Sites`: Most visited thumbnail and page URLs.
- `Favicons`: Cached website icons and icon mappings.
- `Network Action Predictor`: URL prediction and pre-rendering caches.
- `Network/Cookies` & `Cookies`: Session and persistent HTTP cookies across visited domains.

#### 3.1.3 ZIP Archive Layout & Packaging (`writeChromeExportZIP`)

The generated ZIP archive packages the snapshot JSON alongside raw database files:
```text
<profile-name>.zip
├── manifest.json            # Profile metadata and file inventory
├── <profile-name>.json      # Complete chromeExport JSON snapshot
├── History                  # Raw SQLite database
├── Web Data                 # Raw SQLite database
├── Shortcuts                # Raw SQLite database
├── Top Sites                # Raw SQLite database
├── Favicons                 # Raw SQLite database
├── Network Action Predictor # Raw SQLite database
└── Network/
    └── Cookies              # Raw SQLite database
```

#### 3.1.4 Auth Cookies & Token Vault (`ChromeTokenVault`)

To allow seamless test migration without forcing manual credential re-entry:
1. **Cookie Preservation**:
   - Inspects both `<profile>/Network/Cookies` and `<profile>/Cookies`.
   - Reads the raw database binary, encodes it to standard Base64 (`readProfileCookiesBase64`), and embeds it into `CookiesRawBase64`.
   - During import, decodes the binary and writes it to `<dstProfile>/Network/Cookies` with directories created via `os.MkdirAll`.
2. **OAuth Refresh Token Vault**:
   - Connects to `<profile>/Web Data` SQLite database using a read-only URI connection (`file:... ?mode=ro`).
   - Queries `SELECT service, encrypted_token FROM token_service`.
   - Implements multi-tier reversible encryption in `ChromeRefreshTokenEntry`:
     - Layer 1: Double Base64 encoding (`EncodeDoubleBase64`).
     - Layer 2: Caesar cipher transposition (`defaultCaesarShift = 13`).
     - Layer 3: Byte-shift rotation (`defaultByteShift = 7`).
   - Stores entries in `ChromeTokenVault` alongside step-by-step restoration directives in `RevertGuide`.
   - Upon import, reverses encryption passes and reconstitutes the `token_service` table inside the destination profile's `Web Data` database.

#### 3.1.5 Extension Management & Preferences Restoration

Chrome extension installations are tied to machine-specific file paths. Direct file copying of unpackaged extension directories can trigger Chrome's corruption detection. To solve this:
- GitMap inventories all extension directories inside `<profile>/Extensions` via `listExtensionIDs`.
- During export, the 32-character extension IDs (e.g. `cjpalhdlnbpafiamejdnhcphjbkeiagm`) are stored in `ExtensionIDs`.
- During import, GitMap writes `gitmap-pending-extensions.txt` into the destination profile root containing deduplicated extension IDs.
- On first launch, Chrome reads these entries or the operator triggers automated batch extension provisioning via `gitmap chrome extension-install`.

#### 3.1.6 GitMap SQLite Store Tracking (`store.DB`)

All export and import operations are indexed in GitMap's SQLite store (`installation.db` / `gitmap.db`):
- `ChromeProfile`: Master table tracking profile names, source paths, offline status (`IsOffline = 1`), and modification timestamps.
- `ChromeProfileExport`: Child table recording artifact format (`json`, `zip`, `sqlite`, `csv`), destination path, file byte size, and creation timestamp.
- **Fault-Tolerant Persistence**: Failures in SQLite persistence are logged to `os.Stderr` via `constants.MsgChromeProfileDBWarn` without aborting file system writes.

---

### 3.2 VMware Test Environment & Guest Isolation Architecture

Testing browser profile migrations in automated CI/CD pipelines or developer machines requires virtualization boundaries to prevent corrupting local developer sessions.

#### 3.2.1 Host-Guest Isolation Boundary

```
+-------------------------------------------------------------+
|                     Host Workstation (L0)                   |
|  - GitMap CLI Controller                                    |
|  - Repo-Secrets Master Storage                              |
|  - VMware Workstation / Player / ESXi Hypervisor            |
+-------------------------------------------------------------+
                              |
                     [ Shared Folder / SMB ]
                              |
+-------------------------------------------------------------+
|                     VMware Guest VM (L1)                    |
|  - Isolated Windows / Linux Guest OS                        |
|  - Clean Chrome Profile Sandbox: %LOCALAPPDATA%\Google\...   |
|  - GitMap Guest Runner Binary                               |
|  - Ephemeral Profile Fixture Testing                        |
+-------------------------------------------------------------+
```

1. **Process Lock Prevention**:
   - Chrome locks SQLite files with `LOCK` files and exclusive OS locks when running.
   - The test runner queries running processes via `tasklist /FI "IMAGENAME eq chrome.exe"` (Windows) or `pgrep -x chrome` (Linux/macOS).
   - If Chrome processes are active, the runner terminates them gracefully before attempting profile copying or restoration.
2. **Shared Conduit Execution**:
   - The test bundle is transferred to the VMware guest via VMware Shared Folders (`vmware-shared`) or SSH delegation (`gitmap ssh <alias>`).
   - Subtasks verify export creation, profile directory deletion, and subsequent re-import inside the clean guest sandbox.

---

### 3.3 Repo-Secrets Sequencing & Operational Storage

`repo-secrets` serves as the centralized repository for operational configuration, test artifacts, and sensitive machine parameters. To maintain consistency with repository standards:

#### 3.3.1 Sequencing & Directory Structure

- **Location**: Discovered dynamically relative to workspace root (`../repo-secrets/01-gitmap/`) or via environment variable. Committed markdown files must NEVER hardcode absolute paths (such as drive letters or machine home paths).
- **Sequence Assignment**: Must follow a two-digit prefix sequence (`01-`, `02-`, ..., `10-`, etc.), strictly lowercase, with zero underscores and zero hyphens in prefixes.
- **Dedicated Task Folder**: For task 72, the execution workspace is organized under:
  `<repo-secrets-root>/01-gitmap/10-chrome-ext-test-and-vmware/` (or designated next available sequence number).

```text
<repo-secrets-root>/01-gitmap/<seq>-<task-name>/
├── logs/
│   ├── execution.log            # Streamed CLI output and timestamps
│   ├── audit-trace.json         # Structured machine-readable trace
│   └── error.log                # Captured failures and warnings
├── exports/
│   ├── sample-profile.json      # Generated JSON snapshot
│   ├── sample-profile.zip       # Generated ZIP archive
│   └── sample-profile.sqlite    # Exported SQLite databases
├── reports/
│   ├── verification-report.md   # Markdown summary of test execution
│   └── file-diff-matrix.json    # Pre/post restoration file parity matrix
└── artifacts/
    ├── bookmarks-dump.json      # Extracted bookmarks verification
    └── pending-extensions.txt   # Discovered extension IDs
```

#### 3.3.2 Hygiene & Exclusion Rules

- Committed repository code must not contain secret tokens or private keys.
- Secret files in `repo-secrets` are ignored in source repositories.
- All JSON envelopes inside `repo-secrets` adhere to standard schemas with variable interpolation (`${workDir}`, `${keyPath}`).

---

### 3.4 Error Logging & Diagnostic Tracing Engine

Testing profile operations demands exhaustive traceability to guarantee data integrity across export and re-import cycles.

#### 3.4.1 Execution Log Envelope Format (`audit-trace.json`)

The logging engine writes structured execution records:

```json
{
  "taskId": "72-chrome-ext-test-and-vmware-spec",
  "stepNumber": 3,
  "stepName": "chrome-export-test",
  "timestamp": "2026-10-02T12:30:00Z",
  "command": "gitmap cpe sample-profile --format=json,zip",
  "isDryRun": false,
  "hasSucceeded": true,
  "exitCode": 0,
  "durationMs": 420,
  "capturedArtifacts": [
    {
      "format": "json",
      "path": "exports/sample-profile.json",
      "byteSize": 14520,
      "checksumSha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    }
  ],
  "diagnostics": {
    "chromeRunning": false,
    "lockFilesEncountered": 0,
    "tokensExtracted": 2,
    "cookiesExtracted": 1
  },
  "errors": []
}
```

#### 3.4.2 Error Classification & Logging Invariants

- **Fatal Errors**: File I/O permission denied, missing profile source directory, corrupt zip payload. The runner terminates immediately, logs the stack trace to `error.log`, and triggers cleanup.
- **Recoverable Warnings**: Transient Chrome lock files (`LOCK`), SQLite metadata write failure (when file export succeeds), version mismatch notice. These are logged with level `WARN` and allow execution to proceed.
- **Timestamp Integrity**: All timestamps must be formatted according to ISO-8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`).

---

### 3.5 Initial Minor Version Bump Protocol

The project mandates that every major spec implementation wave begins with a minor version bump to ensure clear provenance for generated test artifacts and CLI builds.

#### 3.5.1 Single Source of Truth (`version.json`)

`version.json` at repository root controls all version resolution across the repository:
- **Baseline Version**: `6.458.0`
- **Target Version**: `6.459.0`
- **Cadence**: Minor bump increments the feature track identifier in accordance with the project's versioning conventions.

#### 3.5.2 Atomic Synchronization Targets

When executing the minor bump, the following targets must be updated synchronously:
1. `version.json`:
   - Update `"Version": "6.459.0"`
   - Update `"releaseDate": "<current-date>"`
2. `package.json`:
   - Update `"version": "6.459.0"`
3. `cli/constants/constants.go`:
   - Update `var Version = "6.459.0"`
4. Changelog ledger entries:
   - Record version bump entry under `.ai-memory/plans/ledger.md` and spec changelogs.

---

## 4. Verification & Quality Gates

To validate compliance with this architectural specification:
1. **Compilation Gate**:
   - `go build ./...` must compile cleanly across all packages without syntax errors or broken references.
2. **Relative Path Enforcement**:
   - `python linter-scripts/check-relative-paths.py` must report 0 violations across all committed markdown files.
3. **Coding Guideline Conformance**:
   - `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdchromeprofile cli/constants cli/store --check-only` must pass cleanly.
4. **Boolean Naming Invariant**:
   - All boolean variables, fields, and flags must use positive prefixes (`is...`, `has...`, `should...`). Negative terms (`disable`, `not`, `no`) are forbidden.

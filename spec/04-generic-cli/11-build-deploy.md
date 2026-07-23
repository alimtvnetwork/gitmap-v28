# Build & Deploy

> **Related specs:**
> - [02-project-structure.md](02-project-structure.md) — package layout that the build targets
> - [12-testing.md](12-testing.md) — tests that run as part of the build pipeline
> - [13-checklist.md](13-checklist.md) — implementation phases including build setup

## Build Script Architecture

A single script at the repo root handles the full lifecycle:
**pull → build → deploy → (optional) run**.

### Step-Based Execution

| Step | Action | Skippable |
|------|--------|-----------|
| 1/4 | Git pull latest source | `-NoPull` |
| 2/4 | Resolve dependencies | No |
| 3/4 | Compile binary | No |
| 4/4 | Deploy to target directory | `-NoDeploy` |

### Configuration

Store build/deploy settings in a JSON file:

```json
{
  "deployPath": "/usr/local/bin",
  "buildOutput": "./bin",
  "binaryName": "toolname",
  "copyData": true
}
```

### Build with Embedded Variables

Use linker flags to embed values at compile time:

```
go build -ldflags "-X 'pkg/constants.RepoPath=$path'" -o binary .
```

### Version Verification

After building, run the binary with `version` to confirm:

```
$ ./bin/toolname version
v1.2.0
```

## Build Once, Package Once

Binaries are compiled **exactly once** per pipeline run. All downstream
steps — compression, checksumming, installer generation, and publishing —
operate on the already-built artifacts and **must never trigger a rebuild**.

## Deploy Patterns

### Nested Deploy Structure

```
deploy-target/
└── toolname/
    ├── toolname (binary)
    └── data/
        └── config.json
```

The subfolder should be on the system `PATH`.

### Retry-on-Lock (Windows)

Wrap file copy in a retry loop for locked binaries:

```
maxAttempts = 20
for each attempt:
    try: copy(source, destination); break
    catch locked: wait 500ms, retry
if failed: restore backup
```

## Self-Update Mechanism

### The Problem (Windows)

Running `.exe` files hold a file lock and cannot be overwritten.

### Solution: Copy-and-Handoff

1. **Parent** copies itself to a temp location
2. **Parent** launches the temp copy with a worker command (blocking)
3. **Worker** pulls, builds, and deploys the new binary
4. **Worker** uses rename-first strategy for locked binaries
5. **Worker** compares old vs new version
6. **Worker** runs cleanup to remove temp artifacts

### Critical Rules

- Parent MUST use blocking execution (not async)
- PATH sync MUST use rename-first on Windows
- Generated scripts MUST NOT contain interactive prompts
- Always provide rollback (keep `.old` backup until cleanup)

### Error Handling

| Scenario | Behavior |
|----------|----------|
| No repo path configured | Print error, exit 1 |
| Already up to date | Print message, exit 0 |
| Build fails | Backup remains, exit with error |
| Deploy locked after retries | Restore backup, fail clearly |

## Semantic Logging

Use color-coded logging functions:

| Function | Color | Use Case |
|----------|-------|----------|
| `Write-Step` | Magenta | Step headers `[1/4]` |
| `Write-Success` | Green | Successful operations |
| `Write-Info` | Cyan | Informational messages |
| `Write-Warn` | Yellow | Non-fatal warnings |
| `Write-Fail` | Red | Errors before exit |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)

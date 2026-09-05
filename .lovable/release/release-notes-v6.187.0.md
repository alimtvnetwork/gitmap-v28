## Quick Install v6.187.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.187.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.187.0/install.sh | bash
```

## Changelog v6.187.0

- Introduced pkg/fileutil with type-safe FileActionType and FileModeType enums, FileWrapper encapsulation, and package-level I/O convenience helpers
- Enhanced pkg/appfault with first-class errorId field on AppError and dedicated WrapFailure, WrapWriterFailure, and WrapReaderFailure constructors
- Overhauled 03-ai-scripts/06-cicd-local-runner.py with adaptive parallel worker pool, 3-batch IO throttling, and quiet success tick output
- Fixed SQLite concurrency race condition in gitmap/store/store.go by prioritizing busy_timeout pragma before journal_mode WAL pragma
- Introduced reusable run_worker_pool and add_worker_cli_arguments in 03-ai-scripts/02-shared-engine.py for repository-wide parallel script modernization
- Modernized 03-ai-scripts/16-installer-smoke-tester.py and 28-go-preflight-ci.py to leverage the shared worker pool with quiet tick output and JSON/file exports
- Bumped version to v6.187.0 across all Single Source of Truth manifests

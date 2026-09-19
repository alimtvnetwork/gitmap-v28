## Quick Install v6.263.0

### Windows (PowerShell 5.1+)
```powershell
irm https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.263.0/install.ps1 | iex
```

### Linux / macOS (Bash)
```bash
curl -fsSL https://github.com/alimtvnetwork/gitmap-v28/releases/download/v6.263.0/install.sh | bash
```

## Changelog v6.263.0

- Fix function identifier collisions in cli/cmdautomation: rename collectSearchFiles to collectSpecMigrateFiles in spec_migrate.go and resolveThreshold to resolveSlowThreshold in test_inventory.go
- Implement ListRuntimes, RefreshRuntimes, and reprobeRuntime in cli/cmdautomation/runtime_probe.go and update runtimes_cmd.go to use RuntimeRecord
- Fix RunRelPathsAudit reference to RunRelPathAudit in cli/cmdautomation/preflight.go
- Fix run_cmd.go type mismatch on RunWorkerPool(opts), map opts.Script = args[1], add IsJson to WorkerRunOptions, and render worker results
- Rename cli/cmdautomation/smoke_test.go to cli/cmdautomation/smoke_runner.go to prevent Go compiler build exclusion
- Align BuildFileContext, EncodeToStream, and DecodeFromStream in stream_encoding.go and worker_pool.go with result.Result[T] wrappers and Spec 124
- Fix resolveSpecDir in cli/cmdautomation/spec_migrate.go to inspect 02-spec/21-app and 02-spec in base directory

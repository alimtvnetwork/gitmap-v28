# Subtask 05: Service Driver Resolution Decoupling & Hermetic Testing

## Objective
Decouple `cli/cmdservice` from host OS service manager binaries (`sc.exe`, `systemctl`, `launchctl`) by introducing an injectable `DefaultServiceDriverResolver ServiceDriverResolver` and providing 100% hermetic mock unit tests.

## Assigned Files
- `cli/cmdservice/driver.go`
- `cli/cmdservice/service_cmd.go`
- `cli/cmdservice/service_cmd_test.go`

## Implementation Steps
1. In `cli/cmdservice/driver.go`:
   - Introduce `type ServiceDriverResolver func() ServiceDriver`.
   - Define `var DefaultServiceDriverResolver ServiceDriverResolver = ResolveServiceDriver`.
   - Update `EnsureServiceDriver()` to route through `DefaultServiceDriverResolver()`.
2. In `cli/cmdservice/service_cmd.go`:
   - Replace all 8 direct calls to `ResolveServiceDriver()` with `DefaultServiceDriverResolver()`.
   - Decompose `runServiceCreate` (extract `resolveServiceCreateParams`), `runServiceStatus` (extract `renderServiceStatus`), and `runServiceExport` (extract `buildServiceExportSchema`) to ensure all functions <= 8-15 lines.
3. In `cli/cmdservice/service_cmd_test.go`:
   - Provide `mockDriver` implementing `ServiceDriver`.
   - Implement unit tests covering `ls`, `ls --json`, `status`, `start`, `stop`, `create`, `rm`, `export`, `export --yaml`, `import`, `--help`, and validation error cases.
